package jsqr

import (
	"testing"
)

func Test_parserParse(t *testing.T) {
	tc := []struct {
		give []byte
		when string
		then []byte
	}{
		{
			give: []byte(`{
					"data": { "store": "Grocery" },
					"tags": [
						{ "name": "Fruit", "items": [{ "name": "Apple" }] },
						{ "name": "Snack", "items": [{ "name": "Chips" }] },
						{ "name": "Drink", "items": [{ "name": "Water" }, { "name": "Wine" }] }
					]
				}`),
			when: `.tags.[ .name == "Drink" ].items.[0].name`,
			then: []byte(`"Water"`),
		},
		{
			give: []byte(`{ "a": null }`),
			when: `.a.(exists)`,
			then: []byte(`true`),
		},
		{
			give: []byte(`{}`),
			when: `.a.(exists)`,
			then: []byte(`false`),
		},
		{
			give: []byte(`{ "a": "HELLO world!" }`),
			when: `.a.(lower)`,
			then: []byte(`"hello world!"`),
		},
		{
			give: []byte(`{ "a": "Hello World!" }`),
			when: `.a.(upper)`,
			then: []byte(`"HELLO WORLD!"`),
		},
		{
			give: []byte(`"HellO!"`),
			when: `. ne "hello!"`,
			then: []byte(`false`),
		},
		{
			give: []byte(`"HellO!"`),
			when: `. eq "hello!"`,
			then: []byte(`true`),
		},
		{
			give: []byte(`false`),
			when: `. == false`,
			then: []byte(`true`),
		},
		{
			give: []byte(`true`),
			when: `. == true`,
			then: []byte(`true`),
		},
		{
			give: []byte(`null`),
			when: `. == null`,
			then: []byte(`true`),
		},
		{
			give: []byte(`1`),
			when: `. == 2 | . == 2`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. == 2 | . == 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`1`),
			when: `. == 1 | . == 2`,
			then: []byte(`true`),
		},
		{
			give: []byte(`1`),
			when: `. == 2 & . == 1`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. == 1 & . == 2`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. == 1 & . == 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`3`),
			when: `. <= 2`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. <= 2`,
			then: []byte(`true`),
		},
		{
			give: []byte(`2`),
			when: `. <= 2`,
			then: []byte(`true`),
		},
		{
			give: []byte(`0`),
			when: `. >= 1`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. >= 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`2`),
			when: `. >= 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`1`),
			when: `. > 1`,
			then: []byte(`false`),
		},
		{
			give: []byte(`2`),
			when: `. > 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`2`),
			when: `. != 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`1`),
			when: `. != 1`,
			then: []byte(`false`),
		},
		{
			give: []byte(`2`),
			when: `. == 1`,
			then: []byte(`false`),
		},
		{
			give: []byte(`1`),
			when: `. == 1`,
			then: []byte(`true`),
		},
		{
			give: []byte(`[ { "a": 1 }, { "a": 2 }, { "a": 3 } ]`),
			when: `.[ .a == 2 ]`,
			then: []byte(`{ "a": 2 }`),
		},
		{
			give: []byte(`[ 1, 2, 3 ]`),
			when: `.[3]`,
			then: []byte(``),
		},
		{
			give: []byte(`[ 1, 2, 3 ]`),
			when: `.[2]`,
			then: []byte(`3`),
		},
		{
			give: []byte(`[ 1, 2, 3 ]`),
			when: `.[1]`,
			then: []byte(`2`),
		},
		{
			give: []byte(`[ 1, 2, 3 ]`),
			when: `.[0]`,
			then: []byte(`1`),
		},
		{
			give: []byte(`3`),
			when: `.`,
			then: []byte(`3`),
		},
		{
			give: []byte(`{ "a": { "b": "}" } }`),
			when: `.a`,
			then: []byte(`{ "b": "}" }`),
		},
		{
			give: []byte(`{ "a": [ 3, [ 4 ], "]" ] }`),
			when: `.a`,
			then: []byte(`[ 3, [ 4 ], "]" ]`),
		},
		{
			give: []byte(`{ "a": null }`),
			when: `.b`,
			then: []byte(``),
		},
		{
			give: []byte(`{ "a": null }`),
			when: `.a`,
			then: []byte(`null`),
		},
		{
			give: []byte(`{ "a": false }`),
			when: `.a`,
			then: []byte(`false`),
		},
		{
			give: []byte(`{ "a": true }`),
			when: `.a`,
			then: []byte(`true`),
		},
		{
			give: []byte(`{ "a": 3 }`),
			when: `.a`,
			then: []byte(`3`),
		},
		{
			give: []byte(`{ "a": "b" }`),
			when: `.a`,
			then: []byte(`"b"`),
		},
		{
			give: []byte(`{ "a": "b" }`),
			when: `."a"`,
			then: []byte(`"b"`),
		},
	}
	for _, tc := range tc {

		e := Compile(tc.when)

		p := parser{}
		v := p.parse(e.expr, tc.give)

		if v.String() != string(tc.then) {
			t.Errorf("\nGot:\n%s\nExp:\n%s\nMsg:\n%s\n", v.String(), string(tc.then), tc)
		}
	}
}
