package core

import (
	"fmt"
	"errors"
	"os"
	"strings"
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
	Repeatd		bool
	LocalIndex	int
}

type EnumField struct {
	Index 		int
	Comment 	[]*CommentField
	Name 		string
	Value		int
	LocalIndex	int
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
}

type SPParser struct {

	fileName	string
	maxLine		int
	allComments []*CommentField
	allField 	map[int]interface{}

	Enums 		map[int]*EnumType
	Classes 	map[int]*ClassType
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
			if v, ok := this.Enums[i]; ok {
				v.Comment = append(v.Comment, c)
				break
			}

			if v, ok := this.Classes[i]; ok {
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
				}else {
					return errors.New(fmt.Sprintf("组合枚举时错误，无效的filed,行：%d\n",i+1))
				}
			}
		}
	}

	//填充class区域
	for _, class := range this.Classes {

		for i := class.Index; i <= class.EndIndex; i++ {

			if f, ok := this.allField[i]; ok {

				if rf,ok:=f.(*ClassField);ok{
					rf.LocalIndex = len(class.Fields)+1
					class.Fields = append(class.Fields, rf)

					if _,ok:=enumNameToEnum[rf.Type];ok{
						rf.TypeIsEnum = true
					}else {
						rf.TypeIsEnum = false
					}

				}else {
					return errors.New(fmt.Sprintf("组合类时错误，无效的filed,行：%d\n",i+1))
				}
			}
		}
	}



	return nil
}

func (this*SPParser) Parse(tokenInfos []*TokenInfo,fileName string) {

	fileName = strings.Replace(fileName,".sp","",1)

	fmt.Println("开始解析文件:",fileName)

	this.allComments = make([]*CommentField, 0, 100)
	this.allField = make(map[int]interface{})
	this.Enums = make(map[int]*EnumType)
	this.Classes = make(map[int]*ClassType)
	this.fileName =fileName

	count := len(tokenInfos)

	if count > 0 {
		this.maxLine = tokenInfos[count-1].Line
	} else {
		this.maxLine = 0
	}

	for i := 0; i < count; i++ {
		//fmt.Printf("开始解析line:%d,%s\n", tokenInfos[i].Line, tokenInfos[i].Value)

		switch tokenInfos[i].Value {
		case "/":
			newIndex, err := this.parseComment(tokenInfos, i)

			if err != nil {
				fmt.Println("解析出错，", err.Error())
				os.Exit(2)
			} else {
				i = newIndex
			}

		case "message":
			newIndex, err := this.parseClass(tokenInfos, i)

			if err != nil {
				fmt.Println("解析class出错，", err.Error())
				os.Exit(2)
			} else {
				i = newIndex
			}

		case "enum":
			newIndex, err := this.parseEnum(tokenInfos, i)

			if err != nil {
				fmt.Println("解析枚举出错，", err.Error())
				os.Exit(2)
			} else {
				i = newIndex
			}

		case "*EOF*":

		case "}":

		default:
			newIndex, err := this.parseField(tokenInfos, i)

			if err != nil {
				fmt.Println("解析field出错，", err.Error())
				os.Exit(2)
			} else {
				i = newIndex
			}
		}
	}

	if err:=this.fill();err!=nil{
		fmt.Println("fill时错误：", err.Error())
		os.Exit(2)
	}
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
		return index, errors.New(fmt.Sprintf("解析类错误，没有找到左括号，行：%d\n", tokenInfos[index].Line+1))
	}
	if len(tokenInfos) > leftBracketsIndex+1 {
		newIndex := leftBracketsIndex + 1

		//查找右括号index
		for i := index + 2; i < len(tokenInfos); i++ {

			if tokenInfos[i].Value == "}" && !this.isCommentField(tokenInfos, tokenInfos[i], i) {
				newClass := &ClassType{
					Index:    tokenInfos[index].Line,
					EndIndex: tokenInfos[i].Line,
					Comment:  make([]*CommentField, 0),
					Fields:   make([]*ClassField, 0),
					Name:     tokenInfos[index+1].Value,
				}

				this.Classes[newClass.Index] = newClass

				return newIndex, nil
			}
		}
	}
	return index, errors.New(fmt.Sprintf("解析类错误，括号可能不匹配，行：%d\n", tokenInfos[index].Line+1))

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
		return index, errors.New(fmt.Sprintf("解析枚举错误，没有找到左括号，行：%d\n", tokenInfos[index].Line+1))
	}
	if len(tokenInfos) > leftBracketsIndex+1 {
		newIndex := leftBracketsIndex + 1

		//查找右括号index
		for i := index + 2; i < len(tokenInfos); i++ {

			if tokenInfos[i].Value == "}" && !this.isCommentField(tokenInfos, tokenInfos[i], i) {
				newEnum := &EnumType{
					Index	:tokenInfos[index].Line,
					EndIndex:tokenInfos[i].Line,
					Comment	:make([]*CommentField,0),
					Fields	:make([]*EnumField,0),
					Name	:tokenInfos[index+1].Value,
				}
				this.Enums[newEnum.Index]=newEnum

				return newIndex, nil
			}
		}
	}
	return index, errors.New(fmt.Sprintf("解析枚举错误，括号可能不匹配，行：%d\n", tokenInfos[index].Line+1))

}

