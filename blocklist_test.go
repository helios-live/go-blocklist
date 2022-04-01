package blocklist_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.ideatocode.tech/blocklist"
	"gopkg.in/yaml.v3"
)

var withOutput = flag.Bool("withOutput", false, "When set to true some tests will turn on output")

func TestYaml(t *testing.T) {
	type bl struct {
		Blocklist *blocklist.BL
	}
	ym := `blocklist:
  - facebook.com
  - google.com
  - '*yahoo.com'
`
	x := bl{}
	err := yaml.Unmarshal([]byte(ym), &x)
	assert.NoError(t, err)

	buf, err := yaml.Marshal(&x)
	assert.NoError(t, err)
	fmt.Fprintf(os.Stderr, string(buf))
	assert.Equal(t, ym, string(buf))
}
func TestConcurrency(t *testing.T) {

	assert.NotPanics(t, func() {
		x := blocklist.New()

		x.Add("google.com")
		x.Add("facebook.com")
		x.Add("yahoo.com")
		x.Add("*nike.com")
		x.Add("*go.com")

		wg := sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				time.Sleep(500 * time.Millisecond)
				for i := 0; i < 10000; i++ {
					x.Match("google.com")
					x.Match("facebook.com")
					x.Match("yahoo.com")
					x.Match("tikenike.com")
					x.Match("gogo.com")
					x.Match("asdf.com")
					x.Match("miko.com")
					x.Match("aaaa.es")
					x.Match("oka.de")
				}
				wg.Done()
			}()
		}
		wg.Wait()
		log.Println("Done")
	}, "The code panicked")

}
