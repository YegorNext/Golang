package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type node struct { // структура даних для методу перебору
	data []string // дані
	next *node    // адреси наступих вузлів(нодів)
}

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
func outTree(vertex *tree, file *os.File) { // рекурсивно виводимо значення по зростанню
	if vertex.left != nil {
		outTree(vertex.left, file)
	}
	writeFile(file, vertex)
	if vertex.right != nil {
		outTree(vertex.right, file)
	}
}
func outTreeRev(vertex *tree, file *os.File) { // рекурсивно виводимо значення за спаданням
	if vertex.right != nil {
		outTreeRev(vertex.right, file)
	}
	writeFile(file, vertex)
	if vertex.left != nil {
		outTreeRev(vertex.left, file)
	}
}
func writeFile(file *os.File, temp *tree) { // виводимо масив даних кожної гілки
	for i := 0; i < len(temp.data); i++ { // виводимо масив підстрок у файл
		file.WriteString(temp.data[i] + ";")
	}
	file.WriteString("\n") // новий рядок у файлі
}

/***************** Функції сортування перебером використовуючи однозв'язний список(nodes) *********************/
func createNodeHeader(buffer *string) *node { // створюємо початок списку
	var head node
	head.data = strings.Split(*buffer, ";") // розділяємо Першу уведену строку на масив підстрок
	return &head
}
func addNode(temp **node) { // додаємо вузол
	(*temp).next = new(node)
	(*temp) = (*temp).next
	(*temp).next = nil
}
func listSortRev(head **node, size, value int) {
	var (
		adress  *node
		maximum string = (*head).data[value]
		temp           = (*head)
		isSwap  bool
	)
	for i := 0; i < size; i++ {
		isSwap = false
		for temp.next != nil {
			if temp.next.data[value] > maximum {
				adress = temp
				maximum = temp.next.data[value]
				isSwap = true
			}
			temp = temp.next
		}
		if isSwap {
			temp = adress.next
			adress.next = temp.next
			temp.next = *head
			*head = temp
		} else {
			break
		}
	}
}
func listSort(head **node, size, value int) {
	var (
		adress  *node
		minimum string = (*head).data[value]
		temp           = (*head)
		isSwap  bool
	)
	for i := 0; i < size; i++ {
		isSwap = false
		for temp.next != nil {
			if temp.next.data[value] < minimum {
				adress = temp
				minimum = temp.next.data[value]
				isSwap = true
			}
			temp = temp.next
		}
		if isSwap {
			temp = adress.next
			adress.next = temp.next
			temp.next = *head
			*head = temp
		} else {
			break
		}
	}
}
func main() {
	var (
		buffer  string     // буферна змінна для уведених рядків
		counter int    = 1 // лічильник кількості елементів однозв'язного списку
	)
	var (
		inputFileName  = flag.String("i", "input.csv", "Use a file with the name file-name as an input")
		outputFileName = flag.String("o", "output.csv", "Use a file with the name file-name as an output")
		headOp         = flag.Bool("h", true, "The first line is a header that must be ignored during sorting but included in the output")
		sortByLine     = flag.Int("f", 0, "Sort input lines by value number N")
		revSort        = flag.Bool("r", false, "Sort input lines in reverse order")
		treeSort       = flag.Int("a", 1, "Sorty by tree or default algorithm")
	)
	flag.Parse()

	inpFile, inpErr := os.Create(*inputFileName)  // створюємо файл для вводу
	outFile, outErr := os.Create(*outputFileName) // // створюємо файл для виводу

	if inpErr != nil { // якщо отримали помилку, завершуємо роботу програми
		fmt.Println("Unable to create input file", inpErr)
		os.Exit(1)
	} else if outErr != nil { // якщо отримали помилку, завершуємо роботу програми
		fmt.Println("Unable to create output file", outErr)
		os.Exit(1)
	}
	defer inpFile.Close() // завершуємо роботу з файлом
	defer outFile.Close() // завершуємо роботу з файлом

	fmt.Println("Input CSV data line by line:")
	n, _ := fmt.Fscanln(os.Stdin, &buffer) // вводимо рядок, зберігаємо в змінну
	if n == 0 {                            // інакше завершуємо програму з повідомленням щодо відсутності введених даних
		fmt.Println("You input no data")
		os.Exit(1)
	}
	if *headOp {
		outFile.WriteString(buffer + "\n")
		n, _ = fmt.Fscanln(os.Stdin, &buffer)
	}

	switch *treeSort { // в залежності від обраного типу сортування(прапор -а)

	/***************** СОРТУВАННЯ ПОШУКОМ МІНІМАЛЬНОГО ЕЛЕМЕНТУ *********************/
	/************************************************************/
	case 1:
		_startCell := createNodeHeader(&buffer)                        // початковий вузол списку
		temp := _startCell                                             // змінна для зберігання адреси тимчасового вузла
		for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; counter++ { // створюємо нові елементи списку(nodes)
			addNode(&temp)
			temp.data = strings.Split(buffer, ";") // заповнуємо ноди введеними значеннями в консоль
			n, _ = fmt.Fscanln(os.Stdin, &buffer)  // скануємо наступний уведений рядок
		}

		if *revSort {
			listSortRev(&_startCell, counter, *sortByLine)
		} else {
			listSort(&_startCell, counter, *sortByLine)
		}

		for temp = _startCell; temp != nil; temp = temp.next {
			for i := 0; i < len(temp.data); i++ {
				fmt.Print(temp.data[i] + ";")
			}
		}
	/***************** СОРТУВАННЯ ДЕРЕВОМ *********************/
	/************************************************************/
	case 2:
		if n != 0 { // якщо рядок НЕ пустий
			var vertex *tree = createTreeVertex(&buffer)
			for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; { // створюємо нові елементи списку(nodes)
				addBranch(vertex, &buffer, sortByLine)
				n, _ = fmt.Fscanln(os.Stdin, &buffer) // скануємо наступний уведений рядок
			}
			if *revSort {
				outTreeRev(vertex, outFile)
			} else {
				outTree(vertex, outFile)
			}
		} else { // інакше завершуємо програму з повідомленням щодо відсутності введених даних
			fmt.Println("You input no data")
			os.Exit(1)
		}
	}
}
