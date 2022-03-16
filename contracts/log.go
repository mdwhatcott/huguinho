package contracts

type Logger interface {
	Print(...any)
	Printf(string, ...any)
	Println(...any)
}
