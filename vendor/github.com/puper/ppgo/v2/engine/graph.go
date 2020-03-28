package engine

import "errors"

type graph struct {
	// names contains the keys of the "edges" field.
	// It allows the vertices to be sorted.
	// It makes the structure deterministic.
	names []string
	// vertices ordered by name.
	vertices map[string]*graphVertex
}

// graphVertex contains the vertex data.
type graphVertex struct {
	// numIn in the number of incoming edges.
	numIn int
	// numInTmp is used by the TopologicalOrdering to avoid messing with numIn
	numInTmp int
	// out contains the name the outgoing edges.
	out []string
	// outMap is the same as "out", but in a map
	// to quickly check if a vertex is in the outgoing edges.
	outMap map[string]struct{}
}

// newGraph creates a new graph.
func newGraph() *graph {
	return &graph{
		names:    []string{},
		vertices: map[string]*graphVertex{},
	}
}

// AddVertex adds a vertex to the graph.
func (g *graph) AddVertex(v string) {
	_, ok := g.vertices[v]
	if ok {
		return
	}

	g.names = append(g.names, v)

	g.vertices[v] = &graphVertex{
		numIn:  0,
		out:    []string{},
		outMap: map[string]struct{}{},
	}
}

// AddEdge adds an edge to the graph.
func (g *graph) AddEdge(from, to string) {
	g.AddVertex(from)
	g.AddVertex(to)

	// check if the edge is aleady registered
	if _, ok := g.vertices[from].outMap[to]; ok {
		return
	}

	// update the vertices
	g.vertices[from].out = append(g.vertices[from].out, to)
	g.vertices[from].outMap[to] = struct{}{}
	g.vertices[to].numIn++
}

// TopologicalOrdering returns a valid topological sort.
// It implements Kahn's algorithm.
// If there is a cycle in the graph, an error is returned.
// The list of vertices is also returned even if it is not ordered.
func (g *graph) TopologicalOrdering() ([]string, error) {
	l := []string{}
	q := []string{}

	for _, v := range g.names {
		if g.vertices[v].numIn == 0 {
			q = append(q, v)
		}
		g.vertices[v].numInTmp = g.vertices[v].numIn
	}

	for len(q) > 0 {
		n := q[len(q)-1]
		q = q[:len(q)-1]
		l = append(l, n)

		for _, m := range g.vertices[n].out {
			g.vertices[m].numInTmp--
			if g.vertices[m].numInTmp == 0 {
				q = append(q, m)
			}
		}
	}

	if len(l) != len(g.names) {
		return append([]string{}, g.names...), errors.New("a cycle has been found in the dependencies")
	}

	return l, nil
}
