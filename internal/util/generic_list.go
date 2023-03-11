package util

type List[T any] struct {
	Elements map[string]*T
}

func (c *List[T]) Get(key string) (*T, bool) {
	element, ok := c.Elements[key]
	return element, ok
}

func (c *List[T]) Put(key string, element *T) {
	c.Elements[key] = element
}

func (c *List[T]) Delete(key string) {
	delete(c.Elements, key)
}

func (c *List[T]) Length() int {
	n := len(c.Elements)
	return n
}

func (c *List[T]) All() []*T {
	elements := []*T{}
	for _, el := range c.Elements {
		elements = append(elements, el)
	}
	return elements
}
