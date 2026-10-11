---
date: 2026-09-22
---
## GoDoc

Install Packages

```shell
go install golang.org/x/tools/cmd/godoc@latest
godoc -http:=6060
```

Write Documents

```go
/*
This is description of at package level
*/
package lib

// Average returns the average of a series of numbers
func Average(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}
	return int(total/len(s))
}
```

```go
var Debug bool = false

// Example
func ExampleAverage() {
	v := Average([]int{1, 2, 3, 4, 5, 6, 7})
	fmt.Println(v)
}

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

## Go mod

Add missing and remove unused modules

```shell
go mod tidy
```

## Go fmt

Format Go

```shell
gofmt -w hello.go
```
