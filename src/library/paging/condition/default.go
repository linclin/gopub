package condition

import (
	"library/common"
	"library/paging"
	"strings"
)

func init() {
	paging.PagingConditions["Default"] = Default{}
}

type Default struct{}

/**
 * 如果properties中有多个，以半角 分号 隔开的字段，则需要拆开处理成or关系
 */
func (def Default) GetCondition(condition common.Info, value interface{}) string {
	// 这个扩展暂时只支持一个属性，不支持多属性
	key := condition.String("properties")

	// 根据值格式化
	// 半角分号隔开格式化成 in
	// 数组格式化成 指定key的操作方式
	// 字符串格式化成 =
	switch value.(type) {
	case common.Info:
		ret := ""
		for operator, val := range value.(common.Info) {
			ret += " and " + convertOperator(key, operator, val.(string))
		}
		return ret[5:len(ret)]
	case string:
		values := strings.Split(value.(string), ";")
		if len(values) == 1 {
			if value == "#empty#" {
				return key + " IS NULL"
			} else if value == "#notempty#" {
				return key + " IS NOT NULL"
			} else {
				return key + "='" + value.(string) + "'"
			}
		} else {
			// value中包含";"则返回in条件
			return key + " in ('" + strings.Join(values, "','") + "') "
		}
	default:
		return key + "='" + common.ToString(value) + "'"
	}
}

/**
 * return "column operator value";
 */
func convertOperator(column string, operator string, value string) string {
	if strings.ToLower(operator) == "like" {
		return column + " " + operator + " '%" + value + "%'"
	}
	if strings.ToLower(value) == "now()" || strings.ToLower(value) == "null" {
		return column + " " + operator + " " + value
	}
	return column + " " + operator + " '" + value + "'"
}
