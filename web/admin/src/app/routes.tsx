import { createBrowserRouter, Navigate } from 'react-router-dom';
import { OrgRedirectOrOnboarding } from './routing/OrgRedirectOrOnboarding';
import { OrgLayout } from './routing/OrgLayout';
import { SiteLayout } from './routing/SiteLayout';
import { AdminLayout } from './routing/AdminLayout';
import { RequireOrgAccess, RequirePlatformAdmin, RequireSession } from './routing/guards';
import { SitesPage } from '../pages/SitesPage';
import { CreateSitePage } from '../pages/CreateSitePage';
import { SiteOverviewPage } from '../pages/SiteOverviewPage';
import { RuntimePage } from '../pages/RuntimePage';
import { DeploymentsPage } from '../pages/DeploymentsPage';
import { DeploymentDetailPage } from '../pages/DeploymentDetailPage';
import { AgentsPage } from '../pages/AgentsPage';
import { AuditPage } from '../pages/AuditPage';
import { UsagePage } from '../pages/UsagePage';
import { MembersPage } from '../pages/MembersPage';
import { AdminOverviewPage } from '../pages/AdminOverviewPage';
import { AdminRuntimesPage } from '../pages/AdminRuntimesPage';
import { AdminOrgsPage } from '../pages/AdminOrgsPage';
import { AdminUsersPage } from '../pages/AdminUsersPage';
import { AdminSitesPage } from '../pages/AdminSitesPage';
import { AdminDeploymentsPage } from '../pages/AdminDeploymentsPage';
import { AdminDeploymentDetailPage } from '../pages/AdminDeploymentDetailPage';
import { AdminAgentsPage } from '../pages/AdminAgentsPage';
import { AdminAuditPage } from '../pages/AdminAuditPage';

export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/app" replace /> },
  { path: '/app', element: <RequireSession><OrgRedirectOrOnboarding /></RequireSession> },
  { path: '/admin', element: <RequireSession><RequirePlatformAdmin><AdminLayout /></RequirePlatformAdmin></RequireSession>, children: [
    { index: true, element: <Navigate to="overview" replace /> },
    { path: 'overview', element: <AdminOverviewPage /> },
    { path: 'runtimes', element: <AdminRuntimesPage /> },
    { path: 'orgs', element: <AdminOrgsPage /> },
    { path: 'users', element: <AdminUsersPage /> },
    { path: 'sites', element: <AdminSitesPage /> },
    { path: 'deployments', element: <AdminDeploymentsPage /> },
    { path: 'deployments/:deploymentId', element: <AdminDeploymentDetailPage /> },
    { path: 'agents', element: <AdminAgentsPage /> },
    { path: 'audit', element: <AdminAuditPage /> },
    { path: '*', element: <AdminOverviewPage /> },
  ] },
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
    { path: 'usage', element: <UsagePage /> },
    { path: 'members', element: <MembersPage /> },
  ] },
]);
