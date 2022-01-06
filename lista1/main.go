package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	id          int
	packet      chan Packet
	connections []*Node
	serviced	[]*Packet
	trap		chan bool
}

type Packet struct {
	id       int
	visited  []*Node
	lifetime int
}

func (node *Node) append(n *Node) {
	node.connections = append(node.connections, n)
}

func (node *Node) receive(announce chan *Node) {
	if len(node.packet) == 1 {
		temp := <- node.packet
		temp.lifetime++
		fmt.Printf("%v received packet with id %v.", node.id, temp.id)
		node.serviced = append(node.serviced, &temp)
		temp.visited = append(temp.visited, node)
		//bool timeout := make(chan bool)

		select {
			case <- node.trap:
				fmt.Printf("packet with id %v has been caught by trap in node %v.", temp.id, node.id)
				i := 0
				for i < len(p) {
					if p[i].id == temp.id {
						p[i].visited = temp.visited
						break
					}
					i++
				}
				announce <- node
				break
			default:
				if temp.lifetime >= lifetime {
					fmt.Printf("packet with id %v reached lifetime in node %v.", temp.id, node.id)
					i := 0
					for i < len(p) {
						if p[i].id == temp.id {
							p[i].visited = temp.visited
							break
						}
						i++
					}
					announce <- node
				} else {
					rand.Seed(time.Now().UnixNano())
					sec := rand.Intn(5)
					time.Sleep(time.Second * time.Duration(sec))

					node.send(temp)
				}
		}
	}
}

func (node *Node) send(packet Packet) {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(len(node.connections))

	//fmt.Printf("%v is waiting for free space in node %v", node.id, node.connections[random].id)
	node.connections[random].packet <- packet

	fmt.Printf("Packet %v sent to %v.", packet.id, node.connections[random].id)
}

func (node *Node) procedure(announce chan *Node) {
	for {
		rand.Seed(time.Now().UnixNano())
		sec := rand.Float64() * 5
		time.Sleep(time.Second * time.Duration(sec))
		node.receive(announce)
	}
}

func (node *Node) lastNodeListener(announce chan *Node) {
	for {
		if len(node.packet) == 1 {
			rand.Seed(time.Now().UnixNano())
			sec := rand.Float64() * 5
			time.Sleep(time.Second * time.Duration(sec))
			node.receiveLast(announce)
		}
	}
}

func (node *Node) receiveLast(announce chan *Node) {
	if len(node.packet) == 1 {
		temp := <- node.packet
		temp.lifetime++
		fmt.Printf("%v received packet with id %v.", node.id, temp.id)
		node.serviced = append(node.serviced, &temp)
		temp.visited = append(temp.visited, node)

		select {
		case <- node.trap:
			fmt.Printf("packet with id %v has been caught by trap in node %v.", temp.id, node.id)
			i := 0
			for i < len(p) {
				if p[i].id == temp.id {
					p[i].visited = temp.visited
					break
				}
				i++
			}
			announce <- node
			break
		default:
			if temp.lifetime >= lifetime {
				fmt.Printf("packet with id %v reached lifetime in node %v.", temp.id, node.id)
				i := 0
				for i < len(p) {
					if p[i].id == temp.id {
						p[i].visited = temp.visited
						break
					}
					i++
				}
				announce <- node
			} else {
				rand.Seed(time.Now().UnixNano())
				sec := rand.Intn(5)
				time.Sleep(time.Second * time.Duration(sec))

				fmt.Printf("Received %v pocket.", temp.id)
				i := 0
				for i < len(p) {
					if p[i].id == temp.id {
						p[i].visited = temp.visited
						break
					}
					i++
				}
				announce <- node
			}
		}
	}
}

func makeNode(id int, connections []*Node) *Node {
	node := &Node{id, make(chan Packet, 1), connections, make([]*Packet, 0), make(chan bool)}
	return node
}

func initializePackets(k int) []Packet {
	temp := make([]Packet, k)
	j := 0

	for j < k {
		temp[j] = Packet{j, make([]*Node, 0), 0}
		j++
	}
	return temp
}

func draw(i int,j int) {
	a := 0
	var s strings.Builder

	if i < j {
		for a < i {
			s.WriteString("    ")
			a++
		}
		s.WriteString(strconv.Itoa(i))
		for a < j - 1 {
			s.WriteString("----")
			a++
		}
		s.WriteString("-->" + strconv.Itoa(j))
	} else {
		for a < j {
			s.WriteString("    ")
			a++
		}
		s.WriteString(strconv.Itoa(j) + "<--")

		for a < i - 1 {
			s.WriteString("----")
			a++
		}
		s.WriteString(strconv.Itoa(i))
	}

	fmt.Println(s.String())
}

func sumOfShortcuts(n int) int {
	i := 2
	sum := 0
	for i < n {
		sum += i - 1
		i++
	}
	return sum
}

func sumOfReverseShortcuts(n int) int {
	i := 0
	sum := 0
	for i < n - 1 {
		sum += n - i
		i++
	}
	return sum
}

func klusownik(nodes []*Node) {
	for {
		rand.Seed(time.Now().UnixNano())
		sec := rand.Float64() * 10
		time.Sleep(time.Second * time.Duration(sec))
		random := rand.Intn(len(nodes))
		fmt.Println("klusownik set trap in node", nodes[random].id)
		nodes[random].trap <- true
	}
}

