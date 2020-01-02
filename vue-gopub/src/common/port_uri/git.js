//获取分支
exports.branch = process.env.API_URL+"/api/get/git/branch"
//获取tag
exports.getTag = process.env.API_URL+"/api/get/git/tag"
//获取提交
exports.commit = process.env.API_URL+"/api/get/git/commit"

exports.gitlog = process.env.API_URL+"/api/get/git/gitlog"
exports.gitpull = process.env.API_URL+"/api/get/git/gitpull"

//获取分支
exports.jenkinsBranch = process.env.API_URL+"/api/get/jenkins/commit"
