/**
 * @file: index.
 * @intro: axios配置.
 * @author: zzmhot.
 * @email: zzmhot@163.com.
 * @Date: 2017/4/27 17:48.
 * @Copyright(©) 2017 by thinkive.
 *
 */

//导入模块
import axios from 'axios'
import {port_code} from 'common/port_uri'
import router from 'src/router'
import store from 'store'
import {SET_USER_INFO} from 'store/actions/type'

//设置用户信息action
const setUserInfo = function (user) {
    store.dispatch(SET_USER_INFO, user)
}

const install = function (Vue) {
    if (install.installed) return
    install.installed = true

  //设置默认根地址
  axios.defaults.baseURL = '/'
  //设置请求超时设置
  axios.defaults.timeout = 600000


    // http request 拦截器
    axios.interceptors.request.use(
        config => {
        if (store.state.user_info && store.state.user_info.user && store.state.user_info.user.AuthKey
    )
    {
        config.headers.Authorization = `token ${store.state.user_info.user.AuthKey}`;
    }
    return config;
},
    err =>{
        return Promise.reject(err);
    }
    )
    ;

    /**
     * 添加响应拦截器
     */
    axios.interceptors.response.use(response => {
        //成功时
        let resData = response.data
        let dataCode = resData.code
        let datamsg = resData.msg
        if (dataCode === port_code.success
    )
    {
        return Promise.resolve(response)
    }
    else
    if (dataCode === port_code.unlogin) {
        setUserInfo(null)
        router.replace({name: "login"})
    }
    if (datamsg == null || datamsg == "") {
        return Promise.resolve(response)
    } else {
        Vue.prototype.$message({
            message: datamsg,
            type: 'warning'
        })
    }
    return Promise.reject({code: dataCode, msg: datamsg})
},
    error =>
    {
        //错误时
        if (error.response) {
            let resError = error.response
            let resCode = resError.status
            let resMsg = error.message
            Vue.prototype.$message({
                message: '操作失败！错误原因 ' + resMsg,
                type: 'error'
            })
            return Promise.reject({code: resCode, msg: resMsg})
        }
    }
    )
    ;

    //设置到axios到Vue上
    Vue.axios = axios
    Object.defineProperties(Vue.prototype, {
        axios: {
            get() {
                return axios
            }
        },
        $http: {
            get() {
                return axios
            }
        }
    })
}

export default {
    install
}
