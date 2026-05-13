import type { ConfigResponse, OIDCConfig } from '../services/types';

const tokenStorageKey = 'go-go-host.oidc.tokens';
const pkceStorageKey = 'go-go-host.oidc.pkce';
const refreshSkewMs = 60_000;

export interface StoredOIDCTokens {
  issuer?: string;
  clientId?: string;
  scopes?: string[];
  idToken: string;
  accessToken?: string;
  refreshToken?: string;
  expiresAt?: number;
}

interface DiscoveryDocument {
  authorization_endpoint: string;
  token_endpoint: string;
  end_session_endpoint?: string;
}

interface TokenEndpointResponse {
  id_token?: string;
  access_token?: string;
  refresh_token?: string;
  token_type?: string;
  expires_in?: number;
  scope?: string;
}

export interface BearerTokenOptions {
  forceRefresh?: boolean;
}

let refreshInFlight: Promise<StoredOIDCTokens | null> | null = null;

function base64Url(bytes: Uint8Array): string {
  let binary = '';
  bytes.forEach((byte) => { binary += String.fromCharCode(byte); });
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
}

async function sha256(input: string): Promise<string> {
  const bytes = new TextEncoder().encode(input);
  const digest = await crypto.subtle.digest('SHA-256', bytes);
  return base64Url(new Uint8Array(digest));
}

function randomString(bytes = 32): string {
  const data = new Uint8Array(bytes);
  crypto.getRandomValues(data);
  return base64Url(data);
}

function redirectURI(config: OIDCConfig): string {
  return `${window.location.origin}${config.redirectPath || '/app/auth/callback'}`;
}

function logoutRedirectURI(config: OIDCConfig): string {
  return `${window.location.origin}${config.logoutRedirectPath || '/app'}`;
}

async function discover(config: OIDCConfig): Promise<DiscoveryDocument> {
  const response = await fetch(`${config.issuer.replace(/\/$/, '')}/.well-known/openid-configuration`);
  if (!response.ok) throw new Error(`OIDC discovery failed with HTTP ${response.status}`);
  return response.json();
}

function tokenFrom(tokens: StoredOIDCTokens | null): string | undefined {
  return tokens?.accessToken || tokens?.idToken;
}

function tokenExpiresSoon(tokens: StoredOIDCTokens, now = Date.now()): boolean {
  if (!tokens.expiresAt) return false;
  return tokens.expiresAt - now <= refreshSkewMs;
}

function storeTokens(config: OIDCConfig, tokens: TokenEndpointResponse, previous?: StoredOIDCTokens | null): StoredOIDCTokens {
  const refreshToken = tokens.refresh_token || previous?.refreshToken;
  const idToken = tokens.id_token || previous?.idToken;
  const accessToken = tokens.access_token || previous?.accessToken;
  if (!idToken && !accessToken) throw new Error('OIDC token response did not include an id_token or access_token');
  const scopes = tokens.scope ? tokens.scope.split(/\s+/).filter(Boolean) : config.scopes;
  const stored: StoredOIDCTokens = {
    issuer: config.issuer,
    clientId: config.clientId,
    scopes,
    idToken: idToken || '',
    accessToken,
    refreshToken,
    expiresAt: tokens.expires_in ? Date.now() + tokens.expires_in * 1000 : previous?.expiresAt,
  };
  localStorage.setItem(tokenStorageKey, JSON.stringify(stored));
  return stored;
}

export function getStoredTokens(): StoredOIDCTokens | null {
  const raw = localStorage.getItem(tokenStorageKey);
  if (!raw) return null;
  try { return JSON.parse(raw) as StoredOIDCTokens; } catch { return null; }
}

export function bearerToken(): string | undefined {
  return tokenFrom(getStoredTokens());
}

export async function getValidBearerToken(config?: ConfigResponse, options: BearerTokenOptions = {}): Promise<string | undefined> {
  const tokens = getStoredTokens();
  if (!tokens) return undefined;
  if (!options.forceRefresh && !tokenExpiresSoon(tokens)) return tokenFrom(tokens);
  if (!isOIDCEnabled(config) || !tokens.refreshToken) {
    clearTokens();
    return undefined;
  }
  const refreshed = await refreshStoredTokens(config.oidc, tokens);
  return tokenFrom(refreshed);
}

