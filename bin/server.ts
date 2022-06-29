import { serve } from "https://deno.land/x/sift@0.5.0/mod.ts";
import { bot } from "../src/bot.ts";
import { webhookCallback } from "../deps.deno.ts";

const handleUpdate = webhookCallback(bot, "std/http");

serve({
  ["/" + Deno.env.get("TELEGRAM_TOKEN")]: async (req) => {
    if (req.method == "POST") {
      try {
        return await handleUpdate(req);
      } catch (err) {
        console.error(err);
      }
    }
    return new Response();
  },
  "/register": () => {
    return new Response("Hello world!");
  },
});
