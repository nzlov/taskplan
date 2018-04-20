<template>
  <v-container fluid>
      <!-- 添加编辑对话框 -->
      <v-dialog
        v-model="dialog"
        fullscreen
        hide-overlay
        transition="dialog-bottom-transition"
        scrollable
      >
      <!-- 操作对话框 -->
        <v-card title>
          <v-toolbar card dark color="primary">
            <v-btn icon @click.native="close" dark>
              <v-icon>close</v-icon>
            </v-btn>
            <v-toolbar-title>{{ formTitle }}</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-toolbar-items>
              <v-btn dark flat @click.native="save" :disabled="!valid || commit">Save</v-btn>
            </v-toolbar-items>
          </v-toolbar> 
          <v-progress-linear style="margin-bottom: 0px;" :indeterminate="true" :active="commit"></v-progress-linear>
          <v-alert
              type="error"
              dismissible 
              v-model="alert_error"
              transition="scale-transition"
          >
          {{ message }}
          </v-alert>         
          <v-card-text>
            <v-container grid-list-md>
              <v-layout>
                <v-flex xs12 sm12 md12>
                  <v-form v-model="valid" ref="form" lazy-validation>
                    <v-text-field
                      v-if="editedItem.start > 0 && editedItem.start <= formatTimestamp(new Date())"
                      name="input-7-1"
                      label="已开始的任务修改请填写备注"
                      v-model="editedItem.remark"
                      :rules="remarkRules"
                      required
                    ></v-text-field>
                    <v-text-field
                    label="任务名称"
                    v-model="editedItem.name"
                    :rules="nameRules"
                    required
                    ></v-text-field>
                    <v-flex xs12 sm12 md12>
                      <v-select
                        label="父级任务"
                        :items="tasks"
                        v-model="editedItem.taskid"
                        item-text="name"
                        item-value="id"
                        max-height="600px"
                        hint="选择父级任务"
                      ></v-select>
                    </v-flex>
                    <v-flex xs12 sm12 md12>
                      <v-select
                        label="用户组"
                        :items="usergroups"
                        v-model="editedItem.usergroupid"
                        item-text="name"
                        item-value="id"
                        max-height="600px"
                        hint="选择用户组"
                      ></v-select>
                    </v-flex>
                    <v-flex xs12 sm12 md12>
                      <v-select
                        label="执行人"
                        :items="users"
                        v-model="editedItem.userid"
                        item-text="name"
                        item-value="id"
                        max-height="600px"
                        hint="选择执行人"
                      ></v-select>
                    </v-flex>
                    <v-card v-if="!editedItem.ptask">
                      <v-card-title>
                        <span class="headline">计划开始时间</span>
                      </v-card-title>
                      <v-card-text>
                        <v-date-picker locale="zh-cn" :first-day-of-week="1" :allowed-dates="allowedDates" v-model="startdate"></v-date-picker>
                        <v-time-picker locale="zh-cn" v-model="editedItem.starttime" format="24hr"></v-time-picker>
                      </v-card-text>
                    </v-card>
                    <div v-if="!editedItem.ptask" style="height: 10px">
                    </div>
                    <h2 v-else>父级任务不能手动修改时间</h2>
                    <v-card v-if="!editedItem.ptask">
                      <v-card-title>
                        <span class="headline">计划结束时间</span>
                      </v-card-title>
                      <v-card-text>
                        <v-date-picker locale="zh-cn" :first-day-of-week="1" :allowed-dates="allowedDates" v-model="enddate"></v-date-picker>
                        <v-time-picker locale="zh-cn" v-model="editedItem.endtime" format="24hr"></v-time-picker>
                      </v-card-text>
                    </v-card>
                    <v-text-field
                      name="input-7-1"
                      label="简介"
                      multi-line
                      v-model="editedItem.description"
                    ></v-text-field>
                </v-form>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>
         </v-card>
      </v-dialog>
      <!-- 操作对话框 -->
      <v-dialog v-model="actionDialog" persistent max-width="500px">
        <v-card>
          <v-card-title>
            <span> {{ actionMessage }} </span>
          </v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="error" flat @click="actione" :disabled="commit">Yes</v-btn>
            <v-btn color="primary" flat @click.stop="actionDialog=false;commit=false" :disabled="commit">No</v-btn>
          </v-card-actions>
          <v-alert
              type="error"
              dismissible 
              v-model="alert_error"
              transition="scale-transition"
          >
          {{ message }}
          </v-alert>
          <v-progress-linear style="margin-bottom: 0px;" :indeterminate="true" :active="commit"></v-progress-linear>
        </v-card>
      </v-dialog>
      <!-- 表格 -->
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
            <tr @click="props.expanded = !props.expanded">
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
                <v-tooltip bottom v-if="editP && props.item.start > formatTimestamp(new Date()) || editP && props.item.start > 0 && props.item.start  <= formatTimestamp(new Date()) && expireP">
                  <v-btn icon slot="activator" class="mx-0" @click.stop="editItem(props.item)">
                    <v-icon color="teal">edit</v-icon>
                  </v-btn>
                  <span>编辑</span>
                </v-tooltip>
                <v-tooltip bottom v-if="!props.item.ptask && openP && props.item.status == 2">
                  <v-btn icon slot="activator" class="mx-0" @click.stop="openItem(props.item)">
                    <v-icon color="orange">restore</v-icon>
                  </v-btn>
                  <span>再次打开</span>
                </v-tooltip>
                <v-tooltip bottom v-if="!props.item.ptask && doneP && props.item.status != 2">
                  <v-btn icon slot="activator" class="mx-0" @click.stop="doneItem(props.item)">
                    <v-icon color="primary">check_circle</v-icon>
                  </v-btn>
                  <span>完成</span>
                </v-tooltip>
                <v-tooltip bottom v-if="delP">
                  <v-btn icon slot="activator" class="mx-0" @click.stop="deleteItem(props.item)">
                    <v-icon color="pink">delete</v-icon>
                  </v-btn>
                  <span>删除</span>
                </v-tooltip>
              </td>
            </tr>
          </template>
          <template slot="expand" slot-scope="props">
            <v-card flat style="margin: 20px;">
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
            没有数据
          </v-alert>
        </v-data-table>

      </v-card>
      <v-btn
        v-if="addP"
        fixed
        dark
        fab
        bottom
        right
        color="pink"
        @click.native="newItem"
      >
      <v-icon>add</v-icon>
      </v-btn>
  </v-container>
