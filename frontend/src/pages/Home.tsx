import { useQuery } from '@tanstack/react-query'
import { goodsApi } from '../api/goods'
import { GoodCard } from '../components/GoodCard'
import { Loading } from '../components/Loading'

export function Home() {
  const { data, isLoading } = useQuery({
    queryKey: ['goods'],
    queryFn: () => goodsApi.list(20, 0),
  })

  if (isLoading) return <Loading />

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-tea-800 mb-8 text-center">
        🍵 Магазин чая
      </h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {data?.goods.map((good) => (
          <GoodCard key={good.id} good={good} />
        ))}
      </div>
      {data?.goods.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-600 text-lg">Товары не найдены</p>
        </div>
      )}
    </div>
  )
}

