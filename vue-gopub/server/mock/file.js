/**
 * Created by zzmhot on 2017/3/26.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/26 15:14
 * @Copyright(©) 2017 by zzmhot.
 *
 */

var Mock = require('mockjs')
var port_code = require('../../src/common/port_uri').port_code

exports.image_upload = Mock.mock({
    code: port_code.success,
    msg: "图片上传成功",
    data: {
        'id|10-100': 1,
        "name": "@ctitle",
        "image": "@image"
    }
})
