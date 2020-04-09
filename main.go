package main

import (
	"fmt"
	_ "io"
	"log"
	"net/http"

	"github.com/anhdvu/mock_remote_api/walletutils"
)

func processRemoteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 NOT FOUND\nPLEASE USE POST METHOD", http.StatusNotFound)
	} else {
		fmt.Printf("######## NOTICE: New request received ########\n")
		w.Header().Set("Content-Type", "text/xml")       // Set response header
		payload := walletutils.ParseRemoteRequestBody(r) // Parse XML in request body to struct
		parsedReq := walletutils.ParseMethod(payload)    // From general struct to method-specific struct
		walletutils.DumpJSON(parsedReq)                  // Marshal method-specific struct to JSON and dump to os.Stdout
		walletutils.KLVSplitter(parsedReq.GetKLV())

		if payload.MethodName == "AdministrativeMessage" {
			w.Write(walletutils.GenerateResponse("0", "Message has been received"))
		} else {
			w.Write(walletutils.GenerateResponse("1", "Approved"))
		}
	}
	fmt.Printf("\n######## INFO: Request parse completed ########\n\n\n\n")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\n Welcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/code1", processRemoteAPI)

	fs := http.FileServer(http.Dir("log"))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	logFile := "log/logs.txt"
	mux.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, logFile)
	})

	log.Fatal(http.ListenAndServe(":8888", mux))
}
