package contracts

type Config struct {
	Author      string
	TemplateDir string
	ContentRoot string
	TargetRoot  string
	BasePath    string
	BuildDrafts bool
	BuildFuture bool
}
