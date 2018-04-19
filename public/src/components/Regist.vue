<template>
  <v-app id="inspire">
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md8>
            <v-card class="elevation-12">
              <v-toolbar dark color="primary">
                <v-toolbar-title>注册</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form v-model="valid" ref="form" lazy-validation>
                    <v-text-field
                    label="用户名"
                    v-model="name"
                    :rules="nameRules"
                    required
                    ></v-text-field>
                    <v-text-field
                    label="姓名"
                    v-model="realname"
                    :rules="realRules"
                    required
                    ></v-text-field>
                    <v-text-field
                    label="密码"
                    v-model="password"
                    hint="最少6位"
                    :type="e1 ? 'password' : 'text'"
                    :rules="passwordRules"
                    min="6"
                    :append-icon="e1 ? 'visibility' : 'visibility_off'"
                    :append-icon-cb="() => (e1 = !e1)"
                    required
                    ></v-text-field>
                    <!--
                    <v-text-field
                    label="重复密码"
                    v-model="passwordt"
                    hint="最少6位"
                    :type="e2 ? 'password' : 'text'"
                    :rules="passwordtRules"
                    min="6"
                    :append-icon="e2 ? 'visibility' : 'visibility_off'"
                    :append-icon-cb="() => (e2 = !e2)"
                    required
                    ></v-text-field>
                    -->
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-btn color="primary" @click="back">返回</v-btn>
                <v-spacer></v-spacer>
                <v-btn 
                    color="success" 
                    @click="submit"
                    :disabled="!valid || login">注册</v-btn>
                <v-btn @click="clear">清空</v-btn>
              </v-card-actions>
              <v-alert
                  type="success"
                  :value="alert_success"
                  transition="scale-transition"
                >
                  注册成功
              </v-alert>
              <v-alert
                  type="error"
                  :value="alert_error"
                  transition="scale-transition"
              >
                  注册失败: {{ message }}
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
    e2: true,
    valid: true,
    login: false,
    alert_error: false,
    alert_success: false,
    message: '',
    name: '',
    realname: '',
    password: '',
    passwordt: '',
    nameRules: [
      v => !!v || '用户名必填',
    ],
    realRules: [
      v => !!v || '姓名必填',
    ],
    passwordRules: [
      v => !!v || '密码必填',
      v => !v || v.length >= 6 || '密码最少6个',
    ],
    passwordtRules: [
      v => !!v || '重复密码必填',
      v => !v || v.length >= 6 || '密码最少6个',
    //   (v) => {
    //     console.dir(v);
    //     console.dir(this.password);
    //     console.dir(this);
    //     console.dir(v === this.password);
    //     return !v || v === this.password || '2次密码不一样';
    //   },
    ],
  }),

  methods: {
    submit() {
      this.message = '';
      this.alert_success = false;
      this.alert_error = false;
      if (this.$refs.form.validate()) {
        // Native form submission is not yet supported
        this.login = true;
        this.$http.Post('/user', {
          name: this.name,
          realname: this.realname,
          password: this.password,
        }).then((resp) => {
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
    back() {
      this.$router.back();
    },
    clear() {
      this.alert_success = false;
      this.alert_error = false;
      this.$refs.form.reset();
    },
  },
};
</script>