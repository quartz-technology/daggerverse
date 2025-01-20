package dockerfile

import (
	"fmt"
	"os"
	"strings"

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

// Args does not handle self interpolation for simplicity.
// TODO: handle self interpolation (ARG XXX="XX-${XXXX}")
func (d *Dockerfile) Args() map[string]string {
	args := map[string]string{}

	for _, child := range d.content.AST.Children {
		if child.Value != "ARG" {
			continue
		}

		entry := strings.SplitN(child.Next.Value, "=", 2)
		switch len(entry) {
		case 1:
			args[entry[0]] = ""
		case 2:
			args[entry[0]] = entry[1]
		default:
			panic(fmt.Errorf("invalid ARG: %s", child.Next.Value))
		}
	}

	return args
}

func (d *Dockerfile) Secrets() []string {
	secrets := []string{}

	for _, child := range d.content.AST.Children {
		if child.Value != "RUN" {
			continue
		}

		for _, flag := range child.Flags {
			if !strings.Contains(flag, "--mount=type=secret,id=") {
				continue
			}
			
			secrets = append(secrets, strings.TrimPrefix(flag, "--mount=type=secret,id="))
		}
	}

	return secrets
}

func (d *Dockerfile) String() string {
	var result string

	result += fmt.Sprintf("Filename: %s\n", d.filename)
	result += fmt.Sprintf("Stages: %s\n", strings.Join(d.Stages(), ", "))
	result += fmt.Sprintf("Secrets: %s\n", strings.Join(d.Secrets(), ", "))

	for key, value := range d.Args() {
		result += fmt.Sprintf("ARG %s=%s\n", key, value)
	}

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