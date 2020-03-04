package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anhdvu/mock_remote_api/utils"
)

func processRemoteAPI(w http.ResponseWriter, r *http.Request) {
	responseCode1 := `
<methodResponse>
  <params>
    <param>
      <value>
        <struct>
          <member>
            <name>resultCode</name>
            <value>
              <string>1</string>
            </value>
          </member>
        </struct>
      </value>
    </param>
  </params>
</methodResponse>`

	responseCodeMinus9 := `
<methodResponse>
  <params>
    <param>
      <value>
        <struct>
          <member>
            <name>resultCode</name>
            <value>
              <string>-9</string>
            </value>
          </member>
        </struct>
      </value>
    </param>
  </params>
</methodResponse>`

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
	// s := "0021543093400000053300412000000000700026047299049033440850025003ADJ25110Adjustment25304907025400255002560203"
	s := "00200"
	walletutils.KLVSplitter(s)

	log.Fatal(http.ListenAndServe(":80", nil))
}
