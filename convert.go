package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type TopWord struct {
	Strong      string
	Lemma       string
	Morph       string
	Occurrence  string
	Occurrences string
	Word        string
}

type BottomWord struct {
	Occurrence  string
	Occurrences string
	Word        string
}

type alignment struct {
	TopWords    []TopWord
	BottomWords []BottomWord
}

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

	var results []alignment

	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		log.Fatalln(err)
	}

	hdrs := []string{
		"Book", "Chapter",
		"Verse", "Alignment Sequence",
		"Type", "Word Sequence",
		"Occurence", "Occurences",
		"Word", "Strong",
		"Lemma", "Morph",
	}

	var rows [][]string
	// add the headers first
	rows = append(rows, hdrs)

	for i := 0; i < len(results); i++ {

		// process the slice of Top Words
		for t := 0; t < len(results[i].TopWords); t++ {
			row := make([]string, 12)
			row[0] = "Titus"
			row[1] = "1"
			row[2] = "1"
			row[3] = fmt.Sprintf("%v", i+1)
			row[4] = "T"
			row[5] = fmt.Sprintf("%v", t+1)
			row[6] = results[i].TopWords[t].Occurrence
			row[7] = results[i].TopWords[t].Occurrences
			row[8] = results[i].TopWords[t].Word
			row[9] = results[i].TopWords[t].Strong
			row[10] = results[i].TopWords[t].Lemma
			row[11] = results[i].TopWords[t].Morph

			// add it to the rows
			rows = append(rows, row)
		}

		// process the slice of Bottom Words
		for t := 0; t < len(results[i].BottomWords); t++ {
			row := make([]string, 12)
			row[0] = "Titus"
			row[1] = "1"
			row[2] = "1"
			row[3] = fmt.Sprintf("%v", i+1)
			row[4] = "B"
			row[5] = fmt.Sprintf("%v", t+1)
			row[6] = results[i].BottomWords[t].Occurrence
			row[7] = results[i].BottomWords[t].Occurrences
			row[8] = results[i].BottomWords[t].Word
			row[9] = ""
			row[10] = ""
			row[11] = ""

			// add it to the rows
			rows = append(rows, row)
		}
	}

	// fmt.Println(rows)

	// open output file
	var w *csv.Writer
	fo, foerr := os.Create(*output)
	if foerr != nil {
		log.Fatalln("os.Create() Error:" + foerr.Error())
	}
	defer fo.Close()
	w = csv.NewWriter(fo)
	w.WriteAll(rows) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: convert -i input.csv -o output.csv\n")
	flag.PrintDefaults()
	log.Fatalln("Try again.")
}
