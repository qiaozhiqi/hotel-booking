<template>
  <div class="hotel-detail">
    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <template v-else-if="hotel">
      <div class="hotel-header">
        <div class="hotel-image-large">
          <img :src="hotel.image_url" :alt="hotel.name" />
          <div 
            v-if="hotel.supplier_code" 
            class="supplier-badge-large" 
            :style="{ backgroundColor: getSupplierColor(hotel.supplier_code) }"
          >
            <span class="supplier-icon">{{ getSupplierIcon(hotel.supplier_code) }}</span>
            <span class="supplier-name">{{ hotel.supplier_name }}</span>
          </div>
        </div>
        <div class="hotel-info-large">
          <div class="hotel-title-row">
            <h1 class="hotel-name-large">{{ hotel.name }}</h1>
            <span v-if="hotel.brand" class="brand-tag-large">{{ hotel.brand }}</span>
          </div>
          <div class="hotel-meta">
            <div class="meta-item">
              <span class="meta-icon">📍</span>
              <span class="meta-text">{{ hotel.city }} · {{ hotel.address }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-icon">⭐</span>
              <span class="meta-text rating-high">{{ hotel.rating }} 分</span>
            </div>
          </div>
          <p class="hotel-desc-large">{{ hotel.description }}</p>
          <div class="price-range-large">
            <span class="price-label">价格范围：</span>
            <span class="price-value-large">
              <span v-if="hotel.min_price > 0">¥{{ hotel.min_price }} 起</span>
              <span v-else>{{ hotel.price_range }}</span>
            </span>
          </div>
        </div>
      </div>

      <div class="rooms-section">
        <h2 class="section-title">房型选择</h2>
        
        <div class="room-list">
          <div 
            v-for="room in rooms" 
            :key="room.id" 
            class="room-card"
          >
            <div class="room-image">
              <img :src="room.image_url" :alt="room.name" />
              <div 
                v-if="room.supplier_code" 
                class="room-supplier-badge" 
                :style="{ backgroundColor: getSupplierColor(room.supplier_code) }"
              >
                <span class="supplier-icon-small">{{ getSupplierIcon(room.supplier_code) }}</span>
              </div>
            </div>
            <div class="room-info">
              <div class="room-header-row">
                <h3 class="room-name">{{ room.name }}</h3>
                <div class="room-tags">
                  <span 
                    v-if="room.is_price_controlled" 
                    class="price-control-tag" 
                    :title="room.price_control_reason"
                  >
                    管控价
                  </span>
                  <span v-if="room.promotion_tag" class="promotion-tag">
                    {{ room.promotion_tag }}
                  </span>
                </div>
              </div>
              <div class="room-features">
                <span class="feature-tag">
                  <span class="feature-icon">🛏️</span>
                  {{ room.bed_type }}
                </span>
                <span class="feature-tag">
                  <span class="feature-icon">📐</span>
                  {{ room.area }} ㎡
                </span>
                <span class="feature-tag">
                  <span class="feature-icon">👥</span>
                  最多 {{ room.capacity }} 人
                </span>
                <span v-if="room.payment_type" class="feature-tag">
                  <span class="feature-icon">💳</span>
                  {{ room.payment_type }}
                </span>
              </div>
              <p class="room-desc">{{ room.description }}</p>
              <div class="room-amenities">
                <span class="amenities-label">设施：</span>
                <span class="amenities-text">{{ room.amenities }}</span>
              </div>
              <div class="room-availability-row">
                <div class="room-availability" :class="{ 'available': room.available_count > 0 }">
                  <span v-if="room.available_count > 0" class="available-text">
                    剩余 {{ room.available_count }} 间
                  </span>
                  <span v-else class="unavailable-text">暂无空房</span>
                </div>
                <button 
                  class="btn-compare" 
                  @click="openComparisonModal(room)"
                >
                  🔄 渠道比价
                </button>
              </div>
              <div v-if="room.cancel_policy" class="cancel-policy">
                <span class="policy-label">退改政策：</span>
                <span class="policy-text">{{ room.cancel_policy }}</span>
              </div>
            </div>
            <div class="room-price-box">
              <div v-if="room.original_price > room.price" class="original-price">
                <span class="original-price-num">¥{{ room.original_price }}</span>
              </div>
              <div class="room-price">
                <span class="price-currency">¥</span>
                <span class="price-num">{{ room.price }}</span>
                <span class="price-unit">/晚</span>
              </div>
              <div v-if="room.supplier_name" class="room-supplier-info">
                <span class="supplier-icon-small">{{ getSupplierIcon(room.supplier_code) }}</span>
                <span class="supplier-text">{{ room.supplier_name }}</span>
              </div>
              <button 
                class="btn-book" 
                :disabled="room.available_count <= 0"
                @click="openBookingModal(room)"
              >
                立即预订
              </button>
            </div>
          </div>
        </div>

        <div class="empty-state" v-if="rooms.length === 0">
          <div class="empty-icon">🛏️</div>
          <p class="empty-text">暂无房型信息</p>
        </div>
      </div>
    </template>

    <div class="empty-state" v-else>
      <div class="empty-icon">🏨</div>
      <p class="empty-text">酒店不存在</p>
      <router-link to="/" class="btn-back">返回首页</router-link>
    </div>

    <div v-if="showBookingModal" class="modal-overlay" @click.self="closeBookingModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">预订确认</h3>
          <button class="modal-close" @click="closeBookingModal">×</button>
        </div>
        
        <div class="modal-body" v-if="selectedRoom">
          <div class="booking-summary">
            <h4 class="summary-title">{{ hotel.name }}</h4>
            <p class="summary-room">{{ selectedRoom.name }}</p>
          </div>

          <div class="booking-form">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">入住日期</label>
                <input type="date" v-model="bookingForm.checkIn" class="form-input" :min="today" />
              </div>
              <div class="form-group">
                <label class="form-label">离店日期</label>
                <input type="date" v-model="bookingForm.checkOut" class="form-input" :min="bookingForm.checkIn" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">入住人姓名</label>
                <input type="text" v-model="bookingForm.guestName" class="form-input" placeholder="请输入姓名" />
              </div>
              <div class="form-group">
                <label class="form-label">联系电话</label>
                <input type="tel" v-model="bookingForm.guestPhone" class="form-input" placeholder="请输入手机号" />
              </div>
            </div>
          </div>

          <div class="booking-price-summary" v-if="nights > 0">
            <div class="price-row">
              <span class="price-label">房费：¥{{ selectedRoom.price }} × {{ nights }} 晚</span>
              <span class="price-amount">¥{{ selectedRoom.price * nights }}</span>
            </div>
            <div class="price-row total">
              <span class="price-label">总计</span>
              <span class="price-amount total">¥{{ totalPrice }}</span>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-cancel" @click="closeBookingModal">取消</button>
          <button class="btn-confirm" :disabled="!canSubmit || submitting" @click="submitBooking">
            <span v-if="submitting">提交中...</span>
            <span v-else>确认预订</span>
          </button>
        </div>
      </div>
    </div>

    <div v-if="showSuccessModal" class="modal-overlay" @click.self="closeSuccessModal">
      <div class="modal-content success-modal">
        <div class="success-icon">✅</div>
        <h3 class="success-title">预订成功！</h3>
        <p class="success-text">您的订单已确认，我们会为您预留房间。</p>
        <div class="success-actions">
          <button class="btn-secondary" @click="closeSuccessModal">继续浏览</button>
          <router-link to="/orders" class="btn-primary" @click="closeSuccessModal">查看订单</router-link>
        </div>
      </div>
    </div>

    <div v-if="showComparisonModal" class="modal-overlay comparison-modal" @click.self="closeComparisonModal">
      <div class="modal-content comparison-content">
        <div class="modal-header">
          <h3 class="modal-title">
            <span class="title-icon">📊</span>
            渠道比价 - {{ comparisonData?.room_name || '房型' }}
          </h3>
          <button class="modal-close" @click="closeComparisonModal">×</button>
        </div>
        
        <div class="comparison-body" v-loading="loadingComparison">
          <div v-if="comparisonData && comparisonData.channels.length > 0">
            <div class="comparison-summary">
              <div class="summary-item best-price">
                <span class="summary-label">最优价格</span>
                <span class="summary-value">¥{{ comparisonData.best_price }}</span>
                <span class="summary-supplier">({{ comparisonData.best_price_supplier }})</span>
              </div>
              <div class="summary-item room-type">
                <span class="summary-label">标准房型</span>
                <span class="summary-value">{{ comparisonData.standard_room_type || '标准间' }}</span>
              </div>
            </div>

            <div class="comparison-table">
              <div class="table-header">
                <div class="col supplier-col">供应商</div>
                <div class="col price-col">价格</div>
                <div class="col stock-col">库存</div>
                <div class="col tags-col">标签</div>
                <div class="col action-col">操作</div>
              </div>
              <div 
                v-for="channel in sortedChannels" 
                :key="channel.supplier_code" 
                class="table-row"
                :class="{ 
                  'best-price-row': channel.is_best_price, 
                  'recommended-row': channel.is_recommended 
                }"
              >
                <div class="col supplier-col">
                  <div class="supplier-info">
                    <div 
                      class="supplier-avatar" 
                      :style="{ backgroundColor: getSupplierColor(channel.supplier_code) }"
                    >
                      {{ getSupplierIcon(channel.supplier_code) }}
                    </div>
                    <div class="supplier-details">
                      <div class="supplier-name">{{ channel.supplier_name }}</div>
                      <div class="supplier-priority">
                        优先级: <span class="priority-value">{{ channel.supplier_priority }}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="col price-col">
                  <div class="price-info">
                    <span v-if="channel.original_price > channel.price" class="original-price-cross">
                      ¥{{ channel.original_price }}
                    </span>
                    <span class="current-price" :class="{ 'best-price-text': channel.is_best_price }">
                      ¥{{ channel.price }}
                    </span>
                    <span class="price-unit">/晚</span>
                  </div>
                  <div v-if="channel.promotion_tag" class="price-promo">
                    {{ channel.promotion_tag }}
                  </div>
                </div>
                <div class="col stock-col">
                  <div class="stock-info" :class="{ 'low-stock': channel.available_count <= 3 }">
                    <span class="stock-icon">{{ channel.available_count > 0 ? '✅' : '❌' }}</span>
                    <span class="stock-text">
                      {{ channel.available_count > 0 ? `剩余 ${channel.available_count} 间` : '暂无库存' }}
                    </span>
                  </div>
                </div>
                <div class="col tags-col">
                  <div class="channel-tags">
                    <span v-if="channel.is_recommended" class="tag recommended-tag">
                      ⭐ 推荐
                    </span>
                    <span v-if="channel.is_best_price" class="tag best-price-tag">
                      🔥 最优价
                    </span>
                    <span v-if="channel.is_price_controlled" class="tag controlled-tag" :title="channel.price_control_reason">
                      📋 管控价
                    </span>
                    <span v-if="channel.payment_type" class="tag payment-tag">
                      {{ channel.payment_type }}
                    </span>
                  </div>
                </div>
                <div class="col action-col">
                  <button 
                    class="btn-select-channel"
                    :disabled="channel.available_count <= 0"
                    @click="selectChannel(channel)"
                  >
                    选择此渠道
                  </button>
                </div>
              </div>
            </div>

            <div v-if="comparisonData.channels.length > 0" class="comparison-footer">
              <div class="legend-info">
                <span class="legend-item">
                  <span class="legend-dot recommended"></span>
                  推荐渠道（优先级最高）
                </span>
                <span class="legend-item">
                  <span class="legend-dot best-price"></span>
                  最优价格
                </span>
                <span class="legend-item">
                  <span class="legend-dot controlled"></span>
                  管控价格（有特殊限制）
                </span>
              </div>
            </div>
          </div>
          <div v-else-if="!loadingComparison" class="empty-comparison">
            <div class="empty-icon">🔍</div>
            <p class="empty-text">暂无其他渠道的比价数据</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { hotelApi, orderApi } from '../api'

