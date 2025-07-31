<template>
  <div class="register-container">
    <div class="register-form">
      <div class="title-container">
        <h3 class="subtitle">用户注册</h3>
      </div>
      
      <el-form ref="registerForm" :model="registerForm" :rules="registerRules" class="form" autocomplete="on">
        <el-form-item prop="username">
          <el-input
            ref="username"
            v-model="registerForm.username"
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
            v-model="registerForm.password"
            :type="passwordType"
            placeholder="密码"
            name="password"
            tabindex="2"
            autocomplete="on"
          />
        </el-form-item>
        
        <el-form-item prop="confirmPassword">
          <el-input
            :key="confirmPasswordType"
            ref="confirmPassword"
            v-model="registerForm.confirmPassword"
            :type="confirmPasswordType"
            placeholder="确认密码"
            name="confirmPassword"
            tabindex="3"
            autocomplete="on"
            @keyup.enter.native="handleRegister"
          />
        </el-form-item>
        
        <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleRegister">注册</el-button>
        
        <div class="tips">
          <span>已有账户？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { register } from '@/api/user'
import { Message } from 'element-ui'

export default {
  name: 'UserRegister',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('密码不能少于6位'))
      } else {
        callback()
      }
    }
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== this.registerForm.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    return {
      registerForm: {
        username: '',
        password: '',
        confirmPassword: ''
      },
      registerRules: {
        username: [{ required: true, trigger: 'blur' }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }],
        confirmPassword: [{ required: true, trigger: 'blur', validator: validateConfirmPassword }]
      },
      loading: false,
      passwordType: 'password',
      confirmPasswordType: 'password'
    }
  },
  methods: {
    handleRegister() {
      this.$refs.registerForm.validate(valid => {
        if (valid) {
          this.loading = true
          const data = {
            username: this.registerForm.username,
            password: this.registerForm.password
          }
          register(data)
            .then(response => {
              const { message } = response
              this.loading = false
              Message({
                message: message || '注册成功',
                type: 'success',
                duration: 2000
              })
              this.$router.push('/login')
            })
            .catch(error => {
              console.error('注册失败:', error)
              Message({
                message: '注册失败，请稍后重试',
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
.register-container {
  min-height: 100%;
  width: 100%;
  background-color: #2d3a4b;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.register-form {
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