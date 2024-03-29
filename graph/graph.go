package graph

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"reflect"

	"github.com/Kistor/info_bez/crypto"
)

const (
	redColor = iota
	blueColor
	yellowColor
)

var colors = map[string]int{
	"R": redColor,
	"B": blueColor,
	"Y": yellowColor,
}

var (
	errInvalidColor      = errors.New("invalid vertex color")
	errVertexUncolored   = errors.New("vertex doesn't have any color")
	errColoringNotProper = errors.New("graph coloring isn't proper")
)

type Graph struct {
	Edges    map[string][]string
	vertices map[string]*vertexInfo
}

type vertexInfo struct {
	color string
	r     *big.Int
	P     *big.Int
	Q     *big.Int
	N     *big.Int
	c     *big.Int
	d     *big.Int
}

func New() *Graph {
	return &Graph{
		Edges:    make(map[string][]string),
		vertices: make(map[string]*vertexInfo),
	}
}

func (g *Graph) AddEdge(from, to string) {
	g.Edges[from] = append(g.Edges[from], to)
}

func (g *Graph) AddVectex(v, color string) error {
	if _, ok := colors[color]; !ok { // try to add vertex with non-existing color
		return errInvalidColor
	}
	g.vertices[v] = &vertexInfo{
		color: color,
	}
	return nil
}

func (g *Graph) ShuffleColors() {
	shuffled := reflect.ValueOf(colors).MapKeys()    // get all keys from color map
	mathrand.Shuffle(len(shuffled), func(i, j int) { // shuffle them
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	for k, v := range g.vertices {
		c := colors[v.color]
		g.vertices[k].color = shuffled[c].String()
	}
}

func (g *Graph) CalcVertexParams() error {
	for k, v := range g.vertices {
		r, err := rand.Int(rand.Reader, new(big.Int).
			Exp(big.NewInt(10), big.NewInt(20), nil))
		if err != nil {
			return err
		}
		bits := getReplacementBits(v.color) // get color bits
		for i, bit := range bits {
			r.SetBit(r, i, bit)
		}
		g.vertices[k].r = r
		g.vertices[k].P, err = rand.Prime(rand.Reader, 64)
		if err != nil {
			return err
		}
		g.vertices[k].Q, err = rand.Prime(rand.Reader, 64)
		if err != nil {
			return err
		}
		g.vertices[k].N = new(big.Int).Mul(g.vertices[k].P, g.vertices[k].Q)
		f := new(big.Int).Mul(
			new(big.Int).Add(g.vertices[k].P, big.NewInt(-1)),
			new(big.Int).Add(g.vertices[k].Q, big.NewInt(-1))) // f = (p-1)(q-1)
		g.vertices[k].d, err = rand.Int(rand.Reader, new(big.Int).
			Exp(big.NewInt(10), big.NewInt(20), nil))
		if err != nil {
			return err
		}
		euc, err := crypto.BigEuclidean(g.vertices[k].d, f) // calc c RSA param
		if err != nil {
			return err
		}
		for euc.Gcd.Cmp(big.NewInt(1)) != 0 {
			g.vertices[k].d, err = rand.Int(rand.Reader, new(big.Int).
				Exp(big.NewInt(10), big.NewInt(20), nil))
			if err != nil {
				return err
			}
			euc, err = crypto.BigEuclidean(g.vertices[k].d, f) // calc c RSA param
			if err != nil {
				return err
			}
		}
		if euc.X.Sign() < 0 {
			euc.X.Add(euc.X, f)
		}
		g.vertices[k].c = euc.X
	}
	return nil
}

type publicGraphInfo struct {
	Z *big.Int
	N *big.Int
	D *big.Int
}

func (g *Graph) SendPublicData() map[string]publicGraphInfo {
	data := make(map[string]publicGraphInfo, len(g.vertices))
	for k, v := range g.vertices {
		data[k] = publicGraphInfo{
			Z: new(big.Int).Exp(v.r, v.d, v.N),
			N: v.N,
			D: v.d,
		}
	}
	return data
}

func (g *Graph) C(v string) *big.Int {
	return g.vertices[v].c
}

func getReplacementBits(color string) []uint {
	c := uint(colors[color])                                     // get color number
	bitsCount := int(math.Ceil(math.Log2(float64(len(colors))))) // get bits count
	bits := make([]uint, bitsCount)                              // bits array
	for i := range bits {                                        // setting bits
		bits[i] = c & uint(math.Exp2(float64(i))) >> i
	}
	return bits
}

func (g *Graph) Bfs() error {
	/* just bfs alg to further ensure graph coloring is proper */
	var (
		queue   = []string{getRandVertex(g.Edges)}
		visited = make(map[string]struct{})
	)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if _, ok := visited[cur]; ok {
			continue
		}
		visited[cur] = struct{}{}
		curInfo, ok := g.vertices[cur]
		if !ok {
			return fmt.Errorf("%w: vertex %s", errVertexUncolored, cur)
		}
		for _, v := range g.Edges[cur] {
			if _, ok := visited[v]; ok {
				continue
			}
			info, ok := g.vertices[v]
			if !ok {
				return fmt.Errorf("%w: vertex %s", errVertexUncolored, v)
			}
			if info.color == curInfo.color {
				return errColoringNotProper
			}
			queue = append(queue, v)
		}
	}
	return nil
}

func getRandVertex(edges map[string][]string) string {
	verteces := reflect.ValueOf(edges).MapKeys()
	idx := mathrand.Intn(len(edges))
	return verteces[idx].String()
}

func (g *Graph) GetRandEdge() (string, string) {
	v1 := getRandVertex(g.Edges)
	rndIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.Edges[v1]))))
	if err != nil {
		panic("got rand error")
	}
	// rndIdx := mathrand.Intn(len(g.Edges[v1]))
	v2 := g.Edges[v1][rndIdx.Int64()]
	return v1, v2
}
