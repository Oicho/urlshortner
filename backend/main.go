package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type redirect struct {
	ID  string `json:"ID"`
	Url string `json:"Url"`
}

type allRedirects []redirect

var redirects = allRedirects{
	{
		ID:  "1",
		Url: "https://perdu.com",
	},
	{
		ID:  "2",
		Url: "https://www.lttstore.com/",
	},
}

func redirectFromID(w http.ResponseWriter, r *http.Request) {
	urlID := mux.Vars(r)["id"]
	for _, v := range redirects {
		if v.ID == urlID {
			http.Redirect(w, r, fmt.Sprintf(v.Url), http.StatusFound)
			return
		}
	}
	http.NotFound(w, r)
	// No match

}

func addRedirect(w http.ResponseWriter, r *http.Request) {
	var newRedirect redirect
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(requestBody, &newRedirect)
	redirects = append(redirects, newRedirect)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newRedirect)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{id}", redirectFromID).Methods("GET")
	router.HandleFunc("/url/", addRedirect).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
