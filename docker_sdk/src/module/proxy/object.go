package proxy

import (
	"bytes"
	"fmt"
	"text/template"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
)

type Service struct {
	Service  *dagger.Service
	Name     string
	Alias    string
	Frontend int
	Backend  int
	IsTcp    bool
	Exposed  bool
}

type Proxy struct {
	ctr *dagger.Container
}

func New() *Proxy {
	return &Proxy{
		ctr: dag.Container().
			From("nginx:1.25.3").
			WithNewFile("/etc/nginx/stream.conf", streamConf).
			WithNewFile("/etc/nginx/nginx.conf", nginxConf),
	}
}

func (p *Proxy) WithService(
	service *Service,
) *Proxy {
	config := p.getConfig(service.Backend, service.Name, service.Frontend, service.IsTcp)
	configPath := fmt.Sprintf("/etc/nginx/stream.d/%s.conf", service.Name)
	if service.IsTcp {
		configPath = fmt.Sprintf("/etc/nginx/conf.d/%s.conf", service.Name)
	}

	if service.Exposed {
		p.ctr = p.ctr.
		WithNewFile(configPath, config).
		WithServiceBinding(service.Name, service.Service).
		WithExposedPort(service.Frontend)
	}

	return p
}

func (p *Proxy) Service() *dagger.Container {
	return p.ctr.WithDefaultArgs([]string{"nginx", "-g", "daemon off;"})
}

func (p *Proxy) getConfig(port int, name string, frontend int, isTcp bool) string {
	var result bytes.Buffer
	var config string

	if !isTcp {
		config = `
    server {
      listen {{ .frontend }};
      listen [::]:{{ .frontend }};
      proxy_pass {{ .name }}:{{ .port }};
    }
`
	} else {
		config = `
    server {
      listen {{ .frontend }};
      listen [::]:{{ .frontend }};
    
      server_name {{ .name }};
    
      location / {
        proxy_pass http://{{ .name }}:{{ .port }};
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
      }
    }
`
	}

	templ := template.Must(template.New("nginx").Parse(config))
	templ.Execute(&result, map[string]interface{}{
		"frontend": frontend,
		"name":     name,
		"port":     port,
	})

	return result.String()
}
