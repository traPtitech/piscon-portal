const validChars =
  'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'

export function randomString(length: number) {
  let array = new Uint8Array(length)
  window.crypto.getRandomValues(array)
  array = array.map(x => validChars.charCodeAt(x % validChars.length))
  return String.fromCharCode(...array)
}
const b64Chars: { [key: string]: string } = { '+': '-', '/': '_', '=': '' }

function urlEncodeB64(input: string) {
  return input.replace(/[+/=]/g, m => b64Chars[m])
}

function bufferToBase64UrlEncoded(input: ArrayBuffer) {
  const bytes = new Uint8Array(input)
  return urlEncodeB64(window.btoa(String.fromCharCode(...bytes)))
}
function sha256(message: string) {
  const data = new TextEncoder().encode(message)
  return window.crypto.subtle.digest('SHA-256', data)
}

export async function pkce(verifier: string) {
  return sha256(verifier).then(bufferToBase64UrlEncoded)
}
