package pkg

import (
	"math/rand"
	"strconv"
)

func GetSerialId(n int) string {
	t := "0000000"
	if len(strconv.Itoa(n+1)) == len(strconv.Itoa(n)) {
		return t[len(strconv.Itoa(n)):] + strconv.Itoa(n+1)
	}
	return t[len(strconv.Itoa(n))+1:] + strconv.Itoa(n+1)
}

func GenerateOTP() int {

	return rand.Intn(900000) + 100000
}
