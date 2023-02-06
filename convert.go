package main

import (
	"encoding/json"
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
	for verseKey, verseVal := range m {
		log.Printf("\nWorking on verse: %v", verseKey)
		_verseVal := verseVal.(map[string]any)
		for _, assignmentsVal := range _verseVal {
			__v := assignmentsVal.([]any)
			for _, wordsVal := range __v {
				// log.Printf("wordsKey,wordsVal=%v,%v", wordsKey, wordsVal)
				_alignedWordsKey := wordsVal.(map[string]any)
				for wordsKey, wordsVal := range _alignedWordsKey {
					// log.Printf("wordkey,wordval: %v, %v", wordKey, wordVal)
					if wordsKey == "topWords" {
						log.Printf("process topWords\n")
						_top := wordsVal.([]any)
						for topKey, topVal := range _top {
							log.Printf("....tw: key, val=%v,%v", topKey, topVal)
						}
					} else if wordsKey == "bottomWords" {
						log.Printf("process bottomWords\n")
						_bot := wordsVal.([]any)
						for botKey, botVal := range _bot {
							log.Printf("...bw: key, val=%v,%v", botKey, botVal)
						}
					}
				}
			}
		}
	}
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: convert -i input.csv -o output.csv\n")
	flag.PrintDefaults()
	log.Fatalln("Try again.")
}
