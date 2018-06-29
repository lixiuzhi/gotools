package core

const toJavaMessageUtilTemplate = `
package {{GetJavaPackageName}};
import java.util.*;
import {{GetJavaPackageName}}.sproto.*;

public class MessageUtil {

    @FunctionalInterface
    public interface GenMsgAction { SprotoTypeBase make(byte[]data); }

    static Map idByName = new HashMap();
    static Map NameById = new HashMap();    
    static Map GenMsgById = new HashMap(); 

    public static void Init()
    {
        idByName.clear();
        NameById.clear();
        GenMsgById.clear();
 {{range $i, $class := .Classes}}
        idByName.put({{$class.MsgID}},"{{$class.Name}}");
        NameById.put("{{$class.Name}}",{{$class.MsgID}});
        GenMsgById.put({{$class.MsgID}},(GenMsgAction)(byte[]data)->{return new {{$class.Name}}(data);});
{{end}}
    }
    
    public static SprotoTypeBase GetMsgById(int id,byte[]data)
    {
        GenMsgAction action = (GenMsgAction)GenMsgById.get(id); 
        if(action!=null)
        {
            return action.make(data);
        }
        return null;
    }
}
`