# jsqr

Query or evaluate JSON values.

### Install

```
go get github.com/ofabricio/jsqr
```

# Example

([Playground](https://go.dev/play/p/VoDXHRge1jF))

```go
package main

import "github.com/ofabricio/jsqr"

func main() {

    j := []byte(`
        {
            "data": { "store": "Grocery" },
            "tags": [
                { "name": "Fruit", "items": [{ "name": "Apple" }] },
                { "name": "Snack", "items": [{ "name": "Chips" }] },
                { "name": "Drink", "items": [{ "name": "Water" }, { "name": "Wine" }] }
            ]
        }
    `)

    a := jsqr.Get(j, `.data.store`)
    b := jsqr.Get(j, `.tags.[1].name == "Snack"`)
    c := jsqr.Get(j, `.tags.[ .name == "Drink" ].items.[0].name.(upper)`)

    fmt.Println(a) // "Grocery"
    fmt.Println(b) // true
    fmt.Println(c) // "WATER"
}
```

Also works with Go types:

([Playground](https://go.dev/play/p/F5wk8pc_gkW))

```go
package main

import "github.com/ofabricio/jsqr"

func main() {

    s := Store{
        Data: Data{
            Name: "Grocery",
        },
        Tags: []Tags{
            {Name: "Fruit", Items: []Item{{Name: "Apple"}}},
            {Name: "Snack", Items: []Item{{Name: "Chips"}}},
            {Name: "Drink", Items: []Item{{Name: "Water"}, {Name: "Wine"}}},
        },
    }

    a := jsqr.GetStruct(s, `.Data.Name`)
    b := jsqr.GetStruct(s, `.Tags.[1].Name == "Snack"`)
    c := jsqr.GetStruct(s, `.Tags.[ .Name == "Drink" ].Items.[0].Name.(upper)`)

    fmt.Println(a) // "Grocery"
    fmt.Println(b) // true
    fmt.Println(c) // "WATER"
}
```

Note that `jsqr.Get*(v, expr)` compiles the expression each time it is called. Avoid it with `jsqr.Compile(expr)`.

# Documentation

| Expression | Description |
| --- | --- | 
| `.` | Returns the current context. |
| `.a` | Returns a key value. For keys with characters other than `[a-zA-Z0-9_]` use `."a"`. |
| `.a == 100` | Boolean expression that returns either `true` or `false`. |
| `.[0]` | Returns the array item at the index. |
| `.[ .a > .b ]` | Returns the array item that matches the filter expression. |
| `==` `!=` `>=` `>` `<=` `<` | Comparison operators. |
| `eq` `ne` | Case insensitive comparison operators. |
| `&` `\|` | Logical operators: `.a == .b & .c == 100 \| .d == true`. |
| `.(func)` | Calls a function. See below. |

## Functions

| Function | Description |
| --- | --- | 
| `upper` | Converts a string to uppercase. |
| `lower` | Converts a string to lowercase. |
| `exists` | Tells if a context exists. |
