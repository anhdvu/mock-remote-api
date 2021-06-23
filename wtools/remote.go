package wtools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RemoteMessage struct {
	MessageType               string             `json:"messageType"`
	TrackingNumber            string             `json:"trackingNumber"`
	CustomerReference         string             `json:"customerReference"`
	EventType                 string             `json:"eventType,omitempty"`
	DigitizedDeviceIdentifier string             `json:"digitizedDeviceIdentifier,omitempty"`
	DigitizedPan              string             `json:"digitizedPan,omitempty"`
	DigitizedPanExpiry        string             `json:"digitizedPanExpiry,omitempty"`
	DigitizedFpanMasked       string             `json:"digitizedFpanMasked,omitempty"`
	DigitizedTokenReference   string             `json:"digitizedTokenReference,omitempty"`
	WalletIdentifier          string             `json:"walletIdentifier,omitempty"`
	Challenge                 string             `json:"challenge,omitempty"`
	TokenRequestorId          string             `json:"tokenRequestorId,omitempty"`
	ActivationMethods         []ActivationMethod `json:"activationMethods,omitempty"`
	ResultCode                string             `json:"resultCode,omitempty"`
}

type ActivationMethod struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}

func HandleRemoteMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please use POST method"))
		} else {
			ParseRemoteRequestHeader(r)
			defer r.Body.Close()
			pl := &RemoteMessage{}
			err := pl.ParseJSON(r.Body)
			if err != nil {
				fmt.Printf("Couldn't unmarshal json payload, %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			fmt.Println("\n******** Request body ********")
			pl.JSON(os.Stdout)

			if pl.MessageType == "digitization.activationmethods" {
				am1 := ActivationMethod{1, "1(###) ### 4567"}
				am2 := ActivationMethod{2, "2a***d@anymail.com"}

				pl.ActivationMethods = append(pl.ActivationMethods, am1, am2)
			}

			// Available result code
			// 0000		Message processed and confirmed.
			// 1000		API internal error.
			// 1022		Authorization error.
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

func (m *RemoteMessage) ParseJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(m)
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
