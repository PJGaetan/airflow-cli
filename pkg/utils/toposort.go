package utils

import (
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type Node struct {
	Id             int64
	Name           string
	DownStreamName []string
}

func (n Node) ID() int64 { return n.Id }

// func (n Node) negated() Node { return Node{-n.Id, n.Name} }

// topoSort return a list of node sorted.
func TopoSort(nodes map[string]Node) (sorted []string) {
	g := simple.NewDirectedGraph()

	for _, element := range nodes {
		for _, downStreamTask := range element.DownStreamName {
			e := simple.Edge{
				F: element,
				T: nodes[downStreamTask],
			}
			g.SetEdge(e)
		}
	}
	sortedNodes, err := topo.SortStabilized(g, nil)
	ExitIfError(err)

	nodeMap := make(map[int64]string)
	for _, node := range nodes {
		nodeMap[node.Id] = node.Name
	}

	for _, n := range sortedNodes {
		sorted = append(sorted, nodeMap[n.ID()])
	}

	return sorted
}
