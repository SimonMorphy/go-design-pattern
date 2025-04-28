package creational

// Cloneable <-*Impl Here*-> PrototypePattern
type Cloneable interface {
	Clone() Cloneable
	Print() string
}
