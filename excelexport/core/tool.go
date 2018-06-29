package core

import (
    "fmt"
    "strings"
    "path/filepath"
    "os"
)

func SubString(str string, begin, length int) (substr string) {
    // 将字符串的转换成[]rune
    rs := []rune(str)
    lth := len(rs)

    // 简单的越界判断
    if begin < 0 {
        begin = 0
    }
    if begin >= lth {
        begin = lth
    }
    end := begin + length
    if end > lth {
        end = lth
    }

    // 返回子串
    return string(rs[begin:end])
}

func CreateDir(name string) bool {
    if IsDir(name) {
        //fmt.Printf("%s is already a directory.\n", name)
        return true
    }

    if createDirImpl(name) {
        fmt.Println("Create directory successfully:" + name)
        return true
    } else {
        return false
    }
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
    files = make([]string, 0, 30)
    suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
    err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
        //if err != nil { //忽略错误
        // return err
        //}
        if fi.IsDir() { // 忽略目录
            return nil
        }
        if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
            files = append(files, filename)
        }
        return nil
    })
    return files, err
}

func IsDir(name string) bool {
    fi, err := os.Stat(name)
    if err != nil {
        //fmt.Println("Error: ", err)
        return false
    }

    return fi.IsDir()
}

func createDirImpl(name string) bool {
    err := os.MkdirAll(name, 0666)
    if err == nil {
        return true
    } else {
        fmt.Println("Error: ", err)
        return false
    }
}

func removeEndSapceChar(src string) string{

    if len(src)==0{
        return src
    }

    for i:= len(src)-1;i>=0;i--{
        if src[i]!=' '{
            return src[0:i+1]
        }
    }

    return src
}

func IsStrNull(s string) bool {

    if len(s) == 0 {
        return true
    }
    for _, c := range s {
        if c != ' ' {
            return false
        }
    }
    return true
}