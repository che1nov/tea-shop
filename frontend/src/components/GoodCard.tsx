import { ShoppingCart } from 'lucide-react'
import { useCartStore } from '../store/cartStore'
import type { Good } from '../types'

interface GoodCardProps {
  good: Good
}

export function GoodCard({ good }: GoodCardProps) {
  const addItem = useCartStore((state) => state.addItem)

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow">
      <div className="h-48 bg-gradient-to-br from-tea-100 to-tea-200 flex items-center justify-center">
        <span className="text-6xl">🍵</span>
      </div>
      <div className="p-4">
        <h3 className="text-xl font-semibold text-gray-800 mb-2">{good.name}</h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{good.description}</p>
        <div className="flex items-center justify-between">
          <div>
            <p className="text-2xl font-bold text-tea-600">{good.price} ₽</p>
            <p className="text-xs text-gray-500">В наличии: {good.stock}</p>
          </div>
          <button
            onClick={() => addItem(good)}
            disabled={good.stock === 0}
            className="bg-tea-600 text-white px-4 py-2 rounded-lg hover:bg-tea-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
          >
            <ShoppingCart size={20} />
            В корзину
          </button>
        </div>
      </div>
    </div>
  )
}

