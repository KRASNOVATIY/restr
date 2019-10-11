package restr

import (
	"strings"
)

// MarkovGen - Markov model text generator
type MarkovGen interface {
	ModelApply(title, text string, norma uint)
	Generate(len uint) func() string
}

type markovGen struct {
	depth   uint
	texts   map[string]string
	exclude []rune
	model   []string
	_seq    []string
}

// NewMarkovGen - constructor MarkovGen
// depth >= 2: determines the size of the model element.
//     The more, the higher the accuracy of the reproduction of the original phrases of the text
// exclude: array of runes that should be excluded from the model
func NewMarkovGen(depth uint, exclude []rune) MarkovGen {
	if depth <= 1 {
		panic("the size of the element (state of the text model) must exceed 1")
	}
	mg := new(markovGen)
	mg.depth = depth
	mg.exclude = exclude
	mg.texts = make(map[string]string) // TODO: Add model normalization
	mg._seq = make([]string, 0)
	return mg
}

// ModelApply - apply model for MarkovGen
// text: text model contents
// norma: text multiplier
func (mg *markovGen) ModelApply(title, text string, norma uint) {
	for _, r := range mg.exclude {
		text = strings.Replace(text, string(r), "", -1)
	}
	if len(text) < int(mg.depth) {
		panic("the size of the text after the exclusion of characters must exceed the size of the element")
	}
	text = strings.Repeat(text, int(norma))
	mg.texts[title] = text

	ta := strings.Split(text, "")
	for i := range ta {
		j := i + int(mg.depth)
		if j > len(text) {
			continue
		}
		item := ta[i:j]
		mg.model = append(mg.model, strings.Join(item, ""))
	}
	for len(mg.model) < 1 {
		mg.ModelApply(title, text, norma+1)
	}
}

func (mg *markovGen) next(prefix string, slip int) string {
	mg._seq = []string{}
	for _, s := range mg.model {
		if strings.HasPrefix(s[slip:], prefix) {
			mg._seq = append(mg._seq, s[mg.depth-1:])
		}
	}
	for len(mg._seq) < 1 {
		mg.next(prefix[1:], slip+1)
	}
	return RandomString(mg._seq)()
}

func (mg *markovGen) generate(length uint) string {
	if len(mg.model) == 0 {
		panic("model is empty")
	}
	start := RandomString(mg.model)()[:mg.depth-1]
	nxt := mg.next(start, 0)
	word := start + nxt

	for i := 0; i < int(length-mg.depth); i++ {
		nxt = mg.next(word[i+1:], 0)
		word += nxt
	}

	return word
}

// Generate - generate text based on a model with a length of "length"
func (mg *markovGen) Generate(length uint) func() string {
	return func() string {
		return mg.generate(length)
	}
}
