/* eslint-disable @typescript-eslint/camelcase */
// import { OAuth2ResponseType, OAuth2Scope, OAuth2Token } from '@traptitech/traq'
import { OAuth2Token } from '@traptitech/traq'
import axios from 'axios'
import { randomString, pkce } from '../../utils'
import traqApis from './traq'

const traQBaseURL = 'https://q.trap.jp/api/v3'
// const BASE_PATH = 'https://q.trap.jp/api/v3'
// const REDIRECT_URL = 'https://piscon.trap.jp'

export const traQClientID = 'nmVeJT08KHXIdB8xlrCIwa6YJTkISrP5zWzm'

export function setAuthToken(token: OAuth2Token) {
  if (token) {
    axios.defaults.headers.common[
      'Authorization'
    ] = `Bearer ${token.access_token}`
  } else {
    delete axios.defaults.headers.common['Authorization']
  }
}

export async function redirectAuthorizationEndpoint() {
  const state = randomString(10)
  const codeVerifier = randomString(43)
  const codeChallenge = await pkce(codeVerifier)

  sessionStorage.setItem(`login-code-verifier-${state}`, codeVerifier)
  // traqApis.getOAuth2Authorize(
  //   traQClientID,
  //   OAuth2ResponseType.Code,
  //   REDIRECT_URL,
  //   OAuth2Scope.Read,
  //   state,
  //   codeChallenge,
  //   'S256'
  // )
  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`)
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: traQClientID,
    response_type: 'code',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state
  }).toString()
  window.location.assign(authorizationEndpointUrl.toString())
  return
}

export async function fetchAuthToken(code: string, verifier: string) {
  return axios.post(
    `${traQBaseURL}/oauth2/token`,
    new URLSearchParams({
      client_id: traQClientID,
      grant_type: 'authorization_code',
      code_verifier: verifier,
      code
    })
  )
  // return traqApis.postOAuth2Token(
  //   'authorization_code',
  //   code,
  //   REDIRECT_URL,
  //   traQClientID,
  //   verifier
  // )
}

export function revokeAuthToken(token: OAuth2Token) {
  return traqApis.revokeMyToken(token.access_token)
}
