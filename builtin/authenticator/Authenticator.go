package authenticator

// Authenticator is the interface that wrap the SMK Authentication logic.
type Authenticator interface {
	Init(config map[string]string) error
	Login(user, password string) (bool, error)
}
