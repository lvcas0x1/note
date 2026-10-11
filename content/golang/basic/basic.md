---
date: 2026-09-22
---
## Hello World

```shell
go init example/hello
```

```go
package main

import "fmt"

func main() {
  fmt.Println("hello, world!")
}
```

```shell
go run .
```

Call External Package
[Quote](https://pkg.go.dev/rsc.io/quote/v4)

```go
package main

import (
  "fmt"
  "rsc.io/quote/v4"
)

func main() {
  fmt.Println(quote.Go())
}
```

```shell
go mod tidy # add missing and remove unused modules
go run .
```

## 特殊関数 init

パッケージが初めて読み込まれるときに自動的に実行される  
パッケージレベルの変数の初期化、パッケージのセットアップに必要なリソース確保に使用される

```go
package main

import "fmt"

func init() {
  fmt.Println("func init")
}

func main() {
  fmt.Println("hello, world!")
}
```

## 変数

### var 変数

宣言のみ可能、初期値付与しない場合はデフォルト値入る

```go
// varは関数の外でも宣言できる
var (
  i int // 0
  f64 float64 // 0
  s string // ""
  t, f bool // false false
)
```

### := ショット

ショットは関数内しか使えない、右の値のタイプを関数に付与する

```go
i := 1
```

```go
// 明示的に関数タイプを宣言する
// デフォルトはfloat64になる
func main() {
  var f32 float32 = 1.2
	f64 := 1.2
	fmt.Printf("%T\n", f32)
	fmt.Printf("%T\n", f64)
}
```

### const 変数

```go
const (
	Username string = "username" // constは通常Type指定しない
	Password = "password"
	Pi = 3.14
	Xi
)

func main() {
	fmt.Printf("type=%T value=%v\n", Xi, Xi) // 宣言のみの場合デフォルトは一個上の変数を継承する？
}
```

### シフト

```go
// 2進数
fmt.Println(1 << 0) // 0001 0001
fmt.Println(1 << 1) // 0001 0010
fmt.Println(1 << 2) // 0001 0100
fmt.Println(1 << 3) // 0001 1000
```

### 文字列

```go
fmt.Println("Hello World") // Hello World
fmt.Println(string("Hello World"[0])) // H

var s string = "Hello World"
fmt.Println(strings.Replace(s, "H", "X", 1)) // Xello World
fmt.Println(s) // Hello World

s = strings.Replace(s, "H", "X", 1)
fmt.Println(s) // Xello World

fmt.Println(strings.Contains(s, "World")) // true

fmt.Println(`Test
Test`) // このまま出力される

fmt.Println("\"\"") // ""
fmt.Println(`""`) // ""
```

### 論理値

```go
fmt.Println(true && false) // false
fmt.Println(true || false) // true
fmt.Println(!true) // false
```

### 型変換

```go
var x int =  1
xx := float64(x)
fmt.Printf("%T %v %f\n", xx, xx, xx) // float64 1 1.000000

var y float64 =  1.2
yy := int(y)
fmt.Printf("%T %v %d\n", yy, yy, yy) // int 1 1

var s string = "10"
i, _ := strconv.Atoi(s)
fmt.Printf("%T %v\n", i, i) // int 10

h := "Hello World"
fmt.Println(h[0]) // 72
fmt.Println(string(h[0])) // H
```

### 配列

```go
var a [2]int
var b [2]int = [2]int{1, 2}

fmt.Println(a, b) // [0 0] [1 2]

// Sliceはサイズ変更できる、配列はできない
var c []int = []int{1, 2}
c = append(c, 3)
fmt.Println(c) // [1 2 3]
```

### スライス

```go
var n = []int{1, 2, 3, 4, 5, 6}
fmt.Println(n) // [1 2 3 4 5 6]
fmt.Println(n[2]) // [3]
fmt.Println(n[2:4]) // [3 4]
fmt.Println(n[:2]) // [1 2]
fmt.Println(n[2:]) // [3 4 5 6]
fmt.Println(n[:]) // [1 2 3 4 5 6]

n[0] = 100
fmt.Println(n) // [100 2 3 4 5 6]

var board = [][]int{
  {1, 2, 3},
  {4, 5, 6},
  {7, 8, 9},
}
fmt.Println(board) // [[1 2 3] [4 5 6] [7 8 9]]

n = append(n, 100)
fmt.Println(n) // [100 2 3 4 5 6 100]
```

### スライスの make と cap

```go
m := make([]int, 3)
fmt.Printf("len=%d cap=%d value=%v\n", len(m), cap(m), m) //len=3 cap=3 value=[0 0 0]

n := make([]int, 3, 5)
fmt.Printf("len=%d cap=%d value=%v\n", len(n), cap(n), n) //len=3 cap=5 value=[0 0 0]

n = append(n, 1, 2, 3, 4, 5)
fmt.Printf("len=%d cap=%d value=%v\n", len(n), cap(n), n) //len=8 cap=10 value=[0 0 0 1 2 3 4 5]

a := make([]int, 0) // 0のスライスをメモリに確保する
var b []int // メモリ確保しない
fmt.Printf("len=%d cap=%d value=%v\n", len(a), cap(a), a) //len=0 cap=0 value=[]
fmt.Printf("len=%d cap=%d value=%v\n", len(b), cap(b), b) //len=0 cap=0 value=[]

c := make([]int, 5)
for i := 0; i < 5; i ++ {
  c = append(c, i)
  fmt.Printf("len=%d cap=%d value=%v\n", len(c), cap(c), c)
  // len=6 cap=10 value=[0 0 0 0 0 0]
  // len=7 cap=10 value=[0 0 0 0 0 0 1]
  // len=8 cap=10 value=[0 0 0 0 0 0 1 2]
  // len=9 cap=10 value=[0 0 0 0 0 0 1 2 3]
  // len=10 cap=10 value=[0 0 0 0 0 0 1 2 3 4]
}
fmt.Printf("len=%d cap=%d value=%v\n", len(c), cap(c), c)
// len=10 cap=10 value=[0 0 0 0 0 0 1 2 3 4]

d := make([]int, 0, 5)
for i := 0; i < 5; i ++ {
  d = append(d, i)
  fmt.Printf("len=%d cap=%d value=%v\n", len(d), cap(d), d)
  // len=1 cap=5 value=[0]
  // len=2 cap=5 value=[0 1]
  // len=3 cap=5 value=[0 1 2]
  // len=4 cap=5 value=[0 1 2 3]
  // len=5 cap=5 value=[0 1 2 3 4]
}
fmt.Printf("len=%d cap=%d value=%v\n", len(d), cap(d), d)
// len=5 cap=5 value=[0 1 2 3 4]
```

### map

```go
// key string , value int
m := map[string]int{"apple": 100, "banana": 20}
fmt.Println(m) // map[apple:100 banana:20]
fmt.Println(m["apple"]) // 100

m["banana"] = 300
fmt.Println(m["banana"]) // 300

m["new"] = 500
fmt.Println(m) // map[apple:100 banana:300 new:500]

fmt.Println(m["nothing"]) // 0

v := m["apple"]
fmt.Println(v) // 100

v2, ok := m["nothing"]
fmt.Println(v2, ok) // 0 false

// make宣言、メモリ確保
m2 := make(map[string]int)
m2["orange"] = 10
fmt.Println(m2) // map[orange:10]

// Error: nil dereference in map update
/*
var m3 map[string]int
m3["strawberry"] = 5
fmt.Println(m3)
*/
```

### バイト型

```go
a := []byte{72, 73}
fmt.Println(a) // [72 73]
fmt.Println(string(a)) // HI

b := []byte("HI")
fmt.Println(b) // [72 73]
fmt.Println(string(b)) // HI
```

### 関数

```go
func calc(x int, y int) ( add int, minus int){
	add = x + y
	minus = x - y
	return add, minus
}

func main() {
	r1, r2 := calc(10, 20)
	fmt.Println(r1, r2)

	f := func() {
		fmt.Println("inner func 1")
	}
	f()

	func() {
		fmt.Println("inner func 2")
	}()
}
```

### クロージャー

```go
func circleArea(pi float64) (func(radius float64) float64) {
	return func(radius float64) float64 {
		return pi * radius * radius
	}
}

func main() {
	c1 := circleArea(3.14)
	fmt.Println(c1(2)) // 12.56

	c2 := circleArea(3)
	fmt.Println(c2(2)) // 12
}
```

### 可変長引数

```go
func foo(params ...int) {
	fmt.Println(len(params), params) // 2 [1 2]
	for _, param := range params {
		fmt.Println(param)
		// 1
		// 2
	}
}

func main() {
	// foo(1 ,2)

	s := []int{1, 2, 3}
	fmt.Println(s)

	foo(s...)
}
```

## ステートメント

### If

```go
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here, though
	return lim
}

func main() {
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)
}
```

### For

for, continue, break

```go
sum := 0
for i := 0; i < 10; i++ {
  sum += i
  if i == 3 {
    fmt.Println("continue")
    continue
  }
  if i > 5 {
    fmt.Println("break")
    break
  }
  fmt.Println(i)
}
fmt.Println(sum)
```

省略記法

```go
sum := 1
for sum < 10 {
  sum += sum
  fmt.Println(sum)
}
fmt.Println(sum)
```

無限ループ

```go
for {
  fmt.Println("loop")
}

```

### Range

slice

```go
l := []string{"python", "go", "java"}

for i, v := range l {
  fmt.Println(i, v) // 0 python ...
}

for _, v := range l {
  fmt.Println(v) // python go java
}
```

map

```go
m := map[string]int{"apple": 10, "banana": 20}
for k, v := range m {
  fmt.Println(k, v)
}

for k := range m {
  fmt.Println(k)
}

for _, v := range m {
  fmt.Println(v)
}
```

### Switch

```go
t := time.Now()
switch {
case t.Hour() < 12:
  fmt.Println("morning")
case t.Hour() < 17:
  fmt.Println("afternoon")
}
```

### Defer

関数実行の後に defer 関数を実行する

```go

func foo() {
	defer fmt.Println("defer foo")
	fmt.Println("foo")
}

func main() {
	foo()
	defer fmt.Println("defer main")
	fmt.Println("main")
}

```

Stacking Defer: 最初に defer に入れた関数が最後に実行される

```go
func main() {

	fmt.Println("run")
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)

	fmt.Println("finish")
}
```

File 処理後、Close 忘れずに defer 使う

```go
func main() {
	file, _ := os.Open("./hello.go")
	defer file.Close()

	data := make([]byte, 100)
	file.Read(data)
	fmt.Println(string(data))
}
```

### Log

```go
func LoggingSettings(logFile string) {
	logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}

func main() {
	LoggingSettings("logger.log")
	_, err := os.Open("xxxxxxx")
	if err != nil {
		log.Fatalln("Exit", err)
	}
}
```

### Error Handling

log.fatalln はプログラム終了のため、推奨しない

```go
LoggingSettings("logger.log")
file, err := os.Open("./hello.go")
if err != nil {
  log.Fatalln("Error", err)
}
defer file.Close()

data := make([]byte, 100)
count, err := file.Read(data)
if err != nil {
  log.Fatalln("Error")
}
fmt.Println(count, string(data))

// Chdirの返り値はErrorのみのため、一行で完結
if err = os.Chdir("dir"); err != nil {
  log.Fatalln("Error")
}
```

### Panic と Recover

Panic は推奨しない

```go
func thirdPartyConnectDB() {
	panic("Unable to connect database")
}

func save() {
	// save()関数実行完了後defer関数実行される
	// thirdPartyConnectDB()で起こったPanicをdeferのrecover()でキャッチする
	// キャッチしたPanicの内容をsに入れてPrintする
	// * deferはPanic起こる関数の上に書く必要ある
	defer func() {
		s := recover()
		fmt.Println(s) // Unable to connect database
	}()
	thirdPartyConnectDB()
}

func main() {
	save()
	fmt.Println("OK") // recoverあるため実行される
}
```

## ポインタ

### Pointer

```go
n := 100
fmt.Println(n) // n=100 値をプリントする
fmt.Println(&n) // nが入るメモリのアドレス

var p *int = &n // nのメモリのアドレスをpのメモリに保存する
fmt.Println(p) // nが入るメモリのアドレス
fmt.Println(*p) // 100 pのメモリに保存されているnのアドレスが指しているメモリの中身を表示する
*p = 1
fmt.Println(n) // 1 nの中身が1になったため
```

### make と new の違い

Pointer を返すか返さないかで make と new が違う

```go
	// slick
	s := make([]int, 0)
	fmt.Printf("%T\n", s) // []int
	// map
	m := make(map[string]int)
	fmt.Printf("%T\n", m) // map[string]int
	// channel
	ch := make(chan int)
	fmt.Printf("%T\n", ch) // chan int

	// pointer
	var p *int = new(int)
	fmt.Printf("%T\n", p) // *int
	fmt.Println(*p) // 0
	// struct
  // st := new(struct{})
  // st := &struct{}
	var st = new(struct{})
	fmt.Printf("%T\n", st) // *struct {}
```

### Struct

```go
type Vertex struct {
	X, Y int
	S string
}

func main() {
	v1 := Vertex{X: 1, Y: 2}
	fmt.Println(v1) // {1 2 }
	fmt.Println(v1.X, v1.Y) // 1 2

	v2 := Vertex{}
	fmt.Println(v2) // {0 0 }
	fmt.Printf("%T %v\n", v2, v2) // main.Vertex {0 0 }

	var v3 Vertex
	fmt.Println(v3) // {0 0 } v2と同じ
	fmt.Printf("%T %v\n", v3, v3) // main.Vertex {0 0 }

	v4 := new(Vertex)
	// var v4 = new(Vertex)
	fmt.Println(v4) // &{0 0 } アドレス
	fmt.Printf("%T %v\n", v4, v4) // *main.Vertex &{0 0 }

	v5 := &Vertex{}
	fmt.Println(v5) // &{0 0 } v4と同じ
	fmt.Printf("%T %v\n", v5, v5) // *main.Vertex &{0 0 }
}
```

```go
type Vertex struct {
	X, Y int
	S string
}

func changeVertex1(v Vertex) {
	v.X = 1000
}

func changeVertex2(v *Vertex) {
	// (*v).X = 1000
  v.X = 1000
}


func main() {
	v1 := Vertex{1, 2, "test"}
	changeVertex1(v1)
	fmt.Println(v1) // {1 2 test}

	v2 := &Vertex{1, 2, "test"}
	changeVertex2(v2)
	fmt.Println(v2) // &{1000 2 test}
	fmt.Println(*v2) // {1000 2 test}
}
```

## 構造体

### メソッド（オブジェクト指向）

```go
type Vertex struct {
	X, Y int
}
// メソッド
func (v Vertex) Area() int {
	return v.X * v.Y
}
// 関数
func Area(v Vertex) int {
	return v.X * v.Y
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(Area(v))

	fmt.Println(v.Area())
}
```

### Pointer Receiver & Value Receiver

```go
type Vertex struct {
	X, Y int
}
// Structの値を渡しているためValue Receiver
func (v Vertex) Area() int{
	return v.X * v.Y
}

// *VertexつけることでStructの中身書き換える
// Pointer Receiver
func (v *Vertex) Scale(i int) {
	v.X = v.X * i
	v.Y = v.Y * i
}

func main() {
	v := Vertex{3, 4}
	v.Scale(10)

	fmt.Println(v.Area())
}
```

### Construct

New 関数を定義する

```go
type Vertex struct {
	x, y int
}
// Structの値を渡しているためValue Receiver
func (v Vertex) Area() int{
	return v.x * v.y
}

// *VertexつけることでStructの中身書き換える
// Pointer Receiver
func (v *Vertex) Scale(i int) {
	v.x = v.x * i
	v.y = v.y * i
}

// コンストラクタ
func New(x, y int) *Vertex {
	return &Vertex{x, y}
}

func main() {
	// v := Vertex{3, 4}

	v := New(3, 4)
	v.Scale(10)

	fmt.Println(v.Area())
}
```

### Embedded(継承)

```go
type Vertex struct {
	x, y int
}
// Structの値を渡しているためValue Receiver
func (v Vertex) Area() int{
	return v.x * v.y
}

// *VertexつけることでStructの中身書き換える
// Pointer Receiver
func (v *Vertex) Scale(i int) {
	v.x = v.x * i
	v.y = v.y * i
}

// VertexのStructを継承する
type Vertex3D struct {
	Vertex
	z int
}
// Structの値を渡しているためValue Receiver
func (v Vertex3D) Area3D() int{
	return v.x * v.y * v.z
}

// *VertexつけることでStructの中身書き換える
// Pointer Receiver
func (v *Vertex3D) Scale3D(i int) {
	v.x = v.x * i
	v.y = v.y * i
	v.z = v.z * i
}

// コンストラクタ
func New(x, y ,z int) *Vertex3D {
	return &Vertex3D{Vertex{x, y}, z}
}

func main() {
	// v := Vertex{3, 4}

	v := New(3, 4, 5)
	v.Scale3D(10)

	//fmt.Println(v.Area())
	fmt.Println(v.Area3D()) // 60000
}
```

### Non-Struct

```go
// non-struct
type MyInt int

func (i MyInt) Double() int {
	return int(i * 2)
}

func main() {
	myInt := MyInt(10)
	fmt.Println(myInt.Double())
}
```

### Interface

メソッド

```go
type Human interface {
	Say()
}

type Person struct {
	Name string
}

func (p *Person) Say() {
	p.Name = "Mr." + p.Name
	fmt.Println(p.Name)
}

func main() {
	var mike Human = &Person{"Mike"}
	mike.Say()
}
```

Pointer & Receiver

```go
type Human interface {
	Say()
}

type Person struct {
	Name string
}

func (p *Person) Say() {
	p.Name = "Mr." + p.Name
	fmt.Println(p.Name)
}

func main() {
	var mike Human = &Person{"Mike"}
	mike.Say()
}
```

### Duck Typing

```go
type Human interface {
	Say() string
}

type Person struct {
	Name string
}

func (p *Person) Say() string {
	p.Name = "Mr." + p.Name
	return p.Name
}

func DriveCar(human Human) {
	if human.Say() == "Mr.Mike" {
		fmt.Println("Run")
	} else {
		fmt.Println("Get Out")
	}
}

func main() {
	var mike Human = &Person{"Mike"}
	DriveCar(mike)
}
```

### Type Assertion

```go
// interface{}どんな型でも受け付ける
func do(i interface{}) {
	// iをInterface型からIntに変換、タイプアサーション
	ii := i.(int)
	ii *= 2
	fmt.Println(ii)
}

func main() {
	do(2)
}
```

### Switch Type

```go
// interface{}どんな型でも受け付ける
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println(v * 2)
	case string:
		fmt.Println(v + "!")
	default:
		fmt.Printf("I don't know %T\n", v)
	}
}

func main() {
	do(2)
	do("Mike")
	do(true)
}
```

### Stringer

```go
type Person struct {
	Name string
	Age int
}

func (p Person) String() string{
	return fmt.Sprintf("My name is %v.", p.Name)
}

func main() {
	mike := Person{"Mike", 22}
	fmt.Println(mike)
}
```

### Custom Error

```go
type UserNotFound struct {
	Username string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("User not found: %v", e.Username)
}

func myFunc() error {
	ok := false
	if ok {
		return nil
	}
	return &UserNotFound{Username: "Mike"}
}

func main() {
	if err := myFunc(); err != nil {
		fmt.Println(err)
	}
}
```

## Go Routine

### sync.WaitGroup

```go
func normal(s string) {
	for i := 0; i< 5; i++ {
		fmt.Println(s)
	}
}

func gorouting(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i< 5; i++ {
		fmt.Println(s)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	normal("hello")
	go gorouting("world", &wg)

	wg.Wait()
}
```

### Channel

```go
func goRoutine1(s []int, c chan int) {
	sum := 0
	for _, v := range s{
		sum += v
	}
	c <- sum
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	c := make(chan int)

	go goRoutine1(s,c)
	// c <- sum が帰ってくるまで待つ
	x := <-c
	fmt.Println(x)
}
```

二つの GoRoutine で一つの Channel 共有するとき

make(chan int) は UnBuffered、Queue の役割

```go
func goRoutine1(s []int, c chan int) {
	sum := 0
	for _, v := range s{
		sum += v
	}
	c <- sum
}

func goRoutine2(s []int, c chan int) {
	sum := 0
	for _, v := range s{
		sum += v
	}
	c <- sum
}


func main() {
	s := []int{1, 2, 3, 4, 5}
	c := make(chan int) // 15 15

	go goRoutine1(s,c)
	go goRoutine2(s,c)

	x := <-c
	fmt.Println(x)
	y := <-c
	fmt.Println(y)
}
```

### Buffered Channels

```go
func main() {
	// Channelのバッファを2に指定
	ch := make(chan int, 2)
	ch <- 100
	fmt.Println(len(ch))
	ch <- 200
	fmt.Println(len(ch))

	// Bufferから一つ取り出す
	x := <-ch
	fmt.Println(x)
	// len(ch)が1になる
	fmt.Println(len(ch))

	ch <- 300
	fmt.Println(len(ch))
}
```

For Range で取り出すとき Close()

```go
func goRoutine1(s []int, c chan int) {
	sum := 0
	for _, v := range s{
		sum += v
		c <- sum
	}
	close(c)
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	c := make(chan int) // 15 15

	go goRoutine1(s,c)
	for i:= range c {
		fmt.Println(i)
	}
}
```

### Producer&Consumer

```go
func producer(ch chan int, i int) {
	ch <- i * 2
}

func consumer(ch chan int, wg *sync.WaitGroup) {
	for i := range ch {
		func (){
			defer wg.Done()
			fmt.Println("Process", i * 1000)
		}()
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	// Producer
	for i:=0; i<10; i++ {
		wg.Add(1)
		go producer(ch, i)
	}

	// Consumer
	go consumer(ch, &wg)
	wg.Wait()
	close(ch)
}
```

### FanIn & FanOut

```go
func producer(first chan int) {
	defer close(first)
	for i := 0; i < 10; i ++ {
		first <- i
	}
}

func multi2(first <-chan int, second chan<- int) {
	defer close(second)
	for i := range first {
		second <- i * 2
	}
}

func multi4(second <-chan int, third chan<- int) {
	defer close(third)
	for i := range second {
		third <- i * 4
	}
}

func main() {
	first := make(chan int)
	second := make(chan int)
	third := make(chan int)

	go producer(first)
	go multi2(first, second)
	go multi4(second, third)

	for result := range third {
		fmt.Println(result)
	}
}
```

### Channel & Select

```go
func goroutine1(ch chan string){
	for {
		ch <- "packet from 1"
		time.Sleep(1 * time.Second)
	}
}

func goroutine2(ch chan string){
	for {
		ch <- "packet from 2"
		time.Sleep(1 * time.Second)
	}
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go goroutine1(c1)
	go goroutine2(c2)

	for {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
	}
}
```

### Channel & Break

```go
func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)

	Loop:
		for {
			select {
			case t := <-tick:
				fmt.Println("tick.", t) // tick. 2024-05-11 14:33:18.982668 +0900 JST m=+0.101283376
			case <-boom:
				fmt.Println("BOOM!")
				break Loop
				// return Main関数から抜ける
			default:
				fmt.Println("    .")
				time.Sleep(50 * time.Millisecond)
			}
		}
	fmt.Println("####################")
}
```

### Sync.Mutex

```go
type Counter struct {
	v map[string]int
	mux sync.Mutex
}

func (c *Counter) Inc(key string) {
	c.mux.Lock()
  defer	c.mux.Unlock()

	c.v[key]++
}

func (c *Counter) Value(key string) int{
	c.mux.Lock()
  defer	c.mux.Unlock()

	return c.v[key]
}

func main() {

	c := Counter{v: make(map[string]int)}

	go func() {
		for i := 0; i < 10; i ++ {
			c.Inc("Key")
		}
	}()

	go func() {
		for i := 0; i < 10; i ++ {
			c.Inc("Key")
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println(&c, c.Value("Key")) // {map[Key:20] {0 0}} 20
}
```
