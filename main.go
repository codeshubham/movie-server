package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(movies)
}

func getMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rw.Header().Set("Content-Type", "application/json")
	for _, movie := range movies {
		if movie.ID == id {
			json.NewEncoder(rw).Encode(movie)
			return
		}
	}
}

func deleteMovie(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rw.Header().Set("Content-Type", "application/json")
	for idx, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:idx], movies[idx+1:]...)
			break
		}
	}
}

func createMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000000))
	movies = append(movies, movie)
}

func updateMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	for idx, curr := range movies {
		if curr.ID == id {
			dc := json.NewDecoder(r.Body)
			var mv Movie
			_ = dc.Decode(&mv)
			movies = append(movies[:idx], movies[idx+1:]...)
			mv.ID = strconv.Itoa(rand.Intn(100000000))
			movies = append(movies, mv)
			break
		}
	}
	//http.Error(rw, "Element to be updated not found", http.StatusBadRequest)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "191338", Title: "Movie One", Director: &Director{Firstname: "Shubham", Lastname: "Kumar"}})
	movies = append(movies, Movie{ID: "2", Isbn: "495677", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at post 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
