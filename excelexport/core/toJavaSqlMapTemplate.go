package core

const javaSqlMapTemplateStr = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper  PUBLIC "-//ibatis.apache.org//DTD Mapper 3.0//EN"  "http://ibatis.apache.org/dtd/ibatis-3-mapper.dtd">
<mapper namespace="{{.Name}}">
    <resultMap id="bean" type="data.bean.{{.Name}}Bean" >{{ range $i, $info := .Infos }}{{ if $info.IsExport }}
            <result column="{{$info.Name}}" property="{{$info.Name}}" jdbcType="{{GetDBType $info}}" />{{ end }}{{ end }}
    </resultMap>
    <select id="selectAll" resultMap="bean">
        select * from {{.Name}}  order by t_id
    </select>
</mapper>
`