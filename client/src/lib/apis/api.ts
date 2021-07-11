/* eslint-disable @typescript-eslint/camelcase */
import { OAuth2Token } from '@traptitech/traq'
import axios from 'axios'
import apis from '.'

const traQBaseURL = 'https://q.trap.jp/api/v3'
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
  const pkceParams = (await apis.authCodePost()).data
  const authorizationEndpointUrl = new URL(`${traQBaseURL}/oauth2/authorize`)
  authorizationEndpointUrl.search = new URLSearchParams({
    client_id: pkceParams.client_id,
    response_type: pkceParams.response_type,
    code_challenge: pkceParams.code_challenge,
    code_challenge_method: pkceParams.code_challenge_method
  }).toString()
  window.location.assign(authorizationEndpointUrl.toString())
  return
}
