import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

export function apiErrorMessage(error: unknown): string {
  if (!error) return '';
  if (isRecord(error) && 'status' in error) {
    const fetchError = error as FetchBaseQueryError;
    if ('data' in fetchError && isRecord(fetchError.data) && typeof fetchError.data.error === 'string') return fetchError.data.error;
    if ('error' in fetchError && typeof fetchError.error === 'string') return fetchError.error;
    return `Request failed with status ${String(fetchError.status)}`;
  }
  if (error instanceof Error) return error.message;
  return String(error);
}
