package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/paulcarmichael/cryptop"
)

func index(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("C:\\Users\\Paul\\go\\src\\github.com\\paulcarmichael\\fincrypt\\html\\index.html")

	t.Execute(w, nil)
}

type xorOperation struct {
	Input1 string
	Input2 string
	Result string
}

func xor(w http.ResponseWriter, r *http.Request) {

	// if the method is GET then we have a request for the xor page
	if r.Method == "GET" {
		t, _ := template.ParseFiles("C:\\Users\\Paul\\go\\src\\github.com\\paulcarmichael\\fincrypt\\html\\xor.html")

		t.Execute(w, nil)
	} else {
		// parse the input variables and store them in the result structure
		r.ParseForm()

		xorResult := xorOperation{
			Input1: strings.Join(r.Form["input1"], ""),
			Input2: strings.Join(r.Form["input2"], ""),
		}

		// pack the given inputs
		packedI1, err := cryptop.PackFromHex(xorResult.Input1)

		if err != nil {
			log.Print("xor failed to pack input1: ", xorResult.Input1)
			log.Print(err)
		}

		packedI2, err := cryptop.PackFromHex(xorResult.Input2)

		if err != nil {
			log.Print("xor failed to pack input2: ", xorResult.Input2)
			log.Print(err)
		}

		// xor
		result, err := cryptop.XOR(string(packedI1),
			string(packedI2))

		if err != nil {
			log.Print(err)
			log.Print("xor failed, see earlier log messages")
		}

		// expand the result to the correct form for the frontend
		xorResult.Result, err = cryptop.ExpandToHex(result)

		log.Print("xor result [" + xorResult.Result + "]")

		// generate the response
		t, _ := template.ParseFiles("C:\\Users\\Paul\\go\\src\\github.com\\paulcarmichael\\fincrypt\\html\\xor.html")

		t.Execute(w, xorResult)
	}
}

type aesOperation struct {
	Key    string
	Data   string
	IV     string
	Result string
}

func aes(w http.ResponseWriter, r *http.Request) {

	// if the method is GET then we have a request for the aes page
	if r.Method == "GET" {
		t, _ := template.ParseFiles("C:\\Users\\Paul\\go\\src\\github.com\\paulcarmichael\\fincrypt\\html\\aes.html")

		t.Execute(w, nil)
	} else {
		// parse the input variables and store them in the result structure
		r.ParseForm()

		aesResult := aesOperation{
			Key:  strings.Join(r.Form["key"], ""),
			Data: strings.Join(r.Form["data"], ""),
			IV:   strings.Join(r.Form["iv"], ""),
		}

		// pack the given inputs
		packedKey, err := cryptop.PackFromHex(aesResult.Key)

		if err != nil {
			log.Print("aes failed to pack key: ", aesResult.Key)
			log.Print(err)
		}

		packedData, err := cryptop.PackFromHex(aesResult.Data)

		if err != nil {
			log.Print("aes failed to pack data: ", aesResult.Data)
			log.Print(err)
		}

		packedIV, err := cryptop.PackFromHex(aesResult.IV)

		if err != nil {
			log.Print("aes failed to pack iv: ", aesResult.IV)
			log.Print(err)
		}

		// prepare the operation mode
		mode := cryptop.Encrypt

		if r.Form.Get("mode") == "decrypt" {
			mode = cryptop.Decrypt
		}

		// aes
		result, err := cryptop.AES_CBC(packedKey,
			packedData,
			packedIV,
			mode)

		if err != nil {
			log.Print(err)
			log.Print("aes failed, see earlier log messages")
		}

		// expand the result to the correct form for the frontend
		aesResult.Result, err = cryptop.ExpandToHex(result)

		log.Print("aes result [" + aesResult.Result + "]")

		// generate the response
		t, _ := template.ParseFiles("C:\\Users\\Paul\\go\\src\\github.com\\paulcarmichael\\fincrypt\\html\\aes.html")

		t.Execute(w, aesResult)
	}
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
	http.HandleFunc("/", index)
	http.HandleFunc("/xor", xor)
	http.HandleFunc("/aes", aes)

	// go to work!
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
