package core

const javaDaoTemplateStr = `
/**
 * Auto generated, do not edit it
 *
 * {{.Name}}
 */
package data.dao;

import data.GameDataManager;
import java.util.List;
import org.apache.ibatis.session.SqlSession;
import data.bean.{{.Name}}Bean;

public class {{.Name}}Dao
{
    public List<{{.Name}}Bean> select()
    {
        try
        (SqlSession session = GameDataManager.getInstance().getSqlSessionFactory().openSession())
        {
            List<{{.Name}}Bean> list = session.selectList("{{.Name}}.selectAll");
            return list;
        }
    }
}
`