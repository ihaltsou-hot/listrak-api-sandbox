package view

import (
	"fmt"
	"strconv"
)

func StringToPhone(phoneNumber string) string {
	areaCode := phoneNumber[1:4]
	prefix := phoneNumber[4:7]
	lineNumber := phoneNumber[7:11]

	formattedPhoneNumber := fmt.Sprintf("+1 (%s) %s-%s", areaCode, prefix, lineNumber)
	return formattedPhoneNumber
}

func String(i int) string {
	return strconv.Itoa(i)
}
