import { ErrorCallout } from '../../atoms/ErrorCallout';
import { JsonTree } from '../../atoms/JsonTree';
import type { ValidationReport } from '../../../services/types';
import './ValidationSummary.css';
export interface ValidationSummaryProps { report: ValidationReport; showJson?: boolean; }
export function ValidationSummary({ report, showJson = false }: ValidationSummaryProps) {
  return <section className="validation-summary" data-part="validation-summary" data-valid={report.valid}>
    <header><strong>{report.valid ? 'Validation passed' : 'Validation failed'}</strong><span>{report.files} files / {report.bytes} bytes</span></header>
    {report.errors?.length ? <ErrorCallout title="Validation errors" error={report.errors.join('\n')} /> : null}
    {report.warnings?.length ? <p className="validation-summary__warnings">Warnings: {report.warnings.join(', ')}</p> : null}
    {report.effectiveCapabilities?.length ? <p>Capabilities: {report.effectiveCapabilities.join(', ')}</p> : null}
    {showJson ? <JsonTree value={report} /> : null}
  </section>;
}
