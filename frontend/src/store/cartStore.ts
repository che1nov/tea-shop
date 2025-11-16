import { create } from 'zustand'
import type { CartItem, Good } from '../types'

interface CartStore {
  items: CartItem[]
  addItem: (good: Good, quantity?: number) => void
  removeItem: (goodId: number) => void
  updateQuantity: (goodId: number, quantity: number) => void
  clear: () => void
  getTotal: () => number
  getItemCount: () => number
}

export const useCartStore = create<CartStore>((set, get) => ({
  items: [],

  addItem: (good, quantity = 1) => {
    const items = get().items
    const existingItem = items.find((item) => item.id === good.id)

    if (existingItem) {
      set({
        items: items.map((item) =>
          item.id === good.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        ),
      })
    } else {
      set({
        items: [...items, { ...good, quantity }],
      })
    }
  },

  removeItem: (goodId) => {
    set({
      items: get().items.filter((item) => item.id !== goodId),
    })
  },

  updateQuantity: (goodId, quantity) => {
    if (quantity <= 0) {
      get().removeItem(goodId)
      return
    }

    set({
      items: get().items.map((item) =>
        item.id === goodId ? { ...item, quantity } : item
      ),
    })
  },

  clear: () => {
    set({ items: [] })
  },

  getTotal: () => {
    return get().items.reduce((total, item) => total + item.price * item.quantity, 0)
  },

  getItemCount: () => {
    return get().items.reduce((count, item) => count + item.quantity, 0)
  },
}))

