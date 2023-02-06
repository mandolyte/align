package main

import (
	"encoding/csv"
	"encoding/json"
	"sort"
	"strconv"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	flag.Parse()

	if *output == "" {
		usage("No output file provided!")
	}
	if *input == "" {
		usage("No input file provided!")
	}

	// input file
	fi, fierr := os.Open(*input)
	if fierr != nil {
		log.Fatal("os.Open() Error:" + fierr.Error())
	}
	defer fi.Close()

	data, err := ioutil.ReadAll(fi)

	// var results []alignment
	var results any
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		log.Fatalln(err)
	}
	m := results.(map[string]any)

	var rows [][]string

	for verseKey, verseVal := range m {
		log.Printf("\nWorking on verse: %v", verseKey)
		_verseVal := verseVal.(map[string]any)
		for _, assignmentsVal := range _verseVal {
			__v := assignmentsVal.([]any)
			alignmentSequence := 0
			for _, wordsVal := range __v {
				// log.Printf("wordsKey,wordsVal=%v,%v", wordsKey, wordsVal)
				_alignedWordsKey := wordsVal.(map[string]any)
				alignmentSequence++
				for wordsKey, wordsVal := range _alignedWordsKey {
					// log.Printf("wordkey,wordval: %v, %v", wordKey, wordVal)

					if wordsKey == "topWords" {
						_top := wordsVal.([]any)
						wordSequence := 0
						for _, topVal := range _top {
							row := make([]string, 12)
							row[0] = "Titus"       // book
							row[1] = "1"           // chapter
							row[2] = "" + verseKey // verse
							row[3] = fmt.Sprintf("%v", alignmentSequence)

							row[4] = "T"
							wordSequence++
							row[5] = fmt.Sprintf("%v", wordSequence)
							_topMap := topVal.(map[string]any)
							occurrence := fmt.Sprintf("%v", _topMap["occurrence"])
							occurrences := fmt.Sprintf("%v", _topMap["occurrences"])
							word := _topMap["word"].(string)
							strong := _topMap["strong"].(string)
							lemma := _topMap["lemma"].(string)
							morph := _topMap["morph"].(string)
							row[6] = occurrence
							row[7] = occurrences
							row[8] = word
							row[9] = strong
							row[10] = lemma
							row[11] = morph
							rows = append(rows, row)
						}
					} else if wordsKey == "bottomWords" {
						_bot := wordsVal.([]any)
						wordSequence := 0
						for _, botVal := range _bot {
							row := make([]string, 12)
							row[0] = "Titus"       // book
							row[1] = "1"           // chapter
							row[2] = "" + verseKey // verse
							row[3] = fmt.Sprintf("%v", alignmentSequence)
							row[4] = "B"
							wordSequence++
							row[5] = fmt.Sprintf("%v", wordSequence)
							_botMap := botVal.(map[string]any)
							occurrence := fmt.Sprintf("%v", _botMap["occurrence"])
							occurrences := fmt.Sprintf("%v", _botMap["occurrences"])
							word := _botMap["word"].(string)
							row[6] = occurrence
							row[7] = occurrences
							row[8] = word
							row[9] = ""
							row[10] = ""
							row[11] = ""
							rows = append(rows, row)
						}
					}
				}
			}
		}
	}

	// sort the rows by:
	// book, chapter, verse, alignment sequence, type, and word sequence
	sort.Slice(rows, func(i, j int) bool {
		// book
		if rows[i][0] == rows[j][0] {
			// chapter
			ci, cerri := strconv.Atoi(rows[i][1])
			cj, cerrj := strconv.Atoi(rows[j][1])
			if cerri != nil {
				return true
			}
			if cerrj != nil {
				return true
			}

			if ci == cj {
				// verse
				vi, verri := strconv.Atoi(rows[i][2])
				vj, verrj := strconv.Atoi(rows[j][2])
				if verri != nil {
					return true
				}
				if verrj != nil {
					return true
				}
				if vi == vj {
					// align sequence
					ai, _ := strconv.Atoi(rows[i][3])
					aj, _ := strconv.Atoi(rows[j][3])

					if ai == aj {
						// type
						if rows[i][4] == rows[j][4] {
							// word sequence
							wi, werri := strconv.Atoi(rows[i][5])
							wj, werrj := strconv.Atoi(rows[j][5])
							if werri != nil {
								panic("Word sequence is not a number!" + fmt.Sprintf("%v", wi) + fmt.Sprintf("%v", rows[i]))
							}
							if werrj != nil {
								panic("Word sequence is not a number!" + fmt.Sprintf("%v", wj) + fmt.Sprintf("%v", rows[j]))
							}
							if wi == wj {
								panic("Word seqences are the same!" + fmt.Sprintf("%v and %v", wi, wj))
							} else if wi < wj {
								return true
							} else {
								return false
							}
						} else if rows[i][4] > rows[j][4] { // sort T before B
							return true
						} else {
							return false
						}
					} else if ai < aj {
						return true
					} else {
						return false
					}
				} else if vi < vj {
					return true
				} else {
					return false
				}

			} else if ci < cj {
				return true
			} else {
				return false
			}

		} else if rows[i][0] < rows[j][0] {
			return true
		} else {
			return false
		}
	})

	hdrs := []string{
		"Book", "Chapter",
		"Verse", "Alignment Sequence",
		"Type", "Word Sequence",
		"Occurence", "Occurences",
		"Word", "Strong",
		"Lemma", "Morph",
	}

	var sortedRows [][]string
	// add the headers first
	sortedRows = append(sortedRows, hdrs)
	sortedRows = append(sortedRows, rows...)

	// open output file
	var w *csv.Writer
	fo, foerr := os.Create(*output)
	if foerr != nil {
		log.Fatalln("os.Create() Error:" + foerr.Error())
	}
	defer fo.Close()
	w = csv.NewWriter(fo)
	w.WriteAll(sortedRows) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: convert -i input.json -o output.csv\n")
	flag.PrintDefaults()
	log.Fatalln("Try again.")
}
