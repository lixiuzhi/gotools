package core

import (
    "os"
    "fmt"
    "bytes"
    "io/ioutil"
    "text/template"
    "strings"
    "strconv"
)

type luaHelper struct{}

//参数 sheet 列
func (this*luaHelper) IsNotRowEnd(sheet*DataSheet,col int) bool{

    i:=len(sheet.Infos)-1
    for ;i>=0;i--{
        if sheet.Infos[i].IsExport{
            break
        }
    }

    if col+1<=i{
        return true
    }
    return false
}

//参数 sheet 行
func (this*luaHelper) IsNotLastRow(sheet*DataSheet,row int) bool{

    i:=len(sheet.Data)-1

    if row+1<=i{
        return true
    }
    return false
}

//参数 sheet 行 列
func (this*luaHelper) GetFieldValue(sheet*DataSheet,row,col int) string{

    value:=sheet.Data[row][col]
    if strings.Contains(sheet.Infos[col].SrcTypeName,"text") {
        return strconv.Quote(value)
    }

    if IsStrNull(value){
        return "0"
    }

    return value
}

func GenLua(dataSheets []*DataSheet,  outpath string) error {

    //初始化目录
    CreateDir(outpath + "/design")

    //删除不存在的配置表文件
    files, err := WalkDir(outpath + "/design", ".lua")
    if err!=nil{
        return err
    }
    if len(files) > 0 {
        for _, value := range files {
                os.Remove(value)
        }
    }

    helper := &luaHelper{}
    funcMap := template.FuncMap{
        "GetFieldValue" :   helper.GetFieldValue,
        "IsNotRowEnd"   :   helper.IsNotRowEnd,
        "IsNotLastRow"  :   helper.IsNotLastRow,
    }

    luatpl, err := template.New("genLua").Funcs(funcMap).Parse(toLuaTableTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    luaglobaltabletpl, err := template.New("genLuaGlobal").Parse(toLuaGlobalTableTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    for _, dataSheet := range dataSheets {

        //生成bean
        var tmpBuf bytes.Buffer
        err = luatpl.Execute(&tmpBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }
        ioutil.WriteFile(outpath+"/design/"+dataSheet.Name+".lua", tmpBuf.Bytes(), 0666)
    }

    //生成全局表
    var tmpBuf1 bytes.Buffer
    err = luaglobaltabletpl.Execute(&tmpBuf1,&struct {}{})
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    ioutil.WriteFile(outpath+"/CfgTable.lua", tmpBuf1.Bytes(), 0666)

    return nil
}