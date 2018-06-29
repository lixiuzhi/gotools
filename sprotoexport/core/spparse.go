package core

import (
	"fmt"
	"errors"
	"strings"
	"reflect"
	"encoding/json"
	"strconv"
)

type CommentField struct {
	Tokens 		[]*TokenInfo
	Index 		int
}

func(this * CommentField) String() string{
	str:=""
	for _,v:=range this.Tokens{
		str+=v.Value+" "
	}
	return str
}

type ClassField struct{
	Index 		int
	Comment 	[]*CommentField
	Name 		string
	Type		string
	TypeIsEnum	bool
	TypeIsClass	bool
	Repeatd		bool
	LocalIndex	int
	FileName	string
}

type EnumField struct {
	Index 		int
	Comment 	[]*CommentField
	Name 		string
	Value		int
	LocalIndex	int
	FileName	string
}

type EnumType struct {
	Comment 	[]*CommentField
	Fields		[]*EnumField
	Index 		int
	EndIndex	int
	Name		string
}

type ClassType struct {
	Comment 	[]*CommentField
	Fields		[]*ClassField
	Index 		int
	EndIndex	int
	Name		string
	AutoMsgID	int
	MsgID		int
}

type SPParser struct {

	fileName	string
	maxLine		int
	allComments []*CommentField
	allField 	map[int]interface{}

	Enums 		map[int]*EnumType
	Classes 	map[int]*ClassType


	EnumsByLineID 	map[int]*EnumType
	ClassesByLineID map[int]*ClassType
}

type ClassCommentMetaInfo struct {
	Uid	int
}

func (this*SPParser) GetEnumByName(name string) *EnumType{

	for _,v:=range this.Enums{
		if v.Name==name{
			return v
		}
	}
	return nil
}

func (this*SPParser) PrintfAll()  {

	for _,v:= range this.allComments{
		fmt.Printf("注释 %d:%s\n",v.Index,v.String())
	}

	for _,v:= range this.allField{

		if rv,ok:=v.(*EnumField); ok{
			fmt.Printf("枚举Field %s:%d\n",rv.Name,rv.Index)
		}
		if rv,ok:=v.(*ClassField);ok{
			fmt.Printf("类Field %s:%d repeat:%t type:%s\n",rv.Name,rv.Index,rv.Repeatd,rv.Type)
		}
	}

	fmt.Println("*******************************************************\n")

	for _,v:= range this.Classes{
		//fmt.Printf("类 %d:%d\n",v.Index,v.EndIndex)
		fmt.Println(*v)
	}

	for _,v:= range this.Enums{
		//fmt.Printf("枚举 %d:%d\n",v.Index,v.EndIndex)
		fmt.Println(*v)
	}
}

func (this*SPParser) fill() error{

	//填充注释区域
	for _, c := range this.allComments {

		//如果当前行不是则一直查找下一行
		for i := c.Index; i < this.maxLine; i++ {

			//判断当前行是否是field，枚举 或者类
			if v, ok := this.EnumsByLineID[i]; ok {
				v.Comment = append(v.Comment, c)
				break
			}

			if v, ok := this.ClassesByLineID[i]; ok {
				v.Comment = append(v.Comment, c)
				break
			}

			if v, ok := this.allField[i]; ok {

				if rv, ok := v.(*EnumField); ok {
					rv.Comment = append(rv.Comment, c)
					break
				}

				if rv, ok := v.(*ClassField); ok {
					rv.Comment = append(rv.Comment, c)
					break
				}

				break
			}
		}
	}


	enumNameToEnum:= map[string]*EnumType{}

	//填充enum区域
	for _, enum := range this.Enums {
		enumNameToEnum[enum.Name] = enum
		for i := enum.Index; i <= enum.EndIndex; i++ {

			if f, ok := this.allField[i]; ok {

				if rf,ok:=f.(*EnumField);ok{
					rf.LocalIndex = len(enum.Fields)+1
					enum.Fields = append(enum.Fields, rf)

					this.allField[i] = nil
				}else {
					return errors.New(fmt.Sprintf("组合枚举时错误，无效的filed,行：%d\n",i+1))
				}
			}
		}
	}

	classNameToClass:= map[string]*ClassType{}
	for _, class := range this.Classes {
		classNameToClass[class.Name] = class
	}

	//填充class区域
	for _, class := range this.Classes {

		for i := class.Index; i <= class.EndIndex; i++ {

			if f, ok := this.allField[i]; ok {

				if rf,ok:=f.(*ClassField);ok{
					rf.LocalIndex = len(class.Fields)+1
					class.Fields = append(class.Fields, rf)
					
					this.allField[i] = nil

					if _,ok:=enumNameToEnum[rf.Type];ok{
						rf.TypeIsEnum = true
					}else {
						rf.TypeIsEnum = false
					}

					if _,ok:=classNameToClass[rf.Type];ok{
						rf.TypeIsClass = true
					}else {
						rf.TypeIsClass = false
					}

				}else {
					return errors.New(fmt.Sprintf("组合类时错误，无效的filed,%s,行：%d\n",reflect.TypeOf(f).Name(),i+1))
				}
			}
		}
	}



	return nil
}

