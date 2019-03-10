package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paulcarmichael/fincrypt"
	"golang.org/x/crypto/acme/autocert"
)

type operationResponse struct {
	Res string
	Err string
}

const htmlPath = "/go/src/github.com/paulcarmichael/fincrypt.io/html"
const linksPath = htmlPath + "/links.html"
const navbarPath = htmlPath + "/navbar.html"
const sidebarPath = htmlPath + "/sidebar.html"
const errorPath = htmlPath + "/error.html"
const scriptsPath = htmlPath + "/scripts.html"
const fourOhFourPath = htmlPath + "/404.html"

func serve(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// if the request method is get, it's a request for a webpage
		var buffer bytes.Buffer
		buffer.WriteString(htmlPath)

		if r.URL.String() == "/" {
			buffer.WriteString("/index")
		} else {
			buffer.WriteString(r.URL.String())
		}
		buffer.WriteString(".html")

		// recover and parse the template file
		file, err := filepath.Abs(buffer.String())

		if err != nil {
			log.Println("Abs: ", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		t, err := template.ParseFiles(file, linksPath, navbarPath, sidebarPath, errorPath, scriptsPath)

		// if the template wasn't found then prepare the 404 page
		if err != nil {
			log.Println("ParseFiles: ", err)

			file, err = filepath.Abs(fourOhFourPath)

			if err != nil {
				log.Println("Abs: ", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			t, err = template.ParseFiles(file, linksPath, navbarPath, sidebarPath)

			if err != nil {
				log.Println("ParseFiles: ", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}

		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// if the request method is post, it's a request for a calculation
		var op fincrypt.Operation
		var response operationResponse

		url := r.URL.String()

		if url == "/aes" {
			op = &fincrypt.AESOperation{}
		} else if url == "/base64" {
			op = &fincrypt.Base64Operation{}
		} else if url == "/bertlv" {
			op = &fincrypt.BERTLVParser{}
		} else if url == "/cvv" {
			op = &fincrypt.CVVOperation{}
		} else if url == "/des" {
			op = &fincrypt.DESOperation{}
		} else if url == "/hash" {
			op = &fincrypt.HashOperation{}
		} else if url == "/hmac" {
			op = &fincrypt.HMACOperation{}
		} else if url == "/luhn" {
			op = &fincrypt.LuhnOperation{}
		} else if url == "/pinblock" {
			op = &fincrypt.PINBlockOperation{}
		} else if url == "/pinoffset" {
			op = &fincrypt.PINOffsetOperation{}
		} else if url == "/pvv" {
			op = &fincrypt.PVVOperation{}
		} else if url == "/retailmac" {
			op = &fincrypt.RetailMACOperation{}
		} else if url == "/rsa" {
			op = &fincrypt.RSAOperation{}
		} else if url == "/tagsearch" {
			op = &fincrypt.BERTLVTag{}
		} else if url == "/xor" {
			op = &fincrypt.XOROperation{}
		} else {
			log.Println("Unknown URL:", url)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// unpack the json payload into the operation struct
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println("ReadAll", err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		log.Println(http.MethodPost, url, string(body))

		err = json.Unmarshal(body, &op)

		if err != nil {
			log.Println("Unmarshal", err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		// calculate the result of the operation
		response.Res, err = op.Calculate()

		if err != nil {
			response.Err = err.Error()
		}

		// build the response
		if len(response.Err) > 0 {
			w.WriteHeader(http.StatusNotAcceptable)
		}

		responseBytes, err := json.Marshal(response)

		if err != nil {
			log.Println("Marshal", err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
		}

		log.Println(http.MethodPost, url, string(responseBytes))

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBytes)
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer

	buffer.WriteString("/go/src/github.com/paulcarmichael/fincrypt.io/")

	if r.URL.String() == "/favicon.ico" {
		buffer.WriteString("images/favicon.ico")
	} else if r.URL.String() == "/style.css" {
		buffer.WriteString("css/style.css")
	} else { // serve scripts
		buffer.WriteString(r.URL.String())
	}

	file, err := filepath.Abs(buffer.String())

	if err != nil {
		log.Println("Abs: ", err)
		return
	}

	// check that the file exists
	_, err = os.Stat(file)

	if err != nil {
		log.Println("Stat: ", err)
		return

	}

	http.ServeFile(w, r, file)
}

func main() {
	// start the diagnostic logger
	logFile, err := os.OpenFile("/go/logs/fincrypt.io.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		666)

	if err != nil {
		log.Fatal("Failed to open log file, exiting...")
	}

	defer logFile.Close()
	log.SetOutput(logFile)

	// check for the production flag
	productionFlag := flag.Bool("production", false, "Starts the http server with HTTPS")
	flag.Parse()

	if *productionFlag == true {
		log.Println("fincrypt.io starting in production mode...")

		mux := http.NewServeMux()
		mux.HandleFunc("/", serve)
		mux.HandleFunc("/style.css", fileHandler)
		mux.HandleFunc("/favicon.ico", fileHandler)
		mux.HandleFunc("/scripts/", fileHandler)

		cm := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("certs"),
			HostPolicy: autocert.HostWhitelist("www.fincrypt.io"),
		}

		server := http.Server{
			Addr:      ":443",
			Handler:   mux,
			TLSConfig: cm.TLSConfig(),
		}

		// listening with http only on port 80 allows autocert to handle
		// http-01 challenges and the parameter acts as a http to https
		// redirect for other traffic
		go http.ListenAndServe(":80", cm.HTTPHandler(nil))

		err = server.ListenAndServeTLS("", "")

		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	} else {
		log.Println("fincrypt.io starting...")

		// register the handler functions
		http.HandleFunc("/", serve)
		http.HandleFunc("/style.css", fileHandler)
		http.HandleFunc("/favicon.ico", fileHandler)
		http.HandleFunc("/scripts/", fileHandler)

		// serve up some cryptography
		err = http.ListenAndServe(":80", nil)

		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
