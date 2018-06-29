package core

const javaDBGameDataConfigTemplate = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE configuration PUBLIC "-//ibatis.apache.org//DTD Config 3.0//EN" "http://ibatis.apache.org/dtd/ibatis-3-config.dtd">
<configuration>
    <environments default="development">
	<environment id="development">
	    <transactionManager type="JDBC"/>
	    <dataSource type="POOLED">
		<property name="driver" value="com.mysql.jdbc.Driver"/>
		<property name="url" value="${url}"/>
		<property name="username" value="${username}"/>
		<property name="password" value="${password}"/>
		<property name="poolPingQuery" value="select 1"/>
		<property name="poolPingEnabled" value="true"/>
		<property name="poolPingConnectionsNotUsedFor" value="3600000"/>
	    </dataSource>
	</environment>
    </environments>
    <mappers>{{ range $i, $info := . }}
            <mapper resource="data/sqlmap/{{$info.Name}}.xml" />{{ end }}
    </mappers>
</configuration>

`