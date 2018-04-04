package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMap(t *testing.T) {
	routes := NewRoutesMap()

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
