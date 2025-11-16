import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { goodsApi } from '../api/goods'
import { adminApi } from '../api/admin'
import { deliveriesApi } from '../api/deliveries'
import { useState } from 'react'
import { Plus, Edit, Trash2, X, Package, Truck } from 'lucide-react'
import { Loading } from '../components/Loading'
import type { Good, Delivery } from '../types'

type Tab = 'goods' | 'deliveries'

export function Admin() {
  const [activeTab, setActiveTab] = useState<Tab>('goods')
  const [showModal, setShowModal] = useState(false)
  const [editingGood, setEditingGood] = useState<Good | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    stock: '',
  })
  const [statusFilter, setStatusFilter] = useState<string>('')

  const queryClient = useQueryClient()

  const { data: goodsData, isLoading: goodsLoading, error: goodsError } = useQuery({
    queryKey: ['admin-goods'],
    queryFn: () => goodsApi.list(100, 0),
    enabled: activeTab === 'goods',
  })

  const { data: deliveriesData, isLoading: deliveriesLoading, error: deliveriesError } = useQuery({
    queryKey: ['admin-deliveries', statusFilter],
    queryFn: () => deliveriesApi.list(100, 0, statusFilter || undefined),
    enabled: activeTab === 'deliveries',
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) => deliveriesApi.updateStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-deliveries'] })
    },
    onError: (error: any) => {
      console.error('Ошибка при обновлении статуса:', error)
      alert(`Ошибка обновления статуса: ${error.response?.data?.error || error.message || 'Неизвестная ошибка'}`)
    },
  })

  const createMutation = useMutation({
    mutationFn: adminApi.createGood,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-goods'] })
      queryClient.invalidateQueries({ queryKey: ['goods'] })
      setShowModal(false)
      resetForm()
    },
    onError: (error: any) => {
      console.error('Ошибка при создании товара:', error)
      alert(`Ошибка создания товара: ${error.response?.data?.error || error.message || 'Неизвестная ошибка'}`)
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: any }) => adminApi.updateGood(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-goods'] })
      queryClient.invalidateQueries({ queryKey: ['goods'] })
      setShowModal(false)
      setEditingGood(null)
      resetForm()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: adminApi.deleteGood,
    onSuccess: () => {
      // Инвалидируем кэш для обновления списка
      queryClient.invalidateQueries({ queryKey: ['admin-goods'] })
      queryClient.invalidateQueries({ queryKey: ['goods'] })
    },
    onError: (error: any) => {
      console.error('Ошибка при удалении:', error)
      alert(`Ошибка удаления: ${error.response?.data?.error || error.message || 'Неизвестная ошибка'}`)
    },
  })

  const resetForm = () => {
    setFormData({ name: '', description: '', price: '', stock: '' })
    setEditingGood(null)
  }

  const openCreateModal = () => {
    resetForm()
    setShowModal(true)
  }

  const openEditModal = (good: Good) => {
    setEditingGood(good)
    setFormData({
      name: good.name,
      description: good.description,
      price: good.price.toString(),
      stock: good.stock.toString(),
    })
    setShowModal(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const data = {
      name: formData.name,
      description: formData.description,
      price: parseFloat(formData.price),
      stock: parseInt(formData.stock),
    }

    if (editingGood) {
      updateMutation.mutate({ id: editingGood.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const handleDelete = (id: number, name: string) => {
    if (confirm(`Вы уверены, что хотите удалить товар "${name}" (ID: ${id})?`)) {
      deleteMutation.mutate(id)
    }
  }

  const isLoading = activeTab === 'goods' ? goodsLoading : deliveriesLoading
  const error = activeTab === 'goods' ? goodsError : deliveriesError

  if (isLoading) {
    return <Loading />
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">Ошибка загрузки: {error instanceof Error ? error.message : 'Неизвестная ошибка'}</p>
        </div>
      </div>
    )
  }

  const goods = goodsData?.goods || []
  const deliveries = deliveriesData?.deliveries || []

  const getStatusBadgeColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'in_transit':
        return 'bg-blue-100 text-blue-800'
      case 'delivered':
        return 'bg-green-100 text-green-800'
      case 'cancelled':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'pending':
        return 'Ожидает отправки'
      case 'in_transit':
        return 'В пути'
      case 'delivered':
        return 'Доставлено'
      case 'cancelled':
        return 'Отменено'
      default:
        return status
    }
  }

  const handleUpdateStatus = (id: number, currentStatus: string) => {
    let newStatus: string
    if (currentStatus === 'pending') {
      newStatus = 'in_transit'
    } else if (currentStatus === 'in_transit') {
      newStatus = 'delivered'
    } else {
      return
    }

    if (confirm(`Изменить статус доставки на "${getStatusLabel(newStatus)}"?`)) {
      updateStatusMutation.mutate({ id, status: newStatus })
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-4xl font-bold text-tea-800">Админ-панель</h1>
        {activeTab === 'goods' && (
          <button
            onClick={openCreateModal}
            className="bg-tea-600 text-white px-6 py-3 rounded-lg hover:bg-tea-700 flex items-center gap-2"
          >
            <Plus size={20} />
            Добавить товар
          </button>
        )}
      </div>

      {/* Вкладки */}
      <div className="mb-6 border-b border-gray-200">
        <nav className="flex space-x-8">
          <button
            onClick={() => setActiveTab('goods')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'goods'
                ? 'border-tea-500 text-tea-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <Package size={20} />
              Товары
            </div>
          </button>
          <button
            onClick={() => setActiveTab('deliveries')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'deliveries'
                ? 'border-tea-500 text-tea-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <div className="flex items-center gap-2">
              <Truck size={20} />
              Доставки
            </div>
          </button>
        </nav>
      </div>

      {/* Контент для товаров */}
      {activeTab === 'goods' && (
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Артикул
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Название
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Описание
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Цена
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Остаток
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Действия
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {goods.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                  Товары не найдены. Добавьте первый товар!
                </td>
              </tr>
            ) : (
              goods.map((good) => (
              <tr key={good.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-mono">
                  {good.sku || `GOOD-${String(good.id).padStart(6, '0')}`}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {good.name}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {good.description.substring(0, 50)}
                  {good.description.length > 50 && '...'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {good.price} ₽
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {good.stock}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div className="flex gap-2">
                    <button
                      onClick={() => openEditModal(good)}
                      className="text-tea-600 hover:text-tea-900"
                    >
                      <Edit size={18} />
                    </button>
                    <button
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        handleDelete(good.id, good.name)
                      }}
                      className="text-red-600 hover:text-red-900"
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>
                </td>
              </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      )}

      {/* Контент для доставок */}
      {activeTab === 'deliveries' && (
        <div className="space-y-4">
          {/* Фильтр по статусу */}
          <div className="bg-white rounded-lg shadow-md p-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Фильтр по статусу:
            </label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500"
            >
              <option value="">Все статусы</option>
              <option value="pending">Ожидает отправки</option>
              <option value="in_transit">В пути</option>
              <option value="delivered">Доставлено</option>
              <option value="cancelled">Отменено</option>
            </select>
          </div>

          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Заказ ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Адрес
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Статус
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {deliveries.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                      Доставки не найдены
                    </td>
                  </tr>
                ) : (
                  deliveries.map((delivery) => (
                    <tr key={delivery.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {delivery.id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {delivery.order_id}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {delivery.address.substring(0, 50)}
                        {delivery.address.length > 50 && '...'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusBadgeColor(delivery.status)}`}>
                          {getStatusLabel(delivery.status)}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        {delivery.status === 'pending' && (
                          <button
                            onClick={() => handleUpdateStatus(delivery.id, delivery.status)}
                            disabled={updateStatusMutation.isPending}
                            className="bg-tea-600 text-white px-4 py-2 rounded-lg hover:bg-tea-700 disabled:bg-gray-300 text-sm"
                          >
                            Отправить
                          </button>
                        )}
                        {delivery.status === 'in_transit' && (
                          <button
                            onClick={() => handleUpdateStatus(delivery.id, delivery.status)}
                            disabled={updateStatusMutation.isPending}
                            className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 disabled:bg-gray-300 text-sm"
                          >
                            Отметить доставленным
                          </button>
                        )}
                        {(delivery.status === 'delivered' || delivery.status === 'cancelled') && (
                          <span className="text-gray-400 text-sm">-</span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Модальное окно для товаров */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-bold text-tea-800">
                {editingGood ? 'Редактировать товар' : 'Добавить товар'}
              </h2>
              <button
                onClick={() => {
                  setShowModal(false)
                  resetForm()
                }}
                className="text-gray-500 hover:text-gray-700"
              >
                <X size={24} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Название
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Описание
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  required
                  rows={3}
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Цена (₽)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                    required
                    className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Остаток
                  </label>
                  <input
                    type="number"
                    min="0"
                    value={formData.stock}
                    onChange={(e) => setFormData({ ...formData, stock: e.target.value })}
                    required
                    className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500"
                  />
                </div>
              </div>

              <div className="flex gap-2 pt-4">
                <button
                  type="submit"
                  disabled={createMutation.isPending || updateMutation.isPending}
                  className="flex-1 bg-tea-600 text-white py-2 rounded-lg hover:bg-tea-700 disabled:bg-gray-300"
                >
                  {createMutation.isPending || updateMutation.isPending
                    ? 'Сохранение...'
                    : editingGood
                    ? 'Сохранить'
                    : 'Добавить'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false)
                    resetForm()
                  }}
                  className="px-6 py-2 border rounded-lg hover:bg-gray-50"
                >
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
