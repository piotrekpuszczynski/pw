//Piotr Puszczyński 
//myślę,że dla standardowych danych (utawionych w kodzie) najlepiej widać jak działa algorytm
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var announce = make(chan bool)

func main() {
	n := flag.Int("n", 10, "vertices")
	d := flag.Int("d", 4, "shortcuts")
	h := flag.Int("h", 5, "hosts")
	flag.Parse()

	if *n < 2 || *d < 0 || *d > maxShortcuts(*n) || *h < 0 {
		fmt.Println("Invalid parameters")
		return
	}

	vertices := initializeNodes(*n, *h)
	initializeShortcuts(vertices, *d)
	initializeRoutingTable(vertices)
	initializeHosts(*h, vertices)
	drawGraph(vertices)
	printCosts(vertices)
	goRoutines(vertices)

	for i := 0; i < *h; i++ {
		<- announce
	}
}

func initializeHosts(h int, r []*Vertex) {
	for i := 0; i < h; i++ {
		random := rand.Intn(len(r))
		r[random].hosts = append(r[random].hosts, &Host{random, len(r[random].hosts), make(chan *SPacket)})
	}
}

func printCosts(vertices []*Vertex) {
	var s strings.Builder
	s.WriteString("table of costs changed: \n")
	s.WriteString("   ")
	for i := range vertices {
		s.WriteString(strconv.Itoa(i))
		s.WriteString(" ")
	}
	s.WriteString("\n")
	for _, vertex := range vertices {
		s.WriteString(strconv.Itoa(vertex.id) + " [")
		for i, c := range vertex.routingTable.cost {
			s.WriteString(strconv.Itoa(c))
			if i < len(vertex.routingTable.cost) - 1 {
				s.WriteString(" ")
			}
		}
		s.WriteString("]\n")
	}
	fmt.Println(s.String())
}

func goRoutines(vertices []*Vertex) {
	for _, vertex := range vertices {
		go receiver(vertex, vertices)
	}
	for _, vertex := range vertices {
		go sender(vertex, vertices)
	}
	for _, vertex := range vertices {
		go forwarder(vertex, vertices)
	}

	for _, vertex := range vertices {
		for _, host := range vertex.hosts {
			go hostSender(host, vertices)
		}
	}
}

func forwarder(vertex *Vertex, vertices []*Vertex) {
	for {

		packet := <- vertex.packetChan

		packet.visited = append(packet.visited, vertex)

		if packet.destinationR == vertex.id {
			vertex.hosts[packet.destinationH].packet <- packet
			fmt.Println("packet sent from router " + strconv.Itoa(vertex.id) + " to host " + strconv.Itoa(packet.destinationH))
		} else {
			vertices[vertex.routingTable.nextHop[packet.destinationR].id].packetChan <- packet
			fmt.Println("packet sent from router " + strconv.Itoa(vertex.id) + " to router " +
				strconv.Itoa(vertex.routingTable.nextHop[packet.destinationR].id))
		}
	}
}

func hostSender(host *Host, vertices []*Vertex) {

	var rR int
	var rH int
	for {
		rR = rand.Intn(len(vertices))

		for len(vertices[rR].hosts) == 0 {
			rR = rand.Intn(len(vertices))
		}

		rH = rand.Intn(len(vertices[rR].hosts))

		if !(rR == host.r && rH == host.h) {
			break
		}
	}

	packet := &SPacket{host.r, host.h, rR, rH, make([]*Vertex, 0)}

	vertices[host.r].packetChan <- packet
	fmt.Println("packet sent from host " + strconv.Itoa(host.h) + " to router " + strconv.Itoa(host.r))

	for {
		received := <- host.packet
		fmt.Println("packet received by host " + strconv.Itoa(received.destinationH) +
			" from router " + strconv.Itoa(received.destinationR))

		if received.destinationH == received.sourceH && received.destinationR == received.sourceR {
			var s strings.Builder
			s.WriteString("routers visited by packet: [")
			for i, r := range received.visited {
				s.WriteString(strconv.Itoa(r.id))
				if i < len(received.visited) - 1 {
					s.WriteString(" ")
				}
			}
			s.WriteString("]")
			fmt.Println(s.String())
			announce <- true

		} else {
			rand.Seed(time.Now().UnixNano())
			r := rand.Float64()  * 10 + 5
			time.Sleep(time.Second * time.Duration(r))

			received.destinationH = received.sourceH
			received.destinationR = received.sourceR

			vertices[host.r].packetChan <- received
		}
	}
}

