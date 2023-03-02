package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	keyLength = 4
	alphabet  = []string{"А", "Б", "В", "Г", "Д", "Е", "Ж", "З", "И", "Й", "К",
		"Л", "М", "Н", "О", "П", "Р", "С", "Т", "У", "Ф", "Х",
		"Ц", "Ч", "Ш", "Щ", "Ъ", "Ы", "Ь", "Э", "Ю", "Я", "_"}
)

func main() {
	hack()
}

func prepareData() ([][]string, int, string) {
	data := readFile("2decode")
	letters := strings.Split(data, "")
	log.Println("символов", len(letters))
	columnSymbols := len(letters) / keyLength
	log.Println("символов в одном столбике", columnSymbols)
	if columnSymbols*keyLength != len(letters) {
		letters = letters[:columnSymbols*keyLength]
	}
	log.Println("обновленное количество символов", len(letters))
	return chunkBy(letters), columnSymbols, strings.Join(letters, "")
}

func calcEntries(columns [][]string) map[int][]int {
	entries := map[int]map[string]int{}
	for i, column := range columns {
		// i - вхождение
		entries[i] = map[string]int{}
		for _, letter := range alphabet {
			for _, columnLetter := range column {
				if letter == columnLetter {
					if _, exists := entries[i][letter]; exists {
						entries[i][letter] += 1
					} else {
						entries[i][letter] = 1
					}
				}
			}
		}

	}
	numEntries := map[int][]int{}
	for i, entry := range entries {
		log.Println(fmt.Sprintf("вхождения для %d строки", i+1))
		var numbers []string
		var nums []int
		for _, letter := range alphabet {
			numbers = append(numbers, strconv.Itoa(entry[letter]))
			nums = append(nums, entry[letter])
		}
		numEntries[i] = nums
		log.Println(strings.Join(numbers, " "))
	}
	return numEntries
}

func hack() {
	columns, columnSymbols, encodedText := prepareData()
	// подсчет вхождений для каждого столбца
	entries := calcEntries(columns)
	positions := getPositions(entries, columnSymbols)
	var newAlphabet []string
	newAlphabet = append(newAlphabet, alphabet...)
	newAlphabet = append(newAlphabet, alphabet...)
	newAlphabet = append(newAlphabet, alphabet...)
	keyColumns := map[int][]string{}
	for i := 0; i < keyLength; i++ {
		keyColumns[i] = []string{}
	}
	for i, column := range keyColumns {
		startFrom := len(alphabet) - positions[i] + 1
		for j := 0; j < len(alphabet); j++ {
			column = append(column, newAlphabet[startFrom+j])
		}
		keyColumns[i] = column
	}
	for i := range alphabet {
		var word []string
		for j := 0; j < keyLength; j++ {
			word = append(word, keyColumns[j][i])
		}
		log.Println(i+1, strings.Join(word, ""))
	}
	decode(newAlphabet, positions, encodedText)
}

func decode(newAlphabet []string, positions []int, encodedText string) {
	fmt.Print("слово, которым кодировали текст .... ?\n")
	k1 := ""
	_, err := fmt.Fscan(os.Stdin, &k1)
	if err != nil {
		panic(err)
	}
	for i, el := range strings.Split(k1, "") {
		for j, l := range alphabet {
			if l == el {
				positions[i] = j
			}
		}
	}
	encodedLetters := strings.Split(encodedText, "")
	var decodedText []string
	for i, letter := range encodedLetters {
		letterPos := -1
		for j, l := range alphabet {
			if l == letter {
				letterPos = j
				break
			}
		}
		delta := letterPos - positions[i%4]
		newLetter := newAlphabet[len(alphabet)+delta]
		decodedText = append(decodedText, newLetter)
	}
	log.Println(strings.Join(decodedText, ""))
}

func getPositions(entries map[int][]int, length int) []int {
	firstEntry := entries[0]
	indexes := map[int][]float64{}
	for i := 1; i < len(entries); i++ {
		entry := entries[i]
		var comb [][]int
		for j := 0; j < len(alphabet); j++ {
			comb = append(comb, []int{})
			for k := 0; k < len(alphabet); k++ {
				comb[j] = append(comb[j], 0)
			}
		}
		comb[0] = entry
		for j := 1; j < len(alphabet); j++ {
			for k := 0; k < len(alphabet); k++ {
				if k == 0 {
					comb[j][k] = comb[j-1][len(alphabet)-1]
				} else {
					comb[j][k] = comb[j-1][k-1]
				}
			}
		}
		for k := 0; k < len(alphabet); k++ {
			sum := 0
			for j := 0; j < len(alphabet); j++ {
				sum += firstEntry[j] * comb[k][j]
			}
			indexes[i] = append(indexes[i], float64(sum)/(float64(length)*float64(length)))
		}
	}
	indexesJson, _ := json.MarshalIndent(indexes, " ", " ")
	log.Println("взаимные индексы совпадения", string(indexesJson))
	positions := []int{0}
	for i, numbers := range indexes {
		max := 0.0
		maxPos := -1
		for j, num := range numbers {
			if num > max {
				max = num
				maxPos = j
			}
		}
		positions = append(positions, maxPos)
		log.Println(fmt.Sprintf("максимум в %d %.8f, позиция: %d", i, max, maxPos))
	}
	return positions
}

func chunkBy(items []string) (chunks [][]string) {
	var res [][]string
	for i := 0; i < keyLength; i++ {
		res = append(res, []string{})
	}
	for i, letter := range items {
		if i%keyLength == 0 {
			res[0] = append(res[0], letter)
		} else if i%keyLength == 1 {
			res[1] = append(res[1], letter)
		} else if i%keyLength == 2 {
			res[2] = append(res[2], letter)
		} else if i%keyLength == 3 {
			res[3] = append(res[3], letter)
		}
	}
	return res
}
