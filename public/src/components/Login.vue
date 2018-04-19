<template>
  <v-app id="inspire">
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md8>
            <v-card class="elevation-12">
              <v-toolbar dark color="primary">
                <v-toolbar-title>Login</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form v-model="valid" ref="form" lazy-validation>
                    <v-text-field
                    label="User"
                    v-model="user"
                    :rules="userRules"
                    required
                    ></v-text-field>
                    <v-text-field
                    label="Password"
                    v-model="password"
                    hint="At least 6 characters"
                    :type="e1 ? 'password' : 'text'"
                    :rules="passwordRules"
                    min="6"
                    :append-icon="e1 ? 'visibility' : 'visibility_off'"
                    :append-icon-cb="() => (e1 = !e1)"
                    required
                    ></v-text-field>
                    <v-checkbox
                      color="green"
                      v-model="autologin"
                    >
                    <div slot="label">
                    记住账号密码
                    </div>
                    </v-checkbox>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-btn color="primary" @click="regist">注册</v-btn>
                <v-spacer></v-spacer>
                <v-btn 
                    color="success" 
                    @click="submit"
                    :disabled="!valid || login">登录</v-btn>
                <v-btn @click="clear">清空</v-btn>
              </v-card-actions>
              <v-alert
                  type="success"
                  :value="alert_success"
                  transition="scale-transition"
                >
                  登录成功
              </v-alert>
              <v-alert
                  type="error"
                  :value="alert_error"
                  transition="scale-transition"
              >
                  登录失败
              </v-alert>
              <v-progress-linear :indeterminate="true" :active="login"></v-progress-linear>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
  </v-app>
</template>

<script>

export default {
  data: () => ({
    e1: true,
    valid: true,
    login: false,
    user: '',
    alert_error: false,
    alert_success: false,
    autologin: true,
    userRules: [
      v => !!v || 'User is required',
    ],
    password: '',
    passwordRules: [
      v => !!v || 'Password is required',
      v => !v || v.length >= 6 || 'Password must be more than 6 characters',
    ],
  }),

  mounted() {
    this.user = localStorage.getItem('user');
    this.password = localStorage.getItem('password');
    this.autologin = localStorage.getItem('autologin');
  },

  methods: {
    submit() {
      this.alert_success = false;
      this.alert_error = false;
      if (this.$refs.form.validate()) {
        // Native form submission is not yet supported
        this.login = true;
        this.$http.Post('/login', {
          user: this.user,
          password: this.password,
        }).then((resp) => {
          switch (resp.data.code) {
            case 1000: {
              setTimeout(() => {
                if (this.autologin) {
                  localStorage.setItem('user', this.user);
                  localStorage.setItem('password', this.password);
                  localStorage.setItem('autologin', this.autologin);
                } else {
                  localStorage.removeItem('user');
                  localStorage.removeItem('password');
                  localStorage.removeItem('autologin');
                }
                this.login = false;
                this.alert_success = true;
                const d = resp.data.data;
                this.$store.commit('login', d);
              }, 1000);
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
    regist() {
      this.$router.push('/regist');
    },
    clear() {
      this.alert_success = false;
      this.alert_error = false;
      this.$refs.form.reset();
    },
  },
};
</script>