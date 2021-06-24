import axios from 'axios'
import { randomString, pkce } from './utils'
export const traQBaseURL = 'https://q.trap.jp/api/v3'
export const traQClientID = 'J0RR7Auk9OVa4LZnQ4pD37hupkEkYloEHiIU'

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
    /* eslint-disable @typescript-eslint/camelcase */
    client_id: traQClientID,
    response_type: 'code',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state
  }).toString()
  window.location.assign(authorizationEndpointUrl.toString())
}
