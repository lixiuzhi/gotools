package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func writeStrToBuf(buf *bytes.Buffer, str string) {

	str = strings.TrimSpace(str)
	strbytes := []byte(str)
	binary.Write(buf, binary.BigEndian, uint16(len(strbytes)))
	buf.Write(strbytes)

}

//导出bytes数据
func GenBytes(dataSheets []*DataSheet, outpath string) error {

	cs:=&csharpHelper{}

	if b, _ := PathExists(outpath); !b {
		os.Mkdir(outpath, 0777)
	}

	for _, dataSheet := range dataSheets {

		fmt.Printf("导出表到二进制文件:%s\n", dataSheet.Name)

		var bytebuf bytes.Buffer
		for i, row := range dataSheet.Data {
			for j := 0; j < dataSheet.ColCount; j++ {
				if dataSheet.Infos[j].IsExport {
					if  cs.checkIsModifyType(dataSheet.Infos[j]) {
						writeStrToBuf(&bytebuf, row[j])
					} else {
						rv := row[j]
						if IsStrNull(rv) {
							rv = "0"
						}
						switch dataSheet.Infos[j].TypeName {
						case "long":
							if v, err := strconv.ParseInt(rv, 10,64) ;err == nil{
								binary.Write(&bytebuf, binary.BigEndian, int64(v))

								//var bbb bytes.Buffer
								//binary.Write(&bbb, binary.BigEndian, int64(v))
								//ddddd,_ :=ioutil.ReadAll(&bbb)
								//fmt.Println("dddddddddddd:",ddddd)
							} else {
								return errors.New(fmt.Sprintf("写入二进制错误,不能转换到长整形,表:%s ,行:%d,列:%s,值:%s . %s", dataSheet.Name, i+6, dataSheet.Infos[j].Name, row[j], err.Error()))
							}
						case "bool":
							if v, err := strconv.ParseBool(rv) ;err == nil{

								//var bbb bytes.Buffer
								binary.Write(&bytebuf, binary.BigEndian, v)
								//ddddd,_ :=ioutil.ReadAll(&bbb)
								//fmt.Println("dddddddddddd:",rv,ddddd)
							} else {
								return errors.New(fmt.Sprintf("写入二进制错误,不能转换到bool,表:%s ,行:%d,列:%s,值:%s . %s", dataSheet.Name, i+6, dataSheet.Infos[j].Name, row[j], err.Error()))
							}
						case "text":
							writeStrToBuf(&bytebuf, row[j])
						case "textmult":
							writeStrToBuf(&bytebuf, row[j])
						default: //默认全是整数

							if dataSheet.Infos[j].TypeName == "float" {
								fmt.Print("")
							}

							if v, err := strconv.Atoi(rv); err == nil {
								binary.Write(&bytebuf, binary.BigEndian, int32(v))
							} else {
								//尝试转换成浮点数，在转成整数
								if v2, err := strconv.ParseFloat(rv, 64); err == nil {
									binary.Write(&bytebuf, binary.BigEndian, int32(v2))
								} else {
									return errors.New(fmt.Sprintf("写入二进制错误,不能转换到整数,表:%s ,行:%d,列:%s,值:%s . %s", dataSheet.Name, i+6, dataSheet.Infos[j].Name, row[j], err.Error()))
								}
							}
						}
					}

				}
			}
		}

		ioutil.WriteFile(outpath+"/"+dataSheet.Name+"Bean.bytes", bytebuf.Bytes(), 0666)
	}
	fmt.Print("\n************************************导出表到二进制文件 完毕\n\n")

	return nil
}

////检查 返回类型是否是修改后的类型
//func checkIsModifyType(typeName string,commondstr string) bool {
//	return checkIsSplitClass(typeName) || CheckIsSplitCommond(commondstr)
//}
//
////检查是否为自定义类
//func checkIsSplitClass(typeName string) bool {
//
//	return strings.Contains(typeName,"Vector2",) || strings.Contains(typeName,"Vector3",)
//}
//
//
//func CheckIsSplitCommond(commonds string) bool {
//	return  strings.Contains(commonds,"split")
//}
