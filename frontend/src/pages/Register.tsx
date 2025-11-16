import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuthStore } from '../store/authStore'

export function Register() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      await authApi.register(email, name, password)
      // После регистрации автоматически логинимся
      const { token, user } = await authApi.login(email, password)
      setAuth(user, token)
      navigate('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка регистрации')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-tea-50 to-tea-100">
      <div className="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 className="text-3xl font-bold text-tea-800 mb-6 text-center">
          Регистрация
        </h1>
        
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Имя
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Пароль
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={6}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-tea-500 focus:border-transparent"
            />
            <p className="text-xs text-gray-500 mt-1">Минимум 6 символов</p>
          </div>

          <button
            type="submit"
            className="w-full bg-tea-600 text-white py-2 rounded-lg hover:bg-tea-700 transition-colors font-semibold"
          >
            Зарегистрироваться
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Уже есть аккаунт?{' '}
          <Link to="/login" className="text-tea-600 hover:underline">
            Войти
          </Link>
        </p>
      </div>
    </div>
  )
}

