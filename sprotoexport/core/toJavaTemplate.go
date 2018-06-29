package core

const toJavaEnumTemplate = `
package {{GetJavaPackageName}};

public class {{.Name}}{
 {{range $i,$enumfield :=.Fields}} 
	public static final int {{$enumfield.Name}} = {{$enumfield.LocalIndex}};	{{GetEnumFieldComment	$enumfield}} {{end}}
}
`

const toJavaClassTemplate = `
package {{GetJavaPackageName}};

import com.lxz.sproto.*;
import java.util.List;
import java.util.function.Supplier;
import java.util.ArrayList;

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

  	public int GetId ()
    {
        return {{.MsgID}};
    }

	{{range $fieldIndex, $field := .Fields}} 
	private {{GetClassFieldType $field}} _{{$field.Name}}; // tag {{$fieldIndex}}
	public boolean Has{{$field.Name}}(){
		return super.has_field.has_field({{$fieldIndex}});
	}
	public {{GetClassFieldType $field}} get{{$field.Name}}() { {{if $field.Repeatd}}
		if (_{{$field.Name}}==null){
			_{{$field.Name}} = new ArrayList<>();
		}
		{{end}}return _{{$field.Name}};
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