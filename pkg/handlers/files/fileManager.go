package files

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileManager struct {
	UploadDir string
}

func NewFileManager(uploadDir string) *FileManager {
	return &FileManager{UploadDir: uploadDir}
}

func (fm *FileManager) UploadFile(file io.Reader, filename string) error {
	dst, err := os.Create(filepath.Join(fm.UploadDir, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) ListFiles() ([]string, error) {
	files, err := ioutil.ReadDir(fm.UploadDir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

/*func FilesFS(fsys fs.FS, dir string) ([]string, error) {
	if dir == "" {
		dir = "."
	}
	var fileNames []string
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}
*/
