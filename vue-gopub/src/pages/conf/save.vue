<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body"
         v-loading="load_data"
         element-loading-text="拼命加载中">
      <el-row>
        <el-col :span="24">
          <el-form ref="form" :model="form" :rules="rules">
            <el-row>
              <el-col :span="12">
                <el-form-item label="项目名称:" prop="Name" label-width="100px">
                  <el-input v-model="form.Name" placeholder="请输入项目名称"
                            style="width: 600px;"></el-input>
                </el-form-item>
                <el-form-item label="项目标签:" prop="Tag" label-width="100px">
                  <el-select v-model="form.TagArray" filterable multiple allow-create default-first-option placeholder="请选择" style="width: 400px;">
                    <el-option
                      v-for="item in Tags"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="项目环境:" label-width="100px">
                  <el-radio-group v-model="form.Level">
                    <el-radio :label="2">预发布环境</el-radio>
                    <el-radio :label="3">线上环境</el-radio>
                  </el-radio-group>
                </el-form-item>

                <el-tabs v-if="!route_id" v-model="form.RepoType" type="card" @tab-click="handleClick">
                  <el-tab-pane label="Git" name="git">
                    <el-form-item label="地址:" prop="RepoUrl" label-width="100px">
                      <el-tooltip class="item" effect="dark"
                                  content=" git格式:ssh-url，需要把宿主机php进程用户的ssh-key加入git信任"
                                  placement="top">
                        <el-input v-model="form.RepoUrl"
                                  placeholder="git@gitee.com/dev-ops/gopub.git"
                                  style="width: 600px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                  </el-tab-pane>
                  <el-tab-pane label="File" name="file">
                    <el-form-item label="地址:" prop="RepoUrl" label-width="100px">
                      <el-tooltip class="item" effect="dark"
                                  content="发布包的http地址"
                                  placement="top">
                        <el-input v-model="form.RepoUrl"
                                  placeholder="File模式填入下载文件夹路径如：http://download.aaa.org/a/"
                                  style="width: 600px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                  </el-tab-pane>
                  <el-tab-pane label="Jenkins" name="jenkins">
                    <el-form-item label="jenkins地址:" prop="RepoUrl" label-width="100px">
                      <el-tooltip class="item" effect="dark"
                                  content="job页jenkins地址，类似http://jenkins.xxxxx.com/job/项目名称/"
                                  placement="top">
                        <el-input v-model="form.RepoUrl"
                                  placeholder="job页jenkins地址，类似http://jenkins.xxxxx.com/job/项目名称/"
                                  style="width: 600px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                  </el-tab-pane>
                </el-tabs>
                <el-form-item v-if="route_id" label="地址:" prop="RepoUrl" label-width="100px">
                  <el-tooltip class="item" effect="dark"
                              content="支持git/svn。git格式:ssh-url，需要把宿主机php进程用户的ssh-key加入git信任"
                              placement="top">
                    <el-input v-model="form.RepoUrl"
                              placeholder="git@gitee.com/dev-ops/gopub.git(File模式填入下载文件夹路径如：http://download.xxxx.org/a/)"
                              style="width: 600px;"></el-input>
                  </el-tooltip>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="8">
                <div class="panel-title">
                  宿主机
                  <div class="fr">
                    <slot></slot>
                  </div>
                </div>
                <div class="panel-body el-form el-form--label-top">

                  <el-form-item label="代码检出仓库:" prop="DeployFrom">
                    <el-tooltip class="item" effect="dark" content="代码的检出存放路径" placement="top">
                      <el-input v-model="form.DeployFrom" placeholder="/data/gopub"
                                style="width: 400px;"></el-input>
                    </el-tooltip>
                  </el-form-item>
                  <el-form-item label="排除文件:" prop="Excludes">
                    <el-tooltip class="item" effect="dark" content="剔除不上线的文件、目录，每行一个"
                                placement="top">
                      <el-input type="textarea" autosize v-model="form.Excludes" placeholder=".git
