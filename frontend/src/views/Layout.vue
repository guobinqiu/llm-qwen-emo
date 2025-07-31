<template>
  <div class="layout-root">
    <!-- 顶部导航栏 -->
    <div class="navbar">
      <div class="logo">EMO数字人视频生成系统</div>
      <div class="right-menu">
        <el-dropdown class="avatar-container" trigger="click">
          <div class="avatar-wrapper">
            Welcome, <span class="user-name">{{ userInfo.username }}</span>
            <i class="el-icon-caret-bottom" />
          </div>
          <el-dropdown-menu slot="dropdown" class="user-dropdown">
            <el-dropdown-item @click.native="logout">
              <span style="display:block;">退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </div>
    <!-- 主体区域 -->
    <div class="layout-main">
      <!-- 左侧菜单 -->
      <div class="layout-left">
        <el-menu
          :default-active="activeIndex"
          class="el-menu-vertical"
          @select="handleSelect"
          background-color="#545c64"
          text-color="#fff"
          active-text-color="#ffd04b">
          <el-menu-item index="/images">
            <i class="el-icon-picture"></i>
            <span slot="title">图像管理</span>
          </el-menu-item>
          <el-menu-item index="/audios">
            <i class="el-icon-headset"></i>
            <span slot="title">音频管理</span>
          </el-menu-item>
          <el-menu-item index="/tasks">
            <i class="el-icon-s-order"></i>
            <span slot="title">任务管理</span>
          </el-menu-item>
        </el-menu>
      </div>
      <!-- 右侧内容 -->
      <div class="layout-right">
        <router-view :key="$route.path" />
      </div>
    </div>
  </div>
</template>

<script>
import { removeToken } from '@/utils/auth'
import { Message } from 'element-ui'
import { getUserInfo } from '@/api/user'

export default {
  name: 'LayoutView',
  data() {
    return {
      activeIndex: this.$route.path,
      userInfo: {}
    }
  },
  created() {
    this.fetchUserInfo()
  },
  watch: {
    $route(to) {
      this.activeIndex = to.path
    }
  },
  methods: {
    handleSelect(key) {
      this.$router.push(key)
    },
    async fetchUserInfo() {
      try {
        const res = await getUserInfo()
        this.userInfo = res.data.user || {}
      } catch (e) {
        this.userInfo = {}
      }
    },
    async logout() {
      this.$confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeToken()
        Message({
          message: '退出成功',
          type: 'success',
          duration: 2000
        })
        this.$router.push(`/login?redirect=${this.$route.fullPath}`)
      })
    }
  }
}
</script>

<style scoped>
.layout-root {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.navbar {
  height: 60px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  position: relative;
  z-index: 100;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #409EFF;
}

.right-menu {
  display: flex;
  align-items: center;
}

.avatar-container {
  margin-right: 30px;
}

.avatar-wrapper {
  position: relative;
  cursor: pointer;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
}

.el-icon-caret-bottom {
  cursor: pointer;
  right: -20px;
  top: 25px;
  font-size: 12px;
}

.layout-main {
  flex: 1;
  display: flex;
  height: calc(100vh - 60px);
  overflow: hidden;
}

.layout-left {
  width: 200px;
  background: #545c64;
  height: 100%;
  min-height: 100%;
}

.layout-right {
  flex: 1;
  overflow: auto;
  background-color: #f0f2f5;
  padding: 20px;
  height: 100%;
}

.user-name {
  font-weight: bold;
  margin-right: 8px;
  color: #409EFF;
}
</style>