</template>

<script>
import DateUtil from '../utils/date';

export default {
  data() {
    return {
      dialog: false,
      dialog2: false,
      valid: true,
      addP: false,
      editP: false,
      delP: false,
      expireP: false,
      openP: false,
      doneP: false,
      message: '',
      startdate: '2006-01-01',
      enddate: '2006-01-01',
      action: 0, // 操作类型 1 删除 2 开始 3 完成 4 延期
      actionDialog: false,
      actionMessage: '',
      actionItem: {},
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
          width: '300px',
        },
        {
          text: '计划开始时间',
          align: 'center',
          value: 'start',
          width: '300px',
        },
        {
          text: '计划结束时间',
          align: 'center',
          value: 'end',
          width: '300px',
        },
        {
          text: '真正结束时间',
          align: 'center',
          value: 'real_end',
          width: '300px',
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
        { text: 'Actions',
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
      commit: false,
      search: '',
      alert_error: false,
      totalItems: 0,
      items: [],
      loading: true,
      pagination: {},
      editedIndex: -1,
      nameRules: [
        v => !!v || '任务名字必须填写',
      ],
      remarkRules: [
        v => !!v || '任务已开始请填写备注',
      ],
      oeditedItem: {},
      editedItem: {
        id: '',
        name: '',
        tag: '',
        usergroupid: '',
        startdate: '2016-01-01',
        starttime: '00:00',
        enddate: '2016-01-01',
        endtime: '23:59',
      },
      defaultItem: {
        id: '',
        name: '',
        tag: '',
        usergroupid: '',
        startdate: '2016-01-01',
        starttime: '00:00',
        enddate: '2016-01-01',
        endtime: '23:59',
      },
      otasks: [],
      tasks: [],
      users: [],
      usergroups: [],
      listtype: 'self',
      listtypes: [
        {
          name: '自己',
          value: 'self',
        },
        {
          name: '组',
          value: 'group',
        },
        {
          name: '全部',
          value: '',
        },
      ],
    };
  },

  computed: {
    formTitle() {
      return this.editedIndex === -1 ? '新增' : `编辑-${this.editedItem.name}`;
    },
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

  mounted() {
    const p = this.$store.state.permissions;
    this.addP = p['task.add'] !== undefined;
    this.editP = p['task.update'] !== undefined;
    this.delP = p['task.del'] !== undefined;
    this.expireP = p['task.expire'] !== undefined;
    this.openP = p['task.open'] !== undefined;
    this.doneP = p['task.done'] !== undefined;

    this.users = this.$store.state.users;
    this.usergroups = this.$store.state.usergroups;
  },

  methods: {
    searchData() {
      this.getDataFromApi()
      .then((data) => {
        this.items = data.items;
        this.totalItems = data.total;
      });
    },
    newItem() {
      this.editedIndex = -1;
      this.editedItem = {
        starttime: DateUtil.formatTime(new Date()),
        endtime: DateUtil.formatTime(new Date()),
      };
      this.startdate = DateUtil.formatDate1(new Date());
      this.enddate = DateUtil.formatDate1(new Date());
      this.dialog = true;
      this.valid = true;
    },

    editItem(item) {
      this.editedIndex = this.items.indexOf(item);
      this.oeditedItem = Object.assign({}, item);
      this.editedItem = Object.assign({}, item);
      if (item.starts) {
        const date = new Date(item.starts);
        this.startdate = DateUtil.formatDate1(date);
        this.editedItem.starttime = DateUtil.formatTime(date);
      }
      if (item.ends) {
        const date = new Date(item.ends);
        this.enddate = DateUtil.formatDate1(date);
        this.editedItem.endtime = DateUtil.formatTime(date);
      }
      this.tasks = [];
      this.otasks.forEach((v) => {
        if (v.id !== item.id) {
          this.tasks.push(v);
        }
      });
      this.dialog = true;
    },

    deleteItem(item) {
      this.action = 1;
      this.actionItem = item;
      this.actionDialog = true;
      this.actionMessage = `真的要删除${item.name}?`;
    },

    openItem(item) {
      this.action = 3;
      this.actionItem = item;
      this.actionDialog = true;
      this.actionMessage = `真的要再次打开${item.name}?`;
    },

    doneItem(item) {
      this.action = 2;
      this.actionItem = item;
      this.actionDialog = true;
      this.actionMessage = `确定完成${item.name}?`;
    },

    close() {
      this.dialog = false;
      this.actionDialog = false;
      setTimeout(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.delItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      }, 300);
    },
    actione() {
      this.alert_error = false;
      this.commit = true;
      // 1 删除 2 完成 3 重新打开
      switch (this.action) {
        case 1: {
          this.$http.LDel(this.$store.state, `/task/${this.actionItem.id}`).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.updateData();
                this.close();
                break;
              }
              default:
                this.message = '服务器错误601';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch(() => {
            this.message = '服务器错误606';
            this.alert_error = true;
            this.commit = false;
          });
          break;
        }
        case 2: {
          this.$http.LPost(this.$store.state, `/task/${this.actionItem.id}/done`).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.actionItem.statuss = '已完成';
                this.actionItem.status = 2;
                this.close();
                break;
              }
              default:
                this.message = '服务器错误623';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch(() => {
            this.message = '服务器错误628';
            this.alert_error = true;
            this.commit = false;
          });
          break;
        }
        case 3: {
          this.$http.LPost(this.$store.state, `/task/${this.actionItem.id}/open`).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.actionItem.statuss = '重新打开';
                this.actionItem.status = 3;
                this.close();
                break;
              }
              default:
                this.message = '服务器错误645';
                this.alert_errosr = true;
            }
            this.commit = false;
          }).catch(() => {
            this.message = '服务器错误650';
            this.alert_error = true;
            this.commit = false;
          });
          break;
        }
        default: {
          this.alert_error = true;
          this.message = '不支持的操作';
        }
      }
    },
    save() {
      this.alert_error = false;
      if (this.editedIndex > -1) {
        if (this.$refs.form.validate()) {
          this.commit = true;
          let has = false;
          const d = {};
          if (this.editedItem.name !== this.oeditedItem.name) {
            d.name = this.editedItem.name;
            has = true;
          }
          if (this.editedItem.usergroupid !== this.oeditedItem.usergroupid) {
            d.usergroupid = this.editedItem.usergroupid;
            has = true;
          }
          if (this.editedItem.userid !== this.oeditedItem.userid) {
            if (this.editedItem.userid) {
              d.userid = this.editedItem.userid;
            } else {
              d.userid = 0;
            }
            has = true;
          }
          if (this.editedItem.taskid !== this.oeditedItem.taskid) {
            d.taskid = this.editedItem.taskid;
            has = true;
          }
          if (this.editedItem.description !== this.oeditedItem.description) {
            d.description = this.editedItem.description;
            has = true;
          }
          if (this.startdate && this.editedItem.starttime) {
            const start = this.formatTimestamp(`${this.startdate} ${this.editedItem.starttime}`);
            if (this.oeditedItem.start !== start) {
              d.start = start;
              has = true;
            }
          }
          if (this.enddate && this.editedItem.endtime) {
            const end = this.formatTimestamp(`${this.enddate} ${this.editedItem.endtime}`);
            if (this.oeditedItem.end !== end) {
              d.end = end;
              has = true;
            }
          }
          if (this.editedItem.remark) {
            d.remark = this.editedItem.remark;
            has = true;
          }
          if (!has) {
            this.commit = false;
            this.close();
            return;
          }
          this.$http.LPost(this.$store.state, `/task/${this.editedItem.id}`, d).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.updateData();
                this.close();
                break;
              }
              case 300: {
                this.message = '名字重复';
                this.alert_error = true;
                break;
              }
              default:
                this.message = '服务器错误730';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch((e) => {
            this.message = e.response.data.data;
            this.alert_error = true;
            this.commit = false;
          });
        }
      } else {
        if (this.$refs.form.validate()) {
          this.commit = true;
          const d = {
            name: this.editedItem.name,
            usergroupid: this.editedItem.usergroupid,
            userid: this.editedItem.userid,
            description: this.editedItem.description,
          };
          if (this.startdate && this.editedItem.starttime) {
            d.start = this.formatTimestamp(`${this.startdate} ${this.editedItem.starttime}`);
          }
          if (this.enddate && this.editedItem.endtime) {
            d.end = this.formatTimestamp(`${this.enddate} ${this.editedItem.endtime}`);
          }
          if (this.editedItem.taskid) {
            d.taskid = this.editedItem.taskid;
          }
          this.$http.LPost(this.$store.state, '/task', d).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.updateData();
                this.reloadTasks();
                this.close();
                break;
              }
              case 300: {
                this.message = '名字重复';
                this.alert_error = true;
                break;
              }
              default:
                this.message = '服务器错误773';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch(() => {
            this.message = '服务器错误778';
            this.alert_error = true;
            this.commit = false;
          });
        }
      }
    },
    updateData() {
      this.getDataFromApi()
      .then((data) => {
        this.items = data.items;
        this.totalItems = data.total;
      });
    },
    reloadTasks() {
      new Promise((resolve) => {
        const objs = [{
          name: '无',
          id: '0',
        }];
        this.$http.LGet(this.$store.state, '/task').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                objs.push({
                  id: element.ID,
                  name: element.Name,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            objs,
          });
        }).catch(() => {
          console.dir('服务器报错');
          resolve({
            objs,
          });
        });
      }).then((data) => {
        this.tasks = data.objs;
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
        this.$http.LGet(this.$store.state, `/task?all=t&list=${this.listtype}&${order}&offset=${offset}&limit=${rowsPerPage}${search}`).then((resp) => {
          this.loading = false;
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
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
                  creates: DateUtil.formatDate(element.CreatedAt),
                  start: element.Start,
                  starts: DateUtil.formatDate(element.Start),
                  end: element.End,
                  ends: DateUtil.formatDate(element.End),
                  realend: element.RealEnd,
                  realends: DateUtil.formatDate(element.RealEnd),
                  time: this.formatTimeSince(element.Start, element.RealEnd),
                  status: element.Status,
                  statuss: this.formatStatus(element),
                  history: element.TaskHistory,
                });
              });
              total = resp.data.data.total;
              break;
            }
            default:
              this.alert_error = true;
              this.login = false;
          }
          resolve({
            items,
            total,
          });
        }).catch((e) => {
          this.loading = false;
          console.dir('服务器报错');
          console.dir(e);
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
      switch (v.Status) {
        case 1: {
          if (v.Start > this.formatTimestamp(new Date())) {
            return '计划中';
          }
          return '进行中';
        }
        case 2: {
          return '已完成';
        }
        case 3: {
          return '重新打开';
        }
        default: {
          return `未知类型[${v}]`;
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
  },
};
</script>