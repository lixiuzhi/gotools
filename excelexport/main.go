package main

import (
	 excelTool "lixiuzhi/excelexport/core"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"encoding/json"
)

func pause(){
	fmt.Println("按回车键继续...")
	input := bufio.NewScanner(os.Stdin)
	if input.Scan(){
		return
	}
}

type Config struct {
	ExcelPath		string
	ExportCS		bool
	ExportCSPath	string
	ExportBin		bool
	ExportBinPath	string
}

func main() {

	fmt.Println("开始读取go版导表工具配置文件:goexport.cfg\n")

	cfg:=&Config{}

	if data,err:=ioutil.ReadFile("goexport.cfg");err!=nil{
		fmt.Printf("读取goexport.cfg错误:%s\n",err.Error())
		pause()
		return
	}else{
		if err:=json.Unmarshal(data,cfg);err!=nil{
			fmt.Printf("解析goexport.cfg错误:%s\n",err.Error())
			pause()
			return
		}else {
			fmt.Println(*cfg,"\n")
		}
	}

	if b,_:=excelTool.PathExists(cfg.ExcelPath);!b{
		fmt.Printf("excel表路径不存在，或者错误,路径:%s\n",cfg.ExcelPath)
		pause()
		return
	}

	if tables,err := excelTool.ReadAllExcel(cfg.ExcelPath);err!=nil{
		fmt.Println(err.Error())
	}else {
		datasheets :=excelTool.GetDataSheetInfo(tables)

		//生成cs代码
		if cfg.ExportCS{
			if err:= excelTool.GenCsharp(datasheets,cfg.ExportCSPath);err!=nil{
				fmt.Println(err.Error())
				pause()
				return
			}
		}

		//生成二进制文件
		if cfg.ExportBin{
			if err:= excelTool.GenBytes(datasheets,cfg.ExportBinPath);err!=nil{
				fmt.Println(err.Error())
				pause()
				return
			}
		}

	}
}
