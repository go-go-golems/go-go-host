import { createBrowserRouter, Navigate } from 'react-router-dom';
import { OrgRedirectOrOnboarding } from './routing/OrgRedirectOrOnboarding';
import { OrgLayout } from './routing/OrgLayout';
import { SiteLayout } from './routing/SiteLayout';
import { RequireOrgAccess, RequireSession } from './routing/guards';
import { SitesPage } from '../pages/SitesPage';
import { CreateSitePage } from '../pages/CreateSitePage';
import { SiteOverviewPage } from '../pages/SiteOverviewPage';
import { RuntimePage } from '../pages/RuntimePage';
import { DeploymentsPage } from '../pages/DeploymentsPage';
import { DeploymentDetailPage } from '../pages/DeploymentDetailPage';
import { AgentsPage } from '../pages/AgentsPage';
import { AuditPage } from '../pages/AuditPage';

export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/app" replace /> },
  { path: '/app', element: <RequireSession><OrgRedirectOrOnboarding /></RequireSession> },
  { path: '/app/orgs/:orgId', element: <RequireSession><RequireOrgAccess><OrgLayout /></RequireOrgAccess></RequireSession>, children: [
    { index: true, element: <Navigate to="sites" replace /> },
    { path: 'sites', element: <SitesPage /> },
    { path: 'sites/new', element: <CreateSitePage /> },
    { path: 'sites/:siteId', element: <SiteLayout />, children: [
      { index: true, element: <SiteOverviewPage /> },
      { path: 'runtime', element: <RuntimePage /> },
      { path: 'deployments', element: <DeploymentsPage /> },
      { path: 'deployments/:deploymentId', element: <DeploymentDetailPage /> },
    ] },
    { path: 'agents', element: <AgentsPage /> },
    { path: 'audit', element: <AuditPage /> },
  ] },
]);
