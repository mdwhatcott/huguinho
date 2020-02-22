package contracts

import "errors"

type Config struct {
	TemplateDir string
	ContentRoot string
	TargetRoot  string
	BuildDrafts bool
	BuildFuture bool
}

func (this Config) Validate() error {
	if this.TemplateDir == "" {
		return errors.New("template directory is required")
	}
	if this.ContentRoot == "" {
		return errors.New("content directory is required")
	}
	if this.TargetRoot == "" {
		return errors.New("target directory is required")
	}
	return nil
}
