// Package walletutils implements an additional function to analyze XMLs in Remote API request body
package walletutils

import (
	// "bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type reqKLV interface {
	GetKLV() string
}

// Deduct struct type used to hold data of Deduct calls.
type Deduct struct {
	MethodName string `json:"method name"`
	Terminal   string `json:"terminal"`
	Reference  string `json:"reference"`
	Amount     string `json:"amount"`
	Narrative  string `json:"narrative"`
	TxnType    string `json:"transaction type"`
	KLVData    string `json:"klv string"`
	TxnID      string `json:"transaction id"`
	TxnDate    string `json:"transaction date"`
	Checksum   string `json:"checksum"`
}

// GetKLV gets KLV from Deduct calls.
func (dd *Deduct) GetKLV() string {
	return dd.KLVData
}

// Settlement struct type used to hold data of adjustment or reversal calls.
type Settlement struct {
	MethodName string `json:"method name"`
	Terminal   string `json:"terminal"`
	Reference  string `json:"reference"`
	Amount     string `json:"amount"`
	Narrative  string `json:"narrative"`
	KLVData    string `json:"klv string"`
	RefTxnID   string `json:"reference transaction id"`
	RefTxnDate string `json:"reference transaction date"`
	TxnID      string `json:"transaction id"`
	TxnDate    string `json:"transaction date"`
	Checksum   string `json:"checksum"`
}

// GetKLV gets KLV from adjustment or reversal calls.
func (settle *Settlement) GetKLV() string {
	return settle.KLVData
}

// AdminMessage struct type used to hold data of administrative message calls.
type AdminMessage struct {
	MethodName string `json:"method name"`
	Terminal   string `json:"terminal"`
	Reference  string `json:"reference"`
	MsgType    string `json:"message type"`
	KLVData    string `json:"klv string"`
	TxnID      string `json:"transaction id"`
	TxnDate    string `json:"transaction date"`
}

// GetKLV gets KLV from administrative message calls.
func (admmsg *AdminMessage) GetKLV() string {
	return admmsg.KLVData
}

// Payload struct type used to parse all XML to struct based data type with which Go can work.
type Payload struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
	Params     []struct {
		Value struct {
			StringParam string `xml:"string,omitempty"`
			IntParam    string `xml:"int,omitempty"`
			TimeParam   string `xml:"dateTime.iso8601,omitempty"`
		} `xml:"value"`
	} `xml:"params>param"`
}

// ParseRemoteRequestBody parses incoming requests to JSON and print to stdout.
func ParseRemoteRequestBody(r *http.Request) *Payload {
	rBodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	xmlPayload := &Payload{}
	err = xml.Unmarshal(rBodyBytes, xmlPayload)
	if err != nil {
		fmt.Println(err)
	}

	return xmlPayload
}

//  ParseMethod parses payload struct to a specific api method.
func ParseMethod(payload *Payload) reqKLV {
	var request reqKLV // vô duyên vl
	switch payload.MethodName {
	case "Deduct":
		request = &Deduct{
			MethodName: payload.MethodName,
			Terminal:   payload.Params[0].Value.StringParam,
			Reference:  payload.Params[1].Value.StringParam,
			Amount:     payload.Params[2].Value.IntParam,
			Narrative:  payload.Params[3].Value.StringParam,
			TxnType:    payload.Params[4].Value.StringParam,
			KLVData:    payload.Params[5].Value.StringParam,
			TxnID:      payload.Params[6].Value.StringParam,
			TxnDate:    payload.Params[7].Value.TimeParam,
			Checksum:   payload.Params[8].Value.StringParam,
		}
	case "AdministrativeMessage":
		request = &AdminMessage{
			MethodName: payload.MethodName,
			Terminal:   payload.Params[0].Value.StringParam,
			Reference:  payload.Params[1].Value.StringParam,
			MsgType:    payload.Params[2].Value.StringParam,
			KLVData:    payload.Params[3].Value.StringParam,
			TxnID:      payload.Params[4].Value.StringParam,
			TxnDate:    payload.Params[5].Value.TimeParam,
		}
	default:
		request = &Settlement{
			MethodName: payload.MethodName,
			Terminal:   payload.Params[0].Value.StringParam,
			Reference:  payload.Params[1].Value.StringParam,
			Amount:     payload.Params[2].Value.IntParam,
			Narrative:  payload.Params[3].Value.StringParam,
			KLVData:    payload.Params[4].Value.StringParam,
			RefTxnID:   payload.Params[5].Value.StringParam,
			RefTxnDate: payload.Params[6].Value.TimeParam,
			TxnID:      payload.Params[7].Value.StringParam,
			TxnDate:    payload.Params[8].Value.TimeParam,
			Checksum:   payload.Params[9].Value.StringParam,
		}
	}

	return request
}

// DumpJSON function dumps parsed request to os.Stdout in JSON format.
func DumpJSON(reqBody reqKLV) {
	reqJSON, err := json.MarshalIndent(reqBody, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n******** Request body in JSON ********\n%v\n", string(reqJSON))

	/* // An attemp of using json encoder to print to os.Stdout
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("encoder", "  ")
	enc.Encode(request)
	fmt.Println(buf.String()) */
}

// ParseRemoteRequestHeaders parses incoming requests' headers and print to stdout.
func ParseRemoteRequestHeaders(r *http.Request) {
	fmt.Println("\n******** Request headers ********")
	for k, v := range r.Header {
		fmt.Printf("%q: %v\n", k, v)
	}
}
