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
                        <el-form-item label="url:" prop="CommitId" label-width="100px">
                            <el-tooltip class="item" effect="dark" content='若有http开头则不拼接配置文件的地址' placement="top">
                                <el-input v-model="form.CommitId" @change="change_data" placeholder="请输入url"
                                          style="width: 500px;"></el-input>
                            </el-tooltip>     <el-button @click.stop="get_md5_data" size="small">
                            <i class="fa fa-refresh"></i> 检查MD5
                        </el-button>
                            <div>md5：<span v-text="md5" ></span></div>
                        </el-form-item>

                        <el-form-item label="md5 :" label-width="100px">
                            <el-tooltip class="item" effect="dark" content='当MD5不为空时，会检查下载文件的nd5' placement="top">
                                <el-input v-model="form.FileMd5" @change="change_md5data" placeholder="请输入md5"
                                          style="width: 500px;"></el-input>
                            </el-tooltip>
                        </el-form-item>
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
    import {port_task, port_file, port_code,port_conf} from 'common/port_uri'
    import {tools_verify} from 'common/tools'
    import store from 'store'
    export default{
        data(){
            return {
                isShowHost:false,
                Hosts:[],
                selectHosts:[],
                ProjectData:null,
                md5: null,
                fileName: "",
                isRightMd5:false,
                isChange: true,
                form: {
                    FileMd5: null,
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
                    Title: [{required: true, message: '标题不能为空', trigger: 'blur'}]
                }
            }
        },
        created(){

            if (this.route_id) {
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
                var ss=this.ProjectData.Hosts.match(/(\d+)\.(\d+)\.(\d+)\.(\d+)\:(\d+)/g)
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
            change_data(){
                this.isChange=true;
            },
            change_md5data(){
                this.isChange=true;
            },
            get_md5_data(){
                this.load_data = true
                this.md5 = null
                this.$http.get(port_file.md5, {
                            params: {
                                projectId: this.route_id,
                                url: this.form.CommitId
                            }
                        })
                        .then(({data: {data}}) => {
                    this.md5 = data[0]["message"]
                this.fileName = data[0]["id"]
                this.isChange=false;
                this.check_md5()
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            check_md5(){
                if (this.form.FileMd5 && this.form.FileMd5.trim()&&this.form.FileMd5 != this.md5) {
                    this.isRightMd5 = false
                    this.$message({
                        message: "md5不相等，" + this.fileName + ":" + this.md5,
                        type: 'warning'
                    })
                }else{
                    this.isRightMd5 = true
                    this.$message({
                        message: "验证成功",
                        type: 'success'
                    })
                }

            },
            //提交
            on_submit_form(){
                this.$refs.form.validate((valid) => {
                    if (!valid){
                    return false
                }
                if (!this.isRightMd5 ||  this.isChange) {
                    this.$message({
                        message: "没有检查MD5",
                        type: 'warning'
                    })
                    return false
                }
                  if(this.isShowHost){
                    this.form.Title=this.form.Title+"-灰度"
                    this.form.Hosts=this.selectHosts.toString()
                  }
                this.on_submit_loading = true
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
            }
        },
        components: {
            panelTitle
        }
    }
</script>
