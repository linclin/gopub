/**
 * Created by zzmhot on 2017/3/26.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/26 15:17
 * @Copyright(©) 2017 by zzmhot.
 *
 */

var mock = require('../mock/file')
var uri = require('../../src/common/port_uri').port_file

module.exports = function (apiRouter) {
    //图片上传
    apiRouter.post(uri.image_upload, function (req, res) {
        setTimeout(function () {
            res.json(mock.image_upload)
        }, 1000)
    })
}
