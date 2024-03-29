// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"encoding/xml"
	"fmt"
	"os"
	"sync"
)

type response struct {
	code int
	text string
}

// MethodResponseMap is an in-memory db of responses
type MethodResponseMap struct {
	mu        sync.RWMutex
	responses map[string]*response
}

type methodResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Params  struct {
		Member []member `xml:"member"`
	} `xml:"params>param>value>struct"`
}

type member struct {
	Name  string `xml:"name"`
	Value struct {
		Int    string `xml:"int,omitempty"`
		String string `xml:"string,omitempty"`
	} `xml:"value"`
}

// GetRespCode move response code to env var so the setting can be dynamically changed
func GetRespCode() string {
	respCode := os.Getenv("RESP_CODE")
	if respCode == "" {
		respCode = "1"
	}
	return respCode
}

// RespwText generates XML response.
func RespwText(resultMessage string) []byte {
	response := methodResponse{}

	resultCode := GetRespCode()

	member1 := member{}
	member1.Name = "resultCode"
	member1.Value.Int = resultCode

	member2 := member{}
	member2.Name = "resultText"
	member2.Value.String = resultMessage

	response.Params.Member = append(response.Params.Member, member1, member2)

	responseXML, err := xml.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", string(responseXML))
	return responseXML
}

// RespCodeOnly generates XML response for AdministrativeMessage exclusively
func RespCodeOnly() []byte {
	response := methodResponse{}

	resultCode := GetRespCode()

	member := member{}
	member.Name = "resultCode"
	member.Value.Int = resultCode

	response.Params.Member = append(response.Params.Member, member)

	responseXML, err := xml.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", string(responseXML))
	return responseXML
}

// RespwCVM generates XML response for AdministrativeMessage digitization.activationmethods message type
func RespwCVM(ref string) []byte {
	defaultResp := "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>4</value></member><member><name>value</name><value>+27744704621</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>"

	// Dummy data to test.
	cardMethodList := map[string]string{
		"VTSTest1":    "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+1(###) ### 4567</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>joh***n@anymail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest2":    "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+1(###) ### 1234</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>jan***s@anymail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest123":  "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+1(###) ### 8888</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>xo***a@anymail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest1234": "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+1(###) ### 9999</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>som***s@anymail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest120":  "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+27790241919</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>tumimalesela1@gmail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest121":  "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+27790241919</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>tumimalesela1@gmail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"VTSTest122":  "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>+27790241919</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>tumimalesela1@gmail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"Bettr1":      "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>4</value></member><member><name>value</name><value>+27744704621</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"Bettr2":      "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>4</value></member><member><name>value</name><value>+27744704621</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
		"Bettr3":      "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>4</value></member><member><name>value</name><value>+27744704621</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>",
	}

	var response []byte

	if _, ok := cardMethodList[ref]; ok {
		response = []byte(cardMethodList[ref])
	} else {
		response = []byte(defaultResp)
	}

	fmt.Printf("%v\n", string(response))
	return response
}

// RespwBalance generates XML response for Balance call
func RespwBalance() []byte {
	defaultResp := "<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>balanceAmount</name><value><int>626900</int></value></member></struct></value></param></params></methodResponse>"

	var response []byte

	response = []byte(defaultResp)
	fmt.Printf("%v\n", string(response))
	return response
}
