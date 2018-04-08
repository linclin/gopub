<template>
    <div class="left-side">
        <div class="left-side-inner">
            <router-link to="/" class="logo block">
                <div style="color: white"> 卷皮自动发布系统</div>

            </router-link>
            <el-menu
                    class="menu-box"
                    theme="dark"
                    router
                    :default-active="$route.path">
                <div
                        v-for="(item, index) in nav_menu_data"
                        :key="index">
                    <el-menu-item
                            class="menu-list"
                            v-if="typeof item.child === 'undefined'"
                            :index="item.path">
                        <i class="icon fa" :class="item.icon"></i>
                        <span v-text="item.title" class="text"></span>
                    </el-menu-item>
                    <el-submenu
                            :index="item.path"
                            v-else>
                        <template slot="title">
                            <i class="icon fa" :class="item.icon"></i>
                            <span v-text="item.title" class="text"></span>
                        </template>
                        <el-menu-item
                                class="menu-list"
                                v-for="(sub_item, sub_index) in item.child"
                                :index="sub_item.path"
                                :key="sub_index">
                            <!--<i class="icon fa" :class="sub_item.icon"></i>-->
                            <span v-text="sub_item.title" class="text"></span>
                        </el-menu-item>
                    </el-submenu>
                </div>
            </el-menu>
        </div>
    </div>
</template>
<script type="text/javascript">
  import store from 'store'
    export default{
        data(){
          var Role = store.state.user_info.user.Role
          if(Role==1){
            return {
              nav_menu_data: [{
                title: "主页",
                path: "/home",
                icon: "fa-home"
              }, {
                title: "项目配置",
                path: "/conf/list",
                icon: "el-icon-menu"
              }, {
                title: "上线单",
                path: "/task",
                icon: "fa-table",
                child: [{
                  title: "全部上线单",
                  path: "/task/list"
                }, {
                  title: "我的上线单",
                  path: "/task/mylist"
                }, {
                  title: "创建上线单",
                  path: "/task/create"
                }]
              }, {
                title: "agent状态查询",
                path: "/p2p/check",
                icon: "ace-icon fa fa-desktop"
              }, {
                title: "预发布统计",
                path: "/other/noauto",
                icon: "el-icon-date"
              },{
                title: "其他操作",
                path: "/other",
                icon: "fa-bar-chart-o",
                child: [{
                  title: "刷新版本号",
                  path: "/other/flush"
                },{
                  title: "预发布git版本查看",
                  path: "/other/gitpull"
                }]
              }]
            }
          }else{
            return {
              nav_menu_data: [{
                title: "主页",
                path: "/home",
                icon: "fa-home"
              }, {
                title: "上线单",
                path: "/task",
                icon: "fa-table",
                child: [ {
                  title: "我的上线单",
                  path: "/task/mylist"
                }, {
                  title: "创建上线单",
                  path: "/task/create"
                }]
              },{
                title: "其他操作",
                path: "/other",
                icon: "fa-bar-chart-o",
                child: [{
                  title: "预发布git版本查看",
                  path: "/other/gitpull"
                }]
              }]
            }
          }

        }
    }
</script>
