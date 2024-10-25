package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct{
	ID int              `json:"id"`
	OriginalURL string   `json:"original_url"`
	ShortURL string       `json:"short_url"`
	CreatedAt time.Time    `json:"created_date"`
}

/*
   d973671 ---> {
     ID: 1,
	OriginalURL: "https://www.google.com",
	ShortURL: "http://localhost:8080/d973671",
	CreatedAt: time.Now(),
   }
*/
var urlDB = make(map[string]URL)

func generateShortURL(originalURL string) string{
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hashBytes :=  hasher.Sum(nil)
	fmt.Println("hashBytes",hashBytes)
	shortURL := hex.EncodeToString(hashBytes)[:8]
	fmt.Println("shortURL",shortURL)
	fmt.Println("hasher",hasher)
	return shortURL
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	// Removed conversion to int
	urlDB[shortURL] = URL{
		ID: len(urlDB) + 1, // Use a simple incrementing ID
		OriginalURL: OriginalURL,
		ShortURL: shortURL,
		CreatedAt: time.Now(),
	}
	return shortURL
}
func getURL(shortURL string) (URL,error){
	fmt.Println("urlDB",urlDB)
	url,ok := urlDB[shortURL]
	if !ok{
		return URL{},errors.New("URL not found")
	}
	return url,nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request){
	// fmt.Println("handler called")
	fmt.Fprintf(w,"Hello World")
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Shorten URL Handler")
	var data struct{
		URL string `json:"url"`
	}
	err :=json.NewDecoder(r.Body).Decode(&data)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	json.NewEncoder(w).Encode(map[string]string{"short_url":shortURL})
}

func RedirectToOriginalURL(w http.ResponseWriter, r *http.Request){
	fmt.Println("called")
	id := r.URL.Path[len("/redirect/"):]
	url ,err := getURL(id)
	if err != nil {
		http.Error(w,"Invalid requiest", http.StatusNotFound)
		return
	}	
	http.Redirect(w,r,url.OriginalURL,http.StatusFound)

}
func main() {
	fmt.Println("Starting URL Shortener...")
	// shortURL := createURL("https://www.google.com")

	http.HandleFunc("/",RootPageURL)
    http.HandleFunc("/shorten",ShortenURLHandler)
	http.HandleFunc("/redirect",RedirectToOriginalURL)

	// Convert urlDB to JSON and print it
	// jsonData, err := json.MarshalIndent(urlDB, "", "   ")
	// if err != nil {
	// 	fmt.Println("Error marshaling to JSON:", err)
	// } else {
	// 	fmt.Println("urlDB in JSON format:")
	// 	fmt.Println(string(jsonData))
	// }

	// url, err := getURL(shortURL)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("URL:", url)
	// }

	fmt.Println("starting server at port 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
