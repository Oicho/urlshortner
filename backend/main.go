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
	"github.com/rs/xid"
)

var rdb = redis.NewClient(&redis.Options{
	PoolSize: 60,
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

type RedirectReceived struct {
	Url string `json:"Url"`
}

type Redirect struct {
	Url string `json:"Url"`
	Hit uint   `json:"Hit"`
	Id  string `json:"Id"`
}

var ctx = context.Background()

func redirectFromID(w http.ResponseWriter, r *http.Request) {
	var redirect Redirect
	urlID := mux.Vars(r)["id"]

	status := rdb.Get(urlID)
	err := status.Err()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	json.Unmarshal([]byte(status.Val()), &redirect)

	redirect.Hit += 1
	jsonOut, err := json.Marshal(redirect)
	if err != nil {
		fmt.Println(err)
	}
	if rdb.Set(redirect.Id, jsonOut, 0).Err() != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf(redirect.Url), http.StatusFound)
}

func getRedirect(w http.ResponseWriter, r *http.Request) {
	urlID := mux.Vars(r)["id"]

	status := rdb.Get(urlID)
	err := status.Err()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusAccepted)

	w.Write([]byte(status.Val()))
}

func addRedirect(w http.ResponseWriter, r *http.Request) {
	var receivedJson RedirectReceived
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	err = json.Unmarshal(requestBody, &receivedJson)
	if err != nil || receivedJson.Url == "" {
		panic(err)
	}
	/* This should be:
	a) a bijective function https://stackoverflow.com/questions/742013/how-do-i-create-a-url-shortener
	b) checked for duplicate key
	*/
	rdb.RandomKey()
	key := xid.New().String()
	newRedirect := Redirect{receivedJson.Url, 0, key}

	jsonOut, err := json.Marshal(newRedirect)
	if err != nil {
		fmt.Println(err)
	}
	if err := rdb.Set(key, jsonOut, 0).Err(); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newRedirect)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/{id}", redirectFromID).Methods("GET")

	router.HandleFunc("/url/", addRedirect).Methods("POST")
	router.HandleFunc("/url/{id}", getRedirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
