package core

import (
	"bytes"
	"fmt"
	"text/template"
	"io/ioutil"
	"errors"
)

type CSHelper struct {

}

func (this*CSHelper)getCSClassFieldReadFunc(field * ClassField) string {

	str := "read"

	f := false

	switch field.Type {
	case "int32":
		str += "_int32"
	case "int64":
		str += "_int64"
	case "bool":
		str += "_boolean"
	case "string":
		str += "_string"
	case "binary":
		str += "_bytes"
	case "double":
		str += "_double"
	default:
		f = true
		if field.TypeIsEnum {
			str += "_enum"
		}else{
			str += "_obj"
		}
	}

	if field.Repeatd {
		str += "_list"
	}

	if f {
		str += "<" + field.Type + ">"
	}

	str += "()"
	return str
}

func (this*CSHelper)getCSClassFieldWriteFuncName(field * ClassField) string {

	str := "write"

	if field.TypeIsEnum{
		str+="_enum"
	}else {
		switch field.Type {
		case "int32":
			str += "_int32"
		case "int64":
			str += "_int64"
		case "bool":
			str += "_boolean"
		case "string":
			str += "_string"
		case "binary":
			str += "_bytes"
		default:
			str	+= "_obj"
			if field.Repeatd{
				str +="<" + field.Type + ">"
			}
		}
	}

	return str
}

func (this*CSHelper)getCSClassFieldType(field * ClassField) string {

	str := ""

	var baseType = ""

	switch field.Type {
	case "binary":
		baseType = "byte[]"
	case "int32":
		baseType = "Int32"
	case "int64":
		baseType = "Int64"
	case "bool":
		baseType = "bool"
	case "string":
		baseType = "string"
	case "double":
		baseType = "double"
	default:
		baseType = field.Type
	}

	if field.Repeatd{
		str = "List<"+ baseType +">"
	}else{
		str = baseType
	}

	return str
}

func (this*CSHelper)getCSClassFieldComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}

func (this*CSHelper)getCSEnumFieldComment(field * EnumField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenCS(parser * SPParser,outPath string) error{

	if has,_:=PathExists(outPath);!has{
		return errors.New(fmt.Sprintf("生成cs，目录%s 不存在，或者出错!\n",outPath))
	}

	outfileName := parser.fileName +".cs"
	fmt.Println("开始生成csharp文件:",outfileName)

	helper:=&CSHelper{}

	funcMap := template.FuncMap{
		"GetEnumFieldComment"			:helper.getCSEnumFieldComment,
		"GetClassFieldComment"			:helper.getCSClassFieldComment,
		"GetClassFieldType"				:helper.getCSClassFieldType,
		"getClassFieldReadFunc"			:helper.getCSClassFieldReadFunc,
		"getClassFieldWriteFuncName"	:helper.getCSClassFieldWriteFuncName,
	}

	tpl, err := template.New("genCS").Funcs(funcMap).Parse(toCSTemplate)
	if err != nil {
		return err
	}

	var bf bytes.Buffer
	err = tpl.Execute(&bf, *parser)
	if err != nil {
		return err
	}

	ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)

	return nil
}
