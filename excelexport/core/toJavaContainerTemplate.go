package core

const javaContainerTemplateStr = `
/**
 * Auto generated, do not edit it
 *
 * {{.Name}}
 */
package data.container;

import java.util.HashMap;
import java.util.Map;
import java.util.Iterator;
import java.util.List;
import data.bean.{{.Name}}Bean;
import data.dao.{{.Name}}Dao;

public class {{.Name}}Container
{
    private List<{{.Name}}Bean> list;
    private final Map<Integer, {{.Name}}Bean> map = new HashMap<>();
    private final {{.Name}}Dao dao = new {{.Name}}Dao();

    public void load()
    {
        list = dao.select();
        Iterator<{{.Name}}Bean> iter = list.iterator();
        while (iter.hasNext())
        {
            {{.Name}}Bean bean = ({{.Name}}Bean) iter.next();
            map.put(bean.get{{GetIdUpperName .}}(), bean);
        }
    }

    public List<{{.Name}}Bean> getList()
    {
        return list;
    }

    public Map<Integer, {{.Name}}Bean> getMap()
    {
        return map;
    }
}

`