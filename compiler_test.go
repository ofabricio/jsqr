package jsqr

import (
	"bufio"
	"fmt"
	"strings"
)

func Example_compilerCompile() {

	s := `
		true
		false
		null
		123
		""

		.
		."a"
		.a
		."a"."b"
		.a.b
		.a0_B5_z99.世界.世.界
		.[0]

		.a == .b
		.a == .b & .c == .d
		.a == .b & .c == .d & .e == .f
		.a == .b | .c == .d
		.a == .b | .c == .d | .e == .f
		.a == .b & .c == .d | .e == .f
		.a == .b | .c == .d & .e == .f
		( .a == .b | .c == .d ) & .e == .f
		.a == .b | ( .c == .d & .e == .f )

		0 != 1
		1 >= 2
		3 > 4
		5 <= 6
		7 < 8
		"ab" eq "AB"
		"ab" ne "AB"

		.a.[.a==.b].b.[.a==.b]
		.a.[ .a == .b ].b.[ .a == "b" ]

		.(upper)
		.(upper).(lower).[ .(upper) == "HI" ]
		.(exists)
	`

	for _, line := range lines(s) {
		c := compiler{}
		v := c.Compile(line)
		fmt.Println(line)
		print(v.expr)
	}

	// Output:
	// true
	// [Bool] true
	// false
	// [Bool] false
	// null
	// [Null]
	// 123
	// [Num] 123
	// ""
	// [Str] ""
	// .
	// [Path]
	//     [This]
	// ."a"
	// [Path]
	//     [Key]
	//         [Str] "a"
	// .a
	// [Path]
	//     [Idn] a
	// ."a"."b"
	// [Path]
	//     [Key]
	//         [Str] "a"
	//     [Key]
	//         [Str] "b"
	// .a.b
	// [Path]
	//     [Idn] a
	//     [Idn] b
	// .a0_B5_z99.世界.世.界
	// [Path]
	//     [Idn] a0_B5_z99
	//     [Idn] 世界
	//     [Idn] 世
	//     [Idn] 界
	// .[0]
	// [Path]
	//     [Index] 0
	// .a == .b
	// [EQ]
	//     [Path]
	//         [Idn] a
	//     [Path]
	//         [Idn] b
	// .a == .b & .c == .d
	// [And]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [EQ]
	//         [Path]
	//             [Idn] c
	//         [Path]
	//             [Idn] d
	// .a == .b & .c == .d & .e == .f
	// [And]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [And]
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//         [EQ]
	//             [Path]
	//                 [Idn] e
	//             [Path]
	//                 [Idn] f
	// .a == .b | .c == .d
	// [Or]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [EQ]
	//         [Path]
	//             [Idn] c
	//         [Path]
	//             [Idn] d
	// .a == .b | .c == .d | .e == .f
	// [Or]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [Or]
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//         [EQ]
	//             [Path]
	//                 [Idn] e
	//             [Path]
	//                 [Idn] f
	// .a == .b & .c == .d | .e == .f
	// [Or]
	//     [And]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Path]
	//                 [Idn] b
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//     [EQ]
	//         [Path]
	//             [Idn] e
	//         [Path]
	//             [Idn] f
	// .a == .b | .c == .d & .e == .f
	// [Or]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [And]
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//         [EQ]
	//             [Path]
	//                 [Idn] e
	//             [Path]
	//                 [Idn] f
	// ( .a == .b | .c == .d ) & .e == .f
	// [And]
	//     [Or]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Path]
	//                 [Idn] b
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//     [EQ]
	//         [Path]
	//             [Idn] e
	//         [Path]
	//             [Idn] f
	// .a == .b | ( .c == .d & .e == .f )
	// [Or]
	//     [EQ]
	//         [Path]
	//             [Idn] a
	//         [Path]
	//             [Idn] b
	//     [And]
	//         [EQ]
	//             [Path]
	//                 [Idn] c
	//             [Path]
	//                 [Idn] d
	//         [EQ]
	//             [Path]
	//                 [Idn] e
	//             [Path]
	//                 [Idn] f
	// 0 != 1
	// [NE]
	//     [Num] 0
	//     [Num] 1
	// 1 >= 2
	// [GTE]
	//     [Num] 1
	//     [Num] 2
	// 3 > 4
	// [GT]
	//     [Num] 3
	//     [Num] 4
	// 5 <= 6
	// [LTE]
	//     [Num] 5
	//     [Num] 6
	// 7 < 8
	// [LT]
	//     [Num] 7
	//     [Num] 8
	// "ab" eq "AB"
	// [EQI]
	//     [Str] "ab"
	//     [Str] "AB"
	// "ab" ne "AB"
	// [NEI]
	//     [Str] "ab"
	//     [Str] "AB"
	// .a.[.a==.b].b.[.a==.b]
	// [Path]
	//     [Idn] a
	//     [Array]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Path]
	//                 [Idn] b
	//     [Idn] b
	//     [Array]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Path]
	//                 [Idn] b
	// .a.[ .a == .b ].b.[ .a == "b" ]
	// [Path]
	//     [Idn] a
	//     [Array]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Path]
	//                 [Idn] b
	//     [Idn] b
	//     [Array]
	//         [EQ]
	//             [Path]
	//                 [Idn] a
	//             [Str] "b"
	// .(upper)
	// [Path]
	//     [Fun] upper
	// .(upper).(lower).[ .(upper) == "HI" ]
	// [Path]
	//     [Fun] upper
	//     [Fun] lower
	//     [Array]
	//         [EQ]
	//             [Path]
	//                 [Fun] upper
	//             [Str] "HI"
	// .(exists)
	// [Path]
	//     [Fun] exists
}

func lines(v string) (lines []string) {
	for s := bufio.NewScanner(strings.NewReader(v)); s.Scan(); {
		if line := strings.TrimSpace(s.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	return
}
