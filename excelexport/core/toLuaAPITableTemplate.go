package core

var toLuaAPITableTemplate= `
---@class {{ .Name }}Bean {{ range $i, $info := .Infos }}{{ if $info.IsExport }}
---@field {{ $info.Name }} {{GetTypeName $info}} @{{$info.Describe2}} {{ end }}{{ end }}
local m={}
`