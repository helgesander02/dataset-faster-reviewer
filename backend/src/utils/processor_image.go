package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

const (
	defaultMaxSlots      = 10
	taskCancellationWait = 3 * time.Second
	maxImageWidth        = 400
	imageQuality         = 75
	pageRangeThreshold   = 10
)

var maxWorkers = runtime.NumCPU()

type TaskManager struct {
	mu              sync.RWMutex
	runningTasks    []*TaskInfo
	maxSlots        int
	taskCounter     atomic.Int64
	currentPageHint atomic.Int64
}

type TaskInfo struct {
	ID        string
	PageIndex int
	StartTime time.Time
	Cancel    context.CancelFunc
	Done      chan struct{}
	closed    atomic.Bool
}

var (
	globalTaskManager *TaskManager
	once              sync.Once
)

func initTaskManager() {
	once.Do(func() {
		globalTaskManager = &TaskManager{
			runningTasks: make([]*TaskInfo, 0, defaultMaxSlots),
			maxSlots:     defaultMaxSlots,
		}
		log.Printf("Initialized task manager with %d slots", globalTaskManager.maxSlots)
	})
}

func (tm *TaskManager) addTask(taskID string, pageIndex int, cancel context.CancelFunc) *TaskInfo {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Update current page hint
	tm.currentPageHint.Store(int64(pageIndex))

	// Cancel tasks that are outside the active page range
	tm.cancelOutOfRangeTasks(pageIndex)

	// If still at capacity, cancel oldest task
	if len(tm.runningTasks) >= tm.maxSlots {
		tm.cancelOldestTask()
	}

	newTask := tm.createNewTask(taskID, pageIndex, cancel)
	tm.runningTasks = append(tm.runningTasks, newTask)

	return newTask
}

func (tm *TaskManager) cancelOutOfRangeTasks(currentPage int) {
	// Calculate active page range
	minPage := currentPage - pageRangeThreshold
	maxPage := currentPage + pageRangeThreshold

	tasksToCancel := make([]*TaskInfo, 0)
	remainingTasks := make([]*TaskInfo, 0)

	// Identify tasks outside the range
	for _, task := range tm.runningTasks {
		if task.PageIndex < minPage || task.PageIndex > maxPage {
			tasksToCancel = append(tasksToCancel, task)
		} else {
			remainingTasks = append(remainingTasks, task)
		}
	}

	// Cancel out-of-range tasks
	if len(tasksToCancel) > 0 {
		log.Printf("[PageTracker] Current page: %d, cancelling %d tasks outside range [%d, %d]",
			currentPage, len(tasksToCancel), minPage, maxPage)

		for _, task := range tasksToCancel {
			log.Printf("[PageTracker] Cancelling task %s for page %d (age: %v)",
				task.ID, task.PageIndex, time.Since(task.StartTime))
			task.Cancel()
			tm.waitForTaskCancellation(task)
		}

		tm.runningTasks = remainingTasks
	}
}

func (tm *TaskManager) cancelOldestTask() {
	oldestTask := tm.runningTasks[0]

	if log.Writer() != nil {
		log.Printf("Slots full (%d/%d), cancelling oldest task: %s (page: %d, running for %v)",
			len(tm.runningTasks), tm.maxSlots, oldestTask.ID, oldestTask.PageIndex, time.Since(oldestTask.StartTime))
	}

	oldestTask.Cancel()
	tm.runningTasks = tm.runningTasks[1:]
	tm.waitForTaskCancellation(oldestTask)
}

func (tm *TaskManager) waitForTaskCancellation(task *TaskInfo) {
	go func(t *TaskInfo) {
		select {
		case <-t.Done:
			return
		case <-time.After(taskCancellationWait):
			if log.Writer() != nil {
				log.Printf("Old task cancellation timeout for: %s", t.ID)
			}
		}
	}(task)
}

func (tm *TaskManager) createNewTask(taskID string, pageIndex int, cancel context.CancelFunc) *TaskInfo {
	return &TaskInfo{
		ID:        taskID,
		PageIndex: pageIndex,
		StartTime: time.Now(),
		Cancel:    cancel,
		Done:      make(chan struct{}),
	}
}

func (tm *TaskManager) removeTask(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for i, task := range tm.runningTasks {
		if task.ID == taskID {
			tm.closeTask(task)
			tm.removeTaskAtIndex(i)
			break
		}
	}
}

func (tm *TaskManager) closeTask(task *TaskInfo) {
	if task.closed.CompareAndSwap(false, true) {
		close(task.Done)
	}
}

func (tm *TaskManager) removeTaskAtIndex(index int) {
	tm.runningTasks = append(tm.runningTasks[:index], tm.runningTasks[index+1:]...)
}

func (tm *TaskManager) getTaskStatus() (int, int, []string) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	current := len(tm.runningTasks)
	taskIDs := tm.formatTaskStatuses()

	return current, tm.maxSlots, taskIDs
}

func (tm *TaskManager) formatTaskStatuses() []string {
	taskIDs := make([]string, 0, len(tm.runningTasks))
	for i, task := range tm.runningTasks {
		status := fmt.Sprintf("Slot[%d]: %s (page: %d, running: %v)", i, task.ID, task.PageIndex, time.Since(task.StartTime))
		taskIDs = append(taskIDs, status)
	}
	return taskIDs
}