//分离注释区域
func(this*SPParser) parseComment( tokenInfos []*TokenInfo,index int) (int,error) {

	newIndex := index
	if tokenInfos[index].Value == "/" && len(tokenInfos) > index+1 && tokenInfos[index+1].Value == "/" {
		//产生新的comment
		var newComment = &CommentField{
			Tokens: make([]*TokenInfo, 0),
			Index:  tokenInfos[index].Line,
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

		errorMsg := fmt.Sprintf("解析comment错误，结构不匹配，行：%d\n", tokenInfos[index].Line+1)

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
	}

	return index, nil
}

func (this* SPParser)parseEnumField(tokenInfos []*TokenInfo,index int) (int,error) {

	//fmt.Println("开始解析enum field\n")

	//判断当前行是否已经被当做类区域解析过了

	if _,has:=this.allField[tokenInfos[index].Line];has {
		return index,errors.New(fmt.Sprintf("解析枚举field出错,行：%d,token:%s\n",tokenInfos[index].Line+1,tokenInfos[index].Value))
	}

	newEnumField := &EnumField{
		Index: tokenInfos[index].Line,
		Name:  tokenInfos[index].Value,
	}

	this.allField[newEnumField.Index] = newEnumField

	return index,nil
}

func (this* SPParser)parseClassField(tokenInfos []*TokenInfo,index int) (int,error) {

	//fmt.Println("开始解析class field\n")

	//判断当前行是否已经被当做类区域解析过了
	if _,has:=this.allField[tokenInfos[index].Line];has {
		return index,errors.New(fmt.Sprintf("解析类field出错,行：%d,token:%s\n",tokenInfos[index].Line+1,tokenInfos[index].Value))
	}

	newClassField := &ClassField{
		Index: tokenInfos[index].Line,
		Name:  tokenInfos[index].Value,
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
				return index, errors.New(fmt.Sprintf("解析ClassField错误，结构不匹配，行：%d\n", tokenInfos[index].Line+1))
			}

		}else if tokenInfos[index+1].Value != "/" && tokenInfos[index+1].Value != "*EOF*"{
			newIndex=index+1
			newClassField.Repeatd = false
			newClassField.Type = tokenInfos[index+1].Value
		}

	}else {
		return index, errors.New(fmt.Sprintf("解析ClassField错误，结构不匹配，行：%d\n", tokenInfos[index].Line+1))
	}

	this.allField[newClassField.Index] = newClassField

	return newIndex, nil
}