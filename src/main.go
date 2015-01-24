package main

import (
	"fmt"
	"github.com/JensRantil/go-csv"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	FieldDatumFieldCount = 14
)

type FieldMismatchError struct {
	expected, found int
}

func (e *FieldMismatchError) Error() string {
	return "String array field count mismatch. Expected " +
		strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type UnsupportedTypeError struct {
	Type string
}

func (e *UnsupportedTypeError) Error() string {
	return "Unsupported type: " + e.Type
}

var NutDb = make(map[int32]*FoodDatum)

type FoodDatum struct {
	// 5-digit Nutrient Databank number that uniquely identifies a food item. If
	// this field is defined as numeric, the leading zero will be lost.
	Id int32 // %5d

	// 4-digit code indicating food group to which a food item belongs.
	FoodGroupId int32 // %4d

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
	IsSurvey bool // %b Nil

	// Description of inedible parts of a food item (refuse), such as seeds or
	// bone.
	RefuseDescription string // %135s Nil

	// Percentage of refuse.
	Refuse int32 // %2d Nil

	// Scientific name of the food item. Given for the least processed form of
	// the food (usually raw), if applicable.
	ScientificName string // %65s Nil

	// Factor for converting nitrogen to protein.
	NitrogenFactor float32 // %4.2f Nil

	// Factor for calculating calories from protein.
	ProteinFactor float32 // %4.2f Nil

	// Factor for calculating calories from fat (see p. 13).
	FatFactor float32 // %4.2f Nil

	// Factor for calculating calories from carbohydrates
	CarbohydrateFactor float32 // %4.2f Nil
}

func fieldToInt32(s string) int32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Fatalf("Error: Non integer value <%s>", s)
	}
	return int32(i64)
}

func fieldToFloat32(s string) float32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Fatalf("Error: Non float value <%s>", s)
	}
	return float32(f64)
}

func unwrapFoodDatum(record []string, fd *FoodDatum) (err error) {
	if len(record) != FieldDatumFieldCount {
		return &FieldMismatchError{FieldDatumFieldCount, len(record)}
	}
	fd.Id = fieldToInt32(record[0])
	fd.FoodGroupId = fieldToInt32(record[1])
	fd.Description = record[2]
	fd.BriefDescription = record[3]
	fd.CommonName = record[4]
	fd.ManufacturerName = record[5]
	fd.IsSurvey = (record[6] == "Y")
	fd.RefuseDescription = record[7]
	fd.Refuse = fieldToInt32(record[8])
	fd.ScientificName = record[9]
	fd.NitrogenFactor = fieldToFloat32(record[10])
	fd.ProteinFactor = fieldToFloat32(record[11])
	fd.FatFactor = fieldToFloat32(record[12])
	fd.CarbohydrateFactor = fieldToFloat32(record[13])
	return
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
		fd := new(FoodDatum)
		err = unwrapFoodDatum(l, fd)
		if err != nil {
			fmt.Println("Error unwrapping Food Datum %v", err)
		}
		NutDb[fd.Id] = fd
		if i > 10000 {
			break
		}
	}
	fmt.Printf("NutDb (%T) %d\n", NutDb, len(NutDb))
}