func CompressImageSetToBase64(imagePaths []string, pageIndex int) []string {
	initTaskManager()

	taskID := generateTaskID(pageIndex)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = globalTaskManager.addTask(taskID, pageIndex, cancel)
	defer globalTaskManager.removeTask(taskID)

	result := processImagesWithContext(ctx, taskID, imagePaths)

	if isContextCancelled(ctx) {
		log.Printf("[PageTracker] Task %s for page %d was cancelled", taskID, pageIndex)
		return make([]string, len(imagePaths))
	}

	return result
}

func generateTaskID(pageIndex int) string {
	counter := globalTaskManager.taskCounter.Add(1)
	return fmt.Sprintf("task-page%d-%d-%d", pageIndex, time.Now().Unix(), counter)
}

func isContextCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func processImagesWithContext(ctx context.Context, taskID string, imagePaths []string) []string {
	base64Images := make([]string, len(imagePaths))
	sem := make(chan struct{}, maxWorkers)
	defer close(sem)

	var wg sync.WaitGroup

	// Track if we should stop launching new goroutines
	shouldStop := &atomic.Bool{}
	shouldStop.Store(false)

	// Monitor context in background
	go func() {
		<-ctx.Done()
		shouldStop.Store(true)
	}()

	for i, imagePath := range imagePaths {
		// Check both context and stop flag
		if isContextCancelled(ctx) || shouldStop.Load() {
			log.Printf("[ProcessImages] Task %s cancelled, stopping after launching %d/%d workers", taskID, i, len(imagePaths))
			wg.Wait()
			return base64Images
		}

		wg.Add(1)
		go processImage(ctx, &wg, sem, base64Images, i, imagePath)
	}

	return waitForProcessingCompletion(ctx, &wg, base64Images)
}

func processImage(ctx context.Context, wg *sync.WaitGroup, sem chan struct{}, base64Images []string, index int, path string) {
	defer wg.Done()

	// Early exit if context is already cancelled
	if isContextCancelled(ctx) {
		return
	}

	if !acquireSemaphore(ctx, sem) {
		return
	}
	defer releaseSemaphore(sem)

	// Check again before expensive operation
	if isContextCancelled(ctx) {
		return
	}

	base64Image := compressImageSafely(ctx, path)
	base64Images[index] = base64Image
}

func acquireSemaphore(ctx context.Context, sem chan struct{}) bool {
	select {
	case <-ctx.Done():
		return false
	case sem <- struct{}{}:
		return true
	}
}

func releaseSemaphore(sem chan struct{}) {
	<-sem
}

func compressImageSafely(ctx context.Context, path string) string {
	base64Image, err := CompressImageToBase64WithContext(ctx, path)
	if err != nil {
		if ctx.Err() == nil {
			log.Printf("Failed to compress image %s: %v", path, err)
		}
		return ""
	}
	return base64Image
}

func waitForProcessingCompletion(ctx context.Context, wg *sync.WaitGroup, base64Images []string) []string {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		wg.Wait()
		return base64Images
	case <-done:
		return base64Images
	}
}

func CompressImageToBase64(imagePath string) (string, error) {
	return CompressImageToBase64WithContext(context.Background(), imagePath)
}

func CompressImageToBase64WithContext(ctx context.Context, imagePath string) (string, error) {
	// Check before starting
	if isContextCancelled(ctx) {
		return "", ctx.Err()
	}

	srcImage, err := imaging.Open(imagePath)
	if err != nil {
		return "", err
	}

	// Check immediately after file I/O
	if isContextCancelled(ctx) {
		return "", ctx.Err()
	}

	resizedImage := resizeImageIfNeeded(srcImage)

	// Check before expensive encoding operation
	if isContextCancelled(ctx) {
		return "", ctx.Err()
	}

	// Pass context to encoding to allow early termination
	return encodeImageToBase64WithContext(ctx, resizedImage)
}

func resizeImageIfNeeded(srcImage image.Image) image.Image {
	bounds := srcImage.Bounds()
	originalWidth := bounds.Dx()

	if originalWidth > maxImageWidth {
		resizedImage := imaging.Resize(srcImage, maxImageWidth, 0, imaging.Lanczos)
		return resizedImage
	}

	return srcImage
}

func encodeImageToBase64(img image.Image) (string, error) {
	return encodeImageToBase64WithContext(context.Background(), img)
}

func encodeImageToBase64WithContext(ctx context.Context, img image.Image) (string, error) {
	// Check context before expensive encoding
	if isContextCancelled(ctx) {
		return "", ctx.Err()
	}

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: imageQuality}

	if err := webp.Encode(&buf, img, opts); err != nil {
		return "", err
	}

	// Check context after encoding completes
	if isContextCancelled(ctx) {
		return "", ctx.Err()
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// ImageToBase64 reads an image file and converts it to base64 without compression
func ImageToBase64(imagePath string) (string, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	// Encode to WebP format without resizing (original quality)
	var buf bytes.Buffer
	opts := &webp.Options{Lossless: true, Quality: 100}

	if err := webp.Encode(&buf, img, opts); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func GetTaskStatus() (int, int, []string) {
	initTaskManager()
	return globalTaskManager.getTaskStatus()
}

func SetMaxSlots(maxSlots int) {
	initTaskManager()

	validatedSlots := validateMaxSlots(maxSlots)
	updateMaxSlots(validatedSlots)
}

func validateMaxSlots(maxSlots int) int {
	if maxSlots <= 0 {
		return 1
	}
	return maxSlots
}

func updateMaxSlots(maxSlots int) {
	globalTaskManager.mu.Lock()
	defer globalTaskManager.mu.Unlock()

	globalTaskManager.maxSlots = maxSlots
	log.Printf("Updated max slots to: %d", maxSlots)
}
