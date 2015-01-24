package main

import (
	"fmt"
	"github.com/JensRantil/go-csv"
	"log"
	"os"
	"reflect"
)

type FoodDatum struct {
	// 5-digit Nutrient Databank number that uniquely identifies a food item. If
	// this field is defined as numeric, the leading zero will be lost.
	Id int // %5d

	// 4-digit code indicating food group to which a food item belongs.
	FoodGroupId int // %4d

	// 200-character description of food item.
	Description string // %200s

	// 60-character abbreviated description of food item.  Generated from the
	// 200-character description using abbreviations in Appendix A. If short
	// description is longer than 60 characters, additional abbreviations are
	// made.
	BriefDescription string // %60s

	// Other names commonly used to describe a food, including local or regional
	// names for various foods, for example, “soda” or “pop” for “carbonated
	// beverages.”
	CommonName string // %100s Nil

	// Indicates the company that manufactured the product, when appropriate.
	ManufacturerName string // %65s Nil

	// Indicates if the food item is used in the USDA Food and Nutrient Database
	// for Dietary Studies (FNDDS) and thus has a complete nutrient profile for
	// the 65 FNDDS nutrients.
	Survey bool // %b Nil

	// Description of inedible parts of a food item (refuse), such as seeds or
	// bone.
	RefuseDescription string // %135s Nil

	// Percentage of refuse.
	Refuse int // %2d Nil

	// Scientific name of the food item. Given for the least processed form of
	// the food (usually raw), if applicable.
	SciName int // %65s Nil

	// Factor for converting nitrogen to protein.
	NitrogenFactor float32 // %4.2f Nil

	// Factor for calculating calories from protein.
	ProteinFactor float32 // %4.2f Nil

	// Factor for calculating calories from fat (see p. 13).
	FatFactor // %4.2f Nil

	// Factor for calculating calories from carbohydrates
	CHO_Factor float32 // %4.2f Nil
}

func UnmarshalFoodDatum(strings []string, v interface{}) error {
	s := reflect

}

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
