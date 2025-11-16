import { apiClient } from './client'
import type { Good } from '../types'

export interface CreateGoodRequest {
  name: string
  description: string
  price: number
  stock: number
}

export interface UpdateGoodRequest {
  name?: string
  description?: string
  price?: number
  stock?: number
}

export const adminApi = {
  createGood: async (data: CreateGoodRequest): Promise<Good> => {
    const { data: good } = await apiClient.post('/admin/goods', data)
    return good
  },

  updateGood: async (id: number, data: UpdateGoodRequest): Promise<Good> => {
    const { data: good } = await apiClient.put(`/admin/goods/${id}`, data)
    return good
  },

  deleteGood: async (id: number): Promise<void> => {
    const response = await apiClient.delete(`/admin/goods/${id}`)
    if (!response.data?.success) {
      throw new Error(response.data?.message || 'Ошибка удаления товара')
    }
    return response.data
  },
}
