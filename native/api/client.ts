import { Fetcher, Middleware } from 'openapi-typescript-fetch'

import type { paths, definitions } from "../schema";

const logger: Middleware = async (url, init, next) => {
    console.log(`fetching ${url}`)
    const response = await next(url, init)
    console.log(`fetched ${url}`)
    return response
  }

const client = Fetcher.for<paths>()

client.configure({
    baseUrl: "http://localhost:3000",
    init: {
        headers: {
            "Content-Type": "application/json"
        },
    },
    use: [logger],
})

export { client, type definitions }