export async function refreshStoredTokens(config: OIDCConfig, previous = getStoredTokens()): Promise<StoredOIDCTokens | null> {
  if (!previous?.refreshToken) {
    clearTokens();
    return null;
  }
  if (!refreshInFlight) {
    refreshInFlight = refreshTokens(config, previous)
      .catch((error) => {
        clearTokens();
        throw error;
      })
      .finally(() => { refreshInFlight = null; });
  }
  return refreshInFlight;
}

async function refreshTokens(config: OIDCConfig, previous: StoredOIDCTokens): Promise<StoredOIDCTokens> {
  const discovery = await discover(config);
  const body = new URLSearchParams({
    grant_type: 'refresh_token',
    client_id: config.clientId,
    refresh_token: previous.refreshToken || '',
  });
  const response = await fetch(discovery.token_endpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body,
  });
  if (!response.ok) throw new Error(`OIDC token refresh failed with HTTP ${response.status}`);
  const tokens = await response.json() as TokenEndpointResponse;
  return storeTokens(config, tokens, previous);
}

export function clearTokens() {
  localStorage.removeItem(tokenStorageKey);
}

export function isOIDCEnabled(config?: ConfigResponse): config is ConfigResponse & { oidc: OIDCConfig } {
  return Boolean(config && !config.devAuth && config.oidc?.issuer && config.oidc?.clientId);
}

export async function beginLogin(config: ConfigResponse & { oidc: OIDCConfig }, returnTo = window.location.pathname + window.location.search) {
  const discovery = await discover(config.oidc);
  const verifier = randomString(48);
  const state = randomString(24);
  const challenge = await sha256(verifier);
  sessionStorage.setItem(pkceStorageKey, JSON.stringify({ verifier, state, returnTo }));
  const params = new URLSearchParams({
    client_id: config.oidc.clientId,
    redirect_uri: redirectURI(config.oidc),
    response_type: 'code',
    scope: (config.oidc.scopes?.length ? config.oidc.scopes : ['openid', 'profile', 'email']).join(' '),
    state,
    code_challenge: challenge,
    code_challenge_method: 'S256',
  });
  window.location.assign(`${discovery.authorization_endpoint}?${params.toString()}`);
}

export async function completeLogin(config: ConfigResponse & { oidc: OIDCConfig }, callbackURL = window.location.href): Promise<string> {
  const url = new URL(callbackURL);
  const code = url.searchParams.get('code');
  const state = url.searchParams.get('state');
  const error = url.searchParams.get('error');
  if (error) throw new Error(url.searchParams.get('error_description') || error);
  if (!code || !state) throw new Error('OIDC callback is missing code or state');
  const raw = sessionStorage.getItem(pkceStorageKey);
  if (!raw) throw new Error('OIDC login state was not found; start login again');
  const saved = JSON.parse(raw) as { verifier: string; state: string; returnTo?: string };
  if (saved.state !== state) throw new Error('OIDC state mismatch; refusing login callback');
  const discovery = await discover(config.oidc);
  const body = new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: config.oidc.clientId,
    redirect_uri: redirectURI(config.oidc),
    code_verifier: saved.verifier,
    code,
  });
  const response = await fetch(discovery.token_endpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body,
  });
  if (!response.ok) throw new Error(`OIDC token exchange failed with HTTP ${response.status}`);
  const tokens = await response.json() as TokenEndpointResponse;
  if (!tokens.id_token) throw new Error('OIDC token response did not include an id_token');
  storeTokens(config.oidc, tokens);
  sessionStorage.removeItem(pkceStorageKey);
  return saved.returnTo || '/app';
}

export async function logout(config?: ConfigResponse) {
  const tokens = getStoredTokens();
  clearTokens();
  if (!isOIDCEnabled(config)) {
    window.location.assign('/app');
    return;
  }
  try {
    const discovery = await discover(config.oidc);
    if (discovery.end_session_endpoint) {
      const params = new URLSearchParams({ post_logout_redirect_uri: logoutRedirectURI(config.oidc) });
      if (tokens?.idToken) params.set('id_token_hint', tokens.idToken);
      window.location.assign(`${discovery.end_session_endpoint}?${params.toString()}`);
      return;
    }
  } catch {
    // Fall through to local redirect if discovery/logout is unavailable.
  }
  window.location.assign(logoutRedirectURI(config.oidc));
}
