package core

import (
    "bytes"
    "fmt"
    "text/template"
    "io/ioutil"
    "os"
    "errors"
)

const beanTemplateStr = `

/**
 * Auto generated, do not edit it
 *
 */

using System.Text.RegularExpressions;

namespace Data.Beans
{
    public class {{.Name}}Bean
    {
    {{ range $i, $info := .Infos }}{{ if $info.IsExport }}
        public {{GetTypeName $info}} m_{{ $info.Name }};
        public {{GetTypeName $info}} {{ $info.Name }}{      //{{$info.Describe2}}
            {{GetGetBody $info}}
             set{ m_{{$info.Name}} = value; }
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

const containerTempleteStr = `
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

		public List<{{.Name}}Bean> getList()
		{
			return list;
		}

		public Dictionary<int, {{.Name}}Bean> getMap()
		{
			return map;
		}

		public {{.Name}}Bean GetBean(int id)
        {
            if (map.ContainsKey(id))
            {
                return map[id];
            }

            return null;
        }


		public void loadDataFromBin()
		{
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
					Logger.err("import data error: " + ex.ToString());
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

const gamedataManagerTempleteStr = `
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


func getTypeName(colinfo *DataSheetColInfo) string {
    switch colinfo.TypeName {
    case "int":
        return "int"
    case "text":
        return "string"
    case "textmult":
        return "string"
    default:
        return "int"
    }
}

func getReadBody(colinfo *DataSheetColInfo) string {

    str:="m_%s=dis.%s;"
    switch colinfo.TypeName {
    case "text":
        return fmt.Sprintf(str,colinfo.Name,"ReadUTF()")
    case "textmult":
        return fmt.Sprintf(str,colinfo.Name,"ReadUTF()")
    default:
        return fmt.Sprintf(str,colinfo.Name,"ReadInt()")
    }
}

func getGetBody(colinfo *DataSheetColInfo) string {
    switch colinfo.TypeName {
    case "textmult":
        str:=`get
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
        return fmt.Sprintf(str,colinfo.Name,colinfo.Name,colinfo.Name)
    default:
        str:=` get{return m_%s;}`
        return fmt.Sprintf(str,colinfo.Name)
    }
}

//生成CS代码
func GenCsharp( dataSheets []*DataSheet,outpath string) error {

    if b, _ := PathExists(outpath); !b {
        return errors.New("生成cs，目录不存，或者出错!\n")
    }
    //初始化目录

    os.RemoveAll(outpath + "/bean")
    os.RemoveAll(outpath + "/container")

    os.MkdirAll(outpath+"/bean", 0777)
    os.MkdirAll(outpath+"/container", 0777)

    funcMap := template.FuncMap{
        "GetTypeName": getTypeName,
        "GetGetBody":  getGetBody,
        "GetReadBody": getReadBody,
    }

    beantpl, err := template.New("gencsbean").Funcs(funcMap).Parse(beanTemplateStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    Containertpl, err := template.New("gencscontainer").Funcs(funcMap).Parse(containerTempleteStr)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    gamedataMgrtpl, err := template.New("gencsgamedatamgr").Funcs(funcMap).Parse(gamedataManagerTempleteStr)
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
        err = beantpl.Execute(&beanBuf, *dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        //fmt.Println(string(beanBuf.Bytes()))
        ioutil.WriteFile(outpath+"/bean/"+dataSheet.Name+"Bean.cs", beanBuf.Bytes(), 0666)

        var containerBuf bytes.Buffer
        //生成container
        err = Containertpl.Execute(&containerBuf, *dataSheet)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        //fmt.Println(string(containerBuf.Bytes()))

        ioutil.WriteFile(outpath+"/container/"+dataSheet.Name+"Container.cs", containerBuf.Bytes(), 0666)
    }

    return nil
}