import { dag, Container, Directory, object, func, field } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
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
  website: Container;

  /**
   * Port the website is exposed on.
   */
  @field()
  port: number = 8080;

  constructor(source: Directory, website: Container, port?: number) {
    this.source = source;
    this.website = website;
    this.port = port ?? this.port
  }

  /**
   * Run e2e tests on cypress.
   * 
   * Example usage: `dagger call run`
   */
  @func()
  async run(): Promise<string> {
    return this.base().stdout();
  }


  @func()
  base(): Container {
    return dag
      .container()
      .from('cypress/included:13.6.2')
      .withServiceBinding('website', this.website.asService())
      .withDirectory('/e2e', this.source)
      .withWorkdir('/e2e')
      .withExec(['npm', 'install'], { skipEntrypoint: true })
      .withExec(['--env', `BASE_URL=http://website:${this.port}`]);
  }
}