var p []Packet
var lifetime int
var fmt = log.New(os.Stdout, "", 0)

func main() {
	n := flag.Int("n", 5, "an int")
	d := flag.Int("d", 0, "an int")
	k := flag.Int("k", 8, "an int")
	b := flag.Int("b", 0, "an int")
	h := flag.Int("h", 11, "an int")
	kk := flag.Bool("klusownik", false, "a bool")
	flag.Parse()

	nodes := make([]*Node, *n)
	randoms := make([]int, *n)
	reverseRandoms := make([]int, *n)

	if *d > sumOfShortcuts(*n) || *n <= 0 || *d < 0 || *k < 0 || *b < 0 || *b > sumOfReverseShortcuts(*n) || *h <= 0 {
		fmt.Println("Invalid parameters")
		return
	}
	lifetime = *h

	i := 0
	for i < len(randoms) {
		randoms[i] = 0
		reverseRandoms[i] = 0
		i++
	}

	var node *Node
	for *d > 0 {
		i := 2
		for i < *n {
			var random int

			rand.Seed(time.Now().UnixNano())
			random = rand.Intn(i - 1)
			if random + randoms[i] > i - 1 {
				random = 0
			}
			if random > *d {
				random = *d
			}

			randoms[i] += random
			*d-=random
			i++
		}
	}

	for *b > 0 {
		i := 1
		for i < *n - 1 {
			var random int

			rand.Seed(time.Now().UnixNano())
			random = rand.Intn(2)
			if random + reverseRandoms[i] > *n - i - 1 {
				random = 0
			}
			if random > *b {
				random = *b
			}

			reverseRandoms[i] += random
			*b-=random
			i++
		}
	}

	fmt.Println("Graph:")
	i = 0
	for i < *n {
		tempTabOfShortcuts := make([]*Node, 0)
		busy := make([]int, 0)
		busy = append(busy, i - 1)
		j := 0
		for j < randoms[i] {
			rand.Seed(time.Now().UnixNano())
			random2 := rand.Intn(i)
			was := false

			for _, e := range busy {
				if e == random2 {
					was = true
				}
			}

			if !was {
				busy = append(busy, random2)
				tempTabOfShortcuts = append(tempTabOfShortcuts, nodes[random2])
				j++
			}
		}

		tempTabOfShortcuts = append(tempTabOfShortcuts, node)

		j = 0
		for j < len(tempTabOfShortcuts) && i > 0 {
			draw(*n - i - 1, tempTabOfShortcuts[j].id)
			j++
		}


		node = makeNode(*n - i - 1, tempTabOfShortcuts)
		nodes[i] = node
		i++
	}

	i = 0
	for i < *n {

		tempTabOfShortcuts := make([]*Node, 0)
		busy := make([]int, 0)
		j := 0
		for j < reverseRandoms[i] {
			rand.Seed(time.Now().UnixNano())
			random2 := rand.Intn(*n - i) + i
			if random2 == i {
				continue
			}
			was := false

			for _, e := range busy {
				if e == random2 {
					was = true
				}
			}

			if !was {
				busy = append(busy, random2)
				tempTabOfShortcuts = append(tempTabOfShortcuts, nodes[random2])
				j++
			}
		}

		j = 0
		for j < len(tempTabOfShortcuts) && i > 0 {
			draw(*n - i - 1, tempTabOfShortcuts[j].id)
			j++
		}

		for _, e := range tempTabOfShortcuts {
			nodes[i].append(e)
		}
		i++
	}
	fmt.Println()

	announce := make(chan *Node)
	p = initializePackets(*k)

	i = 1
	for i < len(nodes) {
		go nodes[i].procedure(announce)
		i++
	}

	go func() {
		rand.Seed(time.Now().UnixNano())
		sec := rand.Float64() * 5
		time.Sleep(time.Second * time.Duration(sec))
		for _, packet := range p {
			nodes[*n - 1].packet <- packet
		}
	}()

	go nodes[0].lastNodeListener(announce)
	if *kk {
		go klusownik(nodes)
	}


	i = 0
	for i < *k {
		<-announce
		i++
	}

	fmt.Println("\nPackets serviced by nodes:")
	i = 0
	for i < len(nodes) {
		var s strings.Builder
		s.WriteString(strconv.Itoa(nodes[i].id) + ": [")
		j := 0
		for j < len(nodes[i].serviced) {
			s.WriteString(strconv.Itoa(nodes[i].serviced[j].id))
			j++
			if j < len(nodes[i].serviced) {
				s.WriteString(", ")
			}
		}
		s.WriteString("]")
		fmt.Println(s.String())
		i++
	}

	fmt.Println("\nNodes visited by pockets:")
	i = 0
	for i < len(p) {
		var s strings.Builder
		s.WriteString(strconv.Itoa(p[i].id) + ": [")
		j := 0
		for j < len(p[i].visited) {
			s.WriteString(strconv.Itoa(p[i].visited[j].id))
			j++
			if j < len(p[i].visited) {
				s.WriteString(", ")
			}
		}
		s.WriteString("]")
		fmt.Println(s.String())
		i++
	}
}
