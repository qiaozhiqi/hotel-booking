import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import HotelDetail from '../views/HotelDetail.vue'
import Orders from '../views/Orders.vue'
import Login from '../views/Login.vue'
import Favorites from '../views/Favorites.vue'
import Invoices from '../views/Invoices.vue'
import Profile from '../views/Profile.vue'

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
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: Favorites
  },
  {
    path: '/invoices',
    name: 'Invoices',
    component: Invoices
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
