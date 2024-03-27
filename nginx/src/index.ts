/**
 * A simple NGINX module service to expose a website using the NGINX server.
 *
 * This module is still a work in progress and is not yet ready for production.
 * You can use it to test your website locally but also run cypress tests in a dagger CI.
 * 
 * The module isn't designed for complex uses cases and is limited to a simple website.
 */

import { dag, Directory, object, func, field, Service, File } from "@dagger.io/dagger";

const defaultConfig = `
server {
  listen 8080;
  listen [::]:8080;

  root /app/website;
  
  location / {
    try_files $uri /index.html;  
  }
}
`;

@object()
class Nginx {
  /**
   * Nginx version to use (default to 1.25.3)
   */
  @field()
  version: string = "1.25.3";

  /**
   * HTML source files
   */
  @field()
  source: Directory;

  /**
   * Nginx configuration that overrides the default
   */
  @field()
  config: File;

  /**
   * Port to expose the nginx on
   */
  @field()
  port: number = 8080;

  constructor(source: Directory, port?: number, version?: string, config?: File) {
    this.source = source;
    this.version = version ?? this.version;
    this.port = port ?? this.port;
    this.config = config ?? dag.directory().withNewFile("default.conf", defaultConfig).file("default.conf");
  }

  /**
   * Expose the nginx server
   */
  @func()
  expose(): Service {
    return dag
      .container()
      .from(`nginx:${this.version}`)
      .withMountedDirectory("/app/website", this.source)
      .withMountedFile("/etc/nginx/conf.d/default.conf", this.config)
      .withExposedPort(this.port)
      .asService();
  }
}
