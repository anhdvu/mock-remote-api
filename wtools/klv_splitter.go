// Package wtools implements an additional function to analyze XMLs in Remote API request body
package wtools

import (
	"fmt"
	"strconv"
)

// KLVSplitter prints KLV sets and their corresponding definition to stdout.
func KLVSplitter(s string) {
	klvMap := map[string]string{
		"002": "Tracking Number                 ",
		"004": "Original Transaction Amount     ",
		"010": "Conversion Rate                 ",
		"026": "Merchant Category Code          ",
		"032": "Acquiring Instituion Code       ",
		"037": "Retrieval Reference Number      ",
		"041": "Terminal ID                     ",
		"042": "Merchant Identifier             ",
		"043": "Merchant Description            ",
		"048": "Fraud Scoring Data              ",
		"049": "Original Currency Code          ",
		"050": "From Account                    ",
		"052": "Pin Block                       ",
		"085": "Markup Amount                   ",
		"250": "Capture Mode                    ",
		"251": "Network                         ",
		"252": "Fee Type                        ",
		"253": "Last Four Digits                ",
		"254": "Digitized Pan                   ",
		"255": "Digitized Wallet ID             ",
		"256": "Adjustment Reason               ",
		"257": "Original Deduct Reference ID    ",
		"258": "Markup Type                     ",
		"259": "Acquirer Country                ",
		"261": "Transaction Fee Amount          ",
		"262": "Transaction Subtype             ",
		"263": "placeholder                     ",
		"264": "placeholder                     ",
		"265": "placeholder                     ",
		"266": "placeholder                     ",
		"267": "placeholder                     ",
		"268": "placeholder                     ",
		"269": "placeholder                     ",
		"270": "Security Services Data          ",
		"300": "placeholder                     ",
		"301": "Second Additional Amount        ",
		"302": "Cashback POS Currency Code      ",
		"303": "Cashback POS Amount             ",
		"900": "3D Secure OTP                   ",
		"901": "Digitization Activation         ",
		"910": "Digitized Device ID             ",
		"911": "Digitized PAN Expiry            ",
		"912": "Digitized Masked FPAN           ",
		"913": "Digitized Token Reference       ",
		"915": "Digitized Token Requestor ID    ",
		"916": "Visa Digitized PAN              ",
		"999": "Generic Key                     ",
	}

	klvSlices := make([][]string, 0)

	if len(s) < 5 {
		fmt.Println("NOT a proper KLV - Error: < 5.")
		return
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
		if value, ok := klvMap[klvElem[0]]; ok {
			fmt.Printf("%s --- %s --- %s --- %s\n", value, klvElem[0], klvElem[1], klvElem[2])
		} else {
			fmt.Printf("Unknown --- %s --- %s --- %s\n", klvElem[0], klvElem[1], klvElem[2])
			fmt.Println("NOT a proper KLV - Error: unknown key.")
			break
		}
	}
}