README.md" style="width: 400px;"></el-input>
                    </el-tooltip>
                  </el-form-item>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="panel-title">
                  目标机
                  <div class="fr">
                    <slot></slot>
                  </div>
                </div>
                <div class="panel-body">
                  <div class="el-form--label-top">
                    <el-form-item label="用户:" prop="ReleaseUser">
                      <el-tooltip class="item" effect="dark"
                                  content="代码的部署的用户，一般是运行的服务的用户，如php进程用户www" placement="top">
                        <el-input v-model="form.ReleaseUser" placeholder="www"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="webroot:" prop="ReleaseTo">
                      <el-tooltip class="item" effect="dark"
                                  content="代码的最终部署路径，请不要在目标机新建此目录，walle会自动生成此软链，正确设置父目级录即可"
                                  placement="top">
                        <el-input v-model="form.ReleaseTo" placeholder="/data/wwwroot/xxx"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="发布版本库:" prop="ReleaseLibrary">
                      <el-tooltip class="item" effect="dark"
                                  content="代码发布的版本库，每次发布更新webroot的软链到当前最新版本" placement="top">
                        <el-input v-model="form.ReleaseLibrary" placeholder="/data/gopub_releases"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="版本保留数:" prop="KeepVersionNum">
                      <el-tooltip class="item" effect="dark" content="过多的历史版本将被删除，只可回滚保留的版本"
                                  placement="top">
                        <el-input-number v-model="form.KeepVersionNum" placeholder="20"
                                         style="width: 400px;" :max="100" :min="1" :value="20"
                                         :controls="false"></el-input-number>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="服务器组" v-if="server_groups.length>0">
                      <el-select v-model="form.HostGroupArray" filterable multiple default-first-option placeholder="请选择" style="width: 400px;">
                    <el-option
                      v-for="item in server_groups"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value">
                    </el-option>
                  </el-select>
                    </el-form-item>








                    <el-form-item label="是否开启分批发布:">
                      <el-radio-group v-model="form.IsGroup">
                        <el-radio :label="0">关闭</el-radio>
                        <el-radio :label="1">开启</el-radio>
                      </el-radio-group>
                    </el-form-item>
                    <el-form-item v-if="!form.IsGroup" label="机器列表:" prop="Hosts">
                      <el-tooltip class="item" effect="dark" content="要发布的机器列表，一行一个，默认22端口"
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.Hosts" placeholder="192.168.0.1
192.168.0.2:21" style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item v-if="form.IsGroup" label="机器列表:" prop="Hosts">
                      <div>
                        <el-button type="primary" icon="plus" size="small"   @click="add_data()">添加</el-button>
                      </div>
                      <ul>
                        <li v-for="(value, key) in hosts">
                          <el-tooltip class="item" effect="dark" content="要发布的机器列表，一行一个，22端口"
                                      placement="top">
                            <el-input type="textarea" autosize v-model="hosts[key]" placeholder="192.168.0.1
