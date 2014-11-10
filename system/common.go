package system

import (
	"log"
	"strconv"
)

func toInt(b []byte) (i int) {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		log.Println("Failed to convert string to int")
	}
	return
}

func toFloat(b []byte) (f float64) {
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		log.Println("Failed to convert string to float")
	}
	return
}
