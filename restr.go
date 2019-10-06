package restr

import (
	"fmt"
	"math/rand"
	"regexp/syntax"
	"strings"
	"time"
)

const (
	// MAXREPEAT - limit for *, + literals
	MAXREPEAT = 100
)

var whitespace = []int32(" \t\n\r\v\f")
var asciiLowercase = []int32("abcdefghijklmnopqrstuvwxyz")
var asciiUppercase = []int32("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var asciiLetters = append(asciiLowercase, asciiUppercase...)
var digits = []int32("0123456789")
var hexdigits = append(digits, []int32("abcdefABCDEF")...)
var octdigits = []int32("01234567")
var punctuation = []int32("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
var printable = append(digits, append(asciiLetters, append(whitespace, punctuation...)...)...)
var _other = []int32{161, 895, 913, 1327, 1329, 1366, 1488, 1514}
var all = append(printable, runeSet(_other, true)...)

// Rstr - construct random string by regexpr string
func Rstr(re string) string {
	three, err := syntax.Parse(re, syntax.PerlX)
	if err != nil {
		panic(err)
	}
	return buildString(three)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func buildString(r *syntax.Regexp) string {
	newStr := []string{}
	for _, s := range r.Sub {
		newStr = append(newStr, handleState(s))
	}
	return strings.Join(newStr, "")
}

func intRange(start, stop int32) []int32 {
	var result []int32
	for i := start; i <= stop; i++ {
		result = append(result, i)
	}
	return result
}

func runeSet(set []int32, ranged bool) []int32 {
	if !ranged {
		return set
	}
	var result []int32
	if len(set) > 0 {
		var start int32
		for i, v := range set {
			if i&1 == 0 {
				start = v
			} else {
				result = append(result, intRange(start, v)...)
			}
		}
	}
	return result
}

func handleRepeat(min, max int, r *syntax.Regexp) string {
	var result []string
	times := max
	if max == -1 {
		max = MAXREPEAT
	}
	if max-min > 0 {
		times = rand.Intn(max-min+1) + min // rand.Intn(1) = 0 Always
	}
	for i := 0; i < times; i++ {
		result = append(result, handleState(r.Sub[0]))
	}
	return strings.Join(result, "")
}

func removeRune(s []int32, a int32) []int32 {
	index := 0
	for i, v := range all {
		if v == a {
			index = i
			break
		}
	}
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func handleState(r *syntax.Regexp) string {
	// https://golang.org/pkg/regexp/syntax/#Op
	// r.Op, r.Rune, r.Min, r.Max, r.Name

	result := ""
	switch r.Op {
	case syntax.OpLiteral:
		result = string(r.Rune)
	case syntax.OpAnyChar:
		result = string(all[rand.Intn(len(all))])
	case syntax.OpAnyCharNotNL:
		newAll := removeRune(all, '\n')
		result = string(newAll[rand.Intn(len(newAll))])
	case syntax.OpCharClass:
		charClass := runeSet(r.Rune, true)
		result = string(charClass[rand.Intn(len(charClass))])
	case syntax.OpCapture:
		result = buildString(r)
	case syntax.OpAlternate:
		result = handleState(r.Sub[rand.Intn(len(r.Sub))])
	case syntax.OpConcat:
		result = buildString(r)
	case syntax.OpQuest:
		result = handleRepeat(0, 1, r)
	case syntax.OpRepeat:
		result = handleRepeat(r.Min, r.Max, r)
	case syntax.OpPlus:
		result = handleRepeat(1, MAXREPEAT, r)
	case syntax.OpStar:
		result = handleRepeat(0, MAXREPEAT, r)
	default:
	}
	return result
}

func explain(r *syntax.Regexp, i uint) {
	if r == nil {
		return
	}
	fmt.Printf("%d) sub=%v, rune=%v, min=%d, max=%d, name=%s\n",
		i, r.Op, r.Rune, r.Min, r.Max, r.Name)
	i++
	for _, sub := range r.Sub {
		explain(sub, i)
	}
}

func about(re string) {
	three, _ := syntax.Parse(re, syntax.PerlX)
	explain(three, uint(0))
}
