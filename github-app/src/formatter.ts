import { CommitNote } from './notes-reader';

export const COMMENT_MARKER = '<!-- agit-comment -->';

export interface CommitContext {
  sha: string;
  shortSha: string;
  subject: string;
  note: CommitNote | null;
}

export function formatPRComment(commits: CommitContext[]): string {
  const withNotes = commits.filter((c) => c.note !== null);
  if (withNotes.length === 0) return '';

  const total = commits.length;
  const count = withNotes.length;

  const lines: string[] = [];
  lines.push('## 🤖 Agent Context\n');
  lines.push(
    `**${count}/${total} commit${count === 1 ? '' : 's'} have agent context** · powered by [agit](https://github.com/madhurm/agit)\n`
  );
  lines.push('---\n');

  for (const ctx of withNotes) {
    lines.push(formatCommitNote(ctx));
  }

  return lines.join('\n');
}

export function formatCommitNote(ctx: CommitContext): string {
  const note = ctx.note!;
  const lines: string[] = [];

  lines.push(`### \`${ctx.shortSha}\` ${ctx.subject}\n`);

  // Confidence
  const pct = Math.round(note.confidence * 100);
  const confEmoji = note.confidence >= 0.8 ? '🟢' : note.confidence >= 0.5 ? '🟡' : '🔴';
  const confText = note.confidence_rationale
    ? `${confEmoji} ${pct}% — *"${note.confidence_rationale}"*`
    : `${confEmoji} ${pct}%`;

  // Agent
  let agentText = '';
  if (note.agent?.model || note.agent?.id) {
    const parts = [];
    if (note.agent.model) parts.push(note.agent.model);
    else if (note.agent.id) parts.push(note.agent.id);
    if (note.agent.session_id) parts.push(`(session: ${note.agent.session_id})`);
    agentText = parts.join(' ');
  }

  // Table rows
  const tableRows: string[] = [];
  if (note.confidence > 0) tableRows.push(`| **Confidence** | ${confText} |`);
  if (agentText) tableRows.push(`| **Agent** | ${agentText} |`);
  if (note.task) tableRows.push(`| **Task** | ${note.task} |`);

  if (tableRows.length > 0) {
    lines.push('| | |');
    lines.push('|---|---|');
    lines.push(...tableRows);
    lines.push('');
  }

  if (note.intent) {
    lines.push(`**Intent:** ${note.intent}\n`);
  }

  if (note.alternatives_considered?.length > 0) {
    lines.push('**Alternatives rejected:**');
    for (const alt of note.alternatives_considered) {
      lines.push(`- ✗ ${alt.approach} — ${alt.rejected_reason}`);
    }
    lines.push('');
  }

  if (note.key_decisions?.length > 0) {
    lines.push('**Key decisions:**');
    for (const kd of note.key_decisions) {
      lines.push(`- ${kd.decision}: ${kd.rationale}`);
    }
    lines.push('');
  }

  if (note.risks?.length > 0) {
    lines.push('**Risks:**');
    for (const risk of note.risks) {
      const emoji =
        risk.severity === 'high' ? '🔴' : risk.severity === 'medium' ? '🟡' : '🟢';
      lines.push(`- ${emoji} **[${risk.severity}]** ${risk.area} — ${risk.description}`);
    }
    lines.push('');
  }

  if (note.test_results) {
    const tr = note.test_results;
    lines.push(
      `**Tests:** ✅ ${tr.passed} passed, ❌ ${tr.failed} failed, ⏭ ${tr.skipped} skipped\n`
    );
  }

  if (note.unknowns?.length > 0) {
    lines.push('**❓ Unknowns:**');
    for (const u of note.unknowns) {
      lines.push(`- ${u}`);
    }
    lines.push('');
  }

  if (note.ripple_effects?.length > 0) {
    lines.push('**👀 Check these (ripple effects):**');
    for (const r of note.ripple_effects) {
      lines.push(`- ${r}`);
    }
    lines.push('');
  }

  lines.push('---\n');
  return lines.join('\n');
}
