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

// runes
var whitespace = []int32(" \t\n\r\v\f")
var asciiLowercase = []int32("abcdefghijklmnopqrstuvwxyz")
var asciiUppercase = []int32("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var asciiLetters = append(asciiLowercase, asciiUppercase...)
var digits = []int32("0123456789")
var hexdigits = append(digits, []int32("abcdefABCDEF")...)
var octdigits = []int32("01234567")
var punctuation = []int32("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
var printable = arraysJoin(whitespace, digits, asciiLetters, punctuation)

// intervals
var punctuationExtended = []int32{161, 191}
var latinExtended = []int32{192, 591}
var ipaAlphabet = []int32{592, 687}
var spaceModifiers = []int32{688, 767}
var diacriticalMarks = []int32{768, 879}
var cyrillic = []int32{1024, 1327}
var thaana = []int32{1920, 1969}
var devanagari = []int32{2304, 2431}
var myanmar = []int32{4096, 4255}
var hangulJamo = []int32{4352, 4607}
var canadian = []int32{5120, 5759}
var runic = []int32{5792, 5880}
var vedic = []int32{7376, 7417}
var phonetic = []int32{7424, 7615}
var currency = []int32{8352, 8383}
var letterLike = []int32{8448, 8527}
var number = []int32{8528, 8587}
var arrows = []int32{8592, 8703}
var mathematical = []int32{8704, 8959}
var technical = []int32{8960, 9215}
var enclosed = []int32{9312, 9471}
var miscellaneous = []int32{9472, 10239}
var braille = []int32{10240, 10495}
var supplemental = []int32{10496, 11007}
var kangxi = []int32{12032, 12245}
var hangul = []int32{12593, 12686}
var cjkExtended = []int32{13312, 19893}
var cjk = []int32{19968, 40934}
var egyptian = []int32{77824, 78863}
var alchemical = []int32{128768, 128883}

// all
var _other = arraysJoin(
	punctuationExtended, latinExtended, ipaAlphabet, spaceModifiers,
	diacriticalMarks, cyrillic, runic, vedic,
	phonetic, currency, number, arrows, mathematical, technical, egyptian, alchemical,
)
var allChars = append(printable, runeSet(_other, true)...)
var allCharsNotNL = excludingRune(allChars, '\n')

// name capture
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
		result = string(allCharsNotNL[rand.Intn(len(allCharsNotNL))])
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
	if max == -1 {
		max = MaxRepeat
	}
	if captureName != "" &&
		registry[captureName] != nil &&
		(r.Sub[0].Op == syntax.OpAnyChar || r.Sub[0].Op == syntax.OpAnyCharNotNL) {
		return fixSize(min, max, registry[captureName])
	}
	times := max
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

func arraysJoin(args ...[]int32) []int32 {
	result := []int32{}
	for _, arr := range args {
		result = append(result, arr...)
	}
	return result
}

func excludingRune(s []int32, a rune) []int32 {
	n := make([]int32, len(s))
	copy(n, s)

	index := -1
	for i, v := range n {
		if v == a {
			index = i
			break
		}
	}
	if index == -1 {
		return n
	}

	n[index] = n[len(n)-1]
	return n[:len(n)-1]
}

func excludingRunes(s []int32, a []rune) []int32 { // O(s+a)
	n := make([]int32, len(s))
	copy(n, s)

	ri := make(map[rune]int)
	for _, r := range a {
		ri[r] = -1
	}

	counter := len(ri)
	for i, r := range n {
		if _, ok := ri[r]; ok {
			ri[r] = i
			counter--
		}
		if counter == 0 {
			break
		}
	}

	for k, v := range ri {
		if v == -1 {
			delete(ri, k)
		}
	}

	counter = 0
	for _, i := range ri {
		counter++
		r := n[len(n)-counter]
		n[i] = r
		if _, ok := ri[r]; ok {
			ri[r] = i
		}
		if counter == len(ri) {
			break
		}
	}

	return n[:len(n)-counter]
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
	id := strconv.Itoa(int(i))
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
