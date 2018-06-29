package core

const javaBeanTemplateStr = `
/**
 * Auto generated, do not edit it
 *
 * {{.Name}}
 */
package data.bean;

import data.GameDataManager;
public class {{.Name}}Bean
{
 {{ range $i, $info := .Infos }}{{ if $info.IsExport }}
    private {{GetTypeName $info}} {{ $info.Name }}; //{{$info.Describe2}} {{ end }}{{ end }}

{{ range $i, $info := .Infos }}{{ if $info.IsExport }}
    /**
    *get //{{$info.Describe2}}
    */
    public {{GetTypeName $info}} get{{GetInfoUpperName $info}}()
    {
        return {{$info.Name}};
    }

    /**
    *set //{{$info.Describe2}}
    */
    public void set{{GetInfoUpperName $info}}({{GetTypeNameWithDB $info}} {{ $info.Name }})
    { {{ if NeedConvertToLong $info }}
        this.{{$info.Name}}=Long.parseLong({{ $info.Name }});
    {{else}}    this.{{$info.Name}}={{ $info.Name }};{{ end }}
    }  {{ end }} {{ end }}
}

`