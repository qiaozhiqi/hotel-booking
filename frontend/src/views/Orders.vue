<template>
  <div class="orders-page">
    <div class="orders-header">
      <h1 class="page-title">我的订单</h1>
      <div class="status-tabs">
        <button 
          v-for="tab in statusTabs" 
          :key="tab.value"
          class="status-tab"
          :class="{ active: selectedStatus === tab.value }"
          @click="selectStatus(tab.value)"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <template v-else-if="orders.length > 0">
      <div class="orders-list">
        <div 
          v-for="order in orders" 
          :key="order.id" 
          class="order-card"
        >
          <div class="order-header">
            <div class="order-info">
              <span class="order-no">订单号：{{ order.order_no }}</span>
              <span class="order-date">{{ formatDate(order.created_at) }}</span>
            </div>
            <span class="order-status" :class="getStatusClass(order.status)">
              {{ getStatusText(order.status) }}
            </span>
          </div>

          <div class="order-content">
            <div class="hotel-info">
              <h3 class="hotel-name">{{ order.hotel_name }}</h3>
              <p class="room-name">{{ order.room_name }}</p>
              <div class="stay-dates">
                <span class="date-item">
                  <span class="date-label">入住</span>
                  <span class="date-value">{{ order.check_in }}</span>
                </span>
                <span class="date-arrow">→</span>
                <span class="date-item">
                  <span class="date-label">离店</span>
                  <span class="date-value">{{ order.check_out }}</span>
                </span>
              </div>
              <div class="guest-info">
                <span class="guest-item">入住人：{{ order.guest_name }}</span>
                <span class="guest-item">电话：{{ order.guest_phone }}</span>
              </div>
            </div>
            <div class="price-actions">
              <div class="price-info">
                <span class="price-label">订单金额</span>
                <span class="price-value">¥{{ order.total_amount }}</span>
              </div>
              <div class="order-actions">
                <button 
                  v-if="order.status === 'confirmed' || order.status === 'pending'"
                  class="btn-cancel-order"
                  @click="cancelOrder(order.id)"
                >
                  取消订单
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="pagination" v-if="total > 0">
        <button 
          class="page-btn" 
          :disabled="page <= 1" 
          @click="changePage(page - 1)"
        >
          上一页
        </button>
        <span class="page-info">第 {{ page }} 页 / 共 {{ totalPages }} 页</span>
        <button 
          class="page-btn" 
          :disabled="page >= totalPages" 
          @click="changePage(page + 1)"
        >
          下一页
        </button>
      </div>
    </template>

    <div v-else class="empty-state">
      <div class="empty-icon">📋</div>
      <p class="empty-text">暂无订单记录</p>
      <router-link to="/" class="btn-go-book">去预订酒店</router-link>
    </div>

    <div v-if="showCancelModal" class="modal-overlay" @click.self="closeCancelModal">
      <div class="modal-content confirm-modal">
        <div class="modal-icon">⚠️</div>
        <h3 class="modal-title">确认取消订单？</h3>
        <p class="modal-text">取消后将无法恢复，您确定要取消该订单吗？</p>
        <div class="modal-actions">
          <button class="btn-secondary" @click="closeCancelModal">再想想</button>
          <button class="btn-danger" @click="confirmCancel" :disabled="cancelling">
            <span v-if="cancelling">取消中...</span>
            <span v-else>确认取消</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { orderApi } from '../api'

