/**
 * A Cypress module service to run e2e tests without local installation.
 * 
 * The module takes your cypress tests and run them in a container that point to your website.
 * It's super useful to run e2e tests in a dagger CI.
 * 
 * Warning: The module is still a work in progress and may lack of features.
 */

import { dag, Directory, object, func, field, Service } from '@dagger.io/dagger';

@object()
class Cypress {
  /**
   * The directory containing the cypress tests.
   */
  @field()
  source: Directory;

  /**
   * The website contained in a container pre configured to expose as a service
   * the website.
   */
  @field()
  website?: Service;

  /**
   * Port the website is exposed on.
   */
  @field()
  port: number = 8080;

  constructor(source: Directory, website?: Service, port?: number) {
    this.source = source;
    this.website = website;
    this.port = port ?? this.port;
  }

  /**
   * Run e2e tests on cypress.
   *
   * Note: The end to end test should do test using `BASE_URL` environment
   * to point to the website service.
   * If no website is provided, create a new one (only work for node application).
   *
   * TODO(TomChv): Take an interface Builder with a method `Build` to handle that for any language.
   * TODO(TomChv): Add param to run specific tests.
   * TODO(TomChv): Add param to not only run e2e tests.
   *
   * Example usage: `dagger call run`
   */
  @func()
  async run(): Promise<string> {
    // Build the website if it's not provided by the user
    if (!this.website) {
      const distSrc = dag
        .node()
        .withNpm()
        .withSource(this.source)
        .install()
        .commands()
        .build()
        .directory('dist');

      this.website = dag.nginx(distSrc).expose();
    }

    return dag
      .container()
      .from('cypress/included:13.6.2')
      .withServiceBinding('website', this.website)
      .withDirectory('/e2e', this.source)
      .withWorkdir('/e2e')
      .withExec(['npm', 'install'], { skipEntrypoint: true })
      // TODO(TomChv): Find a way to make it generic (right now it's assume cypress fetch the app URL from that env var)
      .withExec(['--env', `BASE_URL=http://website:${this.port}`])
      .stdout();
  }
}
