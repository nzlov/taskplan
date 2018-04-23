<template>
  <v-container fluid>

    <v-snackbar
      :timeout="5000"
      :top="true"
      v-model="snackbar"
    >
      {{ text }}
      <v-btn flat color="pink" @click.native="snackbar = false">Close</v-btn>
    </v-snackbar>
    <!-- 添加编辑对话框 -->
      <v-dialog v-model="dialog" max-width="85%">
        <v-card>
          <v-card-title>
            <span class="headline">编辑</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout>
                <v-flex xs12 sm12 md12>
                  <v-form v-model="valid" ref="form" lazy-validation>
                    <v-text-field
                    label="名字"
                    v-model="editedItem.realname"
                    :rules="realnameRules"
                    required
                    ></v-text-field>
                    <v-layout>
                      <v-flex xs12 sm12 md12>
                        <v-select
                          label="角色"
                          :items="roles"
                          v-model="editedItem.roleid"
                          item-text="name"
                          item-value="id"
                          max-height="600px"
                          hint="选择角色"
                        ></v-select>
                      </v-flex>
                    </v-layout>
                    <v-layout>
                      <v-flex xs12 sm12 md12>
                        <v-select
                          label="资源组"
                          :items="usergroups"
                          v-model="editedItem.usergroupid"
                          item-text="name"
                          item-value="id"
                          max-height="600px"
                          hint="选择资源组"
                        ></v-select>
                      </v-flex>
                    </v-layout>
                </v-form>
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
        <td>{{ props.item.name }}</td>
        <td class="text-xs-center">{{ props.item.realname }}</td>
        <td class="justify-center layout px-0">
          <v-tooltip bottom v-if="editP">
            <v-btn icon slot="activator" class="mx-0" @click="editItem(props.item)">
              <v-icon color="teal">edit</v-icon>
            </v-btn>
            <span>编辑</span>
          </v-tooltip>
          <v-tooltip bottom v-if="changeP">
            <v-btn icon slot="activator" class="mx-0" @click="changeItem(props.item)" >
              <v-icon color="pink">{{ props.item.status==0 ? 'block' :'check_circle' }}</v-icon>
            </v-btn>
            <span>{{ props.item.status==0 ? '禁用' :'启用' }}</span>
          </v-tooltip>
          <v-tooltip bottom v-if="resetP">
            <v-btn icon slot="activator" class="mx-0" @click="resetPasswordItem(props.item)">
              <v-icon color="pink">build</v-icon>
            </v-btn>
            <span>重置密码</span>
          </v-tooltip>
          <v-tooltip bottom v-if="delP">
            <v-btn icon slot="activator" class="mx-0" @click="deleteItem(props.item)">
              <v-icon color="pink">delete</v-icon>
            </v-btn>
            <span>删除</span>
          </v-tooltip>
        </td>
      </template>
      <v-alert slot="no-data" :value="true" color="error" icon="warning">
        没有数据
      </v-alert>
      </v-data-table>

    </v-card>
  </v-container>
