package sprotoexport

import (
	"bytes"
	"fmt"
	"text/template"
	"strconv"
)

const templateStr = `
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
type {{.Name}} struct{ {{range $_, $field := .Fields}}{{GetClassFiledComment	$field}}
	{{$field.Name}} {{GetClassFiledType $field}} {{GetClassFiledDescribe $field}}{{end}}
}
{{end}}
`

func getClassFiledType(field * ClassField) string {

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
	case "boolean":
		str += "bool"
	case "string":
		str += "string"
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

func getClassFiledDescribe(field * ClassField) string {

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


func getClassFiledComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="\n	//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenGo(parser * SPParser){

	var bf bytes.Buffer

	funcMap := template.FuncMap{
		"GetClassFiledComment"	:getClassFiledComment,
		"GetClassFiledType"		:getClassFiledType,
		"GetClassFiledDescribe"	:getClassFiledDescribe,
	}

	tpl, err := template.New("genGolang").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = tpl.Execute(&bf, *parser)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(bf.Bytes()))
}
