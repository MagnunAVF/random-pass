const getToken = () => localStorage.getItem('token')

const authHeaders = () => ({
  'Content-Type': 'application/json',
  Authorization: `Bearer ${getToken()}`,
})

export async function login(email, password) {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  return { ok: res.ok, data: await res.json() }
}

export async function signup(username, email, password) {
  const res = await fetch('/api/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password }),
  })
  return { ok: res.ok, data: await res.json() }
}

export async function generatePassword() {
  const res = await fetch('/api/generate', {
    method: 'POST',
    headers: authHeaders(),
  })
  return { ok: res.ok, data: await res.json() }
}

export async function fetchHistory() {
  const res = await fetch('/api/history', {
    headers: authHeaders(),
  })
  return { ok: res.ok, data: await res.json() }
}
