package core

import "database/sql"
import (
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "strings"
    "strconv"
    "sync"
)


func GenDB(dataSheets []*DataSheet, dsn string) error {
    db, err := sql.Open("mysql", dsn)
    db.SetMaxOpenConns(20)
    if err!=nil{
        return err
    }

    if err = db.Ping();err!=nil{
        return err
    }

    //多线程导入数据库
    var waitgroup sync.WaitGroup

     var retErr error

    for _,sheet0 :=range dataSheets {
        sheet := sheet0
        go func() {

            waitgroup.Add(1)
            defer waitgroup.Done()

            fmt.Println("开始往数据库导入表:" + sheet.TableName)
            defer fmt.Println("结束往数据库导入表:" + sheet.TableName)


            //删除老表
            _, err = db.Exec("DROP TABLE IF EXISTS " + sheet.Name)
            if err != nil {
                fmt.Println("删除老表出错:" + sheet.TableName + " " + sheet.Name + " "+ err.Error())
                retErr = err
                return
            }
            //创建新表
            //组合表字段
            ctstr := ""
            insertstr := ""
            for i, info := range sheet.Infos {
                name := removeEndSapceChar(info.Name)
                if info.IsExport {
                    if i == 0 {
                        insertstr += "`" + name + "`"
                    } else {
                        insertstr += ",`" + name + "`"
                    }
                    if strings.Contains(info.SrcTypeName, "text") || info.SrcTypeName == "long" {
                        ctstr += "`" + name + "` text,"
                    } else
                    {
                        ctstr += "`" + name + "` int,"
                    }
                }
            }

            createsqlstr := fmt.Sprintf("CREATE TABLE `%s`(%s PRIMARY KEY (t_id))CHARACTER SET utf8", sheet.Name, ctstr)
            _, err = db.Exec(createsqlstr)
            if err != nil {
                fmt.Println(createsqlstr)
                fmt.Println("创建新表出错:" + sheet.TableName + " " + sheet.Name+ " "+ err.Error())
                retErr = err
                return
            }

            //写入数据

            insertsql := "INSERT INTO `" + sheet.Name + "`(%s) VALUES (%s)"

            for i, datas := range sheet.Data {

                datastr := ""

                for j, data := range datas {
                    if !sheet.Infos[j].IsExport {
                        continue
                    }

                    s0 := ""

                    if sheet.Name == "t_trainStep" && sheet.Infos[j].Name == "actionIds" {
                        s0 = ""
                    }

                    if strings.Contains(sheet.Infos[j].SrcTypeName, "text") || sheet.Infos[j].SrcTypeName == "long" {
                        s0 = strconv.Quote(data)
                    } else {
                        if IsStrNull(data) {
                            data = "0"
                        }
                        s0 = data
                    }

                    if j == 0 {
                        datastr += s0
                    } else {
                        datastr += "," + s0
                    }
                }

                strings.Replace(datastr, ",,", ",", -1)
                //插入数据
                s := fmt.Sprintf(insertsql, insertstr, datastr)
                _, err = db.Exec(s)
                if err != nil {
                    fmt.Println(s)
                    fmt.Println("插入数据出错:" + sheet.TableName + " " + sheet.Name + " " + strconv.Itoa(i+5)+ " "+ err.Error())
                    retErr = err
                    return
                }
            }
        }()
    }

    defer db.Close()

    waitgroup.Wait()

    return retErr
}