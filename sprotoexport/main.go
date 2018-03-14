package main

import (
	spexport "github.com/lixiuzhi/gotools/sprotoexport/core"
	"io/ioutil"
	"fmt"
	"strings"
)

func main(){
	scanner := &spexport.Scanner{}
	//data,_:=ioutil.ReadFile("test.sp")

	var basePath ="protos"

	fileinfos,err :=ioutil.ReadDir(basePath)
	if err!=nil{
		fmt.Errorf("读取协议目录出错,",err.Error())
		return
	}

	var allFileStr = ""

	for _,v :=range fileinfos{

		if  strings.HasSuffix(strings.ToLower(v.Name()),".sp"){

			fmt.Println("开始读取文件:",basePath+"/"+v.Name())

			if data,err1:=ioutil.ReadFile(basePath+"/"+v.Name());err1==nil{
				allFileStr+=string(data)
			}else {

				fmt.Println("读取协议文件出错,",err1.Error())
				return
			}
		}
	}

	tokens,_ := scanner.GetTokens(allFileStr)
	parser :=&spexport.SPParser{}
	parser.Parse(tokens,"proto")

	//spexport.GenGo(parser,"bin/out")
	spexport.GenCS(parser,"bin/out")
	spexport.GenJava(parser,"com.lxz","bin/out")
}
