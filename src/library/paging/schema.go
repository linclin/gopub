package paging

import (
	"github.com/linclin/gopub/src/library/common"
)

var pagingSchemaCache map[string]common.Info = map[string]common.Info{}

func ClearPagingConfig() {
	pagingSchemaCache = map[string]common.Info{}
}

var defaultSchema common.Info = common.Info{
	"id":             "",
	"name":           "",
	"keyColumn":      "",
	"useDistinct":    false,
	"groupby":        "",
	"defaultOrderby": "",
	"whereSql":       "",
	"columns":        []common.Info{},
	"conditions":     []common.Info{},
	"joins":          []common.Info{},
}

var defaultColumn common.Info = common.Info{
	"id":           "",
	"name":         "",
	"properties":   "",
	"default":      false,
	"handler":      "",
	"handlerParms": common.Info{},
	"joinId":       "",
}

var defaultCondition common.Info = common.Info{
	"id":             "",
	"name":           "",
	"properties":     "",
	"condition":      "Default",
	"conditionParms": common.Info{},
	"joinId":         "",
}

var defaultJoin common.Info = common.Info{
	"id":          "",
	"joinFormat":  "",
	"beforeJoin":  "",
	"groupConcat": false, // 是否使用group concat
	"useCount":    true,  // 计算count时是否连接此表
}

func GetSchema(schemaId string) common.Info {
	if pagingSchemaCache[schemaId] != nil {
		return pagingSchemaCache[schemaId]
	}

	path := "./schema/search/" + schemaId + ".json"
	ret := common.ReadJson(path)
	// 为指定配置信息补充没有指定的默认信息
	ret.Merge(defaultSchema)
	for _, column := range ret.InfoList("columns") {
		column.Merge(defaultColumn)
		if column.String("handler") == "" {
			column["handler"] = nil
		} else if handler, ok := PagingColumns[column.String("handler")]; ok {
			column["handler"] = handler
		} else {
			panic("schema:[" + schemaId + "]不存在指定ColumnHandler:" + column.String("handler"))
		}
	}
	for _, condition := range ret.InfoList("conditions") {
		condition.Merge(defaultCondition)
		if condition.String("condition") == "" {
			condition["condition"] = PagingConditions["Default"]
		} else if handler, ok := PagingConditions[condition.String("condition")]; ok {
			condition["condition"] = handler
		} else {
			panic("schema:[" + schemaId + "]不存在指定Condition:" + condition.String("condition"))
		}
	}
	for _, join := range ret.InfoList("join") {
		join.Merge(defaultJoin)
	}

	pagingSchemaCache[schemaId] = ret
	return ret
}

/**
 * 根据column的ID数组获得对应的column实体数组
 */
func FillColumn(schema common.Info, columnIds []string) []common.Info {
	columns := []common.Info{}
	if len(columnIds) == 0 {
		return columns
	}

	columnDic := map[string]common.Info{}
	for _, column := range schema.InfoList("columns") {
		columnDic[column.String("id")] = column
	}

	for _, columnId := range columnIds {
		if columnDic[columnId] == nil {
			panic("schema[" + schema.String("id") + "]中不存在指定id[" + columnId + "]的结果列")
		}
		columns = append(columns, columnDic[columnId])
	}
	return columns
}

/**
 * 根据condition的ID数组获得对应的condition实体数组，并新增一个values属性
 */
func FillCondition(schema common.Info, conditionIdValues common.Info) []common.Info {
	conditions := []common.Info{}
	if len(conditionIdValues) == 0 {
		return conditions
	}

	for key, value := range conditionIdValues {
		// 普通条件
		for _, condition := range schema.InfoList("conditions") {
			if key == condition["id"] || key == condition["properties"] {
				condition["values"] = value
				conditions = append(conditions, condition)
			}
		}
	}
	return conditions
}
