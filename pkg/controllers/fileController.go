package controllers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
	"word-search-in-files/pkg/handlers/files"
	"word-search-in-files/pkg/searcher"
)

type FileController struct {
	FileMgr *files.FileManager
	Search  *searcher.Index
}

func NewFileController(uploadDir string) *FileController {
	return &FileController{
		FileMgr: files.NewFileManager(uploadDir),
		Search:  searcher.NewIndex(),
	}
}

func (fc *FileController) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := fc.FileMgr.UploadFile(file, handler.Filename); err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
}

func (fc *FileController) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := fc.FileMgr.ListFiles()
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fc *FileController) SearchWordHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "Word parameter is missing", http.StatusBadRequest)
		return
	}

	files, err := fc.FileMgr.ListFiles()
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	foundFiles := make([]string, 0)

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := fc.Search.IndexFile(filepath.Join(fc.FileMgr.UploadDir, file))
			if err != nil {
				http.Error(w, "Failed to index file", http.StatusInternalServerError)
				return
			}
			if _, found := fc.Search.Search(word); found {
				foundFiles = append(foundFiles, file)
			}
		}(file)
	}

	wg.Wait()

	response, err := json.Marshal(foundFiles)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
