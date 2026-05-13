/**
 * docs-data.ts — Import go-go-host CLI docs as raw markdown strings.
 *
 * Vite's `?raw` suffix imports a file's text content as a default export string.
 * This module gathers every markdown doc from cmd/go-go-host/doc/ and
 * cmd/go-go-host-agent/doc/ into a typed, sorted catalogue.
 *
 * The docs have YAML frontmatter with Title, Slug, Short, and optional
 * SectionType fields. We parse those into a structured DocEntry.
 */

/* ── Raw markdown imports (Vite ?raw) ─────────────────────── */

// go-go-host docs
import hostGettingStarted from '../../../../../cmd/go-go-host/doc/getting-started.md?raw';
import hostDeveloperGuide from '../../../../../cmd/go-go-host/doc/developer-guide.md?raw';
import hostJsApiReference from '../../../../../cmd/go-go-host/doc/js-api-reference.md?raw';
import hostDeployWorkflow from '../../../../../cmd/go-go-host/doc/deploy-workflow.md?raw';
import hostCreateSiteWorkflow from '../../../../../cmd/go-go-host/doc/create-site-workflow.md?raw';
import hostRollbackWorkflow from '../../../../../cmd/go-go-host/doc/rollback-workflow.md?raw';
import hostAgentGuide from '../../../../../cmd/go-go-host/doc/agent-guide.md?raw';
import hostAgentSetup from '../../../../../cmd/go-go-host/doc/agent-setup.md?raw';
import hostLoginAndConfig from '../../../../../cmd/go-go-host/doc/login-and-config.md?raw';

// go-go-host-agent docs
import agentGettingStarted from '../../../../../cmd/go-go-host-agent/doc/getting-started.md?raw';
import agentGuide from '../../../../../cmd/go-go-host-agent/doc/agent-guide.md?raw';
import agentKeygenEnroll from '../../../../../cmd/go-go-host-agent/doc/keygen-enroll-deploy.md?raw';
import agentSignatureTroubleshooting from '../../../../../cmd/go-go-host-agent/doc/signature-troubleshooting.md?raw';

/* ── Types ─────────────────────────────────────────────────── */

export type DocSection = 'Tutorial' | 'GeneralTopic' | 'Example' | 'Application' | 'Other';

export interface DocEntry {
  /** Unique URL-safe slug. */
  slug: string;
  /** Human-readable title. */
  title: string;
  /** One-line description. */
  short: string;
  /** Section category for grouping. */
  section: DocSection;
  /** Source binary — "host" or "agent". */
  source: 'host' | 'agent';
  /** Raw markdown content (including frontmatter). */
  raw: string;
  /** Markdown body (frontmatter stripped). */
  body: string;
}

/* ── Frontmatter parser ────────────────────────────────────── */

function parseFrontmatter(raw: string): Record<string, string> {
  const fm: Record<string, string> = {};
  const match = raw.match(/^---\n([\s\S]*?)\n---/);
  if (!match) return fm;
  for (const line of match[1].split('\n')) {
    const idx = line.indexOf(':');
    if (idx === -1) continue;
    const key = line.slice(0, idx).trim();
    const val = line.slice(idx + 1).trim().replace(/^["']|["']$/g, '');
    fm[key] = val;
  }
  return fm;
}

function stripFrontmatter(raw: string): string {
  return raw.replace(/^---\n[\s\S]*?\n---\n?/, '');
}

function toSection(raw?: string): DocSection {
  if (!raw) return 'Other';
  if (raw === 'Tutorial' || raw === 'GeneralTopic' || raw === 'Example' || raw === 'Application') return raw;
  return 'Other';
}

function makeEntry(raw: string, source: 'host' | 'agent'): DocEntry {
  const fm = parseFrontmatter(raw);
  return {
    slug: fm.Slug ?? fm.slug ?? 'unknown',
    title: fm.Title ?? fm.title ?? 'Untitled',
    short: fm.Short ?? fm.short ?? '',
    section: toSection(fm.SectionType ?? fm.sectionType),
    source,
    raw,
    body: stripFrontmatter(raw),
  };
}

/* ── Catalogue ──────────────────────────────────────────────── */

const hostDocs: [string, string][] = [
  ['getting-started', hostGettingStarted],
  ['developer-guide', hostDeveloperGuide],
  ['js-api-reference', hostJsApiReference],
  ['deploy-workflow', hostDeployWorkflow],
  ['create-site-workflow', hostCreateSiteWorkflow],
  ['rollback-workflow', hostRollbackWorkflow],
  ['agent-guide', hostAgentGuide],
  ['agent-setup', hostAgentSetup],
  ['login-and-config', hostLoginAndConfig],
];

const agentDocs: [string, string][] = [
  ['agent-getting-started', agentGettingStarted],
  ['agent-guide', agentGuide],
  ['agent-keygen-enroll-deploy', agentKeygenEnroll],
  ['agent-signature-troubleshooting', agentSignatureTroubleshooting],
];

/** All docs sorted by section then title. */
export const docs: DocEntry[] = [
  ...hostDocs.map(([, raw]) => makeEntry(raw, 'host')),
  ...agentDocs.map(([, raw]) => makeEntry(raw, 'agent')),
].sort((a, b) => {
  const sectionOrder: DocSection[] = ['Tutorial', 'GeneralTopic', 'Example', 'Application', 'Other'];
  const sa = sectionOrder.indexOf(a.section);
  const sb = sectionOrder.indexOf(b.section);
  if (sa !== sb) return sa - sb;
  return a.title.localeCompare(b.title);
});

/** Lookup by slug. */
export function docBySlug(slug: string): DocEntry | undefined {
  return docs.find((d) => d.slug === slug);
}

/** Docs grouped by section. */
export function docsBySection(): Record<DocSection, DocEntry[]> {
  const groups: Record<string, DocEntry[]> = {};
  for (const d of docs) {
    (groups[d.section] ??= []).push(d);
  }
  return groups as Record<DocSection, DocEntry[]>;
}
