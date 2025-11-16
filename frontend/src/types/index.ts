export interface Good {
  id: number
  sku: string
  name: string
  description: string
  price: number
  stock: number
  created_at: number
}

export interface User {
  id: number
  email: string
  name: string
  role?: string // "user" или "admin"
}

export interface LoginResponse {
  token: string
  user: User
}

export interface OrderItem {
  good_id: number
  quantity: number
  price: number
}

export interface Order {
  id: number
  user_id: number
  items: OrderItem[]
  total: number
  status: string
  address: string
  created_at: number
}

export interface CartItem extends Good {
  quantity: number
}

export interface Delivery {
  id: number
  order_id: number
  address: string
  status: string // pending, in_transit, delivered, cancelled
  created_at: number
  updated_at: number
}

