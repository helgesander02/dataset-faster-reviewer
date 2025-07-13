package main

import (
	"backend/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	root string
	port string
)

func init() {
	// 1. Load environment variables from .env file
	// 2. Get the path to the static folder from the .env file
	// 3. Get the port from the .env file (default to 8080 if not set)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	root = os.Getenv("STATIC_FOLDER")
	if root == "" {
		log.Fatal("STATIC_FOLDER environment variable not set")
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	handle := handlers.NewHandle(root)
	handle.SetupAPI()

	http.HandleFunc("/api/folder-structure", handle.FolderStructureHandler)
	http.HandleFunc("/api/getJobs", handle.CORSGetJobs())
	http.HandleFunc("/api/getDatasets", handle.CORSGetDatasets())
	http.HandleFunc("/api/getImages", handle.CORSGetImages())
	http.HandleFunc("/api/getBase64Images", handle.CORSGetBase64Images())
	http.HandleFunc("/api/getAllPages", handle.CORSGetAllPages())
	http.HandleFunc("/api/savePendingReview", handle.CORSSavePendingReview())
	http.HandleFunc("/api/getPendingReview", handle.CORSGetPendingReview())
	http.HandleFunc("/api/approvedRemove", handle.CORSApprovedRemove())
	http.HandleFunc("/api/unapprovedRemove", handle.CORSUnApprovedRemove())

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(root))))

	log.Printf("Server started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
