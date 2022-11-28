package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

/***************** Функції сортування використовуючи однозв'язний список(nodes) *********************/
func createNodeHeader(buffer *string) *node { // створюємо початок списку
	var head node
	head.data = strings.Split(*buffer, ";") // розділяємо Першу уведену строку на масив підстрок
	return &head
}
func addNode(temp **node) { // додаємо вузол
	(*temp).next = new(node)
	(*temp) = (*temp).next
}
func listSortRev(head **node, value int) *node {

	if *head == nil {
		return nil
	}

	var (
		adress  *node
		maximum string = (*head).data[value]
		temp    *node  = (*head)
		isSwap  bool   = true
	)

	for isSwap {
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
		} else if (*head).next != nil {
			(*head).next = listSortRev(&((*head).next), value)
		}
	}
	return *head
}
func listSort(head **node, value int) *node {

	if *head == nil {
		return nil
	}

	var (
		adress  *node
		minimum string = (*head).data[value]
		temp    *node  = (*head)
		isSwap  bool
	)

	for isSwap = true; isSwap; {
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
		} else if (*head).next != nil {
			(*head).next = listSort(&((*head).next), value)
		}
	}
	return *head
}
func headLine(file *bufio.Writer, reader *bufio.Reader) {
	line, _ := reader.ReadString('\n')
	file.WriteString(strings.TrimSuffix(line, "\n"))
	file.WriteByte('\n')
}
func main() {
	var (
		buffer  string     // буферна змінна для уведених рядків
		counter int    = 1 // лічильник кількості елементів однозв'язного списку
	)
	var (
		inputFileName  = flag.String("i", "CLI", "Use a file with the name file-name as an input")
		outputFileName = flag.String("o", "output.csv", "Use a file with the name file-name as an output")
		headOp         = flag.Bool("h", true, "The first line is a header that must be ignored during sorting but included in the output")
		sortByLine     = flag.Int("f", 0, "Sort input lines by value number N")
		revSort        = flag.Bool("r", false, "Sort input lines in reverse order")
		treeSort       = flag.Int("a", 1, "Sorty by tree or default algorithm")
	)
	flag.Parse()

	outFile, outErr := os.Create(*outputFileName) // // створюємо файл для виводу
	if outErr != nil {                            // якщо отримали помилку, завершуємо роботу програми
		fmt.Println("Unable to create output file", outErr)
		os.Exit(1)
	}
	defer outFile.Close()              // завершуємо роботу з файлом
	writer := bufio.NewWriter(outFile) // створюємо поток запису через буфер

	/***************** СОРТУВАННЯ ВХІДНОГО ФАЙЛУ *********************/
	/************************************************************/
	if *inputFileName != "CLI" {
		inpFile, err := os.Open(*inputFileName)
		if err != nil {
			fmt.Println("Unable to open file:", err)
			os.Exit(1)
		}
		defer inpFile.Close() // завершуємо роботу з файлом
		reader := bufio.NewReader(inpFile)

		if *headOp {
			headLine(writer, reader)
		}

		line, _ := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		switch *treeSort {
		case 1: ///// СОРТУВАННЯ ПОШУКОМ НАЙМЕНШОГО ЕЛЕМЕНТУ СПИСКУ /////
			_startCell := createNodeHeader(&line) // початковий вузол списку
			for temp := _startCell; err != io.EOF; {
				addNode(&temp)
				line, err = reader.ReadString('\n')
				temp.data = strings.Split(strings.TrimSuffix(line, "\n"), ";")
			}
			if *revSort {
				listSortRev(&_startCell, *sortByLine)
			} else {
				listSort(&_startCell, *sortByLine)
			}
			for temp := _startCell; temp != nil; temp = temp.next {
				for i := 0; i < len(temp.data); i++ {
					writer.WriteString(temp.data[i] + ";")
				}
				writer.WriteByte('\n')
			}
		case 2: ///// СОРТУВАННЯ ДЕРЕВОМ ПОШУКУ /////
			var vertex *tree = createTreeVertex(&line) // створюємо вершину
			for err != io.EOF {                        // читаємо рядки до кінця файлу
				line, err = reader.ReadString('\n')
				line = strings.TrimSuffix(line, "\n")
				addBranch(vertex, &line, sortByLine)
			}
			if *revSort {
				outTreeRev(vertex, writer)
			} else {
				outTree(vertex, writer)
			}
		}
		writer.Flush() // записуємо дані у файл
		os.Exit(0)
	}

	/***************** СОРТУВАННЯ З CLI *********************/
	/************************************************************/

	fmt.Println("Input CSV data line by line:")
	n, _ := fmt.Fscanln(os.Stdin, &buffer) // вводимо рядок, зберігаємо в змінну
	if *headOp && n != 0 {
		writer.WriteString(buffer + "\n")
		n, _ = fmt.Fscanln(os.Stdin, &buffer)
	}

	switch *treeSort { // в залежності від обраного типу сортування(прапор -а)

	/***************** СОРТУВАННЯ ПОШУКОМ МІНІМАЛЬНОГО ЕЛЕМЕНТУ *********************/
	/************************************************************/
	case 1:
		if n == 0 { // інакше завершуємо програму з повідомленням щодо відсутності введених даних
			fmt.Println("You input no data")
			return
		}
		_startCell := createNodeHeader(&buffer)                        // початковий вузол списку
		temp := _startCell                                             // змінна для зберігання адреси тимчасового вузла
		for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; counter++ { // створюємо нові елементи списку(nodes)
			addNode(&temp)
			temp.data = strings.Split(buffer, ";") // заповнуємо ноди введеними значеннями в консоль
			n, _ = fmt.Fscanln(os.Stdin, &buffer)  // скануємо наступний уведений рядок
		}

		if *revSort {
			listSortRev(&_startCell, *sortByLine)
		} else {
			listSort(&_startCell, *sortByLine)
		}

		for temp = _startCell; temp != nil; temp = temp.next {
			for i := 0; i < len(temp.data); i++ {
				writer.WriteString(temp.data[i] + ";")
			}
			writer.WriteByte('\n')
		}
	/***************** СОРТУВАННЯ ДЕРЕВОМ *********************/
	/************************************************************/
	case 2:
		if n == 0 { // інакше завершуємо програму з повідомленням щодо відсутності введених даних
			fmt.Println("You input no data")
			return
		}
		var vertex *tree = createTreeVertex(&buffer)
		for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; { // створюємо нові елементи списку(nodes)
			addBranch(vertex, &buffer, sortByLine)
			n, _ = fmt.Fscanln(os.Stdin, &buffer) // скануємо наступний уведений рядок
		}
		if *revSort {
			outTreeRev(vertex, writer)
		} else {
			outTree(vertex, writer)
		}
	}
	writer.Flush()
}