192.168.0.2:21" style="width: 400px;"></el-input>
                          </el-tooltip>
                          <el-button type="danger" size="small" icon="delete"
                                     @click="delete_data(key)">删除
                          </el-button>
                        </li>
                      </ul>


                      <div>

                      </div>
                    </el-form-item>
                  </div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="panel-title">
                  高级任务
                  <el-tooltip class="item" effect="dark"
                              content="使用系统变量更方便处理路径问题：{WORKSPACE}：宿主机的独立部署空间或目标机的webroot    {VERSION}：发布的版本库的当前版本"
                              placement="top">
                    <i class="el-icon-search"></i>
                  </el-tooltip>
                  <div class="fr">
                    <slot></slot>
                  </div>
                </div>

                <div class="panel-body">
                  <div class="panel-body el-form el-form--label-top">
                    <el-form-item label="代码检出前任务:" prop="PreDeploy">
                      <el-tooltip class="item" effect="dark"
                                  content="在部署代码之前的准备工作，如git的一些前置检查、vendor的安装（更新），一行一条"
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.PreDeploy"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="代码检出后任务:" prop="PostDeploy">
                      <el-tooltip class="item" effect="dark"
                                  content="git代码检出之后，可能做一些调整处理，如vendor拷贝，环境适配（mv config-test.php config.php），一行一条"
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.PostDeploy"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="同步完目标机后任务:" prop="PreRelease">
                      <el-tooltip class="item" effect="dark"
                                  content='同步完所有目标机器之后，更改版本软链之前触发任务。java可能要做一些暂停服务的操作'
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.PreRelease"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="更改版本软链接后任务:" prop="PostRelease">
                      <el-tooltip class="item" effect="dark"
                                  content='所有目标机器都部署完毕之后，做一些清理工作，如删除缓存、重启服务（nginx、php、task），一行一条'
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.PostRelease"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>
                    <el-form-item label="版本发布发布完成后执行本地任务:" prop="LastDeploy">
                      <el-tooltip class="item" effect="dark"
                                  content='一行一条,相关参数{WORKSPACE},{VERSION},{HOSTS},{PROJECT_ID},{PROJECT_NAME},{TASK_ID},{TASK_LINKID}'
                                  placement="top">
                        <el-input type="textarea" autosize v-model="form.LastDeploy"
                                  style="width: 400px;"></el-input>
                      </el-tooltip>
                    </el-form-item>

                  </div>
                </div>
              </el-col>
            </el-row>

            <el-form-item label="上线方式:">
              <el-radio-group v-model="form.ReleaseType">
                <el-radio :label="0">软链接</el-radio>
                <el-radio :label="1">移动目录</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="是否开启p2p:">
              <el-radio-group v-model="form.P2p">
                <el-radio :label="0">关闭</el-radio>
                <el-radio :label="1">开启</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="是否启用gzip:">
              <el-radio-group v-model="form.Gzip">
                <el-radio :label="0">关闭</el-radio>
                <el-radio :label="1">开启</el-radio>
              </el-radio-group>
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
  import {port_conf, port_code} from 'common/port_uri'
  import {tools_verify} from 'common/tools'

  export default {
    data() {
      return {
        options: [{
          value: 'git',
          label: 'git'
        }, {
          value: 'file',
          label: 'file'
        }, {
          value: 'jenkins',
          label: 'jenkins'
        }],
        server_groups:[],
        hosts: [],
        pmsOptions: {},
        form: {
          Name: null,
          TagArray:[],
          Level: 2,
          RepoType: "git",
          RepoUrl: null,
          DeployFrom: null,
          Excludes: null,
          ReleaseUser: null,
          ReleaseTo: null,
          ReleaseType: 0,
          ReleaseLibrary: null,
          RepoPassword: null,
          RepoUsername: null,
          KeepVersionNum: 20,
          Hosts: null,
          Status: 1,
          RepoMode: "branch",
          Audit: 0,
          P2p: 0,
          Gzip: 1,
          IsGroup: 0,
          HostGroup:"",
          HostGroupArray:[]
        },
        Tags: [],
        route_id: this.$route.params.id,
        load_data: false,
        on_submit_loading: false,
        rules: {
          Name: [{required: true, message: '项目名称不能为空', trigger: 'blur'}],
          RepoUrl: [{required: true, message: '项目地址不能为空', trigger: 'blur'}],
          DeployFrom: [{required: true, message: '代码检出仓库不能为空', trigger: 'blur'}],
          ReleaseUser: [{required: true, message: '目标机器部署代码用户不能为空。', trigger: 'blur'}],
          ReleaseTo: [{required: true, message: '代码的webroot不能为空', trigger: 'blur'}],
          ReleaseLibrary: [{required: true, message: '发布版本库不能为空。', trigger: 'blur'}],
          Hosts: [{required: true, message: '机器列表不能为空。', trigger: 'blur'}],
        }
      }
    },
    created() {
      this.get_public_data()
      if (this.route_id) {
        this.get_form_data()
      }
    },
    methods: {
      add_data() {
        this.hosts.push("")
      },
      delete_data(key) {
        this.hosts.splice(key,1)
      },
      handleClick(tab, event) {
        this.RepoType = tab.name
      },
      get_public_data(){
        //当前已有标签列表
        this.$http.get(port_conf.tags, {
                params: {
                }
            })
            .then(({data: {data}}) => {
          for (var i in data){
            this.Tags.push({
              'label':data[i],
              'value':data[i]
            })
          }
        })

        //服务器组
        this.$http.get(port_conf.server_groups, {
                params: {
                }
            })
            .then(({data: {data}}) => {
        
          for (var i in data){
            this.server_groups.push({
              value:i+"",
              label:data[i]
            })
          }

        })
      },
      //获取数据
      get_form_data() {
        this.load_data = true
        

        this.$http.get(port_conf.get, {
          params: {
            projectId: this.route_id
          }
        })
          .then(({data: {data}}) => {
            data.TagArray=[]
            if(data.Tag != ""){
              data.TagArray=data.Tag.split(" ")
            }
            data.HostGroupArray=[]
            if(data.HostGroup != ""){
              data.HostGroupArray=data.HostGroup.split(" ")
            }
            this.form = data

              if(this.form.IsGroup){
                var shosts=this.form.Hosts.split("\n");
                for (var i in shosts){
                  var count=shosts[i].split("#");
                  if(count.length==2){
                    count[1]=count[1].replace(new RegExp("\n",'g'),"");
                    count[1]=count[1].replace(new RegExp("\r",'g'),"");
                    var index=(count[0]|0)-1
                    if(index>-1){
                      if(this.hosts[index]){
                        this.hosts[index]=this.hosts[index]+count[1]+"\r\n"
                      }else {
                        this.hosts[index]=count[1]+"\r\n"
                      }

                    }
                  }
                }
                for (var i in this.hosts){
                  this.hosts[i]=this.hosts[i].trim("\r\n")
                }
              }

            this.load_data = false
          })
          .catch(() => {
            this.load_data = false
          })
      },
      //时间选择改变时
      on_change_birthday(val) {
        this.$set(this.form, 'birthday', val)
      },
      //提交
      on_submit_form() {
        this.$refs.form.validate((valid) => {
          if (
            !valid
          )
            return false
          if(this.form.IsGroup){
            this.form.Hosts=""
            for (var i in this.hosts){
              var count= (i|0)+1
              var shosts=this.hosts[i].split("\n");
              for (var j=0 ;j<shosts.length;j++){
                shosts[j]=shosts[j].replace(new RegExp("\n",'g'),"");
                shosts[j]=shosts[j].replace(new RegExp("\r",'g'),"");
                if(shosts[j] && shosts[j]!=""){
                  this.form.Hosts=this.form.Hosts+count+"#"+shosts[j]+"\r\n"
                }
              }
            }
            this.form.Hosts=this.form.Hosts.trim("\r\n")
          }


          this.form.Tag = this.form.TagArray.join(" ")
          this.form.HostGroup = this.form.HostGroupArray.join(" ")
          var tagArray= this.form.TagArray
          delete this.form.TagArray
          this.on_submit_loading = true
          this.$http.post(port_conf.save, this.form)
            .then(({data: {msg}}) => {
              this.$message({
                message: msg,
                type: 'success'
              })
              setTimeout(() => {
                  this.$router.back()
                },
                500
              )
            })
            .catch(() => {
              this.on_submit_loading = false
            })
          this.form.TagArray=tagArray
        })
      }
    },
    components: {
      panelTitle
    }
  }
</script>
