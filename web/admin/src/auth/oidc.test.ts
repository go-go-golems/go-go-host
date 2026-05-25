import { beforeEach, describe, expect, it, vi } from 'vitest';
import { clearTokens, getStoredTokens, getValidBearerToken } from './oidc';
import type { ConfigResponse } from '../services/types';

class MemoryStorage {
  private data = new Map<string, string>();
  getItem(key: string) { return this.data.get(key) ?? null; }
  setItem(key: string, value: string) { this.data.set(key, value); }
  removeItem(key: string) { this.data.delete(key); }
  clear() { this.data.clear(); }
}

const config: ConfigResponse = {
  baseDomain: 'hosting.example',
  publicBaseUrl: 'https://hosting.example',
  devAuth: false,
  oidc: {
    issuer: 'https://auth.example/realms/go-go-host',
    clientId: 'go-go-host-dashboard',
    scopes: ['openid', 'profile', 'email'],
  },
};

function storeTokens(value: unknown) {
  localStorage.setItem('go-go-host.oidc.tokens', JSON.stringify(value));
}

beforeEach(() => {
  vi.restoreAllMocks();
  vi.stubGlobal('localStorage', new MemoryStorage());
  clearTokens();
});

describe('OIDC token refresh', () => {
  it('returns the current access token while it is not close to expiry', async () => {
    storeTokens({ idToken: 'id-old', accessToken: 'access-old', refreshToken: 'refresh-old', expiresAt: Date.now() + 5 * 60_000 });
    const fetchSpy = vi.fn();
    vi.stubGlobal('fetch', fetchSpy);

    await expect(getValidBearerToken(config)).resolves.toBe('access-old');
    expect(fetchSpy).not.toHaveBeenCalled();
  });

  it('refreshes an expiring access token and preserves refresh token when Keycloak does not rotate it', async () => {
    storeTokens({ idToken: 'id-old', accessToken: 'access-old', refreshToken: 'refresh-old', expiresAt: Date.now() + 1_000 });
    const fetchSpy = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url.endsWith('/.well-known/openid-configuration')) {
        return Response.json({ token_endpoint: 'https://auth.example/token', authorization_endpoint: 'https://auth.example/auth' });
      }
      if (url === 'https://auth.example/token') {
        expect(init?.method).toBe('POST');
        expect(String(init?.body)).toContain('grant_type=refresh_token');
        expect(String(init?.body)).toContain('client_id=go-go-host-dashboard');
        expect(String(init?.body)).toContain('refresh_token=refresh-old');
        return Response.json({ access_token: 'access-new', id_token: 'id-new', expires_in: 300 });
      }
      throw new Error(`unexpected fetch ${url}`);
    });
    vi.stubGlobal('fetch', fetchSpy);

    await expect(getValidBearerToken(config)).resolves.toBe('access-new');
    expect(getStoredTokens()).toMatchObject({ accessToken: 'access-new', idToken: 'id-new', refreshToken: 'refresh-old', clientId: 'go-go-host-dashboard' });
  });

  it('clears tokens when refresh is required but no refresh token is available', async () => {
    storeTokens({ idToken: 'id-old', accessToken: 'access-old', expiresAt: Date.now() + 1_000 });

    await expect(getValidBearerToken(config)).resolves.toBeUndefined();
    expect(getStoredTokens()).toBeNull();
  });
});
