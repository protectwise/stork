package helper

import (
	"io/ioutil"
)

func Read_file(path string) []byte {
	data, err := ioutil.ReadFile(path)
	Fatal_error("Can't read file.", err)
	return data
}
