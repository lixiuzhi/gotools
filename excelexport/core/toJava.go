package core

import (
    "bytes"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "text/template"
    "strings"
    "unicode"
)

type javaHelper struct{

}

func (java*javaHelper)getTypeName(colinfo *DataSheetColInfo) string {

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

func (java*javaHelper)needConvertToLong(colinfo *DataSheetColInfo) bool{
    if colinfo.SrcTypeName=="long"{
        return true
    }
    return false
}

func (java*javaHelper)getTypeNameWithDB(colinfo *DataSheetColInfo) string {

    switch colinfo.SrcTypeName {
    case "int":
        return "int"
    case "long":
        return "String"
    case "text":
        return "String"
    case "textmult":
        return "String"
    default:
        return "int"
    }
}

func (java*javaHelper)getDBType(colinfo *DataSheetColInfo) string {

    switch colinfo.SrcTypeName {
    case "int":
        return "INTEGER"
    case "long":
        return "LONGVARCHAR"
    case "text":
        return "LONGVARCHAR"
    case "textmult":
        return "LONGVARCHAR"
    default:
        return "INTEGER"
    }
}


func (java*javaHelper)getInfoUpperName(colinfo *DataSheetColInfo) string{
    s:=colinfo.Name
    s=s[:1]
    s= strings.Map(unicode.ToUpper, s) + colinfo.Name[1:]

    return s
}

func (java*javaHelper)getIdUpperName(sheet *DataSheet) string{

    return java.getInfoUpperName(sheet.Infos[0])
}


func (java*javaHelper)IsContainItem(dataSheets []*DataSheet, path string) bool {

    for _, value1 := range dataSheets {
        if strings.Contains(path, value1.Name+"Bean.java") {
            return true
        } else if strings.Contains(path, value1.Name+"Container.java") {
            return true
        } else if strings.Contains(path, value1.Name+"Dao.java"){
            return true
        }else if strings.Contains(path, value1.Name+".xml"){
            return true
        }
    }

    return false
}


//生成CS代码
func GenJava(dataSheets []*DataSheet, outpath string) error {

     fmt.Println("开始生成java代码..")

    java :=&javaHelper{}

    if b, _ := PathExists(outpath); !b {
        return errors.New("生成java，目录不存，或者出错!\n")
    }
    //初始化目录
    CreateDir(outpath + "/bean")
    CreateDir(outpath + "/container")
    CreateDir(outpath + "/dao")
    CreateDir(outpath + "/sqlmap")

    //删除不存在的配置表文件
    files, err := WalkDir(outpath, ".java")
    if len(files) > 0 {
        for _, value := range files {
            if !java.IsContainItem(dataSheets, value) {
                os.Remove(value)
                //fmt.Print("remove file:" + value + "\n")
            } else {

            }
        }
    }

    files, err = WalkDir(outpath, ".xml")
    if len(files) > 0 {
        for _, value := range files {
            if !java.IsContainItem(dataSheets, value) {
                os.Remove(value)
                //fmt.Print("remove file:" + value + "\n")
            } else {

            }
        }
    }

    funcMap := template.FuncMap{
        "GetTypeName":      java.getTypeName,
        "GetInfoUpperName": java.getInfoUpperName,
        "GetIdUpperName":   java.getIdUpperName,
        "GetDBType":        java.getDBType,
        "GetTypeNameWithDB":java.getTypeNameWithDB,
        "NeedConvertToLong":java.needConvertToLong,
    }

    beantpl, err := template.New("genJavaBean").Funcs(funcMap).Parse(javaBeanTemplateStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    containertpl, err := template.New("genJavaContainer").Funcs(funcMap).Parse(javaContainerTemplateStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    daotpl, err := template.New("genJavaDao").Funcs(funcMap).Parse(javaDaoTemplateStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    sqlmaptpl, err := template.New("genJavaSqlMap").Funcs(funcMap).Parse(javaSqlMapTemplateStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    dbgamedataconfigtpl, err := template.New("genJavaDBGameDataConfig").Funcs(funcMap).Parse(javaDBGameDataConfigTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    gamedatamanagertpl, err := template.New("genJavaGameDataManager").Funcs(funcMap).Parse(javaGameDataManagerTemplate)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    for _, dataSheet := range dataSheets {

        //生成bean
        var beanBuf bytes.Buffer
        err = beantpl.Execute(&beanBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }
        ioutil.WriteFile(outpath+"/bean/"+dataSheet.Name+"Bean.java", beanBuf.Bytes(), 0666)


        //生成container
        var containerBuf bytes.Buffer
        err = containertpl.Execute(&containerBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }
        ioutil.WriteFile(outpath+"/container/"+dataSheet.Name+"Container.java", containerBuf.Bytes(), 0666)

        //生成dao
        var daoBuf bytes.Buffer
        err = daotpl.Execute(&daoBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }
        ioutil.WriteFile(outpath+"/dao/"+dataSheet.Name+"Dao.java", daoBuf.Bytes(), 0666)

        //生成sqlmap
        var sqlMapBuf bytes.Buffer
        err = sqlmaptpl.Execute(&sqlMapBuf, dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }
        ioutil.WriteFile(outpath+"/sqlmap/"+dataSheet.Name+".xml", sqlMapBuf.Bytes(), 0666)
    }

    //生成dbgamedataconfig
    var dbgamedataconfigBuf bytes.Buffer
    err = dbgamedataconfigtpl.Execute(&dbgamedataconfigBuf, dataSheets)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    ioutil.WriteFile(outpath+"/db-game-data-config.xml", dbgamedataconfigBuf.Bytes(), 0666)

    //生成gamedatamanager
    var gamedatamanagerBuf bytes.Buffer
    err = gamedatamanagertpl.Execute(&gamedatamanagerBuf, dataSheets)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    ioutil.WriteFile(outpath+"/GameDataManager.java", gamedatamanagerBuf.Bytes(), 0666)

    fmt.Print("\n************************************生成java文件 完毕\n\n")

    return nil
}