func (this*SPParser) Parse(tokenInfos []*TokenInfo,fileName string) error{

	fileName = strings.Replace(fileName,".sp","",1)

	fmt.Println("开始处理tokens")

	this.allComments 	= make([]*CommentField, 0, 100)
	this.allField 		= make(map[int]interface{})
	this.Enums 			= make(map[int]*EnumType)
	this.Classes 		= make(map[int]*ClassType)
	this.EnumsByLineID	= make(map[int]*EnumType)
	this.ClassesByLineID= make(map[int]*ClassType)
	this.fileName =fileName

	count := len(tokenInfos)

	if count > 0 {
		this.maxLine = tokenInfos[count-1].LocalLine + tokenInfos[count-1].LineOffset
	} else {
		this.maxLine = 0
	}

	for i := 0; i < count; i++ {

		 //fmt.Printf("开始解析line:%d,%s\n", tokenInfos[i].LocalLine, tokenInfos[i].Value)

		switch tokenInfos[i].Value {
		case "/":
			newIndex, err := this.parseComment(tokenInfos, i)

			if err != nil {
				return errors.New("解析tokens出错:"+ err.Error())
			} else {
				i = newIndex
			}

		case "message":
			newIndex, err := this.parseClass(tokenInfos, i)

			if err != nil {
				return errors.New("解析class出错:"+ err.Error())
			} else {
				i = newIndex
			}

		case "enum":
			newIndex, err := this.parseEnum(tokenInfos, i)

			if err != nil {
				return errors.New("解析枚举出错:" + err.Error())
			} else {
				i = newIndex
			}

		case "*EOF*":
		case "}":

		default:
			newIndex, err := this.parseField(tokenInfos, i)

			if err != nil {
				return errors.New(fmt.Sprintf( "解析field出错:%s,文件:%s,行%d",err.Error(),tokenInfos[i].FileName,tokenInfos[i].LocalLine+1))
			} else {
				i = newIndex
			}
		}
	}

	if err:=this.fill();err!=nil{
		return errors.New("fill时错误：" + err.Error())
	}

	//判断是否有没有用到的filed，如果有则是有问题
	for _,v:=range this.allField {
		
		if v==nil {
			continue
		}
		
		if f:=v.(*EnumField);f!=nil{
			return errors.New( fmt.Sprintf("错误的行,文件:%s,行%d,值:%s",f.FileName,f.Index+1,f.Name))
		}

		if f:=v.(*ClassField);f!=nil{
			return errors.New( fmt.Sprintf("错误的行,文件:%s,行%d,值:%s",f.FileName,f.LocalIndex+1,f.Name))
		}
	}

	//判断所有类的field是否有没有声明的类型
	for _,c:=range this.Classes{
		for _,f:=range c.Fields{
			if !f.TypeIsEnum && !f.TypeIsClass &&f.Type != "binary" && f.Type != "int32" &&f.Type != "int64" &&f.Type != "bool" &&f.Type != "string" && f.Type!="double"{
				return errors.New( fmt.Sprintf("没有声明的字段类型,文件:%s,行%d,值:%s",f.FileName,f.Index+1,f.Type))
			}
		}
	}

	//初始化class MsgID
	for _,c:=range this.Classes{
		c.MsgID = c.AutoMsgID
		for _,comment:=range c.Comment{
			var meta = &ClassCommentMetaInfo{}
			var metastr =  comment.String()
			if err := json.Unmarshal([]byte(metastr), meta); err == nil{
				c.MsgID = meta.Uid
				fmt.Println("类:"+c.Name+" 新的msgID为:"+ strconv.Itoa( meta.Uid))
			}
		}
	}

	//判断class MsgID是否有重复
	for _,c1:=range this.Classes{
	 for _,c2:=range this.Classes{
	 	if c1!=c2 && c1.MsgID == c2.MsgID{
			return errors.New( fmt.Sprintf("消息MsgID重复,消息 %s 与消息 %s",c1.Name,c2.Name))
		}
	 }
	}

	return nil
}

