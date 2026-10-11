---
date: 2026-09-22
---
##### DynamoDB

テーブル設計：

```text
キーの種類：
・パーティションキー
・ソートキー
・セカンダリインデックス

グローバルセカンダリインデックス
ローカルセカンダリインデックス
```

テーブル作成：

```text
S - String
N - Number
B - Binary
BOOL - Boolean
NULL - Null
L - List
M - Map
SS - String Set
NS - Number Set
BS - Binary Set

```

dynamodb - boto3 init

```python
client = boto3.client('dynamodb')
```

dynamodb - create table

```python
# 貼り付ける際はインデントに注意してください
  response = ddbClient.create_table(
    AttributeDefinitions=[
      {
        'AttributeName': tableDefinition["partitionKey"],
        'AttributeType': 'S',
      },
      {
        'AttributeName': tableDefinition["sortKey"],
        'AttributeType': 'N',
      },
    ],
    KeySchema=[
      {
        'AttributeName': tableDefinition["partitionKey"],
        'KeyType': 'HASH',
      },
      {
        'AttributeName': tableDefinition["sortKey"],
        'KeyType': 'RANGE',
      },
    ],
    ProvisionedThroughput={
      'ReadCapacityUnits': int(tableDefinition["readCapacity"]),
      'WriteCapacityUnits': int(tableDefinition["writeCapacity"]),
    },
    TableName=tableDefinition["tableName"]
  )

```

dynamodb - waiter

```python
# 貼り付ける際はインデントに注意してください
    waiter = ddbClient.get_waiter('table_exists')
    waiter.wait(TableName=tableName)

```

dynamodb - load data

```python
# 貼り付ける際はインデントに注意してください
    table.put_item(
        Item={
            'UserId': note["UserId"],
            'NoteId': int(note["NoteId"]),
            'Note': note["Note"]
        }
    )

```

dynamodb - query data

```python
# 貼り付ける際はインデントに注意してください
    response = ddbClient.query(
        TableName=tableName,
        KeyConditionExpression='UserId = :userId',
        ExpressionAttributeValues={
            ':userId': {"S": qUserId}
        },
        ProjectionExpression="NoteId, Note"
    )

```

dynamodb - paginator

```python
# 貼り付ける際はインデントに注意してください
    paginator = ddbClient.get_paginator('scan')

    page_iterator = paginator.paginate(
      TableName=tableName,
      PaginationConfig={
        'PageSize': pageSize
      })

    pageNumber = 0
    for page in page_iterator:
      if page["Count"] > 0:
        pageNumber += 1
        print("Starting page " + str(pageNumber))
        printNotes(page['Items'])
        print("End of page " + str(pageNumber) + "\n")

```

dynamodb - add attribute

```python
# 貼り付ける際はインデントに注意してください
    response = ddbClient.update_item(
      TableName=tableName,
      Key={
        'UserId': {'S': qUserId},
        'NoteId': {'N': str(qNoteId)}
      },
      ReturnValues='ALL_NEW',
      UpdateExpression='SET Is_Incomplete = :incomplete',
      ExpressionAttributeValues={
          ':incomplete': {'S': 'Yes'}
      }
    )

```

dynamodb - partiql

```python
# 貼り付ける際はインデントに注意してください
    response = ddbClient.execute_statement(
        Statement="SELECT * FROM " + tableName + " WHERE UserId = ? AND NoteId = ?",
        Parameters=[
            {"S": qUserId},
            {"N": str(qNoteId)}
        ]
    )

```
