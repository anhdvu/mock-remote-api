// Package walletutils implements an additional function to analyze XMLs in Remote API request body
package walletutils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type respStruct interface {
	GetKLV() string
}

// Deduct struct type used to hold data of Deduct calls
type Deduct struct {
	MethodName string
	Terminal   string
	Reference  string
	Amount     string
	Narrative  string
	TxnType    string
	KLVData    string
	TxnID      string
	TxnDate    string
	Checksum   string
}

// GetKLV gets KLV from Deduct calls
func (dd *Deduct) GetKLV() string {
	return dd.KLVData
}

// Settlement struct type used to hold data of adjustment or reversal calls
type Settlement struct {
	MethodName string
	Terminal   string
	Reference  string
	Amount     string
	Narrative  string
	KLVData    string
	RefTxnID   string
	RefTxnDate string
	TxnID      string
	TxnDate    string
	Checksum   string
}

// GetKLV gets KLV from Adjustment calls
func (settle *Settlement) GetKLV() string {
	return settle.KLVData
}

// Call struct type used to parse all XML from Tutuka
type Call struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
	Params     []struct {
		Value struct {
			StringParam string `xml:"string"`
			IntParam    string `xml:"int"`
			TimeParam   string `xml:"dateTime.iso8601"`
		} `xml:"value"`
	} `xml:"params>param"`
}

// ParseRemoteRequestBody parses incoming requests to JSON and print to stdout.
func ParseRemoteRequestBody(r *http.Request) respStruct {
	rBodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	// data := string(bodyBytes)

	xmlCall := &Call{}
	// err := xml.Unmarshal([]byte(data), xmlCall)
	err = xml.Unmarshal(rBodyBytes, xmlCall)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("%+v\n", xmlCall)

	var request respStruct
	switch xmlCall.MethodName {
	case "Deduct":
		request = &Deduct{
			MethodName: xmlCall.MethodName,
			Terminal:   xmlCall.Params[0].Value.StringParam,
			Reference:  xmlCall.Params[1].Value.StringParam,
			Amount:     xmlCall.Params[2].Value.IntParam,
			Narrative:  xmlCall.Params[3].Value.StringParam,
			TxnType:    xmlCall.Params[4].Value.StringParam,
			KLVData:    xmlCall.Params[5].Value.StringParam,
			TxnID:      xmlCall.Params[6].Value.StringParam,
			TxnDate:    xmlCall.Params[7].Value.TimeParam,
			Checksum:   xmlCall.Params[8].Value.StringParam,
		}
	default:
		request = &Settlement{
			MethodName: xmlCall.MethodName,
			Terminal:   xmlCall.Params[0].Value.StringParam,
			Reference:  xmlCall.Params[1].Value.StringParam,
			Amount:     xmlCall.Params[2].Value.IntParam,
			Narrative:  xmlCall.Params[3].Value.StringParam,
			KLVData:    xmlCall.Params[4].Value.StringParam,
			RefTxnID:   xmlCall.Params[5].Value.StringParam,
			RefTxnDate: xmlCall.Params[6].Value.TimeParam,
			TxnID:      xmlCall.Params[7].Value.StringParam,
			TxnDate:    xmlCall.Params[8].Value.TimeParam,
			Checksum:   xmlCall.Params[9].Value.StringParam,
		}
	}

	requestJSON, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n******** Request body in JSON ********\n%v\n", string(requestJSON))
	return request
}

// ParseRemoteRequestHeaders parses incoming requests' headers and print to stdout.
func ParseRemoteRequestHeaders(r *http.Request) {
	fmt.Println("\n******** Request headers ********")
	for k, v := range r.Header {
		fmt.Printf("%q: %v\n", k, v)
	}
}
