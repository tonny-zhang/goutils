package mysql

import (
	"fmt"
	"testing"
)

type IData interface {
	getColumns() map[string]any
}
type Person struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

type Man struct {
	Person
	Type string `db:"type"`
}
type Man2 struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
	Type string `db:"type"`
}

func (p *Person) getColumns() map[string]any {
	data := make(map[string]any)
	data["name"] = &p.Name
	data["age"] = &p.Age
	return data
}

type CareerDetail struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

var columns = []string{
	"name", "age",
}

func structSource(columns []string, data IData) (list []any) {
	valuePointers := make([]any, len(columns))
	tagMap := data.getColumns()
	for i, c := range columns {
		v, ok := tagMap[c]
		if ok {
			valuePointers[i] = v
		} else {
			valuePointers[i] = &nilpointer
		}
	}
	return valuePointers
}
func BenchmarkStructScan(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Man

		strutForScan(columns, &p)
	}

}
func BenchmarkStructSource(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Person
		structSource(columns, &p)
	}
}
func BenchmarkStructScanCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Person

		strutForScanCache(columns, &p)
	}
}

func TestStructScanCache(t *testing.T) {
	var p Person

	fmt.Println(strutForScanCache(columns, &p))

	var p1 Person
	fmt.Println(strutForScanCache(columns, &p1))

	t.Fatal(true)
}
