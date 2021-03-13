package contracts

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}
