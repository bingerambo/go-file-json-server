package utils

import (
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
)

type SimpleJsonParser struct {
}

func NewSimpleJsonParse() *SimpleJsonParser {
	return &SimpleJsonParser{}
}

func (sjp *SimpleJsonParser) Load(filename string) ([]byte, error) {

	if !Exists(filename) {
		return nil, errors.New("Load input file_pathname is empty")
	}

	return ioutil.ReadFile(filename)
}

func (sjp *SimpleJsonParser) Parse(json_datas []byte) (*simplejson.Json, error) {
	sj, err := simplejson.NewJson(json_datas)

	if err != nil || sj == nil {
		//fmt.Printf("%v\n", err)
		log.Println("something wrong when call NewJson")
		return nil, err
	}

	//fmt.Println(sj) //&{map[test:map[array:[1 2 3] arraywithsubs:[map[subkeyone:1] map[subkeytwo:2 subkeythree:3]] bignum:8000000000]]}

	return sj, nil
}
