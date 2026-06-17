import { useState } from 'react'
import { login, signup } from '../api'

const cardStyle = {
  borderRadius: 20,
  border: 'none',
  boxShadow: '0 10px 40px rgba(0,0,0,0.2)',
}

const btnStyle = {
  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  border: 'none',
  fontWeight: 600,
}

export default function AuthPage({ onAuth }) {
  const [tab, setTab] = useState('login')
  const [form, setForm] = useState({ username: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const set = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const switchTab = (t) => {
    setTab(t)
    setError('')
    setForm({ username: '', email: '', password: '' })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    const { ok, data } =
      tab === 'login'
        ? await login(form.email, form.password)
        : await signup(form.username, form.email, form.password)

    setLoading(false)

    if (!ok) {
      setError(data.error || 'Something went wrong')
      return
    }

    onAuth(data.token)
  }

  return (
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-12 col-md-6 col-lg-5">
          <div className="card p-4" style={cardStyle}>
            <div className="text-center mb-4">
              <h1 className="display-5 fw-bold text-primary mb-2">
                <i className="bi bi-key-fill"></i> Random Pass
              </h1>
              <p className="text-muted">Random Password Generator</p>
            </div>

            <ul className="nav nav-pills mb-4 justify-content-center">
              <li className="nav-item">
                <button
                  className={`nav-link ${tab === 'login' ? 'active' : ''}`}
                  onClick={() => switchTab('login')}
                >
                  Login
                </button>
              </li>
              <li className="nav-item">
                <button
                  className={`nav-link ${tab === 'signup' ? 'active' : ''}`}
                  onClick={() => switchTab('signup')}
                >
                  Sign Up
                </button>
              </li>
            </ul>

            {error && (
              <div className="alert alert-danger alert-dismissible fade show" role="alert">
                <i className="bi bi-exclamation-triangle-fill"></i> {error}
                <button type="button" className="btn-close" onClick={() => setError('')}></button>
              </div>
            )}

            <form onSubmit={handleSubmit}>
              {tab === 'signup' && (
                <div className="mb-3">
                  <label className="form-label fw-semibold">Username</label>
                  <input
                    type="text"
                    className="form-control"
                    value={form.username}
                    onChange={set('username')}
                    required
                  />
                </div>
              )}
              <div className="mb-3">
                <label className="form-label fw-semibold">Email</label>
                <input
                  type="email"
                  className="form-control"
                  value={form.email}
                  onChange={set('email')}
                  required
                />
              </div>
              <div className="mb-4">
                <label className="form-label fw-semibold">Password</label>
                <input
                  type="password"
                  className="form-control"
                  value={form.password}
                  onChange={set('password')}
                  required
                  minLength={8}
                />
              </div>
              <button type="submit" className="btn btn-primary btn-lg w-100" style={btnStyle} disabled={loading}>
                {loading ? (
                  <>
                    <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                    {tab === 'login' ? 'Signing in...' : 'Creating account...'}
                  </>
                ) : tab === 'login' ? (
                  'Sign In'
                ) : (
                  'Create Account'
                )}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
