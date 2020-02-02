package contracts

type ContentConverter interface {
	Convert(content string) (string, error)
}
