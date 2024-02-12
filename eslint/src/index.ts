import { dag, Container, Directory, object, func, field } from '@dagger.io/dagger';

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Eslint {
  /**
   * The version of eslint to use (default to: 8.56.0).
   */
  @field()
  version = "8.56.0"

  /**
   * The version of node to use (default to: 21-alpine3.18).
   */
  @field()
  nodeVersion = "21-alpine3.18"

  /**
   * The files to lint.
   */
  @field()
  files: Directory

  /*
   * Default configuration files .
   */
  defaultConfig = dag.currentModule().source().directory("src/default-config")

  constructor(files: Directory, version?: string, nodeVersion?: string) {
    this.version = version ?? this.version
    this.nodeVersion = nodeVersion ?? this.nodeVersion
    this.files = files ?? this.files
  }

  /**
   * Return a container with eslint installed in it.
   * 
   * Example usage: `dagger --files=. call container terminal`
   */
  @func()
  async container(): Promise<Container> {    
    let ctr = dag
      .container()
      .from(`node:${this.nodeVersion}`)
      .withExec(["npm", "install", "-g", `eslint@${this.version}`])
      .withMountedDirectory("/app", this.files)
      .withWorkdir("/app")

    // Check if there's an existing project configuration
    // and add missing files if they are missing
    const files = await ctr.directory(".").entries()
    if (!files.find((f) => f === "package.json")) {
      ctr = ctr.withFile("package.json", this.defaultConfig.file("package.json"))
    }

    if (!files.find((f) => /\.eslintrc\.(js|json|cjs)$/.test(f))) {
      ctr = ctr.withFile(".eslintrc.js", this.defaultConfig.file(".eslintrc.js"))
    }

    if (!files.find((f) => f === "tsconfig.json")) {
      ctr = ctr.withFile("tsconfig.json", this.defaultConfig.file("tsconfig.json"))
    }

    return ctr
      .withExec(["npm", "install"])
      .withEntrypoint(["eslint"])
      .withDefaultTerminalCmd(["/bin/sh"])
  }

  /**
   * Lint the files.
   * 
   * Example usage: `dagger --files=. call run stdout`
   */
  @func()
  async run(...args: string[]): Promise<Container> {
    return (await this.container())
      .withExec([".", ...args])
  }
}
