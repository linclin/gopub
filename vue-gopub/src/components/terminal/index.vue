<template>
    <div>
        <el-row v-show="isShowStep">
            <el-col :span="24">
                <el-steps :active="action">
                    <el-step title="步骤 1" description="权限、目录检查"></el-step>
                    <el-step title="步骤 2" description="pre-deploy任务"></el-step>
                    <el-step title="步骤 3" description="代码检出"></el-step>
                    <el-step title="步骤 4" description="post-deploy任务"></el-step>
                    <el-step title="步骤 5" description="同步至服务器"></el-step>
                    <el-step title="步骤 6" description="全量更新(pre-release、更新版本、post-release)"></el-step>
                </el-steps>
            </el-col>
        </el-row>
        <div   style="margin: 5px 5px 0px;
                      padding: 3px;
                      border: 1px dashed rgb(0, 160, 198);
                      background-color: rgb(0,0,0);">
            <code style="background-color: rgb(0, 0, 0);color:#00ff00">
                <br>
                <span v-for="n in showText" :style="{'color': n.color}"> <pre style=" white-space: pre-wrap;" v-html="n.text"></pre> <br></span>
                <br>
            </code>
        </div>
    </div>

</template>
<script type="text/javascript">
    import {port_record} from 'common/port_uri'
    export default{
        props: ['taskId','isJson'],
        data(){
            return {
                showText: [],
                taskid: this.taskId * 1,
                action: 0,
                time: (Date.parse(new Date()) / 1000)-10
            }
        },
        computed: {
            test1: function () {
                this.get_data()
            },
            isShowStep:function () {
                if(this.taskId * 1>0){
                    return true
                } else{
                    return false
                }
            }
        },
        created () {
            this.get_data()
            var _this = this;
            this.intervalid1 = setInterval(function () {
                _this.get_data()
            }, 2000)
        },
        beforeDestroy () {
            this.time=(Date.parse(new Date()) / 1000);
            clearInterval(this.intervalid1)
        },
        methods: {
            get_data() {
                this.$http.get(port_record.list, {
                    params: {
                        taskId: this.taskId,
                        time: this.time
                    }
                }).then(({data: {data}}) => {
                    this.showText = [];
                var action = 0
                for (var i = 0; i < data.length; i++) {
                    //var text=data[i].command+"<br>"+data[i].memo
                    var color = "#00ff00"
                    if (data[i].status == "0") {
                        color = "red"
                    }
                    this.showText.push({text: data[i].command, "color": color})
                    var text=data[i].memo;
                    try{
                        var text= JSON.parse(data[i].memo);
                    }catch (e){
                    }
                    if(typeof text == "string"){
                        this.showText.push({text: data[i].memo, "color": color})
                    }else if(Object.prototype.toString.call(text)=='[object Array]') {
                        for(var j=0;j<text.length;j++){
                            try{
                                this.showText.push({text: "IP:"+text[j].Host, "color": color})
                                if(text[j].ErrorInfo) {
                                    this.showText.push({text: "错误结果:\n" + text[j].Result, "color": color})
                                }else{
                                    this.showText.push({text: "执行结果:\n" + text[j].Result, "color": color})
                                }
                                if(text[j].ErrorInfo){
                                    this.showText.push({text: "错误:"+text[j].ErrorInfo, "color": color})
                                }
                                this.showText.push({text: "=============", "color": color})
                            }catch (e){
                            }
                        }
                        //this.showText.push({text: this.formatJson(text), "color": color})
                    }else{
                        this.showText.push({text: "执行结果:\n"+text.Result, "color": color})
                        if(text.ErrorInfo){
                            this.showText.push({text: "错误:"+text.ErrorInfo, "color": color})
                        }
                        this.showText.push({text: "=============", "color": color})
                      //  this.showText.push({text: this.formatJson(text), "color": color})
                    }
                    if (action < (data[i].action / 10)) {
                        action = data[i].action / 10
                    }
                }
                if(action ==10){
                    setTimeout( clearInterval(this.intervalid1),5000)
                    this.showText.push({text: "发布完成", "color": color})
                }
                this.action = action
            }).catch(() => {

                })
            },
           formatJson(json, options) {
                var reg = null,
                        formatted = '',
                        pad = 0,
                        PADDING = '    '; // one can also use '\t' or a different number of spaces

                // optional settings
                options = options || {};
                // remove newline where '{' or '[' follows ':'
                options.newlineAfterColonIfBeforeBraceOrBracket = (options.newlineAfterColonIfBeforeBraceOrBracket === true) ? true : false;
                // use a space after a colon
                options.spaceAfterColon = (options.spaceAfterColon === false) ? false : true;

                // begin formatting...
                if (typeof json !== 'string') {
                    // make sure we start with the JSON as a string
                    json = JSON.stringify(json);
                } else {
                    // is already a string, so parse and re-stringify in order to remove extra whitespace
                    json = JSON.parse(json);
                    json = JSON.stringify(json);
                }

                // add newline before and after curly braces
                reg = /([\{\}])/g;
                json = json.replace(reg, '\r\n$1\r\n');

                // add newline before and after square brackets
                reg = /([\[\]])/g;
                json = json.replace(reg, '\r\n$1\r\n');

                // add newline after comma
                reg = /(\,)/g;
                json = json.replace(reg, '$1\r\n');

                // remove multiple newlines
                reg = /(\r\n\r\n)/g;
                json = json.replace(reg, '\r\n');

                // remove newlines before commas
                reg = /\r\n\,/g;
                json = json.replace(reg, ',');

                // optional formatting...
                if (!options.newlineAfterColonIfBeforeBraceOrBracket) {
                    reg = /\:\r\n\{/g;
                    json = json.replace(reg, ':{');
                    reg = /\:\r\n\[/g;
                    json = json.replace(reg, ':[');
                }
                if (options.spaceAfterColon) {
                    reg = /\:/g;
                    json = json.replace(reg, ':');
                }
                for(var index in json.split('\r\n')){
                    var node=json.split('\r\n')[index]
                    var i = 0,
                            indent = 0,
                            padding = '';

                    if (node.match(/\{$/) || node.match(/\[$/)) {
                        indent = 1;
                    } else if (node.match(/\}/) || node.match(/\]/)) {
                        if (pad !== 0) {
                            pad -= 1;
                        }
                    } else {
                        indent = 0;
                    }

                    for (i = 0; i < pad; i++) {
                        padding += PADDING;
                    }

                    formatted += padding + node + '\r\n';
                    pad += indent;
                }
                return formatted;
            }
        }
    }
</script>
