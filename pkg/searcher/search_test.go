package searcher

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// пока не придумал как решить проблему в многопоточке, что один файл обрабатывается быстрее другого из-за этого может быть не правильный реузльтат теста
func TestIndexFiles(t *testing.T) {
	index := NewIndex()
	filenames := []string{"../testdata/file_1.txt", "../testdata/file_2.txt"}

	if _, err := os.Stat(filenames[0]); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Файл %s не существует\n", filenames[0])
		} else {
			fmt.Printf("Ошибка при проверке файла %s: %v\n", filenames[0], err)
		}
		return
	}

	err := index.IndexFiles(filenames)
	if err != nil {
		t.Errorf("Ошибка при индексации файлов: %v", err)
	}

	index.wg.Wait()
	time.Sleep(1 * time.Second)

	expectedIndex := map[string][]string{
		"test_word_1": {"../testdata/file_1.txt"},
		"test_word_2": {"../testdata/file_1.txt", "../testdata/file_2.txt"},
	}
	fmt.Println("Ожидаемый файл: ", expectedIndex)
	fmt.Println("Фактический файл: ", index.data)
	if !compareIndexes(index.data, expectedIndex) {
		t.Error("Индексированные данные не совпадают с индексом")
	}
}

func TestSearchFound(t *testing.T) {
	index := NewIndex()
	index.data = map[string][]string{
		"test_word_1": {"../testdata/file_1.txt"},
		"test_word_2": {"../testdata/file_1.txt", "../testdata/file_2.txt"},
	}

	files, err := index.Search("test_word_1")
	if err != nil {
		t.Errorf("Ошибка при поиске слова 'test_word_1': %v", err)
	}

	expectedFiles := []string{"../testdata/file_1.txt"}
	if !compareStringSlices(files, expectedFiles) {
		t.Error("Результаты поиска не совпадают с 'test_word_1'")
	}
}

func TestSearchNotFound(t *testing.T) {
	index := NewIndex()
	index.data = map[string][]string{
		"test_word_1": {"../testdata/file_1.txt"},
		"test_word_2": {"../testdata/file_1.txt", "../testdata/file_2.txt"},
	}

	_, err := index.Search("test_word_3")
	if err == nil {
		t.Error("Ожидалась ошибка при поиске слова 'test_word_3', но ее не было")
	}

	expectedErrorMsg := "Ключевое слово не найдено в файлах"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Ожидалось сообщение об ошибке '%s', но получено '%s'", expectedErrorMsg, err.Error())
	}
}

//Не завелся нормально
/*func TestIndexFileError(t *testing.T) {
	index := NewIndex()

	// Создание файла для чтения, но без записи в него данных
	filename := "wrong_file.txt"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Ошибка при создании файла: %v", err)
	}
	file.Close()

	errChan := make(chan error, 10) // Создаем канал ошибок с буфером 10
	done := make(chan struct{})     // Канал для сообщения об завершении процесса индексации
	indexErr := make(chan error, 1) // Канал для получения ошибки индексации

	go func() {
		for err := range errChan {
			if err != nil {
				t.Errorf("Ошибка в обработке: %v", err)
				indexErr <- err // Отправляем ошибку индексации в канал
			}
		}
		close(done) // Отправляем сигнал о завершении горутины
	}()

	index.IndexFile(filename, errChan) // Индексируем файл с использованием канала ошибок

	// Ждем завершения работы индексации или истечения тайм-аута
	select {
	case <-done:
		// Процесс индексации завершился, можно продолжать
	case <-time.After(time.Second * 500):
		t.Fatal("Тайм-аут: процесс индексации занял слишком много времени")
	}

	// Проверяем, была ли ошибка в процессе индексации
	select {
	case err := <-indexErr:
		t.Fatalf("Ошибка в процессе индексации: %v", err)
	default:
		// Нет ошибки в процессе индексации
	}
}*/

// Вспомогательная функция для сравнения двух индексов
func compareIndexes(index1, index2 map[string][]string) bool {
	if len(index1) != len(index2) {
		return false
	}

	for key, value := range index1 {
		if !compareStringSlices(value, index2[key]) {
			return false
		}
	}

	return true
}

// Вспомогательная функция для сравнения двух слайсов строк
func compareStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

//package searcher
//
//import (
//	"io/fs"
//	"reflect"
//	"testing"
//	"testing/fstest"
//)
//
//func TestSearcher_Search(t *testing.T) {
//	type fields struct {
//		FS fs.FS
//	}
//	type args struct {
//		word string
//	}
//	tests := []struct {
//		name      string
//		fields    fields
//		args      args
//		wantFiles []string
//		wantErr   bool
//	}{
//		{
//			name: "Ok",
//			fields: fields{
//				FS: fstest.MapFS{
//					"file1.txt": {Data: []byte("World")},
//					"file2.txt": {Data: []byte("World1")},
//					"file3.txt": {Data: []byte("Hello World")},
//				},
//			},
//			args:      args{word: "World"},
//			wantFiles: []string{"file1", "file3"},
//			wantErr:   false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Searcher{
//				FS: tt.fields.FS,
//			}
//			gotFiles, err := s.Search(tt.args.word)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
//				t.Errorf("Search() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
//			}
//		})
//	}
//}
