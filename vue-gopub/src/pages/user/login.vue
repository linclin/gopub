<template>
    <div class="login-body">
        <div class="loginWarp"
             v-loading="load_data"
             element-loading-text="正在登陆中..."
             @keyup.enter="submit_form">
            <div class="login-title">
                <div> 自动发布系统</div>
            </div>
            <div class="login-form">
                <el-form ref="form" :model="form" :rules="rules" label-width="0">
                    <el-form-item prop="user_name" class="login-item">
                        <el-input v-model="form.user_name" placeholder="请输入账户名：" class="form-input"
                                  :autofocus="true"></el-input>
                    </el-form-item>
                    <el-form-item prop="user_password" class="login-item">
                        <el-input type="password" v-model="form.user_password" placeholder="请输入账户密码："
                                  class="form-input"></el-input>
                    </el-form-item>
                    <el-form-item class="login-item">
                        <el-button size="large"  class="form-submit" @click="submit_form">登录</el-button>
                    </el-form-item>
                    <el-form-item class="login-item">
                        <el-button size="large"  class="form-submit" @click="to_tasklist">上线单查询</el-button>
                    </el-form-item>
                </el-form>
            </div>
        </div>
    </div>
</template>
<script type="text/javascript">
    import {port_user} from 'common/port_uri'
    import {mapActions} from 'vuex'
    import {SET_USER_INFO} from 'store/actions/type'

    export default{
        data(){
            return {
                form: {
                    user_name: '',
                    user_password: ''
                },
                rules: {
                    user_name: [{required: true, message: '请输入账户名！', trigger: 'blur'}],
                    user_password: [{required: true, message: '请输入账户密码！', trigger: 'blur'}]
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
            this.$refs.form.validate((valid) => {
                if (valid) {
                    this.load_data = true
                    //登录提交
                    this.$http.post(port_user.login, this.form)
                            .then(({data: {data, msg}}) => {
                        let isNull = data !== null
                        this.set_user_info({
                        user: data,
                        login: true
                    })
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
                .
                    catch(() => {
                        this.load_data = false
                })
                } else {
                    return false
                }
            }
        )
    },
    to_tasklist(){
          this.$router.push({path: '/task/searchlist'})
    }
    }
    }
</script>
<style lang="scss" type="text/css" rel="stylesheet/scss">
    .login-body {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        background-image: url(./images/login_bg.jpg);
        background-repeat: no-repeat;
        background-position: center;
        background-size: cover;

        .loginWarp {
            width: 300px;
            padding: 25px 15px;
            margin: 100px auto;
            background-color: #fff;
            border-radius: 5px;

            .login-title {
                margin-bottom: 25px;
                text-align: center;
            }

            .login-item {

                .el-input__inner {
                    margin: 0 !important;
                }

            }
            .form-input {

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
            .form-submit {
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
