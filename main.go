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
		fmt.Printf("######## NOTICE: New request received ########\n")
		w.Header().Set("Content-Type", "text/xml")  // Set response header
		payload := wtools.ParseRemoteRequestBody(r) // Parse XML in request body to struct
		parsedReq := wtools.ParseMethod(payload)    // From general struct to method-specific struct
		wtools.DumpJSON(parsedReq)                  // Marshal method-specific struct to JSON and dump to os.Stdout
		wtools.KLVSplitter(parsedReq.KLV())
		fmt.Printf("\n******** String to hash ********\n%q\n", parsedReq.String())

		switch payload.MethodName {
		case "AdministrativeMessage", "Stop":
			w.Write(wtools.GenerateResponseCodeOnly("1"))
		default:
			w.Write(wtools.GenerateResponse("1", "Approved"))
		}
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n\n\n")
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
