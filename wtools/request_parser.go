// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	// "bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReqData interface {
	KLV() string
	String() string
}

// Deduct struct type used to hold data of Deduct calls.
type DeductLoad struct {
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

// KLV gets KLV from Deduct calls.
func (dd *DeductLoad) KLV() string {
	return dd.KLVData
}

// String returns string to hash.
func (dd *DeductLoad) String() string {
	return fmt.Sprint(dd.MethodName + dd.Terminal + dd.Reference + dd.Amount + dd.Narrative + dd.TxnType + dd.KLVData + dd.TxnID + dd.TxnDate)
}

// LoadAdvice struct used to hold data of LoadAdvice calls.
type LoadAdvice struct {
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

// KLV gets KLV from LoadAdvice calls.
func (la *LoadAdvice) KLV() string {
	return la.KLVData
}

// String returns string to hash.
func (la *LoadAdvice) String() string {
	return fmt.Sprint(la.MethodName + la.Terminal + la.Reference + la.Amount + la.Narrative + la.TxnType + la.KLVData + la.TxnID + la.TxnDate)
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

// KLV gets KLV from adjustment or reversal calls.
func (st *Settlement) KLV() string {
	return st.KLVData
}

// String returns string to hash.
func (st *Settlement) String() string {
	return fmt.Sprint(st.MethodName + st.Terminal + st.Reference + st.Amount + st.Narrative + st.KLVData + st.RefTxnID + st.RefTxnDate + st.TxnID + st.TxnDate)
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

// KLV gets KLV from administrative message calls.
func (adm *AdminMessage) KLV() string {
	return adm.KLVData
}

// String returns message type
func (adm *AdminMessage) String() string {
	return adm.MsgType
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
	rawBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	xmlPayload := &Payload{}
	err = xml.Unmarshal(rawBody, xmlPayload)
	if err != nil {
		fmt.Println(err)
	}

	return xmlPayload
}

//  ParseMethod parses payload struct to a specific api method.
func ParseMethod(payload *Payload) ReqData {
	var request ReqData // vô duyên vl
	switch payload.MethodName {
	case "Deduct", "LoadAdvice":
		request = &DeductLoad{
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
func DumpJSON(reqBody ReqData) {
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
