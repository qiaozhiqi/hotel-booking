<template>
  <div class="profile-page">
    <div class="profile-header">
      <div class="user-info">
        <div class="avatar">
          <span class="avatar-icon">👤</span>
        </div>
        <div class="user-details">
          <h2 class="user-name">{{ user?.username || '未登录' }}</h2>
          <p class="user-email" v-if="user?.email">{{ user.email }}</p>
          <p class="user-phone" v-if="user?.phone">{{ user.phone }}</p>
        </div>
      </div>
      <div class="user-stats">
        <div class="stat-item">
          <span class="stat-value">{{ stats.orders }}</span>
          <span class="stat-label">订单</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ stats.favorites }}</span>
          <span class="stat-label">收藏</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ stats.invoices }}</span>
          <span class="stat-label">发票</span>
        </div>
      </div>
    </div>

    <div class="profile-menu">
      <div class="menu-section">
        <h3 class="section-title">我的订单</h3>
        <div class="menu-items">
          <router-link to="/orders" class="menu-item">
            <span class="menu-icon">📋</span>
            <span class="menu-text">全部订单</span>
            <span class="menu-arrow">›</span>
          </router-link>
        </div>
      </div>

      <div class="menu-section">
        <h3 class="section-title">我的收藏</h3>
        <div class="menu-items">
          <router-link to="/favorites" class="menu-item">
            <span class="menu-icon">❤️</span>
            <span class="menu-text">收藏酒店</span>
            <span class="menu-arrow">›</span>
          </router-link>
        </div>
      </div>

      <div class="menu-section">
        <h3 class="section-title">发票管理</h3>
        <div class="menu-items">
          <router-link to="/invoices" class="menu-item">
            <span class="menu-icon">📄</span>
            <span class="menu-text">开发票</span>
            <span class="menu-arrow">›</span>
          </router-link>
        </div>
      </div>

      <div class="menu-section">
        <h3 class="section-title">账户设置</h3>
        <div class="menu-items">
          <div class="menu-item" @click="handleLogout" v-if="user">
            <span class="menu-icon">🚪</span>
            <span class="menu-text">退出登录</span>
            <span class="menu-arrow">›</span>
          </div>
          <router-link to="/login" class="menu-item" v-else>
            <span class="menu-icon">🔑</span>
            <span class="menu-text">登录账户</span>
            <span class="menu-arrow">›</span>
          </router-link>
        </div>
      </div>
    </div>

    <div v-if="!user" class="login-prompt">
      <div class="prompt-icon">🔒</div>
      <p class="prompt-text">您还未登录，请先登录以查看个人信息</p>
      <router-link to="/login" class="btn-login">去登录</router-link>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { orderApi, favoriteApi, invoiceApi, userApi } from '../api'

export default {
  name: 'Profile',
  setup() {
    const router = useRouter()
    const user = ref(null)
    const stats = ref({
      orders: 0,
      favorites: 0,
      invoices: 0
    })

    const loadUserInfo = () => {
      const userData = localStorage.getItem('user')
      if (userData) {
        user.value = JSON.parse(userData)
      }
    }

    const loadStats = async () => {
      if (!user.value) return

      try {
        const [ordersRes, favoritesRes, invoicesRes] = await Promise.all([
          orderApi.getList({ page: 1, page_size: 1 }),
          favoriteApi.getList({ page: 1, page_size: 1 }),
          invoiceApi.getList({ page: 1, page_size: 1 })
        ])

        stats.value.orders = ordersRes.data?.total || 0
        stats.value.favorites = favoritesRes.data?.total || 0
        stats.value.invoices = invoicesRes.data?.total || 0
      } catch (error) {
        console.error('加载统计数据失败:', error)
      }
    }

    const handleLogout = () => {
      localStorage.removeItem('user')
      user.value = null
      window.dispatchEvent(new Event('logout'))
      router.push('/')
    }

    onMounted(() => {
      loadUserInfo()
      loadStats()
    })

    return {
      user,
      stats,
      handleLogout
    }
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 30px 20px;
}

.profile-header {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-icon {
  font-size: 40px;
}

.user-details {
  color: #fff;
}

.user-name {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.user-email,
.user-phone {
  font-size: 14px;
  opacity: 0.8;
  margin: 4px 0;
}

.user-stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  text-align: center;
  color: #fff;
}

.stat-value {
  display: block;
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  opacity: 0.8;
}

.profile-menu {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.menu-section {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.menu-section:last-child {
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #999;
  padding: 0 24px;
  margin-bottom: 8px;
}

.menu-items {
  display: flex;
  flex-direction: column;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  cursor: pointer;
  transition: background 0.2s;
  text-decoration: none;
  color: #333;
}

.menu-item:hover {
  background: #f8f9fa;
}

.menu-icon {
  font-size: 24px;
  margin-right: 16px;
}

.menu-text {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
}

.menu-arrow {
  font-size: 18px;
  color: #ccc;
}

.login-prompt {
  display: none;
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 12px;
  margin-top: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.prompt-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.prompt-text {
  font-size: 15px;
  color: #666;
  margin-bottom: 24px;
}

.btn-login {
  display: inline-block;
  padding: 12px 40px;
  background: #1a73e8;
  color: #fff;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  transition: background 0.2s;
}

.btn-login:hover {
  background: #1557b0;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .user-info {
    flex-direction: column;
  }

  .user-stats {
    justify-content: center;
    gap: 24px;
  }

  .stat-value {
    font-size: 24px;
  }

  .login-prompt {
    display: block;
  }
}
</style>
