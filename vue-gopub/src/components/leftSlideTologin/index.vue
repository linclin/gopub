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
                </div>
            </el-menu>
        </div>
    </div>
</template>
<script type="text/javascript">
    export default{
        data(){
            return {
                nav_menu_data: [{
                    title: "退出",
                    path: "/user/login",
                    icon: "fa-arrow-left"
                }]
            }
        }
    }
</script>
