package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/anhdvu/mock-remote-api/wtools"
)

func handleRemoteCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("######## INFO: New request received ########")
	fmt.Println("\n******** Request header ********")
	wtools.ParseRemoteRequestHeader(r)

	// Parse request body
	payload := wtools.ParseRemoteRequestBody(r) // Parse XML in request body to a payload struct
	reqObj := wtools.ParseMethod(payload)       // From the general payload struct to a method-specific struct
	fmt.Println("\n******** XML Request converted to JSON format ********")
	wtools.ToJSON(reqObj) // Marshal a method-specific struct to JSON and dump to os.Stdout
	fmt.Println("\n******** KLV breakdown ********")
	wtools.KLVSplitter(reqObj.KLV())
	fmt.Println("\n******** String to hash ********")
	fmt.Println(reqObj.String())

	// Handling response to HTTP request
	w.Header().Set("Content-Type", "text/xml") // Set response header
	fmt.Println("\n******** XML Response ********")
	switch payload.MethodName {
	case "AdministrativeMessage":
		if payload.Params[2].Value.StringParam == "digitization.activationmethods" {
			w.Write(wtools.RespwCVM(payload.Params[1].Value.StringParam))
		} else {
			w.Write(wtools.RespCodeOnly())
		}
	case "Stop":
		w.Write(wtools.RespCodeOnly())
	case "Balance":
		w.Write(wtools.RespwBalance())
	default:
		w.Write(wtools.RespCodeOnly())
	}
	fmt.Printf("\n######## INFO: Request parse completed ########\n\n")
}

func main() {
	terminal := flag.String("terminal", "0059479238", "Terminal")
	password := flag.String("password", "DA43A9C5F8", "Terminal password")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\nWelcome %s", r.URL.Path[1:])
	})

	mux.HandleFunc("/wallet", handleRemoteCall)
	mux.Handle("/remote", wtools.LogRemoteMessage(wtools.HandleRemoteMessage(*terminal, *password)))

	fs := http.FileServer(http.Dir("./log/"))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	adm := http.FileServer(http.Dir("./adm/"))
	mux.Handle("/adm/", http.StripPrefix("/adm/", adm))

	fmt.Println("Wallet v0.4.220101 is running...")
	log.Println("Wallet (Companion + MPQR) v0.4.220121-MPQR is listening on port 80")
	log.Fatal(http.ListenAndServe(":80", mux))
}