func (this* SPParser) isCommentField(tokenInfos []*TokenInfo,token *TokenInfo,index int) bool {

	for i:= index ;i>0;i--{
		if tokenInfos[i].Value=="*EOF*"{
			return false
		}

		if tokenInfos[i].Value == "/" &&  i>0 &&tokenInfos[i-1].Value == "/"{
			return true
		}
	}
	 return  false
}


func(this*SPParser) getIDByName(name string) int {

	v := (int)(GetHash([]byte(name))& 0x7FFFFFFF)
	if v < 0 {
		v = -v
	}
	return v
}

func(this*SPParser) getIDByComment(comment,name string) int {

	v := (int)(GetHash([]byte(name))& 0x7FFFFFFF)
	if v < 0 {
		v = -v
	}
	return v
}

func(this*SPParser) parseClass(tokenInfos []*TokenInfo,index int) (int,error) {

	//查找左括号
	leftBracketsIndex := -1
	for j := index + 2; j < len(tokenInfos); j++ {
		if tokenInfos[j].Value == "*EOF*" {
			continue
		} else if tokenInfos[j].Value == "{" {
			leftBracketsIndex = j
		} else {
			break
		}
	}
	if leftBracketsIndex == -1 {
		return index, errors.New(fmt.Sprintf("解析类错误，没有找到左括号,文件%s,行：%d\n",tokenInfos[index].FileName, tokenInfos[index].LocalLine+1))
	}
	if len(tokenInfos) > leftBracketsIndex+1 {
		newIndex := leftBracketsIndex + 1

		//查找右括号index
		for i := index + 2; i < len(tokenInfos); i++ {

			if tokenInfos[i].Value == "}" && !this.isCommentField(tokenInfos, tokenInfos[i], i) {
				newClass := &ClassType{
					Index:		tokenInfos[index].LocalLine + tokenInfos[index].LineOffset,
					EndIndex:	tokenInfos[i].LocalLine + tokenInfos[i].LineOffset ,
					Comment:	make([]*CommentField, 0),
					Fields:		make([]*ClassField, 0),
					Name:		tokenInfos[index+1].Value,
					AutoMsgID:	this.getIDByName(tokenInfos[index+1].Value),
				}

				//判断id是否重复
				if v,ok:=this.Classes[newClass.AutoMsgID];ok{
					return newIndex,errors.New(fmt.Sprintf("解析类错误,自动生成的协议ID重复,协议%s与协议%s\n",v.Name,newClass.Name))
				}
				this.Classes[newClass.AutoMsgID] = newClass
				this.ClassesByLineID[newClass.Index] = newClass

				return newIndex, nil
			}
		}
	}
	return index, errors.New(fmt.Sprintf("解析类错误，括号可能不匹配,文件%s,行：%d\n", tokenInfos[index].FileName, tokenInfos[index].LocalLine+1))

}

func(this*SPParser) parseEnum(tokenInfos []*TokenInfo,index int) (int,error) {


	//查找左括号
	leftBracketsIndex := -1
	for j := index + 2; j < len(tokenInfos); j++ {
		if tokenInfos[j].Value == "*EOF*" {
			continue
		} else if tokenInfos[j].Value == "{" {
			leftBracketsIndex = j
		} else {
			break
		}
	}
	if leftBracketsIndex == -1 {
		return index, errors.New(fmt.Sprintf("解析枚举错误，没有找到左括号,文件%s,行：%d\n", tokenInfos[index].FileName,tokenInfos[index].LineOffset+1))
	}
	if len(tokenInfos) > leftBracketsIndex+1 {
		newIndex := leftBracketsIndex + 1

		//查找右括号index
		for i := index + 2; i < len(tokenInfos); i++ {

			if tokenInfos[i].Value == "}" && !this.isCommentField(tokenInfos, tokenInfos[i], i) {
				newEnum := &EnumType{
					Index	:tokenInfos[index].LocalLine + tokenInfos[index].LineOffset,
					EndIndex:tokenInfos[i].LocalLine + tokenInfos[i].LineOffset,
					Comment	:make([]*CommentField,0),
					Fields	:make([]*EnumField,0),
					Name	:tokenInfos[index+1].Value,
				}

				if v,ok:=this.Enums[this.getIDByName(newEnum.Name)];ok{
					return newIndex,errors.New(fmt.Sprintf("解析枚举错误,相同的id,协议%s与协议%s\n",v.Name,newEnum.Name))
				}

				this.Enums[this.getIDByName(newEnum.Name)]=newEnum //newEnum.Index
				this.EnumsByLineID[newEnum.Index] = newEnum

				return newIndex, nil
			}
		}
	}
	return index, errors.New(fmt.Sprintf("解析枚举错误，括号可能不匹配,文件%s,行：%d\n",tokenInfos[index].FileName, tokenInfos[index].LocalLine+1))

}

