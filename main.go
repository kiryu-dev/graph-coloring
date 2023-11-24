package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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
		g.AddVectex(v[0], v[1])
	}
	fmt.Println(g)
}
