package core

var toLuaGlobalTableTemplate=`
CfgTable = {}

setmetatable(CfgTable,{
    __index = function(t, k)
        CfgTable[k] = require(k)
    end
})
`
