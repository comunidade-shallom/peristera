// run the bot locally
import { bot } from "../src/bot.ts";

await bot.api.deleteWebhook();

bot.start();
