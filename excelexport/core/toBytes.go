package core

import (
    "bytes"
    "io/ioutil"
    "strconv"
    "encoding/binary"
    "errors"
    "fmt"
    "strings"
    "os"
)

func writeStrToBuf(buf*bytes.Buffer,str string){

    str = strings.TrimSpace(str)
    strbytes :=[]byte(str)
    binary.Write(buf, binary.BigEndian, uint16(len(strbytes)))
    buf.Write(strbytes)

}

//导出bytes数据
func GenBytes( dataSheets []*DataSheet,outpath string) error {

   if b,_:=PathExists(outpath);!b{
       os.Mkdir(outpath,0777)
   }

    for _, dataSheet := range dataSheets {

        fmt.Printf("导出表到二进制文件:%s\n",dataSheet.Name)

        var bytebuf bytes.Buffer
        for i,row := range dataSheet.Data{
            for j:=0;j< dataSheet.ColCount;j++{
                if dataSheet.Infos[j].IsExport{
                    switch dataSheet.Infos[j].TypeName {

                    case "text":
                        writeStrToBuf(&bytebuf,row[j])

                    case "textmult":
                        writeStrToBuf(&bytebuf,row[j])

                    default:        //默认全是整数
                        rv := row[j]
                        if IsStrNull(rv){
                            rv = "0"
                        }

                        if v,err:= strconv.Atoi(rv);err==nil{
                            binary.Write(&bytebuf, binary.BigEndian, int32(v))
                        }else {
                            //尝试转换成浮点数，在转成整数
                            if v2,err:= strconv.ParseFloat(rv,64);err==nil {
                                binary.Write(&bytebuf, binary.BigEndian, int32(v2))
                            }else {
                                return errors.New(fmt.Sprintf("写入二进制错误,不能转换到整数,表:%s ,行:%d,列:%s,值:%s . %s",dataSheet.Name,i+6,dataSheet.Infos[j].Name,row[j],err.Error()))
                            }
                        }
                    }
                }
            }
        }

        ioutil.WriteFile(outpath+"/"+dataSheet.Name+"Bean.bytes", bytebuf.Bytes(), 0666)
    }

    return nil
}