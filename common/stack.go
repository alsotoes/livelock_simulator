package common

type Node struct {
	Value string
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	Nodex []*Node
	count int
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Node) {
	if s.count >= len(s.Nodex) {
		Nodex := make([]*Node, len(s.Nodex)*2)
		copy(Nodex, s.Nodex)
		s.Nodex = Nodex
	}
	s.Nodex[s.count] = n
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Node {
	if s.count == 0 {
		return nil
	}
	node := s.Nodex[s.count-1]
	s.count--
	return node
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	Nodex []*Node
	head  int
	tail  int
	count int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		Nodex := make([]*Node, len(q.Nodex)*2)
		copy(Nodex, q.Nodex[q.head:])
		copy(Nodex[len(q.Nodex)-q.head:], q.Nodex[:q.head])
		q.head = 0
		q.tail = len(q.Nodex)
		q.Nodex = Nodex
	}
	q.Nodex[q.tail] = n
	q.tail = (q.tail + 1) % len(q.Nodex)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.Nodex[q.head]
	q.head = (q.head + 1) % len(q.Nodex)
	q.count--
	return node
}

/*
func main() {
	s := &Stack{Nodex: make([]*Node, 3)}
	s.Push(&Node{1})
	s.Push(&Node{2})
	s.Push(&Node{3})
	fmt.Printf("%v, %v, %v\n", s.Pop().Value, s.Pop().Value, s.Pop().Value)

	q := &Queue{Nodex: make([]*Node, 3)}
	q.Push(&Node{1})
	q.Push(&Node{2})
	q.Push(&Node{3})
	fmt.Printf("%v, %v, %v\n", q.Pop().Value, q.Pop().Value, q.Pop().Value)
}
*/
