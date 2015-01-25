package main

import (
	"davebalmain.com/noshgobble/src/nutdb"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	foodId, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving foodId from path %s\n", r.URL.Path)
	} else {
		if foodId < 0 || int(foodId) >= len(nutdb.FoodDb) {
			fmt.Fprintf(w, "Error retrieving food. FoodId %d does not exist!", foodId)
			fmt.Fprintf(w, "FoodId must be between 0 and %d!", len(nutdb.FoodDb))
		} else {
			food := nutdb.FoodDb[int32(foodId)]
			d := struct {
				Description string
			}{
				food.Description,
			}
			t, _ := template.ParseFiles("templates/view.html")
			t.Execute(w, d)
		}
	}
}

func main() {
	nutdb.InitializeFoodDb()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	fmt.Printf("Now listening at http://localhost:8080/")
}
