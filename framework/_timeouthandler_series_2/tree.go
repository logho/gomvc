package framework_1

import (
	"errors"
	"strings"
)

type TraversalMode int

const (
	CHECK_MODE TraversalMode = 1
	MATCH_MODE TraversalMode = 2
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: NewNode("/")}

}
func (tree *Tree) GetHandler(url string) (ControllerHandler, int) {
	matchNode := tree.root.matchNode(url, MATCH_MODE)
	if matchNode != nil {
		return matchNode.handler, matchNode.timeout
	}
	return nil, 0
}

func (tree *Tree) AddRouter(url string, handler ControllerHandler, timeout int) error {
	n := tree.root
	url = strings.TrimPrefix(url, "/")
	matchNode := tree.root.matchNode(url, CHECK_MODE)
	if matchNode != nil {
		return errors.New("current route is conflict: " + url)
	}
	segments := strings.Split(url, "/")

	for index, segment := range segments {

		segment = strings.ToUpper(segment)
		isLeaf := index == len(segments)-1
		if n.wildChild && isWildCardSegment(segment) && segment != n.children[len(n.children)-1].segment {
			return errors.New("current route is conflict: " + url)
		}

		var candidateNode *node

		childNodes := n.getMatchedChildNodes(segment, MATCH_MODE)
		if childNodes != nil && len(childNodes) > 0 {
			for _, childNode := range childNodes {
				if childNode.segment == segment {
					candidateNode = childNode
					break
				}
			}
		}

		if candidateNode == nil {
			newChild := NewNode(segment)
			if isLeaf {
				newChild.isLeaf = true
				newChild.handler = handler
				newChild.timeout = timeout
			}

			n.addChild(newChild)
			if isWildCardSegment(segment) {
				n.wildChild = true
			}
			candidateNode = newChild

		}
		n = candidateNode
	}
	return nil
}

type node struct {
	segment   string
	handler   ControllerHandler
	children  []*node
	isLeaf    bool
	timeout   int //millisecond timeout
	wildChild bool
}

func NewNode(segment string) *node {
	return &node{
		segment: segment,
	}
}

func isWildCardSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) addChild(child *node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

func (n *node) getMatchedChildNodes(segment string, mode TraversalMode) []*node {
	if len(n.children) == 0 {
		return nil
	}

	//path为通配符
	if isWildCardSegment(segment) && mode == MATCH_MODE {
		return n.children
	}
	matchedChildNode := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		//pattern为通配符或者 path和pattern的路径相同
		if (isWildCardSegment(child.segment) && MATCH_MODE == mode) || child.segment == segment {
			matchedChildNode = append(matchedChildNode, child)
		}
	}
	return matchedChildNode
}

func (n *node) matchNode(url string, mode TraversalMode) *node {
	url = strings.TrimPrefix(url, "/")
	segments := strings.SplitN(url, "/", 2)
	segment := strings.ToUpper(segments[0])

	matchedChildNodes := n.getMatchedChildNodes(segment, mode)

	if matchedChildNodes == nil || len(matchedChildNodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		for _, node := range matchedChildNodes {
			if node.isLeaf {
				return node
			}
		}
		return nil
	}

	for _, childNode := range matchedChildNodes {
		matchNode := childNode.matchNode(segments[1], mode)
		return matchNode
	}

	return nil
}

// ---------- old version -------------

// node' matchAllNode return all nodes ([]*node)  whether it is wildCardNode
func (n *node) matchAllNode(url string) []*node {
	segments := strings.SplitN(url, "/", 2)
	segment := strings.ToUpper(segments[0])

	matchedChildNodes := n.getMatchedChildNodes(segment, MATCH_MODE)
	if matchedChildNodes == nil || len(matchedChildNodes) == 0 {
		return nil
	}
	matchedNodes := make([]*node, 0, len(matchedChildNodes))

	if len(segments) == 1 {
		for _, node := range matchedChildNodes {
			if node.isLeaf {
				matchedNodes = append(matchedNodes, node)
			}
		}
		return matchedNodes
	}
	// user/:id/:age/info  user/id/name
	for _, node := range matchedChildNodes {
		nodes := node.matchAllNode(segments[1])
		if nodes != nil && len(nodes) > 0 {
			matchedNodes = append(matchedNodes, nodes...)
		}
	}
	return matchedNodes

}

//node's childNodes is unsorted, to find the best router we need sorted
//all matchedNode by full path and  the first one is the best.
//but the performance is bad when we get a handler
func (n *node) bestMatchedNode(url string) *node {
	matchNode := n.matchAllNode(url)
	if len(matchNode) == 0 {
		return nil
	} else if len(matchNode) == 1 {
		return matchNode[0]
	} else {
		// a b c
		// a b c
		// x b c
		// a x c
		// a b x
		// x x c
		// x b x
		// a x x
		// x x x

		// sort nodes by node.segment string and choose the first one to return

	}
	return nil
}

// ---------- old version -------------
