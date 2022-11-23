package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
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
func createNodeHeader(buffer *string) node { // створюємо початок списку
	var _startCell node
	_startCell.data = strings.Split(*buffer, ";") // розділяємо Першу уведену строку на масив підстрок
	_startCell.next = nil
	return _startCell
}
func addNode(temp **node) { // додаємо вузол
	(*temp).next = new(node)
	(*temp) = (*temp).next
	(*temp).next = nil
}
func headLineOptionSet(temp **node, _startCell *node, file *os.File, i *int) { // встановлюємо опцію заголовку(-h)
	(*temp) = _startCell.next
	writeIn(file, _startCell)
	(*i)--
}
func nodeBegin(temp **node, _startCell *node) { // повертаємо список до початку
	(*temp) = _startCell
}
func nextNode(temp **node, _startCell *node, headOp *bool) { // перемикання вузла в залежності від статусу обраної опції
	if *headOp {
		if (*temp).next == nil { // робимо наступні умови, доки не відшукаємо усі елементи
			(*temp) = _startCell.next
		} else {
			(*temp) = (*temp).next
		}
	} else {
		if (*temp).next == nil { // робимо наступні умови, доки не відшукаємо усі елементи
			(*temp) = _startCell
		} else {
			(*temp) = (*temp).next
		}
	}
}

func writeIn(file *os.File, temp *node) { // виводимо масив підстрок у файл
	for i := 0; i < len(temp.data); i++ {
		file.WriteString(temp.data[i] + ";")
	}
	file.WriteString("\n")
}
func sortUp(temp *node, counter int, outFile *os.File, headOp *bool, _startCell *node, sortByLine *int, str []string) { // сортування за зростанням перебором
	if *headOp {
		temp = _startCell.next
	} else {
		temp = _startCell
	}
	for i := 0; i < counter; {
		if temp.data[*sortByLine] == str[i] {
			writeIn(outFile, temp)
			i++
		}
		nextNode(&temp, _startCell, headOp)
	}
}
func sortRev(temp *node, counter int, outFile *os.File, headOp *bool, _startCell *node, sortByLine *int, str []string) { // сортування за спаданням перебором
	if *headOp { // призначуємо для temp адрес вузла списку
		temp = _startCell.next
	} else {
		temp = _startCell
	}
	for i := counter - 1; i >= 0; {
		if temp.data[*sortByLine] == str[i] { // шукаємо i елемент відсортованого масиву у списку(методом перебору)
			writeIn(outFile, temp)
			i--
		}
		nextNode(&temp, _startCell, headOp)
	}
}
func main() {
	var (
		_startCell node                 // початковий вузол списку
		temp       *node  = &_startCell // змінна для зберігання адреси тимчасового вузла
		buffer     string               // буферна змінна для уведених рядків
		counter    int    = 1           // лічильник кількості елементів однозв'язного списку
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
	if n != 0 {                            // якщо рядок НЕ пустий
		inpFile.WriteString(buffer + "\n")
	} else { // інакше завершуємо програму з повідомленням щодо відсутності введених даних
		fmt.Println("You input no data")
		os.Exit(1)
	}

	switch *treeSort { // в залежності від обраного типу сортування(прапор -а)

	/***************** СОРТУВАННЯ ПЕРЕБОРОМ *********************/
	/************************************************************/
	case 1:
		for _startCell = createNodeHeader(&buffer); n != 0; counter++ { // створюємо нові елементи списку(nodes)
			inpFile.WriteString(buffer + "\n")
			addNode(&temp)
			temp.data = strings.Split(buffer, ";") // заповнуємо ноди введеними значеннями в консоль
			n, _ = fmt.Fscanln(os.Stdin, &buffer)  // скануємо наступний уведений рядок
		}

		if *headOp { // вмикаємо опцію заголовку -h
			headLineOptionSet(&temp, &_startCell, outFile, &counter)
		} else { // інакше повертаємо список до початку
			nodeBegin(&temp, &_startCell)
		}

		str := make([]string, counter) // створюємо масив за кількістю елементів однозв'язного списку
		for i := 0; i < counter; i++ { // заповнюємо масив першими значеннями масиву підстрок з кожного елементу списку
			str[i] = temp.data[*sortByLine]
			temp = temp.next // гортаємо список
		}
		sort.Strings(str) // сортуємо елементи

		switch *revSort { // сортування за зростанням або спаданням в залежності від опції
		case true:
			sortRev(temp, counter, outFile, headOp, &_startCell, sortByLine, str)
		case false:
			sortUp(temp, counter, outFile, headOp, &_startCell, sortByLine, str)
		}

	/***************** СОРТУВАННЯ ДЕРЕВОМ *********************/
	/************************************************************/
	case 2:
		var vertex *tree
		if *headOp {
			outFile.WriteString(buffer + "\n")
			n, _ = fmt.Fscanln(os.Stdin, &buffer)
			inpFile.WriteString(buffer + "\n")
			vertex = createTreeVertex(&buffer)
		} else {
			vertex = createTreeVertex(&buffer)
		}
		for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; { // створюємо нові елементи списку(nodes)
			inpFile.WriteString(buffer + "\n")
			addBranch(vertex, &buffer, sortByLine)
			n, _ = fmt.Fscanln(os.Stdin, &buffer) // скануємо наступний уведений рядок
		}
		if *revSort {
			outTreeRev(vertex, outFile)
		} else {
			outTree(vertex, outFile)
		}
	}
}
