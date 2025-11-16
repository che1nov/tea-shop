import { apiClient } from './client'
import type { Order, OrderItem } from '../types'

export const ordersApi = {
  create: async (items: OrderItem[], address: string): Promise<Order> => {
    const { data } = await apiClient.post('/orders', { items, address })
    return data
  },

  get: async (id: number): Promise<Order> => {
    const { data } = await apiClient.get(`/orders/${id}`)
    return data
  },
}

