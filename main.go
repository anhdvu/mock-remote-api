package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anhdvu/mock_remote_api/walletutils"
)

func processRemoteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 NOT FOUND\nPLEASE USE POST METHOD", http.StatusNotFound)
	} else {

		fmt.Printf("######## NOTICE: New request received ########\n")
		walletutils.ParseRemoteRequestHeaders(r)
		requestKLV := walletutils.ParseRemoteRequestBody(r).GetKLV()
		walletutils.KLVSplitter(requestKLV)
		codePath := r.URL.Path[(len(r.URL.Path) - 1):]
		switch codePath {
		case "1":
			fmt.Fprintf(w, walletutils.GenerateResponse("1"))
		case "9":
			fmt.Fprintf(w, walletutils.GenerateResponse("-9"))
		default:
			fmt.Fprintf(w, "Couldn't find corresponding response code.")
		}
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n\n\n")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\n Welcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/code1", processRemoteAPI)

	mux.HandleFunc("/code-9", processRemoteAPI)

	fs := http.FileServer(http.Dir("log"))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	logFile := "log/logs.txt"
	mux.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, logFile)
	})

	log.Fatal(http.ListenAndServe(":8888", mux))
}
