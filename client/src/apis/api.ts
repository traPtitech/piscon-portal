import axios from 'axios'
import { randomString, pkce } from '../utils'

export const BASE_PATH = 'https://q.trap.jp/api/1.0'

export const traQClientID = 'J0RR7Auk9OVa4LZnQ4pD37hupkEkYloEHiIU'

const CallbackURL = 'http://localhost:8080/auth/callback' //todo:分岐
const RedirectURL = 'http://localhost:8080'

/* eslint-disable @typescript-eslint/camelcase */
export function setAuthToken(token: string) {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
  } else {
    delete axios.defaults.headers.common['Authorization']
  }
}

export async function redirectAuthorizationEndpoint() {
  const state = randomString(10)
  const codeVerifier = randomString(43)
  const codeChallenge = await pkce(codeVerifier)

  sessionStorage.setItem(`login-code-verifier-${state}`, codeVerifier)
  const authorizationEndpointUrl = new URL(`${BASE_PATH}/oauth2/authorize`)
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: traQClientID,
    response_type: 'code',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state: state
  }).toString()
  window.location.assign(authorizationEndpointUrl.toString())
  return
}

export function fetchAuthToken(code: string, verifier: string) {
  return axios.post(
    `${BASE_PATH}/oauth2/token`,
    new URLSearchParams({
      client_id: traQClientID,
      grant_type: 'authorization_code',
      code_verifier: verifier,
      code: code
    })
  )
  return
}

export function revokeAuthToken(token: string) {
  return axios.post(`$/oauth2/revoke`, new URLSearchParams({ token }))
}

export function getMe() {
  return axios.get(`$/users/me`)
}

export function getMeGroup() {
  return axios.get(`$/users/me/groups`)
}

export function getRsults() {
  return axios.get(`/api/results`)
}

export function getNewer() {
  return axios.get(`/api/newer`)
}

export function getTeam(id: string) {
  return axios.get(`/api/team/${id}`)
}

export function getUser(id: string) {
  return axios.get(`/api/user/${id}`)
}

export function getQueue() {
  // TODO: Fix
  return axios.get(`/api/benchmark/queue`)
}
