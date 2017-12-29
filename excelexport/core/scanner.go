package core

import (
    "io/ioutil"
    "errors"
    "fmt"
    "strings"
    "github.com/tealeg"
    "os"
)

func ReadAllExcel(path string) ([]*Table, error){

    if dir, err := ioutil.ReadDir(path);err==nil {

        tables := make([]*Table, 0, len(dir))

        for _, v := range dir {

            if !v.IsDir() && !strings.Contains(v.Name(), "~$") {

                if newfile, err1 := xlsx.OpenFile(path + "/" + v.Name()); err1 == nil {

                    //读取到table
                    table := &Table{
                        Name:  v.Name(),
                        Sheet: make([]*Sheet, 0, len(newfile.Sheet)),
                    }

                   for name,sheet:=range newfile.Sheet{

                       rows:=make([][]string,sheet.MaxRow)

                       for i,row:=range sheet.Rows{

                           rows[i] = make([]string,sheet.MaxCol)

                           for j, cell := range row.Cells {
                               rows[i][j] = cell.Value
                           }

                           //fmt.Println(len(rows[i]),name,rows[i],"\n")
                       }

                        newsheet := &Sheet{
                            Name:     name,
                            Data:     rows,
                            RowCount: sheet.MaxRow,
                            ColCount: sheet.MaxCol,
                        }

                        //fmt.Println(newfile.GetSheetName(i),"  ",i,"\n")

                        table.Sheet = append(table.Sheet, newsheet)
                    }

                    tables = append(tables,table)

                    fmt.Printf("读取文件: %s\n", v.Name())

                } else {
                    return nil, err1
                }
            }
        }

        return tables, nil
    }

    return nil,errors.New(fmt.Sprintf("读取目录出错%d\n",path))
}


func GetDataSheetInfo(tables []*Table) ([]*DataSheet){

    dss:=make([]*DataSheet,0)

    for _, table := range tables {

        for _,sheet :=range table.Sheet{

            //只有数据表才导出
            if strings.Contains(sheet.Name,"t_"){

                ds:=&DataSheet{
                    Name:sheet.Name,
                    RowCount:sheet.RowCount,
                    ColCount:sheet.ColCount,
                    Infos:make([]*DataSheetColInfo,0,sheet.ColCount),
                    Data:sheet.Data[5:],
                }

                for i:=0;i<sheet.ColCount;i++ {

                    dsci:= &DataSheetColInfo{
                        IsKey:     strings.Compare(sheet.Data[0][i], "1") == 0,
                        IsExport:  !IsStrNull(sheet.Data[1][i]),
                        Name:      sheet.Data[1][i],
                        TypeName:  sheet.Data[2][i],
                        Describe1: sheet.Data[3][i],
                        Describe2: sheet.Data[4][i],
                    }

                    //fmt.Println(*dsci)

                    ds.Infos = append(ds.Infos,dsci)
                }
                dss = append(dss,ds)
            }
        }
    }
    return dss
}

func  IsStrNull(s string) bool{

    if len(s)==0{
        return true
    }
    for _,c:=range s{
        if c!=' '{
            return false
        }
    }
    return true
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}