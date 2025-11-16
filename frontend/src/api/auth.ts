import { apiClient } from './client'
import type { LoginResponse, User } from '../types'

export const authApi = {
  register: async (email: string, name: string, password: string) => {
    const { data } = await apiClient.post('/auth/register', {
      email,
      name,
      password,
    })
    return data
  },

  login: async (email: string, password: string): Promise<LoginResponse> => {
    const { data } = await apiClient.post('/auth/login', {
      email,
      password,
    })
    return data
  },

  getMe: async (): Promise<User> => {
    const { data } = await apiClient.get('/users/me')
    return data
  },
}