//分离注释区域
func(this*SPParser) parseComment( tokenInfos []*TokenInfo,index int) (int,error) {

	newIndex := index
	if tokenInfos[index].Value == "/" && len(tokenInfos) > index+1 && tokenInfos[index+1].Value == "/" {
		//产生新的comment
		var newComment = &CommentField{
			Tokens: make([]*TokenInfo, 0),
			Index:  tokenInfos[index].LocalLine + tokenInfos[index].LineOffset,
		}

		this.allComments = append(this.allComments, newComment)

		for i := index + 2; i < len(tokenInfos); i++ {

			if tokenInfos[i].Value == "*EOF*" {
				newIndex = i
				break
			}
			//加入到注释tokens中
			newComment.Tokens = append(newComment.Tokens, tokenInfos[i])
		}

		return newIndex, nil

	} else { //如果是非法字符，则忽略这行

		errorMsg := fmt.Sprintf("解析comment错误，结构不匹配,文件%s,行：%d\n",tokenInfos[index].FileName, tokenInfos[index].LocalLine+1)

		return index, errors.New(errorMsg)
	}
}

//解析field
func (this* SPParser)parseField(tokenInfos []*TokenInfo,index int) (int,error) {

	//判断区域类型
	if index+1 < len(tokenInfos) {

		if tokenInfos[index+1].Value == "/" || tokenInfos[index+1].Value == "*EOF*" { //枚举类型
			return this.parseEnumField(tokenInfos, index)

		} else {

			return this.parseClassField(tokenInfos, index)
		}
	}else {
		return index, errors.New(fmt.Sprintf("错误的行!",))
	}

	return index, nil
}

func (this* SPParser)parseEnumField(tokenInfos []*TokenInfo,index int) (int,error) {

	//fmt.Println("开始解析enum field\n")

	//判断当前行是否已经被当做类区域解析过了

	if _,has:=this.allField[tokenInfos[index].LocalLine + tokenInfos[index].LineOffset];has {
		return index,errors.New(fmt.Sprintf("解析枚举field出错,文件%s,行：%d,token:%s\n",tokenInfos[index].FileName,tokenInfos[index].LocalLine+1,tokenInfos[index].Value))
	}

	newEnumField := &EnumField{
		Index: 		tokenInfos[index].LocalLine + tokenInfos[index].LineOffset,
		Name:  		tokenInfos[index].Value,
		FileName:	tokenInfos[index].FileName,
	}

	this.allField[newEnumField.Index] = newEnumField

	return index,nil
}

func (this* SPParser)parseClassField(tokenInfos []*TokenInfo,index int) (int,error) {

	//fmt.Println("开始解析class field\n")

	//判断当前行是否已经被当做类区域解析过了
	if _,has:=this.allField[tokenInfos[index].LocalLine + tokenInfos[index].LineOffset];has {
		return index,errors.New(fmt.Sprintf("解析类field出错,文件%s,行：%d,token:%s\n",tokenInfos[index].FileName,tokenInfos[index].LocalLine+1,tokenInfos[index].Value))
	}

	newClassField := &ClassField{
		Index: 		tokenInfos[index].LocalLine +  tokenInfos[index].LineOffset,
		Name:  		tokenInfos[index].Value,
		FileName:	tokenInfos[index].FileName,
	}
	count := len(tokenInfos)
	newClassField.Name = tokenInfos[index].Value

	newIndex :=index

	if (index + 1) < count {

		if tokenInfos[index+1].Value == "[" {

			if index+3 < count && tokenInfos[index+2].Value == "]" &&
				tokenInfos[index+3].Value != "/" &&
				tokenInfos[index+3].Value != "*EOF*" {
				newIndex = index+3
				newClassField.Repeatd = true
				newClassField.Type = tokenInfos[index+3].Value

			}else {
				return index, errors.New(fmt.Sprintf("解析ClassField错误，结构不匹配，文件%s,行：%d\n", tokenInfos[index].FileName,tokenInfos[index].LocalLine+1))
			}

		}else if tokenInfos[index+1].Value != "/" && tokenInfos[index+1].Value != "*EOF*"{
			newIndex=index+1
			newClassField.Repeatd = false
			newClassField.Type = tokenInfos[index+1].Value
		}

	}else {
		return index, errors.New(fmt.Sprintf("解析ClassField错误，结构不匹配，文件%s,行：%d\n", tokenInfos[index].FileName,tokenInfos[index].LocalLine+1))
	}

	this.allField[newClassField.Index] = newClassField

	return newIndex, nil
}