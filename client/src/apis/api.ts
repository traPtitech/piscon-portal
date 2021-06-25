import axios from 'axios'
import { randomString, pkce } from './utils'
export const traQBaseURL = 'https://q.trap.jp/api/v3'
export const traQClientID = 'J0RR7Auk9OVa4LZnQ4pD37hupkEkYloEHiIU'
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

  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`)
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: traQClientID,
    response_type: 'code',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state
  }).toString()
  window.location.assign(authorizationEndpointUrl.toString())
}

export function fetchAuthToken(code: string, verifier: string) {
  return axios.post(
    `${traQBaseURL}/oauth2/token`,
    new URLSearchParams({
      client_id: traQClientID,
      grant_type: 'authorization_code',
      code_verifier: verifier,
      code
    })
  )
}

export function revokeAuthToken(token: string) {
  return axios.post(
    `${traQBaseURL}/oauth2/revoke`,
    new URLSearchParams({ token })
  )
}

export function getMe() {
  return axios.get(`${traQBaseURL}/users/me`)
}

export function getMeGroup() {
  return axios.get(`${traQBaseURL}/users/me/groups`)
}

export function getRsults() {
  return axios.get(`/api/results`)
}

export function getNewer() {
  return axios.get(`/api/newer`)
}

export function getTeam(id) {
  return axios.get(`/api/team/${id}`)
}

export function getUser(id) {
  return axios.get(`/api/user/${id}`)
}

export function getQueue() {
  // TODO: Fix
  return axios.get(`/api/benchmark/queue`)
}
