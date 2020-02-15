package contracts

type Config struct {
	TemplateDir string
	StylesDir   string
	ContentRoot string
	TargetRoot  string
	BuildDrafts bool
	BuildFuture bool
}
