package core

const javaGameDataManagerTemplate = `
/**
 * Auto generated, do not edit it
 */
package data;

{{ range $i, $info := . }}
import data.container.{{$info.Name}}Container;{{ end }}
import org.apache.log4j.Logger;
import org.apache.ibatis.session.SqlSessionFactory;

public class GameDataManager
{
    private final Logger log = Logger.getLogger(GameDataManager.class);

{{ range $i, $info := . }}
    public volatile {{$info.Name}}Container {{$info.Name}}Container = new {{$info.Name}}Container();{{ end }}

    private SqlSessionFactory sessionFactory;

    public void loadAll()
    {
	    log.info("Start load all game data ...");
{{ range $i, $info := . }}
        {{$info.Name}}Container.load();{{ end }}
    }

    public GameDataManager setSqlSessionFactory(SqlSessionFactory sessionFactory)
    {
        this.sessionFactory = sessionFactory;
        return this;
    }

    public SqlSessionFactory getSqlSessionFactory()
    {
        return sessionFactory;
    }

    public static GameDataManager getInstance()
    {
        return Singleton.INSTANCE.getProcessor();
    }

    private enum Singleton
    {
        INSTANCE;
        GameDataManager manager;

        Singleton()
        {
            this.manager = new GameDataManager();
        }

        GameDataManager getProcessor()
        {
            return manager;
        }
    }
}

`