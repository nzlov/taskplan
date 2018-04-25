<template>
    <v-card>
        <v-card-title>
            <v-select
            label="范围"
            :items="listtypes"
            v-model="listtype"
            item-text="name"
            item-value="value"
            max-height="600px"
            hint="选择任务范围"
            ></v-select>
            <v-spacer></v-spacer>
            <v-text-field
            append-icon="search"
            label="搜索"
            single-line
            hide-details
            v-model="search"
            @keyup.enter="searchData"
            ></v-text-field>
        </v-card-title>
        <v-data-table
            :headers="headers"
            :items="items"
            :pagination.sync="pagination"
            :total-items="totalItems"
            :loading="loading"
            class="elevation-1"
        >
            <template slot="items" slot-scope="props">
            <tr @click="props.expanded = !props.expanded" :style="genbackground(props.item)">
                <td>{{ props.item.name }}</td>
                <td class="text-xs-center">{{ props.item.usergroup }}</td>
                <td class="text-xs-center">{{ props.item.user }}</td>
                <td class="text-xs-center">{{ props.item.createuser }}</td>
                <td class="text-xs-center">{{ props.item.creates }}</td>
                <td class="text-xs-center">{{ props.item.starts }}</td>
                <td class="text-xs-center">{{ props.item.ends }}</td>
                <td class="text-xs-center">{{ props.item.realends }}</td>
                <td class="text-xs-center"> {{ props.item.statuss }}</td>
                <td class="text-xs-center">{{ props.item.time }}</td>
                <td class="justify-center layout px-0">
                <v-tooltip bottom v-if="props.item.edit && ( editP && props.item.start > formatTimestamp(new Date()) || editP && props.item.start > 0 && props.item.start  <= formatTimestamp(new Date()) && expireP )">
                    <v-btn icon slot="activator" class="mx-0" @click.stop="editItem(items,props.item)">
                    <v-icon color="teal">edit</v-icon>
                    </v-btn>
                    <span>编辑</span>
                </v-tooltip>
                <v-tooltip bottom v-if="props.item.edit && ( !props.item.ptask && openP && props.item.status == 2 )">
                    <v-btn icon slot="activator" class="mx-0" @click.stop="openItem(props.item)">
                    <v-icon color="orange">restore</v-icon>
                    </v-btn>
                    <span>再次打开</span>
                </v-tooltip>
                <v-tooltip bottom v-if="props.item.edit && ( !props.item.ptask && doneP && props.item.status != 2 )">
                    <v-btn icon slot="activator" class="mx-0" @click.stop="doneItem(props.item)">
                    <v-icon color="primary">check_circle</v-icon>
                    </v-btn>
                    <span>完成</span>
                </v-tooltip>
                <v-tooltip bottom v-if="props.item.edit && delP">
                    <v-btn icon slot="activator" class="mx-0" @click.stop="delItem(props.item)">
                    <v-icon color="pink">delete</v-icon>
                    </v-btn>
                    <span>删除</span>
                </v-tooltip>
                </td>
            </tr>
            </template>
            <template slot="expand" slot-scope="props">
                <v-card flat style="margin: 20px;" v-if="props.item.ptask">
                  <v-card-title>
                    <span class="headline">子任务列表</span>
                  </v-card-title>
                  <task-table
                    :editP="editP" 
                    :openP="openP" 
                    :doneP="doneP" 
                    :expireP="expireP" 
                    :delP="delP" 
                    :pid="props.item.id"
                    :showcolor="showcolor"
                    :editItem="editItem"
                    :openItem="openItem"
                    :doneItem="doneItem"
                    :delItem="delItem"
                  ></task-table>
                </v-card>
                <v-card flat style="margin: 20px;" v-if="!props.item.ptask">
                  <v-card-title>
                    <span class="headline">任务详细时间</span>
                  </v-card-title>
                  <v-list style="margin: 20px;">
                    <v-list-tile avatar v-for="item in props.item.timerect" :key="item.Start">
                      <v-list-tile-content>
                        <v-list-tile-title v-text="item.Str"></v-list-tile-title>
                      </v-list-tile-content>
                    </v-list-tile>
                  </v-list>
                </v-card>
                <v-card flat style="margin: 20px;">
                  <v-card-title>
                    <span class="headline">任务变更记录</span>
                  </v-card-title>
                  <v-data-table
                  :headers="headers2"
                  :items="props.item.history"
                  hide-actions
                  >
                  <template slot="items" slot-scope="props">
                      <td class="text-xs-center">{{ formatDate(props.item.CreatedAt) }}</td>
                      <td class="text-xs-center">{{ props.item.User.Name }}</td>
                      <td class="text-xs-center">{{ formatAction(props.item.Action) }}</td>
                      <td class="text-xs-center" v-html="formatActionItems(props.item.Items)"></td>
                  </template>
                  </v-data-table>
                </v-card>
            </template>
            <v-alert slot="no-data" :value="true" color="error" icon="warning">
                {{nodata}}
            </v-alert>
        </v-data-table>
    </v-card>
