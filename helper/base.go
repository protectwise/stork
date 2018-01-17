package helper

import (
	"log"
)

func Fatal_error(msg string, e error) {
	if e != nil {
		log.Fatal(msg + e.Error())
	}
}
