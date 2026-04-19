<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <div class="logo">
          <span class="logo-icon">🏨</span>
          <span class="logo-text">酒店预订</span>
        </div>
        <p class="login-subtitle">欢迎回来，请登录您的账号</p>
      </div>

      <div class="tab-switch">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'login' }"
          @click="activeTab = 'login'"
        >
          登录
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'register' }"
          @click="activeTab = 'register'"
        >
          注册
        </button>
      </div>

      <form v-if="activeTab === 'login'" class="login-form" @submit.prevent="handleLogin">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input 
            type="text" 
            v-model="loginForm.username" 
            class="form-input"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">密码</label>
          <input 
            type="password" 
            v-model="loginForm.password" 
            class="form-input"
            placeholder="请输入密码"
            required
          />
        </div>
        <button 
          type="submit" 
          class="btn-submit"
          :disabled="loggingIn"
        >
          <span v-if="loggingIn">登录中...</span>
          <span v-else>登录</span>
        </button>
        <p class="form-tip">测试账号：admin / admin123 或 testuser / test123</p>
      </form>

      <form v-else class="login-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input 
            type="text" 
            v-model="registerForm.username" 
            class="form-input"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">密码</label>
          <input 
            type="password" 
            v-model="registerForm.password" 
            class="form-input"
            placeholder="请输入密码"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">确认密码</label>
          <input 
            type="password" 
            v-model="registerForm.confirmPassword" 
            class="form-input"
            placeholder="请再次输入密码"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">邮箱（可选）</label>
          <input 
            type="email" 
            v-model="registerForm.email" 
            class="form-input"
            placeholder="请输入邮箱"
          />
        </div>
        <div class="form-group">
          <label class="form-label">手机号（可选）</label>
          <input 
            type="tel" 
            v-model="registerForm.phone" 
            class="form-input"
            placeholder="请输入手机号"
          />
        </div>
        <button 
          type="submit" 
          class="btn-submit"
          :disabled="registering"
        >
          <span v-if="registering">注册中...</span>
          <span v-else>注册</span>
        </button>
      </form>

      <div class="login-footer">
        <router-link to="/" class="link-home">← 返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { userApi } from '../api'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()

    const activeTab = ref('login')
    const loggingIn = ref(false)
    const registering = ref(false)

    const loginForm = ref({
      username: '',
      password: ''
    })

    const registerForm = ref({
      username: '',
      password: '',
      confirmPassword: '',
      email: '',
      phone: ''
    })

    const handleLogin = async () => {
      if (!loginForm.value.username || !loginForm.value.password) {
        alert('请输入用户名和密码')
        return
      }

      loggingIn.value = true
      try {
        const res = await userApi.login({
          username: loginForm.value.username,
          password: loginForm.value.password
        })

        if (res.code === 200) {
          localStorage.setItem('user', JSON.stringify(res.data))
          window.dispatchEvent(new CustomEvent('login-success'))
          alert('登录成功！')
          router.push('/')
        } else {
          alert(res.message || '登录失败')
        }
      } catch (error) {
        console.error('登录失败:', error)
        alert('登录失败，请稍后重试')
      } finally {
        loggingIn.value = false
      }
    }

    const handleRegister = async () => {
      if (!registerForm.value.username || !registerForm.value.password) {
        alert('请输入用户名和密码')
        return
      }

      if (registerForm.value.password !== registerForm.value.confirmPassword) {
        alert('两次输入的密码不一致')
        return
      }

      registering.value = true
      try {
        const res = await userApi.register({
          username: registerForm.value.username,
          password: registerForm.value.password,
          email: registerForm.value.email,
          phone: registerForm.value.phone
        })

        if (res.code === 200) {
          alert('注册成功！请登录')
          activeTab.value = 'login'
          loginForm.value.username = registerForm.value.username
          loginForm.value.password = ''
          registerForm.value = {
            username: '',
            password: '',
            confirmPassword: '',
            email: '',
            phone: ''
          }
        } else {
          alert(res.message || '注册失败')
        }
      } catch (error) {
        console.error('注册失败:', error)
        alert('注册失败，请稍后重试')
      } finally {
        registering.value = false
      }
    }

    return {
      activeTab,
      loggingIn,
      registering,
      loginForm,
      registerForm,
      handleLogin,
      handleRegister
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: calc(100vh - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
}

.login-container {
  width: 100%;
  max-width: 420px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.logo-icon {
  font-size: 40px;
  margin-right: 8px;
}

.logo-text {
  font-size: 24px;
  font-weight: 600;
  color: #1a73e8;
}

.login-subtitle {
  font-size: 14px;
  color: #999;
}

.tab-switch {
  display: flex;
  margin-bottom: 24px;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 4px;
}

.tab-btn {
  flex: 1;
  padding: 10px 20px;
  background: transparent;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #1a73e8;
}

.tab-btn.active {
  background: #fff;
  color: #1a73e8;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.form-input {
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: #1a73e8;
}

.form-input::placeholder {
  color: #aaa;
}

.btn-submit {
  padding: 14px 24px;
  background: #1a73e8;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  margin-top: 8px;
}

.btn-submit:hover:not(:disabled) {
  background: #1557b0;
}

.btn-submit:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.form-tip {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-top: -8px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
}

.link-home {
  font-size: 14px;
  color: #1a73e8;
  text-decoration: none;
}

.link-home:hover {
  text-decoration: underline;
}
</style>
