package dockerfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Dockerfile represents a parsed Dockerfile.
type Dockerfile struct {
	// filename is the name of the Dockerfile.
	filename string
	// content is the parsed result of the Dockerfile.
	content  *parser.Result

	// stages are the defined build stages in the Dockerfile.
	stages  []string
	// args are the build arguments in the Dockerfile.
	args    map[string]string
	// secrets are the secrets used in the Dockerfile.
	secrets []string
}

// NewDockerfile parses a Dockerfile from a given file and returns a Dockerfile
// object with extracted stages, args, and secrets.
func NewDockerfile(filename string, file *os.File) (*Dockerfile, error) {
	content, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	stages := []string{}
	args := map[string]string{}
	secrets := []string{}

	for _, child := range content.AST.Children {
		switch child.Value {
		case "FROM":
			args := []string{}
			for next := child.Next; next != nil; next = next.Next {
				args = append(args, next.Value)
			}

			// We skip if there's no stage defined
			if len(args) != 3 {
				continue
			}

			stages = append(stages, args[2])
		case "ARG":
			// Args does not handle self interpolation for simplicity.
			// TODO: handle self interpolation (ARG XXX="XX-${XXXX}")
			entry := strings.SplitN(child.Next.Value, "=", 2)
			switch len(entry) {
			case 1:
				args[entry[0]] = ""
			case 2:
				args[entry[0]] = entry[1]
			default:
				panic(fmt.Errorf("invalid ARG: %s", child.Next.Value))
			}
		case "RUN":
			// Parse RUN command to extract secrets if it exists
			for _, flag := range child.Flags {
				if !strings.Contains(flag, "--mount=type=secret,id=") {
					continue
				}

				secrets = append(secrets, strings.TrimPrefix(flag, "--mount=type=secret,id="))
			}
		}
	}

	return &Dockerfile{
		filename: filename,
		content:  content,
		stages:   stages,
		args:     args,
		secrets:  secrets,
	}, nil
}

// Filename returns the filename of the Dockerfile.
func (d *Dockerfile) Filename() string {
	return d.filename
}

// Stages returns build stages defined in the Dockerfile.
func (d *Dockerfile) Stages() []string {
	return d.stages
}

// Args returns build arguments defined in the Dockerfile.
func (d *Dockerfile) Args() map[string]string {
	return d.args
}

// Secrets returns secrets defined in the Dockerfile.
func (d *Dockerfile) Secrets() []string {
	return d.secrets
}

// String displays the Dockerfile content.
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
