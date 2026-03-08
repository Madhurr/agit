import { Octokit } from '@octokit/core';

export interface AgentInfo {
  id: string;
  model: string;
  session_id: string;
}

export interface Alternative {
  approach: string;
  rejected_reason: string;
}

export interface KeyDecision {
  decision: string;
  rationale: string;
}

export interface Risk {
  area: string;
  description: string;
  severity: 'low' | 'medium' | 'high';
}

export interface TestResults {
  passed: number;
  failed: number;
  skipped: number;
  command: string;
}

export interface CommitNote {
  schema_version: string;
  commit_hash: string;
  agent: AgentInfo;
  task: string;
  intent: string;
  confidence: number;
  confidence_rationale: string;
  alternatives_considered: Alternative[];
  key_decisions: KeyDecision[];
  risks: Risk[];
  context_consulted: string[];
  test_results?: TestResults;
  ripple_effects: string[];
  unknowns: string[];
}

export async function readNote(
  octokit: Octokit,
  owner: string,
  repo: string,
  commitSha: string
): Promise<CommitNote | null> {
  try {
    // Step 1: Get notes ref SHA (singular /git/ref/ endpoint returns a single object)
    const { data: refData } = await octokit.request('GET /repos/{owner}/{repo}/git/ref/{ref}', {
      owner,
      repo,
      ref: 'notes/agit',
    });

    if (!refData.object?.sha) return null;

    // Step 2: List all blobs in notes tree
    const { data: treeData } = await octokit.request('GET /repos/{owner}/{repo}/git/trees/{sha}', {
      owner,
      repo,
      sha: refData.object.sha,
      recursive: '1' as unknown as boolean,
    });

    if (!treeData.tree) return null;

    // Step 3: Find blob for specific commit SHA
    const noteEntry = treeData.tree.find(
      (entry: any) => entry.path === commitSha && entry.type === 'blob'
    );

    if (!noteEntry?.sha) return null;

    // Step 4: Get blob content
    const { data: blobData } = await octokit.request('GET /repos/{owner}/{repo}/git/blobs/{sha}', {
      owner,
      repo,
      sha: noteEntry.sha,
    });

    if (blobData.encoding !== 'base64') return null;

    // Step 5-6: Decode and parse JSON
    const content = Buffer.from(blobData.content, 'base64').toString('utf8');
    return JSON.parse(content) as CommitNote;
  } catch (error: any) {
    if (error.status === 404) return null;
    throw error;
  }
}
