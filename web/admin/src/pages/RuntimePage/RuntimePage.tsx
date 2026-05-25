import { useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { RuntimeStatusPanel } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetRuntimeQuery } from '../../services/goGoHostApi';
export function RuntimePage() { const { siteId } = useParams(); const runtime = useGetRuntimeQuery(siteId ?? '', { skip: !siteId, pollingInterval: 10000 }); if (runtime.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>; if (runtime.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load runtime" error={apiErrorMessage(runtime.error)} onRetry={runtime.refetch} /></section>; return <div className="runtime-page"><RuntimeStatusPanel runtime={runtime.data!} /><button type="button" data-part="btn" onClick={runtime.refetch}>Refresh runtime</button></div>; }
