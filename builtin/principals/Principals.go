package principals

// Principals is the interface that wrap the get of SMK Principals.
type Principals interface {
	Init(config map[string]string) error
	Get(user string) ([]string, error)
}
