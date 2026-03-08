import { Probot, Context } from 'probot';
import { readNote } from './notes-reader';
import { formatPRComment, CommitContext, COMMENT_MARKER } from './formatter';

async function handlePullRequest(app: Probot, context: Context<'pull_request'>) {
  try {
    const owner = context.payload.repository.owner.login;
    const repo = context.payload.repository.name;
    const pull_number = context.payload.pull_request.number;

    const commitsResponse = await context.octokit.request(
      'GET /repos/{owner}/{repo}/pulls/{pull_number}/commits',
      { owner, repo, pull_number, per_page: 50 }
    );
    
    const commitContexts: CommitContext[] = [];
    for (const commit of commitsResponse.data) {
      const note = await readNote(context.octokit as any, owner, repo, commit.sha);
      const subject = commit.commit.message.split('\n')[0];
      commitContexts.push({
        sha: commit.sha,
        shortSha: commit.sha.slice(0, 7),
        subject,
        note,
      });
    }

    const commentBody = formatPRComment(commitContexts);
    if (!commentBody) return;

    // Prepend marker so we can find and update this comment later
    const body = `${COMMENT_MARKER}\n${commentBody}`;

    const commentsResponse = await context.octokit.request(
      'GET /repos/{owner}/{repo}/issues/{issue_number}/comments',
      { owner, repo, issue_number: pull_number, per_page: 100 }
    );
    
    const existingComment = commentsResponse.data.find(comment => 
      comment.body?.includes(COMMENT_MARKER)
    );

    if (existingComment) {
      await context.octokit.request(
        'PATCH /repos/{owner}/{repo}/issues/comments/{comment_id}',
        { owner, repo, comment_id: existingComment.id, body }
      );
    } else {
      await context.octokit.request(
        'POST /repos/{owner}/{repo}/issues/{issue_number}/comments',
        { owner, repo, issue_number: pull_number, body }
      );
    }
  } catch (error) {
    app.log.error(error);
  }
}

export default (app: Probot) => {
  app.on(['pull_request.opened', 'pull_request.synchronize'], async (context) => {
    await handlePullRequest(app, context);
  });
};
