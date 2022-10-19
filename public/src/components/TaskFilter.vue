<template>
    <v-card fluid>
        <v-card-text>
            <v-list>
                <v-list-tile>
                  <v-text-field
                  append-icon="search"
                  label="搜索"
                  single-line
                  hide-details
                  v-model="filters.search"
                  ></v-text-field>
                </v-list-tile>
                <div style="height:20px"></div>
                <v-list-tile>
                  <v-select
                      v-model="filters.timerect"
                      label="时间范围"
                      item-text="name"
                      item-value="value"
                      :items="timerect"
                  ></v-select>
                </v-list-tile>
                <v-list-tile v-if="divtimerect">
                  <v-menu
                      ref="divtimerectmenustart"
                      lazy
                      :close-on-content-click="false"
                      v-model="divtimerectmenustart"
                      transition="scale-transition"
                      offset-y
                      style="width:100%"
                      :nudge-right="40"
                      min-width="290px"
                      :return-value.sync="filters.timerectstart"
                  >
                      <v-text-field
                      slot="activator"
                      label="开始时间"
                      v-model="filters.timerectstart"
                      readonly
                      ></v-text-field>
                      <v-date-picker v-model="filters.timerectstart" no-title scrollable>
                      <v-spacer></v-spacer>
                      <v-btn flat color="primary" @click="divtimerectmenustart = false">Cancel</v-btn>
                      <v-btn flat color="primary" @click="$refs.divtimerectmenustart.save(filters.timerectstart)">OK</v-btn>
                      </v-date-picker>
                  </v-menu>
                </v-list-tile>
                <v-list-tile v-if="divtimerect">
                  <v-menu
                      ref="divtimerectmenuend"
                      lazy
                      :close-on-content-click="false"
                      v-model="divtimerectmenuend"
                      transition="scale-transition"
                      offset-y
                      style="width:100%"
                      :nudge-right="40"
                      min-width="290px"
                      :return-value.sync="filters.timerectend"
                  >
                      <v-text-field
                      slot="activator"
                      label="结束时间"
                      v-model="filters.timerectend"
                      readonly
                      ></v-text-field>
                      <v-date-picker v-model="filters.timerectend" no-title scrollable>
                      <v-spacer></v-spacer>
                      <v-btn flat color="primary" @click="divtimerectmenuend = false">Cancel</v-btn>
                      <v-btn flat color="primary" @click="$refs.divtimerectmenuend.save(filters.timerectend)">OK</v-btn>
                      </v-date-picker>
                  </v-menu>
                </v-list-tile>
                <div style="height:20px"></div>
                <v-list-tile>
                  <v-select
                      label="资源"
                      :items="users"
                      v-model="filters.users"
                      multiple
                      multi-line
                      item-text="name"
                      item-value="value"
                      max-height="400"
                      persistent-hint
                  ></v-select>
                </v-list-tile>
                <div style="height:20px"></div>
                <v-list-tile>
                  <v-select
                      label="组"
                      :items="usergroups"
                      v-model="filters.usergroups"
                      multiple
                      multi-line
                      item-text="name"
                      item-value="value"
                      max-height="400"
                      persistent-hint
                  ></v-select>
                </v-list-tile>
                <div style="height:20px"></div>
                <v-list-tile>
                  <v-select
                    label="任务"
                    :items="tasks"
                    v-model="filters.tasks"
                    item-text="name"
                    item-value="id"
                    max-height="600px"
                    cache-items
                    :search-input.sync="ptasksearch"
                    :loading="ptaskloading"
                    chips
                    autocomplete
                    multiple
                  ></v-select>
                </v-list-tile>
                <div style="height:20px"></div>
                <v-list-tile>
                  <v-select
                      label="状态"
                      :items="status"
                      v-model="filters.status"
                      multiple
                      multi-line
                      item-text="name"
                      item-value="value"
                      max-height="400"
                      persistent-hint
                  ></v-select>
                </v-list-tile>
            </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn flat @click="d">默认</v-btn>
          <v-btn flat @click="cancel">取消</v-btn>
          <v-btn color="primary" flat @click="ok">确定</v-btn>
        </v-card-actions>
    </v-card>
