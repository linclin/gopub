<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body"
         v-loading="load_data"
         element-loading-text="拼命加载中">
      <terminal :taskId="-1"></terminal>
    </div>
  </div>
</template>
<script type="text/javascript">
  import {panelTitle, terminal} from 'components'
  import {port_conf, port_code} from 'common/port_uri'
  import {tools_verify} from 'common/tools'

  export default {
    data() {
      return {
        route_id: this.$route.params.id,
        load_data: false,
      }
    },
    created() {

      if (this.route_id) {
        this.detection_data()
      } else {
        this.$message({
          message: "项目id不存在",
          type: 'warning'
        })
        setTimeout(() => {
            this.$router.push({
              name: 'confList'
            })
          },
          500
        )
      }
    },
    methods: {
      detection_data() {
        this.$http.get(port_conf.detection, {
          params: {
            projectId: this.route_id,
          }
        })
          .then(({data: {data}}) => {
            this.$message({
              message: "检测成功",
              type: 'success'
            })
          })
          .catch(() => {
            }
          )
      }
    },
    components: {
      panelTitle,
      terminal
    }
  }
</script>
