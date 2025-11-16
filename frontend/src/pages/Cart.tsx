import { useCartStore } from '../store/cartStore'
import { useAuthStore } from '../store/authStore'
import { useNavigate } from 'react-router-dom'
import { ordersApi } from '../api/orders'
import { Trash2, Plus, Minus } from 'lucide-react'
import { useState } from 'react'

export function Cart() {
  const { items, removeItem, updateQuantity, getTotal, clear } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()
  const [isCreating, setIsCreating] = useState(false)
  const [showAddressModal, setShowAddressModal] = useState(false)
  const [address, setAddress] = useState('')

  const handleCheckout = async () => {
    if (!isAuthenticated()) {
      navigate('/login')
      return
    }

    // Показываем модальное окно для ввода адреса
    setShowAddressModal(true)
  }

  const handleConfirmOrder = async () => {
    if (!address.trim()) {
      alert('Пожалуйста, укажите адрес доставки')
      return
    }

    setIsCreating(true)
    try {
      const orderItems = items.map((item) => ({
        good_id: item.id,
        quantity: item.quantity,
        price: item.price,
      }))

      const order = await ordersApi.create(orderItems, address)
      clear()
      setShowAddressModal(false)
      setAddress('')
      navigate(`/orders/${order.id}`)
    } catch (error) {
      console.error('Ошибка создания заказа:', error)
      alert('Не удалось создать заказ')
    } finally {
      setIsCreating(false)
    }
  }

  if (items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        <h1 className="text-4xl font-bold text-tea-800 mb-4">Корзина пуста</h1>
        <p className="text-gray-600 mb-8">Добавьте товары в корзину</p>
        <button
          onClick={() => navigate('/')}
          className="bg-tea-600 text-white px-6 py-3 rounded-lg hover:bg-tea-700"
        >
          Перейти к каталогу
        </button>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-tea-800 mb-8">Корзина</h1>
      
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        {items.map((item) => (
          <div key={item.id} className="flex items-center justify-between py-4 border-b last:border-b-0">
            <div className="flex-1">
              <h3 className="font-semibold text-lg">{item.name}</h3>
              <p className="text-gray-600">{item.price} ₽ × {item.quantity}</p>
            </div>
            
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <button
                  onClick={() => updateQuantity(item.id, item.quantity - 1)}
                  className="p-1 rounded hover:bg-gray-100"
                >
                  <Minus size={20} />
                </button>
                <span className="w-8 text-center">{item.quantity}</span>
                <button
                  onClick={() => updateQuantity(item.id, item.quantity + 1)}
                  className="p-1 rounded hover:bg-gray-100"
                >
                  <Plus size={20} />
                </button>
              </div>
              
              <p className="font-bold w-24 text-right">
                {item.price * item.quantity} ₽
              </p>
              
              <button
                onClick={() => removeItem(item.id)}
                className="p-2 text-red-600 hover:bg-red-50 rounded"
              >
                <Trash2 size={20} />
              </button>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <span className="text-xl font-semibold">Итого:</span>
          <span className="text-3xl font-bold text-tea-600">{getTotal()} ₽</span>
        </div>
        <button
          onClick={handleCheckout}
          disabled={isCreating}
          className="w-full bg-tea-600 text-white py-3 rounded-lg hover:bg-tea-700 disabled:bg-gray-300 transition-colors text-lg font-semibold"
        >
          {isCreating ? 'Оформление...' : 'Оформить заказ'}
        </button>
      </div>

      {/* Модальное окно для ввода адреса */}
      {showAddressModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-2xl font-bold text-tea-800 mb-4">Адрес доставки</h2>
            <textarea
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Введите полный адрес доставки (город, улица, дом, квартира)"
              rows={4}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500 mb-4"
            />
            <div className="flex gap-2">
              <button
                onClick={handleConfirmOrder}
                disabled={isCreating || !address.trim()}
                className="flex-1 bg-tea-600 text-white py-2 rounded-lg hover:bg-tea-700 disabled:bg-gray-300"
              >
                {isCreating ? 'Оформление...' : 'Подтвердить заказ'}
              </button>
              <button
                onClick={() => {
                  setShowAddressModal(false)
                  setAddress('')
                }}
                disabled={isCreating}
                className="px-6 py-2 border rounded-lg hover:bg-gray-50 disabled:bg-gray-100"
              >
                Отмена
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

