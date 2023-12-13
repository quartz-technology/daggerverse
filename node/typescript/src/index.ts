import { dag, Container, Directory, object, func, field, Secret } from '@dagger.io/dagger';

const yarnCachePath = "/usr/local/share/.cache/yarn"
const npmCachePath = "/root/.npm"

@object
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Node {
  @field
  ctr?: Container = undefined

  /**
   * WithVersion returns Node container with given image version.
   *
   * @param version version to use
   */
  @func
  withVersion(version: string): Node {
    if (!this.ctr) {
      this.ctr = dag.container()
    }

    this.ctr = this.ctr
      .from(`node:${version}`)
      .withEntrypoint(["node"])

    return this
  }

  /**
   * WithContainer returns Node container with the given container.
   *
   * @param ctr container to use in Node.
   */
  @func
  withContainer(ctr: Container): Node {
    this.ctr = ctr

    return this
  }

  /**
   * Container returns Node container.
   */
  @func
  container(): Container {
    return this.ctr
  }

  /**
   * WithSource returns the Node container with source and cache set in it.
   *
   * @param source source code to add to the Node container
   */
  @func
  withSource(source: Directory): Node {
    const workdir = "/src"

    this.ctr = this.ctr
      .withWorkdir(workdir)
      .withMountedDirectory(workdir, source)
      .withMountedCache(`${workdir}/node_modules`,
        dag.cacheVolume("node-typescript-node-modules"))

    return this
  }

  /**
   * WithYarn returns Node container with yarn configured as package manager.
   */
  @func
  withYarn(): Node {
    this.ctr = this.ctr
      .withEntrypoint(["yarn"])
      .withMountedCache(yarnCachePath,
        dag.cacheVolume("node-typescript-yarn-cache")
      )

    return this
  }

  /**
   * WithNpm returns Node container with npm configured as package manager.
   */
  @func
  withNpm(): Node {
    this.ctr = this.ctr
      .withEntrypoint(["npm"])
      .withMountedCache(npmCachePath,
        dag.cacheVolume("node-typescript-npm-cache")
      )

    return this
  }

  /**
   * Run returns the container with the command set.
   *
   * @param args arguments to execute
   */
  @func
  run(args: string[]): Container {
    return this.ctr.withExec(args)
  }

  /**
   * Start returns the output of the code executed.
   */
  @func
  async start(): Promise<string> {
    return this.run(["run", "start"]).stdout()
  }

  /**
   * Lint returns the output of the lint command.
   */
  @func
  async lint(): Promise<string> {
    return this.run(["run", "lint"]).stdout()
  }

  /**
   * Lint returns the output of the lint command.
   */
  @func
  async test(): Promise<string> {
    return this.run(["run", "test"]).stdout()
  }

  /**
   * Build returns the Container with the source built.
   */
  @func
  build(): Node {
    this.ctr = this.run(["run", "build"])

    return this
  }

  /**
   * Publish publishes the source code to npm registry.
   *
   * @param tokenSecret secret token to register to npm registry
   * @param version version of the package
   * @param access access permission of the package
   * @param dryRun whether to do a dry run instead of an actual publish
   */
  @func
  async publish(tokenSecret: Secret, version: string, access?: string, dryRun = false): Promise<string> {
    const token = await tokenSecret.plaintext()
    const npmrc = `//registry.npmjs.org/:_authToken=${token}
registry=https://registry.npmjs.org/
always-auth=true`
    const publishCmd = ["publish"]

    if (dryRun) {
      publishCmd.push("--dry-run")
    }

    if (access) {
      publishCmd.push("--access", access)
    }

    return this.ctr
      .withNewFile(".npmrc", { contents: npmrc, permissions: 0o600 })
      .withExec(["version", "--new-version", version])
      .withExec(publishCmd)
      .stdout()
  }

  /**
   * Install adds given package.
   *
   * @param pkgs packages to additionally install
   */
  @func
  install(pkgs: string[]): Node {
    return this.withContainer(this.run(["install", ...pkgs]))
  }
}
