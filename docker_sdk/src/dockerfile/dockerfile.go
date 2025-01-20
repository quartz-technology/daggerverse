package dockerfile

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Dockerfile struct {
	filename string
	content *parser.Result
}

func New(filename string, file *os.File) (*Dockerfile, error) {
	content, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	return &Dockerfile{
		filename: filename,
		content: content,
	}, nil
}

func (d *Dockerfile) Filename() string {
	return d.filename
}

func (d *Dockerfile) Stages() []string {
	stages := []string{}

	for _, child := range d.content.AST.Children {
		if child.Value != "FROM" {
			continue
		}

		args := []string{}
		for next := child.Next; next != nil; next = next.Next {
			args = append(args, next.Value)
		}

		// We skip if there's no stage defined
		if len(args) != 3 {
			continue
		}

		stages = append(stages, args[2])
	}

	return stages
}

func (d *Dockerfile) String() string {
	var result string

	for _, child := range d.content.AST.Children {
		result += fmt.Sprintf("Command: %s\n", child.Value)
		for _, flag := range child.Flags {
			result += fmt.Sprintf("  Flag: %s\n", flag)
		}

		argIndex := 1
		for next := child.Next; next != nil; next = next.Next {
			result += fmt.Sprintf("  Argument %d: %s\n", argIndex, next.Value)
			argIndex++
		}
	}

	return result
}