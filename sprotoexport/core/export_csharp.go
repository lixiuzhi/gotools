package core

import (
	"bytes"
	"fmt"
	"text/template"
	"io/ioutil"
)

const templateCSStr = `
using System;
using Sproto;
using System.Collections.Generic;

namespace proto
{
{{range $i, $enum := .Enums}}
	public enum {{$enum.Name}} { {{range $i,$enumfield :=$enum.Fields}} 
		{{$enumfield.Name}} = {{$enumfield.LocalIndex}},	{{GetEnumFieldComment	$enumfield}} {{end}}
	}
{{end}}
	
{{range $i, $class := .Classes}}
	public class {{$class.Name}} : SprotoTypeBase {
		private static int max_field_count = {{len $class.Fields}};
	{{range $fieldIndex, $field := $class.Fields}}
		[SprotoHasField]
		public bool Has{{$field.Name}}{
			get { return base.has_field.has_field({{$fieldIndex}}); }
		}
	
		private {{GetClassFieldType $field}} _{{$field.Name}}; // tag {{$fieldIndex}} 
		public {{GetClassFieldType $field}} {{$field.Name}} { {{GetClassFieldComment	$field}}
			get{ return _{{$field.Name}}; }
			set{ base.has_field.set_field({{$fieldIndex}},true); _{{$field.Name}} = value; }
		}
	{{end}}
	
		public {{$class.Name}}() : base(max_field_count) {}
	
		public {{$class.Name}}(byte[] buffer) : base(max_field_count, buffer) {
			this.decode ();
		} 
	
		protected override void decode () {
			int tag = -1;
			while (-1 != (tag = base.deserialize.read_tag ())) {
				switch (tag) {	
	{{range $fieldIndex, $field := $class.Fields}}
					case {{$fieldIndex}}:
						this.{{$field.Name}} = base.deserialize.{{getClassFieldReadFunc $field}};
						break;
	{{end}}
					default:
						base.deserialize.read_unknow_data ();
						break;
					}
				}
			}
	
	public override int encode (SprotoStream stream) {
				base.serialize.open (stream);
	{{range $fieldIndex, $field := $class.Fields}}
				if (base.has_field.has_field ({{$fieldIndex}})) {
					base.serialize.{{getClassFieldWriteFuncName $field}}(this.{{$field.Name}}, {{$fieldIndex}});
				} 
	{{end}}
				return base.serialize.close ();
			}
	} 
{{end}}
}
`

func getCSClassFieldReadFunc(field * ClassField) string {

	str := "read"

	f := false

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
	case "double":
		str += "_double"
	default:
		f = true
		if field.TypeIsEnum {
			str += "_enum"
		}else{
			str += "_obj"
		}
	}

	if field.Repeatd {
		str += "_list"
	}

	if f {
		str += "<" + field.Type + ">"
	}

	str += "()"
	return str
}

func getCSClassFieldWriteFuncName(field * ClassField) string {

	str := "write"

	if field.TypeIsEnum{
		str+="_enum"
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
			if field.Repeatd{
				str +="<" + field.Type + ">"
			}
		}
	}

	return str
}

func getCSClassFieldType(field * ClassField) string {

	str := ""

	var baseType = ""

	switch field.Type {
	case "binary":
		baseType = "byte[]"
	case "int32":
		baseType = "Int32"
	case "int64":
		baseType = "Int64"
	case "bool":
		baseType = "bool"
	case "string":
		baseType = "string"
	case "double":
		baseType = "double"
	default:
		baseType = field.Type
	}

	if field.Repeatd{
		str = "List<"+ baseType +">"
	}else{
		str = baseType
	}

	return str
}

func getCSClassFieldComment(field * ClassField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}

func getCSEnumFieldComment(field * EnumField) string {
	var str =""
	for _,v:= range field.Comment  {
		str+="//"
		for _,v1:=range v.Tokens{
			str+=v1.Value+" "
		}
	}
	return str
}


func GenCS(parser * SPParser,outPath string){

	outfileName := parser.fileName +".cs"
	fmt.Println("开始生成csharp文件:",outfileName)

	var bf bytes.Buffer

	funcMap := template.FuncMap{
		"GetEnumFieldComment"			:getCSEnumFieldComment,
		"GetClassFieldComment"			:getCSClassFieldComment,
		"GetClassFieldType"				:getCSClassFieldType,
		"getClassFieldReadFunc"			:getCSClassFieldReadFunc,
		"getClassFieldWriteFuncName"	:getCSClassFieldWriteFuncName,
	}

	tpl, err := template.New("genGolang").Funcs(funcMap).Parse(templateCSStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = tpl.Execute(&bf, *parser)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ioutil.WriteFile(outPath+"/"+outfileName,bf.Bytes(),0666)
}
