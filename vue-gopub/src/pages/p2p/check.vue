<template>
    <div >
        <div class="panel-body"  style="clear: both;">
          <el-row>
            <el-col :span="12">
             <div class="panel">
              <panel-title :title="msg1"> </panel-title>
               <div class="panel-body">
                <div >

                <el-table
                    :data="table_data"
                    v-loading="load_data"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;"
                     max-height="750">
                <el-table-column
                        prop="Host"
                        label="Agent's IP"
                        width="450">
                </el-table-column>
                       <el-table-column
                        prop="Status"
                        label="状态"
                        width="200">
                </el-table-column>

                       <el-table-column
                        prop="Pid"
                        label="项目Id"
                       >
                </el-table-column>
                  <el-table-column
                    prop="Pname"
                    label="项目名称"
                  >
                  </el-table-column>
             </el-table>

                </div>
               </div>
        </div>

        </el-col>
        <el-col :span="12">
            <div class="panel" style="margin-left: 10px">
            <panel-title :title="msg2"><div style="float: left;margin-right: 10px;">
                <search @search="submit_search"></search>
               </div></panel-title>
            <div class="panel-body">
                <div >

                <el-table
                    :data="project_data"
                    v-loading="load_data1"
                    element-loading-text="拼命加载中"
                    border
                    style="width: 100%;"
                     max-height="750"
                    >
                <el-table-column
                        prop="ip"
                        label="Agent's IP"
                        width="390">
                </el-table-column>
                <el-table-column
                        prop="status"
                        label="状态"
                        width="200">
                </el-table-column>
                 <el-table-column
                        label="操作"
                     >

                <template scope="props">
                        <el-button type="warning" size="small" icon="document" @click="send_data(props.row.id)">发送
                        </el-button>
                        <el-button style="margin-left:0px"type="danger" size="small" icon="delete" @click="restart_agent(props.row.id)">重启
                        </el-button>
                </template>
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
    import {panelTitle, bottomToolBar, search} from 'components'
    import {port_p2p} from 'common/port_uri'
    export default{
        data(){
            return {
                table_data: null,
                //数据总条目
                total: 0,
                //请求时的loading效果
                load_data: false,
                load_data1: false,
                //批量选择数组
                batch_select: [],
                //批量选择数组
                select_info: "",
                //项目详情
                project_data:[],
                msg1:'Agent',
                msg2:'按项目ID查询',
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
                this.project_data=[]
                this.select_info = value
                  this.load_data1 = true
                    this.$http.get(port_p2p.check, {
                            params: {
                                type:"1",
                                projectId: this.select_info,
                            }
                        })
                        .then(({data: {data}}) => {
                            for(var key in data){
                                this.project_data.push({"ip":key,"status":data[key]})
                            }
                this.load_data1 = false
            })
            .
                catch(() => {
                    this.load_data1 = false
            })
            },
            //send
            send_data() {
                  this.load_data1 = true
                    this.$http.get(port_p2p.send, {
                            params: {
                                projectId: this.select_info,
                            }
                        })
                        .then(({data: {msg}}) => {
                            console.log(msg)
                          this.$message({
                          message: msg,
                          type: 'success'
                            })
                      this.submit_search( this.select_info)
                      this.load_data1 = false
            })
            .
                catch(() => {
                    this.load_data1 = false
            })
            },
            //agent
             restart_agent() {
                  this.load_data1 = true
                    this.$http.get(port_p2p.agent, {
                            params: {
                                projectId: this.select_info,
                            }
                        })
                        .then(({data: {msg}}) => {
                           console.log(msg)
                           this.$message({
                           message: msg,
                           type: 'success'
                            })
                this. submit_search(this.select_info)
                this.load_data1 = false
            })
            .
                catch(() => {
                    this.load_data1 = false
            })
            },
            //获取数据
            get_table_data(){
                this.load_data = true
                this.$http.get(port_p2p.check, {
                            params: {
                                type:"0",
                            }
                        })
                        .then(({data: {data}}) => {
                this.table_data=data
                console.log(this.table_data)
                this.load_data = false
            })
            .
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
