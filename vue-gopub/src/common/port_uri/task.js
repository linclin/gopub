//配置数据列表
exports.list = process.env.API_URL+"/api/get/task/list"
//根据id查询数据
exports.get = process.env.API_URL+"/api/get/task/get"
//根据id删除数据
exports.del = process.env.API_URL+"/api/get/task/del"
//添加或修改数据
exports.save = process.env.API_URL+"/api/post/task/save"
//上线接口
exports.release = process.env.API_URL+"/api/get/walle/release"
//创建回滚任务接口
exports.rollback = process.env.API_URL+"/api/get/task/rollback"
//创建回滚任务接口
exports.flush = process.env.API_URL+"/api/get/walle/flush"
//图表接口
exports.chart = process.env.API_URL+"/api/get/task/chart"
//版本变更列表
exports.changes = process.env.API_URL+"/api/get/task/changes"
   
