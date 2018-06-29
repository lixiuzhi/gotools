package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const csBeanTemplateStr = `

/**
 * Auto generated, do not edit it
 *
 */
using UnityEngine;
using System.Collections.Generic;
using System.Text.RegularExpressions;

namespace Data.Beans
{
    public class {{.Name}}Bean
    {
    {{ range $i, $info := .Infos }}{{ if $info.IsExport }}
		{{getModifyPropertyIsInitBody $info}}
 		{{GetModifyPropertyBody $info}}
        public {{GetTypeName $info}} m_{{ $info.Name }};
        public {{GetModifyTypeName $info}} {{ $info.Name }}
		{      //{{$info.Describe2}}
            {{GetGetBody $info}}         
        }  {{ end }}  	
    {{ end }}

        public void LoadData(DataInputStream dis)
        {
            if (dis != null)
            { {{ range $i, $info := .Infos }}{{ if $info.IsExport }}
                {{GetReadBody $info }} {{ end }}{{ end }}
            }
        }
    }
}`

const csContainerTempleteStr = `
/**
 * Auto generated, do not edit it
 */

using System;
using System.Collections.Generic;
using System.IO;
using Data.Beans;
using UnityEngine;

namespace Data.Containers
{
	public class {{.Name}}Container
	{
		private List<{{.Name}}Bean> list = new List<{{.Name}}Bean>();
		private Dictionary<int, {{.Name}}Bean> map = new Dictionary<int, {{.Name}}Bean>();

		private bool isLoad = false;

		public List<{{.Name}}Bean> getList()
		{
			if(!isLoad)
			{
				loadDataFromBin();
			}
			return list;
		}

		public Dictionary<int, {{.Name}}Bean> getMap()
		{
			if(!isLoad)
			{
				loadDataFromBin();
			}
			return map;
		}

		public {{.Name}}Bean GetBean(int id)
        {
			if(!isLoad)
			{
				loadDataFromBin();
			}

            if (map.ContainsKey(id))
            {
                return map[id];
            }

            return null;
        }


		public void loadDataFromBin()
		{
			isLoad = true;
			list.Clear();
			map.Clear();
			Stream ms = ConfLoader.Singleton.getStreamByteName(typeof({{.Name}}Bean).Name + ".bytes");
			if(ms != null)
			{
				DataInputStream dis = new DataInputStream(ms);
				try
				{
					while (dis.Available() != 0)
					{
						{{.Name}}Bean bean = new {{.Name}}Bean();
						bean.LoadData(dis);
						list.Add(bean);
						if (map.ContainsKey(bean.t_id))
                        {
                            Debug.LogError("{{.Name}}Container same key:" + bean.t_id);
                        }

						map.Add(bean.t_id, bean);
					}
				}
				catch (Exception ex)
				{
					Logger.err("import data error: " + ex.ToString() + typeof({{.Name}}Bean).Name + ".bytes");
				}

				dis.Close();
				ms.Dispose();
			}
			else
			{
				Logger.err("找不到配置数据：" + typeof({{.Name}}Bean).Name + ".bytes");
			}
		}
	}

}
`

const csGamedataManagerTempleteStr = `
/**
 * Auto generated, do not edit it
 */
using Data.Beans;

namespace Data.Containers
{
	public class GameDataManager
	{
{{ range $i, $sheet := . }}
		public {{$sheet.Name}}Container {{$sheet.Name}}Container = new {{$sheet.Name}}Container();
{{ end }}

		public void loadAll()
		{
{{ range $i, $sheet := . }}
		{{$sheet.Name}}Container.loadDataFromBin();
{{ end }}
		}

	   private GameDataManager()
	   {
	   }

	   public static readonly GameDataManager Instance = new GameDataManager();

	}
}
`

type csharpHelper struct{

}

func (cs*csharpHelper)getTypeName(colinfo *DataSheetColInfo) string {

	if cs.checkIsModifyType(colinfo) {
		return "string"
	} else {
		switch colinfo.TypeName {
		case "int":
			return "int"
		case "float":
			return "float"
		case "long":
			return "long"
		case "bool":
			return "bool"
		case "text":
			return "string"
		case "textmult":
			return "string"
		case "vec2":
			return "string"
		case "vec3":
			return "string"
		default:
			return "int"
		}
	}

}

func (cs*csharpHelper)getModifyTypeName(colinfo *DataSheetColInfo) string {

	if cs.checkIsModifyType(colinfo) {
		return cs.getModifyTypeArrayName(colinfo)
	} else {
		return cs.getModifyNormalTypeName(colinfo)
	}

}

