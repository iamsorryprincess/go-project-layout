package redis

import (
	"encoding/json"
	"fmt"
	"testing"
)

// go test -bench . -benchmem -count=5 -cpuprofile cpu.out

type testStruct struct {
	UserID uint64 `redis:"user_id"`
	SiteID uint64 `redis:"site_id"`
	Name   string `redis:"name"`
	Phone  string `redis:"phone"`
	Email  string `redis:"email"`
}

func createData(count int) ([]string, error) {
	items := make([]string, count)

	for i := 0; i < count; i++ {
		item := testStruct{
			UserID: uint64(i),
			SiteID: uint64(i),
			Name:   fmt.Sprintf("name%d", i),
			Phone:  fmt.Sprintf("phone%d123 234 455", i),
			Email:  fmt.Sprintf("email%d@yandex.ru", i),
		}
		jsonData, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		items[i] = string(jsonData)
	}

	return items, nil
}

var testData, _ = createData(1000000)

func BenchmarkParse1(b *testing.B) {
	_, cErr := parse1[testStruct](testData)
	if cErr != nil {
		b.Fatal(cErr)
	}
}

func BenchmarkParse2(b *testing.B) {
	_, cErr := parse2[testStruct](testData)
	if cErr != nil {
		b.Fatal(cErr)
	}
}
