package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	excelTool "lixiuzhi/excelexport/core"
	"os"
	"os/exec"
	"strings"
)


func pause() {
	fmt.Println("按回车键继续...")
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		return
	}
}

type Config struct {
	ExcelPath     	string
	ExportCS      	bool
	ExportCSPath  	string
	ExportJava      bool
	ExportJavaPath	string
	ExportBin     	bool
	ExportBinPath 	string
	ExportDB		bool
	DSN				string
}

func PrintfError(s string ,a ...interface{}) {

	fmt.Printf("error******error*****error******error*****:::"+s,a...)
}


func main() {
	fmt.Println("开始读取go版导表工具配置文件:goexport.cfg\n")

	cfg := &Config{}

	if data, err := ioutil.ReadFile("excel_export.cfg"); err != nil {
		PrintfError("读取excel_export.cfg错误:%s\n", err.Error())
		os.Exit(1)
		return
	} else {
		if err := json.Unmarshal(data, cfg); err != nil {
			PrintfError("解析excel_export.cfg错误:%s\n", err.Error())
			pause()
			os.Exit(1)
			return
		} else {
			fmt.Println(*cfg, "\n")
		}
	}

	if b, _ := excelTool.PathExists(cfg.ExcelPath); !b {
		PrintfError("excel表路径不存在，或者错误,路径:%s\n", cfg.ExcelPath)
		pause()
		os.Exit(1)
		return
	}

	if tables, err := excelTool.ReadAllExcel(cfg.ExcelPath); err != nil {
		PrintfError("读取配置表出错:"+err.Error())
		os.Exit(1)
		pause()
		return
	} else {

		datasheets:= excelTool.GetDataSheetInfo(tables)

		//tempPath :=  getCurrentPath() + "data"
		//fmt.Println("tempPath" + tempPath)
		//生成cs代码
		if cfg.ExportCS {
			if err := excelTool.GenCsharp(datasheets, cfg.ExportCSPath); err != nil {
				PrintfError(err.Error())
				pause()
				os.Exit(1)
				return
			}
		}

		//生成二进制文件
		if cfg.ExportBin {
			if err := excelTool.GenBytes(datasheets, cfg.ExportBinPath); err != nil {
				PrintfError(err.Error())
				pause()
				os.Exit(1)
				return
			}
		}

		//生成java
		if cfg.ExportJava{
			if err := excelTool.GenJava(datasheets, cfg.ExportJavaPath); err != nil {
				PrintfError(err.Error())
				pause()
				os.Exit(1)
				return
			}
		}

		//导入db
		if cfg.ExportDB{
			if err := excelTool.GenDB(datasheets, cfg.DSN); err != nil {
				PrintfError(err.Error())
				pause()
				os.Exit(1)
				return
			}
		}
	}
}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
