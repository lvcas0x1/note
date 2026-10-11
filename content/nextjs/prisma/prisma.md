---
date: 2026-09-22
---
### Installation

npm install prisma --save-dev

npx prisma init --datasource-provider sqlite

```javascript
# prisma/schema.prisma

model User {
  id    Int     @id @default(autoincrement())
  email String  @unique
  name  String?
  posts Post[]
}

model Post {
  id        Int     @id @default(autoincrement())
  title     String
  content   String?
  published Boolean @default(false)
  author    User    @relation(fields: [authorId], references: [id])
  authorId  Int
}
```

npx prisma migrate dev --name init

### reconfigure prisma schema

```shell
npx prisma generate --schema=./path/to/schema.prisma
```
