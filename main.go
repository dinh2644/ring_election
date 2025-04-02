package main

import (
	"fmt"
	"math/rand/v2"
)

/* **************************************************** Classes (& Methods) **************************************************** */
type Message struct {
	id   int
	attr int
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
		successor_node := &nodes[(i+1)%n]

		// Set all participating state to false
		if successor_node.state {
			successor_node.state = false
			successor_node.leader_ack = true
		}

		// Stop case: Elected position has been reached
		if successor_node.msg.id == starting_node.msg.id {
			fmt.Println("STOPPED")
			return
		}

		// Move to next node in ring
		i = (i + 1) % n
	}
}

func StartElection(initiator int, nodes []Node, n int, max_id int, successor_node Node) Node {
	i := initiator
	elected_node := &nodes[i]
	elected_node.leader_ack = true // leader of course acknowledges itself as leader
	for i < 1000000 {
		// Determine the next node in a circular manner
		successor_node := nodes[(i+1)%n]

		// FOR DEBUGGING:
		//fmt.Printf("successor_node id: %d\n", successor_node.msg.id)

		/* Handle passing/receiving message */

		// Stop case: traveled the whole ring, so Pi declares itself the leader
		if successor_node.msg.attr == max_id {
			elected_leader := successor_node
			elected_pos := i
			SendElectedMessage(elected_leader, n, elected_pos, nodes)
			return successor_node
		}
		// Case 1: If the received <attribute, ID> is less than its own:
		if successor_node.msg.attr > max_id {
			max_id = successor_node.msg.attr
			// set state to participating
			nodes[(i+1)%n].state = true
		}
		// Case 2: If the received <attribute, ID> is greater than its own
		if successor_node.msg.attr < max_id {
			// set state to participating
			nodes[(i+1)%n].state = true
		}

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
				id:   i + 1,
				attr: randRange(0, 100),
			},
			state:      false,
			leader_ack: false,
		}
		nodes = append(nodes, node)
	}

	// This test case failed
	// Node 1 (attr: 24)
	// Node 2 (attr: 80)
	// Node 3 (attr: 61)
	// Node 4 (attr: 35)
	// Node 5 (attr: 54)

	// Randomly select P(i)
	initiator := randRange(1, n)

	max_id := nodes[initiator].msg.attr
	successor_node := nodes[(initiator+1)%n]

	// FOR DEBUGGING:
	for i := range n {
		fmt.Printf("Node %d (attr: %d)\n", i+1, nodes[i].msg.attr)
	}
	start_node := nodes[initiator]

	// Start election
	leader_node := StartElection(initiator, nodes, n, max_id, successor_node)

	// FOR DEBUGGING:
	fmt.Printf("Start node: %d\n", start_node.msg.id)
	fmt.Printf("Leader node thread id: %d\n", leader_node.msg.attr)

	// SEE IF ALL NODES ACKNOWLEDGED NEW LEADER
	for i := range n {
		fmt.Printf("Node %d (leader_ack = %t)\n", nodes[i].msg.id, nodes[i].leader_ack)
	}

}
