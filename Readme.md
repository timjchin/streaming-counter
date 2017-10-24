# Top N Streaming Counter 
It's often desirable to get top-N counts to answer questions like "What are the most frequent User IDs in this stream of application logs?". In most situations, it's not practical to keep all the counts of ids seen in memory. This leads to solutions like Count-Min Sketch, which is a probailistic data structure that can approximate actual counts, with tuneable accuracy and probability that comes at the expense of memory and time.

The issue with using a count-min sketch is the API. To get the count of an item, you must know what the item is. When dealing with a stream of data, if all you have is the count-min sketch loaded with values, it's impossible to know what values to use to get the counts are. That's where this module comes in. 

`StreamingCounter` combines a sorted set with a count-min sketch. 

When an item is added to the set: 
1) The item is added to the Count-Min Sketch. O(n)
2) The approximate count is queried from the Count-Min Sketch. O(n)
3) The item (`string`) and count (`int`) are added to the sorted set O(log(n))
4) If the number of items in the sorted set are greater than the number of results to store, all of the lowest nodes in the sorted set are removed. O(log(n))

## API
### Configuration
```
type StreamingCounterConfig struct {
	// Number of items to store, must be greater than 0
	NumResults int
	// Accuracy of the counts within a factor of Epsilon.
	Epsilon float64
	// Probability that the counts are within a factor of Epsilon.
	Delta float64
}
```

### Methods
`StreamingCounter.Add(string)`: Add a value to the counter.

`StreamingCounter.GetAll() []Item`: Get a descending list of the Top Counts (configured through `StreamingCounterConfig.NumResults`)

### Benchmarks
```
BenchmarkStreamingCounter-4       500000              2983 ns/op             146 B/op          6 allocs/op
```