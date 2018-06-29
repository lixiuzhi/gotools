package main

import (
	spexport "lixiuzhi/sprotoexport/core"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"bufio"
)

type Config struct {
	ProtoPath				string
	ExportGolang			bool
	ExportGolangPath		string
	ExportCS				bool
	ExportCSPath			string
	ExportJava				bool
	JavaPackageName			string
	ExportJavaPath			string
	ExportLua				bool
	ExportLuaPath			string
	ExportEmmyluaAPI		bool
	ExportEmmyluaAPIPath	string
}

func pause() {
	fmt.Println("按回车键继续...")
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		return
	}
}

func main(){

	cfg := &Config{}

	if data, err := ioutil.ReadFile("sproto_export.cfg"); err != nil {
		fmt.Println("读取sproto_export.cfg错误:%s\n", err.Error())
		os.Exit(1)
		return
	} else {
		if err := json.Unmarshal(data, cfg); err != nil {
			fmt.Println("解析sproto_export.cfg错误:%s\n", err.Error())
			pause()
			os.Exit(1)
			return
		} else {
			fmt.Println(*cfg, "\n")
		}
	}

	scanner := &spexport.Scanner{}
	//data,_:=ioutil.ReadFile("test.sp")


	fileinfos,err :=ioutil.ReadDir(cfg.ProtoPath)
	if err!=nil{
		fmt.Errorf("读取协议目录出错,",err.Error())
		return
	}

	//var allFileStr = ""

	var tokens = make([]*spexport.TokenInfo,0,10)

	var offset = 0

	for _,v :=range fileinfos{

		if  strings.HasSuffix(strings.ToLower(v.Name()),".sp"){

			fmt.Println("开始读取文件:",cfg.ProtoPath+"/"+v.Name())

			if data,err1:=ioutil.ReadFile(cfg.ProtoPath+"/"+v.Name());err1==nil{

				var newtokens = scanner.GetTokens(string(data),v.Name())

				var j = 0

				for _,v:=range newtokens{
					v.LineOffset = offset
					j = v.LocalLine+1
				}

				offset +=j

				tokens = append(tokens,newtokens...)

			}else {

				fmt.Println("读取协议文件出错,",err1.Error())
				return
			}
		}
	}

	parser :=&spexport.SPParser{}
	if err:=parser.Parse(tokens,"proto");err!=nil{
		fmt.Println(err.Error())
		pause()
		os.Exit(1)
	}

	if cfg.ExportGolang{
		if err:=spexport.GenGo(parser,cfg.ExportGolangPath);err!=nil{
			fmt.Println(err.Error())
			pause()
			os.Exit(1)
		}
	}

	if cfg.ExportCS{
		if err:= spexport.GenCS(parser,cfg.ExportCSPath);err!=nil{
			fmt.Println(err.Error())
			pause()
			os.Exit(1)
		}
	}

	if cfg.ExportJava {
		if err := spexport.GenJava(parser, cfg.JavaPackageName, cfg.ExportJavaPath); err != nil {
			fmt.Println(err.Error())
			pause()
			os.Exit(1)
		}
	}

	if cfg.ExportLua {
		if err := spexport.GenLua(parser, cfg.ExportLuaPath); err != nil {
			fmt.Println(err.Error())
			pause()
			os.Exit(1)
		}
	}

	if cfg.ExportEmmyluaAPI {
		if err := spexport.GenLuaAPI(parser, cfg.ExportEmmyluaAPIPath); err != nil {
			fmt.Println(err.Error())
			pause()
			os.Exit(1)
		}
	}
}
