---
date: 2026-09-22
---
#### Installation

```shell
sudo apt update && sudo apt upgrade -y

sudo apt install mysql-server

sudo systemctl is-active mysql
```

#### Post-Installation

```shell
sudo mysql_secure_installation
```

output:

```text
Securing the MySQL server deployment.

Connecting to MySQL using a blank password.

VALIDATE PASSWORD COMPONENT can be used to test passwords
and improve security. It checks the strength of password
and allows the users to set only those passwords which are
secure enough. Would you like to setup VALIDATE PASSWORD component?

Press y|Y for Yes, any other key for No: y

There are three levels of password validation policy:

LOW    Length >= 8
MEDIUM Length >= 8, numeric, mixed case, and special characters
STRONG Length >= 8, numeric, mixed case, special characters and dictionary file

Please enter 0 = LOW, 1 = MEDIUM and 2 = STRONG: 0

Skipping password set for root as authentication with auth_socket is used by default.
If you would like to use password authentication instead, this can be done with the "ALTER_USER" command.
See https://dev.mysql.com/doc/refman/8.0/en/alter-user.html#alter-user-password-management for more information.

By default, a MySQL installation has an anonymous user,
allowing anyone to log into MySQL without having to have
a user account created for them. This is intended only for
testing, and to make the installation go a bit smoother.
You should remove them before moving into a production
environment.

Remove anonymous users? (Press y|Y for Yes, any other key for No) : y
Success.


Normally, root should only be allowed to connect from
'localhost'. This ensures that someone cannot guess at
the root password from the network.

Disallow root login remotely? (Press y|Y for Yes, any other key for No) : y
Success.

By default, MySQL comes with a database named 'test' that
anyone can access. This is also intended only for testing,
and should be removed before moving into a production
environment.


Remove test database and access to it? (Press y|Y for Yes, any other key for No) : y
 - Dropping test database...
Success.

 - Removing privileges on test database...
Success.

Reloading the privilege tables will ensure that all changes
made so far will take effect immediately.

Reload privilege tables now? (Press y|Y for Yes, any other key for No) : y
Success.

All done!
```

#### Login to mysql

```shell
sudo mysql

SHOW DATABASES;
```

system databases:

```text
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+

1.information_schemaデータベースは、データベースのメタデータを保持しています。
データベースの構造（テーブル、カラム、インデックス、制約など）や設定に関する情報が含まれています。
データベースの構造を問い合わせるために使用される標準化された方法を提供します。

2.mysqlデータベースは、MySQLサーバー自体の設定とユーザーアカウント情報を保持しています。
ユーザーの権限、グローバル変数、およびその他のシステムレベルの設定が含まれています。

3.performance_schemaデータベースは、サーバーのパフォーマンス情報を提供します。
実行中のクエリ、サーバーのステータス、パフォーマンスカウンターなどが含まれています。
DBAやデベロッパーがパフォーマンスの問題を診断し、最適化するのに役立ちます。

4.sysデータベースは、パフォーマンス分析とトラブルシューティングを支援するためのビューとストアドプロシージャを提供します。
performance_schemaとinformation_schemaのデータをより使いやすく、わかりやすくするために設計されています。
```

#### Install MySQL Workbench Community (option)

```shell
sudo snap install mysql-workbench-community
```

#### Comment

```sql
# コメントです

-- コメントです

/*
コメントです
コメントです
コメントです
*/
```

#### DATABASE: CREATE, USE, SELECT, SHOW, DROP

```sql
CREATE DATABASES sample_db; -- DB作成
USE sample_db; -- 指定のDBを利用する
SELECT DATABASE() -- 利用中のDB名を表示する

SHOW DATABASES; -- DBを表示する

DROP DATABASE sample_db;　-- DBを削除する
```

#### TABLE: CREATE, SHOW, DROP, DESCRIBE

| Type     | Description                        |
| :------- | :--------------------------------- |
| INT      | 整数値                             |
| DECIMAL  | 小数                               |
| BOOLEAN  | 真偽値                             |
| CHAR     | 固定長文字列（最大：255 バイト）   |
| VARCHAR  | 可変長文字列（最大：16383 バイト） |
| DATE     | 日付                               |
| DATETIME | 時刻                               |
| TEXT     | 655355 文字までの文字列            |

```sql
CREATE TABLE users(
  id INT PRIMARY KEY, -- Primary Key: 重複が存在せず、行を一意に識別でき、空の行が存在しない列
  name VARCHAR(10),
  age INT,
  phone_number CHAR(13),
  message TEXT
);

SHOW TABLES;

DROP TABLE users;

DESCRIBE users; -- Table定義確認
```

```text
+--------------+-------------+------+-----+---------+-------+
| Field        | Type        | Null | Key | Default | Extra |
+--------------+-------------+------+-----+---------+-------+
| id           | int         | NO   | PRI | NULL    |       |
| name         | varchar(10) | YES  |     | NULL    |       |
| age          | int         | YES  |     | NULL    |       |
| phone_number | char(13)    | YES  |     | NULL    |       |
| message      | text        | YES  |     | NULL    |       |
+--------------+-------------+------+-----+---------+-------+

1.Field: テーブルの列の名前を示します。
2.Type: データ型を示します。各フィールドは特定のデータ型を持っており、それに応じて異なる種類のデータを格納できます。
3.Null: フィールドがNULL値を許可するかどうかを示します。YESはNULL値を許可し、NOはNULL値を許可しないことを示します。
4.Key: フィールドがテーブルのキー（例えば、主キー、外部キー、インデックス）の一部であるかどうかを示します。この列に PRI が表示されている場合、そのフィールドは主キーであることを示します。
5.Default: フィールドのデフォルト値を示します。新しいレコードを作成するときに値が提供されていない場合、このデフォルト値が使用されます。
6.Extra: フィールドに関連する追加情報を提供します。例えば、フィールドがオートインクリメントである場合、auto_incrementと表示されます。
```
