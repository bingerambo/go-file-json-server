package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestSimpleJsonParser(t *testing.T) {
	file_path := "D:/Go_projects/src/github.com/bingerambo/go-file-json-server/tmp/xihu-xorigin.json"

	sjp := NewSimpleJsonParse()
	datas, err := sjp.Load(file_path)
	if err != nil {
		log.Fatal(err)
		return 
	}

	js, err  := sjp.Parse(datas)
	if err !=nil{
		log.Fatal(err)
	}
	ma1, err := js.Get("data").Get("schedule_strategies").Float64();
	if err !=nil{
		log.Fatal(err)
	}
	//ma1 := js.Get("data").Get("schedule_strategies").MustArray([]interface{}{"1", 2, "3"})
	//assert.Equal(t, ma2, []interface{}{"1", 2, "3"})
	fmt.Println(ma1)
	ma2 := js.Get("data").Get("schedule_strategies").MustStringArray([]string{"general-global"})
	fmt.Println(ma2)
}
