import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminUsersQuery } from '../../services/goGoHostApi';
import '../AdminOrgsPage/AdminOrgsPage.css';

export function AdminUsersPage() {
  const users = useListAdminUsersQuery();
  if (users.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (users.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load user inventory" error={apiErrorMessage(users.error)} /></section>;
  const rows = users.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Users</h1><p>Known users, platform-admin status, and organization membership counts.</p></header><table><thead><tr><th>User</th><th>ID</th><th>Platform admin</th><th>Orgs</th><th>Created</th><th>Last login</th></tr></thead><tbody>{rows.map((user) => <tr key={user.id}><td><strong>{user.displayName || user.email}</strong><br /><small>{user.email}</small></td><td><code>{user.id}</code></td><td>{user.platformAdmin ? 'yes' : 'no'}</td><td>{user.orgCount}</td><td><Timestamp value={user.createdAt} /></td><td><Timestamp value={user.lastLoginAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No users found.</p> : null}</section>;
}