func (cs*csharpHelper)getSplitCommond(colinfo *DataSheetColInfo) *DataSheetColCommond {

	for _, comm := range colinfo.Commonds {

		if comm.Name == "split" {
			return comm
		}
	}

	return nil
}

func (cs*csharpHelper)getModifyNormalTypeName(colinfo *DataSheetColInfo) string {
	switch strings.Replace(colinfo.TypeName, "[]", "", len(colinfo.TypeName)) {
	case "int":
		return "int"
	case "float":
		return "float"
	case "long":
		return "long"
	case "bool":
		return "bool"
	case "text":
		return "string"
	case "textmult":
		return "string"
	case "vec2":
		return "Vector2"
	case "vec3":
		return "Vector3"
	default:
		return "int"
	}
}

func (cs*csharpHelper)getModifyTypeArrayName(colinfo *DataSheetColInfo) string {

	typeName := cs.getModifyNormalTypeName(colinfo)
	if !cs.checkIsSplitClass(colinfo) {
		typeName = fmt.Sprintf("List<%s>", typeName)
	}

	if strings.Contains(colinfo.TypeName, "[]") {
		typeName = fmt.Sprintf("List<%s>", typeName)
	}

	return typeName
}

func (cs*csharpHelper)getReadBody(colinfo *DataSheetColInfo) string {
	str := "m_%s=dis.%s;"
	if cs.checkIsModifyType(colinfo) {
		return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
	} else {
		switch colinfo.TypeName {
		case "int":
			return fmt.Sprintf(str, colinfo.Name, "ReadInt()")
		case "float":
			return fmt.Sprintf(str, colinfo.Name, "ReadFloat()")
		case "long":
			return fmt.Sprintf(str, colinfo.Name, "ReadLong()")
		case "bool":
			return fmt.Sprintf(str, colinfo.Name, "ReadBoolean()")
		case "text":
			return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
		case "textmult":
			return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
		case "vec2":
			return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
		case "vec3":
			return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
		default:
			if strings.Contains(colinfo.TypeName, "[]") {
				return fmt.Sprintf(str, colinfo.Name, "ReadUTF()")
			} else {
				return fmt.Sprintf(str, colinfo.Name, "ReadInt()")
			}
		}
	}

}

func (cs*csharpHelper)getGetBody(colinfo *DataSheetColInfo) string {
	var str string
	if colinfo.TypeName == "textmult" {
		str := `get
        {
            int ret;
            bool flag = int.TryParse(m_%s, out ret);
            if(flag)
            {
                string tempstr = BeanFactory.getLanguageContent(ret);
                if (tempstr == "") return m_%s;
                else return tempstr;
            }
            else
                return m_%s;
        }`
		return fmt.Sprintf(str, colinfo.Name, colinfo.Name, colinfo.Name)
	} else {
		splitComm := cs.getSplitCommond(colinfo)
		modifyTypeName := cs.getModifyNormalTypeName(colinfo)
		upperTypename := strings.ToUpper(SubString(modifyTypeName, 0, 1)) + SubString(modifyTypeName, 1, strings.Count(modifyTypeName, ""))
		//fmt.Printf(modifyTypeName + ":" + upperTypename + "\n")
		isSplitClass := cs.checkIsSplitClass(colinfo)

		if splitComm != nil || isSplitClass {
			if strings.Contains(colinfo.TypeName, "[]") {
				upperTypename += "List2"
			} else {
				if !isSplitClass {
					upperTypename += "List"
				}
			}

			str = `get
			{
				if (!mm_IsInit_%s) { mm_%s = m_%s.splitTo%s(%s); mm_IsInit_%s = true;}  
				return mm_%s;
			}`

			params := splitComm.ParamStr
			if !IsStrNull(params) {
				if len(params) > 1 {
					params = "\"" + params + "\""
				} else {
					params = "'" + params + "'"
				}
			}

			return fmt.Sprintf(str, colinfo.Name, colinfo.Name, colinfo.Name, upperTypename, params, colinfo.Name, colinfo.Name)

		} else {
			switch colinfo.TypeName {
			default:
				str := ` get{return m_%s;}`
				return fmt.Sprintf(str, colinfo.Name)
			}
		}
	}
}

func (cs*csharpHelper)getSetBody(colinfo *DataSheetColInfo) string {
	if cs.checkIsModifyType(colinfo) {
		return fmt.Sprintf(`set{ mm_%s = value; }`, colinfo.Name)
	} else {
		return ""
	}
}

