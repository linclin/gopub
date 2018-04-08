<template>

    <div class="panel">
        <panel-title :title="$route.meta.title">
            <div style="float: left;margin-right: 10px">
                <search @search="submit_search"></search>
            </div>
            <el-button @click.stop="on_refresh" size="small">
                <i class="fa fa-refresh"></i>
            </el-button>
        </panel-title>

        <div class="panel-body">
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
                        prop="realname"
                        label="创建人"
                        width="100"
                >
                </el-table-column>
                <el-table-column
                        prop="name"
                        label="项目"
                >
                </el-table-column>
                <el-table-column
                        prop="title"
                        label="上线单标题">
                </el-table-column>

                <el-table-column
                        prop="updated_at"
                        label="上线时间"
                        width="180">
                </el-table-column>
                <el-table-column
                        prop="branch"
                        label="分支">
                </el-table-column>
                <el-table-column
                        prop="commit_id"
                        label="	上线commit号">
                </el-table-column>
                <el-table-column
                        prop="pms_batch_id"
                        label="	发布批次ID" width="60">
                </el-table-column>
                <el-table-column
                        prop="pms_uwork_id"
                        label="	uwork流程号" width="80">
                </el-table-column>
                <el-table-column
                        prop="status"
                        label="	当前状态" width="100">
                </el-table-column>
                <el-table-column
                        label="操作"
                        width="300">
                    <template scope="props">
                        <router-link :to="{name: 'taskRelease', params: {id: props.row.id,is_log:1}}"
                                     v-if="props.row.status!='新建提交'" tag="span">
                            <el-button size="small" icon="edit">查看日志</el-button>
                        </router-link>
                        <router-link :to="{name: 'taskRelease', params: {id: props.row.id}}"
                                     v-if="props.row.status!='上线中'&&props.row.status!='上线完成'" tag="span">
                            <el-button type="info" size="small" icon="edit">上线</el-button>
                        </router-link>
                        <el-button type="warning" size="small" icon="share" @click="create_rollback(props.row.id)"
                                   v-if="props.row.status=='上线完成'&&props.row.action=='0'&&props.row.enable_rollback=='1' ">
                            回滚
                        </el-button>
                      <el-button type="warning" size="small" icon="share" @click="create_rollback_this(props.row.id)"
                                 v-if="props.row.status=='上线完成'&&props.row.action=='0'&&props.row.enable_rollback=='1' ">
                        回滚当前
                      </el-button>
                        <el-button type="danger" size="small" icon="delete" v-if="props.row.status=='新建提交'|| props.row.status=='上线失败'"
                                   @click="delete_data(props.row.id)">删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <bottom-tool-bar>

                <div slot="page">
                    <el-pagination
                            @current-change="handleCurrentChange"
                            :current-page="currentPage"
                            :page-size="15"
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
    import {port_task} from 'common/port_uri'
    import store from 'store'
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
                select_info: ""
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
            //刷新
            on_refresh(){
                this.select_info = ""
                this.get_table_data()
            },
            //获取数据
            get_table_data(){
                console.log(store.state.user_info.user.Id)
                this.load_data = true
                this.$http.get(port_task.list, {
                            params: {
                                page: this.currentPage,
                                length: this.length,
                                select_info: this.select_info,
                                my: store.state.user_info.user.Id
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
                this.$http.get(port_task.del, {
                            params: {
                                taskId: id,
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
            },
          create_rollback_this(id){
            this.$confirm('此操作将回滚项目到当前版本, 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
            }).then(() => {
              this.load_data = true
              this.$http.get(port_task.rollback, {
                params: {
                  taskId: id,
                  this: "this",
                }
              })
                .then(({data: {data}}) => {
                  this.$message({
                    message: "success",
                    type: 'success'
                  })
                  setTimeout(() => {
                      this.$router.push({
                        name: 'taskRelease',
                        params: {id: data.Id}

                      })
                    },
                    500
                  )
                }).
              catch(() => {
                this.load_data = false
              })
            }).
            catch(() => {
              this.load_data = false
            })
          },
            //创建回滚
            create_rollback(id){
                this.$confirm('此操作将回滚任务, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.load_data = true
                this.$http.get(port_task.rollback, {
                            params: {
                                taskId: id,
                            }
                        })
                        .then(({data: {data}}) => {
                    this.$message({
                    message: "success",
                    type: 'success'
                })
                setTimeout(() => {
                    this.$router.push({
                    name: 'taskRelease',
                    params: {id: data.Id}

                })
            },
                500
            )
            }).
                catch(() => {
                    this.load_data = false
            })
            }).
                catch(() => {
                    this.load_data = false
            })
            }
        }
    }
</script>
