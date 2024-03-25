import { Router } from "express";

import healthcheck from "./healthcheck";
import user from "./user";

const router = Router();

router.use("/users", user)
router.use("/healthcheck", healthcheck);

export default router;
