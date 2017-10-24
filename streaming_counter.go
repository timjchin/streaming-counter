package counter

import (
	"fmt"

	"github.com/tylertreat/BoomFilters"
	"github.com/wangjia184/sortedset"
)

type StreamingCounterConfig struct {
	// number of items to store
	NumResults int
	// Accuracy of the counts within a factor of Epsilon.
	Epsilon float64
	// Probability that the counts are within a factor of Epsilon.
	Delta float64
}

// StreamingCounters combine a Count-Min-Sketch and a Sorted Set to maintain the
// a configurable number of highest seen values, while minimizing the memory
// required to keep and store this state.
//
// This is useful for times where absolute accuracy is not required, and you'd like
// to understand both the keys and approximate counts for a large amount of data.
type StreamingCounter struct {
	config *StreamingCounterConfig
	counts *boom.CountMinSketch
	set    *sortedset.SortedSet
}

func NewStreamingCounter(config *StreamingCounterConfig) (*StreamingCounter, error) {
	if config.Epsilon == float64(0) {
		config.Epsilon = 0.001
	}
	if config.Delta == float64(0) {
		config.Delta = 0.99
	}
	if config.NumResults == 0 {
		return nil, fmt.Errorf("Unknown amount of results to store.")
	}
	return &StreamingCounter{
		config: config,
		counts: boom.NewCountMinSketch(config.Epsilon, config.Delta),
		set:    sortedset.New(),
	}, nil
}

type Item struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// Add a single value to the counter
func (p *StreamingCounter) Add(val string) {
	b := []byte(val)
	p.counts.Add(b)

	currCount := p.counts.Count(b)
	p.set.AddOrUpdate(val, sortedset.SCORE(currCount), nil)

	if p.set.GetCount() > p.config.NumResults {
		p.set.GetByRankRange(p.config.NumResults, -1, true)
	}
}

// Return all known counts, in descending order.
func (p *StreamingCounter) GetAll() []Item {
	nodes := p.set.GetByRankRange(-1, 0, false)
	out := make([]Item, len(nodes))
	for i, node := range nodes {
		out[i] = Item{
			Key:   node.Key(),
			Value: int64(node.Score()),
		}
	}
	return out
}
