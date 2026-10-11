---
date: 2026-09-22
---
## GraphQL

REST API の問題点

1. エンドポイントが大量に作成される
2. リソースの取得に柔軟性がない

GraphQL スキーマを独自のスキーマ定義言語で定義し、独自のクエリ言語でリクエストを行う  
特定の言語やフレームワークに依存しない  
三つの操作をサポートしている

- Query データの問い合わせ
- Mutation データの書き換え
- Subscription データの購読

GraphQL の特徴

1. 単一エンドポイント
2. 柔軟なデータ取得でオーバーフェッチ、アンダーフェッチの回避
3. 型情報を活用できる

GraphQL の問題点

1. メジャーではないため、情報が少ない
2. サーバ側の負荷予測が難しい
3. キャッシュが複雑になりやすい

### GraphQL スキーマ

- GraphQL API が受け付けるクエリの内容が定義された仕様書
- GraphQL Schema Language という独自の言語で記述される
- 拡張子には.graphqls, .graphql, .gql などがある

#### オブジェクト型

!つけることで Null を許容しない

```gql
type Object {
  Field_1: Type
  Field_2: Type!
  ...
}
```

#### スカラー型

- Int: 整数型
- Float: 浮動小数点数型
- String: 文字列型
- Boolean: 論理型
- ID: ユニークな文字列型

```
Type Task {
  id: Int
  name: String!
  dueDate: String!
  Status: String!
  description: String
}
```

#### リスト型

フィールドが複数の値のリストからなる場合に使用される

```gql
type User {
  tasks: [Task]!
}
```

!マークの場所
|定義方法|意味|
|-|-|
|[String]!|リスト自体必ず存在するが、Null 要素を持つ可能性がある|
|[String!]|リスト自体は Null の可能性があるが、Null 要素を持つことはない|
|[String!]!|リスト自体は必ず存在し、Null 要素も持たない|
|[String]|リスト自体は Null の可能性があり、Null 要素を含む可能性がある|

#### 列挙型

あらかじめ列挙された値のみ取ることができる

```gql
enum Color {
  RED
  BLUE
  GREEN
}
```

#### ユニオン型

複数の型のうちいずれかの型であることを表す

```gql
union Vehicle = Car | Bike
```

#### Query 型

GraphQL API が提供する Query のエントリポイントを定義

```gql
type Query {
  getTasks: [Task!]!
  getTask(id: Int!): Task //引数渡すことができる
}
```

#### Mutation 型

GraphQL API が提供する Mutation のエントリポイントを定義

```gql
type Mutation {
  createTask(id: Int!, name: String!): Task!
  deleteTask(id: Int!): Task!
}
```
