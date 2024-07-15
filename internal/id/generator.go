package id

type Generator struct {
	nextId int
}

func NewGenerator() *Generator {
	return &Generator{
		nextId: 1,
	}
}

func (g *Generator) NextID() int {
	defer func() {
		g.nextId += 1
	}()

	return g.nextId
}
