package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type Tile int

const (
	WATER Tile = iota
	LAND
)

func max(a ...int) int {
	m := a[0]
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}

type Vect struct {
	x, y int
}

// Get the coordinates of the surrounding 8 tiles.
func nearby(v Vect) []Vect {
	return []Vect{
		Vect{v.x - 1, v.y + 1},
		Vect{v.x - 1, v.y},
		Vect{v.x - 1, v.y - 1},
		Vect{v.x, v.y + 1},
		Vect{v.x, v.y - 1},
		Vect{v.x + 1, v.y + 1},
		Vect{v.x + 1, v.y},
		Vect{v.x + 1, v.y - 1}}
}

func randVect(w, h int) Vect {
	return Vect{rand.Intn(w), rand.Intn(h)}
}

type TileMap struct {
	W, H  int
	tiles []Tile
	ids   []int
}

func CreateMap(w, h int) *TileMap {
	return &TileMap{w, h, make([]Tile, w*h), make([]int, w*h)}
}

func (m *TileMap) seedPos(v Vect, t Tile) {
	m.Set(v, t)
	for _, n := range m.GetNeighbors(v) {
		if rand.Intn(8) != 0 {
			m.Set(n, t)
		}
	}
}

func (m *TileMap) Seed() {
	m.seedPos(randVect(m.W, m.H), LAND)
	for i := 0; i < 4; i++ {
		t := WATER
		if i%2 == 0 {
			t = LAND
		}
		for j := 0; j < max(1, (m.W*m.H)/20); j++ {
			m.seedPos(randVect(m.W, m.H), t)
		}
	}
}

func (m *TileMap) Get(v Vect) Tile {
	return m.tiles[v.y*m.W+v.x]
}

func (m *TileMap) Set(v Vect, t Tile) {
	m.tiles[v.y*m.W+v.x] = t
}

func (m *TileMap) GetId(v Vect) int {
	return m.ids[v.y*m.W+v.x]
}

func (m *TileMap) SetId(v Vect, id int) {
	m.ids[v.y*m.W+v.x] = id
}

func (m *TileMap) String() string {
	s := []rune{}
	for i := range m.ids {
		if i%m.W == 0 {
			s = append(s, '\n')
		}
		r := '^'
		if m.tiles[i] != WATER {
			v := m.ids[i] - 1
			if v < 26 {
				// Use uppercase first...
				r = rune(v + 'A')
			} else if v < 52 {
				// then lowercase...
				r = rune(v - 26 + 'a')
			} else {
				// then a ? if we run out.
				r = '?'
			}
		}
		s = append(s, r, ' ')
	}
	return string(s)
}

func (m *TileMap) GetNeighbors(v Vect) []Vect {
	r := make([]Vect, 0, 8)
	for _, n := range nearby(v) {
		// Borders don't wrap.
		if n.x >= 0 && n.x < m.W && n.y >= 0 && n.y < m.H {
			r = append(r, n)
		}
	}
	return r
}

func (m *TileMap) CalcIds() {
	id := 0
	// Assign a unique id to all homogenous nodes connected to each boundary node.
	bnd := []Vect{Vect{0, 0}}
	for len(bnd) != 0 {
		v := bnd[0]
		bnd = bnd[1:]
		if m.GetId(v) != 0 {
			// Already visited; also handled below, but
			// this shortcuts before creating new ids on revisits.
			continue
		}
		id += 1
		t := m.Get(v)
		// Visit all connected homogeneous nodes, noting location of heterogeneous nodes.
		stack := []Vect{v}
		for len(stack) != 0 {
			pos := stack[0]
			stack = stack[1:]
			if m.GetId(pos) != 0 {
				// Already visited, skip.
			} else if m.Get(pos) != t {
				// Heterogenous (aka boundary) node, add to set of boundary nodes to visit.
				bnd = append(bnd, pos)
			} else {
				// Homogeneous node, set id and add neighbors.
				m.SetId(pos, id)
				for _, neighbor := range m.GetNeighbors(pos) {
					stack = append(stack, neighbor)
				}
			}
		}
	}
}

func main() {
	var w, h int
	var seed int64
	flag.IntVar(&w, "width", 15, "the map width")
	flag.IntVar(&h, "height", 15, "the map height")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "a seed value, defaults to nanoseconds since unix epoch")
	flag.Parse()

	rand.Seed(seed)
	island := CreateMap(w, h)
	island.Seed()
	island.CalcIds()
	fmt.Println(island)
}
