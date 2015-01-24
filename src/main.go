package main

import (
	"fmt"
	"github.com/JensRantil/go-csv"
	"log"
	"os"
)

func main() {

	datafile, err := os.Open("data/usda/FOOD_DES.txt")

	if err != nil {
		log.Fatalf("Error opening data file: %v", err)
	}

	defer datafile.Close()

	reader := csv.NewDialectReader(datafile, csv.Dialect{Delimiter: '^', QuoteChar: '~'})

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		log.Fatalf("Error reading data file", err)
	}

	// sanity check, display to standard output
	for i, l := range rawCSVdata {
		fmt.Printf("line %d: len: %d => (%T) %v\n", i, len(l), l, l)
		if i > 10 {
			break
		}
	}
}
