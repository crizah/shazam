package structs

type Information struct {
	Anchor_time uint32
	SongID      uint32
}

type OMap struct {
	Map   map[uint32]Information
	Order []uint32
}

type Helper struct {
	Artist string
	Name   string
}
