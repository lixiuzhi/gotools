package core

var toCSTemplate = `
using System;
using Sproto;
using System.Collections.Generic;

namespace proto
{ {{range $i, $enum := .Enums}}
	public enum {{$enum.Name}} { {{range $i,$enumfield :=$enum.Fields}} 
		{{$enumfield.Name}} = {{$enumfield.LocalIndex}},	{{GetEnumFieldComment	$enumfield}} {{end}}
	} {{end}}
	
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