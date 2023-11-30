package main

import (
	"bufio"
	"flag"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/kiryu-dev/graph-coloring/graph"
)

const (
	maxVertexCount = 1001
	a              = 1
)

const (
	zeroKnowledgeMode = "zero-knowledge"
	bfsMode           = "bfs"
)

func main() {
	filename := flag.String("f", "", "graph info")
	mode := flag.String("m", "zero-knowledge", "proof mode")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	/* read vertex and edge count from file */
	scanner.Scan()
	graphInfo := strings.Split(scanner.Text(), ",")
	if len(graphInfo) != 2 {
		log.Fatal("invalid graph info size")
	}
	n, err := strconv.Atoi(graphInfo[0])
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(graphInfo[1])
	if err != nil {
		log.Fatal(err)
	}
	if !(n < maxVertexCount && m <= n*n) {
		log.Fatalf("invalid graph info: 'n' must be less than %d and 'm' must be less or equal 'n^2'", maxVertexCount)
	}
	/* read graph edges info from file */
	g := graph.New()
	for i := 0; i < m; i++ {
		scanner.Scan()
		edge := strings.Split(scanner.Text(), ",")
		if len(edge) != 2 {
			log.Fatal("invalid edge format")
		}
		g.AddEdge(edge[0], edge[1])
	}
	/* read vertex color info from file */
	for i := 0; i < n; i++ {
		scanner.Scan()
		v := strings.Split(scanner.Text(), ",")
		if len(v) != 2 {
			log.Fatal("invalid vertex format")
		}
		if err := g.AddVectex(v[0], v[1]); err != nil {
			log.Fatal(err)
		}
	}
	switch strings.ToLower(*mode) {
	case zeroKnowledgeMode: // main mode. zero-knowledge proof
		zeroKnowledgeProof(g, a, m)
	case bfsMode: // extra mode. just to further ensure graph coloring is proper
		if err := g.Bfs(); err != nil {
			log.Fatal(err)
		}
		log.Println("graph coloring is proper!")
	default:
		log.Fatal("invalid method to proof proper graph coloring")
	}
}

func zeroKnowledgeProof(g *graph.Graph, a, e int) {
	if a < 1 {
		panic("NOOOOO a MUST BE GREATER THAN ZEROOOOO")
	}
	for i := 0; i < a*e; i++ { // iterate a|E| times
		g.ShuffleColors()                            // shuffle colors
		if err := g.CalcVertexParams(); err != nil { // for every vertex calc r r param & RSA params (P, Q, N, c, d) & set color bits
			log.Fatal(err)
		}
		graphData := g.SendPublicData() // calc Z = r^d mod N & send N, d, Z for every vertex
		v1, v2 := g.GetRandEdge()       // choose rand edge
		// fmt.Printf("%d: %s %s\n", i, v1, v2)
		Z1 := new(big.Int).Exp(graphData[v1].Z, g.C(v1), graphData[v1].N)
		Z2 := new(big.Int).Exp(graphData[v2].Z, g.C(v2), graphData[v2].N)
		/* color bits of Z1 & Z2 must be different */
		if Z1.Bit(0) == Z2.Bit(0) && Z1.Bit(1) == Z2.Bit(1) {
			log.Fatal("alice is caught CHEATING!")
		}
		// for k, v := range g.Edges {
		// 	rand.Shuffle(len(v), func(i, j int) {
		// 		v[i], v[j] = v[j], v[i]
		// 	})
		// 	Z1 := new(big.Int).Exp(graphData[k].Z, g.C(k), graphData[k].N)
		// 	for i := range v {
		// 		Z2 := new(big.Int).Exp(graphData[v[i]].Z, g.C(v[i]), graphData[v[i]].N)
		// 		if Z1.Bit(0) == Z2.Bit(0) && Z1.Bit(1) == Z2.Bit(1) {
		// 			log.Fatal("alice is caught CHEATING!")
		// 		}
		// 	}
		// }
	}
	log.Println("alice proved she knows proper graph coloring")
}
