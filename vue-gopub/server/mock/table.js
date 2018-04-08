/**
 * Created by zzmhot on 2017/3/24.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/24 16:50
 * @Copyright(©) 2017 by zzmhot.
 *
 */

var Mock = require('mockjs')
var port_code = require('../../src/common/port_uri').port_code

// 姓名	性别	年龄	生日	地址

var data_list = [{
    'id|10-100': 1,
    "name": "@cname",
    "sex": "@pick([1, 2])",//1男，2女
    'age|20-25': 20,
    "birthday": "@date('yyyy-MM-dd')",
    "address": "@county(true)",
    "zip": "@zip"
}]

exports.list = Mock.mock({
    code: port_code.success,
    msg: "获取成功",
    "data|15": data_list,
    page: 1,
    "total|50-200": 50
})

exports.get = Mock.mock({
    code: port_code.success,
    msg: "获取成功",
    data: data_list[0]
})

exports.del = Mock.mock({
    code: port_code.success,
    msg: "删除成功"
})

exports.save = Mock.mock({
    code: port_code.success,
    msg: "操作成功"
})

exports.batch_del = Mock.mock({
    code: port_code.success,
    msg: "批量删除成功"
})
