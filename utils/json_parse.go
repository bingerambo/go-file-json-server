package utils

import (
	"encoding/json"
	"io/ioutil"
)

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

/*
将内容为json格式的文件，读取并转为结构体对象
*/
func (jst *JsonStruct) Load(filename string, v interface{}) error {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

/*
将json字符串bytes转为结构体对象
*/
func (jst *JsonStruct) ParseBytes(read_bytes []byte, v interface{}) (err error) {
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(read_bytes, v)
	if err != nil {
		return err
	}

	return nil
}

/*
将json字符串str转为结构体对象
*/
func (jst *JsonStruct) ParseStr(json_str string, v interface{}) (err error) {
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal([]byte(json_str), v)
	if err != nil {
		return err
	}

	return nil
}
