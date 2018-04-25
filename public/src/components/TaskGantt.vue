<template>
  <v-container style="height:100%;">
    <gantt id="gantt" style="height:100%;" ref="gantt"></gantt>
  </v-container>
</template>

<script>
import DateUtil from '../utils/date';

export default {
  data() {
    return {
    };
  },
  mounted() {
    this.load();
  },
  methods: {
    load() {
      new Promise((resolve) => {
        const items = [];
        this.$http.LGet(this.$store.state, '/task?all=t&order=start').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                items.push({
                  id: element.ID,
                  text: element.Name,
                  parent: element.ParentTaskID,
                  user: element.User.RealName,
                  start_date: this.formatDate(element.Start),
                  duration: this.getDuration(element),
                  open: true,

                  description: element.Description,
                  usergroup: element.UserGroup.Name,
                  usergroupid: element.UserGroup.ID,
                  userid: element.User.ID,
                  ptask: element.PTask,
                  createuser: element.CreateUser.RealName,
                  createuserid: element.CreateUser.ID,
                  start: element.Start,
                  end: element.End,
                  realend: element.RealEnd,
                  status: element.Status,
                  statuss: this.formatStatus(element),
                  history: element.TaskHistory,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            items,
          });
        }).catch((e) => {
          console.dir(e);
          if (e.response.data.code === 101) {
            this.$store.commit('logout');
            this.$router.replace('/login');
          }
          resolve({
            objs,
          });
        });
      }).then((data) => {
        this.$refs.gantt.load({
          data: data.items,
          link: [],
        });
      });
    },
    getDuration(v) {
      // 获取任务间隔
      if (v.Start) {
        if (v.RealEnd) {
          return Math.ceil((v.RealEnd - v.Start) / 86400);
        }
        if (v.End) {
          return Math.ceil((v.End - v.Start) / 86400);
        }
      }
      return 0;
    },
    abs(a) {
      let s = a;
      if (s < 0) {
        s *= -1;
      }
      return s;
    },
    // 格式化任务Status
    formatStatus(v) {
      const curr = this.formatTimestamp(new Date());
      switch (v.Status) {
        case 1: {
          if (v.Start > curr) {
            return '计划中';
          }
          if (v.End < curr + 7200) {
            if (v.End < curr) {
              return '已超期';
            }
            return '临近过期';
          }
          return '进行中';
        }
        case 2: {
          if (v.End <= v.RealEnd - 3600) {
            return '超期完成';
          }
          if (v.End > v.RealEnd + 3600) {
            return '提前完成';
          }
          return '按时完成';
        }
        case 3: {
          if (v.End < curr) {
            return '重新打开并超期';
          }
          return '重新打开';
        }
        default: {
          return `未知类型[${v}]`;
        }
      }
    },
    formatDate(v) {
      if (v <= 0) {
        return '';
      }
      const now = new Date(v * 1000);
      const year = now.getFullYear();
      let month = now.getMonth() + 1;
      let date = now.getDate();
      if (month <= 9) {
        month = `0${month}`;
      }
      if (date <= 9) {
        date = `0${date}`;
      }
      return `${date}-${month}-${year}`;
    },
    formatTimestamp(now) {
      return DateUtil.formatTimestamp(now);
    },
  },
};
</script>