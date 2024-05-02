package controllers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
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
		http.Error(w, "Не удалось получить файл с запроса", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := fc.FileMgr.UploadFile(file, handler.Filename); err != nil {
		http.Error(w, "Не удалось загрузить файл", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Файл загружен"))
}

func (fc *FileController) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := fc.FileMgr.ListFiles()
	if err != nil {
		http.Error(w, "Не удалось получить список файлов", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Не удалось маршрутизировать запрос", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fc *FileController) SearchWordHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "Нет слова для поиска", http.StatusBadRequest)
		return
	}

	// Получаем список файлов
	files, err := fc.FileMgr.ListFiles()
	if err != nil {
		http.Error(w, "Не удалось получить список файлов", http.StatusInternalServerError)
		return
	}

	// Формируем полные пути к файлам
	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = filepath.Join(fc.FileMgr.UploadDir, file)
	}

	// Индексируем файлы
	err = fc.Search.IndexFiles(filePaths)
	if err != nil {
		http.Error(w, "Не удалось проиндексировать файл", http.StatusInternalServerError)
		return
	}

	// Ищем файлы по ключевому слову
	foundFiles, err := fc.Search.Search(word)
	if err != nil {
		http.Error(w, "Не удалось найти слово", http.StatusInternalServerError)
		return
	}

	// Возвращаем найденные файлы в JSON формате
	response, err := json.Marshal(foundFiles)
	if err != nil {
		http.Error(w, "Не удалось маршрутизировать запрос", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
