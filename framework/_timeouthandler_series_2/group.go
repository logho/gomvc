package framework_1

type IGroup interface {
	Get(string, ControllerHandler, ...int)
	Post(string, ControllerHandler, ...int)
	Put(string, ControllerHandler, ...int)
	Delete(string, ControllerHandler, ...int)

	Group(string) IGroup
}

type Group struct {
	core   *Core  //point to  core
	parent *Group // point to Parent group
	prefix string // group prefix
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

func (g *Group) Get(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	url = g.getCompletePath() + url
	g.core.Get(url, handler, limit)
}

func (g *Group) Post(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	url = g.getCompletePath() + url
	g.core.Post(url, handler, limit)
}

func (g *Group) Put(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	url = g.getCompletePath() + url
	g.core.Put(url, handler, limit)
}

func (g *Group) Delete(url string, handler ControllerHandler, timeout ...int) {
	limit := checkTimeout(timeout...)
	url = g.getCompletePath() + url
	g.core.Delete(url, handler, limit)
}

func (g *Group) Group(prefix string) IGroup {
	group := NewGroup(g.core, prefix)
	group.parent = g
	return group
}

func (g *Group) getCompletePath() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getCompletePath() + g.prefix
}
