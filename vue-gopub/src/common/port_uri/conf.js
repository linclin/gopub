//配置数据列表
exports.list = process.env.API_URL+"/api/get/conf/list"
exports.mylist = process.env.API_URL+"/api/get/conf/mylist"
//根据id查询数据
exports.get = process.env.API_URL+"/api/get/conf/get"
//根据id删除数据
exports.del = process.env.API_URL+"/api/get/conf/del"
//添加或修改数据
exports.save = process.env.API_URL+"/api/post/conf/save"
//检查
exports.detection = process.env.API_URL+"/api/get/walle/detection"
//复制项目
exports.copy = process.env.API_URL+"/api/get/conf/copy"
exports.tags = process.env.API_URL+"/api/get/conf/tags"
exports.lock = process.env.API_URL+"/api/get/conf/lock"
exports.server_groups = process.env.API_URL+"/api/get/conf/server_groups"
exports.groupinfo = process.env.API_URL+"/api/get/conf/groupinfo"

