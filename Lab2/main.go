package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var abort bool // глобальна змінна для припинення пошуку файлу за ім'ям в каталозі і його підкаталогах

/***************** Рекурсивні функції для роботи з папками *********************/
func DirFileRec(path string, filesChan chan string) { // рекурсивний прохід по всім файлам заданого дерева каталогу формату .csv
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
}

// /////////////////// Знахождення одного файлу за заданим ім'ям /////////////////////////////
func FindFileDir(path string, filesChan chan string, inputfileName string, abort bool) {
	if abort {
		return
	}

	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == inputfileName && !abort {
			filesChan <- filepath.Join(path, file.Name())
			abort = true
		} else if file.IsDir() {
			if abort {
				break
			}
			FindFileDir(filepath.Join(path, file.Name()), filesChan, inputfileName, abort)
		}
	}
}

/***************** Структура бінарного дерева пошуку*********************/
type tree struct { // структура дерева
	data  []string // дані
	left  *tree    // ліва гілка
	right *tree    // права гілка
}

/***************** Функції бінарного дерева пошуку*********************/
func createTreeVertex(buffer string) *tree { // створеня вершини гілки
	var vertex = new(tree)
	vertex.data = strings.Split(buffer, ";") // заповнення вершини початковим значенням
	return vertex
}

// /////////////////// Додаємо нову гілку /////////////////////////////
func addBranch(vertex *tree, buffer string, sortByLine *int) {
	var compare = strings.Split(buffer, ";")
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

// //////////////// Вставка елементу між гілками /////////////////////////////
func insertElem(vertex *tree, data []string) {
	temp := vertex.left
	vertex.left = new(tree)
	vertex.left = &tree{left: temp, data: data}
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

// //////////////// Виведення значень у консоль /////////////////////////////
func outCLI(vertex *tree) { // рекурсивно виводимо значення по зростанню
	if vertex.left != nil {
		outCLI(vertex.left)
	}
	println(strings.Join(vertex.data, ";"))
	if vertex.right != nil {
		outCLI(vertex.right)
	}
}
func outCLIRev(vertex *tree) { // рекурсивно виводимо значення за спаданням
	if vertex.right != nil {
		outCLIRev(vertex.right)
	}
	println(strings.Join(vertex.data, ";"))
	if vertex.left != nil {
		outCLIRev(vertex.left)
	}
}

// //////////////// Виведення значень у файл /////////////////////////////
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
	file.WriteString(strings.Join(temp.data, ";"))
	file.WriteByte('\n')
}

/***********************************************************************/

func main() {
	const go_size int = 3 // кількість горутин у другому етапі конвеєру

	var (
		path           = flag.String("d", ".", "Use a file with the name file-name as an input")
		inputFileName  = flag.String("i", "", "Use a file with the name file-name as an input")
		sortByLine     = flag.Int("f", 0, "Sort input lines by value number N")
		outputFileName = flag.String("o", "", "Use a file with the name file-name as an output")
		revSort        = flag.Bool("r", false, "Sort input lines in reverse order")
	)
	flag.Parse()

	filesChan := make(chan string)       // шляхи до файлів
	isProcessed := make(chan struct{})   // сигнал стану обробки
	filesContent := make(chan string, 3) // зміст знайдених файлів
	buildTree := make(chan *tree)        // вершина побудованого бінарного дерева пошуку

	/***************** Stage one: Directory Reading*********************/
	go func() {
		if *inputFileName != "" {
			if *path != "." {
				log.Fatal("You can't use -i and -d options at the same time")
			}
			FindFileDir(".", filesChan, *inputFileName, abort) // знаходимо один заданий  файл
		} else {
			DirFileRec(*path, filesChan) // знаходимо усі файли .csv
		}
		close(filesChan)
	}()

	/***************** Stage two: File Reading*********************/
	for i := 0; i < go_size; i++ {
		go func() {
			var line string
			var reader *bufio.Reader

			for path := range filesChan { // зчитуємо зміст файлів з диску
				file, err := os.Open(path) // відкриваємо файл
				if err != nil {
					log.Fatal(err)
				}
				reader = bufio.NewReader(file)
				for {
					line, _ = reader.ReadString('\n') // зчитуємо рядок за рядком
					line = strings.Trim(line, "\n")   // прибраємо зайві байти
					if line == "" {                   // якщо рядок пустий
						break // виходимо з циклу
					}
					filesContent <- line // надсилаємо рядок у буфер
				}
				file.Close() // закриваємо файл
			}
			isProcessed <- struct{}{} // даємо сигнал, що горутину завершила роботу
		}()
	}

	/***************** Stage three: Sorting*********************/
	go func() { // сортуємо бінарним деревом
		var vertex *tree = createTreeVertex(<-filesContent)
		for cont := range filesContent {
			addBranch(vertex, cont, sortByLine)
		}
		buildTree <- vertex // надсилаємо адрес вершини бінарного дерева
		close(buildTree)    // закриваємо канал
	}()
	for i := 0; i < go_size; i++ { // очікуємо завершення усіх горутин з етапу File Reading
		<-isProcessed
	}
	close(filesContent) // закриваємо канал

	if *outputFileName != "" { // якщо задано ім'я файлу для виводу
		outFile, outErr := os.Create(*outputFileName) // створюємо файл за ім'ям

		if outErr != nil {
			log.Fatal(outErr)
		}
		defer outFile.Close() // закриваємо файл після завершення роботи з ним
		writer := bufio.NewWriter(outFile)

		if !*revSort { // параметр -r(сортування за спаданням або зростанням)
			outTree(<-buildTree, writer)
		} else {
			outTreeRev(<-buildTree, writer)
		}
		writer.Flush() // скидуємо дані у файл
	} else { // якщо ім'я файлу для виводу не задано
		if !*revSort { // виводу у консоль результат в залежності від значення флагу -r
			outCLI(<-buildTree)
		} else {
			outCLIRev(<-buildTree)
		}
	}

}
