package blink

type Metatile struct {
	Zoom     int
	X        int
	Y        int
	Size     int
	Members  []*Tile
	Children []*Tile
}

type Tile struct {
	Zoom     int
	X        int
	Y        int
	Key      string
	Members  []*Metatile
	Children []*Metatile
	Metatile *Metatile
}

type Bounds struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

// NewMetatile allocates and returns a new Tile
func NewMetatile(zoom int, x int, y int, size int) *Metatile {
	metatile := new(Metatile)
	metatile.Zoom = zoom
	metatile.X = x
	metatile.Y = y
	metatile.Size = size
	return metatile
}

// NewTile allocates and returns a new Tile
func NewTile(zoom int, x int, y int) *Tile {
	tile := new(Tile)
	tile.Zoom = zoom
	tile.X = x
	tile.Y = y
	return tile
}

func (metatile *Metatile) AddChildrenTo() {

	if len(metatile.Children) == 0 {
		return
	}

	var tl, tr, bl, br *Metatile

	var midX = metatile.X*2 + metatile.Size
	var midY = metatile.Y*2 + metatile.Size

	for _, child := range metatile.Children {

		// Are we on the left side?
		if child.X < midX {

			// Are we on the top left side?
			if child.Y < midY {
				if tl == nil {
					tl = NewMetatile(metatile.Zoom+1, metatile.X*2, metatile.Y*2, metatile.Size)
					child.Metatile = tl
					tl.Members = append(tl.Members, child)
				}
			} else {
				if bl == nil {
					bl = NewMetatile(metatile.Zoom+1, metatile.X*2, midY, metatile.Size)
					child.Metatile = bl
					bl.Members = append(bl.Members, child)
				}
			}
		} else {

			// Are we on the top right side?
			if child.Y < midY {
				if tr == nil {
					tr = NewMetatile(metatile.Zoom+1, midX, metatile.Y*2, metatile.Size)
					child.Metatile = tr
					tr.Members = append(tr.Members, child)
				}
			} else {
				if br == nil {
					br = NewMetatile(metatile.Zoom+1, midX, midY, metatile.Size)
					child.Metatile = br
					br.Members = append(br.Members, child)
				}
			}
		}
	}
}

func (tile *Tile) AddChildrendInBoundsTo(bounds *Bounds, tiles []*Tile) {
	var zoom = tile.Zoom + 1
	var x = tile.X * 2
	var y = tile.Y * 2
	if y >= bounds.MinY {
		if x >= bounds.MinX {
			tiles = append(tiles, NewTile(zoom, x, y))
		}
		if x+1 <= bounds.MaxX {
			tiles = append(tiles, NewTile(zoom, x+1, y))
		}
	}
	if y+1 <= bounds.MaxY {
		if x >= bounds.MinX {
			tiles = append(tiles, NewTile(zoom, x, y+1))
		}
		if x+1 <= bounds.MaxX {
			tiles = append(tiles, NewTile(zoom, x+1, y+1))
		}
	}

}
