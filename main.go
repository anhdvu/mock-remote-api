package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anhdvu/mock_remote_api/wtools"
)

func procReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 NOT FOUND\nPLEASE USE POST METHOD", http.StatusNotFound)
	} else {
		log.Println("######## INFO: New request received ########")
		fmt.Println("######## INFO: New request received ########")

		// Parse request header
		fmt.Println("\n******** Request header ********")
		wtools.ParseRemoteRequestHeader(r)

		// Parse request body
		payload := wtools.ParseRemoteRequestBody(r) // Parse XML in request body to a payload struct
		reqObj := wtools.ParseMethod(payload)       // From the general payload struct to a method-specific struct
		fmt.Println("\n******** Request body in JSON ********")
		wtools.JSONize(reqObj) // Marshal a method-specific struct to JSON and dump to os.Stdout
		fmt.Println("\n******** KLV breakdown ********")
		wtools.KLVSplitter(reqObj.KLV())
		fmt.Println("\n******** String to hash ********")
		fmt.Println(reqObj.String())

		// Handling response to HTTP request
		w.Header().Set("Content-Type", "text/xml") // Set response header
		fmt.Println("\n******** Response ********")
		switch payload.MethodName {
		case "AdministrativeMessage", "Stop":
			w.Write(wtools.GenerateResponseCodeOnly("1"))
		default:
			w.Write(wtools.GenerateResponse("1", "Approved"))
		}
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\nWelcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/code1", procReq)

	fs := http.FileServer(http.Dir("log"))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	logFile := "log/logs.txt"
	mux.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, logFile)
	})
	log.Println("Wallet (Companion Remote API) v0.1 is listening on port 8888")
	log.Fatal(http.ListenAndServe(":8888", mux))
}
