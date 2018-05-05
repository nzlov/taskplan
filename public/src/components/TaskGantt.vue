<template>
  <v-container style="height:100%;">
    <v-card style="height:100%;">
      <v-card-title>
        <v-spacer></v-spacer>
        <v-btn @click.stop="exportdata" dark color="indigo">
          导出<v-icon right>cloud_download</v-icon>
        </v-btn>
        <v-menu
          offset-x
          :close-on-content-click="false"
          :nudge-width="400"
          v-model="menu"
        >
          <v-btn icon  slot="activator"><v-icon>more_vert</v-icon></v-btn>
          <task-filter ref="filters" :show="menu" :okfunc="filter" :cancelfunc="() => {menu = false}"></task-filter>
        </v-menu>
      </v-card-title>
      <gantt id="gantt" style="height:96%;" ref="gantt"></gantt>
    </v-card>
  </v-container>
</template>

<script>
import DateUtil from '../utils/date';

export default {
  data() {
    return {
      menu: false,
      filters: {},
    };
  },
  mounted() {
    this.filters = this.$refs.filters.init();
    this.load();
  },
  methods: {
    filter(v) {
      console.dir(v);
      this.filters = v;
      this.menu = false;
      this.load();
    },
    exportdata() {
      this.$refs.gantt.exportdata();
    },
    load() {
      new Promise((resolve) => {
        const items = [];
        this.$http.LGet(this.$store.state, `/taskn?all=t&order=start&filters=${JSON.stringify(this.filters)}`).then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                if (element.Start === 0 || element.End === 0) {
                  return;
                }
                items.push(this.taskToGantt(items, element));
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
    taskchild(items, tasks) {
      tasks.forEach((element) => {
        if (element.Start === 0 || element.End === 0) {
          return;
        }
        items.push(this.taskToGantt(items, element));
      });
    },
    taskToGantt(items, element) {
      let task = {
        id: element.ID,
        text: element.Name,
        parent: element.ParentTaskID,
        user: element.User.RealName,
        start_date: this.formatDate(element.Start),
        end_date: this.formatDate(element.End),
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
        real_date: this.formatDate(element.RealEnd),
        status: element.Status,
        statuss: this.formatStatus(element),
        history: element.TaskHistory,
      };
      if (task.ptask) {
        this.taskchild(items, element.Tasks);
      }
      return task;
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