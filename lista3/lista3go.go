//Piotr Puszczyński 
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var announce = make(chan string)

func main() {
	n := flag.Int("n", 5, "vertexes")
	d := flag.Int("d", 2, "shortcuts")
	flag.Parse()

	if *n < 2 || *d < 0 || *d > maxShortcuts(*n) {
		fmt.Println("Invalid parameters")
		return
	}

	vertexes := initializeNodes(*n)
	initializeShortcuts(vertexes, *d)
	initializeRoutingTable(vertexes)
	drawGraph(vertexes)
	goReceiversAndSenders(vertexes)

	for {
		select {
		case s := <-announce:
			fmt.Println(s)
		case <-time.After(3 * time.Second):
			terminate(vertexes)
			os.Exit(0)
		}
	}
}

func terminate(vertexes []*Vertex) {
	fmt.Println()
	fmt.Println("table of costs: ")
	fmt.Printf("   ")
	for i := range vertexes {
		fmt.Printf(strconv.Itoa(i))
		fmt.Printf(" ")
	}
	fmt.Println()
	for _, vertex := range vertexes {
		fmt.Println(vertex.id, vertex.routingTable.cost)
	}
}

func goReceiversAndSenders(vertexes []*Vertex) {
	for _, vertex := range vertexes {
		go receiver(vertex, vertexes)
	}
	for _, vertex := range vertexes {
		go sender(vertex, vertexes)
	}
}

func sender(vertex *Vertex, vertexes []*Vertex) {
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Float64() * 2
		time.Sleep(time.Second * time.Duration(r))

		jTab := make([]int, 0)
		costTab := make([]int, 0)

		send := checkChanged(vertex)
		for i := 0; i < len(vertex.routingTable.cost); i++ {
			if vertex.routingTable.changed[i] {
				jTab = append(jTab, i)
				costTab = append(costTab, vertex.routingTable.cost[i])
				vertex.routingTable.changed[i] = false
			}
		}

		if send {
			for j := 0; j < len(vertex.connections); j++ {
				go func(vertex *Vertex, j int, vertexes []*Vertex) {
					p := &Packet{jTab, costTab, vertex.id}
					var s string
					//s = "waiting to send packet from " + strconv.Itoa(vertex.id) + " to " + strconv.Itoa(i)
					//announce <- s
					vertexes[vertex.connections[j].id].packet <- p
					s = "sent packet from " + strconv.Itoa(vertex.id) + " to " + strconv.Itoa(vertex.connections[j].id)
					announce <- s
				}(vertex, j, vertexes)
			}
		}
	}
}

func receiver(vertex *Vertex, vertexes []*Vertex) {
	for {
		p := <- vertex.packet
		s := strconv.Itoa(vertex.id) + " received packet from " + strconv.Itoa(p.vertex)
		announce <- s

		for i := 0; i < len(p.cost); i++ {
			newCost := 1 + p.cost[i]

			if newCost < vertexes[p.j[i]].routingTable.cost[vertex.id] {
				vertexes[p.j[i]].routingTable.cost[vertex.id] = newCost
				vertexes[p.j[i]].routingTable.nextHop[vertex.id] = vertexes[p.vertex]
				vertexes[p.j[i]].routingTable.changed[vertex.id] = true
			}
		}
	}
}

func checkChanged(vertex *Vertex) bool {
	for _, changed := range vertex.routingTable.changed {
		if changed {
			return true
		}
	}
	return false
}

type Vertex struct {
	id           int
	connections  []*Vertex
	routingTable *R
	packet       chan *Packet
}

type R struct {
	nextHop     []*Vertex
	cost        []int
	changed		[]bool
}

type Packet struct {
	j      []int
	cost   []int
	vertex int
}

func initializeRoutingTable(vertexes []*Vertex) {
	for _, vertex := range vertexes {
		for i := 0; i < len(vertexes); i++ {

			if vertex.id == i {
				vertex.routingTable.cost[i] = 0
				vertex.routingTable.nextHop[i] = nil
				vertex.routingTable.changed[i] = false
				continue
			}

			contains := false
			for _, connection := range vertex.connections {
				if connection.id == i {
					contains = true
					break
				}
			}

			if contains {
				vertex.routingTable.cost[i] = 1
				vertex.routingTable.nextHop[i] = vertexes[i]
			} else {
				vertex.routingTable.cost[i] = abs(vertex.id - i)
				if vertex.id < i {
					vertex.routingTable.nextHop[i] = vertexes[vertex.id + 1]
				} else {
					vertex.routingTable.nextHop[i] = vertexes[vertex.id - 1]
				}
			}
			vertex.routingTable.changed[i] = true
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func drawGraph(vertexes []*Vertex) {
	fmt.Println("GRAPH: ")
	for _, vertex := range vertexes {
		for _, connection := range vertex.connections {
			if vertex.id < connection.id {
				a := 0
				var s strings.Builder
				for a < vertex.id {
					s.WriteString("    ")
					a++
				}
				s.WriteString(strconv.Itoa(vertex.id))
				for a < connection.id - 1 {
					s.WriteString("----")
					a++
				}
				s.WriteString("---" + strconv.Itoa(connection.id))
				fmt.Println(s.String())
			}
		}
	}
}

func initializeShortcuts(vertexes []*Vertex, d int) {
	for i := 0; i < len(vertexes) - 1; i++ {
		vertexes[i].connections = append(vertexes[i].connections, vertexes[i + 1])
	}
	for i := 1; i < len(vertexes); i++ {
		vertexes[i].connections = append(vertexes[i].connections, vertexes[i - 1])
	}
	for i := 0; i < d; i++ {
		rand.Seed(time.Now().UnixNano())
		r1 := rand.Intn(len(vertexes))
		rand.Seed(time.Now().UnixNano())
		r2 := rand.Intn(len(vertexes))

		correctRandoms := false
		if !(r1 == r2) {
			for _, connection := range vertexes[r1].connections {
				if connection.id == r2 {
					correctRandoms = true
					break
				}
			}
		} else {
			correctRandoms = true
		}

		if !correctRandoms {
			vertexes[r1].connections = append(vertexes[r1].connections, vertexes[r2])
			vertexes[r2].connections = append(vertexes[r2].connections, vertexes[r1])
		} else {
			i--
		}
	}
}

func initializeNodes(n int) []*Vertex {
	vertexes := make([]*Vertex, n)

	for i := 0; i < n; i++ {
		vertexes[i] = &Vertex{i, make([]*Vertex, 0),
			&R{make([]*Vertex, n), make([]int, n), make([]bool, n)},
			make(chan *Packet)}
	}

	return vertexes
}

func maxShortcuts(n int) int {
	return ((n - 3) * n / 2) + 1
}