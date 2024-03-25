import dotenv from "dotenv";

// Load .env
dotenv.config();

interface Config {
  redis: {
    host: string;
    port: number;
  };
  server: {
    port: number;
    sessionSecret: string;
  };
}

function getEnv(key: string, defaultValue?: string): string {
  const value = process.env[key] || defaultValue;
  if (!value) {
    throw new Error(`Environment variable ${key} is not set`);
  }

  return value;
}

const config: Config = {
  redis: {
    host: getEnv("REDIS_HOST", "localhost"),
    port: parseInt(getEnv("REDIS_PORT", "6379")),
  },
  server: {
    port: parseInt(getEnv("SERVER_PORT", "8080")),
    sessionSecret: getEnv("SESSION_SECRET"),
  },
};

export default config;
