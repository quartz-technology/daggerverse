import { dag, Container, Directory, object, func, field, Service, Platform } from "@dagger.io/dagger";
import { Infra } from "./infra";

/**
 * A simple tier app CI example for a bun application that use a Redis and a PostgreSQL database.
 * 
 * All the local and CI infrastructure is defined using Dagger, nothing will be run in the user localhost
 * system. 
 */
@object()
class TierAppCi {
  /**
   * We de not expose the container so we don't use `@field()` decorator
   */
  ctr: Container;

  /**
   * Access specific infrastructure commands
   */
  @func()
  infra(): Infra {
    return new Infra();
  }

  /**
   * Execute build, lint, and test in concurrency
   * 
   * Usage:
   * ```
   * dagger call app --source=. ci
   * ```
   */
  @func()
  async ci(): Promise<string> {
    await Promise.all([await this.build().sync(), await this.lint().stdout(), await this.test()]);

    return "All checks passed!";
  }

  /**
   * Return the built Bun application
   *
   * Usage:
   * ```
   * dagger call app --source=. build file --path api -o ./api
   * ```
   */
  @func()
  build(): Container {
    return this.ctr
      .pipeline("build")
      .withFocus()
      .withExec(["bun", "build", "index.ts", "--compile", "--outfile", "/api"]);
  }

  /**
   * Execute lint on the app
   *
   * Usage:
   * ```
   * dagger call app --source=. lint stdout
   * ```
   */
  @func()
  lint(): Container {
    return this.ctr.pipeline("lint").withFocus().withExec(["bun", "lint"]);
  }

  /**
   * Runs the integration tests of the application in an isolate environment
   *
   * Usage:
   * ```
   * dagger call app --source=. test
   * ```
   */
  @func()
  async test(): Promise<string> {
    return this.stack().withExec(["bun", "test", "--bail"]).withFocus().stdout();
  }

  /**
   * Create the development stack for the Bun application
   * 
   * This provide a Postgres and Redis service for the application and returns
   * a container ready to be use.
   * 
   * You can open a Shell in that container:
   * ```
   * dagger call app --source=. stack terminal
   * ```
   */
  @func()
  stack(): Container {
    const dbPass = dag.setSecret("dbPass", "XkfhrPGef");
    const dbUser = dag.setSecret("dbUser", "dagger");
    const dbPort = 5432;
    const redisPort = 6379;
    const sessionSecret = dag.setSecret("sessionSecret", "dFPGejege");

    const db = this.infra().postgres(dbUser, dbPass, dbPort);
    const redis = this.infra().redis(redisPort);

    // Add environment configuration and binds services to the Bun application
    return this.ctr
      .withSecretVariable("DB_PASS", dbPass)
      .withSecretVariable("DB_USER", dbUser)
      .withEnvVariable("DB_PORT", dbPort.toString())
      .withEnvVariable("DB_HOST", "postgres")
      .withEnvVariable("DB_NAME", "daggerdb")
      .withEnvVariable("DATABASE_URL", "postgresql://dagger:XkfhrPGef@postgres:5432/daggerdb?schema=public")
      .withEnvVariable("REDIS_PORT", redisPort.toString())
      .withEnvVariable("REDIS_HOST", "redis")
      .withSecretVariable("SESSION_SECRET", sessionSecret)
      .withServiceBinding("postgres", db.asService())
      .withServiceBinding("redis", redis.asService())
      .withEnvVariable("CACHE_BUSTER", new Date().toString()) // Force migration on every run
      .withExec(["bunx", "prisma", "migrate", "dev", "--skip-generate"])
      .withExec(["bunx", "prisma", "generate"]);
  }

  /**
   * Start the application on port 8080 with the development stack
   * 
   * Usage:
   * ```
   * dagger call app --source=. run up
   * ```
   */
  @func()
  run(): Service {
    return this.stack().withExec(["bun", "run", "index.ts"]).withExposedPort(8080).asService();
  }

  /**
   * Provide the source code and a Bun container with dependencies installed
   *
   * @param source The directory to build the Bun app
   */
  @func()
  app(source: Directory): TierAppCi {
    this.ctr = dag
      .container({ platform: "linux/amd64" as Platform }) // force amd64 to avoid platform issue with prisma
      .from("oven/bun:1") // FROM
      .withWorkdir("/app") // WORKDIR
      .withDirectory("/app", source) // COPY
      .withExec(["bun", "install"]); // RUN

    return this;
  }
}
