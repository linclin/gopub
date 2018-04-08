/**
 * Created by zzmhot on 2017/3/21.
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/21 10:49
 * @Copyright(©) 2017 by zzmhot.
 *
 */

var express = require('express')
var apiRouter = express.Router()

require('./user')(apiRouter)
require('./table')(apiRouter)
require('./file')(apiRouter)

module.exports = apiRouter
