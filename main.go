package word_search_in_files

import (
	"net/http"
	"word-search-in-files/pkg/controllers"
)

func main() {
	controller := controllers.NewFileController("./uploads")

	http.HandleFunc("/api/v1/upload", controller.UploadFileHandler)
	http.HandleFunc("/api/v1/list", controller.ListFilesHandler)
	http.HandleFunc("/api/v1/search", controller.SearchWordHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
