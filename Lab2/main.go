package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DirFileRec(path string, filesChan chan string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".csv" {
			filesChan <- filepath.Join(path, file.Name())
		} else if file.IsDir() {
			DirFileRec(filepath.Join(path, file.Name()), filesChan)
		}
	}
	//close(filesChan) // закриваємо канал
}

/***************** Структура бінарного дерева пошуку*********************/
type tree struct { // структура дерева
	data  []string // дані
	left  *tree    // ліва гілка
	right *tree    // права гілка
}

/***************** Функції бінарного дерева пошуку*********************/
func createTreeVertex(buffer *string) *tree { // створеня вершини гілки
	var vertex = new(tree)
	vertex.data = strings.Split(*buffer, ";") // заповнення вершини початковим значенням
	return vertex
}
func addBranch(vertex *tree, buffer *string, sortByLine *int) { // додаємо нову гілку
	var compare = strings.Split(*buffer, ";")
	if compare[*sortByLine] > vertex.data[*sortByLine] {
		if vertex.right == nil {
			vertex.right = new(tree)
			vertex.right.data = compare
		} else {
			addBranch(vertex.right, buffer, sortByLine)
		}
	} else if compare[*sortByLine] < vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			addBranch(vertex.left, buffer, sortByLine)
		}
	} else if compare[*sortByLine] == vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			insertElem(vertex, compare)
		}
	}
}
func insertElem(vertex *tree, data []string) { // додаємо елемент між гілками
	temp := vertex.left
	vertex.left = new(tree)
	vertex.left.data = data
	vertex.left.left = temp
}

/***************** Функції виведення значень дерева*********************/
func outTree(vertex *tree, file *bufio.Writer) { // рекурсивно виводимо значення по зростанню
	if vertex.left != nil {
		outTree(vertex.left, file)
	}
	writeFile(file, vertex)
	if vertex.right != nil {
		outTree(vertex.right, file)
	}
}
func outTreeRev(vertex *tree, file *bufio.Writer) { // рекурсивно виводимо значення за спаданням
	if vertex.right != nil {
		outTreeRev(vertex.right, file)
	}
	writeFile(file, vertex)
	if vertex.left != nil {
		outTreeRev(vertex.left, file)
	}
}
func writeFile(file *bufio.Writer, temp *tree) { // виводимо масив даних кожної гілки
	for i := 0; i < len(temp.data); i++ { // виводимо масив підстрок у файл
		file.WriteString(temp.data[i] + ";")
	}
	file.WriteString("\n") // новий рядок у файлі
}

/***********************************************************************/

func main() {
	const go_size int = 3
	var (
		path = flag.String("d", ".", "Use a file with the name file-name as an input")
	)
	flag.Parse()

	filesChan := make(chan string)
	isProcessed := make(chan bool)
	filesContent := make(chan string)

	go func() {
		DirFileRec(*path, filesChan)
		close(filesChan)
	}()
	for i := 0; i < go_size; i++ {
		go func() {
			for path := range filesChan {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				reader := bufio.NewReader(file)
				for {
					line, _ := reader.ReadString('\n')
					strings.TrimSuffix(line, "\n")
					if line == "" {
						break
					}
					filesContent <- line
					//fmt.Print(line)
				}
				file.Close()
			}
			isProcessed <- true
		}()
	}
	go func() {
		for content := range filesContent {
			print(content)
		}
	}()
	for i := 0; i < go_size; i++ {
		<-isProcessed
	}
}
