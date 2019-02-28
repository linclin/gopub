<template>
    <div class="panel">
        <panel-title :title="$route.meta.title"></panel-title>
        <div class="panel-body"
             v-loading="load_data"
             element-loading-text="拼命加载中">
            <el-row>
                <el-col :span="20">
                    <el-form ref="form" :model="form" :rules="rules" label-width="100px">
                        <el-form-item label="上线单标题:" prop="Title">
                            <el-input v-model="form.Title" placeholder="请输入标题" style="width: 500px;"></el-input>
                        </el-form-item>
                      <el-form-item label="发布方式:" prop="isShowValue">
                        <el-radio-group v-model="isShowValue" @change="showStatus">
                          <el-radio :label="1">tag发布</el-radio>
                          <el-radio :label="2">分支发布</el-radio>
                        </el-radio-group>
                      </el-form-item>
                        <div v-if="isShowStatus">
                        <el-form-item label="选取tag:" prop="CommitId" label-width="100px">
                          <el-select v-model="form.CommitId" filterable placeholder="请选择"
                                     style="width: 400px;">
                            <el-option
                              v-for="item in tagData"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value">
                            </el-option>
                          </el-select>
                          <el-button @click.stop="get_tag_data" size="small">
                            <i class="fa fa-refresh"></i>
                          </el-button>
                        </el-form-item>
                        </div>
                        <div v-else>
                        <el-form-item label="选取分支:" prop="Branch" label-width="100px">
                            <el-select v-model="form.Branch" filterable placeholder="请选择" @change="get_commit"
                                       style="width: 400px;">
                                <el-option
                                        v-for="item in branchData"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                </el-option>
                            </el-select>
                            <el-button @click.stop="get_branch_data" size="small">
                                <i class="fa fa-refresh"></i>
                            </el-button>
                        </el-form-item>
                        <el-form-item label="版本选取 :" label-width="100px" prop="CommitId">
                            <el-select v-model="form.CommitId" filterable placeholder="请选择" style="width: 400px;">
                                <el-option
                                        v-for="item in commitData"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        </div>
                      <el-form-item  label="灰度发布 :" >
                        <el-switch v-model="isShowHost" on-text="on" off-text="off">灰度发布</el-switch>
                        <div>
                          <el-select v-if="isShowHost" v-model="selectHosts"  multiple filterable placeholder="请选择" style="width: 400px;">
                            <el-option
                              v-for="item in Hosts"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value">
                            </el-option>
                          </el-select>
                        </div>

                      </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading">立即提交
                            </el-button>
                            <el-button @click="$router.back()">取消</el-button>
                        </el-form-item>
                    </el-form>
                </el-col>
            </el-row>
        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle} from 'components'
    import {port_task, port_git, port_code,port_conf} from 'common/port_uri'
    import {tools_verify} from 'common/tools'
    import store from 'store'
    export default{
        data(){
            return {
                commitData: [],
                branchData: [],
                tagData: [],
                isShowHost:false,
                isShowStatus:false,
                Hosts:[],
                selectHosts:[],
                ProjectData:null,
                isShowValue: 2,
                form: {
                    Branch: null,
                    Tag: null,
                    Title: null,
                    CommitId: null,
                    Hosts:null,
                    ProjectId: this.$route.query.id * 1,
                    UserId: store.state.user_info.user.Id

                },
                route_id: this.$route.query.id,
                load_data: false,
                on_submit_loading: false,
                rules: {
                    Tag: [{required: true, message: '分支不能为空', trigger: 'blur'}],
                    Branch: [{required: true, message: '分支不能为空', trigger: 'blur'}],
                    CommitId: [{required: true, message: 'Commit不能为空', trigger: 'blur'}],
                    Title: [{required: true, message: '标题不能为空', trigger: 'blur'}]
                }
            }
        },
        created(){

            if (this.route_id) {
                this.get_tag_data()
                this.get_branch_data()
                this.get_Project_data()
            } else {
                this.$message({
                    message: "项目id不存在",
                    type: 'warning'
                })
                setTimeout(() => {
                    this.$router.push({
                    name: 'taskMyList'
                })
            },
                500
            )
            }
        },
        methods: {
          get_Project_data(){
            this.load_data = true
            this.$http.get(port_conf.get, {
              params: {
                projectId: this.form.ProjectId
              }
            })
              .then(({data: {data}}) => {
                this.ProjectData = data
                this.Hosts=[]
                var ss=this.ProjectData.Hosts.match(/((\d+)\.(\d+)\.(\d+)\.(\d+)\:(\d+))|((\d+)\.(\d+)\.(\d+)\.(\d+))/g)
                for(var i=0;i<ss.length;i++){
                  this.Hosts.push({label:  ss[i], value:  ss[i]})
                }
                this.load_data = false
              })
              .
              catch(() => {
                this.load_data = false
              })
          },
            get_commit(){
                this.load_data = true
                this.$http.get(port_git.commit, {
                            params: {
                                projectId: this.route_id,
                                branch: this.form.Branch
                            }
                        })
                        .then(({data: {data}}) => {
                    var commitData = []
                    var CommitId = ""
                    for(var i=0;i<data.length;i++){
                        if(i==0){
                            CommitId=data[i].id
                        }
                        commitData.push({label: data[i].message, value: data[i].id})
                    }
                this.commitData = commitData
                this.CommitId = CommitId
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            get_branch_data(){
                this.load_data = true
                this.commitData = []
                this.form.Branch = null
                this.form.CommitId = null
                this.$http.get(port_git.branch, {
                            params: {
                                projectId: this.route_id
                            }
                        })
                        .then(({data: {data}}) => {
                    var branchData = []
                    for(var i in data
            )
                {
                    branchData.push({label: data[i].message, value: data[i].id})
                }
                this.branchData = branchData
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            get_tag_data(){
              this.load_data = true
              this.$http.get(port_git.getTag, {
                params: {
                  projectId: this.form.ProjectId
                }
              })
                .then(({data: {data}}) => {
                  var tagData = []
                  for(var i in data)
                  {
                    if( data[i].id !== "") {
                      tagData.push({label: data[i].message, value: data[i].id})
                    }
                  }
                  console.log('--->', tagData)
                  this.tagData = tagData
                  this.load_data = false
                })
                .
                catch(() => {
                  this.load_data = false
                })
            },
            //提交
            on_submit_form(){
                this.$refs.form.validate((valid) => {
                    if (
                !valid
            )
                return false
                this.on_submit_loading = true
                  if(this.isShowHost){
                      this.form.Title=this.form.Title+"-灰度"
                      this.form.Hosts=this.selectHosts.toString()
                  }
                this.$http.post(port_task.save, this.form)
                        .then(({data: {data}}) => {
                    this.$message({
                    message: "修改成功",
                    type: 'success'
                })
                setTimeout(() => {
                    this.$router.push({
                    name: 'taskMyList'
                })
            },
                500
            )
            })
            .
                catch(() => {
                    this.on_submit_loading = false
            })
            })
            },
          showStatus(){
            if (this.isShowValue === 1){
              this.isShowStatus = true
              this.form.CommitId = null
              this.form.Branch = null
            } else {
              this.isShowStatus = false
              this.form.CommitId = null
              this.form.Branch = null
            }
          }
        },
        components: {
            panelTitle
        }
    }
</script>
