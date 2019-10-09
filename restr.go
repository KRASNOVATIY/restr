package restr

import (
	"fmt"
	"math/rand"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
	"time"
)

// MaxRepeat - limit for * and + literals
var MaxRepeat = 100

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
var allChars = append(printable, runeSet(_other, true)...)

var registry = make(map[string]func() string)
var captureName = ""

// Rstr - construct random string by regexpr string
func Rstr(re string) string {
	three, err := syntax.Parse(re, syntax.PerlX)
	if err != nil {
		panic(err)
	}
	return handleState(three)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleState(r *syntax.Regexp) string {
	// https://golang.org/pkg/regexp/syntax/#Op
	// r.Op, r.Rune, r.Min, r.Max, r.Name

	result := ""
	switch r.Op {
	case syntax.OpLiteral:
		result = string(r.Rune)
	case syntax.OpAnyChar:
		result = string(allChars[rand.Intn(len(allChars))])
	case syntax.OpAnyCharNotNL:
		newAllChars := removeRune(allChars, '\n')
		result = string(newAllChars[rand.Intn(len(newAllChars))])
	case syntax.OpCharClass:
		charClass := runeSet(r.Rune, true)
		result = string(charClass[rand.Intn(len(charClass))])
	case syntax.OpCapture:
		result = handleCapture(r)
	case syntax.OpAlternate:
		result = handleState(r.Sub[rand.Intn(len(r.Sub))])
	case syntax.OpConcat:
		result = handleConcat(r)
	case syntax.OpQuest:
		result = handleRepeat(0, 1, r)
	case syntax.OpRepeat:
		result = handleRepeat(r.Min, r.Max, r)
	case syntax.OpPlus:
		result = handleRepeat(1, MaxRepeat, r)
	case syntax.OpStar:
		result = handleRepeat(0, MaxRepeat, r)
	default:
	}
	return result
}

func handleCapture(r *syntax.Regexp) string {
	if r.Name != "" {
		oldName := captureName
		captureName = r.Name
		defer func() { captureName = oldName }()
	}
	return handleConcat(r)
}

func handleConcat(r *syntax.Regexp) string {

	newStr := []string{}
	for _, s := range r.Sub {
		newStr = append(newStr, handleState(s))
	}
	return strings.Join(newStr, "")
}

func handleRepeat(min, max int, r *syntax.Regexp) string {
	if captureName != "" &&
		registry[captureName] != nil &&
		(r.Sub[0].Op == syntax.OpAnyChar || r.Sub[0].Op == syntax.OpAnyCharNotNL) {
		return fixSize(min, max, registry[captureName])
	}
	times := max
	if max == -1 {
		max = MaxRepeat
	}
	if max-min > 0 {
		times = rand.Intn(max-min+1) + min // rand.Intn(1) = 0 Always
	}
	var result []string
	for i := 0; i < times; i++ {
		result = append(result, handleState(r.Sub[0]))
	}
	return strings.Join(result, "")
}

func fixSize(min, max int, fn func() string) string {
	word := fn()
	for len(word) < min {
		word += fn()
	}
	for len(word) > max {
		word = word[:len(word)-1]
	}
	if min == 0 {
		return RandomString([]string{word, ""})()
	}
	return word
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

func removeRune(s []int32, a rune) []int32 {
	index := 0
	for i, v := range allChars {
		if v == a {
			index = i
			break
		}
	}
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

// named capture

// RandomString - returns a random string from a sequence "array"
func RandomString(array []string) func() string {
	return func() string {
		return array[rand.Intn(len(array))]
	}
}

// RegisterName - register data source for Named Capture Group (?P<name>.{10})
// handles AnyChar and OpAnyCharNotNL literal "." with repeat literals "? * + {}"
// name: Gapture Group Name
// fn: generator func, which returns a string
func RegisterName(name string, fn func() string) {
	registry[name] = fn
}

// test

func explain(r *syntax.Regexp, i uint) {
	if r == nil {
		return
	}
	id := fmt.Sprintf("%s", strconv.Itoa(int(i)))
	if r.Name != "" {
		id = fmt.Sprintf("%s %s", strconv.Itoa(int(i)), r.Name)
	}
	fmt.Printf("%s%s) %v%v min=%d, max=%d\n",
		strings.Repeat("\t", int(i)), id, r.Op, r.Rune, r.Min, r.Max)
	i++
	for _, sub := range r.Sub {
		explain(sub, i)
	}
}

func about(re string) {
	three, _ := syntax.Parse(re, syntax.PerlX)
	explain(three, uint(0))
}

func all(a []bool) bool {
	for _, i := range a {
		if i != true {
			return false
		}
	}
	return true
}

func test() bool {
	RegisterName("xname", RandomString([]string{"xa", "xb"}))

	mg := NewMarkovGen(3, []rune{' ', ','})
	mg.ModelApply("title1", "Tiny love is my favorite toy, I love it and cannot live without it", 1)
	mg.ModelApply("title2", "I am your best friend, I love sweets", 1)
	RegisterName("yname", mg.Generate(25))

	tests := []string{
		`(?P<xname>\w{5}(?P<t>\+.{5}\+).{5})\_(?P<xname>.+ (S|s)(?P<f>\d)) \d{2,5} \:(?P<yname>.*\+.?\+.+ @R)`,
		`(?P<xname>\w{5})`,
		`(?P<xname>(Dd|rr|\dU|.\.){5} .{4} .{3,7})`,
		`^[\da-zA-Z\-](\&|\*|\^) \d+ ..(?P<xname>\d{2})$`,
		`\d{4,8}`,
		`\w{10}`,
		`[AbC6]`,
		`.+`,
		`(Tik|Tak|Tok)`,
		`\A(?P<help>\[\d\])[\w]+\s\d*\(\w\)(Xor)?>\09{2,5}(\x55|[^a0-9]\*).{2}[a-f][[:space:]]Uu$`,
		`S?`,
		`S?7`,
		`F+`,
	}
	var results []bool
	for _, re := range tests {
		match, _ := regexp.MatchString(re, Rstr(re))
		results = append(results, match)
	}
	return all(results)
}
