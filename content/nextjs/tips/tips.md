---
date: 2026-09-22
---
#### Instrumentation

```javascript
/**
 * register という名前の関数をプロジェクトのルート ディレクトリ
 * (または src フォルダーを使用している場合は src フォルダー内)
 * にあるinstrumentation.ts (または .js) ファイルからエクスポートすると、
 * 新しい Next.js サーバー インスタンスが作成されるたびにその関数が呼び出されます。
 */

// project/instrumentation.ts
import { init } from "package-init";

export function register() {
  init();
}
```

#### Environemnt Variables

```javascript
// load environment variables
.env // all environment
.env.local // general
.env.development // for next dev
.env.production // for next start

// bundle environment variables for the browser > need build with docker image
NEXT_PUBLIC_

```

#### Environment load order

```
1. process.env
2. .env.$(NODE_ENV).local
3. .env.local (not checked when NODE_ENV is test)
4. .env.$(NODE_ENV)
5. .env
```
