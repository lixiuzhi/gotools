package core

import (
    "errors"
    "fmt"
    "bytes"
    "io/ioutil"
    "text/template"
    "strconv"
)

type LuaHelper struct {
    parser *SPParser
}

func (this*LuaHelper)getClassFieldType(field * ClassField) string {

    str := ""

    var baseType = ""

    switch field.Type {
    case "binary":
        baseType = "binary"
    case "int32":
        baseType = "integer"
    case "int64":
        baseType = "integer"
    case "bool":
        baseType = "boolean"
    case "string":
        baseType = "string"
    case "double":
        baseType = "integer"
    default:
        baseType = field.Type
    }

    if field.TypeIsEnum{
        baseType = "integer"
    }

    if field.Repeatd{
        str = "*"+ baseType
    }else{
        str = baseType
    }

    return str
}

func (this*LuaHelper)getFieldDefaultValue(field * ClassField) string{

    if field.Repeatd{
        return "nil"
    }

    if field.TypeIsEnum {
        //得到枚举默认值
        var enum = this.parser.GetEnumByName(field.Type)
        if len(enum.Fields) > 0 {
            return strconv.Itoa(enum.Fields[0].LocalIndex)
        } else {
            return "0"
        }
    }

    switch field.Type {
    case "binary":
        return "nil"
    case "int32":
        return "0"
    case "int64":
        return "0"
    case "bool":
        return "false"
    case "string":
        return "\"\""
    case "double":
        return "0"
    }

    return "nil"
}

func (this*LuaHelper)getClassFieldComment(field * ClassField) string {
    var str =""
    for _,v:= range field.Comment  {
        str+="--"
        for _,v1:=range v.Tokens{
            str+=v1.Value+" "
        }
    }
    return str
}

func (this*LuaHelper)getEnumFieldComment(field * EnumField) string {
    var str =""
    for _,v:= range field.Comment  {
        str+="--"
        for _,v1:=range v.Tokens{
            str+=v1.Value+" "
        }
    }
    return str
}

func GenLua(parser * SPParser,outPath string) error {

    if has, _ := PathExists(outPath); !has {
        return errors.New(fmt.Sprintf("生成lua，目录%s 不存在，或者出错!\n", outPath))
    }

    outfileName := parser.fileName +".lua"
    fmt.Println("开始生成lua文件:",outfileName)

    helper:=&LuaHelper{}
    helper.parser = parser

    funcMap := template.FuncMap{
        "GetEnumFieldComment"			:helper.getEnumFieldComment,
        "GetClassFieldComment"			:helper.getClassFieldComment,
        "GetClassFieldType"				:helper.getClassFieldType,
        "GetFieldDefaultValue"          :helper.getFieldDefaultValue,
    }

    tpl, err := template.New("genLua").Funcs(funcMap).Parse(toluaTemplate)
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