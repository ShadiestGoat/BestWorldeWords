package solver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Poppularity struct {
	String string
	Poppularity int
}

type BoardStateType int8

const (
	UNKNOWN BoardStateType = iota
	YELLOW
	GREEN
	FALSE
)

type CurGuess struct {
	WordStatus [5]LetterInfoGuess
	Greeeners map[string]bool
}

var curWordStatus = CurGuess{}

type LetterInfoGuess struct {
	Letter string
	Found bool
}

type LetterInfoMap struct {
	Letter string
	ImoossiblePlaces [5]bool
	IsFalse bool
}

type BoardState map[string]LetterInfoMap

var wordsleft []Poppularity
var letterInfo = BoardState{}
var yelloWLetters = []string{}
// Add known yellows!

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
    var s string
    r := bufio.NewReader(os.Stdin)
    for {
        fmt.Fprint(os.Stderr, label+" ")
        s, _ = r.ReadString('\n')
        if s != "" {
            break
        }
    }
    return strings.TrimSpace(s)
}

func MainSolver() {
	wordsRaw, err := ioutil.ReadFile("output.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(wordsRaw, &wordsleft)
	curWordStatus.Greeeners = map[string]bool{}
	curWordStatus.WordStatus = [5]LetterInfoGuess{}
	for i := 0; i < 6; i++ {
		fmt.Println("Try... " + wordsleft[0].String)
		status := StringPrompt("What was the output?")

		wordCur := []string{}

		for _, letter := range wordsleft[0].String {
			wordCur = append(wordCur, string(letter))
		}

		guess := WordGuess{}

		if len(status) == 1 {
			// This is fine, because we know that the characters are only g, y & f.
			letter := string(status[0])
			for i := 0; i < 4; i++ {
				status += letter
			}
		}

		if strings.ToLower(status) == "ggggg" {
			fmt.Println("Got em ^^")
			break
		}

		for i, letter := range status {
			guess.Letters[i] = wordCur[i]
			switch strings.ToLower(string(letter)) {
			case "g":
				guess.Output[i] = GREEN
			case "y":
				guess.Output[i] = YELLOW
			case "f":
				guess.Output[i] = FALSE
			default: 
				panic("This.. this isn't a good status. You can only use g, y and f!")
			}
		}

		tryWord(guess)
	}
}

type WordGuess struct {
	Letters [5]string
	Output [5]BoardStateType
}

func tryWord(guess WordGuess) {
	for i, letter := range guess.Letters {
		newInfo := letterInfo[letter]
		switch guess.Output[i] {
		case GREEN:
			newLetters := []string{}
			curWordStatus.WordStatus[i].Letter = letter
			curWordStatus.WordStatus[i].Found = true
			curWordStatus.Greeeners[letter] = true
			for _, yL := range yelloWLetters {
				if letter == yL {
					continue
				}
				newLetters = append(newLetters, yL)
			}
			yelloWLetters = newLetters
		case YELLOW:
			newInfo.ImoossiblePlaces[i] = true
			yelloWLetters = append(yelloWLetters, letter)
		case FALSE:
			if !curWordStatus.Greeeners[letter] {
				newInfo.IsFalse = true
			}
		}
		letterInfo[letter] = newInfo
	}

	newWords := []Poppularity{}

	for _, word := range wordsleft {
		curWordInfo := true
		for letterI, letter := range word.String {
			if letterInfo[string(letter)].IsFalse {
				curWordInfo = false
				break
			} else if letterInfo[string(letter)].ImoossiblePlaces[letterI] {
				curWordInfo = false
				break
			} else if curWordStatus.WordStatus[letterI].Found && (string(letter) != curWordStatus.WordStatus[letterI].Letter) {
				curWordInfo = false
				break
			}
		}
		for _, yellowLetter := range yelloWLetters {
			includeThis := false
			for _, letter := range word.String {
				if yellowLetter == string(letter) {
					includeThis = true
					break
				}
			}
			if !includeThis {
				curWordInfo = false
				break
			}
		}
		if curWordInfo {
			newWords = append(newWords, word)
		}
	}
	wordsleft = newWords
}