package contracts

type Config struct {
	TemplateDir string
	ContentRoot string
	TargetRoot  string
	BuildDrafts bool
	BuildFuture bool
}
