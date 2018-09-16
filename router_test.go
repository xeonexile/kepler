package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMapWithAllways(t *testing.T) {
	r := NewRouter(false).(*router)

	r.AddRoute(DefaultRouteGroup, Allways)
	r.AddRoute(DefaultRouteGroup, Allways)

	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(r.routes))
	assert.Equal(t, 2, len(r.byName(DefaultRouteGroup)))
	assert.Equal(t, 2, len(r.byCond(NewMessage(".", 0))))

	assert.Equal(t, 3, len(r.byName("b")))
	assert.Equal(t, 3, len(r.byCond(NewMessage("b", 1))))
}

func TestRouterClose(t *testing.T) {
	r := NewRouter(false).(*router)

	r.AddRoute(DefaultRouteGroup, Allways)
	r.AddRoute(DefaultRouteGroup, Allways)

	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(r.routes))
	assert.Equal(t, 2, len(r.byName(DefaultRouteGroup)))
	assert.Equal(t, 3, len(r.byName("b")))
	r.Close()
	assert.Equal(t, 2, len(r.routes))

	assert.Equal(t, 0, len(r.byName(DefaultRouteGroup)))
	assert.Equal(t, 0, len(r.byName("b")))
}

func TestRouteClose(t *testing.T) {
	r := NewRouter(false).(*router)

	a1 := r.AddRoute(DefaultRouteGroup, Allways)
	a2 := r.AddRoute(DefaultRouteGroup, Allways)

	b1 := r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	b2 := r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	b3 := r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(r.routes))
	assert.Equal(t, 2, len(r.byName(DefaultRouteGroup)))
	assert.Equal(t, 3, len(r.byName("b")))

	a1.Close()
	assert.Equal(t, 1, len(r.byName(DefaultRouteGroup)))

	a1.Close()
	a2.Close()
	b1.Close()
	b2.Close()
	assert.Equal(t, 0, len(r.byName(DefaultRouteGroup)))
	assert.Equal(t, 1, len(r.byName("b")))

	b3.Close()
	assert.Equal(t, 0, len(r.byName("b")))
}

func TestRoutesMap(t *testing.T) {
	r := NewRouter(false).(*router)

	r.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })
	r.AddRoute("a", func(m Message) bool { return m.Value().(int) == 0 })

	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })
	r.AddRoute("b", func(m Message) bool { return m.Value().(int) == 1 })

	assert.Equal(t, 2, len(r.routes))
	assert.Equal(t, 0, len(r.byName("c")))
	assert.Equal(t, 0, len(r.byCond(NewMessage("c", -1))))

	assert.Equal(t, 2, len(r.byName("a")))
	assert.Equal(t, 2, len(r.byCond(NewMessage("a", 0))))

	assert.Equal(t, 3, len(r.byName("b")))
	assert.Equal(t, 3, len(r.byCond(NewMessage("b", 1))))

	r.RemoveRoute(r.byName("b")[0])
	assert.Equal(t, 2, len(r.byName("b")))

	r.byName("b")[0].Close()
	assert.Equal(t, 1, len(r.byName("b")))
}
