// Package walletutils implements an additional function to analyze XMLs in Remote API request body
package walletutils

import (
	"encoding/xml"
	"fmt"
)

// GenerateResponse generates XML response.
func GenerateResponse(resCode string) string {
	type Response struct {
		XMLName xml.Name `xml:"methodResponse"`
		Params  struct {
			Param struct {
				Value struct {
					Struct struct {
						Member struct {
							Name  string `xml:"name"`
							Value struct {
								Int string `xml:"int"`
							} `xml:"value"`
						} `xml:"member"`
					} `xml:"struct"`
				} `xml:"value"`
			} `xml:"param"`
		} `xml:"params"`
	}

	response := &Response{}
	response.Params.Param.Value.Struct.Member.Name = "resultCode"
	response.Params.Param.Value.Struct.Member.Value.Int = resCode

	responseXML, err := xml.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n******** Response ********\n%v\n", string(responseXML))
	return string(responseXML)
}
