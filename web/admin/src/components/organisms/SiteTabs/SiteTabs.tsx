import { NavLink } from 'react-router-dom';
import './SiteTabs.css';
export function SiteTabs({ basePath }: { basePath: string }) {
  return <nav className="site-tabs" aria-label="Site sections"><NavLink end to={basePath}>Overview</NavLink><NavLink to={`${basePath}/runtime`}>Runtime</NavLink><NavLink to={`${basePath}/settings`}>Settings</NavLink><NavLink to={`${basePath}/deployments`}>Deployments</NavLink></nav>;
}
