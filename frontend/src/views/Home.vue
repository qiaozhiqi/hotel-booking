<template>
  <div class="home">
    <div class="search-section">
      <div class="search-content">
        <h1 class="search-title">找到理想的酒店</h1>
        <p class="search-subtitle">精选全球优质酒店，为您提供舒适的住宿体验</p>
        <div class="search-box">
          <div class="search-item">
            <label class="search-label">目的地</label>
            <select v-model="selectedCity" class="search-select">
              <option value="">全部城市</option>
              <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
            </select>
          </div>
          <div class="search-item">
            <label class="search-label">入住日期</label>
            <input type="date" v-model="checkIn" class="search-input" :min="today" />
          </div>
          <div class="search-item">
            <label class="search-label">离店日期</label>
            <input type="date" v-model="checkOut" class="search-input" :min="checkIn" />
          </div>
          <button class="search-btn" @click="searchHotels">
            <span class="search-icon">🔍</span>
            搜索
          </button>
        </div>
      </div>
    </div>

    <div class="hotel-section">
      <div class="section-header">
        <h2 class="section-title">推荐酒店</h2>
        <p class="section-subtitle">共找到 {{ total }} 家酒店</p>
      </div>
      
      <div class="hotel-list" v-if="hotels.length > 0">
        <div 
          v-for="hotel in hotels" 
          :key="hotel.id" 
          class="hotel-card"
          @click="goToHotel(hotel.id)"
        >
          <div class="hotel-image">
            <img :src="hotel.image_url" :alt="hotel.name" />
            <div class="hotel-rating">
              <span class="rating-star">⭐</span>
              <span class="rating-value">{{ hotel.rating }}</span>
            </div>
          </div>
          <div class="hotel-info">
            <h3 class="hotel-name">{{ hotel.name }}</h3>
            <p class="hotel-location">
              <span class="location-icon">📍</span>
              {{ hotel.city }} · {{ hotel.address }}
            </p>
            <p class="hotel-desc" v-if="hotel.description">{{ hotel.description }}</p>
            <div class="hotel-price">
              <span class="price-label">起价</span>
              <span class="price-value">{{ hotel.price_range }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="empty-state" v-else>
        <div class="empty-icon">🏨</div>
        <p class="empty-text">暂无酒店数据</p>
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
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { hotelApi } from '../api'

export default {
  name: 'Home',
  setup() {
    const router = useRouter()
    
    const today = new Date().toISOString().split('T')[0]
    const tomorrow = new Date(Date.now() + 86400000).toISOString().split('T')[0]
    
    const cities = ref([])
    const selectedCity = ref('')
    const checkIn = ref(today)
    const checkOut = ref(tomorrow)
    
    const hotels = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(6)
    
    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

    const loadCities = async () => {
      try {
        const res = await hotelApi.getCities()
        if (res.code === 200) {
          cities.value = res.data
        }
      } catch (error) {
        console.error('加载城市列表失败:', error)
      }
    }

    const loadHotels = async () => {
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        if (selectedCity.value) {
          params.city = selectedCity.value
        }
        
        const res = await hotelApi.getList(params)
        if (res.code === 200) {
          hotels.value = res.data.hotels || []
          total.value = res.data.total || 0
        }
      } catch (error) {
        console.error('加载酒店列表失败:', error)
      }
    }

    const searchHotels = () => {
      page.value = 1
      loadHotels()
    }

    const changePage = (newPage) => {
      page.value = newPage
      loadHotels()
    }

    const goToHotel = (id) => {
      router.push(`/hotel/${id}`)
    }

    onMounted(() => {
      loadCities()
      loadHotels()
    })

    return {
      today,
      cities,
      selectedCity,
      checkIn,
      checkOut,
      hotels,
      total,
      page,
      pageSize,
      totalPages,
      searchHotels,
      changePage,
      goToHotel
    }
  }
}
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.search-section {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  margin: 30px -20px;
  padding: 60px 20px;
  border-radius: 12px;
}

.search-content {
  max-width: 900px;
  margin: 0 auto;
  text-align: center;
}

.search-title {
  color: #fff;
  font-size: 36px;
  font-weight: 600;
  margin-bottom: 12px;
}

.search-subtitle {
  color: rgba(255, 255, 255, 0.8);
  font-size: 16px;
  margin-bottom: 40px;
}

.search-box {
  display: flex;
  gap: 16px;
  background: #fff;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.search-item {
  flex: 1;
  text-align: left;
}

.search-label {
  display: block;
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.search-select,
.search-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
}

.search-select:focus,
.search-input:focus {
  border-color: #1a73e8;
}

.search-btn {
  align-self: flex-end;
  display: flex;
  align-items: center;
  gap: 8px;
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

.search-btn:hover {
  background: #1557b0;
}

.search-icon {
  font-size: 18px;
}

.hotel-section {
  margin-bottom: 40px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.section-subtitle {
  font-size: 14px;
  color: #999;
}

.hotel-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.hotel-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.hotel-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.hotel-image {
  position: relative;
  height: 200px;
  background: #f5f5f5;
}

.hotel-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hotel-rating {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.rating-star {
  font-size: 14px;
}

.rating-value {
  font-size: 14px;
  font-weight: 600;
  color: #f59e0b;
}

.hotel-info {
  padding: 16px;
}

.hotel-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hotel-location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.location-icon {
  font-size: 14px;
}

.hotel-desc {
  font-size: 13px;
  color: #999;
  line-height: 1.5;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.hotel-price {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.price-label {
  font-size: 13px;
  color: #999;
}

.price-value {
  font-size: 20px;
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
  font-size: 16px;
  color: #999;
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

@media (max-width: 1024px) {
  .hotel-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .search-box {
    flex-direction: column;
  }
  
  .hotel-list {
    grid-template-columns: 1fr;
  }
}
</style>
