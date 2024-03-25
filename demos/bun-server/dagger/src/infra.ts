import { dag, object, func, Service, Secret, Container } from "@dagger.io/dagger";

@object()
export class Infra {
  /**
   * Return a ready to use Redis instance exposed on the given port.
   *
   * @param port Port to expose the Redis instance on
   * 
   * Usage:
   * ```
   * dagger call infra redis --port 6379 as-service up
   * ```
   */
  @func()
  redis(port: number): Container {
    return dag.redis({ port }).server().withExposedPort(port);
  }

  /**
   * Return a ready to use database exposed on the given port.
   *
   * @param user Username for the database
   * @param pwd Password for the database
   * @param port Port to expose the database on
   *
   * Usage:
   * ```
   * dagger call infra postgres --user env:DB_USER --pwd env:DB_PASS --port 5432 as-service up
   * ```
   */
  @func()
  postgres(user: Secret, pwd: Secret, port: number): Container {
    return dag.postgres(user, pwd, port, { cache: false, dbName: "daggerdb" }).database().withExposedPort(port);
  }

  /**
   * Start up a dev infrastructure with a Redis and Postgres instance.
   * This is useful for local development and testing.
   * 
   * @param dbPass Database password
   * @param dbUser Database user
   * @param dbPort Database port
   * @param redisPort Redis port
   * 
   * Usage:
   * ```
   * dagger call infra dev --dbPass env:DB_PASS --dbUser env:DB_USER --dbPort $DB_PORT --redisPort $REDIS_PORT up --ports $DB_PORT:$DB_PORT,$REDIS_PORT:$REDIS_PORT 
   * ```
   */
  @func()
  dev(dbPass: Secret, dbUser: Secret, dbPort: number, redisPort: number): Service {
    const redis = this.redis(redisPort);
    const postgres = this.postgres(dbUser, dbPass, dbPort);

    return dag
      .proxy()
      .withService(redis.asService(), "redis", redisPort, redisPort)
      .withService(postgres.asService(), "postgres", dbPort, dbPort)
      .service();
  }
}
