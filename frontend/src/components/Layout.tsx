import { Outlet, Link, useNavigate } from 'react-router-dom'
import { ShoppingCart, User, LogOut } from 'lucide-react'
import { useCartStore } from '../store/cartStore'
import { useAuthStore } from '../store/authStore'

export function Layout() {
  const itemCount = useCartStore((state) => state.getItemCount())
  const { user, logout, isAuthenticated, isAdmin } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-md">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="text-2xl font-bold text-tea-600">
              🍵 Tea Shop
            </Link>
            
            <div className="flex items-center gap-4">
              <Link
                to="/cart"
                className="relative flex items-center gap-2 text-gray-700 hover:text-tea-600"
              >
                <ShoppingCart size={24} />
                {itemCount > 0 && (
                  <span className="absolute -top-2 -right-2 bg-tea-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                    {itemCount}
                  </span>
                )}
              </Link>

              {isAuthenticated() ? (
                <div className="flex items-center gap-4">
                  {isAdmin() && (
                    <Link
                      to="/admin"
                      className="text-gray-700 hover:text-tea-600"
                    >
                      Админка
                    </Link>
                  )}
                  <div className="flex items-center gap-2 text-gray-700">
                    <User size={20} />
                    <span>{user?.name}</span>
                  </div>
                  <button
                    onClick={handleLogout}
                    className="flex items-center gap-2 text-gray-700 hover:text-red-600"
                  >
                    <LogOut size={20} />
                    Выйти
                  </button>
                </div>
              ) : (
                <Link
                  to="/login"
                  className="text-gray-700 hover:text-tea-600"
                >
                  Войти
                </Link>
              )}
            </div>
          </div>
        </div>
      </nav>

      <main>
        <Outlet />
      </main>
    </div>
  )
}

