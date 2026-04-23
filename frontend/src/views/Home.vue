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
            <input type="date" v-model="checkIn" class="search-input" :min="today" :max="maxDate" />
          </div>
          <div class="search-item">
            <label class="search-label">离店日期</label>
            <input type="date" v-model="checkOut" class="search-input" :min="checkIn" :max="maxDate" />
          </div>
          <button class="search-btn" @click="searchHotels">
            <span class="search-icon">🔍</span>
            搜索
          </button>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="filter-section">
        <div class="filter-header">
          <h3 class="filter-title">筛选条件</h3>
          <button class="reset-btn" @click="resetFilters">重置</button>
        </div>

        <div class="filter-group">
          <div class="filter-group-header" @click="toggleFilter('price')">
            <span class="filter-group-title">价格范围</span>
            <span class="filter-toggle">{{ filterExpanded.price ? '−' : '+' }}</span>
          </div>
          <div v-show="filterExpanded.price" class="filter-group-content">
            <div class="price-slider-container">
              <div class="price-range-display">
                <span class="price-value">¥{{ priceMin }}</span>
                <span class="price-separator">至</span>
                <span class="price-value">¥{{ priceMax }}</span>
              </div>
              <div class="slider-wrapper">
                <div class="slider-track">
                  <div class="slider-fill" :style="sliderFillStyle"></div>
                </div>
                <input 
                  type="range" 
                  v-model="priceMin" 
                  :min="filterOptions.price_range?.min || 100" 
                  :max="priceMax" 
                  class="slider-input slider-input-min"
                  @input="updateSlider"
                />
                <input 
                  type="range" 
                  v-model="priceMax" 
                  :min="priceMin" 
                  :max="filterOptions.price_range?.max || 5000" 
                  class="slider-input slider-input-max"
                  @input="updateSlider"
                />
              </div>
              <div class="price-presets">
                <button 
                  v-for="preset in pricePresets" 
                  :key="preset.label"
                  class="price-preset-btn"
                  :class="{ active: isPricePresetActive(preset) }"
                  @click="setPricePreset(preset)"
                >
                  {{ preset.label }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="filter-group">
          <div class="filter-group-header" @click="toggleFilter('rating')">
            <span class="filter-group-title">酒店星级</span>
            <span class="filter-toggle">{{ filterExpanded.rating ? '−' : '+' }}</span>
          </div>
          <div v-show="filterExpanded.rating" class="filter-group-content">
            <div class="rating-options">
              <label 
                v-for="rating in ratingOptions" 
                :key="rating.value"
                class="rating-option"
                :class="{ active: selectedRating === rating.value }"
                @click="selectRating(rating.value)"
              >
                <span class="rating-stars">{{ rating.label }}</span>
                <span class="rating-text">{{ rating.text }}</span>
              </label>
            </div>
          </div>
        </div>

        <div class="filter-group">
          <div class="filter-group-header" @click="toggleFilter('bedType')">
            <span class="filter-group-title">床型</span>
            <span class="filter-toggle">{{ filterExpanded.bedType ? '−' : '+' }}</span>
          </div>
          <div v-show="filterExpanded.bedType" class="filter-group-content">
            <div class="bed-type-options">
              <label 
                class="bed-type-option"
                :class="{ active: selectedBedType === '' }"
                @click="selectedBedType = ''"
              >
                <span class="bed-type-text">全部</span>
              </label>
              <label 
                v-for="bedType in filterOptions.bed_types || []" 
                :key="bedType"
                class="bed-type-option"
                :class="{ active: selectedBedType === bedType }"
                @click="selectedBedType = bedType"
              >
                <span class="bed-type-icon">{{ getBedTypeIcon(bedType) }}</span>
                <span class="bed-type-text">{{ bedType }}</span>
              </label>
            </div>
          </div>
        </div>

        <div class="filter-group">
          <div class="filter-group-header" @click="toggleFilter('amenities')">
            <span class="filter-group-title">酒店设施</span>
            <span class="filter-toggle">{{ filterExpanded.amenities ? '−' : '+' }}</span>
          </div>
          <div v-show="filterExpanded.amenities" class="filter-group-content">
            <div class="amenities-options">
              <label 
                v-for="amenity in (filterOptions.amenities || []).slice(0, 12)" 
                :key="amenity"
                class="amenity-option"
                :class="{ active: selectedAmenities.includes(amenity) }"
                @click="toggleAmenity(amenity)"
              >
                <span class="amenity-icon">{{ getAmenityIcon(amenity) }}</span>
                <span class="amenity-text">{{ amenity }}</span>
              </label>
            </div>
          </div>
        </div>

        <button class="apply-filter-btn" @click="applyFilters">
          应用筛选
        </button>
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
              <div v-if="hotel.supplier" class="hotel-supplier-badge">
                <span class="supplier-icon">🏢</span>
                <span class="supplier-name">{{ getSupplierShortName(hotel.supplier.name) }}</span>
              </div>
            </div>
            <div class="hotel-info">
              <h3 class="hotel-name">
                {{ hotel.name }}
                <span v-if="hotel.supplier" class="supplier-tag">
                  {{ getSupplierShortName(hotel.supplier.name) }}
                </span>
              </h3>
              <p class="hotel-location">
                <span class="location-icon">📍</span>
                {{ hotel.city }} · {{ hotel.address }}
              </p>
              <div v-if="hotel.supplier" class="hotel-supplier-info">
                <span class="supplier-label">来源：</span>
                <span class="supplier-value">{{ hotel.supplier.name }}</span>
                <span v-if="hotel.supplier.priority > 0" class="priority-tag">
                  优先级 {{ hotel.supplier.priority }}
                </span>
              </div>
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
          <p class="empty-text">暂无符合条件的酒店</p>
          <button class="reset-search-btn" @click="resetFilters">清除筛选条件</button>
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
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { hotelApi } from '../api'

export default {
  name: 'Home',
  setup() {
    const router = useRouter()
    
    const today = new Date().toISOString().split('T')[0]
    const tomorrow = new Date(Date.now() + 86400000).toISOString().split('T')[0]
    const maxDate = new Date(Date.now() + 90 * 86400000).toISOString().split('T')[0]
    
    const cities = ref([])
    const selectedCity = ref('')
    const checkIn = ref(today)
    const checkOut = ref(tomorrow)
    
    const hotels = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(6)
    
    const filterOptions = ref({
      bed_types: [],
      amenities: [],
      price_range: { min: 100, max: 5000 }
    })
    
    const filterExpanded = ref({
      price: true,
      rating: true,
      bedType: true,
      amenities: true
    })
    
    const priceMin = ref(100)
    const priceMax = ref(5000)
    const selectedRating = ref(0)
    const selectedBedType = ref('')
    const selectedAmenities = ref([])
    
    const ratingOptions = [
      { value: 0, label: '不限', text: '全部星级' },
      { value: 3, label: '⭐⭐⭐', text: '三星及以上' },
      { value: 4, label: '⭐⭐⭐⭐', text: '四星及以上' },
      { value: 4.5, label: '⭐⭐⭐⭐⭐', text: '五星及以上' }
    ]
    
    const pricePresets = [
      { label: '¥300以下', min: 0, max: 300 },
      { label: '¥300-800', min: 300, max: 800 },
      { label: '¥800-1500', min: 800, max: 1500 },
      { label: '¥1500以上', min: 1500, max: 999999 }
    ]
    
    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
    
    const sliderFillStyle = computed(() => {
      const min = filterOptions.value.price_range?.min || 100
      const max = filterOptions.value.price_range?.max || 5000
      const range = max - min
      const left = ((priceMin.value - min) / range) * 100
      const right = ((priceMax.value - min) / range) * 100
      return {
        left: `${left}%`,
        right: `${100 - right}%`
      }
    })
    
    const toggleFilter = (filter) => {
      filterExpanded.value[filter] = !filterExpanded.value[filter]
    }
    
    const updateSlider = () => {
    }
    
    const isPricePresetActive = (preset) => {
      return priceMin.value === preset.min && priceMax.value === preset.max
    }
    
    const setPricePreset = (preset) => {
      priceMin.value = preset.min
      priceMax.value = preset.max
    }
    
    const selectRating = (value) => {
      selectedRating.value = selectedRating.value === value ? 0 : value
    }
    
    const toggleAmenity = (amenity) => {
      const index = selectedAmenities.value.indexOf(amenity)
      if (index > -1) {
        selectedAmenities.value.splice(index, 1)
      } else {
        selectedAmenities.value.push(amenity)
      }
    }
    
    const getBedTypeIcon = (bedType) => {
      if (bedType.includes('大')) return '🛏️'
      if (bedType.includes('双')) return '🛏️🛏️'
      if (bedType.includes('单')) return '🛏️'
      return '🛏️'
    }
    
    const getAmenityIcon = (amenity) => {
      if (amenity.includes('WiFi')) return '📶'
      if (amenity.includes('空调')) return '❄️'
      if (amenity.includes('电视')) return '📺'
      if (amenity.includes('热水')) return '🚿'
      if (amenity.includes('酒廊')) return '🍸'
      if (amenity.includes('泳池')) return '🏊'
      if (amenity.includes('管家')) return '👔'
      if (amenity.includes('迷你')) return '🧊'
      if (amenity.includes('办公')) return '💼'
      if (amenity.includes('阳台')) return '🌅'
      if (amenity.includes('智能')) return '🤖'
      if (amenity.includes('行政')) return '👑'
      return '✨'
    }
    
    const resetFilters = () => {
      priceMin.value = filterOptions.value.price_range?.min || 100
      priceMax.value = filterOptions.value.price_range?.max || 5000
      selectedRating.value = 0
      selectedBedType.value = ''
      selectedAmenities.value = []
      selectedCity.value = ''
      page.value = 1
      loadHotels()
    }
    
    const applyFilters = () => {
      page.value = 1
      loadHotels()
    }
    
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
    
    const loadFilterOptions = async () => {
      try {
        const res = await hotelApi.getFilterOptions()
        if (res.code === 200) {
          filterOptions.value = res.data
          priceMin.value = res.data.price_range?.min || 100
          priceMax.value = res.data.price_range?.max || 5000
        }
      } catch (error) {
        console.error('加载筛选选项失败:', error)
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
        if (priceMin.value > 0) {
          params.min_price = priceMin.value
        }
        if (priceMax.value < 999999) {
          params.max_price = priceMax.value
        }
        if (selectedRating.value > 0) {
          params.rating_min = selectedRating.value
        }
        if (selectedBedType.value) {
          params.bed_type = selectedBedType.value
        }
        if (selectedAmenities.value.length > 0) {
          params.amenities = selectedAmenities.value[0]
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
      loadCities()
      loadFilterOptions()
      loadHotels()
    })

    return {
      today,
      maxDate,
      cities,
      selectedCity,
      checkIn,
      checkOut,
      hotels,
      total,
      page,
      pageSize,
      totalPages,
      filterOptions,
      filterExpanded,
      priceMin,
      priceMax,
      selectedRating,
      selectedBedType,
      selectedAmenities,
      ratingOptions,
      pricePresets,
      sliderFillStyle,
      toggleFilter,
      updateSlider,
      isPricePresetActive,
      setPricePreset,
      selectRating,
      toggleAmenity,
      getBedTypeIcon,
      getAmenityIcon,
      resetFilters,
      applyFilters,
      searchHotels,
      changePage,
      goToHotel,
      getSupplierShortName
    }
  }
}
</script>

<style scoped>
.home {
  max-width: 1400px;
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

.main-content {
  display: flex;
  gap: 24px;
  margin-bottom: 40px;
}

.filter-section {
  width: 280px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  padding: 20px;
  height: fit-content;
  position: sticky;
  top: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.reset-btn {
  font-size: 13px;
  color: #1a73e8;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background 0.2s;
}

.reset-btn:hover {
  background: rgba(26, 115, 232, 0.1);
}

.filter-group {
  margin-bottom: 16px;
}

.filter-group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  cursor: pointer;
  user-select: none;
}

.filter-group-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.filter-toggle {
  font-size: 18px;
  color: #999;
  font-weight: 300;
}

.filter-group-content {
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.price-slider-container {
  padding: 8px 0;
}

.price-range-display {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.price-value {
  font-size: 16px;
  font-weight: 600;
  color: #1a73e8;
}

.price-separator {
  font-size: 14px;
  color: #999;
}

.slider-wrapper {
  position: relative;
  height: 6px;
  margin: 20px 0;
}

.slider-track {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 6px;
  background: #e0e0e0;
  border-radius: 3px;
}

.slider-fill {
  position: absolute;
  top: 0;
  height: 6px;
  background: linear-gradient(90deg, #1a73e8, #0d47a1);
  border-radius: 3px;
}

.slider-input {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 100%;
  height: 6px;
  background: transparent;
  -webkit-appearance: none;
  appearance: none;
  pointer-events: none;
}

.slider-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 20px;
  height: 20px;
  background: #fff;
  border: 2px solid #1a73e8;
  border-radius: 50%;
  cursor: pointer;
  pointer-events: auto;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s, box-shadow 0.2s;
}

.slider-input::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(26, 115, 232, 0.3);
}

.slider-input::-moz-range-thumb {
  width: 20px;
  height: 20px;
  background: #fff;
  border: 2px solid #1a73e8;
  border-radius: 50%;
  cursor: pointer;
  pointer-events: auto;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.price-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.price-preset-btn {
  padding: 6px 12px;
  font-size: 12px;
  color: #666;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.price-preset-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.price-preset-btn.active {
  background: linear-gradient(135deg, #1a73e8, #0d47a1);
  color: #fff;
  border-color: transparent;
}

.rating-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rating-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.rating-option:hover {
  background: #f5f5f5;
}

.rating-option.active {
  background: rgba(26, 115, 232, 0.1);
}

.rating-stars {
  font-size: 14px;
}

.rating-text {
  font-size: 13px;
  color: #666;
}

.bed-type-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.bed-type-option {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.bed-type-option:hover {
  border-color: #1a73e8;
}

.bed-type-option.active {
  background: linear-gradient(135deg, #1a73e8, #0d47a1);
  color: #fff;
  border-color: transparent;
}

.bed-type-icon {
  font-size: 14px;
}

.bed-type-text {
  font-size: 13px;
}

.amenities-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.amenity-option {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 12px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.amenity-option:hover {
  border-color: #1a73e8;
}

.amenity-option.active {
  background: rgba(26, 115, 232, 0.1);
  border-color: #1a73e8;
}

.amenity-icon {
  font-size: 14px;
}

.amenity-text {
  font-size: 12px;
  color: #666;
}

.apply-filter-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #1a73e8, #0d47a1);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 12px;
}

.apply-filter-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(26, 115, 232, 0.3);
}

.hotel-section {
  flex: 1;
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
  grid-template-columns: repeat(2, 1fr);
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

.hotel-supplier-badge {
  position: absolute;
  bottom: 12px;
  left: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(26, 115, 232, 0.3);
}

.supplier-icon {
  font-size: 12px;
}

.supplier-name {
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.hotel-info {
  padding: 16px;
}

.hotel-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.supplier-tag {
  font-size: 11px;
  font-weight: 500;
  color: #1a73e8;
  background: rgba(26, 115, 232, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid rgba(26, 115, 232, 0.2);
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

.hotel-supplier-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 12px;
  flex-wrap: wrap;
}

.supplier-label {
  color: #999;
}

.supplier-value {
  color: #333;
  font-weight: 500;
}

.priority-tag {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
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
  margin-bottom: 16px;
}

.reset-search-btn {
  padding: 10px 24px;
  background: linear-gradient(135deg, #1a73e8, #0d47a1);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.reset-search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(26, 115, 232, 0.3);
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

@media (max-width: 1200px) {
  .main-content {
    flex-direction: column;
  }
  
  .filter-section {
    width: 100%;
    position: static;
  }
  
  .amenities-options {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 1024px) {
  .hotel-list {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .amenities-options {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .search-box {
    flex-direction: column;
  }
  
  .hotel-list {
    grid-template-columns: 1fr;
  }
  
  .amenities-options {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
