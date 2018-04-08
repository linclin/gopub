<template>
  <div class="login-bodya">
    <div class="loginWarpa">
      <div class="login-titlea">
        <div> 用户注册</div>
      </div>
      <div class="login-forma">
        <el-form ref="form" :model="form" :rules="rules" label-width="0">
          <el-form-item prop="register_username" class="login-itema">
            <label class="labela">账户名 ：</label>
            <el-input v-model="form.register_username" placeholder="请输入账户名：" class="form-inputa"
                      :autofocus="true"></el-input>
          </el-form-item>

          <el-form-item prop="register_realname" class="login-itema">
            <label class="labela">花名.实名 ：</label>
            <el-input v-model="form.register_realname" placeholder="输入规范如：春哥.李宇春" class="form-inputa"
                      :autofocus="true"></el-input>
          </el-form-item>

          <el-form-item prop="register_email" class="login-itema">
            <label class="labela">邮箱 ：</label>
            <el-input v-model="form.register_email" placeholder="请输入联系邮箱：" class="form-inputa"
                      :autofocus="true"></el-input>
          </el-form-item>
          <el-form-item label="用户类型:" label-width="100px">
            <el-radio-group v-model="form.Role">
              <el-radio :label="1">管理员</el-radio>
              <el-radio :label="10">全部预发布用户</el-radio>
              <el-radio :label="20">单个项目用户</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="form.Role==20" label="项目名称:" label-width="100px">
            <el-select v-model="pro_id" filterable placeholder="请选择" multiple style="width: 400px;">
              <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item class="login-item">
            <el-button size="large" class="form-submita" @click="submit_forma">确认注册</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>
<script type="text/javascript">
  import {port_user, port_conf, port_code} from 'common/port_uri'
  import {mapActions} from 'vuex'
  import {SET_USER_INFO} from 'store/actions/type'

  export default {
    data() {
      return {
        projects: null,
        options: [],
        pro_id: [],
        form: {
          Role: 1,
          pro_ids: ""
        },
        rules: {
          register_username: [{required: true, message: '请输入账户名！', trigger: 'blur'}],
          register_realname: [{required: true, message: '请输入花名.实名！', trigger: 'blur'}],
          register_email: [{required: true, message: '请输入邮箱！', trigger: 'blur'}]

        },
        //请求时的loading效果
        load_data: false
      }
    },
    created() {
      this.get_project_data()
    },
    methods: {
      ...mapActions({
        set_user_info: SET_USER_INFO
      }),
      //提交


      get_project_data() {
        this.load_data = true
        this.$http.get(port_conf.list)
          .then(({data: {data}}) => {
            var opData = []
            for (var i in data.table_data) {
              var value = data.table_data[i].id
              var env = ""
              if (data.table_data[i].level == 1) {
                env = "测试环境"
                var lable = env + "-" + data.table_data[i].name
                opData.push({label: lable, value: value})
              }
              if (data.table_data[i].level == 2) {
                env = "预发布环境"
                var lable = env + "-" + data.table_data[i].name
                opData.push({label: lable, value: value})
              }
              if (data.table_data[i].level == 3) {
                env = "正式环境"
                var lable = env + "-" + data.table_data[i].name
                opData.push({label: lable, value: value})
              }


            }
            this.projects = data.table_data
            this.options = opData
          }).catch(() => {
          this.load_data = false
        })
      },
      submit_forma() {
        if (this.form.Role === 20) {
          this.form.pro_ids=this.pro_id.toString()
        }
        this.$http.post(port_user.register, this.form)
          .then(({data: {msg}}) => {
            this.$message({
              message: msg,
              type: 'success'
            })
            setTimeout(() => {
                this.$router.push({path: '/'})
              },
              500
            )
          })
      }
    }
  }
</script>
<style lang="scss" type="text/css" rel="stylesheet/scss">
  .login-bodya {
    position: relative;
    left: 0;
    top: 0;
    width: auto;
    height: auto;
    margin: 0 auto;

  .loginWarpa {
    width: 500px;
    padding: 25px 15px;
    margin: 0 auto;
    background-color: #fff;
    border-radius: 5px;

  .login-titlea {
    margin-bottom: 25px;
    text-align: center;
  }

  .login-itema {

  .el-input__inner {
    margin: 0 !important;
  }

  }
  .form-inputa {

  input {
    margin-bottom: 15px;
    font-size: 12px;
    height: 40px;
    border: 1px solid #eaeaec;
    background: #eaeaec;
    border-radius: 5px;
    color: #555;
  }

  }
  .form-submita {
    width: 100%;
    color: #fff;
    border-color: #6bc5a4;
    background: #6bc5a4;

  &
  :active,

  &
  :hover {
    border-color: #6bc5a4;
    background: #6bc5a4;
  }

  }
  }
  }
</style>
