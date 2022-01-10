package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type Items struct {
	Items []Item
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

type Item struct {
	Name     string
	Quantity int
	Price    float32
}

func main() {
	rezFileName := "IFK-8_Gabrielius_Kojelis"
	//jsonFile, err := os.Open("nieks_neatitinka.json")
	//jsonFile, err := os.Open("dalis_atitinka.json")
	jsonFile, err := os.Open("visi_atitinka.json")
	check(err)
	//perskaito json kaip bytu  array
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var items Items
	json.Unmarshal(byteValue, &items)
	jsonFile.Close()
	fileLength := len(items.Items)
	ResultChan := make(chan Item)
	workerChan := make(chan Item)
	deleteChan := make(chan int)
	InsertChan := make(chan Item)
	ResultsSortedChan := make(chan Item, fileLength)
	passed := 0    // kiek is viso praejo filtra duomenu
	doneCount := 0 // kiek isviso istrinta is masyvo duomenu
	//uzdarymas visu channel
	defer close(ResultChan)
	defer close(workerChan)
	defer close(deleteChan)
	defer close(InsertChan)
	defer close(ResultsSortedChan)
	//worker gija siuncia delete , duomenu procesui,
	for i := 0; i < 5; i++ {
		go workerThread(fileLength, workerChan, deleteChan, ResultChan)
	}
	go dataThread(&doneCount, fileLength, deleteChan, InsertChan, workerChan)
	go resultsThread(&passed, ResultChan, ResultsSortedChan)
	for i := 0; i < fileLength; i++ {
		InsertChan <- items.Items[i] //lygiagreciai siuncia insert zinutes..data threadas priimineja
	}
	printResults(ResultsSortedChan, rezFileName, &passed)
}

func workerThread(fileLength int, worker <-chan Item, deleteMessage chan<- int, results chan<- Item) {
	for i := 0; i < fileLength; i++ {
		deleteMessage <- 1
		itemToCheck := <-worker        // priima workeri
		if itemToCheck.Quantity > 40 { //ar tas item kur prieme atitinka filtra
			results <- itemToCheck //ideda i rezultatus ResultThreade jei atitinka
		}
	}
	results <- Item{}
}

func dataThread(doneCount *int, fileLength int, deleteMessage <-chan int, insertZinute <-chan Item, worker chan<- Item) {
	var duom []Item //sukuriam vidini masyva, kuris yra tik sau matomas
	for i := 0; i < fileLength/2; i++ {
		duom = append(duom, Item{}) //Pripildomas puse duomenu faile esancio duomenu kieko masyva tusciais duomenimis, kaip prasoma salygoje
	}
	duomSize := 0
	for *doneCount < fileLength { //(kiek issiusta workeriui) kol mazesnis uz viso failo dydi
		switch size := duomSize; {
		case size <= 0:
			select {
			case item := <-insertZinute: //insert zinute priimta
				duom[duomSize] = item
				duomSize++

			}
			//Priima tik deleteMessage
		case size >= len(duom): //jei yra 15 ar daugiau td ignoruoja insert , o daro delete
			select {
			case <-deleteMessage:
				worker <- duom[duomSize-1]
				duomSize--
				*doneCount++

			}
		default: //tiek insert tiek delete zinutes priema, kuria priims greiciau ta ir darys
			select {
			case item := <-insertZinute:
				duom[duomSize] = item

				duomSize++

			case <-deleteMessage:
				worker <- duom[duomSize-1]
				duomSize--

				*doneCount++

			}
		}
	}
}

func resultsThread(passed *int, results <-chan Item, sortedResults chan<- Item) {
	var resultArray []Item
	for {
		sortedItem := <-results     //Priimamos zinutes is workerthreado
		if (sortedItem == Item{}) { // Jei tuscias itemas, kuri prieme tada yra breakinamas ciklas ir pradedamas kitas ciklas
			break
		}
		index := sort.Search(len(resultArray), func(i int) bool {
			return resultArray[i].Quantity < sortedItem.Quantity
		})
		resultArray = append(resultArray, Item{})
		copy(resultArray[index+1:], resultArray[index:])
		resultArray[index] = sortedItem
		sortedResults <- resultArray[index]
		*passed++
	}
}

func printResults(results <-chan Item, fileName string, size *int) {
	outFile, _ := os.Create(fileName + ".txt")
	header := fmt.Sprintf("|%-32s|%-10s|%5s|\n", "uzsakyta", "kiekis", "kaina")
	fmt.Fprintf(outFile, header)
	str := strings.Repeat("-", utf8.RuneCountInString(header))
	fmt.Fprintf(outFile, str)
	for i := 0; i < *size; i++ { //pagal passed maine size
		item := <-results
		fmt.Fprintf(outFile, "\n|%-32s|%10d|%5.2f|", item.Name, item.Quantity, item.Price)
	}
	fmt.Fprintf(outFile, "\n{Praejo}: %d", *size)
}
