package handlers

import (
	"log"
	"net/http"
	"strconv"
)

func (handle *Handle) FolderStructureHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	data := handle.DM.ParentData
	writeJSON(w, http.StatusOK, data)
}

func (handle *Handle) getJobs(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	jobNames := handle.DM.GetParentDataAllJobs()
	result := map[string]interface{}{
		"total_jobs": len(jobNames),
		"job_names":  jobNames,
	}
	writeJSON(w, http.StatusOK, result)
}

func (handle *Handle) getDatasets(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	jobName := r.URL.Query().Get("job")
	if jobName == "" {
		http.Error(w, "Missing job parameter", http.StatusBadRequest)
		return
	}
	if !handle.DM.ParentDataJobExists(jobName) {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	datasetNames := handle.DM.GetParentDataAllDatasets(jobName)
	result := map[string]interface{}{
		"total_datasets": len(datasetNames),
		"dataset_names":  datasetNames,
	}
	writeJSON(w, http.StatusOK, result)
}

func (handle *Handle) getImages(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	jobName := r.URL.Query().Get("job")
	datasetName := r.URL.Query().Get("dataset")
	if jobName == "" || datasetName == "" {
		http.Error(w, "Missing job or dataset parameter", http.StatusBadRequest)
		return
	}
	if !handle.DM.ParentDataJobExists(jobName) {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	images := handle.DM.GetParentDataAllImages(jobName, datasetName)
	result := map[string]interface{}{
		"total_images": len(images),
		"images":       images,
	}
	writeJSON(w, http.StatusOK, result)
}

func (handle *Handle) getBase64Images(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	jobName := r.URL.Query().Get("job")
	datasetName := r.URL.Query().Get("dataset")
	pageIndex := r.URL.Query().Get("pageIndex")
	pageNumber := r.URL.Query().Get("pageNumber")
	log.Printf("Fetching images for job: %s, dataset: %s, page_index: %s, page_number: %s", jobName, datasetName, pageIndex, pageNumber)
	if jobName == "" || datasetName == "" {
		http.Error(w, "Missing job or dataset parameter", http.StatusBadRequest)
		return
	}
	if !handle.DM.ParentDataJobExists(jobName) {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	Index, _ := strconv.Atoi(pageIndex)
	number, _ := strconv.Atoi(pageNumber)
	cachedImages, exists := handle.DM.GetImagesCache(jobName, Index, number)
	if exists {
		result := map[string]interface{}{
			"max_page": handle.DM.ImagesPerPageCache.MaxPage,
			"images":   cachedImages,
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if handle.DM.ImagesPerPageCache.MaxPage <= Index {
		http.Error(w, "Page index out of range", http.StatusBadRequest)
		return
	}

	handle.DM.GetBase64ImagesCache(jobName, Index)
	cachedImages, exists = handle.DM.GetImagesCache(jobName, Index, number)
	if !exists {
		http.Error(w, "No images found for the specified job and dataset", http.StatusNotFound)
		return
	}
	result := map[string]interface{}{
		"max_page": handle.DM.ImagesPerPageCache.MaxPage, // 傳遞 MaxPage
		"images":   cachedImages,
	}
	writeJSON(w, http.StatusOK, result)
}

func (handle *Handle) getAllPages(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	jobName := r.URL.Query().Get("job")
	if jobName == "" {
		http.Error(w, "Missing job parameter", http.StatusBadRequest)
		return
	}
	if !handle.DM.ImageCacheJobExists(jobName) {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	pages := handle.DM.GetAllPageDetail()
	result := map[string]interface{}{
		"total_pages": len(pages),
		"pages":       pages,
	}
	writeJSON(w, http.StatusOK, result)
}
