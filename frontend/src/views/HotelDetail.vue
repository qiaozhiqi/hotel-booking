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
        </div>
        <div class="hotel-info-large">
          <h1 class="hotel-name-large">{{ hotel.name }}</h1>
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
            <span class="price-value-large">{{ hotel.price_range }}</span>
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
            </div>
            <div class="room-info">
              <h3 class="room-name">{{ room.name }}</h3>
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
              </div>
              <p class="room-desc">{{ room.description }}</p>
              <div class="room-amenities">
                <span class="amenities-label">设施：</span>
                <span class="amenities-text">{{ room.amenities }}</span>
              </div>
              <div class="room-availability" :class="{ 'available': room.available_count > 0 }">
                <span v-if="room.available_count > 0" class="available-text">
                  剩余 {{ room.available_count }} 间
                </span>
                <span v-else class="unavailable-text">暂无空房</span>
              </div>
            </div>
            <div class="room-price-box">
              <div class="room-price">
                <span class="price-currency">¥</span>
                <span class="price-num">{{ room.price }}</span>
                <span class="price-unit">/晚</span>
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
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { hotelApi, orderApi } from '../api'

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
      bookingForm,
      nights,
      totalPrice,
      canSubmit,
      openBookingModal,
      closeBookingModal,
      submitBooking,
      closeSuccessModal
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
