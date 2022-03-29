package blocklist

import (
	"flag"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var withOutput = flag.Bool("withOutput", false, "When set to true some tests will turn on output")

func TestConcurrency(t *testing.T) {

	assert.NotPanics(t, func() {
		x := New()

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
