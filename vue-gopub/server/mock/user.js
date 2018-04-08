/**
 * Created by zzmhot on 2017/3/21.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/21 10:55
 * @Copyright(©) 2017 by zzmhot.
 *
 */

var Mock = require('mockjs')
var port_code = require('../../src/common/port_uri').port_code

var user_info = {
    'name': '@cname',
    'avatar': 'https://avatars0.githubusercontent.com/u/16893554?v=3&s=240',
    'age|20-25': 20,
    'desc': '@csentence()'
}

var is_login = Math.random() >= 0.5

exports.login = Mock.mock({
    code: port_code.success,
    msg: "登录成功",
    data: user_info
})
exports.logout = Mock.mock({
    code: port_code.success,
    msg: "退出成功"
})

exports.info = Mock.mock({
    code: is_login ? port_code.success : port_code.unlogin,
    msg: is_login ? "获取成功" : "您还没有登录，请登录！",
    data: is_login ? user_info : null
})
