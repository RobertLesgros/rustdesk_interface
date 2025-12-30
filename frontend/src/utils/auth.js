const TokenKey = 'access_token'
const OidcCode = 'oidc_code'
const OidcCodeExpiry = 'oidc_code_expiry';

export function getToken() {
  return localStorage.getItem(TokenKey)
}

export function setToken(token) {
  localStorage.setItem(`wc-option:local:access_token`, token)
  return localStorage.setItem(TokenKey, token)
}

export function removeToken() {
  return localStorage.removeItem(TokenKey)
}

// Set code and store current timestamp (in milliseconds)
export function setCode(code) {
  const now = Date.now(); // Current timestamp (milliseconds)
  const expiry = now + 60 * 1000; // Expires after 60 seconds

  localStorage.setItem(OidcCode, code); // Store code
  localStorage.setItem(OidcCodeExpiry, expiry); // Store expiry timestamp
}

// Get code, delete and return null if expired
export function getCode() {
  const expiry = localStorage.getItem(OidcCodeExpiry); // Get expiry timestamp
  const now = Date.now(); // Current timestamp

  if (expiry && now > parseInt(expiry)) {
    // If expired, delete code and expiry time
    removeCode();
    return null;
  }
  return localStorage.getItem(OidcCode); // Return code (if not expired)
}

// Delete code and expiry time
export function removeCode() {
  localStorage.removeItem(OidcCode);
  localStorage.removeItem(OidcCodeExpiry);
}
