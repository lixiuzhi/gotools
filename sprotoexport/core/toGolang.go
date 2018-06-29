package core

import (
	"bytes"
	"fmt"
	"text/template"
	"strconv"
	"io/ioutil"
	"errors"
)

const templateGoStr = `
package proto

import (
	"reflect"
)

//枚举区域
{{range $i, $enum := .Enums}}
const (
{{range $i,$enumfield :=$enum.Fields}}	{{$enum.Name}}_{{$enumfield.Name}} = {{$enumfield.LocalIndex}}
{{end}})
{{end}}

//结构体区域
{{range $i, $class := .Classes}}
type {{.Name}} struct{ {{range $_, $field := .Fields}}{{GetClassFieldComment	$field}}
	{{$field.Name}} {{GetClassFieldType $field}} {{GetClassFieldDescribe $field}}{{end}}
}
{{end}}
`

func getGoClassFieldType(field * ClassField) string {

	str := ""
	if field.Repeatd {
		str += "[]"
	}

	switch field.Type {
	case "binary":
		str += "[]byte"
	case "int32":
		str += "int32"
	case "int64":
		str += "int64"
	case "bool":
		str += "bool"
	case "string":
		str += "string"
	case "double":
		str +="float64"
	default:
		//如果是枚举类型
		if field.TypeIsEnum{
			str+="int32"
		}else {
			str += "*"
			str += field.Type
		}
	}
	return str
}

func getGoClassFieldDescribe(field * ClassField) string {

	str := " `sproto:\""

	switch field.Type {
	case "binary":
		str += "string"
	case "int32":
		str += "integer"
	case "int64":
		str += "integer"
	case "boolean":
		str += "boolean"
	case "string":
		str += "string"
	default:
		if field.TypeIsEnum{
			str+="integer"
		}else {
			str += field.Type
		}
	}

	str += ","

	str += strconv.Itoa(field.LocalIndex)

	str += ","

	if field.Repeatd && field.Type != "binary" {
		str += "array,"
	}

	str += "name="
	str += field.Name
	str += "\"`"

	return str
}


func getGoClassFieldComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="\n	//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenGo(parser * SPParser,outPath string) error{

	if has,_:=PathExists(outPath);!has{
		return errors.New(fmt.Sprintf("生成golang，目录%s 不存在，或者出错!\n",outPath))
	}

	outfileName :=parser.fileName+".go"

	fmt.Println("开始生成go文件:",outfileName)

	var bf bytes.Buffer

	funcMap := template.FuncMap{
		"GetClassFieldComment"	:getGoClassFieldComment,
		"GetClassFieldType"		:getGoClassFieldType,
		"GetClassFieldDescribe"	:getGoClassFieldDescribe,
	}

	tpl, err := template.New("genGolang").Funcs(funcMap).Parse(templateGoStr)
	if err != nil {
		return nil
	}

	err = tpl.Execute(&bf, *parser)
	if err != nil {
		return nil
	}

	ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)

	return nil
}
