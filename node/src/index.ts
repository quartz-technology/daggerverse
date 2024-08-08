import { dag, Container, Directory, object, func, enumType } from "@dagger.io/dagger";

const workdir = "/src";
const yarnCachePath = "/usr/local/share/.cache/yarn";
const npmCachePath = "/root/.npm";

@enumType()
class PackageManager {
  static readonly YARN: string = "YARN";
  static readonly NPM: string = "NPM";
}

@object()
class Node {
  @func()
  ctr: Container;

  constructor(source: Directory, packageManager?: string, version?: string) {
    if (!packageManager) {
      packageManager = PackageManager.NPM;
    }

    if (!version) {
      version = "20.9.0"; // LTS
    }

    this.ctr = dag
      .container()
      .from(`node:${version}`)
      .withWorkdir(workdir)
      .withMountedDirectory(workdir, source)
      .withMountedCache(`${workdir}/node_modules`, dag.cacheVolume(`node-${version}-typescript-node-modules`));

    switch (packageManager.toUpperCase()) {
      case PackageManager.YARN:
        this.ctr = this.ctr
          .withEntrypoint(["yarn"])
          .withMountedCache(yarnCachePath, dag.cacheVolume(`node-${version}-typescript-yarn-cache`));
        break;
      case PackageManager.NPM:
        this.ctr = this.ctr
          .withEntrypoint(["npm"])
          .withMountedCache(npmCachePath, dag.cacheVolume(`node-${version}-typescript-npm-cache`));
        break;
      default:
        throw new Error(`Unknown package manager ${packageManager}, only "npm" and "yarn" are supported`);
    }
  }

  @func()
  run(args: string[]): Container {
    return this.ctr.withExec(args, { useEntrypoint: true });
  }

  @func()
  async start(): Promise<string> {
    return this.run(["run", "start"]).stdout();
  }

  @func()
  async lint(): Promise<string> {
    return this.run(["run", "lint"]).stdout();
  }

  @func()
  async test(): Promise<string> {
    return this.run(["run", "test"]).stdout();
  }

  @func()
  build(): Node {
    this.ctr = this.run(["run", "build"]);

    return this;
  }

  @func()
  install(pkgs?: string[]): Node {
    if (!pkgs) {
      this.ctr = this.ctr.withExec(["install"], { useEntrypoint: true });
    } else {
      this.ctr = this.ctr.withExec(["install", ...pkgs], { useEntrypoint: true });
    }

    return this;
  }
}
