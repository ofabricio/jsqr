package jsqr

import (
	"reflect"
	"testing"
)

func TestGetStruct(t *testing.T) {

	give := struct {
		Name   string
		Age    int
		Zero   int
		Score  float64
		Nil    *string
		NotNil any
		Bool   bool
		Char   string
		Quote  string
		Data   struct{ Name string }
		Tags   []struct {
			Name  string
			Items []struct{ Name string }
		}
	}{
		Name:   "Alice",
		Age:    33,
		Zero:   0,
		Score:  100.5,
		Nil:    nil,
		NotNil: t,
		Bool:   true,
		Char:   "C",
		Quote:  `"Hello"`,
		Data: struct{ Name string }{
			Name: "Grocery",
		},
		Tags: []struct {
			Name  string
			Items []struct{ Name string }
		}{
			{Name: "Fruit", Items: []struct{ Name string }{{Name: "Apple"}}},
			{Name: "Snack", Items: []struct{ Name string }{{Name: "Chips"}}},
			{Name: "Drink", Items: []struct{ Name string }{{Name: "Water"}, {Name: "Wine"}}},
		},
	}

	tc := []struct {
		when string
		then any
	}{
		{when: `.Quote == "\"Hello\""`, then: true},
		{when: `.Nil.(exists)`, then: false},
		{when: `.NotNil.(exists)`, then: true},
		{when: `.Zero.(exists)`, then: true},
		{when: `.Age.(exists)`, then: true},
		{when: `.Name.(exists)`, then: true},
		{when: `.Name.(upper)`, then: "ALICE"},
		{when: `.Name.(lower)`, then: "alice"},
		{when: `.Char ne "C"`, then: false},
		{when: `.Char ne "c"`, then: false},
		{when: `.Char eq "C"`, then: true},
		{when: `.Char eq "c"`, then: true},
		{when: `.Char > "D"`, then: false},
		{when: `.Char > "C"`, then: false},
		{when: `.Char > "B"`, then: true},
		{when: `.Char >= "D"`, then: false},
		{when: `.Char >= "C"`, then: true},
		{when: `.Char >= "B"`, then: true},
		{when: `.Char > "D"`, then: false},
		{when: `.Char > "C"`, then: false},
		{when: `.Char > "B"`, then: true},
		{when: `.Char != "C"`, then: false},
		{when: `.Char == "C"`, then: true},
		{when: `.Age <= 34`, then: true},
		{when: `.Age <= 33`, then: true},
		{when: `.Age <= 32`, then: false},
		{when: `.Age < 34`, then: true},
		{when: `.Age < 33`, then: false},
		{when: `.Age < 32`, then: false},
		{when: `.Age > 34`, then: false},
		{when: `.Age > 33`, then: false},
		{when: `.Age > 32`, then: true},
		{when: `.Age >= 34`, then: false},
		{when: `.Age >= 33`, then: true},
		{when: `.Age >= 32`, then: true},
		{when: `.Tags.[0].Name == .Nil`, then: false},
		{when: `.Tags.[0].Name == .Tags.[1].Name`, then: false},
		{when: `.Tags.[0].Name == .Tags.[0].Name`, then: true},
		{when: `.Bool == true`, then: true},
		{when: `.Bool != true`, then: false},
		{when: `.Nil == null`, then: true},
		{when: `.Nil != null`, then: false},
		{when: `.NotNil != null`, then: true},
		{when: `.NotNil == null`, then: false},
		{when: `.Data.Name`, then: "Grocery"},
		{when: `.Tags.[1].Name`, then: "Snack"},
		{when: `.Tags.[1].Name == "Snack"`, then: true},
		{when: `.Tags.[2].Items.[0].Name`, then: "Water"},
		{when: `.Tags.[ .Name == "Drink" ].Items.[0].Name`, then: "Water"},
		{when: `.Score == 100.5`, then: true},
		{when: `.Score == 100.5001`, then: false},
		{when: `.Age == 33`, then: true},
		{when: `.Age == 31`, then: false},
		{when: `.Age != 33`, then: false},
		{when: `.Age != 31`, then: true},
		{when: `.Age`, then: 33},
		{when: `.Unknown`, then: nil},
		{when: `.Name`, then: "Alice"},
	}
	for _, tc := range tc {

		v := GetStruct(&give, tc.when)

		if !reflect.DeepEqual(v, tc.then) {
			t.Errorf("\nGot:\n%v\nExp:\n%v\nMsg:\n%s\n", v, tc.then, tc)
		}
	}
}
