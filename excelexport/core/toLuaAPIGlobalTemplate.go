package core

var toLuaAPIGlobalTemplate= `
---@class CfgTable {{ range $i, $sheet := . }}
---@field public {{$sheet.Name}} {{$sheet.Name}}Bean[]  @{{ $sheet.TableName }}{{ end }}
local m = {}
`