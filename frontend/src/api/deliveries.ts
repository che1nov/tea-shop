import { apiClient } from './client'
import type { Delivery } from '../types'

export const deliveriesApi = {
  list: async (limit: number = 100, offset: number = 0, status?: string): Promise<{ deliveries: Delivery[]; total: number }> => {
    const params = new URLSearchParams()
    params.append('limit', limit.toString())
    params.append('offset', offset.toString())
    if (status) {
      params.append('status', status)
    }
    const { data } = await apiClient.get(`/admin/deliveries?${params.toString()}`)
    return data
  },

  updateStatus: async (id: number, status: string): Promise<Delivery> => {
    const { data } = await apiClient.put(`/admin/deliveries/${id}/status`, { status })
    return data
  },
}

