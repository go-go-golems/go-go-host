import { createBrowserRouter, Navigate } from 'react-router-dom';
import { OrgRedirectOrOnboarding } from './routing/OrgRedirectOrOnboarding';
import { OrgLayout } from './routing/OrgLayout';
import { RequireOrgAccess, RequireSession } from './routing/guards';
import { SitesPage } from '../pages/SitesPage';
import { CreateSitePage } from '../pages/CreateSitePage';

export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/app" replace /> },
  { path: '/app', element: <RequireSession><OrgRedirectOrOnboarding /></RequireSession> },
  { path: '/app/orgs/:orgId', element: <RequireSession><RequireOrgAccess><OrgLayout /></RequireOrgAccess></RequireSession>, children: [
    { index: true, element: <Navigate to="sites" replace /> },
    { path: 'sites', element: <SitesPage /> },
    { path: 'sites/new', element: <CreateSitePage /> },
    { path: 'agents', element: <div className="dashboard-panel"><h1>Agents</h1><p>Agents page wiring comes next.</p></div> },
    { path: 'audit', element: <div className="dashboard-panel"><h1>Audit</h1><p>Audit page wiring comes next.</p></div> },
  ] },
]);
