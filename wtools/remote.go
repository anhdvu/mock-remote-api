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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please use POST method"))
		} else {
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

			// Available result code
			// 0000		Message processed and confirmed.
			// 1000		API internal error.
			// 1101		Balance limit exceeded.
			// 1102		Moving annual top up limit exceeded.
			pl.ResultCode = "0000"

			fmt.Println("\n******** Response body ********")
			pl.JSON(os.Stdout)
			w.WriteHeader(http.StatusOK)
			pl.JSON(w)
		}
	}
}

func (m *RemoteMessage) JSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func LogRemoteMessage(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Println("######## INFO: New request received ########")
		fmt.Println("\n******** Request header ********")
		next(w, r)
		fmt.Printf("\n######## INFO: Request parse completed ########\n\n")
	}
}