const supplierConfig = {
  huazhu: {
    code: 'huazhu',
    name: '华住酒店集团',
    color: '#FF6B35',
    icon: '🏨'
  },
  jinjiang: {
    code: 'jinjiang',
    name: '锦江国际酒店集团',
    color: '#1E88E5',
    icon: '🏩'
  },
  rujia: {
    code: 'rujia',
    name: '如家酒店集团',
    color: '#27AE60',
    icon: '🏠'
  }
}

export default {
  name: 'HotelDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()

    const today = new Date().toISOString().split('T')[0]
    const tomorrow = new Date(Date.now() + 86400000).toISOString().split('T')[0]

    const loading = ref(true)
    const hotel = ref(null)
    const rooms = ref([])

    const showBookingModal = ref(false)
    const selectedRoom = ref(null)
    const showSuccessModal = ref(false)
    const submitting = ref(false)

    const showComparisonModal = ref(false)
    const loadingComparison = ref(false)
    const comparisonData = ref(null)
    const selectedChannel = ref(null)

    const bookingForm = ref({
      checkIn: today,
      checkOut: tomorrow,
      guestName: '',
      guestPhone: ''
    })

    const nights = computed(() => {
      if (!bookingForm.value.checkIn || !bookingForm.value.checkOut) return 0
      const checkIn = new Date(bookingForm.value.checkIn)
      const checkOut = new Date(bookingForm.value.checkOut)
      const diff = checkOut.getTime() - checkIn.getTime()
      return Math.max(0, Math.floor(diff / (1000 * 60 * 60 * 24)))
    })

    const totalPrice = computed(() => {
      if (!selectedRoom.value) return 0
      return selectedRoom.value.price * nights.value
    })

    const canSubmit = computed(() => {
      return bookingForm.value.checkIn && 
             bookingForm.value.checkOut && 
             bookingForm.value.guestName && 
             bookingForm.value.guestPhone &&
             nights.value > 0
    })

    const sortedChannels = computed(() => {
      if (!comparisonData.value || !comparisonData.value.channels) return []
      return [...comparisonData.value.channels].sort((a, b) => {
        if (a.is_recommended !== b.is_recommended) {
          return a.is_recommended ? -1 : 1
        }
        if (a.is_best_price !== b.is_best_price) {
          return a.is_best_price ? -1 : 1
        }
        return a.supplier_priority - b.supplier_priority
      })
    })

    const getSupplierColor = (code) => {
      return supplierConfig[code]?.color || '#666'
    }
    
    const getSupplierIcon = (code) => {
      return supplierConfig[code]?.icon || '🏨'
    }

    const loadHotelDetail = async () => {
      loading.value = true
      try {
        const id = route.params.id
        const res = await hotelApi.getDetail(id)
        if (res.code === 200) {
          hotel.value = res.data.hotel
          rooms.value = res.data.rooms || []
        }
      } catch (error) {
        console.error('加载酒店详情失败:', error)
      } finally {
        loading.value = false
      }
    }

    const openBookingModal = (room) => {
      selectedRoom.value = room
      bookingForm.value = {
        checkIn: today,
        checkOut: tomorrow,
        guestName: '',
        guestPhone: ''
      }
      showBookingModal.value = true
    }

    const closeBookingModal = () => {
      showBookingModal.value = false
      selectedRoom.value = null
    }

    const submitBooking = async () => {
      if (!canSubmit.value) return
      
      const user = localStorage.getItem('user')
      if (!user) {
        alert('请先登录')
        closeBookingModal()
        router.push('/login')
        return
      }

      submitting.value = true
      try {
        const res = await orderApi.create({
          hotel_id: hotel.value.id,
          room_id: selectedRoom.value.id,
          check_in: bookingForm.value.checkIn,
          check_out: bookingForm.value.checkOut,
          guest_name: bookingForm.value.guestName,
          guest_phone: bookingForm.value.guestPhone
        })

        if (res.code === 200) {
          closeBookingModal()
          showSuccessModal.value = true
        } else {
          alert(res.message || '预订失败')
        }
      } catch (error) {
        console.error('预订失败:', error)
        alert('预订失败，请稍后重试')
      } finally {
        submitting.value = false
      }
    }

    const closeSuccessModal = () => {
      showSuccessModal.value = false
    }

    const loadComparison = async (hotelId, roomId) => {
      loadingComparison.value = true
      try {
        const res = await hotelApi.getRoomComparison(hotelId, roomId)
        if (res.code === 200) {
          comparisonData.value = res.data
        }
      } catch (error) {
        console.error('加载渠道比价失败:', error)
        comparisonData.value = null
      } finally {
        loadingComparison.value = false
      }
    }

    const openComparisonModal = (room) => {
      selectedRoom.value = room
      comparisonData.value = null
      showComparisonModal.value = true
      if (hotel.value && room.id) {
        loadComparison(hotel.value.id, room.id)
      }
    }

    const closeComparisonModal = () => {
      showComparisonModal.value = false
      comparisonData.value = null
      selectedChannel.value = null
    }

    const selectChannel = (channel) => {
      selectedChannel.value = channel
      alert(`已选择 ${channel.supplier_name} 渠道，价格 ¥${channel.price}`)
      closeComparisonModal()
    }

    onMounted(() => {
      loadHotelDetail()
    })

    return {
      today,
      loading,
      hotel,
      rooms,
      showBookingModal,
      selectedRoom,
      showSuccessModal,
      submitting,
      showComparisonModal,
      loadingComparison,
      comparisonData,
      selectedChannel,
      sortedChannels,
      bookingForm,
      nights,
      totalPrice,
      canSubmit,
      getSupplierColor,
      getSupplierIcon,
      openBookingModal,
      closeBookingModal,
      submitBooking,
      closeSuccessModal,
      openComparisonModal,
      closeComparisonModal,
      selectChannel
    }
  }
}
</script>

