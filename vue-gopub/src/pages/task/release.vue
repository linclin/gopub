<template>
    <div class="panel">
        <panel-title :title="$route.meta.title"></panel-title>
        <div class="panel-body"
             v-loading="load_data"
             element-loading-text="拼命加载中">
            <el-form label-width="100px">
                <el-row>
                    <el-col :span="12">
                        <el-form-item label="任务标题:">
                            {{task.Title}}
                        </el-form-item>
          <span v-if='project.RepoType=="git"'>
            <div v-if="task.Branch===''">
                       <el-form-item label="Tag标签:">
                           {{task.CommitId}}
                       </el-form-item>
            </div>
            <div v-else>
                     <el-form-item label="分支:">
                          {{task.Branch}}
                      </el-form-item>
                       <el-form-item label="哈希:">
                           {{task.CommitId}}
                       </el-form-item>
            </div>
          </span>
         <span v-if='project.RepoType=="file"'>
                        <el-form-item label="包地址:">
                            {{task.CommitId}}
                        </el-form-item>
         </span>
          <span v-if='project.RepoType=="jenkins"'>
                      <el-form-item label="构建名称:">
                            {{task.Branch}}
                        </el-form-item>
                        <el-form-item label="包地址:">
                            {{task.CommitId}}
                        </el-form-item>
         </span>
                        <el-form-item v-if="!is_log">
                            <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading">部署
                            </el-button>
                            <el-button @click="$router.back()">取消</el-button>
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">

                        <el-form-item label="项目名称:">
                            {{project.Name}}
                        </el-form-item>
                        <el-form-item label="环境:">
                            {{levelEnv}}
                        </el-form-item>
                        <el-form-item label="部署路径:">
                            {{project.ReleaseTo}}
                        </el-form-item>
                        <el-form-item label="发布版本库:">
                            {{project.ReleaseLibrary}}
                        </el-form-item>
                        <el-form-item label="发布服务器组:" v-if="groups.length>0">
                            <span v-for="n in groups">{{ n }} <br></span>
                        </el-form-item>
                        <el-form-item label="发布ip:">
                            <span v-for="n in getHost">{{ n }} <br></span>
                        </el-form-item>

                    </el-col>
                </el-row>
            </el-form>
            <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
                <el-tab-pane label="版本区别" name="verLog">
                    <el-table
                            :data="changes"
                            v-loading="load_data"
                            element-loading-text="拼命加载中"
                            border
                            style="width: 100%;">
                        <el-table-column
                                prop="path"
                                label="文件">
                        </el-table-column>
                        <el-table-column
                                prop="name"
                                label="修改人">
                        </el-table-column>
                        <el-table-column
                                prop="time"
                                label="时间">
                        </el-table-column>
                    </el-table>    
                </el-tab-pane>
                <el-tab-pane label="上线进度"  name="publishProcess">
                    <terminal :taskId="task.Id"></terminal>
                </el-tab-pane>
            </el-tabs>

            
        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle, terminal} from 'components'
    import {port_task, port_conf, port_code} from 'common/port_uri'
    import {tools_verify} from 'common/tools'
    import store from 'store'
    export default{
        data(){
            return {
                task: {},
                project: {},
                changes: [],
                activeName:"verLog",
                ips:[],
                groups:[],
                form: {
                    Branch: null,
                    Title: null,
                    CommitId: null,
                    ProjectId: this.$route.params.id * 1,
                    UserId: store.state.user_info.user.Id
                },
                route_id: this.$route.params.id,
                is_log: this.$route.params.is_log,
                load_data: false,
                on_submit_loading: false,
                rules: {
                    Branch: [{required: true, message: '分支不能为空', trigger: 'blur'}],
                    CommitId: [{required: true, message: 'Commit不能为空', trigger: 'blur'}],
                    Title: [{required: true, message: '标题不能为空', trigger: 'blur'}]
                }
            }
        },
        computed: {
            getHost: function () {
                //在get_task中会取基于服务器组的ip，如果取不得则原有的hosts ip列表
                if(this.ips.length>0){
                    return this.ips
                }else{
                    var hosts=[]
                    if(this.task.Hosts && this.task.Hosts!=""){
                        hosts=this.task.Hosts.split("\r\n")
                    }else{
                      if(this.project.Hosts && this.project.Hosts!=""){
                        hosts=this.project.Hosts.split("\r\n")
                      }
                    }
                    return hosts
                }
            },
            levelEnv: function () {
                var env = ""
                if (this.project.Level == 1) {
                    env = "测试环境"
                }
                if (this.project.Level == 2) {
                    env = "预发布环境"
                }
                if (this.project.Level == 3) {
                    env = "正式环境"
                }
                return env
            }
        },
        created(){

            if (this.route_id) {
                this.get_task()
                this.get_changes()
            } else {
                this.$message({
                    message: "任务id不存在",
                    type: 'warning'
                })
                setTimeout(() => {
                    this.$router.back()
            },
                500
            )
            }
        },
        methods: {
            get_task(){
                this.load_data = true
                this.$http.get(port_task.get, {
                            params: {
                                taskId: this.route_id,
                            }
                        })
                        .then(({data: {data}}) => {
                    this.task = data
                    
                    //如果HostGroup有数据，说明使用jumpserver接口管理服务器，则生成ip列表
                    if(data.HostGroup!=""){
                        this.$http.get(port_conf.groupinfo, {
                                    params: {
                                        hostgroup: data.HostGroup
                                    }
                                })
                                .then(({data: {data}}) => {
                                this.ips=data.ips
                                this.groups=[]
                                for (var i in data.id2groupname){
                                    this.groups.push(data.id2groupname[i])
                                }
                        })
                    }
                this.get_project()
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            get_project(){
                this.load_data = true
                this.$http.get(port_conf.get, {
                            params: {
                                projectId: this.task.ProjectId
                            }
                        })
                        .then(({data: {data}}) => {
                    this.project = data
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            get_changes(){
                this.load_data = true
                this.$http.get(port_task.changes, {
                            params: {
                                taskId: this.route_id
                            }
                        })
                        .then(({data: {data}}) => {
                    this.changes = data
                this.load_data = false
                })       
            },
            //提交
            on_submit_form(){
                this.on_submit_loading = true
                this.$http.get(port_task.release, {
                            params: {
                                taskId: this.route_id
                            }
                        })
                        .then(({data: {data}}) => {
                    this.$message({
                    message: "部署开始",
                    type: 'success'
                })
                this.activeName = "publishProcess"
                this.on_submit_loading = false
            })
            .
                catch(() => {
                    this.on_submit_loading = false
            })
            }
        },
        components: {
            panelTitle,
            terminal
        }
    }
</script>
