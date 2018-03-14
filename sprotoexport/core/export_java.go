package core

import (
	"bytes"
	"fmt"
	"text/template"
	"io/ioutil"
)

const templateJavaEnum = `
package {{GetJavaPackageName}};

public class {{.Name}}{
 {{range $i,$enumfield :=.Fields}} 
	public static final int {{$enumfield.Name}} = {{$enumfield.LocalIndex}};	{{GetEnumFieldComment	$enumfield}} {{end}}
}
`

const templateJavaClass = `
package {{GetJavaPackageName}};

import com.lxz.sproto.*;
import java.util.List;
import java.util.function.Supplier;

public class {{.Name}} extends SprotoTypeBase {

	private static int max_field_count = {{len .Fields}};
	public static Supplier<{{.Name}}> proto_supplier = ()->new {{.Name}}();

	public {{.Name}}(){
			super(max_field_count);
	}
	
	public {{.Name}}(byte[] buffer){
			super(max_field_count, buffer);
			this.decode ();
	} 

	{{range $fieldIndex, $field := .Fields}} 
	private {{GetClassFieldType $field}} _{{$field.Name}}; // tag {{$fieldIndex}}
	public boolean Has{{$field.Name}}(){
		return super.has_field.has_field({{$fieldIndex}});
	}
	public {{GetClassFieldType $field}} get{{$field.Name}}() {
		return _{{$field.Name}};
	}
	public void set{{$field.Name}}({{GetClassFieldType $field}} value){
		super.has_field.set_field({{$fieldIndex}},true);
		_{{$field.Name}} = value;
	}
 
	{{end}}
	protected void decode () {
		int tag = -1;
		while (-1 != (tag = super.deserialize.read_tag ())) {
			switch (tag) {	
	{{range $fieldIndex, $field := .Fields}}
			case {{$fieldIndex}}:
				this.set{{$field.Name}}(super.deserialize.{{getClassFieldReadFunc $field}});
				break;
	{{end}}
			default:
				super.deserialize.read_unknow_data ();
				break;
			}
		}
	}
	
	public int encode (SprotoStream stream) {
			super.serialize.open (stream);
	{{range $fieldIndex, $field := .Fields}}
			if (super.has_field.has_field ({{$fieldIndex}})) {
				super.serialize.{{getClassFieldWriteFuncName $field}}(this._{{$field.Name}}, {{$fieldIndex}});
			} 
	{{end}}
			return super.serialize.close ();
	}
}
`

func getJavaClassFieldReadFunc(field * ClassField) string {

	str := "read"

	isclass := false

	switch field.Type {
	case "int32":
		str += "_int32"
	case "int64":
		str += "_int64"
	case "bool":
		str += "_boolean"
	case "string":
		str += "_string"
	case "binary":
		str += "_bytes"
	default:

		if field.TypeIsEnum {
			str += "_int32"
		}else{
			isclass = true
			str += "_obj"
		}
	}

	if field.Repeatd {
		str += "_list"
	}

	if isclass {
		str += "("+ field.Type +".proto_supplier)"
	}else {
		str += "()"
	}
	return str
}

func getJavaClassFieldWriteFuncName(field * ClassField) string {

	str := "write"

	if field.TypeIsEnum{
		str+="_int32"
	}else {
		switch field.Type {
		case "int32":
			str += "_int32"
		case "int64":
			str += "_int64"
		case "bool":
			str += "_boolean"
		case "string":
			str += "_string"
		case "binary":
			str += "_bytes"
		default:
			str	+= "_obj"
		}
	}

	return str
}

func getJavaClassFieldType(field * ClassField) string {

	str := ""

	var baseType = ""
	var ListBaseType = ""

	switch field.Type {
	case "binary":
		baseType = "byte[]"
		ListBaseType = "byte[]"
	case "int32":
		baseType = "int"
		ListBaseType = "Integer"
	case "int64":
		baseType = "long"
		ListBaseType = "Long"
	case "bool":
		baseType = "boolean"
		ListBaseType = "Boolean"
	case "string":
		baseType = "String"
		ListBaseType = "String"
	default:
		baseType = field.Type
		ListBaseType = field.Type
	}

	if field.TypeIsEnum{
		baseType = "int"
		ListBaseType = "Integer"
	}

	if field.Repeatd{
		str = "List<"+ ListBaseType +">"
	}else{
		str = baseType
	}

	return str
}

func getJavaClassFieldComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}

func getJavaEnumFieldComment(field * EnumField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenJava(parser * SPParser,packageName string,outPath string){

	funcMap := template.FuncMap{
		"GetJavaPackageName"			:func() string{return packageName},

		"GetEnumFieldComment"			:getJavaEnumFieldComment,
		"GetClassFieldComment"			:getJavaClassFieldComment,
		"GetClassFieldType"				:getJavaClassFieldType,
		"getClassFieldReadFunc"			:getJavaClassFieldReadFunc,
		"getClassFieldWriteFuncName"	:getJavaClassFieldWriteFuncName,
	}

	//导出枚举类
	enumTpl, err := template.New("genGolang").Funcs(funcMap).Parse(templateJavaEnum)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _,enum :=range parser.Enums{

		outfileName := enum.Name +".java"
		fmt.Println("开始生成java文件:",outfileName)

		var bf bytes.Buffer
		err = enumTpl.Execute(&bf, enum)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)
	}

	//导出class
	classTpl, err := template.New("genGolang").Funcs(funcMap).Parse(templateJavaClass)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _,class :=range parser.Classes{

		outfileName := class.Name +".java"
		fmt.Println("开始生成java文件:",outfileName)

		var bf bytes.Buffer
		err = classTpl.Execute(&bf, class)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)
	}
}
