import { create } from 'zustand'
import type { User } from '../types'

interface AuthStore {
  user: User | null
  token: string | null
  setAuth: (user: User, token: string) => void
  logout: () => void
  isAuthenticated: () => boolean
  isAdmin: () => boolean
  getUserRole: () => string | null
}

// Извлекает роль из JWT токена
function getRoleFromToken(token: string | null): string | null {
  if (!token) return null
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role || null
  } catch {
    return null
  }
}

// Инициализация пользователя из localStorage с извлечением роли из токена
function initUserFromStorage() {
  const storedUser = localStorage.getItem('user')
  const storedToken = localStorage.getItem('token')
  
  if (!storedUser || !storedToken) {
    return { user: null, token: null }
  }
  
  try {
    const user = JSON.parse(storedUser)
    const role = getRoleFromToken(storedToken)
    return {
      user: { ...user, role: role || user.role || 'user' },
      token: storedToken,
    }
  } catch {
    return { user: null, token: null }
  }
}

export const useAuthStore = create<AuthStore>((set, get) => {
  const { user, token } = initUserFromStorage()
  
  return {
    user,
    token,

    setAuth: (user, token) => {
      // Извлекаем роль из токена и добавляем в user
      const role = getRoleFromToken(token)
      const userWithRole = { ...user, role: role || 'user' }
      
      localStorage.setItem('user', JSON.stringify(userWithRole))
      localStorage.setItem('token', token)
      set({ user: userWithRole, token })
    },

    logout: () => {
      localStorage.removeItem('user')
      localStorage.removeItem('token')
      set({ user: null, token: null })
    },

    isAuthenticated: () => {
      return !!get().token && !!get().user
    },

    isAdmin: () => {
      const state = get()
      // Проверяем роль из токена или из user объекта
      const role = getRoleFromToken(state.token) || state.user?.role
      return role === 'admin'
    },

    getUserRole: () => {
      const state = get()
      return getRoleFromToken(state.token) || state.user?.role || null
    },
  }
})

