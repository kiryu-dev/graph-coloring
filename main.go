package main

import (
	"bufio"
	"flag"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/kiryu-dev/graph-coloring/graph"
)

const maxVertexCount = 1001

func main() {
	filename := flag.String("f", "", "graph info")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
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
	g := graph.New()
	for i := 0; i < m; i++ {
		scanner.Scan()
		edge := strings.Split(scanner.Text(), ",")
		if len(edge) != 2 {
			log.Fatal("invalid edge format")
		}
		g.AddEdge(edge[0], edge[1])
	}
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
	g.ShuffleColors()
	if err := g.CalcVertexParams(); err != nil {
		log.Fatal(err)
	}
	graphData := g.SendPublicData()
	for k, v := range g.Edges {
		rand.Shuffle(len(v), func(i, j int) {
			v[i], v[j] = v[j], v[i]
		})
		Z1 := new(big.Int).Exp(graphData[k].Z, g.C(k), graphData[k].N)
		for i := range v {
			Z2 := new(big.Int).Exp(graphData[v[i]].Z, g.C(v[i]), graphData[v[i]].N)
			if Z1.Bit(0) == Z2.Bit(0) && Z1.Bit(1) == Z2.Bit(1) {
				log.Fatal("alice is caught CHEATING!")
			}
		}
	}
	log.Println("alice proved she knows the correct graph coloring")
}
