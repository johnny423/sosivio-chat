package auth

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

//DummyAuthenticator is a type that simulate authentication using a const mapping of users credentials
type DummyAuthenticator struct{}

//Authenticate authenticates the credentials using a const mapping of users
func (auth *DummyAuthenticator) Authenticate(credentials Credentials) bool {
	expectedPassword, ok := users[credentials.Username]
	return ok && expectedPassword == credentials.Password
}
