package graph

type graph struct {
	edges    map[string][]string
	vertices map[string]string
}

func New() *graph {
	return &graph{
		edges:    make(map[string][]string),
		vertices: make(map[string]string),
	}
}

func (g *graph) AddEdge(from, to string) {
	g.edges[from] = append(g.edges[from], to)
}

func (g *graph) AddVectex(v, color string) {
	g.vertices[v] = color
}
