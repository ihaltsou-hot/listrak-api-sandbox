package view

import (
	"fmt"
	"strconv"
)

func StringToPhone(phoneNumber string) string {
	areaCode := phoneNumber[0:3]
	prefix := phoneNumber[3:6]
	lineNumber := phoneNumber[6:10]

	formattedPhoneNumber := fmt.Sprintf("+1 (%s) %s-%s", areaCode, prefix, lineNumber)
	return formattedPhoneNumber
}

func String(i int) string {
	return strconv.Itoa(i)
}
