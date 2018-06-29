package core

import (
	"errors"
	"fmt"
	"github.com/tealeg"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

func ReadAllExcel(path string) ([]*Table, error) {

	if dir, err := ioutil.ReadDir(path); err == nil {

		tables := make([]*Table, 0, len(dir))

		for _, v := range dir {

			if !v.IsDir() && !strings.Contains(v.Name(), "~$") && strings.Contains(v.Name(),".xlsx") {

				if newfile, err1 := xlsx.OpenFile(path + "/" + v.Name()); err1 == nil {

					//读取到table
					table := &Table{
						Name:  v.Name(),
						Sheet: make([]*Sheet, 0, len(newfile.Sheet)),
					}

					for name, sheet := range newfile.Sheet {

						rows := make([][]string, sheet.MaxRow)

						for i, row := range sheet.Rows {

							rows[i] = make([]string, sheet.MaxCol)

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

					tables = append(tables, table)

					fmt.Printf("读取文件: %s\n", v.Name())

					if err:= CheckErrorData(table);err!=nil{
						return nil,err
					}

				} else {
					return nil, err1
				}
			}
		}

		return tables, nil
	}

	return nil, errors.New(fmt.Sprintf("读取目录出错%d\n", path))
}

func GetDataSheetInfo(tables []*Table) ([]*DataSheet,error) {

	dss := make([]*DataSheet, 0)

	for _, table := range tables {

		for _, sheet := range table.Sheet {

			//只有数据表才导出
			if strings.Contains(sheet.Name, "t_") {

				ds := &DataSheet{
					TableName:table.Name,
					Name:     sheet.Name,
					RowCount: sheet.RowCount,
					ColCount: sheet.ColCount,
					Infos:    make([]*DataSheetColInfo, 0, sheet.ColCount),
					Data:     sheet.Data[5:],
				}

				for i := 0; i < sheet.ColCount; i++ {

					dsci := &DataSheetColInfo{
						IsKey:     		strings.Compare(sheet.Data[0][i], "1") == 0,
						IsExport:  		!IsStrNull(sheet.Data[1][i]),
						Name:      		sheet.Data[1][i],
						SrcTypeName:	sheet.Data[2][i],
						TypeName:  		sheet.Data[2][i],
						Describe1: 		sheet.Data[3][i],
						Describe2: 		sheet.Data[4][i],
						Commonds:  		getCommonds(sheet.Data[3][i]),
					}

					//修正类型（所有修正类型的原始类型必须为string）
					for _, value := range dsci.Commonds {
						if value.Name == "type" {
							dsci.TypeName = value.ParamStr
							break
						}
					}

					//fmt.Println(*dsci)

					ds.Infos = append(ds.Infos, dsci)

					ds.indexs = make(map[string] map[string] int)
				}
				dss = append(dss, ds)
			}
		}
	}

	//建立索引
	for _,sheet := range dss{
		for i,info:=range sheet.Infos{
			if info.IsKey{
				 keymap := make(map[string] int)
				 for _,datas:=range sheet.Data{
					 keymap[datas[i]] = i
				 }
				 sheet.indexs[info.Name] = keymap
			}
		}
	}

	err:=CheckConnectIndex(dss)

	return dss,err
}

//检查关联的id是否存在
func CheckConnectIndex(sheets []*DataSheet) error{

	var sheetmap = make(map[string] *DataSheet)

	for _,sheet:=range sheets{
		sheetmap[sheet.Name] = sheet
	}

	for _,sheet:=range sheets{
		for i,info:=range sheet.Infos{
			if info.Commonds!=nil{
				for _,cmd:=range info.Commonds  {
					if cmd.Name=="connect"{
						if len(cmd.paras)!=2{
							return errors.New("表命令参数错误:"+sheet.Name+" 列:"+ strconv.Itoa(i) + "  参数:"+ cmd.ParamStr)
						}else {
							sheetName 	:= cmd.paras[0]
							colName 	:= cmd.paras[1]
							if connSheet, ok := sheetmap[sheetName]; ok {
								if vmap, ok := connSheet.indexs[colName]; ok {
									//遍历查找是否有不存在的
									for j,data:=range sheet.Data{

										if IsStrNull(data[i]){
											continue
										}

										if _,ok:=vmap[data[i]];!ok{
											return errors.New("关联条目不存在,当前表:"+sheet.Name+",当前列:"+strconv.Itoa(i+1)+",当前行:"+ strconv.Itoa(j+6) +",值:"+data[i] +",关联表:"+connSheet.Name+",关联字段:"+colName)
										}
									}
								} else {
									return errors.New("关联命令没有找到对应的sheet列,表:" + sheet.Name + ",列:" + strconv.Itoa(i+1) + ",参数:" + cmd.ParamStr)
								}
							} else
							{
								return errors.New("关联命令没有找到对应的sheet,表:" + sheet.Name + ",列:" + strconv.Itoa(i+1) + ",参数:" + cmd.ParamStr)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func CheckErrorData( table *Table) error {

	for _, sheet := range table.Sheet{

	    if !strings.Contains(sheet.Name,"t_"){
	    	continue
		}

		if len(sheet.Data) == 0 {
			return nil
		}
		//检查第一列是否存在空，存在则报错
		tmpMap := make(map[string]int)
		for i, data := range sheet.Data {
			if i<=4{
				continue
			}
			s := data[0]
			isnull := true
			for _, c := range s {
				if c != ' ' {
					isnull = false
					break
				}
			}
			if isnull {
				return errors.New(fmt.Sprintf("检测表格数据出错,表 %s中的sheet: %s的第%d行的第一列为空!", table.Name, sheet.Name, i+1))
			}
			v, ok := tmpMap[s]
			if ok {
				return errors.New(fmt.Sprintf("检测表格数据出错,表 %s中的sheet %s的第%d行与第%d行id重复!", table.Name, sheet.Name, i+1, v))
			}
			tmpMap[s] = i + 1
		}
	}

	return nil
}

func getCommonds(s string) []*DataSheetColCommond {
	str := strings.TrimSpace(s)
	if IsStrNull(str) {
		return nil
	} else {
		comms := make([]*DataSheetColCommond, 0)
		var arr = strings.Split(str, "|")
		for _, value := range arr {
			length := strings.Index(value,":")
			//var item = strings.Split(value, ":")

			if length >= 0 {
				paramStr := SubString(value, length + 1, strings.Count(value,""))
				com := &DataSheetColCommond{
					Name:     SubString(value, 0, length),
					paras:    strings.Split(paramStr, "+"),
					ParamStr: paramStr,
				}
				comms = append(comms, com)
			}
		}

		return comms
	}
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
