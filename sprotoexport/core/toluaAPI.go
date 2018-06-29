package core

import (
    "errors"
    "fmt"
    "bytes"
    "io/ioutil"
    "text/template"
)

type LuaAPIHelper struct {
    parser *SPParser
}

func (this*LuaAPIHelper)getClassFieldType(field * ClassField) string {

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

    if field.Repeatd{
        str = baseType+"[]"
    }else{
        str = baseType
    }

    return str
}

func (this*LuaAPIHelper)hasClassFieldComment(field * ClassField) bool {
    return len(field.Comment)>0
}

func (this*LuaAPIHelper)getClassFieldComment(field * ClassField) string {
    var str =""
    for _,v:= range field.Comment  {
        str+="  "
        for _,v1:=range v.Tokens{
            str+=v1.Value+" "
        }
    }
    return str
}

func GenLuaAPI(parser * SPParser,outPath string) error {

    if has, _ := PathExists(outPath); !has {
        return errors.New(fmt.Sprintf("生成lua emmylua api，目录%s 不存在，或者出错!\n", outPath))
    }

    outfileName := parser.fileName +".lua"
    fmt.Println("开始生成luaAPI文件:",outfileName)

    helper:=&LuaAPIHelper{}
    helper.parser = parser

    funcMap := template.FuncMap{
        "GetClassFieldComment"			:helper.getClassFieldComment,
        "GetClassFieldType"				:helper.getClassFieldType,
        "hasClassFieldComment"          :helper.hasClassFieldComment,
    }

    tpl, err := template.New("genLuaAPI").Funcs(funcMap).Parse(toluaAPITemplate)
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