package paging

import (
	"github.com/linclin/gopub/src/library/common"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
 * 根据要显示的列以及用到的查询条件，拼凑并执行SQL，拿到原始数据
 * 返回由property指定的列的数据内容，无表头。
 *
 * 如果传了排序字段和类型，则表明用户想要按自己定义的方式来排序；
 * 否则以schema中指定的默认排序方式为准；
 * 如果都未指定则不排序。
 */
func GetListTable(db orm.Ormer, schema common.Info, columns []common.Info, conditions []common.Info, start int64, pageSize int64, orderby []string, isAsc []string) []orm.ParamsList {
	// 用于存放columns或者conditions中可能存在的joinId
	joinIds := map[string]string{}
	sqlSelect := buildSelect(schema, columns, joinIds, false)
	sqlWhere := buildWhere(schema, conditions, joinIds)
	sqlFrom, useGroupConcat := buildFrom(schema, joinIds)

	sql := sqlSelect + sqlFrom + sqlWhere
	if useGroupConcat == true {
		sql += " group by " + schema.String("keyColumn")
	} else {
		if schema["groupby"] != "" {
			sql += " group by " + schema.String("groupby")
		}
	}
	// 如果传了排序字段和类型，则表明用户想要按自己定义的方式来排序；
	// 否则以schema中指定的默认排序方式为准；
	// 如果都未指定则不排序。
	if len(orderby) > 0 {
		if len(isAsc) > 0 {
			sql += " order by " + strings.Join(orderby, ",") + " asc "
		} else {
			sql += " order by " + strings.Join(orderby, ",") + "  " + strings.Join(isAsc, ",")
		}
	} else if schema["defaultOrderby"] != "" {
		sql += " order by " + schema.String("defaultOrderby")
	}
	if pageSize != -1 {
		sql += " limit " + strconv.Itoa(int(start)) + "," + strconv.Itoa(int(pageSize))
	}
	beego.Debug("GetListTable:", sql)
	var list []orm.ParamsList
	_, err := db.Raw(sql).ValuesList(&list)
	if err != nil {
		beego.Error(sql+"  获取失败:", err.Error())
	}
	return list
}

/**
 * 返回结果集总记录数量
 */
func GetCount(db orm.Ormer, schema common.Info, columns []common.Info, conditions []common.Info) int {
	// 用于存放columns或者conditions中可能存在的joinId
	joinIds := map[string]string{}
	distinct := " "
	if schema.Bool("useDistinct") {
		distinct = "distinct "
	}
	sqlSelect := "select count(*) from (select " + distinct + schema.String("keyColumn")
	// 如果columns中有joinid也要找出来
	buildSelect(schema, columns, joinIds, true)
	sqlWhere := buildWhere(schema, conditions, joinIds)
	sqlFrom, useGroupConcat := buildFrom(schema, joinIds)
	sql := sqlSelect + sqlFrom + sqlWhere
	if useGroupConcat == true {
		sql += " group by " + schema.String("keyColumn")
	} else {
		if schema["groupby"] != "" {
			sql += " group by " + schema.String("groupby")
		}
	}
	sql += ") as t"
	beego.Debug("GetCount:", sql)
	var count []orm.ParamsList
	_, err := db.Raw(sql).ValuesList(&count)
	if err != nil {
		beego.Error(sql+"  获取失败:", err.Error())
	}
	return common.GetInt(count[0][0])
}

/**
 * 拼成select 表名.字段
 */
func buildSelect(schema common.Info, columns []common.Info, joinIds map[string]string, isCount bool) string {
	joinDic := map[string]common.Info{}
	for _, join := range schema.InfoList("joins") {
		joinDic[join.String("id")] = join
	}

	sql := "select "
	if schema.Bool("useDistinct") {
		sql += " distinct "
	}
	sql += schema.String("keyColumn")
	for _, column := range columns {
		propertys := strings.Split(column.String("properties"), ";")
		if len(propertys) > 0 {
			sql += "," + strings.Join(propertys, ",")
		} else {
			sql += "," + column.String("properties")
		}

		if column["joinId"] != "" {
			joinTemps := strings.Split(column.String("joinId"), ";")
			for _, joinTemp := range joinTemps {
				_, ok := joinIds[joinTemp]
				if ok {
					// 已经添加的不需要再添加
					continue
				}
				joinInfo, ok := joinDic[joinTemp]
				if !ok {
					panic("schema[" + schema.String("id") + "]中的column[" + column.String("id") + "]的joinId[" + joinTemp + "]不存在")
				}
				if isCount {
					if joinInfo.Bool("useCount") {
						joinIds[joinTemp] = joinTemp
					}
					// 因为查总数时会根据useCount属性过滤掉一批表，所以这里必须先处理beforeJoin，否则在buildFrom处理时，有些join节点已经被过滤掉了
					// buildWhere对应的表尽管是查总数也必须关联，所以不需要先处理beforeJoin
					if joinInfo.String("beforeJoin") != "" {
						beforeJoins := strings.Split(joinInfo.String("beforeJoin"), ";")
						for _, beforeJoin := range beforeJoins {
							if joinDic[beforeJoin] != nil && joinDic[beforeJoin].Bool("useCount") {
								joinIds[beforeJoin] = beforeJoin
							}
						}
					}
				} else {
					// buildFrom会处理beforeJoin，这里就不需要再处理了
					joinIds[joinTemp] = joinTemp
				}
			}
		}
	}
	return sql
}

/**
 * 拼成where ***
 */
func buildWhere(schema common.Info, conditions []common.Info, joinIds map[string]string) string {
	conditionStr := ""
	if schema["whereSql"] != "" {
		conditionStr = " where " + schema.String("whereSql")
	} else {
		conditionStr = " where 1=1 "
	}
	for _, condition := range conditions {
		conditionStr += " and " + condition["condition"].(IPagingCondition).GetCondition(condition, condition["values"])

		// 将condition里的join添加到joinIds
		if condition.String("joinId") != "" {
			joinTemps := strings.Split(condition.String("joinId"), ";")
			for _, joinTemp := range joinTemps {
				if joinTemp != "" {
					joinIds[joinTemp] = joinTemp
				}
			}
		}
	}
	return conditionStr
}

/**
 * 根据join拼成from ... left join ...
 */
func buildFrom(schema common.Info, joinIds map[string]string) (string, bool) {
	actualJoinIds := map[string]string{}
	// 默认需要将第一个join拼上
	if len(schema.InfoList("joins")) > 0 {
		actualJoinIds[schema.InfoList("joins")[0].String("id")] = schema.InfoList("joins")[0].String("id")
	}
	// 处理beforejoin节点
	for _, joinId := range joinIds {
		for _, join := range schema.InfoList("joins") {
			if join["id"] == joinId && join["beforeJoin"] != "" {
				beforeJoins := strings.Split(join.String("beforeJoin"), ";")
				for _, beforeJoin := range beforeJoins {
					actualJoinIds[beforeJoin] = beforeJoin
				}
			}
		}
		actualJoinIds[joinId] = joinId
	}
	fromStr := " from "
	useGroupConcat := false
	// 顺序很重要
	for _, join := range schema.InfoList("joins") {
		if _, ok := actualJoinIds[join.String("id")]; ok {
			fromStr += join.String("joinFormat")
			if join.Bool("useGroupConcat") {
				useGroupConcat = true
			}
		}
	}
	return fromStr, useGroupConcat
}
