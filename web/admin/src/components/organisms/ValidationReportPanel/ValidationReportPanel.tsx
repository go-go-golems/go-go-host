import { ValidationSummary } from '../../molecules/ValidationSummary';
import { ManifestSummary } from '../../molecules/ManifestSummary';
import type { Deployment } from '../../../services/types';
import { parseManifest, parseValidationReport } from '../../../services/json';
import './ValidationReportPanel.css';
export function ValidationReportPanel({ deployment }: { deployment: Deployment }) {
  const report = parseValidationReport(deployment.validationJson);
  const manifest = parseManifest(deployment.manifestJson);
  return <section className="validation-report-panel"><h2>Deployment v{deployment.version}</h2><div className="validation-report-panel__grid"><ManifestSummary manifest={manifest} /><ValidationSummary report={report} showJson /></div></section>;
}
