package main

import (
	"fmt"
	"github.com/akolybelnikov/aoc-2015/internal/utils"
	"os"
	"strconv"
	"strings"
)

type graph [][]int
type node struct {
	name string
	idx  int
	adj  []node
}

func main() {
	data, err := os.ReadFile("inputs/day09.txt")
	if err != nil {
		fmt.Println("Failed to read input file.", err)
		os.Exit(0)
	}
	input := string(data)

	fmt.Println("--- Part One ---")
	fmt.Println("Result:", part1(input))

	fmt.Println("--- Part Two ---")
	fmt.Println("Result:", part2(input))

	os.Exit(0)
}

// part one
func part1(input string) int {
	lines, err := utils.ParseLines(input)
	utils.HandleErr(err)

	nodes, g, err := newGraph(lines)
	utils.HandleErr(err)

	minCost, shortestPath, _, _ := tsp(g)
	cities := make([]string, len(shortestPath))
	for _, idx := range shortestPath {
		cities[idx] = nodes[idx]
	}
	fmt.Println(strings.Join(cities, " -> "))

	return minCost
}

// part two
func part2(input string) int {
	lines, err := utils.ParseLines(input)
	utils.HandleErr(err)

	nodes, g, err := newGraph(lines)
	utils.HandleErr(err)

	_, _, maxCost, longestPath := tsp(g)
	cities := make([]string, len(longestPath))
	for _, idx := range longestPath {
		cities[idx] = nodes[idx]
	}
	fmt.Println(strings.Join(cities, " -> "))

	return maxCost
}

func newGraph(lines []string) (map[int]string, graph, error) {
	nodes := make(map[string]int)
	adj := make([][3]int, 0)
	cur := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		var idx1, idx2 int
		if a, ok := nodes[fields[0]]; !ok {
			idx1 = cur
			cur++
		} else {
			idx1 = a
		}
		if b, ok := nodes[fields[2]]; !ok {
			idx2 = cur
			cur++
		} else {
			idx2 = b
		}
		dist, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, nil, err
		}
		adj = append(adj, [3]int{idx1, idx2, dist})
		nodes[fields[0]] = idx1
		nodes[fields[2]] = idx2
	}

	g := make(graph, len(nodes))
	for _, a := range adj {
		idx1, idx2, dist := a[0], a[1], a[2]
		if g[idx1] == nil {
			g[idx1] = make([]int, len(nodes))
		}
		g[idx1][idx2] = dist
		if g[idx2] == nil {
			g[idx2] = make([]int, len(nodes))
		}
		g[idx2][idx1] = dist
	}

	return invertMap(nodes), g, nil
}

func invertMap(m map[string]int) map[int]string {
	res := make(map[int]string)
	for k, v := range m {
		res[v] = k
	}
	return res
}

func tsp(g [][]int) (int, []int, int, []int) {
	n := len(g)
	const Inf = int(1e9) // A large value to represent infinity

	// DP tables for shortest and longest paths
	minDp := make([][]int, 1<<n)
	maxDp := make([][]int, 1<<n)
	parent := make([][]int, 1<<n)
	parentMax := make([][]int, 1<<n)

	for i := 0; i < (1 << n); i++ {
		minDp[i] = make([]int, n)
		maxDp[i] = make([]int, n)
		parent[i] = make([]int, n)
		parentMax[i] = make([]int, n)
		for j := 0; j < n; j++ {
			minDp[i][j] = Inf  // Initialize shortest path DP table with infinity
			maxDp[i][j] = -Inf // Initialize longest path DP table with negative infinity
			parent[i][j] = -1  // Initialize parent arrays with -1 (undefined)
			parentMax[i][j] = -1
		}
	}

	// Initialize DP for each starting node
	for i := 0; i < n; i++ {
		minDp[1<<i][i] = 0 // Starting cost is 0 for shortest path
		maxDp[1<<i][i] = 0 // Starting cost is 0 for longest path
	}

	// Main DP loop
	for mask := 1; mask < (1 << n); mask++ { // Iterate over all subsets of nodes
		for u := 0; u < n; u++ { // Current ending node
			if mask&(1<<u) == 0 { // Skip if u is not in the current subset
				continue
			}
			for v := 0; v < n; v++ { // Try to transition to another node
				if mask&(1<<v) != 0 || g[u][v] == 0 { // Skip if v is already visited or no edge
					continue
				}

				nextMask := mask | (1 << v)

				// Shortest Path DP Update
				newCost := minDp[mask][u] + g[u][v]
				if newCost < minDp[nextMask][v] {
					minDp[nextMask][v] = newCost
					parent[nextMask][v] = u // Record the parent node for shortest path reconstruction
				}

				// Longest Path DP Update
				newCost = maxDp[mask][u] + g[u][v]
				if newCost > maxDp[nextMask][v] {
					maxDp[nextMask][v] = newCost
					parentMax[nextMask][v] = u // Record the parent node for longest path reconstruction
				}
			}
		}
	}

	// Identify the final nodes and costs
	endMask := (1 << n) - 1 // All nodes visited
	minCost, maxCost := Inf, -Inf
	lastNodeForMin, lastNodeForMax := -1, -1

	for u := 0; u < n; u++ {
		// Find the minimum cost path end node
		if minDp[endMask][u] < minCost {
			minCost = minDp[endMask][u]
			lastNodeForMin = u
		}
		// Find the maximum cost path end node
		if maxDp[endMask][u] > maxCost {
			maxCost = maxDp[endMask][u]
			lastNodeForMax = u
		}
	}

	// Ensure valid paths exist for both shortest and longest paths
	if lastNodeForMin == -1 || lastNodeForMax == -1 {
		panic("Unable to reconstruct paths: check graph connectivity.")
	}

	// Reconstruct the shortest path directly
	curNode := lastNodeForMin
	curMask := endMask
	shortestPath := []int{}

	for curNode != -1 {
		shortestPath = append([]int{curNode}, shortestPath...) // Add current node to the front of the path
		nextNode := parent[curMask][curNode]
		curMask ^= 1 << curNode // Remove current node from the mask
		curNode = nextNode      // Move to the parent node
	}

	// Reconstruct the longest path directly
	curNode = lastNodeForMax
	curMask = endMask
	longestPath := []int{}

	for curNode != -1 {
		longestPath = append([]int{curNode}, longestPath...) // Add current node to the front of the path
		nextNode := parentMax[curMask][curNode]
		curMask ^= 1 << curNode // Remove current node from the mask
		curNode = nextNode      // Move to the parent node
	}

	return minCost, shortestPath, maxCost, longestPath
}
