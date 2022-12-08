package tool

import (
	"encoding/json"
	"fmt"
)

// s必须传指针类型
func Map2struct(m interface{}, s interface{}) {
	fmt.Printf("map=%v\n", m)
	marshal, err := json.Marshal(m)
	if err != nil {
		fmt.Println("marshal:", err)
		panic(err)
	}
	err = json.Unmarshal(marshal, &s)
	if err != nil {
		fmt.Println("unmarshal:", err)
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("struct=%v\n", s)
}
