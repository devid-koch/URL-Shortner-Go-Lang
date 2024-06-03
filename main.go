package main

import (
	"fmt"
	"crypto/md5"
	"time"
	"encoding/hex"
	"errors"
	"net/http"
	"encoding/json"
)

type URL struct {
	ID string `json:"id"`
	OriginalURL string `jsonL:"original_url"`
	ShortURL string `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // It converts the original URL string to a byte slice
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8];
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	id := shortURL
    urlDB[id] = URL{
        ID: id,
        OriginalURL: OriginalURL,
        ShortURL: shortURL,
        CreationDate: time.Now(),
    }
	return shortURL
}

func getURL (id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Hello world!")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL_ := createURL(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{
        ShortURL: shortURL_,
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	id:= r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err!= nil {
        http.Error(w, "URL not found", http.StatusNotFound)
    }
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}


func main(){

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	//start the HTTP server
	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error on staring server: ", err)
	}

}