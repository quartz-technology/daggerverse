import { Router } from "express";
import bodyParser from "body-parser";
import { validateRequest } from "zod-express-middleware";
import { z } from "zod";

import ctx from "@context";
import httpStatus from "http-status";

const router = Router();

router.get("/", async (_, res) => {
  try {
    const users = await ctx.db.user.findMany();

    return res.status(httpStatus.OK).json(users);
  } catch (error) {
    throw new Error(`Error fetching users: ${error}`);
  }
});

router.get("/:id", async (req, res) => {
  try {
    const { id } = req.params;
    const user = await ctx.db.user.findUnique({ where: { id } });
    if (!user) {
      throw new Error(`User ${id} not found`);
    }

    return res.status(httpStatus.OK).json(user);
  } catch (error) {
    throw new Error(`Error fetching user ${req.params.id}: ${error}`);
  }
});

router.post(
  "/",
  bodyParser.json(),
  validateRequest({
    body: z.object({
      email: z.string().email(),
    }),
  }),
  async (req, res) => {
    try {
      const { email } = req.body;
      const user = await ctx.db.user.create({ data: { email } });

      return res.status(httpStatus.CREATED).json(user);
    } catch (error) {
      console.log(error)
      throw new Error(`Error creating user: ${error}`);
    }
  }
);

router.put(
  "/:id",
  bodyParser.json(),
  validateRequest({
    body: z.object({
      email: z.string().email().optional(),
    }),
  }),
  async (req, res) => {
    try {
      const { id } = req.params;
      const { email } = req.body;
      const user = await ctx.db.user.update({ where: { id }, data: { email } });

      return res.status(httpStatus.OK).json(user);
    } catch (error) {
      throw new Error("Error updating user");
    }
  }
);

router.delete("/:id", async (req, res) => {
    try {
        const { id } = req.params;
        await ctx.db.user.delete({ where: { id } });

        return res.status(httpStatus.OK).json({ message: "User successfully deleted" });
    } catch (error) {
        throw new Error("Error deleting user");
    }
})

export default router;
