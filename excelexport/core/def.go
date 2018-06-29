package core

//原始表sheet数据
type Sheet struct {
	Name     string     //sheet name
	Data     [][]string //原始数据
	RowCount int        //行数量
	ColCount int        //列数量
}

//table原始数据
type Table struct {
	Name  string //table name
	Sheet []*Sheet
}

//命令行结构
type DataSheetColCommond struct {
	Name     string   //命令名称
	paras    []string //命令参数
	ParamStr string   //命令参数字符串
}

//数据表头信息
type DataSheetColInfo struct {
	IsKey     	bool
	IsExport  	bool
	Name      	string
	SrcTypeName string
	TypeName  	string
	Describe1 	string
	Describe2 	string
	Commonds  	[]*DataSheetColCommond
}

//数据表表信息
type DataSheet struct {
	TableName string
	Name     string
	Infos    []*DataSheetColInfo
	Data     [][]string
	RowCount int //行数量
	ColCount int
	indexs	map[string] map[string] int	//建立索引 [idkey](key index)
}
