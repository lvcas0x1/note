---
date: 2026-09-22
---
##### Json to Parquet 変換ジョブ定義

```python
import sys
from awsglue.transforms import *
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext
from awsglue.context import GlueContext
from awsglue.job import Job

# ジョブの初期化
args = getResolvedOptions(sys.argv, ['JOB_NAME'])
sc = SparkContext()
glueContext = GlueContext(sc)
spark = glueContext.spark_session
job = Job(glueContext)
job.init(args['JOB_NAME'], args)

# 変換元と変換先のパスリストを定義
source_paths = [
    "s3://BUCKET_NAME/PREFIX",
]

destination_paths = [
    "s3://BUCKET_NAME/PREFIX",
]

# パスのリストをループしてそれぞれのパスに対してETL処理を行う
for source_path, destination_path in zip(source_paths, destination_paths):
    dynamic_frame = glueContext.create_dynamic_frame.from_options(
        format_options={"multiline": False},
        connection_type="s3",
        format="json",
        connection_options={"paths": [source_path], "recurse": True},
        transformation_ctx="dynamic_frame_{}".format(source_path)
    )

    glueContext.write_dynamic_frame.from_options(
        frame=dynamic_frame,
        connection_type="s3",
        format="glueparquet",
        connection_options={"path": destination_path, "partitionKeys": []},
        format_options={"compression": "snappy"},
        transformation_ctx="dynamic_frame_write_{}".format(destination_path)
    )

# ジョブのコミット
job.commit()
```

##### Glue Crawler で S3 からテーブル作成

```

```
