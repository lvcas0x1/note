---
date: 2026-09-22
---
## Testing

math.go

```go
package lib

func Average(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}
	return int(total/len(s))
}
```

math_test.go

```go
package lib

import "testing"

var Debug bool = false

func TestAverage(t *testing.T) {
	if Debug {
		t.Skip("Write skip reason here")
	}

	v := Average([]int{1, 2, 3, 4, 5, 6, 7})
	if v !=3 {
		t.Error("Expected 3, got", v)
	}
}
```

```shell
go test ./...
```

```text
?       example/hello   [no test files]
--- FAIL: TestAverage (0.00s)
    math_test.go:14: Expected 3, got 4
FAIL
FAIL    example/hello/lib       0.198s
FAIL
```

### Time

```go
t := time.Now()
fmt.Println(t) // 2024-05-13 09:14:13.095881 +0900 JST m=+0.000175167
fmt.Println(t.Format(time.RFC3339)) // 2024-05-13T09:14:13+09:00
fmt.Println(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()) // 2024 May 13 9 15 9
fmt.Println(t.Date()) // 2024 May 13
```

### Regex

```go
match, _ := regexp.MatchString("a([a-z]+)e", "apple")
fmt.Println(match)

r := regexp.MustCompile("a([a-z]+)e")
ms := r.MatchString("apple")
fmt.Println(ms)

// s := "/view/test"
r2 := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
fs := r2.FindString("/view/test")
fmt.Println(fs)

url1 := r2.FindStringSubmatch("/view/test")
fmt.Println(url1, url1[0], url1[1], url1[2])

url2 := r2.FindStringSubmatch("/save/test")
fmt.Println(url2, url2[0], url2[1], url2[2])

url3 := r2.FindStringSubmatch("/edit/test")
fmt.Println(url3, url3[0], url3[1], url3[2])
```

### Sort

```go
i := []int{5, 3, 7, 2, 8, 1, 0}
s := []string{"t", "a", "y", "q"}
p := []struct {
	Name string
	Age int
}{
	{"Nancy", 20},
	{"Vera", 40},
	{"Mike", 30},
	{"Bob", 10},
}
fmt.Println(i, s, p)

sort.Ints(i)
fmt.Println(i)

sort.Strings(s)
fmt.Println(s)

sort.Slice(p, func(i, j int) bool { return p[i].Name < p[j].Name})
fmt.Println(p)

sort.Slice(p, func(i, j int) bool { return p[i].Age < p[j].Age})
fmt.Println(p)
```

### Iota

```go
const (
	c1 = iota
	c2 = iota
	c3 = iota
	c4
	c5
	c6
)

const (
	_ = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

func main() {
	fmt.Println(c1, c2, c3, c4, c5, c6) // 0 1 2 3 4 5
	fmt.Println(KB, MB, GB)
}
```

### Context

```go
func longProcess(ctx context.Context, ch chan string) {
	fmt.Println("run")
	time.Sleep(2 * time.Second)
	fmt.Println("finish")
	ch <- "result"
}

func main() {
	ch := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	// ctx := context.TODO()
	go longProcess(ctx, ch)

	CTXLOOP:
		for {
			select {
			case <- ctx.Done():
				fmt.Println(ctx.Err())
				// return
				break CTXLOOP
			case <- ch:
				fmt.Println("success")
				// return
				break CTXLOOP
			}
		}
	fmt.Println("context timeout only affect go routine func")
}
```

### IO & OS Read/Write

```go
content, err := os.ReadFile("main.go")
if err != nil {
	log.Fatal(err)
}
os.Stdout.Write(content)
// fmt.Println(string(content))

err = os.WriteFile("Hello.txt", []byte("Hello, Gophers!"), 0666)
if err != nil {
	log.Fatal(err)
}

err = os.WriteFile("tmp.go", content, 0666)
if err != nil {
	log.Fatal(err)
}

r := bytes.NewBuffer([]byte("abc"))
content, _ = io.ReadAll(r)
fmt.Println(string(content))
```

### Http

