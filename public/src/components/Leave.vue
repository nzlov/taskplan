<template>
  <v-container fluid>
      <!-- 添加编辑对话框 -->
      <v-dialog v-model="dialog" max-width="85%">
        <v-card>
          <v-card-title>
            <span class="headline">{{ formTitle }}</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout>
                <v-flex xs12 sm12 md12>
                  <v-form v-model="valid" ref="form" lazy-validation>
                    <v-select
                      label="请假人"
                      :items="users"
                      v-model="editedItem.userid"
                      item-text="name"
                      item-value="id"
                      max-height="600px"
                      hint="选择请假人"
                      :rules="useridRules"
                      required
                    ></v-select>
                  </v-form>
                  <v-card>
                    <v-card-title>
                      <span class="headline">开始时间</span>
                    </v-card-title>
                    <v-card-text>
                      <v-date-picker locale="zh-cn" :first-day-of-week="1" :allowed-dates="allowedDates" v-model="startdate"></v-date-picker>
                      <v-time-picker locale="zh-cn" v-model="editedItem.starttime" format="24hr"></v-time-picker>
                    </v-card-text>
                  </v-card>
                  <v-card>
                    <v-card-title>
                      <span class="headline">结束时间</span>
                    </v-card-title>
                    <v-card-text>
                      <v-date-picker locale="zh-cn" :first-day-of-week="1" :allowed-dates="allowedDates" v-model="enddate"></v-date-picker>
                      <v-time-picker locale="zh-cn" v-model="editedItem.endtime" format="24hr"></v-time-picker>
                    </v-card-text>
                  </v-card>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat @click.native="close">Cancel</v-btn>
            <v-btn color="blue darken-1" flat @click.native="save" :disabled="!valid || commit">Save</v-btn>
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
      <!-- 删除对话框 -->
      <v-dialog v-model="delDialog" persistent max-width="500px">
        <v-card>
          <v-card-title>
            <span> 真的要删除 {{ delItem }} </span>
          </v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="error" flat @click="del" :disabled="commit">Yes</v-btn>
            <v-btn color="primary" flat @click.stop="delDialog=false;commit=false" :disabled="commit">No</v-btn>
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
          <td>{{ props.item.realname }}</td>
          <td class="text-xs-center">{{ props.item.starts }}</td>
          <td class="text-xs-center">{{ props.item.ends }}</td>
          <td class="justify-center layout px-0">
            <v-btn v-if="delP" icon class="mx-0" @click="deleteItem(props.item)">
              <v-icon color="pink">delete</v-icon>
            </v-btn>
          </td>
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
  data: () => ({
    addP: false,
    delP: false,
    dialog: false,
    delDialog: false,
    delItem: {},
    valid: true,
    message: '',
    headers: [
      {
        text: '请假人',
        align: 'left',
        value: 'user_id',
      },
      {
        text: '开始时间',
        align: 'center',
        value: 'start',
        width: '300px',
      },
      {
        text: '结束时间',
        align: 'center',
        value: 'end',
        width: '300px',
      },
      { text: 'Actions',
        align: 'center',
        value: 'name',
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
    useridRules: [
      v => !!v || '请假人 is required',
    ],
    users: [],
    startdate: '2016-01-01',
    enddate: '2016-01-01',
    editedItem: {
      id: '',
      userid: '',
      starttime: '00:00',
      endtime: '23:59',
    },
    defaultItem: {
      id: '',
      userid: '',
      starttime: '00:00',
      endtime: '23:59',
    },
  }),

  computed: {
    formTitle() {
      return this.editedIndex === -1 ? '新增' : '编辑';
    },
  },


  mounted() {
    const p = this.$store.state.permissions;
    this.addP = p['leave.add'] !== undefined;
    this.delP = p['leave.del'] !== undefined;

    this.users = this.$store.state.users;
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
      };
      this.startdate = DateUtil.formatDate1(new Date());
      this.enddate = DateUtil.formatDate1(new Date());
      this.dialog = true;
      this.valid = true;
    },

    editItem(item) {
      this.editedIndex = this.items.indexOf(item);
      this.editedItem = Object.assign({}, item);
      console.dir(item);
      this.dialog = true;
    },

    deleteItem(item) {
      this.delItem = item;
      this.delDialog = true;
    },

    close() {
      this.dialog = false;
      this.delDialog = false;
      setTimeout(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.delItem = Object.assign({}, this.defaultItem);
        this.sp = [];
        this.sm = [];
        this.editedIndex = -1;
      }, 300);
    },

    save() {
      this.alert_error = false;
      if (this.editedIndex === -1) {
        if (this.$refs.form.validate()) {
          const d = {
            userid: this.editedItem.userid,
            start: '',
            end: '',
          };
          if (this.startdate && this.editedItem.starttime) {
            d.start = this.formatTimestamp(`${this.startdate} ${this.editedItem.starttime}`);
          }
          if (this.enddate && this.editedItem.endtime) {
            d.end = this.formatTimestamp(`${this.enddate} ${this.editedItem.endtime}`);
          }
          if (!d.userid) {
            console.dir(d.userid);
            return;
          }
          if (!d.start) {
            this.alert_error = true;
            this.message = '请选择开始时间';
            console.dir(d.start);
            return;
          }
          if (!d.end) {
            this.alert_error = true;
            this.message = '请选择结束时间';
            console.dir(d.end);
            return;
          }
          this.commit = true;
          this.$http.LPost(this.$store.state, '/leave', d).then((resp) => {
            this.commit = false;
            switch (resp.data.code) {
              case 0: {
                this.updateData();
                this.close();
                break;
              }
              case 300: {
                this.message = '日期重复';
                this.alert_error = true;
                break;
              }
              default:
                this.message = '服务器错误';
                this.alert_error = true;
            }
            this.commit = false;
          }).catch(() => {
            this.message = '服务器错误';
            this.alert_error = true;
            this.commit = false;
          });
        }
      }
    },
    del() {
      this.alert_error = false;
      this.commit = true;
      this.$http.LDel(this.$store.state, `/leave/${this.delItem.id}`).then((resp) => {
        this.commit = false;
        switch (resp.data.code) {
          case 0: {
            this.updateData();
            this.close();
            break;
          }
          default:
            this.message = '服务器错误';
            this.alert_error = true;
        }
        this.commit = false;
      }).catch(() => {
        this.message = '服务器错误';
        this.alert_error = true;
        this.commit = false;
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
        this.$http.LGet(this.$store.state, `/leave?${order}&offset=${offset}&limit=${rowsPerPage}${search}`).then((resp) => {
          this.loading = false;
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                items.push({
                  id: element.ID,
                  userid: element.UserID,
                  realname: element.User.RealName,
                  start: element.Start,
                  starts: DateUtil.formatDate(element.Start),
                  end: element.End,
                  ends: DateUtil.formatDate(element.End),
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
    formatTimestamp(now) {
      return DateUtil.formatTimestamp(now);
    },
  },
};
</script>