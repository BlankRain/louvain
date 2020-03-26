package louvain

type WeightType float32

type Edge struct {
	destId int
	weight WeightType
}

type Edges [][]Edge

type Graph struct {
	Incidences Edges
	Selfs      []WeightType
}

func (this *Graph) AddUndirectedEdge(sourceId int, destId int, weight WeightType) {
	this.Incidences[sourceId] = append(this.Incidences[sourceId], Edge{destId, weight})
	this.Incidences[destId] = append(this.Incidences[destId], Edge{sourceId, weight})
}

func (this *Graph) AddDirectedEdge(sourceId int, destId int, weight WeightType) {
	this.Incidences[sourceId] = append(this.Incidences[sourceId], Edge{destId, weight})
}

func (this *Graph) AddSelfEdge(nodeId int, weight WeightType) {
	this.Selfs[nodeId] = weight
}

func (this *Graph) GetIncidentEdges(nodeId int) []Edge {
	return this.Incidences[nodeId]
}

func (this *Graph) GetSelfWeight(nodeId int) WeightType {
	return this.Selfs[nodeId]
}

func (this *Graph) GetNodeSize() int {
	return len(this.Selfs)
}
