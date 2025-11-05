package role

type Role int

const (
	Buyer   Role = iota // 0
	Manager             // 1
	Admin               // 2
)

func (r Role) String() string {
	return [...]string{"buyer", "manager", "admin"}[r]
}
