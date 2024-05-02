package searcher

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Поисковый индекс файлов.
type Index struct {
	data  map[string][]string // Мапа для хранения индекса (слово/имя файла)
	mutex sync.Mutex
	wg    sync.WaitGroup
}

// Создание экземпляра Index.
func NewIndex() *Index {
	return &Index{
		data: make(map[string][]string),
	}
}

// Индексация списка файлов асинхронно.
func (index *Index) IndexFiles(filenames []string) error {
	errChan := make(chan error)
	defer close(errChan)

	index.wg.Add(len(filenames))

	for _, filename := range filenames {
		go index.IndexFile(filename, errChan)
	}

	var firstError error
	go func() {
		for err := range errChan {
			if firstError == nil {
				firstError = err
			}
		}
		fmt.Println("Обработка ошибок завершена")
	}()

	index.wg.Wait()

	return firstError
}

// Индесирует файл
func (index *Index) IndexFile(filename string, errChan chan<- error) {
	defer index.wg.Done() // Уменьшение счетчика горутин при завершении

	file, err := os.Open(filename) // Открываем файл
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		errChan <- err // Выбрасываем ошибку
		return
	}
	defer file.Close()

	content, _ := ioutil.ReadFile(filename)
	fmt.Println(string(content))

	scanner := bufio.NewScanner(file) // Создание сканера для чтения файла

	for scanner.Scan() {
		fmt.Println("Содержимое файла:", scanner.Text())
		words := strings.Fields(scanner.Text()) // Разбиваем строку на слова
		fmt.Println("Добавление слова в индекс:", words)
		for _, word := range words {
			fmt.Println("Добавление слова в индекс:", word)
			fmt.Println("Добавление слова:", word, "из файла:", filename)
			index.mutex.Lock() // Захватываем мьютекс для потокобезопасного доступа
			index.data[word] = append(index.data[word], filename)
			index.mutex.Unlock() // Освобождаем мьютекс
		}
	}

	if err := scanner.Err(); err != nil { // Проверка сканера на ошибки
		fmt.Println("Ошибка сканера:", err)
		errChan <- err // Выбрасываем ошибку
		index.wg.Add(1)
		return
	}
	fmt.Println("Файл успешно проиндексирован:", filename) // Добавим вывод перед успешным завершением
}

// Поиск по ключевому слову в индексе.
func (index *Index) Search(keyword string) ([]string, error) {
	index.mutex.Lock()
	defer index.mutex.Unlock()

	files, found := index.data[keyword] // Получаем имена файлов для заданного ключевого слова из индекса
	if !found {
		return nil, errors.New("Ключевое слово не найдено в файлах") // Возвращаем ошибку если слово не найдено
	}
	return files, nil // Возвращаем список файлов для найденного ключевого слова
}
