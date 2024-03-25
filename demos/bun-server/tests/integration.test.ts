import { test, expect } from "bun:test";
import supertest from "supertest";

import { createApp } from "@server";

const app = createApp();

test("GET /healthcheck", async () => {
  await supertest(app).get("/healthcheck").expect(200).expect({ status: "ok" });
});

test("User flow", async () => {
  for (let i = 0; i < 10; i++) {
    await supertest(app)
      .post("/users")
      .send({ email: `${i}-@email.com` })
      .expect(201);
  }

  const users = await supertest(app).get("/users");
  expect(users.body.length).toBe(10);
});