</template>

<script>
export default {
  props: {
    show: {
      type: Boolean,
      default: false,
    },
    okfunc: {
      type: Function,
    },
    cancelfunc: {
      type: Function,
    },
  },
  data() {
    return {
      timerect: [
        {
          name: '当前星期',
          value: 0,
        },
        {
          name: '本月',
          value: 1,
        },
        {
          name: '自定义',
          value: 2,
        },
      ],
      divtimerect: false,
      divtimerectmenustart: false,
      divtimerectmenuend: false,
      filters: {},
      users: [],
      usergroups: [],
      status: [
        {
          name: '计划中',
          value: 1,
        },
        {
          name: '进行中',
          value: 2,
        },
        {
          name: '临近过期',
          value: 3,
        },
        {
          name: '提前完成',
          value: 4,
        },
        {
          name: '已完成',
          value: 5,
        },
        {
          name: '已过期',
          value: 6,
        },
        {
          name: '过期完成',
          value: 7,
        },
        {
          name: '重新打开',
          value: 8,
        },
        {
          name: '重新打开并超期',
          value: 9,
        },
      ],
      ptaskloading: false,
      ptasksearch: null,
      tasks: [],
      tasknames: [],
      taskids: [],
    };
  },
  watch: {
    'filters.timerect': 'settimerect',
    ptasksearch(val) {
      val && this.reloadTasks();
    },
  },
  mounted() {
    this.$store.state.users.forEach((element) => {
      if (element.name === '无') {
        return;
      }
      this.users.push({
        name: element.name,
        value: element.id,
      });
    });
    this.$store.state.usergroups.forEach((element) => {
      if (element.name === '无') {
        return;
      }
      this.usergroups.push({
        name: element.name,
        value: element.id,
      });
    });
    this.init();
  },
  methods: {
    d() {
      this.filters = {
        search: '',
        users: [],
        tasks: [],
        usergroups: [],
        timerect: 0,
        status: [1, 2, 3, 4, 5, 6, 7, 8, 9],
      };
    },
    init() {
      this.d();
      const filters = localStorage.getItem('filters');
      if (filters && filters.length > 0) {
        this.filters = JSON.parse(filters);
      }
      return this.filters;
    },
    ok() {
      this.okfunc(this.filters);
      localStorage.setItem('filters', JSON.stringify(this.filters));
    },
    cancel() {
      this.cancelfunc();
    },
    settimerect(v) {
      if (v === 2) {
        this.divtimerect = true;
      } else {
        this.divtimerect = false;
      }
    },

    reloadTasks(id) {
      if (!id && !this.ptasksearch) {
        return;
      }
      if (id) {
        if (this.taskids.indexOf(id) > -1) {
          return;
        }
      }
      if (this.tasknames.indexOf(this.ptasksearch) > -1) {
        return;
      }

      this.ptaskloading = true;
      new Promise((resolve) => {
        let url = `/task?filter=${this.ptasksearch}`;
        if (id) {
          url = `/task/${id}`;
        }

        this.$http.LGet(this.$store.state, url).then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                if (this.tasknames.indexOf(element.Name) === -1) {
                  this.tasknames.push(element.Name);
                  this.taskids.push(element.ID);
                  this.tasks.push({
                    id: element.ID,
                    name: element.Name,
                  });
                }
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve();
        }).catch((e) => {
          console.dir('服务器报错');
          if (e.response.data.code === 101) {
            this.$store.commit('logout');
            this.$router.replace('/login');
          }
          resolve();
        });
      }).then(() => {
        this.ptaskloading = false;
      });
    },
  },
};
</script>
