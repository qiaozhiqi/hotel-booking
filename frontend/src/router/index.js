import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import HotelDetail from '../views/HotelDetail.vue'
import Orders from '../views/Orders.vue'
import Login from '../views/Login.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/hotel/:id',
    name: 'HotelDetail',
    component: HotelDetail
  },
  {
    path: '/orders',
    name: 'Orders',
    component: Orders
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
