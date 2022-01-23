package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Data struct {
	commonLetter map[string]int
	commonLetterWord []map[string]int
}

type PoppularityInfo struct {
	String string
	Poppularity int
}

func main() {
	parser("")
	parser("2")
	parser("3")
}

func parser(prefix string) {
	wordsRaw, err := ioutil.ReadFile("./words" + prefix + ".json")
	if err != nil {
		panic(err)
	}
	words := []string{}
	data := Data{}

	data.commonLetter = map[string]int{}
	data.commonLetterWord = make([]map[string]int, 5)

	err = json.Unmarshal(wordsRaw, &words)
	if err != nil {
		panic(err)
	}

	for i, word := range words {
		for letterI, letter := range word {
			data.commonLetter[string(letter)]++
			if i == 0 {
				data.commonLetterWord[letterI] = map[string]int{}		
			}
			data.commonLetterWord[letterI][string(letter)]++
		}
	}

	fmt.Println(`Most common letters: "` + strings.Join(commonLetters(data.commonLetter), "") + `"`)

	poppularityByLetterIndex := [5]map[string]int{} // letter -> index

	for i, dat := range data.commonLetterWord {
		letters := []PoppularityInfo{}
		for lett, pop := range dat {
			letters = append(letters, PoppularityInfo{
				String: lett,
				Poppularity: pop,
			})
		}
		sort.SliceStable(letters, func(i, j int) bool {
			return letters[i].Poppularity > letters[j].Poppularity
		})
		lettersWithIndex := map[string]int{}
		for letterI, info := range letters {
			lettersWithIndex[info.String] = letterI
		}
		poppularityByLetterIndex[i] = lettersWithIndex
	}

	wordsSorted := []PoppularityInfo{}

	for _, word := range words {
		popp := 0
		for lettI, lett := range word {
			popp += poppularityByLetterIndex[lettI][string(lett)]
		}

		wordsSorted = append(wordsSorted, PoppularityInfo{
			String: word,
			Poppularity: popp,
		})
	}

	sort.SliceStable(wordsSorted, func(i, j int) bool {
		return wordsSorted[i].Poppularity < wordsSorted[j].Poppularity
	})

	jsonOutput, err := json.MarshalIndent(wordsSorted, "", "  ")
	
	if err != nil {
		panic(err)
	}

	file, err := os.Create("output" + prefix + ".json")
	
	if err != nil {
		panic(err)
	}

	file.Write(jsonOutput)
}

func commonLetters(data map[string]int) []string {
	commonLetters := []string{}

	commonLetter := ""
	
	for letter, common := range data {
		if data[commonLetter] == common {
			commonLetters = append(commonLetters, letter)
		} else if data[commonLetter] < common {
			commonLetters = []string{letter}
			commonLetter = letter
		}
	}

	return commonLetters
}