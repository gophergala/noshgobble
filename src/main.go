package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	foodId, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving foodId from path %s\n", r.URL.Path)
	} else {
		if foodId < 0 || int(foodId) >= len(FoodDb) {
			fmt.Fprintf(w, "Error retrieving food. FoodId %d does not exist!", foodId)
			fmt.Fprintf(w, "FoodId must be between 0 and %d!", len(FoodDb))
		} else {
			food := FoodDb[int32(foodId)]
			fmt.Fprintf(w, "Hi there, I love %+v!", food)
		}
	}
}

func main() {
	InitializeFoodDb()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
