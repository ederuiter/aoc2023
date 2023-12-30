package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Plane uint8

const (
	H Plane = iota
	V
)

type Node struct {
	Y           int
	X           int
	HeatLoss    int
	Connections []*Connection
	Entry       Plane
}

type Connection struct {
	HeatLoss    int
	Source      *Node
	Destination *Node
	NumMoves    int
	Plane       Plane
}

type PriorityQueue[T any, P *T] struct {
	items      []P
	priorities map[P]int
	sorted     bool
}

func (p *PriorityQueue[T, P]) Add(item P, priority int) {
	p.items = append(p.items, item)
	p.priorities[item] = priority
	p.sorted = false
}

func (p *PriorityQueue[T, P]) UpdatePriority(item P, priority int) {
	if _, ok := p.priorities[item]; ok {
		if !p.sorted {
			p.priorities[item] = priority
			return
		}

		oldIndex := slices.Index(p.items, item)
		p.priorities[item] = priority
		if oldIndex == 0 {
			return
		}
		var i = 0
		for i = oldIndex - 1; i >= 0; i-- {
			if priority >= p.priorities[p.items[i]] {
				break
			} else {
				p.items[i+1] = p.items[i]
			}
		}
		p.items[i+1] = item
	}

}

func (p *PriorityQueue[T, P]) Remove() (P, int) {
	if len(p.items) == 0 {
		return nil, math.MaxInt
	}
	if !p.sorted {
		slices.SortFunc(p.items, func(a P, b P) int {
			return p.priorities[a] - p.priorities[b]
		})
		p.sorted = true
	}

	var res P
	res, p.items = p.items[0], p.items[1:]
	priority := p.priorities[res]
	delete(p.priorities, res)
	return res, priority
}

func (p *PriorityQueue[T, P]) HasItems() bool {
	return len(p.items) > 0
}

func (p *PriorityQueue[T, P]) HasItem(item P) bool {
	_, ok := p.priorities[item]
	return ok
}

func NewPriorityQueue[T any]() *PriorityQueue[T, *T] {
	return &PriorityQueue[T, *T]{
		items:      []*T{},
		priorities: map[*T]int{},
		sorted:     false,
	}
}

func CalcMinHeatLoss(nodes []*Node, source *Node, destination *Node) int {
	q := NewPriorityQueue[Node]()
	dist := map[*Node]int{}
	prev := map[*Node]*Connection{}
	for _, node := range nodes {
		dist[node] = math.MaxInt
		q.Add(node, math.MaxInt)
	}
	dist[source] = 0
	q.UpdatePriority(source, 0)
	for q.HasItems() {
		u, _ := q.Remove()
		if _, hasPrev := prev[u]; u == destination && hasPrev {
			res := 0
			for {
				conn := prev[u]
				fmt.Printf("=> %d, %d => %d\n", conn.Source.Y, conn.Source.X, conn.HeatLoss)
				res += conn.HeatLoss
				u = conn.Source
				if _, ok := prev[u]; !ok {
					return res
				}
			}
		}
		for _, connection := range u.Connections {
			v := connection.Destination
			if q.HasItem(v) {
				alt := dist[u] + connection.HeatLoss
				if alt < dist[v] {
					dist[v] = alt
					q.UpdatePriority(v, alt)
					prev[v] = connection
				}
			}
		}
	}
	return math.MaxInt
}

func newNode(y int, x int, p Plane, heatLoss int) *Node {
	return &Node{
		Y:           y,
		X:           x,
		Entry:       p,
		HeatLoss:    heatLoss,
		Connections: []*Connection{},
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	grid := map[Plane][][]*Node{
		V: {},
		H: {},
	}
	y := 0
	cols := 0
	scanner := bufio.NewScanner(file)
	nodes := []*Node{}
	for scanner.Scan() {
		grid[V] = append(grid[V], []*Node{})
		grid[H] = append(grid[H], []*Node{})
		line := strings.Split(scanner.Text(), "")
		for x, point := range line {
			p, _ := strconv.ParseInt(point, 10, 64)
			nodeH := newNode(y, x, H, int(p))
			nodeV := newNode(y, x, V, int(p))
			grid[H][y] = append(grid[H][y], nodeH)
			grid[V][y] = append(grid[V][y], nodeV)
			nodes = append(nodes, nodeV, nodeH)
		}
		cols = len(line)
		y++
	}
	rows := y

	for plane, g := range grid {
		newPlane := V
		if plane == V {
			newPlane = H
		}
		for y, row := range g {
			for x, node := range row {
				dx := 0
				dy := 0
				heatLoss := 0
				for i := 1; i <= 3; i++ {
					if plane == V {
						dx -= 1
					} else {
						dy -= 1
					}
					if y+dy >= 0 && y+dy < rows && x+dx >= 0 && x+dx < cols {
						heatLoss += grid[newPlane][y+dy][x+dx].HeatLoss
						conn := &Connection{
							Source:      node,
							Destination: grid[newPlane][y+dy][x+dx],
							Plane:       newPlane,
							NumMoves:    i,
							HeatLoss:    heatLoss,
						}
						node.Connections = append(
							node.Connections,
							conn,
						)
					} else {
						break
					}
				}
				dx = 0
				dy = 0
				heatLoss = 0
				for i := 1; i <= 3; i++ {
					if plane == V {
						dx += 1
					} else {
						dy += 1
					}
					if y+dy >= 0 && y+dy < rows && x+dx >= 0 && x+dx < cols {
						heatLoss += grid[newPlane][y+dy][x+dx].HeatLoss
						conn := &Connection{
							Source:      node,
							Destination: grid[newPlane][y+dy][x+dx],
							Plane:       newPlane,
							NumMoves:    i,
							HeatLoss:    heatLoss,
						}
						node.Connections = append(
							node.Connections,
							conn,
						)
					}
				}
			}
		}
	}

	minHeatLoss := min(
		CalcMinHeatLoss(nodes, grid[H][0][0], grid[H][rows-1][cols-1]),
		CalcMinHeatLoss(nodes, grid[H][0][0], grid[V][rows-1][cols-1]),
		CalcMinHeatLoss(nodes, grid[V][0][0], grid[H][rows-1][cols-1]),
		CalcMinHeatLoss(nodes, grid[V][0][0], grid[V][rows-1][cols-1]),
	)

	fmt.Printf("min heatloss: %d\n", minHeatLoss)
}
