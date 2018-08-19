package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMapWithAllways(t *testing.T) {
	routes := NewRouter(false).(*router)

	routes.AddRoute(DefaultRouteGroup, Allways)
	routes.AddRoute(DefaultRouteGroup, Allways)

	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 2, len(routes.byName(DefaultRouteGroup)))
	assert.Equal(t, 2, len(routes.byCond(NewMessage(".", 0))))

	assert.Equal(t, 3, len(routes.byName("b")))
	assert.Equal(t, 3, len(routes.byCond(NewMessage("b", 1))))
}

func TestRoutesMap(t *testing.T) {
	routes := NewRouter(false).(*router)

	routes.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })
	routes.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })

	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	routes.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(routes.routes))
	assert.Equal(t, 0, len(routes.byName("c")))
	assert.Equal(t, 0, len(routes.byCond(NewMessage("c", -1))))

	assert.Equal(t, 2, len(routes.byName("a")))
	assert.Equal(t, 2, len(routes.byCond(NewMessage("a", 0))))

	assert.Equal(t, 3, len(routes.byName("b")))
	assert.Equal(t, 3, len(routes.byCond(NewMessage("b", 1))))

	routes.RemoveRoute(routes.byName("b")[0])
	assert.Equal(t, 2, len(routes.byName("b")))

	routes.byName("b")[0].Close()
	assert.Equal(t, 1, len(routes.byName("b")))
}
