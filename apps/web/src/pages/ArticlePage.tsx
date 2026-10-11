import { useParams } from 'react-router';
import { Placeholder } from '../components/Placeholder.tsx';
import { ApiError } from '../lib/api.ts';
import { ArticleView } from '../components/article/ArticleView.tsx';
import { readConfig } from '../lib/config.ts';
import { useArticle } from '../lib/queries.ts';

/** Article ID from the /a/* route; a trailing slash (static hosts) is ignored. */
export function useArticleId(): string {
  return (useParams()['*'] ?? '').replace(/\/+$/, '');
}

/** Shows one article, or a not-found message. */
export function ArticlePage() {
  const id = useArticleId();
  const article = useArticle(id);
  if (article.isPending) return <p className="text-muted">Loading…</p>;
  if (article.isError) {
    return article.error instanceof ApiError && article.error.status === 404 ? (
      <Placeholder title="Article not found">There is no article at this address.</Placeholder>
    ) : (
      <Placeholder title="Could not load this article">Try again later.</Placeholder>
    );
  }
  return <ArticleView article={article.data} basePath={readConfig().basePath} />;
}
