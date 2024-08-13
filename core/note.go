package core

type Note struct {
	Channel  uint8
	Key      uint8
	Velocity uint8
}

func (n Note) IsValid() bool {
	return n.Key == 0
}
