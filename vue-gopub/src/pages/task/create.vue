<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body"
         v-loading="load_data"
         element-loading-text="拼命加载中">
      <el-row>
        <el-col :span="20">
          <el-form label-width="100px">
            <div class="panel-title">
              预发布
              <div class="fr">

              </div>
            </div>
            <div class="panel-body">
              <el-form-item label="项目名称:" label-width="100px">
                <el-select v-model="pro_id" filterable placeholder="请选择" style="width: 400px;">
                  <el-option
                    v-for="item in options"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value">
                  </el-option>
                </el-select>
              </el-form-item>
            </div>
            <div class="panel-title">
              正式环境
              <div class="fr">

              </div>
            </div>
            <div class="panel-body">
              <el-form-item label="项目名称:" label-width="100px">
                <el-select v-model="pro_id1" filterable placeholder="请选择" style="width: 400px;">
                  <el-option
                    v-for="item in options1"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value">
                  </el-option>
                </el-select>
              </el-form-item>
            </div>
            <el-form-item>
              <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading">创建
              </el-button>
              <el-button @click="$router.back()">返回</el-button>
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
        projects: null,
        options: [],
        options1: [],
        pro_id: null,
        pro_id1: null,
        load_data: false,
        on_submit_loading: false
      }
    },
    created() {
      this.get_project_data()
    },
    methods: {
      //获取数据
      get_project_data() {
        this.load_data = true
        this.$http.get(port_conf.mylist)
          .then(({data: {data}}) => {
            var opData = []
            var opData1 = []
            var opData2 = []
            for (var i in data.table_data) {
              var value = data.table_data[i].id
              var env = ""
              if (data.table_data[i].level == 1) {
                env = "测试环境"
                var lable = env + "-" + data.table_data[i].name
                opData2.push({label: lable, value: value})
              }
              if (data.table_data[i].level == 2) {
                env = "预发布环境"
                var lable = env + "-" + data.table_data[i].name
                opData.push({label: lable, value: value})
              }
              if (data.table_data[i].level == 3) {
                env = "正式环境"
                var lable = env + "-" + data.table_data[i].name
                opData1.push({label: lable, value: value})
              }


            }
            this.projects = data.table_data
            this.options = opData
            this.options1 = opData1
            this.load_data = false
            if (this.$route.query.id) {
              this.pro_id=this.$route.query.id
              this.on_submit_form()
            }
          }).catch(() => {
          this.load_data = false
        })
      },
      //提交
      on_submit_form() {
        var proId = 0
        if (this.pro_id) {
          proId = this.pro_id
        } else {
          proId = this.pro_id1
        }
        if (proId) {
          for (var i in  this.projects) {
            var project = this.projects[i]
            if (proId == project.id) {
              if (project.repo_type == "git") {
                this.$router.push({
                  name: 'taskGit',
                  query: {id: proId}
                })
                return
              }
              if (project.repo_type == "file") {
                this.$router.push({
                  name: 'taskFile',
                  query: {id: proId}
                })
                return
              }
            }
          }

        }
      }
    },
    components: {
      panelTitle
    }
  }
</script>
