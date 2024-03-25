import express from "express";

import { init as initCtx } from "@context";
import middlewares from "@middlewares";
import endpoints from "@endpoints";

export function createApp(): express.Application {
  // Init my context application
  initCtx();

  const app = express();

  app.use("/", endpoints);

  app.use(middlewares.promiseCatcher);

  return app;
}
