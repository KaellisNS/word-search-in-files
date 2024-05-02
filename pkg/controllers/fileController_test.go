package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileController_UploadFileHandler(t *testing.T) {
	tempDir := t.TempDir()
	fileController := NewFileController(tempDir)
	server := httptest.NewServer(http.HandlerFunc(fileController.UploadFileHandler))
	defer server.Close()

	fileContents := "Test file contents"
	fileBuffer := bytes.NewBufferString(fileContents)
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Ошибка создания файла: %v", err)
	}
	_, err = io.Copy(part, fileBuffer)
	if err != nil {
		t.Fatalf("Ошибка копирования содержимого файла: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", server.URL, reqBody)
	if err != nil {
		t.Fatalf("Ошибка создания реквеста: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Ошибка отправки реквеста: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидаемый код %d, полученный %d", http.StatusCreated, resp.StatusCode)
	}

	// Verify that the file has been uploaded
	uploadedFiles, err := fileController.FileMgr.ListFiles()
	if err != nil {
		t.Fatalf("Ошибка получения списка файлов: %v", err)
	}
	if len(uploadedFiles) != 1 {
		t.Errorf("Ожидался 1 файл, было получено %d", len(uploadedFiles))
	}

	// Verify the contents of the uploaded file
	uploadedFilePath := filepath.Join(tempDir, uploadedFiles[0])
	uploadedFileContents, err := os.ReadFile(uploadedFilePath)
	if err != nil {
		t.Fatalf("Ошибка чтения загруженного файла %v", err)
	}
	if string(uploadedFileContents) != fileContents {
		t.Errorf("Ожидание содержимого файла: %q, получено: %q", fileContents, uploadedFileContents)
	}
}

func TestFileController_ListFilesHandler(t *testing.T) {
	tempDir := t.TempDir()
	fileController := NewFileController(tempDir)
	server := httptest.NewServer(http.HandlerFunc(fileController.ListFilesHandler))
	defer server.Close()

	// Make a request to list files
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый код %d, полученный %d", http.StatusOK, resp.StatusCode)
	}

	// Decode response JSON
	var files []string
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		t.Fatalf("Ошибка декодирования json: %v", err)
	}

	// Verify the response
	if len(files) != 0 {
		t.Errorf("Ожидали 0 файлов, получено %d", len(files))
	}
}

func TestFileController_SearchWordHandler(t *testing.T) {
	tempDir := t.TempDir()
	fileController := NewFileController(tempDir)
	server := httptest.NewServer(http.HandlerFunc(fileController.SearchWordHandler))
	defer server.Close()

	// Upload a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFilePath, []byte("test file"), 0644)
	if err != nil {
		t.Fatalf("Ошибка записи файла: %v", err)
	}

	// Make a request to search for a word
	resp, err := http.Get(server.URL + "?word=test")
	if err != nil {
		t.Fatalf("Ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидал статус %d, получил %d", http.StatusOK, resp.StatusCode)
	}

	// Decode response JSON
	var foundFiles []string
	err = json.NewDecoder(resp.Body).Decode(&foundFiles)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	// Verify the response
	if len(foundFiles) != 1 || filepath.Base(foundFiles[0]) != "test.txt" {
		t.Errorf("Ожидал ['test.txt'], получил %v", foundFiles)
	}
}
