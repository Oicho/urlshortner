package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type redirectReceived struct {
	Url string `json:"Url"`
	Hit uint   `json:"Hit"`
}

type redirect struct {
	Url string `json:"Url"`
	Hit uint   `json:"Hit"`
}

type allRedirects []redirect

var ctx = context.Background()

/*
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
*/
func addRedirect(w http.ResponseWriter, r *http.Request) {
	var receivedJson redirectReceived
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(requestBody, &receivedJson)

	newRedirect := redirect{receivedJson.Url, 0}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ret := rdb.Set("4", newRedirect, 0)
	err = ret.Err()
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newRedirect)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	//router.HandleFunc("/{id}", redirectFromID).Methods("GET")
	router.HandleFunc("/url/", addRedirect).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
