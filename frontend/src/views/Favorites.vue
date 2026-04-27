<template>
  <div class="favorites-page">
    <div class="favorites-header">
      <h1 class="page-title">我的收藏</h1>
    </div>

    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <template v-else-if="favorites.length > 0">
      <div class="favorites-list">
        <div 
          v-for="favorite in favorites" 
          :key="favorite.id" 
          class="favorite-card"
        >
          <router-link :to="`/hotel/${favorite.hotel_id}`" class="hotel-link">
            <div class="hotel-image">
              <img :src="favorite.image_url" :alt="favorite.hotel_name" />
            </div>
            <div class="hotel-info">
              <h3 class="hotel-name">{{ favorite.hotel_name }}</h3>
              <div class="hotel-location">
                <span class="location-icon">📍</span>
                <span>{{ favorite.city }} · {{ favorite.address }}</span>
              </div>
              <div class="hotel-rating">
                <span class="rating-star">⭐</span>
                <span class="rating-value">{{ favorite.rating }}</span>
              </div>
              <div class="hotel-price">
                <span class="price-range">{{ favorite.price_range }}</span>
              </div>
            </div>
          </router-link>
          <button 
            class="btn-remove-favorite"
            @click="removeFavorite(favorite.hotel_id, favorite.hotel_name)"
          >
            <span class="remove-icon">❤️</span>
            取消收藏
          </button>
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
      <div class="empty-icon">💝</div>
      <p class="empty-text">暂无收藏的酒店</p>
      <router-link to="/" class="btn-go-explore">去探索酒店</router-link>
    </div>

    <div v-if="showRemoveModal" class="modal-overlay" @click.self="closeRemoveModal">
      <div class="modal-content confirm-modal">
        <div class="modal-icon">💔</div>
        <h3 class="modal-title">确认取消收藏？</h3>
        <p class="modal-text">您确定要取消收藏 "{{ hotelToRemoveName }}" 吗？</p>
        <div class="modal-actions">
          <button class="btn-secondary" @click="closeRemoveModal">再想想</button>
          <button class="btn-danger" @click="confirmRemove" :disabled="removing">
            <span v-if="removing">处理中...</span>
            <span v-else>确认取消</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { favoriteApi } from '../api'

export default {
  name: 'Favorites',
  setup() {
    const loading = ref(true)
    const favorites = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(10)

    const showRemoveModal = ref(false)
    const hotelToRemoveId = ref(null)
    const hotelToRemoveName = ref('')
    const removing = ref(false)

    const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

    const loadFavorites = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        const res = await favoriteApi.getList(params)
        if (res.code === 200) {
          favorites.value = res.data.favorites || []
          total.value = res.data.total || 0
        }
      } catch (error) {
        console.error('加载收藏列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    const removeFavorite = (hotelId, hotelName) => {
      hotelToRemoveId.value = hotelId
      hotelToRemoveName.value = hotelName
      showRemoveModal.value = true
    }

    const closeRemoveModal = () => {
      showRemoveModal.value = false
      hotelToRemoveId.value = null
      hotelToRemoveName.value = ''
    }

    const confirmRemove = async () => {
      if (!hotelToRemoveId.value) return
      
      removing.value = true
      try {
        const res = await favoriteApi.delete(hotelToRemoveId.value)
        if (res.code === 200) {
          closeRemoveModal()
          loadFavorites()
        } else {
          alert(res.message || '取消收藏失败')
        }
      } catch (error) {
        console.error('取消收藏失败:', error)
        alert('取消收藏失败，请稍后重试')
      } finally {
        removing.value = false
      }
    }

    const changePage = (newPage) => {
      page.value = newPage
      loadFavorites()
    }

    onMounted(() => {
      loadFavorites()
    })

    return {
      loading,
      favorites,
      total,
      page,
      pageSize,
      totalPages,
      showRemoveModal,
      hotelToRemoveId,
      hotelToRemoveName,
      removing,
      removeFavorite,
      closeRemoveModal,
      confirmRemove,
      changePage
    }
  }
}
</script>

<style scoped>
.favorites-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 30px 20px;
}

.favorites-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
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

.favorites-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.favorite-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: stretch;
}

.hotel-link {
  display: flex;
  flex: 1;
  text-decoration: none;
  color: inherit;
}

.hotel-image {
  width: 200px;
  min-width: 200px;
  height: 150px;
  overflow: hidden;
}

.hotel-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hotel-info {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.hotel-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.hotel-location {
  font-size: 13px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.location-icon {
  font-size: 14px;
}

.hotel-rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.rating-star {
  font-size: 14px;
}

.rating-value {
  font-size: 14px;
  font-weight: 600;
  color: #f59e0b;
}

.hotel-price {
  margin-top: auto;
}

.price-range {
  font-size: 16px;
  font-weight: 600;
  color: #e74c3c;
}

.btn-remove-favorite {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid #e74c3c;
  color: #e74c3c;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  margin: 16px;
  align-self: flex-start;
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-remove-favorite:hover {
  background: #e74c3c;
  color: #fff;
}

.remove-icon {
  font-size: 14px;
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

.btn-go-explore {
  display: inline-block;
  padding: 12px 32px;
  background: #1a73e8;
  color: #fff;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  transition: background 0.2s;
}

.btn-go-explore:hover {
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
  .favorite-card {
    flex-direction: column;
  }

  .hotel-link {
    flex-direction: column;
  }

  .hotel-image {
    width: 100%;
    min-width: auto;
    height: 200px;
  }

  .btn-remove-favorite {
    margin: 16px;
    align-self: stretch;
    justify-content: center;
  }
}
</style>
