package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anhdvu/mock-remote-api/wtools"
)

func procReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	} else {
		fmt.Println("######## INFO: New request received ########")

		// Parse request header
		fmt.Println("\n******** Request header ********")
		wtools.ParseRemoteRequestHeader(r)

		// Parse request body
		payload := wtools.ParseRemoteRequestBody(r) // Parse XML in request body to a payload struct
		reqObj := wtools.ParseMethod(payload)       // From the general payload struct to a method-specific struct
		fmt.Println("\n******** Request body in JSON ********")
		wtools.ToJSON(reqObj) // Marshal a method-specific struct to JSON and dump to os.Stdout
		fmt.Println("\n******** KLV breakdown ********")
		wtools.KLVSplitter(reqObj.KLV())
		fmt.Println("\n******** String to hash ********")
		fmt.Println(reqObj.String())

		// Handling response to HTTP request
		w.Header().Set("Content-Type", "text/xml") // Set response header
		fmt.Println("\n******** Response ********")
		switch payload.MethodName {
		case "AdministrativeMessage":
			if payload.Params[2].Value.StringParam == "digitization.activationmethods" {
				w.Write(wtools.GenerateResponsewActivationMethods())
			} else {
				w.Write(wtools.GenerateResponseCodeOnly())
			}
		case "Stop":
			w.Write(wtools.GenerateResponseCodeOnly())
		default:
			w.Write(wtools.GenerateResponseCodeOnly())
		}
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\nWelcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/wallet", procReq)
	mux.HandleFunc("/remote", wtools.LogRemoteMessage(wtools.HandleRemoteMessage()))

	fs := http.FileServer(http.Dir("./log/"))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	adm := http.FileServer(http.Dir("adm"))
	mux.Handle("/adm/", http.StripPrefix("/adm/", adm))

	logFile := "./log/logs.txt"
	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, logFile)
	})
	fmt.Println("Wallet v0.2.20210624 is running...")
	log.Println("Wallet (Companion + MPQR) v0.2.20210813-MPQR is listening on port 8888")
	log.Fatal(http.ListenAndServe(":8888", mux))
}
