package handlers

import (
	"net/http"
)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (handle *Handle) CORSGetJobs() http.HandlerFunc {
	return corsMiddleware(handle.getJobs)
}

func (handle *Handle) CORSGetDatasets() http.HandlerFunc {
	return corsMiddleware(handle.getDatasets)
}

func (handle *Handle) CORSGetImages() http.HandlerFunc {
	return corsMiddleware(handle.getImages)
}

func (handle *Handle) CORSGetBase64Images() http.HandlerFunc {
	return corsMiddleware(handle.getBase64Images)
}

func (handle *Handle) CORSSavePendingReview() http.HandlerFunc {
	return corsMiddleware(handle.savePendingReview)
}

func (handle *Handle) CORSGetPendingReview() http.HandlerFunc {
	return corsMiddleware(handle.getPendingReview)
}

func (handle *Handle) CORSApprovedRemove() http.HandlerFunc {
	return corsMiddleware(handle.approvedRemove)
}

func (handle *Handle) CORSUnApprovedRemove() http.HandlerFunc {
	return corsMiddleware(handle.unApprovedRemove)
}

func (handle *Handle) CORSGetAllPages() http.HandlerFunc {
	return corsMiddleware(handle.getAllPages)
}
