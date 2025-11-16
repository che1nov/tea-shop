import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useAuthStore } from './store/authStore'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { Cart } from './pages/Cart'
import { Login } from './pages/Login'
import { Register } from './pages/Register'
import { OrderDetail } from './pages/OrderDetail'
import { Admin } from './pages/Admin'

const queryClient = new QueryClient()

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticatedFn = useAuthStore((state) => state.isAuthenticated)
  const isAuthenticated = isAuthenticatedFn()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticatedFn = useAuthStore((state) => state.isAuthenticated)
  const isAdminFn = useAuthStore((state) => state.isAdmin)
  
  const isAuthenticated = isAuthenticatedFn()
  const isAdmin = isAdminFn()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  
  if (!isAdmin) {
    return <Navigate to="/" replace />
  }
  
  return <>{children}</>
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Home />} />
            <Route path="cart" element={<Cart />} />
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register />} />
            <Route
              path="admin"
              element={
                <AdminRoute>
                  <Admin />
                </AdminRoute>
              }
            />
            <Route
              path="orders/:id"
              element={
                <PrivateRoute>
                  <OrderDetail />
                </PrivateRoute>
              }
            />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App

