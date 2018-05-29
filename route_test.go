package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMapWithAllways(t *testing.T) {
	routes := NewRouter()

	routes.Add(NewRoute("a", Allways))
	routes.Add(NewRoute("a", Allways))

	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))
	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))
	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 2, len(routes.ByName("a")))
	assert.Equal(t, 2, len(routes.ByCond(NewValueMessage("a", 0))))

	assert.Equal(t, 3, len(routes.ByName("b")))
	assert.Equal(t, 3, len(routes.ByCond(NewValueMessage("b", 1))))
}

func TestRoutesMap(t *testing.T) {
	routes := NewRouter()

	routes.Add(NewRoute("a", func(m Message) bool { return m.Value().(int) == 0 }))
	routes.Add(NewRoute("a", func(m Message) bool { return m.Value().(int) == 0 }))

	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))
	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))
	routes.Add(NewRoute("b", func(m Message) bool { return m.Value().(int) == 1 }))

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 0, len(routes.ByName("c")))
	assert.Equal(t, 0, len(routes.ByCond(NewValueMessage("c", -1))))
	assert.Equal(t, 0, len(routes.Cases(NewValueMessage("c", -1))))

	assert.Equal(t, 2, len(routes.ByName("a")))
	assert.Equal(t, 2, len(routes.ByCond(NewValueMessage("a", 0))))
	assert.Equal(t, 2, len(routes.Cases(NewValueMessage("a", 0))))

	assert.Equal(t, 3, len(routes.ByName("b")))
	assert.Equal(t, 3, len(routes.ByCond(NewValueMessage("b", 1))))
	assert.Equal(t, 3, len(routes.Cases(NewValueMessage("b", 1))))
}
