package counter

import (
	"math/rand"
	"testing"
	"time"

	"github.com/timjchin/logcounter/random"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func BenchmarkStreamingCounter(b *testing.B) {
	text := random.NewWeightedChoices(random.RandomChoices(50000, 10))
	counter, _ := NewStreamingCounter(&StreamingCounterConfig{
		NumResults: 15,
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		currString := text.Get()
		b.StartTimer()
		counter.Add(currString)
	}
}
