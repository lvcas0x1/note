---
date: 2026-09-22
---
lambda - update function code

```shell
cd ~/environment/api/dictate-function

zip dictate-function.zip app.py

aws lambda update-function-code \
--function-name dictate-function \
--zip-file fileb://dictate-function.zip

```

lambda - update function handler

```shell
aws lambda update-function-configuration \
--function-name dictate-function \
--handler app.lambda_handler

```

lambda - invoke function

```shell
aws lambda invoke \
--function-name dictate-function \
--payload fileb://event.json response.txt

```

lambda - get role arn

```shell
roleArn=$(aws iam list-roles --output text --query 'Roles[?contains(RoleName, `lambdaPollyRole`) == `true`].Arn')
```

lambda - create function

```python
folderName=createUpdate-functionfolderName=createUpdate-function

cd ~/environment/api/$folderName

zip $folderName.zip app.py

aws lambda create-function \
--function-name $folderName  \
--handler app.lambda_handler \
--runtime python3.11 \
--role $roleArn \
--environment Variables={TABLE_NAME=$notesTable} \
--zip-file fileb://$folderName.zip

```

lambda - update function

```shell
aws lambda update-function-configuration \
--function-name dictate-function \
--environment Variables="{MP3_BUCKET_NAME=$apiBucket, TABLE_NAME=$notesTable}"

```

lambda - amazon polly

```shell
# 貼り付ける際はインデントに注意してください
    pollyResponse = pollyClient.synthesize_speech(
        OutputFormat='mp3',
        Text = text,
        VoiceId = VoiceId
    )

```

lambda - upload file to s3

```shell
# 貼り付ける際はインデントに注意してください
    s3Client.upload_file(filePath,
    mp3Bucket,
    UserId+'/'+NoteId+'.mp3')

```
