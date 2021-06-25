import axios from 'axios'
import { randomString, pkce } from '../utils'
import {
  Apis,
  Configuration,
  OAuth2ResponseType,
  OAuth2Scope
} from '@traptitech/traq'

export const BASE_PATH = '/api/v3'

export const traQClientID = 'J0RR7Auk9OVa4LZnQ4pD37hupkEkYloEHiIU'

const RedirectURL = 'http://localhost:8080/dashboard' //todo:分岐

export const api = new Apis(
  new Configuration({
    basePath: BASE_PATH
  })
)
/* eslint-disable @typescript-eslint/camelcase */
export function setAuthToken(token: string) {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
  } else {
    delete axios.defaults.headers.common['Authorization']
  }
}

export async function redirectAuthorizationEndpoint(): Promise<void> {
  const state = randomString(10)
  const codeVerifier = randomString(43)
  const codeChallenge = await pkce(codeVerifier)

  sessionStorage.setItem(`login-code-verifier-${state}`, codeVerifier)
  await api.getOAuth2Authorize(
    traQClientID,
    OAuth2ResponseType.Code,
    RedirectURL,
    OAuth2Scope.Write,
    state,
    codeChallenge
  )
  return
}

export function fetchAuthToken(code: string, verifier: string) {
  return axios.post(
    `$/oauth2/token`,
    new URLSearchParams({
      client_id: traQClientID,
      grant_type: 'authorization_code',
      code_verifier: verifier,
      code
    })
  )
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
