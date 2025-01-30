package proxy

import (
	"bytes"
	"fmt"
	"text/template"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
)

// Service represents a service to be proxied.
type Service struct {
	// Service refers to a dagger.Service that this proxy will forward traffic to.
	Service *dagger.Service

	// Name is the unique name of the service.
	Name string

	// Alias is an alternative name for the service.
	Alias string

	// Frontend is the port number exposed by the proxy.
	Frontend int

	// Backend is the port number where the actual service is running.
	Backend int

	// IsTcp specifies whether the service uses the TCP protocol.
	IsTcp bool

	// Exposed signals if the service should be exposed by the proxy.
	Exposed bool
}

// Proxy is a struct that sets up a reverse proxy using an Nginx container.
type Proxy struct {
	// ctr is the internal container running the Nginx proxy.
	ctr *dagger.Container
}

// New initializes a Proxy with default Nginx configuration.
func New() *Proxy {
	return &Proxy{
		ctr: dag.Container().
			From("nginx:1.25.3").
			WithNewFile("/etc/nginx/stream.conf", streamConf).
			WithNewFile("/etc/nginx/nginx.conf", nginxConf),
	}
}

// WithService configures adds the given service to the proxy.
//
// If the service isn't exposed, it will not be added to the proxy.
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

// Service returns the configured proxy container ready to start services.
func (p *Proxy) Service() *dagger.Container {
	return p.ctr.WithDefaultArgs([]string{"nginx", "-g", "daemon off;"})
}

// getConfig generates the Nginx configuration for a service.
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
