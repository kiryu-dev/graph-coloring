package graph

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	mathrand "math/rand"
	"reflect"

	crypto "github.com/Kistor/info_bez/crypto"
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
	errInvalidColor = errors.New("invalid vertex color")
)

type graph struct {
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

func New() *graph {
	return &graph{
		Edges:    make(map[string][]string),
		vertices: make(map[string]*vertexInfo),
	}
}

func (g *graph) AddEdge(from, to string) {
	g.Edges[from] = append(g.Edges[from], to)
}

func (g *graph) AddVectex(v, color string) error {
	if _, ok := colors[color]; !ok {
		return errInvalidColor
	}
	g.vertices[v] = &vertexInfo{
		color: color,
	}
	return nil
}

func (g *graph) ShuffleColors() {
	shuffled := reflect.ValueOf(colors).MapKeys()
	mathrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	for k, v := range g.vertices {
		c := colors[v.color]
		g.vertices[k].color = shuffled[c].String()
	}
}

func (g *graph) CalcVertexParams() error {
	for k, v := range g.vertices {
		r, err := rand.Int(rand.Reader, new(big.Int).
			Exp(big.NewInt(10), big.NewInt(256), nil))
		if err != nil {
			return err
		}
		bits := getReplacementBits(v.color)
		for i, bit := range bits {
			r.SetBit(r, i, bit)
		}
		g.vertices[k].r = r
		g.vertices[k].P, err = rand.Prime(rand.Reader, 1024)
		if err != nil {
			return err
		}
		g.vertices[k].Q, err = rand.Prime(rand.Reader, 1024)
		if err != nil {
			return err
		}
		g.vertices[k].N = new(big.Int).Mul(g.vertices[k].P, g.vertices[k].Q)
		f := new(big.Int).Mul(
			new(big.Int).Add(g.vertices[k].P, big.NewInt(-1)),
			new(big.Int).Add(g.vertices[k].Q, big.NewInt(-1)))
		g.vertices[k].d, err = rand.Int(rand.Reader, new(big.Int).
			Exp(big.NewInt(10), big.NewInt(256), nil))
		if err != nil {
			return err
		}
		euc, err := crypto.BigEuclidean(g.vertices[k].d, f)
		if err != nil {
			return err
		}
		for euc.Gcd.Cmp(big.NewInt(1)) != 0 {
			g.vertices[k].d, err = rand.Int(rand.Reader, new(big.Int).
				Exp(big.NewInt(10), big.NewInt(256), nil))
			if err != nil {
				return err
			}
			euc, err = crypto.BigEuclidean(g.vertices[k].d, f)
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

func (g *graph) SendPublicData() map[string]publicGraphInfo {
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

func (g *graph) C(v string) *big.Int {
	return g.vertices[v].c
}

func getReplacementBits(color string) []uint {
	c := uint(colors[color])
	bitsCount := int(math.Ceil(math.Log2(float64(len(colors)))))
	bits := make([]uint, bitsCount)
	for i := range bits {
		bits[i] = c & uint(math.Exp2(float64(i))) >> i
	}
	return bits
}
