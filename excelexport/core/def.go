package core

//原始表sheet数据
type Sheet struct {
    Name        string  //sheet name
    Data[][]    string  //原始数据
    RowCount    int     //行数量
    ColCount    int     //列数量
}

//table原始数据
type Table struct {
    Name        string         //table name
    Sheet[]     *Sheet
}

//数据表头信息
type DataSheetColInfo struct{
    IsKey       bool
    IsExport    bool
    Name        string
    TypeName    string
    Describe1   string
    Describe2   string
}

//数据表表信息
type DataSheet struct{
    Name        string
    Infos[]     *DataSheetColInfo
    Data        [][]string
    RowCount    int     //行数量
    ColCount    int
}