package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func serve(w http.ResponseWriter, r *http.Request) {
	// build the file path of the template we require based on the request URL
	var buffer bytes.Buffer

	buffer.WriteString("/go/src/github.com/paulcarmichael/fincrypt/html")

	if r.URL.String() == "/" {
		buffer.WriteString("/index")
	} else {
		buffer.WriteString(r.URL.String())
	}

	buffer.WriteString(".html")

	// recover and parse the template file
	file, err := filepath.Abs(buffer.String())

	if err != nil {
		log.Fatal("Abs: ", err)
	}

	t, err := template.ParseFiles(file)

	// if the template wasn't found then prepare the 404 page
	if err != nil {
		log.Println("ParseFiles: ", err)

		file, err = filepath.Abs("/go/src/github.com/paulcarmichael/fincrypt/html/404.html")

		if err != nil {
			log.Fatal("Abs: ", err)
		}

		t, err = template.ParseFiles(file)

		if err != nil {
			log.Fatal("ParseFiles: ", err)
		}

		w.WriteHeader(http.StatusNotFound)
	}

	t.Execute(w, nil)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer

	buffer.WriteString("/go/src/github.com/paulcarmichael/fincrypt/")

	if r.URL.String() == "/favicon.ico" {
		buffer.WriteString("images/favicon.ico")
	} else if r.URL.String() == "/style.css" {
		buffer.WriteString("css/style.css")
	} else { // serve scripts
		buffer.WriteString(r.URL.String())
	}

	file, err := filepath.Abs(buffer.String())

	if err != nil {
		log.Fatal("Abs: ", err)
	}

	// check that the file exists
	_, err = os.Stat(file)

	if err != nil {
		log.Fatal("Stat: ", err)
	}

	http.ServeFile(w, r, file)
}

func main() {
	// start the logger
	logFile, err := os.OpenFile("/go/logs/fincrypt.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		666)

	if err != nil {
		log.Fatal("Failed to open log file, exiting...")
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("fincrypt started...")

	// register the handler functions
	http.HandleFunc("/", serve)
	http.HandleFunc("/style.css", fileHandler)
	http.HandleFunc("/favicon.ico", fileHandler)
	http.HandleFunc("/scripts/", fileHandler)

	// serve up some web pages
	err = http.ListenAndServe(":80", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}