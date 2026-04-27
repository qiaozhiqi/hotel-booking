import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(
  (config) => {
    const user = localStorage.getItem('user')
    if (user) {
      const userData = JSON.parse(user)
      config.headers['X-User-ID'] = userData.id
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    return Promise.reject(error)
  }
)

export const userApi = {
  login: (data) => api.post('/login', data),
  register: (data) => api.post('/register', data),
  getInfo: () => api.get('/user')
}

export const hotelApi = {
  getList: (params) => api.get('/hotels', { params }),
  getDetail: (id) => api.get(`/hotels/${id}`),
  getCities: () => api.get('/cities'),
  getFilterOptions: () => api.get('/filter-options')
}

export const orderApi = {
  create: (data) => api.post('/orders', data),
  getList: (params) => api.get('/orders', { params }),
  getDetail: (id) => api.get(`/orders/${id}`),
  cancel: (id) => api.post(`/orders/${id}/cancel`)
}

export const guestApi = {
  getList: () => api.get('/guests'),
  getDetail: (id) => api.get(`/guests/${id}`),
  create: (data) => api.post('/guests', data),
  update: (id, data) => api.put(`/guests/${id}`, data),
  delete: (id) => api.delete(`/guests/${id}`),
  setDefault: (id) => api.post(`/guests/${id}/default`)
}

export const favoriteApi = {
  create: (data) => api.post('/favorites', data),
  delete: (hotelId) => api.delete(`/favorites/${hotelId}`),
  getList: (params) => api.get('/favorites', { params }),
  checkStatus: (hotelId) => api.get(`/favorites/check/${hotelId}`)
}

export const invoiceApi = {
  create: (data) => api.post('/invoices', data),
  getList: (params) => api.get('/invoices', { params }),
  getDetail: (id) => api.get(`/invoices/${id}`),
  getAvailableOrders: () => api.get('/invoices/orders/available')
}

export default api
