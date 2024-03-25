import { PrismaClient } from "@prisma/client";

interface Context {
    db: PrismaClient
}

const context: Partial<Context> = {}

export function init() {
    context.db = new PrismaClient()
}

export default context as Context