import { useState, useEffect, useCallback } from 'react'
import { generatePassword, fetchHistory } from '../api'

const cardStyle = {
  borderRadius: 20,
  border: 'none',
  boxShadow: '0 10px 40px rgba(0,0,0,0.2)',
}

const generateBtnStyle = {
  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  border: 'none',
  fontWeight: 600,
  transition: 'transform 0.2s',
}

const passwordInputStyle = {
  fontFamily: '"Courier New", monospace',
  fontSize: '1.1rem',
  letterSpacing: 1,
  background: '#f8f9fa',
  border: '2px solid #dee2e6',
}

export default function Dashboard({ onLogout }) {
  const [password, setPassword] = useState('')
  const [history, setHistory] = useState([])
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState('')

  const loadHistory = useCallback(async () => {
    const { ok, data } = await fetchHistory()
    if (ok) setHistory(data.history || [])
  }, [])

  useEffect(() => {
    loadHistory()
  }, [loadHistory])

  const handleGenerate = async () => {
    setLoading(true)
    setError('')
    const { ok, data } = await generatePassword()
    setLoading(false)
    if (ok) {
      setPassword(data.password)
      loadHistory()
    } else {
      setError(data.error || 'Failed to generate password')
    }
  }

  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      setError('Failed to copy to clipboard')
    }
  }

  return (
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-12 col-md-8 col-lg-6">
          <div className="card p-4 position-relative" style={cardStyle}>
            <span className="badge bg-success" style={{ position: 'absolute', top: 20, right: 20 }}>
              <i className="bi bi-shield-check"></i> Authenticated
            </span>

            <div className="text-center mb-4">
              <h1 className="display-5 fw-bold text-primary mb-2">
                <i className="bi bi-key-fill"></i> Random Pass
              </h1>
              <p className="text-muted">Random Password Generator</p>
            </div>

            {error && (
              <div className="alert alert-danger alert-dismissible fade show" role="alert">
                <i className="bi bi-exclamation-triangle-fill"></i> {error}
                <button type="button" className="btn-close" onClick={() => setError('')}></button>
              </div>
            )}

            <div className="mb-4">
              <label className="form-label fw-semibold">
                <i className="bi bi-lock-fill"></i> Generated Password
              </label>
              <div className="input-group">
                <input
                  type="text"
                  className="form-control"
                  style={passwordInputStyle}
                  value={password}
                  readOnly
                  placeholder="Click generate to create a password..."
                />
                {password && (
                  <button
                    className="btn btn-outline-secondary"
                    onClick={() => copyToClipboard(password)}
                    title="Copy to clipboard"
                    style={{ transition: 'all 0.2s' }}
                  >
                    <i className={`bi ${copied ? 'bi-check-lg text-success' : 'bi-clipboard'}`}></i>
                  </button>
                )}
              </div>
            </div>

            <button
              className="btn btn-primary btn-lg w-100 mb-4"
              onClick={handleGenerate}
              disabled={loading}
              style={generateBtnStyle}
            >
              {loading ? (
                <>
                  <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                  Generating...
                </>
              ) : (
                <>
                  <i className="bi bi-arrow-repeat"></i> Generate Password
                </>
              )}
            </button>

            <hr className="my-4" />

            <div>
              <div className="d-flex justify-content-between align-items-center mb-3">
                <h5 className="mb-0">
                  <i className="bi bi-clock-history"></i> Recent Passwords
                </h5>
                <span className="badge bg-secondary">{history.length}</span>
              </div>

              {history.length === 0 ? (
                <div className="text-center text-muted py-4">
                  <i className="bi bi-inbox" style={{ fontSize: '2rem' }}></i>
                  <p className="mt-2 mb-0">No passwords generated yet</p>
                </div>
              ) : (
                <ul className="list-group list-group-flush">
                  {history.map((pwd, i) => (
                    <li
                      key={i}
                      className="list-group-item d-flex justify-content-between align-items-center"
                      style={{ transition: 'all 0.2s', cursor: 'pointer' }}
                      onClick={() => copyToClipboard(pwd)}
                      onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#f8f9fa')}
                      onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '')}
                    >
                      <span className="font-monospace small text-truncate me-2">{pwd}</span>
                      <i className="bi bi-clipboard text-muted"></i>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            <div className="mt-4 text-center">
              <button className="btn btn-sm btn-outline-danger" onClick={onLogout}>
                <i className="bi bi-box-arrow-right"></i> Logout
              </button>
            </div>
          </div>

          <div className="text-center mt-4 text-white">
            <small>
              <i className="bi bi-shield-lock-fill"></i> Your passwords are securely stored
            </small>
          </div>
        </div>
      </div>
    </div>
  )
}
