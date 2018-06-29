package core

import (
    "os"
    "fmt"
    "bytes"
    "io/ioutil"
    "text/template"
)

type luaAPIHelper struct{}

func (this*luaAPIHelper)getTypeName(colinfo *DataSheetColInfo) string {

    switch colinfo.SrcTypeName {
    case "int":
        return "int"
    case "long":
        return "long"
    case "text":
        return "String"
    case "textmult":
        return "String"
    default:
        return "int"
    }
}

func GenLuaAPI(dataSheets []*DataSheet,outpath string) error {

    //初始化目录
    CreateDir(outpath)

    //删除不存在的配置表文件
    files, err := WalkDir(outpath, ".lua")
    if err != nil {
        return err
    }
    if len(files) > 0 {
        for _, value := range files {
            os.Remove(value)
        }
    }

    helper := &luaAPIHelper{}
    funcMap := template.FuncMap{
        "GetTypeName": helper.getTypeName,
    }

    luaAPItpl, err := template.New("genLuaAPI").Funcs(funcMap).Parse(toLuaAPITableTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    luaapiglobaltpl, err := template.New("genLuaAPIGlobal").Funcs(funcMap).Parse(toLuaAPIGlobalTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    for _, dataSheet := range dataSheets {

        //生成bean
        var tmpBuf bytes.Buffer
        err = luaAPItpl.Execute(&tmpBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        ioutil.WriteFile(outpath+"/"+dataSheet.Name+".lua", tmpBuf.Bytes(), 0666)
    }

    //生成GlobalAPI

    var tmpBuf bytes.Buffer
    err = luaapiglobaltpl.Execute(&tmpBuf, dataSheets)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    ioutil.WriteFile(outpath+"/"+"CfgTable.lua", tmpBuf.Bytes(), 0666)
    return nil
}