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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey there!\n Welcome %s", r.URL.Path[1:])
	})

	http.HandleFunc("/code1", processRemoteAPI)

	http.HandleFunc("/code-9", processRemoteAPI)

	/* s := "0021543093400000053300612000000000700026047299049033440850025003ADJ25110Adjustment25304907025400255002560203"
	walletutils.KLVSplitter(s) */

	/* key := "3983250115"
	message := "CreateLinkedCard0029504320556789731580SniperInvoker3456344513273456789101115151520200211T14:59:43"
	fmt.Println(walletutils.CalculateSHA1Checksum(key, message))
	fmt.Println(walletutils.CalculateSHA256Checksum(key, message)) */

	log.Fatal(http.ListenAndServe(":80", nil))
}
