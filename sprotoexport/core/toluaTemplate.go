package core

var toluaTemplate = `
Enum = { {{range $i, $enum := .Enums}}
	{{$enum.Name}} = { {{range $i,$enumfield :=$enum.Fields}} 
		{{$enumfield.Name}} = {{$enumfield.LocalIndex}}, {{GetEnumFieldComment	$enumfield}} {{end}}
	},
{{end}}
}

local sproto = {
	Schema = [[
{{range $i, $class := .Classes}}
	.{{$class.Name}} { {{range $fieldIndex, $field := $class.Fields}}
	    {{$field.Name}} {{$fieldIndex}} : {{GetClassFieldType $field}}{{end}}
    } 
{{end}}
]],

NameByID = { {{range $i, $class := .Classes}}
	[{{$class.MsgID}}] = "{{$class.Name}}",{{end}}
},

IDByName = { {{range $i, $class := .Classes}}
	["{{$class.Name}}"]={{$class.MsgID}},{{end}}
},

ResetByID = {
  {{range $i, $class := .Classes}}
	[{{$class.MsgID}}] = function(obj)   --{{$class.Name}}
        if obj == nil then return end {{range $fieldIndex, $field := $class.Fields}}
	    obj.{{$field.Name}} = {{GetFieldDefaultValue $field}} {{end}}
    end,
{{end}}
}
}

return sproto
`