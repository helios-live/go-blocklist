package blocklist // import go.ideatocode.tech/blocklist

import (
	"encoding/json"

	"github.com/armon/go-radix"
)

// BL is a blocklist
type BL struct {
	fixed *radix.Tree
	wild  *radix.Tree
	list  []string
}

// New returns a new instancel of BL
func New() *BL {
	return &BL{
		fixed: radix.New(),
		wild:  radix.New(),
	}
}

// Add a pattern to the blocklist
func (b *BL) add(str string, updateList bool) {
	var x struct{}

	wild := false
	if str[0:1] == "*" {
		wild = true
		str = str[1:]
	}

	pattern := reverse(str)

	var updated = false
	if wild {
		_, updated = b.wild.Insert(pattern, x)
	} else {
		_, updated = b.fixed.Insert(pattern, x)
	}
	if updateList && !updated {
		b.list = append(b.list, str)
	}
}

// Add a pattern to the blocklist
func (b *BL) Add(str string) {
	b.add(str, true)
}

// Match checks if domain matches the patterns stored in BL
func (b *BL) Match(domain string) bool {
	domain = reverse(domain)

	if _, ok := b.fixed.Get(domain); ok {
		return true
	}
	if _, _, ok := b.wild.LongestPrefix(domain); ok {
		return true
	}
	return false
}
func reverse(s string) string {

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// UnmarshalJSON is called by the json package
func (b *BL) UnmarshalJSON(data []byte) error {
	b.fixed = radix.New()
	b.wild = radix.New()

	err := json.Unmarshal(data, &b.list)
	if err != nil {
		return err
	}
	for _, v := range b.list {
		b.add(v, false)
	}
	return nil
}

// MarshalJSON is called by the json package
func (b BL) MarshalJSON() ([]byte, error) {
	return json.Marshal(&b.list)
}
