package search

import (
	"context"
	"sort"
	"testing"
	"time"

	"notes/internal/store/storetest"
)

var latencyQueries = []string{"kafka", "bucket policy", "trigram", "deploy pipeline", "sqlite", "react vite", "certificate", "stream topic broker", "golang", "calendar sidebar"}

// TestSearchLatencyP95 checks the server-side p95 on 1,000 generated articles.
func TestSearchLatencyP95(t *testing.T) {
	if testing.Short() {
		t.Skip("latency check")
	}
	r := openBuild(t, storetest.Result(storetest.Articles(3, 1000)))
	ctx := context.Background()
	var durations []time.Duration
	for round := range 10 {
		for _, q := range latencyQueries {
			start := time.Now()
			if _, err := Search(ctx, r.DB(), q, 0); err != nil {
				t.Fatal(err)
			}
			if round > 0 { // first round warms the page cache
				durations = append(durations, time.Since(start))
			}
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p95 := durations[len(durations)*95/100]
	t.Logf("p95 = %s over %d queries (max %s)", p95, len(durations), durations[len(durations)-1])
	if p95 > 100*time.Millisecond {
		t.Fatalf("p95 = %s, want under 100ms", p95)
	}
}

func BenchmarkSearch(b *testing.B) {
	r := openBuildB(b)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if _, err := Search(ctx, r.DB(), latencyQueries[i%len(latencyQueries)], 0); err != nil {
			b.Fatal(err)
		}
	}
}
