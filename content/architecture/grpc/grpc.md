---
date: 2026-09-22
---
## 環境構築

Install Homebrew

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Install golang&protobuf

```bash
brew install go
brew install protobuf

go version
protoc --version
```

Install Go Plugins

```bash
# https://grpc.io/docs/languages/go/quickstart/
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"

protoc-gen-go --version
protoc-gen-go-grpc --version
```

## Protocol Buffers

Protocol Buffers: Google によって 2008 年にオープンソース化されたスキーマ言語

- スキーマ言語: 処理ではなく、要素や属性などの構造を定義するための言語
- gRPC のデータフォーマットとして使用されている
- プログラミング言語からは独立している、様々な言語に変換可能
- バイナリ形式にシリアライズするため、サイズが小さく高速な通信が可能
- 型安全にデータやり取りが可能
- JSON に変換することも可能

Protocol Buffers 使った開発

- スキーマの定義
- 開発言語のオブジェクトを自動生成
- バイナリ形式へシリアライズ

### Proto ファイル

```proto
syntax = "proto3";

// write some comment ...

/*
block
comment
...
*/
```

### Message とは

複数のフィールドを持つことができる型定義

- それぞれのフィールドはスカラ型もしくはコンポジット型

各言語のコードとしてコンパイルした場合、構造体やクラスとして変換される  
一つの proto ファイルに複数の Message 型を定義することも可能

https://protobuf.dev/programming-guides/proto3/#scalar

```
message メッセージの名前 {
  型 フィールド名 = フィールド番号
  型 フィールド名 = フィールド番号
}
```

```proto
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 results_per_page = 3;
}
```

### Tag とは

- ProtocolBuffers ではフィールドはフィールド名ではなく、タグ番号によって識別される
- タグ番号の重複は一意の必要がある
- 最小値は 1、最大値は 2 の 29 - 1
- 19000 - 19999 は Protocol Buffers の予約番号のため使用不可
- 1 - 15 番は 1byte で表すことができるため、よく使うフィールドには 1-15 番を割り当てる
- タグは連番にする必要はないので、あまり使わないフィールドはあえて 16 番以後を割り当てることも可能
- タグ番号を予約するなど、安全に Protocol Buffers を使用する方法もある

### 列挙型

- フィールド名は大文字で記載する。
- 0 を数値のデフォルト値として使用できるようにするため、最初のフィールドのフィールド番号は 0 にする。

```
enum 名前 {
  フィールド名 = タグ番号;
  フィールド名 = タグ番号;
}
```

```proto
syntax = "proto3";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}
```

### Repeated

repeated フィールドでは配列を扱うことができます。

```
repeated 型 フィールド名 = フィールド番号;
```

```proto
syntax = "proto3";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}
```

### Map

- map では key-value の連想配列のような型を定義できます。
- key は string、numeric、bool のいずれかを使用できます。
- repeated フィールドは使用できません。

```
map<keyの型, Valueの型> フィールド名 = フィールド番号;
```

```proto
syntax = "proto3";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
  map<string, Project> project = 6;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}

message Project {}
```

### Oneof

Oneof は定義した複数の型のいずれかがを値として持ちます。

```
oneof フィールド名 {
  型 フィールド名 = フィールド番号;
  型 フィールド名 = フィールド番号;
}
```

```proto
syntax = "proto3";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
  map<string, Project> project = 6;
  oneof profile {
    string text = 7;
    Video video = 8;
  }
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}

message Project {}
message Video {}
```

### デフォルト値

定義しているフィールドに値が格納されていない場合に、型にはデフォルト値が格納されます。
|型|ゼロ値|
|-|-|
|strings|空の文字列|
|bytes|空のバイト|
|bools|false|
|numeric|0|
|enums|最初に定義された列挙値で、0 にする必要がある|
|repeated|空のリスト|

### Message のネスト

```proto
syntax = "proto3";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
  map<string, Company.Project> project = 6;
  oneof profile {
    string text = 7;
    Video video = 8;
  }
  Parent.Children nest = 9;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}

message Company {
  message Project {}
}
message Video {}
message Parent {
  message Children{}
}
```

### Import&Package

date.proto

```proto
syntax = "proto3";

package date;

message Date {
  int32 year = 1;
  int32 month = 2;
  int32 day = 3;
}
```

