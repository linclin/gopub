
<template>
<div class="login-bodya">
    <div class="loginWarpa">
        <div class="login-titlea">
                <div> 个人信息设置</div>
            </div>

     <div class="login-forma"  v-if="get_user_info.login">
                    <el-form ref="form" :model="form" :rules="rules" label-width="0">
                    <el-form-item prop="old_password" class="login-itema">
                    <label >用户名 ：</label><span v-text="get_user_info.user.Username"></span>
                    </el-form-item>

                      <el-form-item prop="old_password" class="login-itema">
                    <label >邮箱 ：</label> <span v-text="get_user_info.user.Email"></span>
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema">
                    <label >花名.实名 ：</label> <span v-text="get_user_info.user.Realname">  </span>    
                    </el-form-item>

                    <el-form-item prop="old_password" class="login-itema">
                        <el-input type="password" v-model="form.old_password" placeholder="请输入旧密码："
                                  class="form-inputa"></el-input>
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
    import {mapGetters, mapActions} from 'vuex'
    import {SET_USER_INFO} from 'store/actions/type'
    import {GET_USER_INFO} from 'store/getters/type'


    export default{
          computed: {
                ...mapGetters({
                    get_user_info: GET_USER_INFO
                })
    }
    ,
        data(){
            return {
                form: {
                
                },
                rules: {
                    old_password: [{required: true, message: '请输入旧密码：', trigger: 'blur'}],
                    newpassword: [{required: true, message: '请输入新密码：', trigger: 'blur'}],
                   repeat_newpassword: [{required: true, message: '确认新密码：', trigger: 'blur'}]
                },
                //请求时的loading效果
                load_data: false
            }
        },
        
        methods: {
                ...mapActions({
                    set_user_info: SET_USER_INFO
                }),
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
