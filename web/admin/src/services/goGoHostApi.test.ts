import { configureStore } from '@reduxjs/toolkit';
import { beforeEach, describe, it, vi } from 'vitest';
import { clearTokens } from '../auth/oidc';
import { goGoHostApi } from './goGoHostApi';

class MemoryStorage {
  private data = new Map<string, string>();
  getItem(key: string) { return this.data.get(key) ?? null; }
  setItem(key: string, value: string) { this.data.set(key, value); }
  removeItem(key: string) { this.data.delete(key); }
  clear() { this.data.clear(); }
}

function makeStore() {
  return configureStore({
    reducer: { [goGoHostApi.reducerPath]: goGoHostApi.reducer },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(goGoHostApi.middleware),
  });
}

function storeTokens(value: unknown) {
  localStorage.setItem('go-go-host.oidc.tokens', JSON.stringify(value));
}

beforeEach(() => {
  vi.restoreAllMocks();
  vi.stubGlobal('localStorage', new MemoryStorage());
  clearTokens();
});

describe('goGoHostApi OIDC refresh transport', () => {
  it('refreshes an expired token before an API request', async () => {
    storeTokens({ idToken: 'id-old', accessToken: 'access-old', refreshToken: 'refresh-old', expiresAt: Date.now() + 1_000 });
    const seenAuthHeaders: string[] = [];
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = input instanceof Request ? input.url : String(input);
      if (url.endsWith('/api/v1/config')) {
        return Response.json({
          baseDomain: 'hosting.example',
          publicBaseUrl: 'https://hosting.example',
          devAuth: false,
          oidc: { issuer: 'https://auth.example/realms/go-go-host', clientId: 'go-go-host-dashboard', scopes: ['openid', 'profile', 'email'] },
        });
      }
      if (url.endsWith('/.well-known/openid-configuration')) {
        return Response.json({ token_endpoint: 'https://auth.example/token', authorization_endpoint: 'https://auth.example/auth' });
      }
      if (url === 'https://auth.example/token') {
        return Response.json({ access_token: 'access-new', id_token: 'id-new', refresh_token: 'refresh-new', expires_in: 300 });
      }
      if (url.endsWith('/api/v1/me')) {
        const headers = input instanceof Request ? input.headers : new Headers(init?.headers as HeadersInit | undefined);
        seenAuthHeaders.push(headers.get('authorization') || headers.get('Authorization') || '');
        return Response.json({ user: { id: 'user_1', email: 'a@example.com', displayName: 'A' }, memberships: [], platformAdmin: false });
      }
      throw new Error(`unexpected fetch ${url}`);
    }));

    const store = makeStore();
    const action = await store.dispatch(goGoHostApi.endpoints.getMe.initiate());
    if ('error' in action) throw new Error(`query failed ${JSON.stringify(action.error)}`);
    const result = action.data;
    if (!result || result.user.email !== 'a@example.com') throw new Error(`unexpected result ${JSON.stringify(result)}`);
    if (seenAuthHeaders.length !== 1 || seenAuthHeaders[0] !== 'Bearer access-new') throw new Error(`unexpected auth headers ${JSON.stringify(seenAuthHeaders)}`);
  });

  it('retries once with a refreshed token after a 401 response', async () => {
    storeTokens({ idToken: 'id-old', accessToken: 'access-old', refreshToken: 'refresh-old', expiresAt: Date.now() + 5 * 60_000 });
    const seenAuthHeaders: string[] = [];
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = input instanceof Request ? input.url : String(input);
      if (url.endsWith('/api/v1/config')) {
        return Response.json({
          baseDomain: 'hosting.example',
          publicBaseUrl: 'https://hosting.example',
          devAuth: false,
          oidc: { issuer: 'https://auth.example/realms/go-go-host', clientId: 'go-go-host-dashboard', scopes: ['openid', 'profile', 'email'] },
        });
      }
      if (url.endsWith('/.well-known/openid-configuration')) {
        return Response.json({ token_endpoint: 'https://auth.example/token', authorization_endpoint: 'https://auth.example/auth' });
      }
      if (url === 'https://auth.example/token') {
        return Response.json({ access_token: 'access-new', id_token: 'id-new', refresh_token: 'refresh-new', expires_in: 300 });
      }
      if (url.endsWith('/api/v1/me')) {
        const headers = input instanceof Request ? input.headers : new Headers(init?.headers as HeadersInit | undefined);
        const header = headers.get('authorization') || headers.get('Authorization') || '';
        seenAuthHeaders.push(header);
        if (header === 'Bearer access-old') return Response.json({ error: 'expired' }, { status: 401 });
        return Response.json({ user: { id: 'user_1', email: 'a@example.com', displayName: 'A' }, memberships: [], platformAdmin: false });
      }
      throw new Error(`unexpected fetch ${url}`);
    }));

    const store = makeStore();
    const action = await store.dispatch(goGoHostApi.endpoints.getMe.initiate());
    if ('error' in action) throw new Error(`query failed ${JSON.stringify(action.error)}`);
    if (seenAuthHeaders.length !== 2 || seenAuthHeaders[0] !== 'Bearer access-old' || seenAuthHeaders[1] !== 'Bearer access-new') {
      throw new Error(`unexpected auth headers ${JSON.stringify(seenAuthHeaders)}`);
    }
  });
});
