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
                        cache-items
                        :search-input.sync="ptasksearch"
                        :loading="ptaskloading"
                        autocomplete
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
      <!-- Main -->
      <v-card>
        <v-card-title>
          <v-switch :label="'颜色'" v-model="showcolor"></v-switch>
        </v-card-title>
        <task-table 
          :editP="editP" 
          :openP="openP" 
          :doneP="doneP" 
          :expireP="expireP" 
          :delP="delP"
          :showcolor="showcolor"
          :editItem="editItem"
          :openItem="openItem"
          :doneItem="doneItem"
          :delItem="deleteItem"
          ref="tasktable"
        ></task-table>
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
      showcolor: false,
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
      commit: false,
      alert_error: false,
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
      ptaskloading: false,
      ptasksearch: null,
      tasks: [{
        name: '无',
        id: '0',
      }],
      tasknames: [],
      taskids: [],
      users: [],
      usergroups: [],
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
    ptasksearch(val) {
      val && this.reloadTasks();
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

    this.reloadTasks();
    this.users = this.$store.state.users;
    this.usergroups = this.$store.state.usergroups;
  },

  methods: {
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

    editItem(items, item) {
      this.editedIndex = items.indexOf(item);
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
      this.reloadTasks(item.taskid);
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
                this.$refs.tasktable.updateData();
                this.close();
                break;
              }
              default:
                this.message = '服务器错误601';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch((e) => {
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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
          }).catch((e) => {
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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
          }).catch((e) => {
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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
                this.$refs.tasktable.updateData();
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
            console.dir(e);
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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
                this.$refs.tasktable.updateData();
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
          }).catch((e) => {
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
            this.message = '服务器错误778';
            this.alert_error = true;
            this.commit = false;
          });
        }
      }
    },
    reloadTasks(id) {
      if (!id && !this.ptasksearch) {
        return;
      }
      if (!id) {
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
        const objs = [{
          name: '无',
          id: '0',
        }];
        let url = `/task?filter=${this.ptasksearch}`;
        if (id) {
          url = `/task/${id}`;
        }

        this.$http.LGet(this.$store.state, url).then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                if (this.tasknames.indexOf(element.name) === -1) {
                  this.tasknames.push(element.name);
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
          resolve({
            objs,
          });
        }).catch((e) => {
          console.dir('服务器报错');
          if (e.response.data.code === 101) {
            this.$store.commit('logout');
            this.$router.replace('/login');
          }
          resolve({
            objs,
          });
        });
      }).then(() => {
        this.ptaskloading = false;
      });
    },
    getDesserts() {
      return this.datas;
    },
    allowedDates(v) {
      return this.$store.state.holidays.indexOf(v) === -1;
    },
    formatTimestamp(now) {
      return DateUtil.formatTimestamp(now);
    },
  },
};
</script>