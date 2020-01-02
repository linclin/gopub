
<template>
<div class="login-bodya">
    <div class="loginWarpa">
        <div class="login-titlea">
                <div> 个人信息设置</div>
            </div>

     <div class="login-forma"  v-if="userinfo.id">
                    <el-form ref="form" :model="form" :rules="rules" label-width="0">
                    <el-form-item prop="old_password" class="login-itema">
                    <label >用户名 ：</label><span v-text="userinfo.username"></span>
                    </el-form-item>

                      <el-form-item prop="old_password" class="login-itema">
                    <label >邮箱 ：</label> <span v-text="userinfo.email"></span>
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema">
                    <label >花名.实名 ：</label> <span v-text="userinfo.realname">  </span>    
                    </el-form-item>

                    <el-form-item prop="newpassword" class="login-itema">
                        <el-input type="password" v-model="form.newpassword" placeholder="请输入新密码："
                                  class="form-inputa"></el-input>
                    </el-form-item>
                     <el-form-item prop="repeat_newpassword" class="login-itema">
                        <el-input type="password" v-model="form.repeat_newpassword" placeholder="确认新密码："
                                  class="form-inputa"></el-input>
                    </el-form-item>

                    <el-form-item class="login-itema">
                        <el-button size="large"  class="form-submita" @click="submit_form">修改密码</el-button>
                    </el-form-item>
                </el-form>
            </div>
        </div>
</div>
</template>

<script type="text/javascript">
    import {port_user} from 'common/port_uri'
    import store from 'store'


    export default{
        data(){
            return {
                uid: this.$route.query.id,
                userinfo:{},
                form: {
                    uid:this.$route.query.id
                },
                rules: {
                    newpassword: [{required: true, message: '请输入新密码：', trigger: 'blur'}],
                   repeat_newpassword: [{required: true, message: '确认新密码：', trigger: 'blur'}]
                },
                //请求时的loading效果
                load_data: false
            }
        },
        created() {
            var uid = store.state.user_info.user.Id
            var role = store.state.user_info.user.Role
            if (this.$route.query.id && (uid == this.$route.query.id || role == 1)) {
                this.get_user_info()
            }else{
                this.$message({
                    message: "无权访问",
                    type: 'warning'
                })
            }
        },
        methods: {
        get_user_info(){
            this.$http.get(port_user.users, {
                  params: {
                    id:this.uid
                  }
              })
              .then(({data: {data}}) => {
                    this.userinfo=data
          })
        },
        //提交
        submit_form() {
             this.$http.post(port_user.changepasswd,this.form)
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
        margin:0 auto;
        
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
