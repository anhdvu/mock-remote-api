// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"encoding/xml"
	"fmt"
	"os"
)

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

// Move response code to env var so the setting can be dynamically changed
func GetRespCode() string {
	respCode := os.Getenv("RESP_CODE")
	if respCode == "" {
		respCode = "1"
	}
	return respCode
}

// GenerateResponse generates XML response.
func GenerateResponsewText(resultMessage string) []byte {
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

// GenerateResponseCodeOnly generates XML response for AdministrativeMessage exclusively
func GenerateResponseCodeOnly() []byte {
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

// GenerateResponsewActivationMethods generates XML response for AdministrativeMessage digitization.activationmethods message type
func GenerateResponsewActivationMethods() []byte {
	response := []byte("<methodResponse><params><param><value><struct><member><name>resultCode</name><value><int>1</int></value></member><member><name>activationMethods</name><value><array><data><value><struct><member><name>type</name><value>1</value></member><member><name>value</name><value>1(###) ### 4567</value></member></struct></value><value><struct><member><name>type</name><value>2</value></member><member><name>value</name><value>2a***d@anymail.com</value></member></struct></value></data></array></value></member></struct></value></param></params></methodResponse>")

	// responseXML, err := xml.Marshal(response)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%v\n", string(responseXML))
	// return responseXML
	fmt.Printf("%v\n", string(response))
	return response
}