```go
// access url
resp, err := http.Get("http://example.com")
if err != nil {
	log.Fatal()
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))

// url check
base, err := url.Parse("http://example.com")
fmt.Println(base, err)

// create url
reference, _ := url.Parse("/test?a=1&b=2")
endpoint := base.ResolveReference(reference).String()
fmt.Println(endpoint)

// add header
req, _ := http.NewRequest("GET", endpoint, nil)
req.Header.Add("IF-None-Match", `W/"wyzzy"`)

// add parameter
q := req.URL.Query()
q.Add("c", "3&%")
fmt.Println(q)
fmt.Println(q.Encode())
req.URL.RawQuery = q.Encode()

var client *http.Client = &http.Client{}
fmt.Println(req)
// &{GET http://example.com/test?a=1&b=2&c=3 HTTP/1.1 1 1 map[If-None-Match:[W/"wyzzy"]] <nil> <nil> 0 [] false example.com map[] map[] <nil> map[]   <nil> <nil> <nil> {{}} <nil> [] map[]}
resp, _ = client.Do(req)
body, _ = io.ReadAll(resp.Body)
fmt.Println(string(body))
```

### Json Unmarshal & Marshal

```go
type T struct{}

type Person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age,string"` // Unmarshalの時はInt、MarshalのときはString
	Nicknames []string `json:"nicknames"`
	City      string   `json:"-"` // 非表示
	Number    int      `json:"number,omitempty"` // 0あるいは""、[]の時のみ非表示
	T         *T       `json:"T,omitempty"` // pointerにする必要ある
}

func main() {
	b := []byte(`{"name": "mike", "age":"20", "nicknames": ["a", "b", "c"], "city": "Tokyo", "number": 0}`)

	var p Person

	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Name, p.Age, p.Nicknames, p.City, p.Number, p.T) // mike 20 [a b c]  0 {}

	// Marshalは`json:"name"`で名前設定できる
	v, _ := json.Marshal(p)
	fmt.Println(string(v))
}
```

### Json Marshal Customize

```go
type T struct{}

type Person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age,string"` // Unmarshalの時はInt、MarshalのときはString
	Nicknames []string `json:"nicknames"`
	City      string   `json:"-"` // 非表示
	Number    int      `json:"number,omitempty"` // 0あるいは""、[]の時のみ非表示
	T         *T       `json:"T,omitempty"` // pointerにする必要ある
}

// MarshalJSONでJSONをカスタマイズする
func (p Person) MarshalJSON() ([]byte, error) {
	v, err := json.Marshal(&struct{
		Name string
	}{
		Name: "Mr." + p.Name,
	})
	return v, err
}


func main() {
	b := []byte(`{"name": "mike", "age":"20", "nicknames": ["a", "b", "c"], "city": "Tokyo", "number": 0}`)

	var p Person

	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Name, p.Age, p.Nicknames, p.City, p.Number, p.T) // mike 20 [a b c]  0 {}

	// Marshalは`json:"name"`で名前設定できる
	v, _ := json.Marshal(p)
	fmt.Println(string(v))
}
```

### Json Unmarshal Customize

```go
type T struct{}

type Person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age,string"` // Unmarshalの時はInt、MarshalのときはString
	Nicknames []string `json:"nicknames"`
	City      string   `json:"-"` // 非表示
	Number    int      `json:"number,omitempty"` // 0あるいは""、[]の時のみ非表示
	T         *T       `json:"T,omitempty"` // pointerにする必要ある
}

func (p *Person) UnmarshalJSON(b []byte) error {
	type Person2 struct{
		Name string
	}
	var p2 Person2
	err := json.Unmarshal(b, &p2)
	if err != nil {
		fmt.Println(err)
	}
	p.Name = p2.Name + "!"
	return err
}

func main() {
	b := []byte(`{"name": "mike", "age":"20", "nicknames": ["a", "b", "c"], "city": "Tokyo", "number": 0}`)

	var p Person

	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Name, p.Age, p.Nicknames, p.City, p.Number, p.T) // mike 20 [a b c]  0 {}

	// Marshalは`json:"name"`で名前設定できる
	v, _ := json.Marshal(p)
	fmt.Println(string(v))
}
```

### HMAC

```go
var DB = map[string]string {
	"User1Key": "User1Secret",
	"User2Key": "User2Secret",
}

func Server(apiKey, sign string, data []byte) {
	apiSecret := DB[apiKey]

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)

	expectedHMAC := hex.EncodeToString(h.Sum(nil))

	fmt.Println(sign == expectedHMAC)
}

func main() {
	const apiKey = "User1Key"
	const apiSecret = "User1Secret"

	// apiSecretでSHA256 Hashを作成
	h := hmac.New(sha256.New, []byte(apiSecret))
	// データ書き込む
	data := []byte("data")
	h.Write(data)
	// Signatureを作成する
	sign := hex.EncodeToString(h.Sum(nil))

	fmt.Println(sign)

	Server(apiKey, sign, data)
}
```
