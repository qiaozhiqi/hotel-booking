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
          <h1 class="hotel-name-large">
            {{ hotel.name }}
            <span v-if="hotel.supplier" class="supplier-badge-large">
              <span class="supplier-icon-large">🏢</span>
              {{ getSupplierShortName(hotel.supplier.name) }}
            </span>
          </h1>
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
          
          <div v-if="hotel.supplier" class="supplier-info-section">
            <div class="supplier-info-header">
              <span class="supplier-label">房源供应商：</span>
              <span class="supplier-name-large">{{ hotel.supplier.name }}</span>
            </div>
            <div class="supplier-info-details">
              <span v-if="hotel.supplier.priority > 0" class="priority-badge">
                <span class="priority-icon">🏆</span>
                优先级 {{ hotel.supplier.priority }}
              </span>
              <span class="price-control-badge">
                <span class="price-control-icon">💰</span>
                控价系数 {{ hotel.supplier.price_control }}
              </span>
            </div>
            <p v-if="hotel.supplier.description" class="supplier-desc">
              {{ hotel.supplier.description }}
            </p>
          </div>
          
          <p class="hotel-desc-large">{{ hotel.description }}</p>
          <div class="price-range-large">
            <span class="price-label">价格范围：</span>
            <span class="price-value-large">{{ hotel.price_range }}</span>
          </div>
        </div>
      </div>

      <div class="date-filter-section">
        <div class="date-filter-header">
          <h2 class="section-title">选择入住日期</h2>
          <div class="date-summary" v-if="calendarNights > 0">
            <span class="date-range-display">
              <span class="date-display-item">
                <span class="date-display-label">入住</span>
                <span class="date-display-value">{{ checkInDate }}</span>
              </span>
              <span class="date-arrow">→</span>
              <span class="date-display-item">
                <span class="date-display-label">离店</span>
                <span class="date-display-value">{{ checkOutDate }}</span>
              </span>
            </span>
            <span class="nights-badge">共 {{ calendarNights }} 晚</span>
          </div>
        </div>
        
        <div class="calendar-wrapper">
          <div class="calendar-nav">
            <button 
              class="nav-btn" 
              :disabled="!canShowPrevMonth"
              @click="prevMonth"
            >
              ‹
            </button>
            <span class="month-label">{{ currentMonthLabel }}</span>
            <button class="nav-btn" @click="nextMonth">›</button>
          </div>
          
          <div class="calendar-weekdays">
            <span class="weekday">日</span>
            <span class="weekday">一</span>
            <span class="weekday">二</span>
            <span class="weekday">三</span>
            <span class="weekday">四</span>
            <span class="weekday">五</span>
            <span class="weekday">六</span>
          </div>
          
          <div class="calendar-days">
            <div 
              v-for="(day, index) in calendarDays" 
              :key="index"
              class="calendar-day"
              :class="{
                'other-month': day.isOtherMonth,
                'disabled': !day.isSelectable && !day.isOtherMonth,
                'in-range': day.isInRange,
                'check-in': day.isCheckIn,
                'check-out': day.isCheckOut
              }"
              @click="selectDate(day)"
            >
              <span v-if="day.day !== ''" class="day-number">{{ day.day }}</span>
              <span 
                v-if="day.day !== '' && day.isSelectable && day.minPrice > 0" 
                class="day-price"
              >
                ¥{{ day.minPrice }}
              </span>
              <span 
                v-if="day.day !== '' && !day.isSelectable" 
                class="day-disabled-text"
              >
                不可选
              </span>
            </div>
          </div>
        </div>
        
        <div class="date-hint">
          <span class="hint-item">
            <span class="hint-icon hint-blue"></span>
            <span class="hint-text">入住日期</span>
          </span>
          <span class="hint-item">
            <span class="hint-icon hint-light-blue"></span>
            <span class="hint-text">入住期间</span>
          </span>
          <span class="hint-item">
            <span class="hint-icon hint-blue"></span>
            <span class="hint-text">离店日期</span>
          </span>
          <span class="hint-item">
            <span class="hint-icon hint-gray"></span>
            <span class="hint-text">不可选日期</span>
          </span>
        </div>
      </div>

      <div class="rooms-section">
        <h2 class="section-title">房型选择 <span class="section-subtitle">支持多渠道比价预订</span></h2>
        
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
              <div class="room-header-row">
                <h3 class="room-name">{{ room.name }}</h3>
                <div v-if="room.best_price < room.price" class="best-price-tag">
                  <span class="best-price-icon">🔥</span>
                  最优价比
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
              </div>
              <p class="room-desc">{{ room.description }}</p>
              <div class="room-amenities">
                <span class="amenities-label">设施：</span>
                <span class="amenities-text">{{ room.amenities }}</span>
              </div>
              
              <div v-if="room.channel_prices && room.channel_prices.length > 0" class="channel-prices-section">
                <div class="channel-prices-header">
                  <span class="channel-prices-title">📊 渠道比价</span>
                  <span class="channel-prices-count">共 {{ room.channel_prices.length }} 个渠道</span>
                </div>
                <div class="channel-prices-list">
                  <div 
                    v-for="(cp, index) in room.channel_prices" 
                    :key="index" 
                    class="channel-price-item"
                    :class="{ 
                      'best-price': cp.is_best_price && cp.available_count > 0,
                      'out-of-stock': cp.available_count <= 0 
                    }"
                  >
                    <div class="channel-price-info">
                      <div class="channel-name-row">
                        <span class="channel-name">{{ getSupplierShortName(cp.supplier_name) }}</span>
                        <span v-if="cp.is_best_price && cp.available_count > 0" class="best-price-badge">最低价</span>
                        <span v-if="cp.available_count <= 0" class="out-of-stock-badge">暂无房</span>
                      </div>
                      <div class="channel-details">
                        <span v-if="cp.priority > 0" class="channel-priority">优先级 {{ cp.priority }}</span>
                        <span class="channel-original-price" v-if="cp.price !== cp.original_price">
                          原价 ¥{{ cp.original_price }}
                        </span>
                      </div>
                    </div>
                    <div class="channel-price-right">
                      <div class="channel-price-display">
                        <span class="channel-price-currency">¥</span>
                        <span class="channel-price-value">{{ cp.price }}</span>
                        <span class="channel-price-unit">/晚</span>
                      </div>
                      <div class="channel-availability">
                        <span v-if="cp.available_count > 0" class="channel-available">
                          剩余 {{ cp.available_count }} 间
                        </span>
                        <span v-else class="channel-unavailable">
                          暂无空房
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
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
                <span class="price-num">{{ room.best_price || room.price }}</span>
                <span class="price-unit">/晚</span>
                <span v-if="room.best_price < room.price" class="price-discount">
                  省 ¥{{ Math.round(room.price - room.best_price) }}
                </span>
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
                <input type="date" v-model="bookingForm.checkIn" class="form-input" :min="today" :max="maxDate" />
              </div>
              <div class="form-group">
                <label class="form-label">离店日期</label>
                <input type="date" v-model="bookingForm.checkOut" class="form-input" :min="bookingForm.checkIn" :max="maxDate" />
              </div>
            </div>

            <div class="guest-section">
              <div class="guest-section-header">
                <label class="form-label">入住人信息</label>
                <button v-if="!showGuestForm" class="btn-toggle-guest" @click="openAddGuestForm">
                  + 新增常用入住人
                </button>
              </div>

              <div v-if="showGuestForm" class="guest-form-wrapper">
                <div class="guest-form-header">
                  <span class="guest-form-title">添加常用入住人</span>
                  <button class="btn-close-guest-form" @click="closeGuestForm">×</button>
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label class="form-label">姓名 <span class="required">*</span></label>
                    <input type="text" v-model="guestForm.name" class="form-input" placeholder="请输入姓名" />
                  </div>
                  <div class="form-group">
                    <label class="form-label">手机号 <span class="required">*</span></label>
                    <input type="tel" v-model="guestForm.phone" class="form-input" placeholder="请输入手机号" />
                  </div>
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label class="form-label">证件类型</label>
                    <select v-model="guestForm.id_type" class="form-input">
                      <option value="">请选择</option>
                      <option value="身份证">身份证</option>
                      <option value="护照">护照</option>
                      <option value="港澳通行证">港澳通行证</option>
                      <option value="台胞证">台胞证</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-label">证件号码</label>
                    <input type="text" v-model="guestForm.id_number" class="form-input" placeholder="请输入证件号码" />
                  </div>
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label class="checkbox-label">
                      <input type="checkbox" v-model="guestForm.is_default" />
                      <span>设为默认入住人</span>
                    </label>
                  </div>
                </div>
                <div class="guest-form-actions">
                  <button class="btn-cancel-guest" @click="closeGuestForm">取消</button>
                  <button class="btn-save-guest" :disabled="savingGuest" @click="saveGuest">
                    <span v-if="savingGuest">保存中...</span>
                    <span v-else>保存</span>
                  </button>
                </div>
              </div>

              <div v-if="guests.length > 0 && !showGuestForm" class="guest-list-wrapper">
                <p class="guest-list-hint">点击选择常用入住人（可自动填充）：</p>
                <div class="guest-list">
                  <div 
                    v-for="guest in guests" 
                    :key="guest.id" 
                    class="guest-item"
                    :class="{ 'selected': selectedGuest && selectedGuest.id === guest.id }"
                    @click="selectGuest(guest)"
                  >
                    <div class="guest-info">
                      <span class="guest-name">{{ guest.name }}</span>
                      <span class="guest-phone">{{ guest.phone }}</span>
                      <span v-if="guest.id_type && guest.id_number" class="guest-id">
                        {{ guest.id_type }}: {{ guest.id_number }}
                      </span>
                    </div>
                    <div class="guest-actions">
                      <span v-if="guest.is_default" class="default-badge">默认</span>
                      <span v-if="selectedGuest && selectedGuest.id === guest.id" class="check-icon">✓</span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="guest-input-section" v-if="!showGuestForm">
                <div class="guest-input-header">
                  <span class="guest-input-label">入住人信息</span>
                  <button v-if="selectedGuest" class="btn-unselect-guest" @click="clearSelectedGuest">
                    清除选择，手动填写
                  </button>
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label class="form-label">姓名 <span class="required">*</span></label>
                    <input type="text" v-model="bookingForm.guestName" class="form-input" placeholder="请输入姓名" />
                  </div>
                  <div class="form-group">
                    <label class="form-label">手机号 <span class="required">*</span></label>
                    <input type="tel" v-model="bookingForm.guestPhone" class="form-input" placeholder="请输入手机号" />
                  </div>
                </div>
                <div class="save-guest-option" v-if="!bookingForm.guestID">
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="saveAsGuest" />
                    <span>保存为常用入住人</span>
                  </label>
                </div>
                <div class="save-guest-option" v-if="saveAsGuest && !bookingForm.guestID">
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="setAsDefault" />
                    <span>设为默认入住人</span>
                  </label>
                </div>
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
import { hotelApi, orderApi, guestApi } from '../api'

