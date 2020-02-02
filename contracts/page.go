package contracts

type Page struct {
	Metadata JSONFrontMatter
	Content  Content
}
type Content struct {
	Original  string
	Converted string
}
