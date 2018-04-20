<template>
  <v-container fluid>
      <!-- 添加编辑对话框 -->
      <v-dialog v-model="dialog" max-width="500px">
        <v-card>
          <v-card-title>
            <span class="headline">{{ formTitle }}</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout>
                <v-flex xs12 sm12 md12>
                  <v-date-picker
                    class="mt-3"
                    :landscape="true"
                    v-model="editedItem.day"
                    locale="zh-cn"
                    :first-day-of-week="1"
                    :allowed-dates="allowedDates"
                  ></v-date-picker>
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
          <td>{{ props.item.day }}</td>
          <td class="justify-center layout px-0">
            <v-btn v-if="editP" icon class="mx-0" @click="editItem(props.item)">
              <v-icon color="teal">edit</v-icon>
            </v-btn>
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

  export default {
    data: () => ({
      addP: false,
      editP: false,
      delP: false,
      dialog: false,
      delDialog: false,
      delItem: {},
      valid: true,
      message: '',
      headers: [
        {
          text: '假期',
          align: 'left',
          value: 'day',
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
      nameRules: [
        v => !!v || 'Day is required',
      ],
      editedItem: {
        id: '',
        day: '',
      },
      defaultItem: {
        id: '',
        day: '',
      },
    }),

    computed: {
      formTitle() {
        return this.editedIndex === -1 ? '新增' : '编辑';
      },
    },


    mounted() {
      const p = this.$store.state.permissions;
      this.addP = p['holiday.add'] !== undefined;
      this.editP = p['holiday.update'] !== undefined;
      this.delP = p['holiday.del'] !== undefined;
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
          day: '',
        };
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
        if (!this.editedItem.day) {
          return;
        }
        this.alert_error = false;
        if (this.editedIndex > -1) {
          this.commit = true;
          this.$http.LPost(this.$store.state, `/holiday/${this.editedItem.id}`, {
            day: Date.parse(new Date(this.editedItem.day)) / 1000,
          }).then((resp) => {
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
        } else {
          this.commit = true;
          this.$http.LPost(this.$store.state, '/holiday', {
            day: Date.parse(new Date(this.editedItem.day)) / 1000,
          }).then((resp) => {
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
      },
      del() {
        this.alert_error = false;
        this.commit = true;
        this.$http.LDel(this.$store.state, `/holiday/${this.delItem.id}`).then((resp) => {
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
        this.$store.commit('reloadholidays');
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
          this.$http.LGet(this.$store.state, `/holiday?${order}&offset=${offset}&limit=${rowsPerPage}${search}`).then((resp) => {
            this.loading = false;
            switch (resp.data.code) {
              case 0: {
                resp.data.data.data.forEach((element) => {
                  items.push({
                    id: element.ID,
                    day: this.formatDate(element.Day),
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
      formatDate(d) {
        if (d <= 0) {
          return '';
        }
        const now = new Date(d * 1000);
        const year = now.getFullYear();
        let month = now.getMonth() + 1;
        let date = now.getDate();
        if (month <= 9) {
          month = `0${month}`;
        }
        if (date <= 9) {
          date = `0${date}`;
        }
        return `${year}-${month}-${date}`;
      },
    },
  };
</script>