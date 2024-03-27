/**
 * A EsLint module service to lint and fix your code without local installation.
 * 
 * You can use it to lint any typescript files in your project.
 * It does not support custom configuration files.
 */

import {
  dag,
  Container,
  Directory,
  object,
  func,
  field,
} from "@dagger.io/dagger";

@object()
class Eslint {
  /**
   * The version of eslint to use (default to: 8.56.0).
   */
  @field()
  version = "8.56.0";

  /**
   * The version of node to use (default to: 21-alpine3.18).
   */
  @field()
  nodeVersion = "21-alpine3.18";

  /**
   * The files to lint.
   */
  @field()
  files: Directory;

  /*
   * Default configuration files .
   */
  defaultConfig = dag.currentModule().source().directory("src/default-config");

  constructor(files: Directory, version?: string, nodeVersion?: string) {
    this.version = version ?? this.version;
    this.nodeVersion = nodeVersion ?? this.nodeVersion;
    this.files = files ?? this.files;
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
      .withWorkdir("/app");

    // Check if there's an existing project configuration
    // and add missing files if they are missing
    const files = await ctr.directory(".").entries();
    if (!files.find((f) => f === "package.json")) {
      ctr = ctr.withFile(
        "package.json",
        this.defaultConfig.file("package.json")
      );
    }

    if (!files.find((f) => /\.eslintrc\.(js|json|cjs)$/.test(f))) {
      ctr = ctr.withFile(
        ".eslintrc.js",
        this.defaultConfig.file(".eslintrc.js")
      );
    }

    if (!files.find((f) => f === "tsconfig.json")) {
      ctr = ctr.withFile(
        "tsconfig.json",
        this.defaultConfig.file("tsconfig.json")
      );
    }

    return ctr
      .withExec(["npm", "install"])
      .withEntrypoint(["eslint"])
      .withDefaultTerminalCmd(["/bin/sh"]);
  }

  /**
   * Return the container with the linter executed in it.
   *
   * Example usage: `dagger --files=. call run stdout`
   */
  @func()
  async run(...args: string[]): Promise<Container> {
    return (await this.container()).withExec([".", ...args]);
  }

  /**
   * Lint the files and return the result
   *
   * Example usage: `dagger --files=. call lint`
   */
  @func()
  async lint(...args: string[]): Promise<string> {
    return (await this.run(...args)).stdout();
  }

  /**
   * Lint the files with auto fix and returns the result directory
   *
   * Example usage: `dagger --files=src call fix -o src`
   */
  @func()
  async fix(): Promise<Directory> {
    let resultCtr: Container;

    try {
      resultCtr = await (await this.run("--fix")).sync();
    } catch (e) {
      console.log("Error while running linter with fix options:", e);
    }

    return this.files.diff(
      dag.directory().withDirectory(".", resultCtr.directory("."), {
        include: ["**/*.ts"],
        exclude: ["node_modules", "dist"],
      })
    );
  }
}
