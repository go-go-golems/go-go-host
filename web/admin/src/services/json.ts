import type { ValidationReport } from './types';

export function parseJson<T>(value: string | undefined, fallback: T): T {
  if (!value) return fallback;
  try { return JSON.parse(value) as T; } catch { return fallback; }
}

export function parseValidationReport(value: string | undefined): ValidationReport {
  return parseJson<ValidationReport>(value, { valid: false, files: 0, bytes: 0, errors: ['Unable to parse validation report JSON'] });
}

export function parseManifest(value: string | undefined): Record<string, unknown> {
  return parseJson<Record<string, unknown>>(value, { parseError: 'Unable to parse manifest JSON' });
}
