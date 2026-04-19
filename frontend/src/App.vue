<template>
  <div class="app">
    <header class="header">
      <div class="header-content">
        <div class="logo" @click="goHome">
          <span class="logo-icon">🏨</span>
          <span class="logo-text">酒店预订</span>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link" :class="{ active: $route.path === '/' }">首页</router-link>
          <router-link to="/orders" class="nav-link" :class="{ active: $route.path === '/orders' }">我的订单</router-link>
        </nav>
        <div class="user-area">
          <template v-if="user">
            <span class="user-name">{{ user.username }}</span>
            <button class="btn-logout" @click="logout">退出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="btn-login">登录</router-link>
          </template>
        </div>
      </div>
    </header>
    <main class="main">
      <router-view />
    </main>
    <footer class="footer">
      <p>&copy; 2026 酒店预订系统. 版权所有.</p>
    </footer>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'App',
  setup() {
    const router = useRouter()
    const user = ref(null)

    const checkLogin = () => {
      const userData = localStorage.getItem('user')
      if (userData) {
        user.value = JSON.parse(userData)
      }
    }

    const logout = () => {
      localStorage.removeItem('user')
      user.value = null
      router.push('/')
    }

    const goHome = () => {
      router.push('/')
    }

    onMounted(() => {
      checkLogin()
      window.addEventListener('login-success', checkLogin)
      window.addEventListener('logout', checkLogin)
    })

    return {
      user,
      logout,
      goHome
    }
  }
}
</script>

<style>
.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.logo-icon {
  font-size: 28px;
  margin-right: 8px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #1a73e8;
}

.nav {
  display: flex;
  gap: 30px;
}

.nav-link {
  font-size: 16px;
  color: #666;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: #1a73e8;
}

.user-area {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-name {
  color: #333;
  font-size: 14px;
}

.btn-login {
  padding: 8px 20px;
  background: #1a73e8;
  color: #fff;
  border-radius: 4px;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-login:hover {
  background: #1557b0;
}

.btn-logout {
  padding: 6px 16px;
  background: transparent;
  border: 1px solid #ddd;
  color: #666;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-logout:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.main {
  min-height: calc(100vh - 120px);
}

.footer {
  background: #f8f9fa;
  padding: 20px 0;
  text-align: center;
  color: #999;
  font-size: 14px;
}
</style>
