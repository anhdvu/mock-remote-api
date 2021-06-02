package wtools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RemoteMessage struct {
	MessageType       string `json:"messageType"`
	TrackingNumber    string `json:"trackingNumber"`
	CustomerReference string `json:"customerReference"`
	WalletIdentifier  string `json:"walletIdentifier"`
	Challenge         string `json:"challenge"`
	ResultCode        string `json:"resultCode,omitempty"`
}

func HandleRemoteMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ParseRemoteRequestHeader(r)

		raw, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		pl := &RemoteMessage{}
		err = json.Unmarshal(raw, pl)
		if err != nil {
			fmt.Printf("Couldn't unmarshal json payload, %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println("\n******** Request body ********")
		pl.JSON(os.Stdout)

		pl.ResultCode = "1000"
		fmt.Println("\n******** Response body ********")
		pl.JSON(os.Stdout)
		pl.JSON(w)
	}
}

func (m *RemoteMessage) JSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func LogRemoteMessage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Println("######## INFO: New request received ########")
		fmt.Println("\n******** Request header ********")
		next(w, r)
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n")
	})
}
