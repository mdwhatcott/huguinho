package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/mdwhatcott/huguinho/contracts"
)

func main() {
	log.SetFlags(log.Lshortfile)
	root := "/Users/mike/src/github.com/mdwhatcott/blog/content"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		all, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		parts := strings.Split(string(all), "+++")
		meta := parts[1]
		if len(parts) > 3 {
			return fmt.Errorf("too many parts: %d [%s]", len(meta), path)
		}
		log.Println(path)
		var old contracts.FrontMatter__DEPRECATED
		_, err = toml.Decode(meta, &old)
		if err != nil {
			return err
		}

		slug := path
		slug = strings.TrimSuffix(slug, ".md") + "/"
		slug = strings.TrimPrefix(slug, root)

		new := contracts.ArticleMetadata{
			Draft: old.IsDraft,
			Slug:  slug,
			Title: old.Title,
			Intro: old.Description,
			Tags:  old.Tags,
			Date:  old.Date,
		}

		formatted := format(new)

		//log.Printf("META:\n\n%s\n\n--------\n\n", formatted + parts[2])

		return ioutil.WriteFile(path, []byte(formatted+parts[2]), 0644)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func format(metadata contracts.ArticleMetadata) string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "date:  %s\n", metadata.Date.Format("2006-01-02"))
	_, _ = fmt.Fprintf(&builder, "slug:  %s\n", metadata.Slug)
	_, _ = fmt.Fprintf(&builder, "tags:  %s\n", strings.Join(metadata.Tags, " "))
	_, _ = fmt.Fprintf(&builder, "title: %s\n", metadata.Title)
	_, _ = fmt.Fprintf(&builder, "intro: %s\n", metadata.Intro)
	_, _ = fmt.Fprintf(&builder, "draft: %t\n", metadata.Draft)
	_, _ = fmt.Fprint(&builder, "\n+++")
	return builder.String()
}
