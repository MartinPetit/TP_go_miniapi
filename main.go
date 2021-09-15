package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

//Get ("/") ==> Print current hours and minutes
func hour(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		currentTime := time.Now()
		hours := currentTime.Hour()
		minutes := currentTime.Minute()
		fmt.Fprintln(w, hours, "h", minutes)
	} else {
		fmt.Fprintln(w, "Une erreur est survenue")
	}
}

//Post ("/add") ==> Add entry and author to file
func add(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintf(w, "Something went bad")
			return
		}

		file, err := os.OpenFile("data.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		defer file.Close()

		if err != nil {
			panic(err)
		}

		author := req.Form.Get("author")
		entry := req.Form.Get("entry")

		if author != "" && entry != "" {
			file.WriteString(author + ":" + entry + "\n")
		} else {
			fmt.Fprintf(w, "You must write author and entry")
		}

	} else {
		fmt.Fprintf(w,"Error has occured")
	}
}

//Get("/entries") ==> Get all entries
func entries(w http.ResponseWriter, req *http.Request)  {
	// On ouvre le fichier et on vérifie s'il est non null
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Fprintln(w, "Une erreur est survenue")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(strings.Split(scanner.Text(), ":")) > 1 {
			fmt.Fprintf(w, "%s\n", strings.Split(scanner.Text(), ":")[1])
		}
	}

}


func main() {
	http.HandleFunc("/", hour)
	http.HandleFunc("/add", add)
	http.HandleFunc("/entries", entries)
	http.ListenAndServe(":4567", nil)
}