import axios from 'axios'
import { randomString, pkce } from './utils'

// axios.defaults.withCredentials = true
export const traQBaseURL = 'https://q.trap.jp/api/1.0'
export const traQClientID = process.env.VUE_APP_API_CLIENT_ID || 'CySPaqKWiUXvwechb66dk0yubIlDqCcK07DV'
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : process.env.VUE_APP_API_ENDPOINT

export function setAuthToken (token) {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
  } else {
    delete axios.defaults.headers.common['Authorization']
  }
}

export async function redirectAuthorizationEndpoint () {
  const state = randomString(10)
  const codeVerifier = randomString(43)
  const codeChallenge = await pkce(codeVerifier)

  sessionStorage.setItem(`login-code-verifier-${state}`, codeVerifier)

  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`)
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: traQClientID,
    response_type: 'code',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state
  })
  window.location.assign(authorizationEndpointUrl)
}

export function fetchAuthToken (code, verifier) {
  return axios.post(`${traQBaseURL}/oauth2/token`, new URLSearchParams({
    client_id: traQClientID,
    grant_type: 'authorization_code',
    code_verifier: verifier,
    code
  }))
}

export function revokeAuthToken (token) {
  return axios.post(`${traQBaseURL}/oauth2/revoke`, new URLSearchParams({ token }))
}

export function getMe () {
  return axios.get(`${traQBaseURL}/users/me`)
}

export function getMeGroup () {
  return axios.get(`${traQBaseURL}/users/me/groups`)
}

export function getUser (id) {
  return axios.get(`${traQBaseURL}/users/${id}`)
}

export function getRsults () {
  return axios.get(`/api/results`)
}

export function getNewer () {
  return axios.get(`/api/newer`)
}

export function getTeam (name) {
  return axios.get(`/api/team/${name}`)
}

export function getQueue () {
  // TODO: Fix
  return axios.get(`/api/benchmark/queue`)
}
