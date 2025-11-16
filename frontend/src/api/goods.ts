import { apiClient } from './client'
import type { Good } from '../types'

export const goodsApi = {
  list: async (limit = 10, offset = 0): Promise<{ goods: Good[]; total: number }> => {
    const { data } = await apiClient.get('/goods', {
      params: { limit, offset },
    })
    return data
  },

  get: async (id: number): Promise<Good> => {
    const { data } = await apiClient.get(`/goods/${id}`)
    return data
  },
}