func sender(vertex *Vertex, vertices []*Vertex) {
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Float64() * 8 + 2
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
					vertexes[vertex.connections[j].id].packet <- p
				}(vertex, j, vertices)
			}
		}
	}
}

func receiver(vertex *Vertex, vertices []*Vertex) {
	for {
		p := <- vertex.packet

		for i := 0; i < len(p.cost); i++ {
			newCost := 1 + p.cost[i]

			if newCost < vertices[p.j[i]].routingTable.cost[vertex.id] {
				vertices[p.j[i]].routingTable.cost[vertex.id] = newCost
				vertices[p.j[i]].routingTable.nextHop[vertex.id] = vertices[p.vertex]
				vertices[p.j[i]].routingTable.changed[vertex.id] = true
				printCosts(vertices)
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

type Host struct {
	r 	   int
	h      int
	packet chan *SPacket
}

type SPacket struct {
	sourceR 	 int
	sourceH 	 int
	destinationR int
	destinationH int
	visited      []*Vertex
}

type Vertex struct {
	id           int
	connections  []*Vertex
	routingTable *R
	packet       chan *Packet
	hosts		 []*Host
	packetChan   chan *SPacket
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

func initializeRoutingTable(vertices []*Vertex) {
	for _, vertex := range vertices {
		for i := 0; i < len(vertices); i++ {

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
				vertex.routingTable.nextHop[i] = vertices[i]
			} else {
				vertex.routingTable.cost[i] = abs(vertex.id - i)
				if vertex.id < i {
					vertex.routingTable.nextHop[i] = vertices[vertex.id + 1]
				} else {
					vertex.routingTable.nextHop[i] = vertices[vertex.id - 1]
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

func drawGraph(vertices []*Vertex) {
	fmt.Println("GRAPH: ")
	for _, vertex := range vertices {
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
	var s strings.Builder
	for _, vertex := range vertices {
		s.WriteString(strconv.Itoa(len(vertex.hosts)) + "   ")
	}
	fmt.Println(s.String())
}

func initializeShortcuts(vertices []*Vertex, d int) {
	for i := 0; i < len(vertices) - 1; i++ {
		vertices[i].connections = append(vertices[i].connections, vertices[i + 1])
	}
	for i := 1; i < len(vertices); i++ {
		vertices[i].connections = append(vertices[i].connections, vertices[i - 1])
	}
	for i := 0; i < d; i++ {
		rand.Seed(time.Now().UnixNano())
		r1 := rand.Intn(len(vertices))
		rand.Seed(time.Now().UnixNano())
		r2 := rand.Intn(len(vertices))

		correctRandoms := false
		if !(r1 == r2) {
			for _, connection := range vertices[r1].connections {
				if connection.id == r2 {
					correctRandoms = true
					break
				}
			}
		} else {
			correctRandoms = true
		}

		if !correctRandoms {
			vertices[r1].connections = append(vertices[r1].connections, vertices[r2])
			vertices[r2].connections = append(vertices[r2].connections, vertices[r1])
		} else {
			i--
		}
	}
}

func initializeNodes(n int, h int) []*Vertex {
	vertexes := make([]*Vertex, n)

	for i := 0; i < n; i++ {
		vertexes[i] = &Vertex{i, make([]*Vertex, 0),
			&R{make([]*Vertex, n), make([]int, n), make([]bool, n)},
			make(chan *Packet), make([]*Host, 0), make(chan *SPacket, h)}
	}

	return vertexes
}

func maxShortcuts(n int) int {
	return ((n - 3) * n / 2) + 1
}