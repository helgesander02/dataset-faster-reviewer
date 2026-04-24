package models_verify_viewer

import (
	"fmt"
	"os"
)

// Merge merges another PendingReview into this one, preserving items that exist in both
func (pr *PendingReview) Merge(other *PendingReview) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	// Create map of new items by key
	newMap := make(map[string]PendingReviewItem, len(other.items))
	for _, item := range other.items {
		newMap[item.Key()] = item
	}

	// Keep old items that exist in new set
	merged := make([]PendingReviewItem, 0, len(other.items))
	for _, oldItem := range pr.items {
		if _, exists := newMap[oldItem.Key()]; exists {
			merged = append(merged, oldItem)
			delete(newMap, oldItem.Key())
		}
	}

	// Add new items that weren't in old set
	for _, item := range other.items {
		if _, exists := newMap[item.Key()]; exists {
			merged = append(merged, item)
		}
	}

	pr.items = merged
}

func (pr *PendingReview) Clear() {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.items = pr.items[:0]
}

func (pr *PendingReview) Items() []PendingReviewItem {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	items := make([]PendingReviewItem, len(pr.items))
	copy(items, pr.items)
	return items
}

func (pr *PendingReview) Replace(items []PendingReviewItem) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.items = items
}

func (pr *PendingReview) Add(item PendingReviewItem) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	for _, existing := range pr.items {
		if existing.Key() == item.Key() {
			return
		}
	}

	pr.items = append(pr.items, item)
}

func (pr *PendingReview) Remove(item PendingReviewItem) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	key := item.Key()
	newItems := make([]PendingReviewItem, 0, len(pr.items))
	for _, existing := range pr.items {
		if existing.Key() != key {
			newItems = append(newItems, existing)
		}
	}
	pr.items = newItems
}

func (pr *PendingReview) Len() int {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	return len(pr.items)
}

func DeleteImageFile(imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", imagePath)
	}

	return os.Remove(imagePath)
}
