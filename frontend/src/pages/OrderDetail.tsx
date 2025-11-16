import { useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { ordersApi } from '../api/orders'
import { Loading } from '../components/Loading'
import { CheckCircle } from 'lucide-react'

export function OrderDetail() {
  const { id } = useParams<{ id: string }>()
  const { data: order, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.get(Number(id)),
    enabled: !!id,
  })

  if (isLoading) return <Loading />

  if (!order) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        <h1 className="text-2xl font-bold text-gray-800 mb-4">Заказ не найден</h1>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex items-center gap-3 mb-4">
            <CheckCircle className="text-tea-600" size={32} />
            <div>
              <h1 className="text-2xl font-bold text-tea-800">Заказ успешно создан!</h1>
              <p className="text-gray-600">Номер заказа: #{order.id}</p>
            </div>
          </div>

          <div className="border-t pt-4">
            <h2 className="text-lg font-semibold mb-4">Детали заказа</h2>
            <div className="space-y-3">
              {order.items.map((item, index) => (
                <div key={index} className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Товар #{item.good_id}</p>
                    <p className="text-sm text-gray-600">
                      Количество: {item.quantity} × {item.price} ₽
                    </p>
                  </div>
                  <p className="font-semibold">{item.quantity * item.price} ₽</p>
                </div>
              ))}
            </div>
          </div>

          <div className="border-t mt-4 pt-4">
            <div className="flex justify-between items-center">
              <span className="text-xl font-semibold">Итого:</span>
              <span className="text-2xl font-bold text-tea-600">{order.total} ₽</span>
            </div>
          </div>

          <div className="mt-6">
            <p className="text-sm text-gray-600">
              Статус: <span className="font-semibold text-tea-600">{order.status}</span>
            </p>
          </div>
        </div>

        <div className="text-center">
          <a
            href="/"
            className="inline-block bg-tea-600 text-white px-6 py-3 rounded-lg hover:bg-tea-700"
          >
            Вернуться к каталогу
          </a>
        </div>
      </div>
    </div>
  )
}

