import { Bot } from "../deps.deno.ts";

export const bot = new Bot(Deno.env.get("TELEGRAM_TOKEN") || "");

const sanitizeString = (val: string) => {
  return val.replace(
    /(\[[^\][]*]\(http[^()]*\))|[_*[\]()~>#+=|{}.!-]/gi,
    (x, y) => (y ? y : "\\" + x)
  );
}

bot.command("start", (ctx) => ctx.reply("Welcome! Up and running."));

bot.command("me", async (ctx) => {
  try {
    const autor = await ctx.getAuthor();

    const response = `*Your Data*

*Name:* ${autor.user.first_name} ${autor.user.last_name}
*Username:* \`${autor.user.username}\`
*UserID:* \`${autor.user.id}\`
*Status:* ${autor.status}

Now: ${sanitizeString(new Date().toLocaleString())}`;
    await ctx.reply(response, {
      reply_to_message_id: ctx.msg.message_id,
      parse_mode: "MarkdownV2"
    });
  } catch (err) {
    ctx.reply(err.message, {
      reply_to_message_id: ctx.msg.message_id,
    });
  }
});
