<template>
  <v-app>
    <v-navigation-drawer
      fixed
      v-model="drawer"
      app
      v-if="login"
    >
      <v-toolbar flat class="transparent elevation-3">
        <v-list class="pa-0">
          <v-list-tile avatar>
            <v-list-tile-avatar>
              <img src="https://randomuser.me/api/portraits/men/85.jpg" >
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title value="username">{{ this.$store.state.realname }}</v-list-tile-title>
            </v-list-tile-content>

              <v-tooltip bottom>
                <v-btn flat icon color="primary" @click="edit" slot="activator">
                  <v-icon dark>edit</v-icon>
                </v-btn>
                <span>编辑信息</span>
              </v-tooltip>
              <v-tooltip bottom>
                <v-btn flat icon color="pink" @click="logout" slot="activator">
                  <v-icon dark>exit_to_app</v-icon>
                </v-btn>
                <span>退出</span>
              </v-tooltip>
          </v-list-tile>
        </v-list>
      </v-toolbar>
      <v-divider></v-divider>
      <v-list>
        <v-list-tile
          value="true"
          v-for="(item, i) in items"
          :to="item.path"
          :key="i"
        >
          <v-list-tile-action>
            <v-icon v-html="item.icon"></v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title v-text="item.title"></v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar
      app
    >
      <v-toolbar-side-icon v-if="login" @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title v-text="title"></v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content>
      <router-view/>
    </v-content>
    <v-footer fixed app>
      <span>&copy; 2017</span>
    </v-footer>
  <!-- 添加编辑对话框 -->
    <v-dialog v-model="dialog" max-width="85%">
      <v-card>
        <v-card-title>
          <span class="headline">编辑用户信息</span>
        </v-card-title>
        <v-card-text>
          <v-container grid-list-md>
            <v-layout>
              <v-flex xs12 sm12 md12>
                <v-form v-model="valid" ref="form" lazy-validation>
                  <v-text-field
                  label="姓名"
                  v-model="realname"
                  :rules="realRules"
                  required
                  ></v-text-field>
                  <v-text-field
                  label="密码(空白不修改)"
                  v-model="password"
                  hint="最少6位"
                  :type="e1 ? 'password' : 'text'"
                  :rules="passwordRules"
                  min="6"
                  :append-icon="e1 ? 'visibility' : 'visibility_off'"
                  :append-icon-cb="() => (e1 = !e1)"
                  ></v-text-field>
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
  </v-app>
</template>

<script>
export default {

  data() {
    return {
      e1: true,
      dialog: false,
      drawer: true,
      valid: true,
      alert_error: false,
      commit: false,
      items: [],
      username: '',
      message: '',
      realname: '',
      password: '',
      realRules: [
        v => !!v || '姓名必填',
      ],
      passwordRules: [
        v => !v || v.length >= 6 || '密码最少6个',
      ],
      title: 'TaskPlan',
      login: this.$store.state.login,
    };
  },
  watch: {
    '$store.state.login': 'init',
  },
  created() {
    this.$store.commit('init');
  },
  methods: {
    init() {
      this.items = [];
      this.login = this.$store.state.login;
      this.username = this.$store.state.username;
      this.realname = this.$store.state.realname;
      if (this.$store.state.login) {
        this.$store.state.menu.forEach((element) => {
          for (const v of this.$store.state.constmenu) {
            if (v.path === element) {
              this.items.push(v);
              break;
            }
          }
        });
        this.$store.commit('reloadusers');
        this.$store.commit('reloadusergroups');
        this.$store.commit('reloadroles');
        this.$store.commit('reloadholidays');
        this.$router.replace('/');
      } else {
        this.$router.replace('/login');
      }
    },
    close() {
      this.password = '';
      this.dialog = false;
    },
    edit() {
      this.password = '';
      this.dialog = true;
    },
    save() {
      this.message = '';
      this.alert_error = false;
      if (this.$refs.form.validate()) {
        // Native form submission is not yet supported
        this.login = true;
        const d = {
          realname: this.realname,
        };
        if (!this.password) {
          d.password = this.password;
        }
        this.$http.Post(`/user/${this.$store.state.id}`, d).then((resp) => {
          switch (resp.data.code) {
            case 0: {
              setTimeout(() => {
                this.login = false;
                this.alert_success = true;
                this.$router.back();
              }, 2000);
              break;
            }
            case 300: {
              this.message = '用户名重复';
              this.alert_error = true;
              this.login = false;
              break;
            }
            default:
              this.alert_error = true;
              this.login = false;
          }
        }).catch(() => {
          this.alert = true;
          this.login = false;
        });
      }
    },
    logout() {
      this.$store.commit('logout');
    },
  },
  name: 'App',
};
</script>
