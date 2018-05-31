package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMapWithAllways(t *testing.T) {
	routes := NewRouter().(*router)

	routes.AddRoute("a", Allways)
	routes.AddRoute("a", Allways)

	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 2, len(routes.byName("a")))
	assert.Equal(t, 2, len(routes.byCond(NewValueMessage("a", 0))))

	assert.Equal(t, 3, len(routes.byName("b")))
	assert.Equal(t, 3, len(routes.byCond(NewValueMessage("b", 1))))
}

func TestRoutesMap(t *testing.T) {
	routes := NewRouter().(*router)

	routes.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })
	routes.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })

	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 0, len(routes.byName("c")))
	assert.Equal(t, 0, len(routes.byCond(NewValueMessage("c", -1))))

	assert.Equal(t, 2, len(routes.byName("a")))
	assert.Equal(t, 2, len(routes.byCond(NewValueMessage("a", 0))))

	assert.Equal(t, 3, len(routes.byName("b")))
	assert.Equal(t, 3, len(routes.byCond(NewValueMessage("b", 1))))
}
