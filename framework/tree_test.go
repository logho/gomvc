package framework

import (
	"sort"
	"strings"
	"testing"
)

func TestMatchRoute(t *testing.T) {

}
func TestString(t *testing.T) {
	url := "user"
	segments := strings.Split(url, "/")
	for i, segment := range segments {
		t.Log(i, segment)
	}
	t.Log(strings.ToUpper(":ooooxxx"))

	url = "/xx/xx/xx"
	prefix := strings.TrimPrefix(url, "/")
	t.Log(prefix)

}
func TestStringSort(t *testing.T) {
	arr := []string{"a/b/c", "a/b/:x", "a/:x/c", ":x/b/c", "a/:x/:x", ":x/b/:x", ":x/:x/c", ":x/:x/:x"}
	for _, i := range arr {
		t.Log(i)
	}
	sort.Strings(arr)
	////sort.Strings(arr)
	//sort.Slice(arr, func(i, j int) bool {
	//	return i < j
	//
	//})
	//sort.Reverse(arr)
	t.Log("===========")
	for _, i := range arr {
		t.Log(i)
	}
}
func TestSlice(t *testing.T) {
	arr := []int{2, 1, 3, 41, 123, -1}
	for _, i := range arr {
		t.Log(i)
	}
	sort.Ints(arr)
	t.Log("===========")
	for _, i := range arr {
		t.Log(i)
	}

}

func Test_filterChildNodes(t *testing.T) {
	root := &node{
		isLeaf:  false,
		segment: "",
		handler: func(*Context) error { return nil },
		children: []*node{
			{
				isLeaf:   true,
				segment:  "FOO",
				handler:  func(*Context) error { return nil },
				children: nil,
			},
			{
				isLeaf:   false,
				segment:  ":id",
				handler:  nil,
				children: nil,
			},
		},
	}

	{
		nodes := root.getMatchedChildNodes("FOO", MATCH_MODE)
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}

	{
		nodes := root.getMatchedChildNodes(":foo", MATCH_MODE)
		if len(nodes) != 2 {
			t.Error(":foo error")
		}
	}

}

func Test_matchNode(t *testing.T) {
	context := &Context{}

	root := &node{
		isLeaf:  false,
		segment: "",
		handler: func(*Context) error { return nil },
		children: []*node{
			{
				isLeaf:  false,
				segment: "FOO",
				handler: nil,
				children: []*node{
					&node{
						isLeaf:  true,
						segment: "BAR",
						handler: func(ctx *Context) error {
							t.Log("foo/bar")
							return nil
						},
						children: []*node{},
					},
				},
			},
			{
				isLeaf:  false,
				segment: ":ooo",
				handler: nil,
				children: []*node{
					&node{
						isLeaf:  false,
						segment: ":xxx",
						handler: func(ctx *Context) error {
							t.Log(":ooo/:xxx")
							return nil
						},
						children: []*node{
							&node{
								isLeaf:  true,
								segment: "1",
								handler: func(ctx *Context) error {
									t.Log(":ooo/:xxx/1")
									return nil
								},
								children: []*node{},
							},
						},
					},
				},
			},
		},
	}

	{
		nodes := root.matchAllNode("foo/bar")
		t.Log("path is foo/bar", "the matchedNode is")
		for _, node := range nodes {
			if node != nil {
				node.handler(context)
			}
		}
	}

	{
		nodes := root.matchAllNode("test")
		t.Log("path is test", "the matchedNode is")
		for _, node := range nodes {
			if node != nil {
				node.handler(context)
			}
		}
	}

	{
		nodes := root.matchAllNode("foo")
		t.Log("path is foo", "the matchedNode is")
		for _, node := range nodes {
			if node != nil {
				node.handler(context)
			}
		}
	}

	{
		nodes := root.matchAllNode("1/2")
		t.Log("path is 1/2", "the matchedNode is")
		for _, node := range nodes {
			if node != nil {
				node.handler(context)
			}
		}
	}
	{
		nodes := root.matchAllNode("far/bar/1")
		t.Log("path is far/bar/1", "the matchedNode is")
		for _, node := range nodes {
			if node != nil {
				node.handler(context)
			}
		}
	}

}

func TestTree_AddRouter(t *testing.T) {
	//tree := &Tree{
	//	root: NewNode("/"),
	//}
	////tree.AddRouter("/user/info/1", nil)
	////tree.AddRouter("user/info/1", nil)
	////tree.AddRouter("user/:info/1", nil)
	////tree.AddRouter("user/:info/2", nil)
	////tree.AddRouter("user/:info/2", nil)
	////tree.AddRouter("user/:test/2", nil)

}

func deferA() (i int) {
	defer func() {
		i++
		i++
		i++

	}()
	return i
}
func TestDefer(t *testing.T) {
	m := map[string]string{}
	m["mysql"] = "mysql"
	i, err := m["mysql"]
	t.Log(i, err)

}