export default {
  name: 'Orders',
  setup() {
    const loading = ref(true)
    const orders = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(10)
    const selectedStatus = ref('')

    const showCancelModal = ref(false)
    const cancellingOrderId = ref(null)
    const cancelling = ref(false)

    const statusTabs = [
      { label: '全部', value: '' },
      { label: '待确认', value: 'pending' },
      { label: '已确认', value: 'confirmed' },
      { label: '已入住', value: 'checked_in' },
      { label: '已完成', value: 'checked_out' },
      { label: '已取消', value: 'cancelled' }
    ]

    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

    const getStatusText = (status) => {
      const map = {
        'pending': '待确认',
        'confirmed': '已确认',
        'checked_in': '已入住',
        'checked_out': '已完成',
        'cancelled': '已取消'
      }
      return map[status] || status
    }

    const getStatusClass = (status) => {
      const map = {
        'pending': 'status-pending',
        'confirmed': 'status-confirmed',
        'checked_in': 'status-checked',
        'checked_out': 'status-done',
        'cancelled': 'status-cancelled'
      }
      return map[status] || ''
    }

    const formatDate = (dateStr) => {
      const date = new Date(dateStr)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hour = String(date.getHours()).padStart(2, '0')
      const minute = String(date.getMinutes()).padStart(2, '0')
      return `${year}-${month}-${day} ${hour}:${minute}`
    }

    const loadOrders = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        if (selectedStatus.value) {
          params.status = selectedStatus.value
        }

        const res = await orderApi.getList(params)
        if (res.code === 200) {
          orders.value = res.data.orders || []
          total.value = res.data.total || 0
        }
      } catch (error) {
        console.error('加载订单列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    const selectStatus = (status) => {
      selectedStatus.value = status
      page.value = 1
      loadOrders()
    }

    const changePage = (newPage) => {
      page.value = newPage
      loadOrders()
    }

    const cancelOrder = (orderId) => {
      cancellingOrderId.value = orderId
      showCancelModal.value = true
    }

    const closeCancelModal = () => {
      showCancelModal.value = false
      cancellingOrderId.value = null
    }

    const confirmCancel = async () => {
      if (!cancellingOrderId.value) return
      
      cancelling.value = true
      try {
        const res = await orderApi.cancel(cancellingOrderId.value)
        if (res.code === 200) {
          closeCancelModal()
          loadOrders()
        } else {
          alert(res.message || '取消失败')
        }
      } catch (error) {
        console.error('取消订单失败:', error)
        alert('取消失败，请稍后重试')
      } finally {
        cancelling.value = false
      }
    }

    onMounted(() => {
      loadOrders()
    })

    return {
      loading,
      orders,
      total,
      page,
      pageSize,
      selectedStatus,
      showCancelModal,
      cancellingOrderId,
      cancelling,
      statusTabs,
      totalPages,
      getStatusText,
      getStatusClass,
      formatDate,
      selectStatus,
      changePage,
      cancelOrder,
      closeCancelModal,
      confirmCancel
    }
  }
}
</script>

<style scoped>
.orders-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 30px 20px;
}

.orders-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
}

.status-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.status-tab {
  padding: 8px 20px;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.status-tab:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.status-tab.active {
  background: #1a73e8;
  border-color: #1a73e8;
  color: #fff;
}

.loading {
  text-align: center;
  padding: 80px 0;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e0e0e0;
  border-top-color: #1a73e8;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fafafa;
  border-bottom: 1px solid #eee;
}

.order-info {
  display: flex;
  gap: 20px;
}

.order-no {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.order-date {
  font-size: 13px;
  color: #999;
}

.order-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
}

.status-pending {
  background: #fff3cd;
  color: #856404;
}

.status-confirmed {
  background: #d1ecf1;
  color: #0c5460;
}

.status-checked {
  background: #d4edda;
  color: #155724;
}

.status-done {
  background: #e8f0fe;
  color: #1a73e8;
}

.status-cancelled {
  background: #f8d7da;
  color: #721c24;
}

.order-content {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 20px;
}

.hotel-info {
  flex: 1;
}

.hotel-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.room-name {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}

.stay-dates {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.date-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.date-label {
  font-size: 12px;
  color: #999;
}

.date-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.date-arrow {
  font-size: 16px;
  color: #ccc;
}

.guest-info {
  display: flex;
  gap: 20px;
}

.guest-item {
  font-size: 13px;
  color: #666;
}

.price-actions {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-end;
  border-left: 1px solid #eee;
  padding-left: 20px;
}

.price-info {
  text-align: right;
}

.price-label {
  display: block;
  font-size: 13px;
  color: #999;
  margin-bottom: 4px;
}

.price-value {
  font-size: 20px;
  font-weight: 600;
  color: #e74c3c;
}

.order-actions {
  margin-top: 16px;
}

.btn-cancel-order {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid #e74c3c;
  color: #e74c3c;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel-order:hover {
  background: #e74c3c;
  color: #fff;
}

.empty-state {
  text-align: center;
  padding: 80px 0;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 15px;
  color: #999;
  margin-bottom: 24px;
}

.btn-go-book {
  display: inline-block;
  padding: 12px 32px;
  background: #1a73e8;
  color: #fff;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-go-book:hover {
  background: #1557b0;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-top: 40px;
}

.page-btn {
  padding: 10px 24px;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: #1a73e8;
  color: #1a73e8;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: #fff;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
}

.confirm-modal {
  text-align: center;
  padding: 40px;
}

.modal-icon {
  font-size: 56px;
  margin-bottom: 16px;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.modal-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 30px;
}

.modal-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.btn-secondary {
  padding: 12px 28px;
  background: #f5f5f5;
  border: none;
  color: #666;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-danger {
  padding: 12px 28px;
  background: #e74c3c;
  border: none;
  color: #fff;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-danger:hover:not(:disabled) {
  background: #c0392b;
}

.btn-danger:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .order-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .order-info {
    flex-direction: column;
    gap: 8px;
  }

  .order-content {
    flex-direction: column;
    gap: 20px;
  }

  .price-actions {
    border-left: none;
    border-top: 1px solid #eee;
    padding-left: 0;
    padding-top: 20px;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }

  .order-actions {
    margin-top: 0;
  }
}
</style>
