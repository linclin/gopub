/**
 * Created by zzmhot on 2017/3/24.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/24 14:56
 * @Copyright(©) 2017 by zzmhot.
 *
 */

//获取用户信息
exports.info = process.env.API_URL+"/api/get/user/info"
//用户登录
exports.login = process.env.API_URL+"/login"
//用户登出
exports.logout = process.env.API_URL+"/logout"
//修改用户密码
exports.changepasswd = process.env.API_URL+"/changePasswd"
    //用户注册
exports.register = process.env.API_URL+"/register"

exports.users = process.env.API_URL+"/api/get/user"
exports.usersProject = process.env.API_URL+"/api/get/user/project"
