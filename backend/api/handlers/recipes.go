package handlers

import (
	"api/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleRecipes(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		typeValue := queryValues.Get("type")

		switch r.Method {
		case http.MethodGet:
			switch typeValue {
			case "all":
				getAllRecipes(w, r, db)
			case "":
				getRecipeById(w, r, db)
			default:
				log.Print("no type for type " + typeValue)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		default:
			log.Print("no implementation for method " + r.Method)
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		}
	}
}

func getAllRecipes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var recipesFound []structs.Recipe
	results, err := db.Query("SELECT id, name, mealTime, information,time_it_takes, roasting_time,rest_time,rating,difficulty FROM recipes;")
	if err != nil {
		log.Println(err.Error())
		return
	}

	for results.Next() {
		var recipe structs.Recipe
		err = results.Scan(&recipe.Id, &recipe.Name, &recipe.MealTime,
			&recipe.Information, &recipe.TimeItTakes, &recipe.RoastingTime,
			&recipe.RestTime, &recipe.Rating, &recipe.Difficulty)

		if err != nil {
			log.Println(err.Error())
			return
		}
		recipesFound = append(recipesFound, recipe)
	}
	w.Header().Add("content-type", "application/json")
	jsonEncodedData, err := json.Marshal(recipesFound)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("failed to encode data: ", err.Error())
	}

	_, err = fmt.Fprint(w, string(jsonEncodedData))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("failed to print data: ", err.Error())
	}

}

func getRecipeById(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idInUrl := r.URL.Query().Get("id")
	fmt.Println("id" + idInUrl)
	result, err := db.Query("SELECT id, name, mealTime, information,time_it_takes, roasting_time,rest_time,rating,difficulty FROM recipes WHERE id= ?;", idInUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var recipe structs.Recipe
	for result.Next() {
		err = result.Scan(&recipe.Id, &recipe.Name, &recipe.MealTime,
			&recipe.Information, &recipe.TimeItTakes, &recipe.RoastingTime,
			&recipe.RestTime, &recipe.Rating, &recipe.Difficulty)

		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	w.Header().Add("content-type", "application/json")
	jsonEncodedData, err := json.Marshal(recipe)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("failed to encode data: ", err.Error())
	}

	_, err = fmt.Fprint(w, string(jsonEncodedData))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("failed to print data: ", err.Error())
	}

}
