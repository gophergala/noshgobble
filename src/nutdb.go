package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	foodFieldCount     = 14
	nutrientFieldCount = 6
)

var FoodDb = make(map[int32]*Food)
var NutrientDb = make(map[int32]*Nutrient)

type FieldMismatchError struct {
	expected, found int
}

func (e *FieldMismatchError) Error() string {
	return "String array field count mismatch. Expected " +
		strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
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

type Food struct {
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

func unwrapFood(record []string, fd *Food) (err error) {
	if len(record) != foodFieldCount {
		return &FieldMismatchError{foodFieldCount, len(record)}
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

func loadFood(record []string, i interface{}) {
	db := i.(map[int32]*Food)
	fd := new(Food)
	err := unwrapFood(record, fd)
	if err != nil {
		fmt.Println("Error unwrapping Food Datum %v", err)
	}
	db[fd.Id] = fd
}

type Nutrient struct {
	// Unique 3-digit identifier code for a nutrient.
	Id int32 // %3d

	// Units of measure (mg, g, μg, and so on).
	Units string // %7s

	// International Network of Food Data Systems (INFOODS) Tagnames. A unique
	// abbreviation for a nutrient/food component developed by INFOODS to aid in
	// the interchange of data.
	Tagname string // %20s Nil

	// Name of nutrient/food component.
	Description string // %60s

	// Number of decimal places to which a nutrient value is rounded.
	Precision int32 // %1d

	// Used to sort nutrient records in the same order as various reports
	// produced from SR.
	SortOrder int32 // %6d
}

func unwrapNutrient(record []string, n *Nutrient) (err error) {
	if len(record) != nutrientFieldCount {
		return &FieldMismatchError{nutrientFieldCount, len(record)}
	}
	n.Id = fieldToInt32(record[0])
	n.Units = record[1]
	n.Tagname = record[2]
	n.Description = record[3]
	n.Precision = fieldToInt32(record[4])
	n.SortOrder = fieldToInt32(record[5])
	return
}

func loadNutrient(record []string, i interface{}) {
	db := i.(map[int32]*Nutrient)
	n := new(Nutrient)
	err := unwrapNutrient(record, n)
	if err != nil {
		log.Fatalf("Error unwrapping Nutrient Datum %v", err)
	}
	db[n.Id] = n
}

func loadFile(filename string, db interface{}, loadRecord func(record []string, db interface{})) {
	datafile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening data file \"%s\": %v", filename, err)
	}
	defer datafile.Close()

	reader := csv.NewReader(datafile)
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading data file", err)
	}

	for _, r := range rawCSVdata {
		loadRecord(r, db)
	}
}

func InitNutritionDatabase() {
	loadFile("data/nutcsv/FOOD_DES.csv", FoodDb, loadFood)
	loadFile("data/nutcsv/NUTR_DEF.csv", NutrientDb, loadNutrient)
}
