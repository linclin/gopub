
<template>

    <div class="panel">
        <panel-title :title="$route.meta.title">
            <div style="float: left;margin-right: 10px;margin-top: 5px">
                <search @search="submit_search"></search>
            </div>
            <el-button @click.stop="on_refresh" size="small">
                <i class="fa fa-refresh"></i>
            </el-button>
            <router-link :to="{name: 'confAdd'}" tag="span">
                <el-button type="primary" icon="plus" size="small">创建项目</el-button>
            </router-link>

        </panel-title>

        <div class="panel-body" style="clear: both;">
            <el-table
                    :data="table_data"
                    v-loading="load_data"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;">
                <el-table-column
                        prop="id"
                        label="id"
                        width="80">
                </el-table-column>
                <el-table-column
                        prop="name"
                        label="项目名称">
                </el-table-column>
                <el-table-column
                        prop="realname"
                        label="创建人"
                        width="120">
                </el-table-column>
                <el-table-column
                        prop="level"
                        label="环境"
                        width="150">
                    <template scope="props">
                        <span v-text="props.row.level == 3 ? '线上环境' : '预发布环境'"></span>
                    </template>
                </el-table-column>
                <el-table-column
                        prop="repo_type"
                        label="发布方式"
                        width="100">
                </el-table-column>
                <el-table-column
                        prop="keep_version_num"
                        label="保留版本数量"
                        width="80">
                </el-table-column>
                <el-table-column
                        prop="p2p"
                        label="是否使用p2p"
                        width="80">
                    <template scope="props">
                        <span v-text="props.row.p2p == 0 ? '否' : '是'"></span>
                    </template>
                </el-table-column>
                <el-table-column
                        prop="updated_at"
                        label="更新时间"
                        width="200">
                </el-table-column>
                <el-table-column
                        prop="created_at"
                        label="创建时间"
                        width="200">
                </el-table-column>
                <el-table-column
                        label="操作"
                         width="370">
                <template scope="props">

                <el-popover ref="popover4" placement="left-start" width="620" trigger="click" >
                 <div style="margin-left:10px;font-size:15px">
              <i class="el-icon-date" style="color:#20A0FF"></i>  <span style="margin-left:8px;color:#20A0FF">项目详情</span>
                </div>
                  <div class="login-bodya">
                   <div class="loginWarpa">
                     <div class="login-forma">
                    <el-form ref="form"  label-width="0">
                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >项目名称 ：</label><span style="color:teal">{{project_data.Name}}</span>
                    </el-form-item>
                      <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >项目环境 ：</label> <span style="color:teal">{{project_data.Level == 3 ? '线上环境' : '预发布环境'}}</span>
                    </el-form-item>
                     <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                      <label >项目地址 ：</label> <span style="color:teal"> {{project_data.RepoUrl}} </span>
                      </el-form-item>
                       </el-form>
                    </div>
                  </div>
                </div>

                 <div style="margin-left:10px;font-size:15px">
                <i class="ace-icon fa fa-desktop" style="color:#F7BA21"></i> <span style="margin-top:-20px;margin-left:8px;color:#F7BA2A">宿主机</span>
                </div>
                  <div class="login-bodya">
                      <div class="loginWarpa">
                       <div class="login-forma">
                    <el-form ref="form"  label-width="0">
                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >代码检出仓库 ：</label><span style="color:teal">{{project_data.DeployFrom}}</span>
                    </el-form-item>
                      <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >排除文件 ：</label> <span style="color:teal">{{project_data.Excludes}}</span>
                    </el-form-item>
                 </el-form>
                </div>
              </div>
             </div>

                <div style="margin-left:10px;font-size:15px">
                <i class="el-icon-menu" style="color:#13CE66"></i> <span style="margin-top:-20px;margin-left:8px;color:#13CE66">机器列表</span>
                </div>
                  <div class="login-bodya">
                      <div class="loginWarpa">
                       <div class="login-forma">
                    <el-form ref="form"  label-width="0">
                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >用户 ：</label><span style="color:teal">{{project_data.ReleaseUser}}</span>
                    <label style="margin-left:120px">保留版本数 ：</label><span style="color:teal">{{project_data.KeepVersionNum}}</span>
                    </el-form-item>

                      <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >webroot ：</label> <span style="color:teal">{{project_data.ReleaseTo}}</span>
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >发布版本库 ：</label> <span style="color:teal"> {{project_data.ReleaseLibrary}} </span>
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >机器列表 ：</label> <span style="color:teal"> {{project_data.Hosts}} </span>
                    </el-form-item>
                 </el-form>
                </div>
              </div>
             </div>

                 <div style="margin-left:10px;font-size:15px">
                 <i class="el-icon-setting" style="color:#FF4949"></i><span style="margin-top:-20px;margin-left:8px;color:#FF4949">高级任务</span>
                </div>
                  <div class="login-bodya">
                      <div class="loginWarpa">
                       <div class="login-forma">
                    <el-form ref="form"  label-width="0">
                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >代码检出前任务 ：</label><span style="color:teal">{{project_data.PreDeploy}}</span>
                    </el-form-item>

                      <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >代码检出后任务 ：</label> <span style="color:teal">{{project_data.PostDeploy}}</span>
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >同步完目标机后任务 ：</label> <span style="color:teal"> {{project_data.PreRelease}} </span>
                    </el-form-item>

                     <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                    <label >更改版本软链接后任务  ：</label> <span style="color:teal"> {{project_data.PostRelease}} </span>
                    </el-form-item>

                     <el-form-item prop="old_password" class="login-itema" style="margin-top:-20px;margin-left:-25px">
                     <el-form-item><label >部属方式  ：</label> <span style="color:teal"> {{project_data.ReleaseType== 0 ? '软链接' : '移动目录'}} </span></el-form-item>
                     <el-form-item><label >是否开启p2p  ：</label> <span style="color:teal"> {{project_data.P2p== 0 ? 'No' : 'Yes'}} </span></el-form-item>
                     <el-form-item><label>是否开启gzip  ：</label> <span style="color:teal"> {{project_data.Gzip == 0 ? 'No' : 'Yes'}} </span></el-form-item>
                    </el-form-item>
                 </el-form>
                </div>
              </div>
             </div>


                </el-popover>
                        <el-button type="info" v-popover:popover4 size="small" icon="search" @click="open(props.row.id)">查看
                        </el-button>
                        <router-link :to="{name: 'confDetection', params: {id: props.row.id}}" tag="span">
                            <el-button type="success" size="small" icon="setting">检测</el-button>
                        </router-link>
                        <router-link :to="{name: 'confUpdate', params: {id: props.row.id}}" tag="span">
                            <el-button type="info" size="small" icon="edit">修改</el-button>
                        </router-link>
                        <el-button type="warning" size="small" icon="document" @click="copy_data(props.row.id)">复制
                        </el-button>
                        <el-button style="margin-left:0px"type="danger" size="small" icon="delete" @click="delete_data(props.row.id)">删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <bottom-tool-bar>

                <div slot="page">
                    <el-pagination
                            @current-change="handleCurrentChange"
                            :current-page="currentPage"
                            :page-size="10"
                            layout="total, prev, pager, next"
                            :total="total">
                    </el-pagination>
                </div>
            </bottom-tool-bar>
        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle, bottomToolBar, search} from 'components'
    import {port_conf} from 'common/port_uri'
    export default{
        data(){
            return {
                table_data: null,
                //当前页码
                currentPage: 1,
                //数据总条目
                total: 0,
                //每页显示多少条数据
                length: 15,
                //请求时的loading效果
                load_data: true,
                //批量选择数组
                batch_select: [],
                //批量选择数组
                select_info: "",
                //项目详情
                project_data:[]
            }
        },
        components: {
            panelTitle,
            bottomToolBar,
            search

        },
        created(){
            this.get_table_data()
        },
        methods: {
            submit_search(value) {
                this.select_info = value
                this.$message({
                    message: value,
                    type: 'success'
                })
                this.get_table_data()
            },
            open(value){
                this.$http.get(port_conf.get, {
                     params: {
                             projectId:value
                            }
                        })
                   .then(({data: {data}}) => {
                         console.log(data)
                  this.project_data=data

            })
             .
                catch(() => {})
            },
            //刷新
            on_refresh(){
                this.select_info = ""
                this.get_table_data()
            },
            //获取数据
            get_table_data(){
                this.load_data = true
                this.$http.get(port_conf.list, {
                            params: {
                                page: this.currentPage,
                                length: this.length,
                                select_info: this.select_info
                            }
                        })
                        .then(({data: {data}}) => {
                this.table_data = data.table_data
                this.currentPage = data.currentPage
                this.total = data.total
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            //根据id删除数据
            delete_data(id){
                this.$confirm('此操作将删除该数据, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.load_data = true
                this.$http.get(port_conf.del, {
                            params: {
                                projectId: id,
                            }
                        })
                        .then(({data: {msg}}) => {
                    this.get_table_data()
                this.$message({
                    message: msg,
                    type: 'success'
                })
            }).
                catch(() => {
                    this.load_data = false
            })
            }).
                catch(() => {
                    this.load_data = false
            })
            },
            //复制项目
            copy_data(id){
                this.$confirm('是否复制项目?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.load_data = true
                this.$http.get(port_conf.copy, {
                            params: {
                                projectId: id,
                            }
                        })
                        .then(({data: {msg}}) => {
                    this.get_table_data()
                this.$message({
                    message: msg,
                    type: 'success'
                })
            }).
                catch(() => {
                    this.load_data = false
            })
            }).
                catch(() => {
                    this.load_data = false
            })
            },
            //页码选择
            handleCurrentChange(val) {
                this.currentPage = val
                this.get_table_data()
            }
        }
    }
</script>
