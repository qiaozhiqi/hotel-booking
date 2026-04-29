<template>
  <div class="invoices-page">
    <div class="invoices-header">
      <h1 class="page-title">发票管理</h1>
    </div>

    <div class="tab-nav">
      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'list' }"
        @click="activeTab = 'list'"
      >
        我的发票
      </button>
      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'create' }"
        @click="activeTab = 'create'"
      >
        申请发票
      </button>
    </div>

    <div v-if="activeTab === 'list'" class="invoice-list-section">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <template v-else-if="invoices.length > 0">
        <div class="status-filter">
          <button 
            v-for="filter in statusFilters" 
            :key="filter.value"
            class="filter-btn"
            :class="{ active: selectedStatus === filter.value }"
            @click="selectStatus(filter.value)"
          >
            {{ filter.label }}
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
              <span class="invoice-status" :class="getStatusClass(invoice.status)">
                {{ getStatusText(invoice.status) }}
              </span>
            </div>

            <div class="invoice-content">
              <div class="invoice-details">
                <div class="detail-row">
                  <span class="detail-label">订单号：</span>
                  <span class="detail-value">{{ invoice.order_no }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">发票类型：</span>
                  <span class="detail-value">{{ invoice.invoice_type === 'personal' ? '个人发票' : '企业发票' }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">发票抬头：</span>
                  <span class="detail-value">{{ invoice.invoice_title }}</span>
                </div>
                <div v-if="invoice.tax_number" class="detail-row">
                  <span class="detail-label">税号：</span>
                  <span class="detail-value">{{ invoice.tax_number }}</span>
                </div>
                <div v-if="invoice.email" class="detail-row">
                  <span class="detail-label">接收邮箱：</span>
                  <span class="detail-value">{{ invoice.email }}</span>
                </div>
              </div>
              <div class="invoice-amount">
                <span class="amount-label">开票金额</span>
                <span class="amount-value">¥{{ invoice.amount }}</span>
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
        <div class="empty-icon">📄</div>
        <p class="empty-text">暂无发票记录</p>
        <button class="btn-go-create" @click="activeTab = 'create'">去申请发票</button>
      </div>
    </div>

    <div v-else class="invoice-create-section">
      <div class="create-form">
        <h2 class="form-title">申请发票</h2>

        <div class="form-group">
          <label class="form-label">选择订单 <span class="required">*</span></label>
          <select v-model="form.order_id" class="form-select" :disabled="loadingOrders">
            <option value="">请选择需要开票的订单</option>
            <option v-for="order in availableOrders" :key="order.id" :value="order.id">
              {{ order.order_no }} - {{ order.hotel_name }} {{ order.room_name }} (¥{{ order.total_amount }})
            </option>
          </select>
          <p v-if="!loadingOrders && availableOrders.length === 0" class="form-hint">暂无可开票的订单</p>
        </div>

        <div class="form-group">
          <label class="form-label">发票类型 <span class="required">*</span></label>
          <div class="type-radio-group">
            <label class="type-radio">
              <input type="radio" v-model="form.invoice_type" value="personal" />
              <span class="radio-label">个人发票</span>
            </label>
            <label class="type-radio">
              <input type="radio" v-model="form.invoice_type" value="company" />
              <span class="radio-label">企业发票</span>
            </label>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">发票抬头 <span class="required">*</span></label>
          <input 
            v-model="form.invoice_title" 
            type="text" 
            class="form-input" 
            :placeholder="form.invoice_type === 'personal' ? '请输入个人姓名' : '请输入企业名称'"
          />
        </div>

        <template v-if="form.invoice_type === 'company'">
          <div class="form-group">
            <label class="form-label">税号 <span class="required">*</span></label>
            <input v-model="form.tax_number" type="text" class="form-input" placeholder="请输入纳税人识别号" />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">开户银行</label>
              <input v-model="form.bank_name" type="text" class="form-input" placeholder="请输入开户银行（选填）" />
            </div>
            <div class="form-group">
              <label class="form-label">银行账号</label>
              <input v-model="form.bank_account" type="text" class="form-input" placeholder="请输入银行账号（选填）" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">企业地址</label>
              <input v-model="form.address" type="text" class="form-input" placeholder="请输入企业地址（选填）" />
            </div>
            <div class="form-group">
              <label class="form-label">企业电话</label>
              <input v-model="form.phone" type="text" class="form-input" placeholder="请输入企业电话（选填）" />
            </div>
          </div>
        </template>

        <div class="form-group">
          <label class="form-label">接收邮箱</label>
          <input v-model="form.email" type="email" class="form-input" placeholder="请输入发票接收邮箱（选填）" />
        </div>

        <div class="form-actions">
          <button class="btn-cancel" @click="activeTab = 'list'">取消</button>
          <button class="btn-submit" @click="submitInvoice" :disabled="submitting || !canSubmit">
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
    const activeTab = ref('list')
    const loading = ref(true)
    const loadingOrders = ref(false)
    const submitting = ref(false)
    const invoices = ref([])
    const availableOrders = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(10)
    const selectedStatus = ref('')

    const form = ref({
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

    const statusFilters = [
      { label: '全部', value: '' },
      { label: '待开票', value: 'pending' },
      { label: '已开票', value: 'issued' }
    ]

    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

    const canSubmit = computed(() => {
      if (!form.value.order_id) return false
      if (!form.value.invoice_title) return false
      if (form.value.invoice_type === 'company' && !form.value.tax_number) return false
      return true
    })

    const getStatusText = (status) => {
      const map = {
        'pending': '待开票',
        'issued': '已开票',
        'failed': '开票失败'
      }
      return map[status] || status
    }

    const getStatusClass = (status) => {
      const map = {
        'pending': 'status-pending',
        'issued': 'status-issued',
        'failed': 'status-failed'
      }
      return map[status] || ''
    }

    const formatDate = (dateStr) => {
      const date = new Date(dateStr)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
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
      loadingOrders.value = true
      try {
        const res = await invoiceApi.getAvailableOrders()
        if (res.code === 200) {
          availableOrders.value = res.data || []
        }
      } catch (error) {
        console.error('加载可开票订单失败:', error)
      } finally {
        loadingOrders.value = false
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

    const submitInvoice = async () => {
      if (!canSubmit.value) return

      submitting.value = true
      try {
        const data = {
          order_id: form.value.order_id,
          invoice_type: form.value.invoice_type,
          invoice_title: form.value.invoice_title,
          tax_number: form.value.tax_number,
          bank_name: form.value.bank_name,
          bank_account: form.value.bank_account,
          address: form.value.address,
          phone: form.value.phone,
          email: form.value.email
        }
        const res = await invoiceApi.create(data)
        if (res.code === 200) {
          alert('发票申请提交成功！')
          form.value = {
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
          activeTab.value = 'list'
          loadInvoices()
        } else {
          alert(res.message || '提交失败')
        }
      } catch (error) {
        console.error('提交发票申请失败:', error)
        alert('提交失败，请稍后重试')
      } finally {
        submitting.value = false
      }
    }

    watch(activeTab, (newTab) => {
      if (newTab === 'list') {
        loadInvoices()
      } else if (newTab === 'create') {
        loadAvailableOrders()
      }
    })

    onMounted(() => {
      if (activeTab.value === 'list') {
        loadInvoices()
      }
    })

    return {
      activeTab,
      loading,
      loadingOrders,
      submitting,
      invoices,
      availableOrders,
      total,
      page,
      pageSize,
      selectedStatus,
      form,
      statusFilters,
      totalPages,
      canSubmit,
      getStatusText,
      getStatusClass,
      formatDate,
      selectStatus,
      changePage,
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
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.tab-nav {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.tab-btn {
  padding: 10px 24px;
  background: #f5f5f5;
  border: none;
  border-radius: 8px 8px 0 0;
  font-size: 15px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: #e8e8e8;
  color: #333;
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

.status-filter {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.filter-btn {
  padding: 6px 16px;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.filter-btn.active {
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

.status-failed {
  background: #f8d7da;
  color: #721c24;
}

.invoice-content {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 20px;
}

.invoice-details {
  flex: 1;
}

.detail-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  font-size: 14px;
  color: #999;
}

.detail-value {
  font-size: 14px;
  color: #333;
}

.invoice-amount {
  text-align: right;
  border-left: 1px solid #eee;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.amount-label {
  display: block;
  font-size: 13px;
  color: #999;
  margin-bottom: 4px;
}

.amount-value {
  font-size: 22px;
  font-weight: 600;
  color: #e74c3c;
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

.btn-go-create {
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

.btn-go-create:hover {
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

.invoice-create-section {
  background: #fff;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.create-form {
  max-width: 600px;
}

.form-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}

.required {
  color: #e74c3c;
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

.form-input:disabled,
.form-select:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.form-hint {
  font-size: 13px;
  color: #999;
  margin-top: 8px;
}

.form-row {
  display: flex;
  gap: 20px;
}

.form-row .form-group {
  flex: 1;
}

.type-radio-group {
  display: flex;
  gap: 30px;
}

.type-radio {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.type-radio input[type="radio"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.radio-label {
  font-size: 14px;
  color: #333;
}

.form-actions {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #eee;
}

.btn-cancel {
  padding: 12px 32px;
  background: #f5f5f5;
  border: none;
  color: #666;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-cancel:hover {
  background: #e8e8e8;
}

.btn-submit {
  padding: 12px 32px;
  background: #1a73e8;
  border: none;
  color: #fff;
  border-radius: 8px;
  font-size: 14px;
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
  .form-row {
    flex-direction: column;
    gap: 0;
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

  .invoice-content {
    flex-direction: column;
    gap: 20px;
  }

  .invoice-amount {
    border-left: none;
    border-top: 1px solid #eee;
    padding-left: 0;
    padding-top: 20px;
    flex-direction: row;
    align-items: center;
    gap: 16px;
  }

  .amount-label {
    margin-bottom: 0;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions button {
    width: 100%;
  }
}
</style>
