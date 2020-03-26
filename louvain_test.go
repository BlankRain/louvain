package louvain

import (
	"testing"
)

func TestModularity(t *testing.T) {

	graphReader := NewGraphReader()
	graph := graphReader.Load("resource/karate.txt")

	louvain := NewLouvain(graph)
	louvain.Compute()

	actual := louvain.BestModularity()
	expected := WeightType(0.41880345)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

}
