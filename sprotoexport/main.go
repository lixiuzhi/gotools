package main

import (
	spexport "github.com/lixiuzhi/gotools/sprotoexport/core"
)

var testProto = `

message Test1{
    A int64
}

// 注释注释注释
message Test2{

	A int64

	B string

	C int32

	D []Test1

	E []int32 // 注释测试

	F Test1

	G GTest3

	H binary
}


enum Test3 {
	OK					// 成功
	ERROR			    //error
    OTHER			// 其他

}

`

func main(){
	scanner := &spexport.Scanner{}
	//data,_:=ioutil.ReadFile("test.sp")

	tokens,_ := scanner.GetTokens(testProto)

	parser :=&spexport.SPParser{}
	parser.Parse(tokens)

	spexport.GenGo(parser)
}
