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
                    <v-text-field
                    label="Name"
                    v-model="editedItem.name"
                    :rules="nameRules"
                    required
                    ></v-text-field>
                    <v-layout row wrap>
                      <v-flex xs12 sm12 md12>
                        <v-select
                          label="权限"
                          :items="permissions"
                          v-model="editedItem.permission"
                          multiple
                          item-text="name"
                          item-value="tag"
                          max-height="400"
                          hint="选择权限"
                          persistent-hint
                        ></v-select>
                      </v-flex>
                    </v-layout>
                    <v-layout row wrap>
                      <v-flex xs12 sm12 md12>
                        <v-select
                          label="菜单"
                          :items="menu"
                          v-model="editedItem.menu"
                          multiple
                          item-text="title"
                          item-value="path"
                          max-height="auto"
                          hint="选择菜单"
                          persistent-hint
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
          <td>{{ props.item.name }}</td>
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
          text: '角色',
          align: 'left',
          value: 'name',
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
        v => !!v || 'Name is required',
      ],
      editedItem: {
        id: '',
        name: '',
        permission: [],
        menu: [],
      },
      defaultItem: {
        id: '',
        name: '',
        permission: [],
        menu: [],
      },
      permissions: [],
      menu: [],
    }),

    computed: {
      formTitle() {
        return this.editedIndex === -1 ? '新增' : '编辑';
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
    },

    mounted() {
      const p = this.$store.state.permissions;
      this.addP = p['role.add'] !== undefined;
      this.editP = p['role.update'] !== undefined;
      this.delP = p['role.del'] !== undefined;

      this.menu = this.$store.state.constmenu;
      new Promise((resolve) => {
        const permissions = [];
        this.$http.LGet(this.$store.state, '/permission').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                permissions.push({
                  name: element.Name,
                  tag: element.Tag,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            permissions,
          });
        }).catch(() => {
          console.dir('服务器报错');
          resolve({
            permissions,
          });
        });
      }).then((data) => {
        this.permissions = data.permissions;
      });
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
        console.dir(this.editedItem);
        this.alert_error = false;
        if (this.editedIndex > -1) {
          if (this.$refs.form.validate()) {
            this.commit = true;
            this.$http.LPost(this.$store.state, `/role/${this.editedItem.id}`, this.editedItem).then((resp) => {
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
        } else {
          if (this.$refs.form.validate()) {
            this.commit = true;
            this.$http.LPost(this.$store.state, '/role', this.editedItem).then((resp) => {
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
      del() {
        this.alert_error = false;
        this.commit = true;
        this.$http.LDel(this.$store.state, `/role/${this.delItem.id}`).then((resp) => {
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
      },
      updateData() {
        this.$store.commit('reloadroles');
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
          this.$http.LGet(this.$store.state, `/role?${order}&offset=${offset}&limit=${rowsPerPage}${search}`).then((resp) => {
            this.loading = false;
            switch (resp.data.code) {
              case 0: {
                resp.data.data.data.forEach((element) => {
                  items.push({
                    id: element.ID,
                    name: element.Name,
                    permission: element.Permission,
                    menu: element.Menu,
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
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
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