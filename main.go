package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/anhdvu/mock_remote_api/walletutils"
)

func processRemoteAPI(w http.ResponseWriter, r *http.Request) {
	responseCode1 := `<methodResponse><params><param><value><struct><member><name>resultCode</name><value><string>1</string></value></member></struct></value></param></params></methodResponse>`

	responseCodeMinus9 := `<methodResponse><params><param><value><struct><member><name>resultCode</name><value><string>-9</string></value></member></struct></value></param></params></methodResponse>`

	if r.Method != "POST" {
		http.Error(w, "404 NOT FOUND\nPLEASE USE POST METHOD", http.StatusNotFound)
	} else {
		codePath := r.URL.Path[(len(r.URL.Path) - 1):]
		switch codePath {
		case "1":
			fmt.Fprintf(w, responseCode1)
		case "9":
			fmt.Fprintf(w, responseCodeMinus9)
		default:
			fmt.Fprintf(w, "Couldn't find corresponding response code.")
		}
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\n Welcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/code1", processRemoteAPI)

	mux.HandleFunc("/code-9", processRemoteAPI)

	log.Fatal(http.ListenAndServe(":443", mux))
}
