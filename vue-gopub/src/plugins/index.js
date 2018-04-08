/**
 * Created by zzmhot on 2017/3/23.
 *
 * 自定义插件
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/23 18:29
 * @Copyright(©) 2017 by zzmhot.
 *
 */
import DateFormat from 'plugins/date'

const install = function (Vue) {
    if (install.installed) return
    install.installed = true
    Date.prototype.$DateFormat = DateFormat
}

export default {
    install
}
