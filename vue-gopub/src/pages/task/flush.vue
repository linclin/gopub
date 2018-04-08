<template>
    <div class="panel">
        <panel-title :title="$route.meta.title"></panel-title>
        <div class="panel-body"
             v-loading="load_data"
             element-loading-text="拼命加载中">
            <el-row>
                <el-col :span="20">
                    <el-form label-width="100px">
                        <el-form-item label="项目名称:" label-width="100px">
                            <el-select v-model="pro_id" multiple  filterable placeholder="请选择" style="width: 400px;">
                                <el-option
                                        v-for="item in options"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading">刷新
                            </el-button>
                            <el-button @click="$router.back()">返回</el-button>
                        </el-form-item>
                    </el-form>
                    <terminal  v-if="pro_id.length>0" :taskId="-2"></terminal>
                </el-col>
            </el-row>
        </div>
    </div>
</template>
<script type="text/javascript">
    import {panelTitle,terminal} from 'components'
    import {port_conf,port_task, port_code} from 'common/port_uri'
    import {tools_verify} from 'common/tools'

    export default{
        data(){
            return {
                projects: null,
                options: [],
                pro_id: [],
                load_data: false,
                on_submit_loading: false
            }
        },
        created(){
            this.get_project_data()
        },
        methods: {
            //获取数据
            get_project_data(){
                this.load_data = true
                this.$http.get(port_conf.list)
                        .then(({data: {data}}) => {
                    var opData = []
                    for(var i in data.table_data){
                        var value = data.table_data[i].id
                        var env = ""
                        if (data.table_data[i].level == 1) {
                            env = "测试环境"
                        }
                        if (data.table_data[i].level == 2) {
                            env = "预发布环境"
                        }
                        if (data.table_data[i].level == 3) {
                            env = "正式环境"
                        }
                        var lable = env + "-" + data.table_data[i].name
                        opData.push({label: lable, value: value})

                }
                this.projects = data.table_data
                this.options = opData
                this.load_data = false
            })
            .
                catch(() => {
                    this.load_data = false
            })
            },
            //提交
            on_submit_form(){
                this.on_submit_loading = true
                if (this.pro_id.length>0) {
                    this.$http.get(port_task.flush, {
                                params: {
                                    projectIds: this.pro_id.join(","),
                                }
                            })
                            .then(({data: {data}}) => {
                    console.log(data)
                    this.load_data = false
                    this.on_submit_loading = false
                })
                .
                    catch(() => {
                        this.load_data = false
                    this.on_submit_loading = false
                })

                }
            }
        },
        components: {
            panelTitle,
            terminal
        }
    }
</script>
