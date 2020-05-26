// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"encoding/xml"
	"fmt"
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

// GenerateResponse generates XML response.
func GenerateResponse(resultCode, resultMessage string) []byte {
	response := methodResponse{}

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
	fmt.Printf("\n******** Response ********\n%v\n", string(responseXML))
	return responseXML
}

// GenerateResponseAdm generates XML response for AdministrativeMessage exclusively
func GenerateResponseCodeOnly(resultCode string) []byte {
	response := methodResponse{}

	member := member{}
	member.Name = "resultCode"
	member.Value.Int = resultCode

	response.Params.Member = append(response.Params.Member, member)

	responseXML, err := xml.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n******** Response ********\n%v\n", string(responseXML))
	return responseXML
}
