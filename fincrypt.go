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

	buffer.WriteString("../src/github.com/paulcarmichael/fincrypt/html")

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

		file, err = filepath.Abs("../src/github.com/paulcarmichael/fincrypt/html/404.html")

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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	file, err := filepath.Abs("../src/github.com/paulcarmichael/fincrypt/images/favicon.ico")

	if err != nil {
		log.Fatal("Abs: ", err)
	}

	// check that the favicon file exists
	_, err = os.Stat(file)

	if err != nil {
		log.Fatal("Stat: ", err)
	}

	http.ServeFile(w, r, file)
}

func main() {
	// start the logger
	logFile, err := os.OpenFile("fincrypt.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		666)

	if err != nil {
		log.Fatal("Failed to open log file, exiting...")
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("fincrypt starting...")

	// register the handler functions
	http.HandleFunc("/", serve)
	http.HandleFunc("/favicon.ico", faviconHandler)

	// serve up some web pages
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