</template>
<script>

  export default {
    data: () => ({
      dialog: false,
      actionDialog: false,
      snackbar: false,
      text: '',
      message: '',
      editP: false,
      changeP: false,
      resetP: false,
      delP: false,
      action: 0, // 操作类型 1 删除 2 启用禁用 3 重置密码
      actionMessage: '',
      actionItem: {},
      valid: true,
      headers: [
        {
          text: '用户名',
          align: 'left',
          value: 'name',
        },
        {
          text: '姓名',
          align: 'center',
          value: 'realname',
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
      realnameRules: [
        v => !!v || '姓名 is required',
      ],
      editedItem: {
        id: '',
        name: '',
        realname: '',
        roleid: 0,
        usergroupid: 0,
      },
      defaultItem: {
        id: '',
        name: '',
        realname: '',
        roleid: 0,
        usergroupid: 0,
      },
      roles: [],
      usergroups: [],
    }),
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
    mounted() {
      this.editP = Object.prototype.hasOwnProperty.call(this.$store.state.permissions, 'user.update');
      this.changeP = Object.prototype.hasOwnProperty.call(this.$store.state.permissions, 'user.change');
      this.resetP = Object.prototype.hasOwnProperty.call(this.$store.state.permissions, 'user.resetpassword');
      this.delP = Object.prototype.hasOwnProperty.call(this.$store.state.permissions, 'user.del');

      this.roles = this.$store.state.roles;
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
          name: '',
          permission: [],
          menu: [],
        };
        this.dialog = true;
        this.valid = true;
      },

      editItem(item) {
        this.editedIndex = this.items.indexOf(item);
        this.editedItem = Object.assign({}, item);
        this.dialog = true;
      },

      deleteItem(item) {
        this.action = 1;
        this.actionItem = item;
        this.actionDialog = true;
        this.actionMessage = `真的要删除${item.name}?`;
      },

      changeItem(item) {
        this.action = 2;
        this.actionItem = item;
        this.actionDialog = true;
        this.actionMessage = `真的要${item.status === 0 ? '禁用' : '启用'}${item.name}?`;
      },

      resetPasswordItem(item) {
        this.action = 3;
        this.actionItem = item;
        this.actionDialog = true;
        this.actionMessage = `真的要重置${item.name}密码?`;
      },

      close() {
        this.dialog = false;
        this.actionDialog = false;
        setTimeout(() => {
          this.editedItem = Object.assign({}, this.defaultItem);
          this.actionItem = Object.assign({}, this.defaultItem);
          this.editedIndex = -1;
        }, 300);
      },
      actione() {
        this.alert_error = false;
        this.commit = true;
        switch (this.action) {
          case 1: {
            this.$http.LDel(this.$store.state, `/user/${this.actionItem.id}`).then((resp) => {
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
            }).catch((e) => {
              if (e.response.data.code === 101) {
                this.$store.commit('logout');
                this.$router.replace('/login');
              }
              this.message = '服务器错误';
              this.alert_error = true;
              this.commit = false;
            });
            break;
          }
          case 2: {
            this.$http.LPost(this.$store.state, `/user/${this.actionItem.id}/change`, { status: this.actionItem.status === 0 ? 1 : 0 }).then((resp) => {
              this.commit = false;
              switch (resp.data.code) {
                case 0: {
                  this.actionItem.status = this.actionItem.status === 0 ? 1 : 0;
                  this.close();
                  break;
                }
                default:
                  this.message = '服务器错误';
                  this.alert_error = true;
              }
              this.commit = false;
            }).catch((e) => {
              if (e.response.data.code === 101) {
                this.$store.commit('logout');
                this.$router.replace('/login');
              }
              this.message = '服务器错误';
              this.alert_error = true;
              this.commit = false;
            });
            break;
          }
          case 3: {
            this.$http.LPost(this.$store.state, `/user/${this.actionItem.id}/resetpassword`).then((resp) => {
              this.commit = false;
              switch (resp.data.code) {
                case 0: {
                  this.text = `${this.actionItem.name}重置密码为：${resp.data.data}`;
                  this.snackbar = true;
                  this.close();
                  break;
                }
                default:
                  this.message = '服务器错误';
                  this.alert_errosr = true;
              }
              this.commit = false;
            }).catch((e) => {
              if (e.response.data.code === 101) {
                this.$store.commit('logout');
                this.$router.replace('/login');
              }
              this.message = '服务器错误';
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
        console.dir(this.editedItem);
        this.alert_error = false;
        if (this.editedIndex > -1) {
          if (this.$refs.form.validate()) {
            console.dir(this.editedItem);
            this.commit = true;
            this.$http.LPost(this.$store.state, `/user/${this.editedItem.id}`, this.editedItem).then((resp) => {
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
                  this.message = '服务器错误';
                  this.alert_error = true;
              }
              this.commit = false;
            }).catch((e) => {
              if (e.response.data.code === 101) {
                this.$store.commit('logout');
                this.$router.replace('/login');
              }
              this.message = '服务器错误';
              this.alert_error = true;
              this.commit = false;
            });
          }
        }
      },
      updateData() {
        this.$store.commit('reloadusers');
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
          this.$http.LGet(this.$store.state, `/user?${order}&offset=${offset}&limit=${rowsPerPage}${search}`).then((resp) => {
            this.loading = false;
            switch (resp.data.code) {
              case 0: {
                resp.data.data.data.forEach((element) => {
                  items.push({
                    id: element.ID,
                    name: element.Name,
                    realname: element.RealName,
                    status: element.Status,
                    roleid: element.RoleID,
                    usergroupid: element.UserGroupID,
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
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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
    },
  };
</script>