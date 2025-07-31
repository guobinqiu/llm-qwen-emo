<template>
  <div class="login-container">
    <div class="login-form">
      <div class="title-container">
        <h3 class="subtitle">用户登录</h3>
      </div>
      
      <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="form" autocomplete="on">
        <el-form-item prop="username">
          <el-input
            ref="username"
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            :key="passwordType"
            ref="password"
            v-model="loginForm.password"
            :type="passwordType"
            placeholder="密码"
            name="password"
            tabindex="2"
            autocomplete="on"
            @keyup.enter.native="handleLogin"
          />
        </el-form-item>
        
        <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">登录</el-button>
        
        <div class="tips">
          <span>没有账户？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { login } from '@/api/user'
import { setToken } from '@/utils/auth'
import { Message } from 'element-ui'

export default {
  name: 'UserLogin',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('密码不能少于6位'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur' }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  methods: {
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          login(this.loginForm)
            .then(response => {
              const { data } = response
              setToken(data.token)
              this.loading = false
              Message({
                message: '登录成功',
                type: 'success',
                duration: 2000
              })
              this.$router.push({ path: this.redirect || '/' })
            })
            .catch(error => {
              console.error('登录失败:', error)
              Message({
                message: '登录失败，请检查用户名和密码',
                type: 'error',
                duration: 3000
              })
              this.loading = false
            })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100%;
  width: 100%;
  background-color: #2d3a4b;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  position: relative;
  width: 520px;
  max-width: 100%;
  padding: 35px 35px 0;
  margin: 0 auto;
  overflow: hidden;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.title-container {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  font-size: 26px;
  color: #2d3a4b;
  margin: 0px auto 20px auto;
  text-align: center;
  font-weight: bold;
}

.subtitle {
  font-size: 18px;
  color: #8492a6;
  margin: 0px auto 40px auto;
  text-align: center;
}

.form {
  margin: 0 auto;
  width: 80%;
}

.tips {
  font-size: 14px;
  color: #8492a6;
  margin-bottom: 10px;
  text-align: center;
}

.tips a {
  color: #409EFF;
  text-decoration: none;
  margin-left: 5px;
}
</style>