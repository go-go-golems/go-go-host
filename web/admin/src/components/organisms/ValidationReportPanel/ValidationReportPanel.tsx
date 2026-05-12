import { ValidationSummary } from '../../molecules/ValidationSummary';
import { ManifestSummary } from '../../molecules/ManifestSummary';
import type { Deployment, ValidationReport } from '../../../services/types';
import './ValidationReportPanel.css';
function parseJson<T>(value: string, fallback: T): T { try { return JSON.parse(value) as T; } catch { return fallback; } }
export function ValidationReportPanel({ deployment }: { deployment: Deployment }) {
  const report = parseJson<ValidationReport>(deployment.validationJson, { valid: false, files: 0, bytes: 0, errors: ['Unable to parse validationJson'] });
  const manifest = parseJson<Record<string, unknown>>(deployment.manifestJson, { parseError: 'Unable to parse manifestJson' });
  return <section className="validation-report-panel"><h2>Deployment v{deployment.version}</h2><div className="validation-report-panel__grid"><ManifestSummary manifest={manifest} /><ValidationSummary report={report} showJson /></div></section>;
}
