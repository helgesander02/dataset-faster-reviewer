package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func (handle *Handle) savePendingReview(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	body := json.NewDecoder(r.Body)
	itemsLen := handle.DM.SavePendingReviewData(body)
	if itemsLen == 0 {
		log.Println("Failed to save pending review data")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "success", "count": itemsLen})
}

func (handle *Handle) getPendingReview(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	query := r.URL.Query()
	flatten := query.Get("flatten") == "true"
	if flatten {
		writeJSON(w, http.StatusOK, handle.DM.PendingReviewData)
		return
	}
	//writeJSON(w, http.StatusOK, handle.DM.PendingReviewData)
}

func (handle *Handle) approvedRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := handle.DM.GetPendingReviewItems()
	var failedFiles []string
	for _, item := range items {
		filePath := handle.DM.ImageRoot + "/" + item.Job + "/" + item.Dataset + "/" + item.ImagePath
		err := os.Remove(filePath)
		if err != nil {
			log.Printf("Failed to remove file: %s, error: %v", filePath, err)
			failedFiles = append(failedFiles, filePath)
		} else {
			log.Printf("Successfully removed file: %s", filePath)
		}
	}

	if len(failedFiles) > 0 {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":       "partial_success",
			"failed_files": failedFiles,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "success"})
}

func (handle *Handle) unApprovedRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handle.DM.ClearPendingReviewData()
	log.Println("PendingReviewData has been cleared.")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
