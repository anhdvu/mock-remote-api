package wtools

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type RemoteMessage struct {
	ActivationMethods         []ActivationMethod `json:"activationMethods,omitempty"`
	Challenge                 string             `json:"challenge,omitempty"`
	CurrencyCode              string             `json:"currencyCode,omitempty"`
	CustomerReference         string             `json:"customerReference"`
	DigitizedDeviceIdentifier string             `json:"digitizedDeviceIdentifier,omitempty"`
	DigitizedFpanMasked       string             `json:"digitizedFpanMasked,omitempty"`
	DigitizedPan              string             `json:"digitizedPan,omitempty"`
	DigitizedPanExpiry        string             `json:"digitizedPanExpiry,omitempty"`
	DigitizedTokenReference   string             `json:"digitizedTokenReference,omitempty"`
	EventType                 string             `json:"eventType,omitempty"`
	MerchantDescription       string             `json:"merchantDescription,omitempty"`
	MessageType               string             `json:"messageType"`
	ResultCode                string             `json:"resultCode,omitempty"`
	TokenRequestorID          string             `json:"tokenRequestorId,omitempty"`
	TrackingNumber            string             `json:"trackingNumber"`
	TransactionAmount         string             `json:"transactionAmount,omitempty"`
	WalletIdentifier          string             `json:"walletIdentifier,omitempty"`
}

type ActivationMethod struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}

func HandleRemoteMessage(terminal, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please use POST method"))
			return
		}
		ParseRemoteRequestHeader(r)
		receivedTerminal, receivedChecksum := extractAuthorizationData(r)

		defer r.Body.Close()
		rawPl, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error: Couldn't read raw payload.")
		}
		pl := &RemoteMessage{}
		err = json.Unmarshal(rawPl, pl)
		if err != nil {
			fmt.Printf("Couldn't unmarshal json payload, %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println("\n******** Request body ********")
		pl.JSON(os.Stdout)
		fmt.Println("\n******** Raw payload to hash ********")
		rawPayload := fmt.Sprintf("%q", rawPl)
		fmt.Println(strings.Trim(rawPayload, "\""))

		checksum := strings.ToUpper(calculateChecksum(password, rawPl))

		fmt.Println("\n******** Authorization data ********")
		fmt.Printf("Received terminal: %s\n", receivedTerminal)
		fmt.Printf("Configured terminal: %s\n", terminal)
		fmt.Printf("Received checksum: %s\n", receivedChecksum)
		fmt.Printf("Calculated checksum: %s\n", checksum)

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

		if receivedTerminal == terminal && receivedChecksum == checksum {
			pl.ResultCode = "0000"
		} else {
			pl.ResultCode = "1022"
		}

		fmt.Println("\n******** Response body ********")
		pl.JSON(os.Stdout)
		w.WriteHeader(http.StatusOK)
		pl.JSON(w)
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

func calculateChecksum(pk string, payload []byte) string {
	trimmed := bytes.TrimSpace(payload)
	digester := hmac.New(sha256.New, []byte(pk))
	digester.Write(trimmed)
	return hex.EncodeToString(digester.Sum(nil))
}

func extractAuthorizationData(r *http.Request) (string, string) {
	authHeader := r.Header.Get("Authorization")
	headerValues := strings.Split(authHeader, ",")
	terminal := strings.Split(headerValues[0], "=")[1]
	checksum := strings.Split(headerValues[1], "=")[1]

	return terminal, checksum
}
