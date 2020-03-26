package louvain

type Community struct {
	inWeight    WeightType
	totalWeight WeightType
}

type Level struct {
	graph         Graph
	communities   []Community
	inCommunities []int
}

type Louvain struct {
	level   []Level
	current *Level
	m2      WeightType
}

func NewLouvain(graph Graph) *Louvain {
	louvain := Louvain{}
	louvain.level = make([]Level, 1)
	louvain.current = &louvain.level[0]
	louvain.current.graph = graph
	louvain.current.communities = make([]Community, louvain.current.graph.GetNodeSize())
	louvain.current.inCommunities = make([]int, louvain.current.graph.GetNodeSize())
	louvain.m2 = WeightType(0.0)
	for nodeId := 0; nodeId < louvain.current.graph.GetNodeSize(); nodeId++ {
		louvain.current.inCommunities[nodeId] = nodeId
		neigh := WeightType(0.0)
		for _, edge := range louvain.current.graph.GetIncidentEdges(nodeId) {
			neigh += edge.weight
		}
		self := WeightType(louvain.current.graph.GetSelfWeight(nodeId))
		louvain.current.communities[nodeId] = Community{self, neigh + self}
		louvain.m2 += (neigh + 2*self)
	}
	return &louvain
}

func (this *Louvain) BestModularity() WeightType {
	return this.Modularity(len(this.level) - 1)
}

func (this *Louvain) Modularity(level int) WeightType {
	return this.level[level].modularity(this.m2)
}

func (this *Level) modularity(m2 WeightType) WeightType {
	q := WeightType(0.0)
	for _, comm := range this.communities {
		q += comm.inWeight/m2 - (comm.totalWeight/m2)*(comm.totalWeight/m2)
	}
	return q
}

func (this *Louvain) merge() bool {
	improved := false

	q := make([]int, this.current.graph.GetNodeSize())
	mark := make([]bool, this.current.graph.GetNodeSize())
	for nodeId := 0; nodeId < this.current.graph.GetNodeSize(); nodeId++ {
		q[nodeId] = nodeId
	}

	for len(q) != 0 {
		nodeId := q[0]
		q = q[1:] // pop_front
		mark[nodeId] = true

		neighWeights := map[int]WeightType{}

		self := WeightType(this.current.graph.GetSelfWeight(nodeId))
		totalWeight := self

		neighWeightsKeys := make([]int, 0, len(neighWeights))
		for _, edge := range this.current.graph.GetIncidentEdges(nodeId) {
			destCommId := this.current.inCommunities[edge.destId]
			if _, exists := neighWeights[destCommId]; !exists {
				neighWeightsKeys = append(neighWeightsKeys, destCommId)
			}
			neighWeights[destCommId] += edge.weight
			totalWeight += edge.weight
		}

		prevCommunity := this.current.inCommunities[nodeId]
		prevNeighWeight := WeightType(neighWeights[prevCommunity])
		this.remove(nodeId, prevCommunity, 2*prevNeighWeight+self, totalWeight)

		maxInc := WeightType(0.0)
		bestCommunity := prevCommunity
		bestNeighWeight := WeightType(prevNeighWeight)
		for _, community := range neighWeightsKeys {
			weight := neighWeights[community]
			inc := WeightType(weight - this.current.communities[community].totalWeight*totalWeight/this.m2)
			if inc > maxInc {
				maxInc = inc
				bestCommunity = community
				bestNeighWeight = weight
			}
		}

		this.insert(nodeId, bestCommunity, 2*bestNeighWeight+self, totalWeight)

		if bestCommunity != prevCommunity {
			improved = true
			for _, edge := range this.current.graph.GetIncidentEdges(nodeId) {
				if mark[edge.destId] {
					q = append(q, edge.destId)
					mark[edge.destId] = false
				}
			}
		}
	}

	return improved
}

func (this *Louvain) insert(nodeId int, community int, inWeight WeightType, totalWeight WeightType) {
	this.current.inCommunities[nodeId] = community
	this.current.communities[community].inWeight += inWeight
	this.current.communities[community].totalWeight += totalWeight
}

func (this *Louvain) remove(nodeId int, community int, inWeight WeightType, totalWeight WeightType) {
	this.current.inCommunities[nodeId] = -1
	this.current.communities[community].inWeight -= inWeight
	this.current.communities[community].totalWeight -= totalWeight
}

func (this *Louvain) rebuild() {

	renumbers := map[int]int{}
	num := 0

	for nodeId, inCommunity := range this.current.inCommunities {
		if commId, exists := renumbers[inCommunity]; !exists {
			renumbers[inCommunity] = num
			this.current.inCommunities[nodeId] = num
			num++
		} else {
			this.current.inCommunities[nodeId] = commId
		}
	}

	newCommunities := make([]Community, num)
	for nodeId := 0; nodeId < len(this.current.communities); nodeId++ {
		if comm, exists := renumbers[nodeId]; exists {
			newCommunities[comm] = this.current.communities[nodeId]
		}
	}

	communityNodes := make([][]int, len(newCommunities))
	for nodeId := 0; nodeId < this.current.graph.GetNodeSize(); nodeId++ {
		communityNodes[this.current.inCommunities[nodeId]] = append(communityNodes[this.current.inCommunities[nodeId]], nodeId)
	}

	newGraph := Graph{make(Edges, len(newCommunities)), make([]WeightType, len(newCommunities))}

	for commId := 0; commId < newGraph.GetNodeSize(); commId++ {
		newEdges := map[int]WeightType{}
		selfWeight := WeightType(0.0)
		for _, nodeId := range communityNodes[commId] {
			edges := this.current.graph.GetIncidentEdges(nodeId)
			for _, edge := range edges {
				newEdges[this.current.inCommunities[edge.destId]] += edge.weight
			}
			selfWeight += this.current.graph.GetSelfWeight(nodeId)
		}

		newGraph.AddSelfEdge(commId, newEdges[commId]+selfWeight)
		for nodeId, weight := range newEdges {
			if nodeId != commId {
				newGraph.AddDirectedEdge(commId, nodeId, weight)
			}
		}
	}

	newInCommunities := make([]int, newGraph.GetNodeSize())
	for nodeId := 0; nodeId < len(newInCommunities); nodeId++ {
		newInCommunities[nodeId] = nodeId
	}

	this.level = append(this.level, Level{newGraph, newCommunities, newInCommunities})
	this.current = &this.level[len(this.level)-1]
}

func (this *Louvain) GetLevel(n int) *Level {
	return &this.level[n]
}

func (this *Louvain) Compute() {
	for this.merge() {
		this.rebuild()
	}
}

func (this *Louvain) GetPertition(level int) []int {
	nodeToCommunity := make([]int, this.level[0].graph.GetNodeSize())
	for nodeId, _ := range nodeToCommunity {
		commId := nodeId
		for l := 0; l != level; l++ {
			commId = this.level[l].inCommunities[commId]
		}
		nodeToCommunity[nodeId] = commId
	}
	return nodeToCommunity
}

func (this *Louvain) GetBestPertition() []int {
	return this.GetPertition(len(this.level))
}

func (this *Louvain) GetNodeToCommunityInEachLevel(nodeId int) []int {
	nodeToCommunityInEachLevel := make([]int, len(this.level))
	commId := nodeId
	for l := 0; l != len(this.level); l++ {
		commId = this.level[l].inCommunities[commId]
		nodeToCommunityInEachLevel[l] = commId
	}
	return nodeToCommunityInEachLevel
}
