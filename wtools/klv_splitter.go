// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"fmt"
	"strconv"
)

// KLVSplitter prints KLV sets and their corresponding definition to stdout.
func KLVSplitter(s string) {
	klvMap := map[string]string{
		"002": "Tracking Number             ",
		"004": "Original Transaction Amount ",
		"010": "Conversion Rate             ",
		"026": "Merchant Category Code      ",
		"037": "Retrieval Reference Number  ",
		"041": "Terminal ID                 ",
		"042": "Merchant Identifier         ",
		"048": "Fraud Scoring Data          ",
		"049": "Original Currency Code      ",
		"050": "From Account                ",
		"052": "Pin Block                   ",
		"085": "Markup Amount               ",
		"250": "Capture Mode                ",
		"251": "Network                     ",
		"252": "Fee Type                    ",
		"253": "Last Four Digits            ",
		"254": "Digitized Pan               ",
		"255": "Digitized Wallet ID         ",
		"256": "Adjustment Reason           ",
		"257": "Reference ID                ",
		"258": "Markup Type                 ",
		"259": "Acquirer Country            ",
		"900": "3D Secure OTP               ",
		"901": "Digitization Activation     ",
		"910": "Digitized Device ID         ",
		"911": "Digitized PAN Expiry        ",
		"912": "Digitized Masked FPAN       ",
		"913": "Digitized token reference   ",
		"999": "Generic Key                 ",
	}

	klvSlices := make([][]string, 0)

	if len(s) < 5 {
		fmt.Println("NOT a proper KLV - Error: < 5.")
	}

	for len(s) > 4 {
		key := s[:3]
		length := s[3:5]
		lengthInt, _ := strconv.Atoi(length)

		if len(s) < 5+lengthInt {
			fmt.Println("NOT a proper KLV - Error: < 5 + length")
			break
		}

		value := s[5:(5 + lengthInt)]
		klvSet := []string{key, length, value}
		klvSlices = append(klvSlices, klvSet)
		s = s[(5 + lengthInt):]
	}

	for _, klvElem := range klvSlices {
		if value, present := klvMap[klvElem[0]]; present {
			fmt.Printf("%s --- %s --- %s --- %s\n", value, klvElem[0], klvElem[1], klvElem[2])
		} else {
			fmt.Printf("Unknown --- %s --- %s --- %s\n", klvElem[0], klvElem[1], klvElem[2])
			fmt.Println("NOT a proper KLV - Error: unknown key.")
			break
		}
	}
}