employee.proto

```proto
syntax = "proto3";

package employee;

import "proto/date.proto";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
  map<string, Company.Project> project = 6;
  oneof profile {
    string text = 7;
    Video video = 8;
  }
  Parent.Children nest = 9;
  date.Date birthday = 10;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}

message Company {
  message Project {}
}
message Video {}
message Parent {
  message Children{}
}
```

### Compile

```
-IPATH, --proto_path=PATH

--java_out=OUT_DIR          Generate Java source file.
--kotlin_out=OUT_DIR        Generate Kotlin file.
--pyi_out=OUT_DIR           Generate python pyi stub.
--python_out=OUT_DIR        Generate Python source file.
--rust_out=OUT_DIR          Generate Rust sources.
```

Golang Package の名前空間定義を追加する

employee.proto

```proto
syntax = "proto3";

package employee;

option go_package = "./pb";

import "proto/date.proto";

message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  Occupation occupation = 4;
  repeated string phone_number = 5;
  map<string, Company.Project> project = 6;
  oneof profile {
    string text = 7;
    Video video = 8;
  }
  Parent.Children nest = 9;
  date.Date birthday = 10;
}

enum Occupation {
  OCCUPATION_UNKNOW = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}

message Company {
  message Project {}
}
message Video {}
message Parent {
  message Children{}
}
```

date.proto

```proto
syntax = "proto3";

package date;

option go_package = "./pb";

message Date {
  int32 year = 1;
  int32 month = 2;
  int32 day = 3;
}
```

```bash
# compile from root directory

protoc -I. --go_out=. proto/*.proto

# download golang module

go mod tidy
```

以下のようなタイプ定義や Getter メソッド生成される

date.pb.go

```go
type Date struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Year  int32 `protobuf:"varint,1,opt,name=year,proto3" json:"year,omitempty"`
	Month int32 `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Day   int32 `protobuf:"varint,3,opt,name=day,proto3" json:"day,omitempty"`
}

func (x *Date) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Date) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *Date) GetDay() int32 {
	if x != nil {
		return x.Day
	}
	return 0
}
```

### Serialize&Deserialize

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ns1root/grpc/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id: 1,
		Name: "Azaria",
		Email: "demo@example.com",
		Occupation: pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project: map[string]*pb.Company_Project{"ProjectX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My name is Azaria",
		},
		Birthday: &pb.Date{
			Year: 2000,
			Month: 1,
			Day: 1,
		},
	}

	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Println("cannot serialize", err)
	}

	if err := os.WriteFile("demo.bin", binData, 0666); err != nil {
		log.Println("cannot write", err)
	}

	in, err := os.ReadFile("demo.bin")
	if err != nil {
		log.Println("cannot read file", err)
	}

	readEmployee := &pb.Employee{}
	err = proto.Unmarshal(in, readEmployee)
	if err != nil {
		log.Println("cannot deserialize", err)
	}

	fmt.Println(readEmployee)
}
```

### JSON Mapping

```go
package main

import (
	"fmt"
	"log"

	"github.com/ns1root/grpc/pb"
  "google.golang.org/protobuf/encoding/protojson"
)

func main() {
	employee := &pb.Employee{
		Id: 1,
		Name: "Suzuki",
		Email: "test@example.com",
		Occupation: pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project: map[string]*pb.Company_Project{"ProjectX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My name is Suzuki",
		},
		Birthday: &pb.Date{
			Year: 2000,
			Month: 1,
			Day: 1,
		},
	}

	marshaler := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	out, err := marshaler.Marshal(employee)
	if err != nil {
		log.Println("cannot marshal to json", err)
	}

	fmt.Println(string(out))

	readEmployee := &pb.Employee{}
	if err := protojson.Unmarshal(out, readEmployee); err != nil {
		log.Println("cannot unmarshal from json", err)
	}

	fmt.Println(readEmployee)
}
```

## gRPC

gRPC の g はバージョンごとに意味が異なる

- 1.64 'g' stands for 'grateful'

