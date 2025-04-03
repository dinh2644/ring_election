package main

import (
	"fmt"
	"math/rand/v2"
)

/* **************************************************** Classes (& Methods) **************************************************** */
type Message struct {
	id    int
	attrj int
}

type Node struct {
	msg        Message
	state      bool
	leader_ack bool
}

/* **************************************************** Functions **************************************************** */

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func SendElectedMessage(elected_leader Node, n int, elected_pos int, nodes []Node) {
	i := elected_pos
	starting_node := nodes[elected_pos]
	for i < 1000000 {
		// Determine the next node in a circular manner
		receiving_node := &nodes[(i+1)%n]

		// Set all participating state to false
		if receiving_node.state {
			receiving_node.state = false
			receiving_node.leader_ack = true
		}

		// Stop case: Elected position has been reached
		if receiving_node.msg.id == starting_node.msg.id {
			return
		}

		// Move to next node in ring
		i = (i + 1) % n
	}
}

func StartElection(initiator int, nodes []Node, n int, attrx int, receiving_node Node) Node {
	i := initiator
	elected_node := &nodes[i]
	elected_node.leader_ack = true // leader of course acknowledges itself as leader
	for i < 1000000 {
		// Determine the next node in a circular manner
		receiving_node := nodes[(i+1)%n]

		// Stop case: traveled the whole ring, so Pi declares itself the leader
		if receiving_node.msg.attrj == attrx {
			elected_leader := receiving_node
			elected_pos := i
			SendElectedMessage(elected_leader, n, elected_pos, nodes)
			return receiving_node
		}

		// Case 1: If the received <attribute, ID> is greater than its own
		if attrx > receiving_node.msg.attrj {
			// forward message (election, <attrx, x>) to successor (no updating)

			// set state to participating
			nodes[(i+1)%n].state = true
		}

		// Case 2: If the received <attribute, ID> is less than its own:
		if attrx < receiving_node.msg.attrj {
			if !receiving_node.state {
				// send (election, <attrj, j>) to successor
				attrx = receiving_node.msg.attrj
				// set state to participating
				nodes[(i+1)%n].state = true
			}

		}

		// Track nodes visited
		fmt.Printf("Processed node %d (attr: %d)\n", receiving_node.msg.id, receiving_node.msg.attrj)

		// Move to next node in ring
		i = (i + 1) % n
	}

	return Node{msg: Message{}, state: false}
}

/* **************************************************** MAIN **************************************************** */
func main() {
	n := 5
	var nodes []Node

	// Setup n nodes
	for i := range n {
		node := Node{
			msg: Message{
				id:    i + 1,
				attrj: randRange(0, 100),
			},
			state:      false,
			leader_ack: false,
		}
		nodes = append(nodes, node)
	}

	// Randomly select P(i)
	initiator := randRange(1, n)

	attrx := nodes[initiator].msg.attrj
	receiving_node := nodes[(initiator+1)%n]

	// Start election
	fmt.Printf("\nStart Election:\n")
	for i := range n {
		fmt.Printf("Node %d (attr: %d)\n", nodes[i].msg.id, nodes[i].msg.attrj)
	}

	start_node := nodes[initiator]
	fmt.Printf("\nStart node: %d\n\n", start_node.msg.id)

	leader_node := StartElection(initiator, nodes, n, attrx, receiving_node)

	fmt.Printf("\nLeader elected: Node %d (attr: %d)\n\n", leader_node.msg.id, nodes[leader_node.msg.id-1].msg.attrj)

	fmt.Println("Did all node acknowledge new leader?\n")
	for i := range n {
		fmt.Printf("Node %d (leader_ack: %t)\n", nodes[i].msg.id, nodes[i].leader_ack)
	}
	fmt.Println()

}
