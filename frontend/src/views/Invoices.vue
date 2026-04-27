<template>
  <div class="invoices-page">
    <div class="invoices-header">
      <h1 class="page-title">开发票</h1>
      <div class="tab-switcher">
        <button 
          class="tab-btn"
          :class="{ active: currentTab === 'list' }"
          @click="switchTab('list')"
        >
          我的发票
        </button>
        <button 
          class="tab-btn"
          :class="{ active: currentTab === 'create' }"
          @click="switchTab('create')"
        >
          申请发票
        </button>
      </div>
    </div>

    <div v-if="currentTab === 'list'">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <template v-else-if="invoices.length > 0">
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

        <div class="invoices-list">
          <div 
            v-for="invoice in invoices" 
            :key="invoice.id" 
            class="invoice-card"
          >
            <div class="invoice-header">
              <div class="invoice-info">
                <span class="invoice-no">发票号：{{ invoice.invoice_no || '-' }}</span>
                <span class="invoice-date">{{ formatDate(invoice.created_at) }}</span>
              </div>
              <span class="invoice-status" :class="getInvoiceStatusClass(invoice.status)">
                {{ getInvoiceStatusText(invoice.status) }}
              </span>
            </div>

            <div class="invoice-content">
              <div class="invoice-detail">
                <div class="detail-row">
                  <span class="detail-label">订单号</span>
                  <span class="detail-value">{{ invoice.order_no }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">发票类型</span>
                  <span class="detail-value">{{ getInvoiceTypeText(invoice.invoice_type) }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">发票抬头</span>
                  <span class="detail-value">{{ invoice.invoice_title }}</span>
                </div>
                <div class="detail-row" v-if="invoice.tax_number">
                  <span class="detail-label">税号</span>
                  <span class="detail-value">{{ invoice.tax_number }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">开票金额</span>
                  <span class="detail-value price">¥{{ invoice.amount }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">接收邮箱</span>
                  <span class="detail-value">{{ invoice.email }}</span>
                </div>
              </div>
            </div>

            <div class="invoice-footer" v-if="invoice.invoice_url">
              <a :href="invoice.invoice_url" target="_blank" class="btn-download">
                下载发票
              </a>
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
        <div class="empty-icon">📄</div>
        <p class="empty-text">暂无发票记录</p>
        <button class="btn-create-invoice" @click="switchTab('create')">
          去申请发票
        </button>
      </div>
    </div>

    <div v-else>
      <div class="create-invoice-form">
        <div class="form-section">
          <h3 class="section-title">选择订单</h3>
          <div class="order-selector">
            <label class="form-label">可开票订单</label>
            <select 
              v-model="invoiceForm.order_id" 
              class="form-select"
              @change="onOrderSelect"
            >
              <option :value="null">请选择订单</option>
              <option v-for="order in availableOrders" :key="order.id" :value="order.id">
                {{ order.order_no }} - {{ order.hotel_name }} - ¥{{ order.total_amount }}
              </option>
            </select>
            <div v-if="selectedOrder" class="selected-order-info">
              <div class="order-info-item">
                <span class="info-label">酒店：</span>
                <span class="info-value">{{ selectedOrder.hotel_name }}</span>
              </div>
              <div class="order-info-item">
                <span class="info-label">房型：</span>
                <span class="info-value">{{ selectedOrder.room_name }}</span>
              </div>
              <div class="order-info-item">
                <span class="info-label">入住：</span>
                <span class="info-value">{{ selectedOrder.check_in }} 至 {{ selectedOrder.check_out }}</span>
              </div>
              <div class="order-info-item">
                <span class="info-label">金额：</span>
                <span class="info-value price">¥{{ selectedOrder.total_amount }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">发票信息</h3>
          
          <div class="form-group">
            <label class="form-label required">发票类型</label>
            <div class="type-selector">
              <button 
                class="type-btn"
                :class="{ active: invoiceForm.invoice_type === 'personal' }"
                @click="invoiceForm.invoice_type = 'personal'"
              >
                个人
              </button>
              <button 
                class="type-btn"
                :class="{ active: invoiceForm.invoice_type === 'company' }"
                @click="invoiceForm.invoice_type = 'company'"
              >
                企业
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label required">发票抬头</label>
            <input 
              type="text" 
              v-model="invoiceForm.invoice_title" 
              class="form-input"
              :placeholder="invoiceForm.invoice_type === 'personal' ? '请输入个人姓名' : '请输入企业名称'"
            />
          </div>

          <div v-if="invoiceForm.invoice_type === 'company'" class="company-fields">
            <div class="form-group">
              <label class="form-label required">税号</label>
              <input 
                type="text" 
                v-model="invoiceForm.tax_number" 
                class="form-input"
                placeholder="请输入纳税人识别号"
              />
            </div>
            <div class="form-group">
              <label class="form-label">开户银行</label>
              <input 
                type="text" 
                v-model="invoiceForm.bank_name" 
                class="form-input"
                placeholder="请输入开户银行名称（选填）"
              />
            </div>
            <div class="form-group">
              <label class="form-label">银行账号</label>
              <input 
                type="text" 
                v-model="invoiceForm.bank_account" 
                class="form-input"
                placeholder="请输入银行账号（选填）"
              />
            </div>
            <div class="form-group">
              <label class="form-label">企业地址</label>
              <input 
                type="text" 
                v-model="invoiceForm.address" 
                class="form-input"
                placeholder="请输入企业地址（选填）"
              />
            </div>
            <div class="form-group">
              <label class="form-label">企业电话</label>
              <input 
                type="text" 
                v-model="invoiceForm.phone" 
                class="form-input"
                placeholder="请输入企业电话（选填）"
              />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label required">接收邮箱</label>
            <input 
              type="email" 
              v-model="invoiceForm.email" 
              class="form-input"
              placeholder="请输入接收发票的邮箱地址"
            />
          </div>
        </div>

        <div class="form-actions">
          <button class="btn-submit" @click="submitInvoice" :disabled="submitting">
            <span v-if="submitting">提交中...</span>
            <span v-else>提交申请</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { invoiceApi } from '../api'

export default {
  name: 'Invoices',
  setup() {
    const currentTab = ref('list')
    const loading = ref(false)
    const invoices = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(10)
    const selectedStatus = ref('')

    const statusTabs = [
      { label: '全部', value: '' },
      { label: '待处理', value: 'pending' },
      { label: '已开具', value: 'issued' },
      { label: '已作废', value: 'voided' }
    ]

    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

    const availableOrders = ref([])
    const selectedOrder = ref(null)
    const submitting = ref(false)

    const invoiceForm = ref({
      order_id: null,
      invoice_type: 'personal',
      invoice_title: '',
      tax_number: '',
      bank_name: '',
      bank_account: '',
      address: '',
      phone: '',
      email: ''
    })

    const getInvoiceStatusText = (status) => {
      const map = {
        'pending': '待处理',
        'issued': '已开具',
        'voided': '已作废'
      }
      return map[status] || status
    }

    const getInvoiceStatusClass = (status) => {
      const map = {
        'pending': 'status-pending',
        'issued': 'status-issued',
        'voided': 'status-voided'
      }
      return map[status] || ''
    }

    const getInvoiceTypeText = (type) => {
      return type === 'company' ? '企业发票' : '个人发票'
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

    const switchTab = (tab) => {
      currentTab.value = tab
      if (tab === 'list') {
        loadInvoices()
      } else {
        loadAvailableOrders()
      }
    }

    const loadInvoices = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        if (selectedStatus.value) {
          params.status = selectedStatus.value
        }
        const res = await invoiceApi.getList(params)
        if (res.code === 200) {
          invoices.value = res.data.invoices || []
          total.value = res.data.total || 0
        }
      } catch (error) {
        console.error('加载发票列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    const loadAvailableOrders = async () => {
      try {
        const res = await invoiceApi.getAvailableOrders()
        if (res.code === 200) {
          availableOrders.value = res.data || []
        }
      } catch (error) {
        console.error('加载可开票订单失败:', error)
      }
    }

    const selectStatus = (status) => {
      selectedStatus.value = status
      page.value = 1
      loadInvoices()
    }

    const changePage = (newPage) => {
      page.value = newPage
      loadInvoices()
    }

    const onOrderSelect = () => {
      if (invoiceForm.value.order_id) {
        selectedOrder.value = availableOrders.value.find(o => o.id === invoiceForm.value.order_id)
      } else {
        selectedOrder.value = null
      }
    }

    const validateForm = () => {
      if (!invoiceForm.value.order_id) {
        alert('请选择订单')
        return false
      }
      if (!invoiceForm.value.invoice_title) {
        alert('请输入发票抬头')
        return false
      }
      if (invoiceForm.value.invoice_type === 'company' && !invoiceForm.value.tax_number) {
        alert('企业发票请输入税号')
        return false
      }
      if (!invoiceForm.value.email) {
        alert('请输入接收邮箱')
        return false
      }
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(invoiceForm.value.email)) {
        alert('请输入有效的邮箱地址')
        return false
      }
      return true
    }

    const submitInvoice = async () => {
      if (!validateForm()) return

      submitting.value = true
      try {
        const data = {
          order_id: invoiceForm.value.order_id,
          invoice_type: invoiceForm.value.invoice_type,
          invoice_title: invoiceForm.value.invoice_title,
          email: invoiceForm.value.email
        }
        
        if (invoiceForm.value.invoice_type === 'company') {
          data.tax_number = invoiceForm.value.tax_number
          data.bank_name = invoiceForm.value.bank_name
          data.bank_account = invoiceForm.value.bank_account
          data.address = invoiceForm.value.address
          data.phone = invoiceForm.value.phone
        }

        const res = await invoiceApi.create(data)
        if (res.code === 200) {
          alert('发票申请成功！')
          invoiceForm.value = {
            order_id: null,
            invoice_type: 'personal',
            invoice_title: '',
            tax_number: '',
            bank_name: '',
            bank_account: '',
            address: '',
            phone: '',
            email: ''
          }
          selectedOrder.value = null
          switchTab('list')
        } else {
          alert(res.message || '申请失败')
        }
      } catch (error) {
        console.error('提交发票申请失败:', error)
        alert('提交失败，请稍后重试')
      } finally {
        submitting.value = false
      }
    }

    onMounted(() => {
      loadInvoices()
    })

    return {
      currentTab,
      loading,
      invoices,
      total,
      page,
      pageSize,
      selectedStatus,
      statusTabs,
      totalPages,
      availableOrders,
      selectedOrder,
      submitting,
      invoiceForm,
      getInvoiceStatusText,
      getInvoiceStatusClass,
      getInvoiceTypeText,
      formatDate,
      switchTab,
      selectStatus,
      changePage,
      onOrderSelect,
      submitInvoice
    }
  }
}
</script>

<style scoped>
.invoices-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 30px 20px;
}

.invoices-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.tab-switcher {
  display: flex;
  gap: 8px;
}

.tab-btn {
  padding: 8px 20px;
  background: #f5f5f5;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: #e0e0e0;
}

.tab-btn.active {
  background: #1a73e8;
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

.status-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 20px;
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

.invoices-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.invoice-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.invoice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fafafa;
  border-bottom: 1px solid #eee;
}

.invoice-info {
  display: flex;
  gap: 20px;
}

.invoice-no {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.invoice-date {
  font-size: 13px;
  color: #999;
}

.invoice-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
}

.status-pending {
  background: #fff3cd;
  color: #856404;
}

.status-issued {
  background: #d4edda;
  color: #155724;
}

.status-voided {
  background: #f8d7da;
  color: #721c24;
}

.invoice-content {
  padding: 20px;
}

.invoice-detail {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.detail-row {
  display: flex;
  gap: 8px;
}

.detail-label {
  font-size: 13px;
  color: #999;
  min-width: 70px;
}

.detail-value {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.detail-value.price {
  color: #e74c3c;
  font-size: 15px;
}

.invoice-footer {
  padding: 12px 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
}

.btn-download {
  padding: 8px 20px;
  background: #1a73e8;
  color: #fff;
  border-radius: 6px;
  font-size: 13px;
  text-decoration: none;
  transition: background 0.2s;
}

.btn-download:hover {
  background: #1557b0;
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

.btn-create-invoice {
  padding: 12px 32px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-create-invoice:hover {
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

.create-invoice-form {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.form-section {
  margin-bottom: 32px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-label.required::after {
  content: '*';
  color: #e74c3c;
  margin-left: 4px;
}

.form-input,
.form-select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  color: #333;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #1a73e8;
}

.form-input::placeholder,
.form-select::placeholder {
  color: #999;
}

.type-selector {
  display: flex;
  gap: 12px;
}

.type-btn {
  flex: 1;
  padding: 12px;
  background: #f5f5f5;
  border: 2px solid transparent;
  border-radius: 8px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.type-btn:hover {
  background: #e0e0e0;
}

.type-btn.active {
  background: #e8f0fe;
  border-color: #1a73e8;
  color: #1a73e8;
  font-weight: 500;
}

.selected-order-info {
  margin-top: 12px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.order-info-item {
  display: flex;
  gap: 8px;
}

.info-label {
  font-size: 13px;
  color: #999;
}

.info-value {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.info-value.price {
  color: #e74c3c;
  font-size: 14px;
}

.company-fields {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.btn-submit {
  padding: 12px 40px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-submit:hover:not(:disabled) {
  background: #1557b0;
}

.btn-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .invoices-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .invoice-detail {
    grid-template-columns: 1fr;
  }

  .selected-order-info {
    grid-template-columns: 1fr;
  }

  .invoice-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .invoice-info {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