</template>

<script>
import DateUtil from '../utils/date';

export default {
  props: {
    editP: {
      defalut: false,
    },
    openP: {
      defalut: false,
    },
    doneP: {
      defalut: false,
    },
    expireP: {
      defalut: false,
    },
    delP: {
      defalut: false,
    },
    pid: {
      defalut: '0',
    },
    showcolor: {
      defalut: false,
    },
    editItem: Function,
    openItem: Function,
    doneItem: Function,
    delItem: Function,
  },
  data() {
    return {
      headers: [
        {
          text: '任务',
          align: 'left',
          value: 'name',
          width: '400px',
        },
        {
          text: '用户组',
          align: 'center',
          value: 'user_group_id',
          width: '120px',
          sortable: false,
        },
        {
          text: '执行人',
          align: 'center',
          value: 'user_id',
          width: '120px',
          sortable: false,
        },
        {
          text: '创建人',
          align: 'center',
          value: 'create_user_id',
          width: '120px',
          sortable: false,
        },
        {
          text: '创建时间',
          align: 'center',
          value: 'created_at',
          width: '150px',
        },
        {
          text: '计划开始时间',
          align: 'center',
          value: 'start',
          width: '150px',
        },
        {
          text: '计划结束时间',
          align: 'center',
          value: 'end',
          width: '150px',
        },
        {
          text: '真正结束时间',
          align: 'center',
          value: 'real_end',
          width: '150px',
        },
        {
          text: '状态',
          align: 'center',
          value: 'status',
        },
        {
          text: '计时',
          align: 'center',
          value: 'time',
          width: '200px',
        },
        {
          text: 'Actions',
          align: 'center',
          value: 'name',
          sortable: false,
        },
      ],

      headers2: [
        {
          text: '时间',
          align: 'center',
          value: 'CreatedAt',
          width: '200px',
        },
        {
          text: '执行人',
          align: 'center',
          value: 'UserID',
          width: '120px',
        },
        {
          text: '动作',
          align: 'center',
          value: 'Action',
          width: '120px',
        },
        {
          text: '内容',
          align: 'center',
          value: 'Items',
          sortable: false,
        },
      ],
      search: '',
      alert_error: false,
      totalItems: 0,
      items: [],
      loading: true,
      pagination: {
        descending: false,
        page: 1,
        rowsPerPage: 10,
        sortBy: 'created_at',
        totalItems: 0,
      },
      listtype: 'self',
      listtypes: [
        {
          name: '自己',
          value: 'self',
        },
        {
          name: '全部',
          value: '',
        },
      ],
      nodata: '没有数据',
    };
  },

  computed: {
    formTitle() {
      return this.editedIndex === -1 ? '新增' : `编辑-${this.editedItem.name}`;
    },
  },

  created() {
    this.$store.state.usergroups.forEach(element => {
      if (element.name === '无') {
        return;
      }
      this.listtypes.push({
        name: element.name,
        value: element.id,
      })
    });
  },

  watch: {
    dialog(val) {
      val || this.close();
    },
    pagination: {
      handler() {
        this.updateData();
      },
      deep: true,
    },
    listtype() {
      this.updateData();
    },
  },

  methods: {
    searchData() {
      this.getDataFromApi().then((data) => {
        this.items = data.items;
        this.totalItems = data.total;
      });
    },
    updateData() {
      this.getDataFromApi()
      .then((data) => {
        this.items = data.items;
        this.totalItems = data.total;
      });
    },
    getDataFromApi() {
      this.loading = true;
      return new Promise((resolve) => {
        const { sortBy, descending, page, rowsPerPage } = this.pagination;
        const items = [];
        let total = 0;
        let order = 'order=';
        if (sortBy) {
          if (!descending) {
            order = `${order}-`;
          }
          order = `${order}${sortBy}`;
        }
        const offset = (page - 1) * rowsPerPage;
        let search = '';
        if (this.search) {
          search = `&filter=${this.search}`;
        }
        this.$http
          .LGet(this.$store.state, `/task?all=t&pid=${this.pid}&list=${this.listtype}&${order}&offset=${offset}&limit=${rowsPerPage}${search}`)
          .then((resp) => {
            this.loading = false;
            switch (resp.data.code) {
              case 0: {
                resp.data.data.data.forEach((element) => {
                  const ttimerect = [];
                  if (element.TimeRect) {
                    element.TimeRect.forEach((v) => {
                      ttimerect.push({
                        Start: v.Start,
                        End: v.End,
                        Str: `${this.formatDate(v.Start)} - ${this.formatDate(v.End)}`,
                      });
                    });
                  }

                  items.push({
                    id: element.ID,
                    name: element.Name,
                    description: element.Description,
                    usergroup: element.UserGroup.Name,
                    usergroupid: element.UserGroup.ID,
                    user: element.User.RealName,
                    userid: element.User.ID,
                    taskid: element.ParentTaskID,
                    ptask: element.PTask,
                    createuser: element.CreateUser.RealName,
                    createuserid: element.CreateUser.ID,
                    creates: this.formatDate(element.CreatedAt),
                    start: element.Start,
                    starts: this.formatDate(element.Start),
                    end: element.End,
                    ends: this.formatDate(element.End),
                    realend: element.RealEnd,
                    realends: this.formatDate(element.RealEnd),
                    time: this.formatTimeSince(element.Start, element.RealEnd),
                    timerect: ttimerect,
                    status: element.Status,
                    statuss: this.formatStatus(element),
                    history: element.TaskHistory,
                    edit: this.$store.state.usergroup === element.UserGroup.ID,
                  });
                });
                total = resp.data.data.total;
                break;
              }
              default:
                this.alert_error = true;
                this.login = false;
            }
            if (items.length === 0) {
              this.nodata = '没有数据';
            }
            resolve({
              items,
              total,
            });
          }).catch((e) => {
            this.loading = false;
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
            this.nodata = '服务器报错';
            resolve({
              items,
              total,
            });
          });
      });
    },
    getDesserts() {
      return this.datas;
    },
    allowedDates(v) {
      return this.$store.state.holidays.indexOf(v) === -1;
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
    genbackground(v) {
      const curr = this.formatTimestamp(new Date());
      if (!this.showcolor) {
        return {};
      }
      switch (v.status) {
        case 1: {
          if (v.start > curr) {
            // 计划中
            return { background: '#99CCFF' };
          }
          if (v.end < curr + 7200) {
            if (v.end < curr) {
              // 已超期
              return { background: '#FF0033' };
            }
            // 临近过期
            return { background: '#FF6600' };
          }
          // 进行中
          return { background: '#66CC99' };
        }
        case 2: {
          // 超期完成
          if (v.end <= v.realend - 3600) {
            return { background: '#FF6666' };
          }
          // 提前完成
          if (v.end > v.realend + 3600) {
            return { background: '#33CC99' };
          }
          // 按时完成
          return { background: '#009933' };
        }
        case 3: {
          if (v.end < curr) {
            // 重新打开并超期
            return { background: '#CC0066' };
          }
          // 重新打开
          return { background: '#009999' };
        }
        default: {
          return {};
        }
      }
    },
    // 格式化任务记录里的Action
    formatAction(v) {
      switch (v) {
        case 1: {
          return '创建任务';
        }
        case 2: {
          return '编辑任务';
        }
        case 3: {
          return '完成任务';
        }
        case 4: {
          return '重新打开';
        }
        case 5: {
          return '删除任务';
        }
        default: {
          return `未知类型[${v}]`;
        }
      }
    },
    // 格式化任务记录里的Items
    formatActionItems(v) {
      let str = '';
      if (v && v instanceof Array) {
        v.forEach((i) => {
          if (i.Field === 'Remark') {
            str = `${str}<br /> 备注:${i.New}`;
          } else {
            str = `${str}<br /> ${i.Field}:${i.Old} => ${i.New}`;
          }
        });
      }
      return str;
    },
    // TODO 为开始 进行中 延时 完成
    formatTimeSince(a, b) {
      if (a) {
        if (b) {
          return this.formatSecond(b - a);
        }
        return this.formatSecond((Date.parse(new Date()) / 1000) - a);
      }
      return '未开始';
    },
    formatSecond(a) {
      let s = a;
      if (s < 0) {
        s *= -1;
      }
      return `${(s / 86400).toFixed(2)} 天`;
    },
    formatTimestamp(now) {
      return DateUtil.formatTimestamp(now);
    },
    formatDate(d) {
      return DateUtil.formatDate(d);
    },
  },
};
</script>