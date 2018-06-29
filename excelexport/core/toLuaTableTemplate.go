package core

var toLuaTableTemplate= `
local {{.Name}} = { {{ $Datas := .Data }}{{ $Infos := .Infos }}{{ $sheet := . }}{{ range $i, $data:= $Datas }}
    [{{index $data 0}}] = { {{ range $j, $info := $Infos }}{{ if $info.IsExport }}{{$info.Name}}={{ GetFieldValue $sheet $i $j}}{{ if IsNotRowEnd $sheet $j }},{{ end }} {{ end }}{{ end }} }{{ if IsNotLastRow $sheet $i}},{{ end }} {{ end }}
}

{{.Name}}.__index =function(t,k)
    --log.err("配置表项不存在,{{.Name}} :"+tostring(k))
end

setmetatable({{.Name}},{{.Name}})

return {{.Name}}
`