1. データフォーマットに Protocol Buffers を使用
2. IDL(Protocol Buffers)からサーバー側、クライアント側に必要なソースコードを自動生成
3. HTTP/2 使用される
   - ストリームの概念を導入
     - 一つの TCP 接続を用いて、複数のリクエスト/レスポンスのやり取りが可能。サーバの負荷経験
   - ヘッダーの圧縮
     - HPACK という圧縮方式で圧縮し、さらにキャッシュ行い、差分のみ送受信することで効率化
   - サーバプッシュ
     - クライアントからのリクエスト無しにサーバからデータ送信できる

### Service

- RPC(メソッド)の実装単位
  - サービス内に定義するメソッドがエンドポイントになる
  - 一つのサービス内に複数のメソッドを定義できる
- サービス名、メソッド名、引数、戻り値を定義する必要がある
- コンパイルして Go ファイルに変換すると、インターフェースに変換される

### 通信方式

- Unary RPC
- Server Streaming RPC
- Client Streaming RPC
- Bidrectional Streaming RPC

**Unary RPC**

- 1 リクエスト 1 レスポンス
- 通常の関数のコールのように扱う
- 用途: API など

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (SayHelloRequest) returns (SayHelloResponse);
}
```

**Server Streaming RPC**

- 1 リクエスト複数レスポンス
- クライアントはサーバから送信完了の信号が送信されるまでストリームのメッセージを読み続ける
- 用途: サーバからのプッシュ通知など

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (SayHelloRequest) returns (stream SayHelloResponse);
}
```

**Client Streaming RPC**

- 複数リクエスト 1 レスポンス
- サーバはクライアントからリクエスト完了の信号が送信されるまでストリームからメッセージを読み続け、レスポンスを返さない
- 用途: ファイルアップロードなど

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (stream SayHelloRequest) returns (SayHelloResponse);
}
```

**Bidrectional Streaming RPC**

- 複数リクエスト複数レスポンス
- クライアントとサーバのストリームが独立しており、リクエストとレスポンスはどのような順序でも良い
- 用途: チャットやオンライン対戦ゲームなど

```proto
message SayHelloRequest {}
message SayHelloResponse {}

service Greeter {
  rpc SayHello (stream SayHelloRequest) returns (stream SayHelloResponse);
}
```

### Compile

file.proto

```proto
syntax = "proto3";

package file;

option go_package = "./pb";

message ListFilesRequest {};
message ListFilesResponse {
  repeated string filenames = 1;
};

service FileService {
  rpc ListFiles (ListFilesRequest) returns (ListFilesResponse);
}
```

```bash
# compile for golang
protoc -I. --go_out=. --go-grpc_out=. proto/file.proto

# install modules
go mod tidy
```

以下のようなタイプ定義が生成される

```go
type ListFilesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

type ListFilesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filenames []string `protobuf:"bytes,1,rep,name=filenames,proto3" json:"filenames,omitempty"`
}
```

以下のような golang interface が生成される

```go
type FileServiceClient interface {
	ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesResponse, error)
}

type FileServiceServer interface {
	ListFiles(context.Context, *ListFilesRequest) (*ListFilesResponse, error)
	mustEmbedUnimplementedFileServiceServer()
}
```

### Interceptor

- メソッドの前後に処理を行うための仕組み
  - メソッドがリクエストを受け取る前、レスポンスを返した後のタイミングで任意の処理を割り込ませることが可能
- 認証やロギング、監視やバリデーションなど複数の RPC で共通して行たい処理で使用する
- Unary 用とストリーミング用がある
- サーバ側、クライアント側どちらにも対応している
  - UnaryServerInterceptor
  - StreamServerInterceptor
  - UnaryClientInterceptor
  - StreamClientInterceptor

### エラーハンドリング

https://grpc.io/docs/guides/error/

HTTP コード対照表
|gRPC|HTTP|内容|
|-|-|-|
|DEADLINE_EXCEEDED|504 Gateway Timeout|処理が完了前にタイムアウト|
|UNIMPLEMENTED|501 Not Implemented|処理が実装されていないか、サポートされていない|
|UNAVAILABLE|503 Service Unavailable|サービスが一時的に利用できない|
|UNKNOWN|500 Internal Server Error|不明なエラー|
|UNAUTHENTICATED|401 Unauthorized|リクエストに適切な認証情報がない|
|PERMISSION_DENIED|403 Forbidden|呼び出し元に処理を実行する権限がない|

### Sample

SSL+
https://github.com/ns1root/grpc
