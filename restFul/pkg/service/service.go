package service

type Auth interface {
	GenerateToken(username, password string) (string, error)
}

type Service struct {
	Auth
}

//func NewService()
