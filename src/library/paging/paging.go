package paging

import (
	"encoding/json"
	"strings"

	"library/common"

	"github.com/astaxie/beego/orm"
)

/**
 * 存储这个系统中支持的所有列处理类
 */
var PagingColumns map[string]IPagingColumn = map[string]IPagingColumn{}

/**
 * 存储这个系统中支持的所有查询条件处理类
 */
var PagingConditions map[string]IPagingCondition = map[string]IPagingCondition{}

/**
 * 列处理类
 */
type IPagingColumn interface {
	GetColumnValue(column common.Info, properties []string) string
}

/**
 * 查询条件处理类，拼查询sql
 */
type IPagingCondition interface {
	GetCondition(condition common.Info, value interface{}) string
}

/**
 * 单表增删改
 */
type Paging struct {
	db orm.Ormer
}

func NewPaging(db orm.Ormer) Paging {
	return Paging{
		db: db,
	}
}

/**
 * 返回单行数据。
 * 以关联数据的方式返回查询结果以方便程序使用。
 * 兼容DB和API两种数据源
 */
func (paging Paging) Get(schemaId string, columnIds []string, condition common.Info) common.Info {
	data := paging.Select(schemaId, columnIds, condition, 0, 1, []string{}, []string{})
	if len(data) > 0 {
		return data[0]
	}
	return nil
}

/**
 * 返回单行数据
 * 无需传入指定的columns，该方法返回schema中所有default为true的项
 * @param type schemaId
 * @param type condition
 * @return array字典非数组，如果根据查询条件没有结果则返回null
 */
func (paging Paging) GetDefault(schemaId string, condition common.Info) common.Info {
	table := paging.SelectDefault(schemaId, condition)
	if len(table) > 0 {
		return table[0]
	}
	return nil
}

/**
 * 返回多行数据
 * 无需传入指定的columns，该方法返回schema中所有default为true的项
 * @param type schemaId
 * @param type condition
 * @return type array数组非字典，数组中的每一项都是一个字典
 */
func (paging Paging) SelectDefault(schemaId string, condition common.Info) []common.Info {
	//拼出所有default为true的column
	defaultChecks := []string{}
	schema := GetSchema(schemaId)
	for _, column := range schema.InfoList("columns") {
		if column.Bool("default") {
			defaultChecks = append(defaultChecks, column.String("id"))
		}
	}
	return paging.Select(schemaId, defaultChecks, condition, 0, -1, []string{}, []string{})
}

/**
 * 获取多行数据，返回数组，其中每一项是一个字典非数组
 * 经过列处理，兼容DB和API两种数据源
 */
func (paging Paging) Select(schemaId string, columnIds []string, condition common.Info, start int64, pageSize int64, orderby []string, isAsc []string) []common.Info {
	data := paging.SelectList(schemaId, columnIds, condition, start, pageSize, orderby, isAsc)
	ret := []common.Info{}
	for i := 0; i < len(data); i++ {
		rowData := common.Info{}
		// 第一列是PrimaryKey
		rowData["PrimaryKey"] = data[i][0]
		for j := 0; j < len(columnIds); j++ {
			rowData[columnIds[j]] = data[i][j+1]
		}
		ret = append(ret, rowData)
	}
	return ret
}

/**
 * 获取多行数据，返回数组，其中每一项是一个数组非字典，只有数据无表头
 * 经过列处理
 */
func (paging Paging) SelectList(schemaId string, columnIds []string, condition common.Info, start int64, pageSize int64, orderby []string, isAsc []string) [][]string {
	schema := GetSchema(schemaId)
	columns := FillColumn(schema, columnIds)
	conditions := FillCondition(schema, condition)
	rawData := GetListTable(paging.db, schema, columns, conditions, start, pageSize, orderby, isAsc)
	return fillList(schema, rawData, columns)
}

/**
 * 调用通用查询配置文件中设定的handler类对列进行处理
 */
func fillList(schema common.Info, rawData []orm.ParamsList, columns []common.Info) [][]string {
	result := [][]string{}
	for _, row := range rawData {
		rowNew := []string{}
		// 第0列是key列
		rowNew = append(rowNew, common.GetString(row[0]))
		// columns的下标
		num := 0
		for i := 1; i < len(row) && num < len(columns); {
			parms := []string{}
			// 当前column需要显示的列数据，根据properties的个数确定需要取row中的几列
			for j := 0; j < len(strings.Split(columns[num].String("properties"), ";")); j++ {
				parms = append(parms, common.GetString(row[i]))
				i++
			}
			if columns[num]["handler"] != nil {
				rowNew = append(rowNew, columns[num]["handler"].(IPagingColumn).GetColumnValue(columns[num], parms))
			} else {
				if len(parms) == 1 {
					rowNew = append(rowNew, parms[0])
				} else {
					str, _ := json.Marshal(parms)
					rowNew = append(rowNew, string(str))
				}
			}
			num++
		}
		result = append(result, rowNew)
	}
	return result
}

/**
 * 返回结果集总记录数量
 */
func (paging Paging) Count(schemaId string, columnIds []string, condition common.Info) int {
	schema := GetSchema(schemaId)
	columns := FillColumn(schema, columnIds)
	conditions := FillCondition(schema, condition)
	return GetCount(paging.db, schema, columns, conditions)
}
