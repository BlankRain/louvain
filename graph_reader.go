package louvain

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type GraphReader struct {
	nodeLabelToIndex map[string]int
	nodeIndexToLabel []string
}

func NewGraphReader() *GraphReader {
	return &GraphReader{map[string]int{}, make([]string, 0, 1000)}
}

func (this *GraphReader) addNode(nodeLabel string) {
	if _, exists := this.nodeLabelToIndex[nodeLabel]; !exists {
		nodeNum := len(this.nodeLabelToIndex)
		this.nodeLabelToIndex[nodeLabel] = nodeNum
		this.nodeIndexToLabel = append(this.nodeIndexToLabel, nodeLabel)
	}
}

func (this *GraphReader) GetNodeIndex(nodeLabel string) int {
	return this.nodeLabelToIndex[nodeLabel]
}

func (this *GraphReader) GetNodeLabel(nodeIndex int) string {
	return this.nodeIndexToLabel[nodeIndex]
}

func (this *GraphReader) GetNodeSize() int {
	return len(this.nodeLabelToIndex)
}
func (this *GraphReader) AddNode(node string) {
	this.addNode(node)
}
func (this *GraphReader) Load(filename string) Graph {
	edgeNum := 0

	fp, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		this.addNode(line[0])
		this.addNode(line[1])
		edgeNum++
	}

	graph := Graph{make(Edges, edgeNum), make([]WeightType, this.GetNodeSize())}

	fp.Seek(0, 0)
	scanner = bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		source := this.GetNodeIndex(line[0])
		dest := this.GetNodeIndex(line[1])
		weight := WeightType(1.0)
		graph.AddUndirectedEdge(source, dest, weight)
	}

	return graph

}