<style scoped>
.hotel-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 20px;
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

.hotel-header {
  display: flex;
  gap: 30px;
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.hotel-image-large {
  width: 400px;
  height: 280px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
  flex-shrink: 0;
}

.hotel-image-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hotel-info-large {
  flex: 1;
}

.hotel-name-large {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.hotel-meta {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-icon {
  font-size: 16px;
}

.meta-text {
  font-size: 14px;
  color: #666;
}

.meta-text.rating-high {
  color: #f59e0b;
  font-weight: 600;
}

.hotel-desc-large {
  font-size: 14px;
  color: #666;
  line-height: 1.8;
  margin-bottom: 20px;
}

.price-range-large {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.price-label {
  font-size: 14px;
  color: #999;
}

.price-value-large {
  font-size: 24px;
  font-weight: 600;
  color: #e74c3c;
}

.rooms-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.room-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.room-card {
  display: flex;
  gap: 20px;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.room-card:hover {
  border-color: #1a73e8;
  box-shadow: 0 4px 12px rgba(26, 115, 232, 0.1);
}

.room-image {
  width: 220px;
  height: 160px;
  border-radius: 6px;
  overflow: hidden;
  background: #f5f5f5;
  flex-shrink: 0;
}

.room-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.room-info {
  flex: 1;
}

.room-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.room-features {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.feature-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
  color: #666;
}

.feature-icon {
  font-size: 14px;
}

.room-desc {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
}

.room-amenities {
  font-size: 13px;
  margin-bottom: 12px;
}

.amenities-label {
  color: #999;
}

.amenities-text {
  color: #666;
}

.room-availability {
  font-size: 13px;
}

.room-availability.available .available-text {
  color: #27ae60;
}

.unavailable-text {
  color: #e74c3c;
}

.room-price-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
  flex-shrink: 0;
  min-width: 140px;
}

.room-price {
  display: flex;
  align-items: baseline;
  margin-bottom: 12px;
}

.price-currency {
  font-size: 14px;
  color: #e74c3c;
}

.price-num {
  font-size: 28px;
  font-weight: 600;
  color: #e74c3c;
}

.price-unit {
  font-size: 13px;
  color: #999;
}

.btn-book {
  padding: 10px 28px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-book:hover:not(:disabled) {
  background: #1557b0;
}

.btn-book:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}

.empty-icon {
  font-size: 56px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 15px;
  color: #999;
  margin-bottom: 20px;
}

.btn-back {
  display: inline-block;
  padding: 10px 24px;
  background: #1a73e8;
  color: #fff;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
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
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #eee;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  border-radius: 50%;
  transition: background 0.2s;
}

.modal-close:hover {
  background: #f5f5f5;
}

.modal-body {
  padding: 24px;
}

.booking-summary {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.summary-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.summary-room {
  font-size: 14px;
  color: #666;
}

.booking-form {
  margin-bottom: 20px;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.form-group {
  flex: 1;
}

.form-label {
  display: block;
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: #1a73e8;
}

.booking-price-summary {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.price-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 16px;
  font-size: 14px;
}

.price-row.total {
  background: #fafafa;
  font-weight: 600;
}

.price-amount.total {
  color: #e74c3c;
  font-size: 16px;
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 1px solid #eee;
}

.btn-cancel {
  padding: 10px 24px;
  background: transparent;
  border: 1px solid #ddd;
  color: #666;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  border-color: #999;
}

.btn-confirm {
  padding: 10px 24px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-confirm:hover:not(:disabled) {
  background: #1557b0;
}

.btn-confirm:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.success-modal {
  text-align: center;
  padding: 40px;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.success-title {
  font-size: 22px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.success-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 30px;
}

.success-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.btn-secondary {
  padding: 12px 28px;
  background: transparent;
  border: 1px solid #ddd;
  color: #666;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  border-color: #999;
}

.btn-primary {
  padding: 12px 28px;
  background: #1a73e8;
  color: #fff;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #1557b0;
}

.supplier-badge-large {
  position: absolute;
  top: 16px;
  left: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 107, 53, 0.95);
  border-radius: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.supplier-badge-large .supplier-icon {
  font-size: 18px;
}

.supplier-badge-large .supplier-name {
  font-size: 13px;
  font-weight: 600;
  color: #fff;
}

.hotel-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.hotel-title-row .hotel-name-large {
  margin-bottom: 0;
}

.brand-tag-large {
  flex-shrink: 0;
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;
}

.room-supplier-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 107, 53, 0.95);
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.supplier-icon-small {
  font-size: 14px;
}

.room-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.room-header-row .room-name {
  margin-bottom: 0;
}

.room-tags {
  display: flex;
  gap: 6px;
}

.price-control-tag {
  padding: 2px 8px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: #fff;
  font-size: 11px;
  font-weight: 500;
  border-radius: 4px;
}

.promotion-tag {
  padding: 2px 8px;
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: #fff;
  font-size: 11px;
  font-weight: 500;
  border-radius: 4px;
}

.room-availability-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.btn-compare {
  padding: 6px 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-compare:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.cancel-policy {
  margin-top: 8px;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 4px;
  font-size: 12px;
}

.policy-label {
  color: #999;
  font-weight: 500;
}

.policy-text {
  color: #666;
}

.original-price {
  margin-bottom: 4px;
}

.original-price-num {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
}

.room-supplier-info {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
  font-size: 12px;
  color: #666;
}

.supplier-text {
  font-size: 11px;
}

.comparison-modal .modal-content {
  max-width: 900px;
  width: 90vw;
}

.comparison-body {
  padding: 20px 24px;
}

.title-icon {
  margin-right: 8px;
}

.comparison-summary {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.summary-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.summary-item.best-price {
  padding: 12px 20px;
  background: linear-gradient(135deg, #ff6b6b 0%, #ffd93d 100%);
  border-radius: 8px;
}

.summary-item.best-price .summary-label {
  color: rgba(255, 255, 255, 0.9);
}

.summary-item.best-price .summary-value {
  color: #fff;
  font-size: 24px;
}

.summary-item.best-price .summary-supplier {
  color: rgba(255, 255, 255, 0.9);
}

.summary-label {
  font-size: 12px;
  color: #999;
}

.summary-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.summary-supplier {
  font-size: 12px;
  color: #666;
}

.comparison-table {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.table-header {
  display: flex;
  background: #f8f9fa;
  padding: 12px 16px;
  font-weight: 600;
  font-size: 13px;
  color: #666;
}

.table-row {
  display: flex;
  padding: 16px;
  border-top: 1px solid #eee;
  transition: background 0.2s;
}

.table-row:hover {
  background: #fafafa;
}

.table-row.best-price-row {
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.05) 0%, rgba(255, 217, 61, 0.05) 100%);
}

.table-row.recommended-row {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.col {
  display: flex;
  align-items: center;
}

.supplier-col {
  flex: 0 0 200px;
}

.price-col {
  flex: 0 0 140px;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
}

.stock-col {
  flex: 0 0 120px;
}

.tags-col {
  flex: 1;
}

.action-col {
  flex: 0 0 120px;
  justify-content: flex-end;
}

.supplier-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.supplier-avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 20px;
}

.supplier-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.supplier-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.supplier-priority {
  font-size: 12px;
  color: #999;
}

.priority-value {
  font-weight: 600;
  color: #667eea;
}

.price-info {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.original-price-cross {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.current-price {
  font-size: 20px;
  font-weight: 600;
  color: #e74c3c;
}

.current-price.best-price-text {
  color: #ff6b6b;
}

.price-unit {
  font-size: 12px;
  color: #999;
}

.price-promo {
  margin-top: 4px;
  padding: 2px 8px;
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: #fff;
  font-size: 10px;
  font-weight: 500;
  border-radius: 4px;
}

.stock-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.stock-info.low-stock {
  color: #f59e0b;
}

.stock-icon {
  font-size: 14px;
}

.stock-text {
  color: #666;
}

.channel-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.recommended-tag {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.best-price-tag {
  background: linear-gradient(135deg, #ff6b6b 0%, #ffd93d 100%);
  color: #fff;
}

.controlled-tag {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: #fff;
}

.payment-tag {
  background: #f0f0f0;
  color: #666;
}

.btn-select-channel {
  padding: 8px 16px;
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-select-channel:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-select-channel:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.comparison-footer {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.legend-info {
  display: flex;
  gap: 24px;
  font-size: 12px;
  color: #999;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-dot.recommended {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.legend-dot.best-price {
  background: linear-gradient(135deg, #ff6b6b 0%, #ffd93d 100%);
}

.legend-dot.controlled {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.empty-comparison {
  text-align: center;
  padding: 40px 0;
}

.empty-comparison .empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.empty-comparison .empty-text {
  font-size: 14px;
  color: #999;
}

@media (max-width: 900px) {
  .hotel-header {
    flex-direction: column;
  }
  
  .hotel-image-large {
    width: 100%;
  }
  
  .room-card {
    flex-direction: column;
  }
  
  .room-image {
    width: 100%;
    height: 200px;
  }
  
  .room-price-box {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }
  
  .room-price {
    margin-bottom: 0;
  }
}
</style>
