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
	var results interface{}
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		log.Fatalln(err)
	}
	m := results.(map[string]interface{})
	for verseKey, verseVal := range m {
		log.Printf("\nWorking on verse: %v", verseKey)
		_verseVal := verseVal.(map[string]interface{})
		for _, assignmentsVal := range _verseVal {
			__v := assignmentsVal.([]interface{})
			for _, wordsVal := range __v {
				// log.Printf("wordsKey,wordsVal=%v,%v", wordsKey, wordsVal)
				_alignedWordsKey := wordsVal.(map[string]any)
				for wordKey, wordVal := range _alignedWordsKey {
					log.Printf("wordkey,wordval: %v, %v", wordKey, wordVal)
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
