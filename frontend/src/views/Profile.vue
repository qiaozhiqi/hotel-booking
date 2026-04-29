<template>
  <div class="profile-page">
    <div class="profile-container">
      <div class="user-card">
        <div class="user-avatar">
          <span class="avatar-icon">👤</span>
        </div>
        <div class="user-info">
          <h2 class="user-name">{{ userInfo?.username || '用户' }}</h2>
          <p class="user-email">{{ userInfo?.email || '未设置邮箱' }}</p>
        </div>
      </div>

      <div class="menu-section">
        <h3 class="section-title">我的服务</h3>
        <div class="menu-grid">
          <router-link to="/orders" class="menu-item">
            <div class="menu-icon order-icon">
              <span>📋</span>
            </div>
            <div class="menu-content">
              <span class="menu-title">我的订单</span>
              <span class="menu-desc" v-if="orderCount > 0">待处理 {{ orderCount }} 单</span>
              <span class="menu-desc" v-else>暂无订单</span>
            </div>
            <span class="menu-arrow">›</span>
          </router-link>

          <router-link to="/favorites" class="menu-item">
            <div class="menu-icon favorite-icon">
              <span>❤️</span>
            </div>
            <div class="menu-content">
              <span class="menu-title">我的收藏</span>
              <span class="menu-desc" v-if="favoriteCount > 0">已收藏 {{ favoriteCount }} 家</span>
              <span class="menu-desc" v-else>暂无收藏</span>
            </div>
            <span class="menu-arrow">›</span>
          </router-link>

          <router-link to="/invoices" class="menu-item">
            <div class="menu-icon invoice-icon">
              <span>📄</span>
            </div>
            <div class="menu-content">
              <span class="menu-title">发票管理</span>
              <span class="menu-desc" v-if="invoiceCount > 0">待开票 {{ pendingInvoiceCount }} 张</span>
              <span class="menu-desc" v-else>暂无发票</span>
            </div>
            <span class="menu-arrow">›</span>
          </router-link>
        </div>
      </div>

      <div class="quick-stats" v-if="showStats">
        <h3 class="section-title">统计概览</h3>
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-value">{{ totalOrders }}</div>
            <div class="stat-label">总订单数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ totalSpent }}</div>
            <div class="stat-label">累计消费</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ favoriteCount }}</div>
            <div class="stat-label">收藏酒店</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { userApi, orderApi, favoriteApi, invoiceApi } from '../api'

export default {
  name: 'Profile',
  setup() {
    const userInfo = ref(null)
    const orderCount = ref(0)
    const favoriteCount = ref(0)
    const invoiceCount = ref(0)
    const pendingInvoiceCount = ref(0)
    const totalOrders = ref(0)
    const totalSpent = ref('¥0')
    const showStats = ref(false)

    const loadUserInfo = async () => {
      try {
        const user = localStorage.getItem('user')
        if (user) {
          userInfo.value = JSON.parse(user)
        }
      } catch (error) {
        console.error('加载用户信息失败:', error)
      }
    }

    const loadOrderStats = async () => {
      try {
        const params = { page: 1, page_size: 100 }
        const res = await orderApi.getList(params)
        if (res.code === 200) {
          const orders = res.data.orders || []
          totalOrders.value = res.data.total || 0
          orderCount.value = orders.filter(o => 
            o.status === 'pending' || o.status === 'confirmed'
          ).length
          
          let total = 0
          orders.forEach(o => {
            total += o.total_amount || 0
          })
          totalSpent.value = `¥${total}`
          showStats.value = true
        }
      } catch (error) {
        console.error('加载订单统计失败:', error)
      }
    }

    const loadFavoriteCount = async () => {
      try {
        const params = { page: 1, page_size: 1 }
        const res = await favoriteApi.getList(params)
        if (res.code === 200) {
          favoriteCount.value = res.data.total || 0
        }
      } catch (error) {
        console.error('加载收藏数量失败:', error)
      }
    }

    const loadInvoiceCount = async () => {
      try {
        const params = { page: 1, page_size: 100 }
        const res = await invoiceApi.getList(params)
        if (res.code === 200) {
          const invoices = res.data.invoices || []
          invoiceCount.value = res.data.total || 0
          pendingInvoiceCount.value = invoices.filter(i => 
            i.status === 'pending'
          ).length
        }
      } catch (error) {
        console.error('加载发票数量失败:', error)
      }
    }

    onMounted(() => {
      loadUserInfo()
      loadOrderStats()
      loadFavoriteCount()
      loadInvoiceCount()
    })

    return {
      userInfo,
      orderCount,
      favoriteCount,
      invoiceCount,
      pendingInvoiceCount,
      totalOrders,
      totalSpent,
      showStats
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

.profile-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 20px;
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(26, 115, 232, 0.3);
}

.user-avatar {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-icon {
  font-size: 40px;
}

.user-info {
  color: #fff;
}

.user-name {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
}

.user-email {
  font-size: 14px;
  opacity: 0.8;
}

.menu-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.menu-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 10px;
  transition: all 0.2s;
}

.menu-item:hover {
  background: #e8f0fe;
  transform: translateX(4px);
}

.menu-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.menu-icon span {
  font-size: 24px;
}

.order-icon {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
}

.favorite-icon {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
}

.invoice-icon {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.menu-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.menu-desc {
  font-size: 12px;
  color: #999;
}

.menu-arrow {
  font-size: 20px;
  color: #ccc;
}

.quick-stats {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a73e8;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #666;
}

@media (max-width: 768px) {
  .user-card {
    flex-direction: column;
    text-align: center;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    text-align: left;
  }

  .stat-value {
    font-size: 22px;
    margin-bottom: 0;
  }
}
</style>
