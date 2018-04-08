<template>
    <div >
        <div class="panel-body"  style="clear: both;">
          <el-row>
            <el-col :span="8">
             <div class="panel">
              <panel-title :title="msg1"> </panel-title>
               <div class="panel-body">
                <div  >
                <el-table
                    :data="day_data"
                    v-loading="load_data"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;"
                    max-height="750">
                <el-table-column
                        prop="id"
                        label="项目Id"
                        width="80">
                </el-table-column>
                       <el-table-column
                        prop="name"
                        label="项目名称"
                       >
                </el-table-column>
             </el-table>
                </div>
               </div>
        </div>
        </el-col>
        <el-col :span="8">
            <div class="panel" style="margin-left: 10px">
            <panel-title :title="msg2"><div style="float: left;margin-right: 10px;">
               </div></panel-title>
            <div class="panel-body">
                <div >
                <el-table
                    :data="week_data"
                    v-loading="load_data"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;"
                     max-height="750">
                <el-table-column
                        prop="id"
                        label="项目Id"
                        width="80">
                </el-table-column>
                       <el-table-column
                        prop="name"
                        label="项目名称"
                      >
                </el-table-column>
             </el-table>
                </div>
            </div>
            </div>
        </el-col>
          <el-col :span="8">
            <div class="panel" style="margin-left: 10px">
            <panel-title :title="msg3"><div style="float: left;margin-right: 10px;">
               </div></panel-title>
            <div class="panel-body">
                <div >
                <el-table
                    :data="month_data"
                    v-loading="load_data"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;"
                     max-height="750">
                <el-table-column
                        prop="id"
                        label="项目Id"
                        width="80">
                </el-table-column>
                       <el-table-column
                        prop="name"
                        label="项目名称"
                      >
                </el-table-column>
             </el-table>
                </div>
            </div>
            </div>
        </el-col>
    </el-row>

        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle, bottomToolBar} from 'components'
    import {port_other} from 'common/port_uri'
    export default{
        data(){
            return {
                day_data: null,
                week_data: null,
                month_data: null,
                //数据总条目
                total: 0,
                //每页显示多少条数据
                length: 15,
                //请求时的loading效果
                load_data: false,
                load_data1: false,
                select_info: "",
                //项目详情
                project_data:[],
                msg1:'本日共有个项目未自动预发布',
                msg2:'本周共有个项目未自动预发布',
                msg3:'本月共有个项目未自动预发布',
            }
        },
        components: {
            panelTitle,
            bottomToolBar,
        },
        created(){
            this.get_day_data()
            this.get_week_data()
            this.get_month_data()
        },
        methods: {  
            get_day_data(){
                this.load_data = true
                this.$http.get(port_other.noauto, {
                            params: {
                                taskType:"day",               
                            }
                        })
                        .then(({data: {data}}) => {
                this.day_data=data
                this.msg1= "本日共有"+data.length+"个项目未自动预发布"
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
        },
            get_week_data(){
                this.load_data = true
                this.$http.get(port_other.noauto, {
                            params: {
                                taskType:"week",               
                            }
                        })
                        .then(({data: {data}}) => {
                this.week_data=data
                this.msg2= "本周共有"+data.length+"个项目未自动预发布"
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
        },
         get_month_data(){
                this.load_data = true
                this.$http.get(port_other.noauto, {
                            params: {
                                taskType:"month",               
                            }
                        })
                        .then(({data: {data}}) => {
                this.month_data=data
                this.msg3= "本月共有"+data.length+"个项目未自动预发布"
                console.log(this.month_data)
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
        }
        }
    }
</script>