export default {
  name: 'HotelDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()

    const today = new Date().toISOString().split('T')[0]
    const tomorrow = new Date(Date.now() + 86400000).toISOString().split('T')[0]
    const maxDate = new Date(Date.now() + 90 * 86400000).toISOString().split('T')[0]

    const loading = ref(true)
    const hotel = ref(null)
    const rooms = ref([])

    const showBookingModal = ref(false)
    const selectedRoom = ref(null)
    const showSuccessModal = ref(false)
    const submitting = ref(false)

    const checkInDate = ref(today)
    const checkOutDate = ref(tomorrow)
    const currentMonth = ref(new Date())

    const guests = ref([])
    const selectedGuest = ref(null)
    const showGuestForm = ref(false)
    const guestForm = ref({
      name: '',
      phone: '',
      id_type: '',
      id_number: '',
      is_default: false
    })
    const savingGuest = ref(false)
    const saveAsGuest = ref(false)
    const setAsDefault = ref(false)

    const bookingForm = ref({
      checkIn: today,
      checkOut: tomorrow,
      guestName: '',
      guestPhone: '',
      guestID: null
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

    const loadGuests = async () => {
      try {
        const res = await guestApi.getList()
        if (res.code === 200) {
          guests.value = res.data || []
        }
      } catch (error) {
        console.error('加载常用入住人失败:', error)
      }
    }

    const selectGuest = (guest) => {
      selectedGuest.value = guest
      bookingForm.value.guestName = guest.name
      bookingForm.value.guestPhone = guest.phone
      bookingForm.value.guestID = guest.id
    }

    const clearSelectedGuest = () => {
      selectedGuest.value = null
      bookingForm.value.guestName = ''
      bookingForm.value.guestPhone = ''
      bookingForm.value.guestID = null
    }

    const openAddGuestForm = () => {
      guestForm.value = {
        name: bookingForm.value.guestName || '',
        phone: bookingForm.value.guestPhone || '',
        id_type: '',
        id_number: '',
        is_default: false
      }
      showGuestForm.value = true
    }

    const closeGuestForm = () => {
      showGuestForm.value = false
      guestForm.value = {
        name: '',
        phone: '',
        id_type: '',
        id_number: '',
        is_default: false
      }
    }

    const saveGuest = async () => {
      if (!guestForm.value.name || !guestForm.value.phone) {
        alert('请填写姓名和手机号')
        return
      }

      savingGuest.value = true
      try {
        const res = await guestApi.create(guestForm.value)
        if (res.code === 200) {
          await loadGuests()
          const newGuest = {
            id: res.data.guest_id,
            ...guestForm.value
          }
          selectGuest(newGuest)
          closeGuestForm()
        } else {
          alert(res.message || '保存失败')
        }
      } catch (error) {
        console.error('保存入住人失败:', error)
        alert('保存失败，请稍后重试')
      } finally {
        savingGuest.value = false
      }
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

    const openBookingModal = async (room) => {
      selectedRoom.value = room
      bookingForm.value = {
        checkIn: checkInDate.value,
        checkOut: checkOutDate.value,
        guestName: '',
        guestPhone: '',
        guestID: null
      }
      selectedGuest.value = null
      showGuestForm.value = false
      saveAsGuest.value = false
      setAsDefault.value = false
      
      await loadGuests()
      
      const defaultGuest = guests.value.find(g => g.is_default)
      if (defaultGuest) {
        selectGuest(defaultGuest)
      }
      
      showBookingModal.value = true
    }

    const closeBookingModal = () => {
      showBookingModal.value = false
      selectedRoom.value = null
      selectedGuest.value = null
      showGuestForm.value = false
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
        if (saveAsGuest.value && !bookingForm.value.guestID && bookingForm.value.guestName && bookingForm.value.guestPhone) {
          try {
            const guestRes = await guestApi.create({
              name: bookingForm.value.guestName,
              phone: bookingForm.value.guestPhone,
              idType: '',
              idNumber: '',
              isDefault: setAsDefault.value
            })
            if (guestRes.code === 200) {
              await loadGuests()
            }
          } catch (e) {
            console.error('保存入住人失败:', e)
          }
        }

        const orderData = {
          hotel_id: hotel.value.id,
          room_id: selectedRoom.value.id,
          check_in: bookingForm.value.checkIn,
          check_out: bookingForm.value.checkOut
        }
        
        if (bookingForm.value.guestID) {
          orderData.guest_id = bookingForm.value.guestID
        } else {
          orderData.guest_name = bookingForm.value.guestName
          orderData.guest_phone = bookingForm.value.guestPhone
        }

        const res = await orderApi.create(orderData)

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

    const formatDate = (date) => {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }

    const getRoomPrice = (room, dateStr) => {
      if (!room) return 0
      const basePrice = room.best_price || room.price
      const date = new Date(dateStr)
      const dayOfWeek = date.getDay()
      
      let priceMultiplier = 1.0
      if (dayOfWeek === 5 || dayOfWeek === 6) {
        priceMultiplier = 1.15
      }
      if (dayOfWeek === 0) {
        priceMultiplier = 1.05
      }
      
      const dateSeed = dateStr.split('-').join('')
      const pseudoRandom = (parseInt(dateSeed) % 7) * 0.02
      
      return Math.round(basePrice * (priceMultiplier + pseudoRandom - 0.06))
    }

    const isDateSelectable = (dateStr) => {
      return dateStr >= today && dateStr <= maxDate
    }

    const isDateInRange = (dateStr) => {
      return dateStr >= checkInDate.value && dateStr <= checkOutDate.value
    }

    const isCheckInDate = (dateStr) => {
      return dateStr === checkInDate.value
    }

    const isCheckOutDate = (dateStr) => {
      return dateStr === checkOutDate.value
    }

    const getMinPriceForDate = (dateStr) => {
      if (!rooms.value || rooms.value.length === 0) return 0
      let minPrice = Infinity
      for (const room of rooms.value) {
        const price = getRoomPrice(room, dateStr)
        if (price < minPrice) {
          minPrice = price
        }
      }
      return minPrice === Infinity ? 0 : minPrice
    }

    const calendarDays = computed(() => {
      const year = currentMonth.value.getFullYear()
      const month = currentMonth.value.getMonth()
      
      const firstDay = new Date(year, month, 1)
      const lastDay = new Date(year, month + 1, 0)
      
      const days = []
      
      const firstDayOfWeek = firstDay.getDay()
      for (let i = 0; i < firstDayOfWeek; i++) {
        days.push({ day: '', date: '', isOtherMonth: true })
      }
      
      for (let d = 1; d <= lastDay.getDate(); d++) {
        const date = new Date(year, month, d)
        const dateStr = formatDate(date)
        const isSelectable = isDateSelectable(dateStr)
        const minPrice = isSelectable ? getMinPriceForDate(dateStr) : 0
        
        days.push({
          day: d,
          date: dateStr,
          isSelectable,
          isInRange: isDateInRange(dateStr),
          isCheckIn: isCheckInDate(dateStr),
          isCheckOut: isCheckOutDate(dateStr),
          minPrice
        })
      }
      
      return days
    })

    const currentMonthLabel = computed(() => {
      const year = currentMonth.value.getFullYear()
      const month = currentMonth.value.getMonth() + 1
      return `${year}年${month}月`
    })

    const prevMonth = () => {
      const newMonth = new Date(currentMonth.value)
      newMonth.setMonth(newMonth.getMonth() - 1)
      currentMonth.value = newMonth
    }

    const nextMonth = () => {
      const newMonth = new Date(currentMonth.value)
      newMonth.setMonth(newMonth.getMonth() + 1)
      currentMonth.value = newMonth
    }

    const selectDate = (day) => {
      if (!day.isSelectable) return
      
      const dateStr = day.date
      
      if (!checkInDate.value || dateStr < checkInDate.value) {
        checkInDate.value = dateStr
        const nextDay = new Date(dateStr)
        nextDay.setDate(nextDay.getDate() + 1)
        checkOutDate.value = formatDate(nextDay)
      } else if (dateStr === checkInDate.value) {
        const nextDay = new Date(dateStr)
        nextDay.setDate(nextDay.getDate() + 1)
        checkOutDate.value = formatDate(nextDay)
      } else if (dateStr > checkInDate.value) {
        checkOutDate.value = dateStr
      }
    }

    const calendarNights = computed(() => {
      if (!checkInDate.value || !checkOutDate.value) return 0
      const checkIn = new Date(checkInDate.value)
      const checkOut = new Date(checkOutDate.value)
      const diff = checkOut.getTime() - checkIn.getTime()
      return Math.max(0, Math.floor(diff / (1000 * 60 * 60 * 24)))
    })

    const canShowPrevMonth = computed(() => {
      const firstDayOfCurrentMonth = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth(), 1)
      const todayDate = new Date(today)
      todayDate.setDate(todayDate.getDate() - 1)
      return firstDayOfCurrentMonth > todayDate
    })

    const getSupplierShortName = (name) => {
      if (!name) return ''
      const shortNames = {
        '华住酒店集团': '华住',
        '锦江国际酒店集团': '锦江',
        '如家酒店集团': '如家',
        '万豪国际酒店集团': '万豪',
        '希尔顿酒店集团': '希尔顿',
        '洲际酒店集团': '洲际',
        '万达酒店集团': '万达',
        '开元酒店集团': '开元',
        '绿地酒店集团': '绿地',
        '模拟供应商A': '供应商A',
        '模拟供应商B': '供应商B',
        '模拟供应商C': '供应商C'
      }
      for (let key in shortNames) {
        if (name.includes(key) || name.includes(shortNames[key])) {
          return shortNames[key]
        }
      }
      if (name.includes('石基畅联')) {
        if (name.includes('万豪')) return '万豪'
        if (name.includes('希尔顿')) return '希尔顿'
        if (name.includes('洲际')) return '洲际'
        if (name.includes('开元')) return '开元'
        if (name.includes('万达')) return '万达'
        if (name.includes('绿地')) return '绿地'
        return '石基畅联'
      }
      return name.slice(0, 4)
    }

    onMounted(() => {
      loadHotelDetail()
    })

    return {
      today,
      maxDate,
      loading,
      hotel,
      rooms,
      checkInDate,
      checkOutDate,
      currentMonth,
      showBookingModal,
      selectedRoom,
      showSuccessModal,
      submitting,
      bookingForm,
      nights,
      totalPrice,
      canSubmit,
      calendarDays,
      currentMonthLabel,
      calendarNights,
      canShowPrevMonth,
      getSupplierShortName,
      getRoomPrice,
      prevMonth,
      nextMonth,
      selectDate,
      openBookingModal,
      closeBookingModal,
      submitBooking,
      closeSuccessModal,
      guests,
      selectedGuest,
      showGuestForm,
      guestForm,
      savingGuest,
      saveAsGuest,
      setAsDefault,
      selectGuest,
      clearSelectedGuest,
      openAddGuestForm,
      closeGuestForm,
      saveGuest
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
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.supplier-badge-large {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(26, 115, 232, 0.3);
}

.supplier-icon-large {
  font-size: 16px;
}

.supplier-badge-large .supplier-name {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.supplier-info-section {
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.supplier-info-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 8px;
}

.supplier-label {
  font-size: 13px;
  color: #666;
}

.supplier-name-large {
  font-size: 15px;
  font-weight: 600;
  color: #1a73e8;
}

.supplier-info-details {
  display: flex;
  gap: 12px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.priority-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.priority-icon {
  font-size: 12px;
}

.price-control-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.price-control-icon {
  font-size: 12px;
}

.supplier-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.5;
  margin: 0;
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

.date-filter-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.date-filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.date-summary {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.date-range-display {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  border-radius: 8px;
}

.date-display-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.date-display-label {
  font-size: 12px;
  color: #999;
}

.date-display-value {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.date-arrow {
  font-size: 18px;
  color: #1a73e8;
}

.nights-badge {
  padding: 8px 16px;
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  color: #fff;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
}

.calendar-wrapper {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
}

.calendar-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fafafa;
  border-bottom: 1px solid #e0e0e0;
}

.nav-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 50%;
  font-size: 18px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn:hover:not(:disabled) {
  background: #1a73e8;
  border-color: #1a73e8;
  color: #fff;
}

.nav-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.month-label {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  padding: 12px 0;
  background: #fafafa;
  border-bottom: 1px solid #e0e0e0;
}

.weekday {
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: #666;
}

.calendar-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  background: #e0e0e0;
  padding: 1px;
}

.calendar-day {
  background: #fff;
  padding: 12px 8px;
  text-align: center;
  cursor: pointer;
  min-height: 60px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  transition: all 0.2s;
  position: relative;
}

.calendar-day:hover:not(.other-month):not(.disabled) {
  background: #f5f7fa;
}

.calendar-day.other-month {
  background: #fafafa;
  cursor: default;
  opacity: 0.3;
}

.calendar-day.disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.calendar-day.disabled .day-number {
  color: #ccc;
}

.calendar-day.in-range {
  background: rgba(26, 115, 232, 0.1);
}

.calendar-day.check-in,
.calendar-day.check-out {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
}

.calendar-day.check-in .day-number,
.calendar-day.check-out .day-number,
.calendar-day.check-in .day-price,
.calendar-day.check-out .day-price {
  color: #fff;
}

.day-number {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.day-price {
  font-size: 12px;
  font-weight: 500;
  color: #e74c3c;
}

.day-disabled-text {
  font-size: 11px;
  color: #ccc;
}

.date-hint {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  flex-wrap: wrap;
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.hint-icon {
  width: 16px;
  height: 16px;
  border-radius: 4px;
}

.hint-blue {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
}

.hint-light-blue {
  background: rgba(26, 115, 232, 0.1);
  border: 1px solid rgba(26, 115, 232, 0.3);
}

.hint-gray {
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
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

.section-subtitle {
  font-size: 13px;
  font-weight: 500;
  color: #1a73e8;
  background: rgba(26, 115, 232, 0.1);
  padding: 4px 10px;
  border-radius: 4px;
  margin-left: 8px;
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

.room-header-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.best-price-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}

.best-price-icon {
  font-size: 12px;
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

.channel-prices-section {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
  border: 1px solid #e2e8f0;
}

.channel-prices-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.channel-prices-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.channel-prices-count {
  font-size: 12px;
  color: #666;
}

.channel-prices-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.channel-price-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: #fff;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s;
}

.channel-price-item:hover {
  border-color: #1a73e8;
  box-shadow: 0 2px 8px rgba(26, 115, 232, 0.15);
}

.channel-price-item.best-price {
  border-color: #10b981;
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
}

.channel-price-item.out-of-stock {
  opacity: 0.6;
  background: #f8fafc;
}

.channel-price-info {
  flex: 1;
}

.channel-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.channel-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.best-price-badge {
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  padding: 2px 8px;
  border-radius: 10px;
}

.out-of-stock-badge {
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #64748b 0%, #475569 100%);
  padding: 2px 8px;
  border-radius: 10px;
}

.channel-details {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.channel-priority {
  font-size: 12px;
  color: #f59e0b;
}

.channel-original-price {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.channel-price-right {
  text-align: right;
}

.channel-price-display {
  display: flex;
  align-items: baseline;
  gap: 2px;
  margin-bottom: 4px;
}

.channel-price-currency {
  font-size: 12px;
  color: #e74c3c;
}

.channel-price-value {
  font-size: 20px;
  font-weight: 600;
  color: #e74c3c;
}

.channel-price-unit {
  font-size: 11px;
  color: #999;
}

.channel-availability {
  font-size: 12px;
}

.channel-available {
  color: #10b981;
}

.channel-unavailable {
  color: #64748b;
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

.price-discount {
  font-size: 12px;
  font-weight: 600;
  color: #10b981;
  background: rgba(16, 185, 129, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
  margin-left: 8px;
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

.guest-section {
  margin-top: 20px;
}

.guest-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.btn-toggle-guest {
  padding: 4px 12px;
  background: transparent;
  border: 1px solid #1a73e8;
  color: #1a73e8;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-toggle-guest:hover {
  background: #1a73e8;
  color: #fff;
}

.guest-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.guest-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.guest-item:hover {
  border-color: #1a73e8;
  background: rgba(26, 115, 232, 0.02);
}

.guest-item.selected {
  border-color: #1a73e8;
  background: rgba(26, 115, 232, 0.05);
}

.guest-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.guest-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.guest-phone {
  font-size: 13px;
  color: #666;
}

.guest-id {
  font-size: 12px;
  color: #999;
}

.guest-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.default-badge {
  padding: 2px 8px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: #fff;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
}

.check-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1a73e8;
  color: #fff;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
}

.guest-form-wrapper {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.guest-form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.guest-form-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.btn-close-guest-form {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 18px;
  color: #999;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
}

.btn-close-guest-form:hover {
  background: #e2e8f0;
  color: #666;
}

.required {
  color: #e74c3c;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.guest-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.btn-cancel-guest {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid #ddd;
  color: #666;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel-guest:hover {
  border-color: #999;
}

.btn-save-guest {
  padding: 8px 20px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-save-guest:hover:not(:disabled) {
  background: #1557b0;
}

.btn-save-guest:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.guest-empty {
  text-align: center;
  padding: 24px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px dashed #e2e8f0;
}

.guest-empty-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}

.guest-empty-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

.btn-add-guest {
  padding: 8px 20px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-add-guest:hover {
  background: #1557b0;
}

.guest-or {
  font-size: 13px;
  color: #999;
}

.btn-manual-input {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid #1a73e8;
  color: #1a73e8;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-manual-input:hover {
  background: #1a73e8;
  color: #fff;
}

.manual-input-section {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #eee;
}

.manual-input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.manual-input-label {
  font-size: 13px;
  color: #999;
}

.save-guest-option {
  margin-top: 8px;
}

.guest-list-wrapper {
  margin-bottom: 16px;
}

.guest-list-hint {
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.guest-input-section {
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.guest-input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.guest-input-label {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.btn-unselect-guest {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid #ddd;
  color: #666;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-unselect-guest:hover {
  border-color: #999;
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
