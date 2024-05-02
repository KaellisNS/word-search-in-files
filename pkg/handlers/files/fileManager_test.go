package files

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadFile(t *testing.T) {
	uploadDir := "../testdata"
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		t.Fatalf("Ошибка создания директории: %v", err)
	}
	defer os.RemoveAll(uploadDir)

	fileManager := NewFileManager(uploadDir)

	testData := []byte("test data")
	filename := "test.txt"
	err = fileManager.UploadFile(bytes.NewReader(testData), filename)
	if err != nil {
		t.Fatalf("UploadFile вернул ошибку: %v", err)
	}

	filePath := filepath.Join(uploadDir, filename)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatalf("Загруженный файл %s не найден", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Ошибка чтения загруженного файла: %v", err)
	}

	if !bytes.Equal(content, testData) {
		t.Fatalf("Содержимое загруженного файла не правильно")
	}
}

func TestListFiles(t *testing.T) {
	uploadDir := "../testdata"
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		t.Fatalf("Ошибка создания директории: %v", err)
	}
	defer os.RemoveAll(uploadDir)

	fileManager := NewFileManager(uploadDir)

	// Create some test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(uploadDir, filename)
		if _, err := os.Create(filePath); err != nil {
			t.Fatalf("Ошибка создания файла %v", err)
		}
		defer os.Remove(filePath)
	}

	files, err := fileManager.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles вернул ошибку: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Fatalf("Ожидалось %d файлов, вернулось %d", len(testFiles), len(files))
	}

	for _, filename := range testFiles {
		found := false
		for _, file := range files {
			if file == filename {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Ожидаемый файл %s не найден", filename)
		}
	}
}