//检查 返回类型是否是修改后的类型
func (cs*csharpHelper)checkIsModifyType(colinfo *DataSheetColInfo) bool {
	return cs.checkIsSplitClass(colinfo) || cs.getSplitCommond(colinfo) != nil
}

//检查是否为自定义类
func (cs*csharpHelper)checkIsSplitClass(colinfo *DataSheetColInfo) bool {
	modifyType := cs.getModifyNormalTypeName(colinfo)
	switch modifyType {
	case "Vector2", "Vector3":
		return true
	default:
		return false
	}

	return false
}

func (cs*csharpHelper)getModifyPropertyBody(colinfo *DataSheetColInfo) string {
	modifyTypeName := cs.getModifyTypeName(colinfo)
	if cs.checkIsModifyType(colinfo) {
		return fmt.Sprintf(`public %s mm_%s;//原始属性修改过后的属性`, modifyTypeName, colinfo.Name)
	} else {
		return ""
	}
}

func (cs*csharpHelper)getModifyPropertyIsInitBody(colinfo *DataSheetColInfo) string {
	if cs.checkIsModifyType(colinfo) {
		return fmt.Sprintf(`public bool mm_IsInit_%s;//属性是否被初始化赋值`, colinfo.Name)
	} else {
		return ""
	}
}


func (cs*csharpHelper)IsContainItem(dataSheets []*DataSheet, path string) bool {

	for _, value1 := range dataSheets {
		if strings.Contains(path, value1.Name+"Bean.cs") {
			return true
		} else if strings.Contains(path, value1.Name+"Container.cs") {
			return true
		} else {

		}
	}

	return false
}



//生成CS代码
func GenCsharp(dataSheets []*DataSheet, outpath string) error {

	cs :=&csharpHelper{}

	if b, _ := PathExists(outpath); !b {
		return errors.New("生成cs，目录不存，或者出错!\n")
	}
	//初始化目录
	CreateDir(outpath + "/bean")
	CreateDir(outpath + "/container")


	//删除不存在的配置表文件
	files, err := WalkDir(outpath, ".cs")
	if len(files) > 0 {
		for _, value := range files {
			if !cs.IsContainItem(dataSheets, value) {
				os.Remove(value)
				fmt.Print("remove file:" + value + "\n")
			} else {

			}
		}
	}

	//os.RemoveAll(outpath + "/bean")
	//os.RemoveAll(outpath + "/container")

	//os.MkdirAll(outpath+"/bean", 0777)
	//os.MkdirAll(outpath+"/container", 0777)

	funcMap := template.FuncMap{
		"GetTypeName":                 cs.getTypeName,
		"getModifyPropertyIsInitBody": cs.getModifyPropertyIsInitBody,
		"GetModifyPropertyBody":       cs.getModifyPropertyBody,
		"GetGetBody":                  cs.getGetBody,
		"GetSetBody":                  cs.getSetBody,
		"GetReadBody":                 cs.getReadBody,
		"GetModifyTypeName":           cs.getModifyTypeName,
	}

	beantpl, err := template.New("gencsbean").Funcs(funcMap).Parse(csBeanTemplateStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Containertpl, err := template.New("gencscontainer").Funcs(funcMap).Parse(csContainerTempleteStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	gamedataMgrtpl, err := template.New("gencsgamedatamgr").Funcs(funcMap).Parse(csGamedataManagerTempleteStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var dataMgrBuf bytes.Buffer

	err = gamedataMgrtpl.Execute(&dataMgrBuf, dataSheets)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ioutil.WriteFile(outpath+"/GameDataManager.cs", dataMgrBuf.Bytes(), 0666)

	//
	for _, dataSheet := range dataSheets {

		var beanBuf bytes.Buffer
		//生成bean
		err = beantpl.Execute(&beanBuf, dataSheet)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		//fmt.Println(string(beanBuf.Bytes()))
		ioutil.WriteFile(outpath+"/bean/"+dataSheet.Name+"Bean.cs", beanBuf.Bytes(), 0666)

		var containerBuf bytes.Buffer
		//生成container
		err = Containertpl.Execute(&containerBuf, dataSheet)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		//fmt.Println(string(containerBuf.Bytes()))

		ioutil.WriteFile(outpath+"/container/"+dataSheet.Name+"Container.cs", containerBuf.Bytes(), 0666)
	}
	fmt.Print("\n************************************生成CS文件 完毕\n\n")

	return nil
}
