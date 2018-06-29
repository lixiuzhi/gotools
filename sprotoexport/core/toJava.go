package core

import (
	"bytes"
	"fmt"
	"text/template"
	"io/ioutil"
	"errors"
)

func getJavaClassFieldReadFunc(field * ClassField) string {

	str := "read"

	isclass := false

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

		if field.TypeIsEnum {
			str += "_int32"
		}else{
			isclass = true
			str += "_obj"
		}
	}

	if field.Repeatd {
		str += "_list"
	}

	if isclass {
		str += "("+ field.Type +".proto_supplier)"
	}else {
		str += "()"
	}
	return str
}

func getJavaClassFieldWriteFuncName(field * ClassField) string {

	str := "write"

	if field.TypeIsEnum{
		str+="_int32"
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
		}
	}

	return str
}

func getJavaClassFieldType(field * ClassField) string {

	str := ""

	var baseType = ""
	var ListBaseType = ""

	switch field.Type {
	case "binary":
		baseType = "byte[]"
		ListBaseType = "byte[]"
	case "int32":
		baseType = "int"
		ListBaseType = "Integer"
	case "int64":
		baseType = "long"
		ListBaseType = "Long"
	case "bool":
		baseType = "boolean"
		ListBaseType = "Boolean"
	case "string":
		baseType = "String"
		ListBaseType = "String"
	default:
		baseType = field.Type
		ListBaseType = field.Type
	}

	if field.TypeIsEnum{
		baseType = "int"
		ListBaseType = "Integer"
	}

	if field.Repeatd{
		str = "List<"+ ListBaseType +">"
	}else{
		str = baseType
	}

	return str
}

func getJavaClassFieldComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}

func getJavaEnumFieldComment(field * EnumField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenJava(parser * SPParser,packageName string,outPath string)error{

	if has,_:=PathExists(outPath);!has{
		return errors.New(fmt.Sprintf("生成java，目录%s 不存在，或者出错!\n",outPath))
	}

	funcMap := template.FuncMap{
		"GetJavaPackageName"			:func() string{return packageName},

		"GetEnumFieldComment"			:getJavaEnumFieldComment,
		"GetClassFieldComment"			:getJavaClassFieldComment,
		"GetClassFieldType"				:getJavaClassFieldType,
		"getClassFieldReadFunc"			:getJavaClassFieldReadFunc,
		"getClassFieldWriteFuncName"	:getJavaClassFieldWriteFuncName,
	}

	//导出枚举类
	enumTpl, err := template.New("genJavaEnum").Funcs(funcMap).Parse(toJavaEnumTemplate)
	if err != nil {
		return err
	}

	for _,enum :=range parser.Enums{

		outfileName := enum.Name +".java"
		fmt.Println("开始生成java文件:",outfileName)

		var bf bytes.Buffer
		err = enumTpl.Execute(&bf, enum)
		if err != nil {
			return err
		}
		ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)
	}

	//导出class
	classTpl, err := template.New("genJavaClass").Funcs(funcMap).Parse(toJavaClassTemplate)
	if err != nil {
		return err
	}

	for _,class :=range parser.Classes{

		outfileName := class.Name +".java"
		fmt.Println("开始生成java文件:",outfileName)

		var bf bytes.Buffer
		err = classTpl.Execute(&bf, class)
		if err != nil {
			return err
		}
		ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)
	}

	//导出classMap
	classMapTpl, err := template.New("genJavaClassMap").Funcs(funcMap).Parse(toJavaMessageUtilTemplate)
	if err != nil {
		return err
	}

	var bf bytes.Buffer
	err = classMapTpl.Execute(&bf, parser)
	if err != nil {
		return err
	}
	ioutil.WriteFile(outPath+"/MessageUtil.java",bf.Bytes(),0666)

	return nil
}
