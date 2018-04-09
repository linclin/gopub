<template>

    <div class="panel">
        <panel-title :title="$route.meta.title">
          <router-link style="float: left; margin-right: 10px" :to="{name: 'register'}" tag="span">
            <el-button type="primary" icon="plus" size="small">添加</el-button>
          </router-link>
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
                        label="名称"
                        width="100"
                >
                </el-table-column>
                   <el-table-column
                        prop="updated_at"
                        label="上线时间"
                        width="180">
                </el-table-column>
              <el-table-column
                prop="email"
                label="邮箱">
              </el-table-column>
              <el-table-column
                prop="username"
                label="用户名">
              </el-table-column>
                <el-table-column
                        prop="created_at"
                        label="创建时间"
                        width="180">
                </el-table-column>
                <el-table-column
                        prop="role"
                        label="	角色" width="100">
                  <template scope="props">
                    <span  >{{ props.row.role | getRole }}</span>
                  </template>
                </el-table-column>
                <el-table-column
                        label="操作"
                        width="300">
                    <template scope="props">
                        <router-link :to="{name: 'register', query:  {id: props.row.id}}"
                                     tag="span">
                            <el-button size="small" icon="edit">修改</el-button>
                        </router-link>
                        <el-button type="danger" size="small" icon="delete"
                                   @click="delete_data(props.row.id)">删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle, bottomToolBar, search} from 'components'
    import {port_user} from 'common/port_uri'
    import store from 'store'
    export default{
      filters: {
        getRole: function (value) {
          if (value == 1) {
            return "管理员"
          } else if (value == 10) {
            return "全部预发布用户"
          } else if (value == 20) {
            return "单个项目用户"
          }
          return value
        },
      },
        data(){
            return {
                table_data: null,
                load_data: true,
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
                this.load_data = true
                this.$http.get(port_user.users)
                        .then(({data: {data}}) => {
                    this.table_data = data
                    this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },

        }
    }
</script>
