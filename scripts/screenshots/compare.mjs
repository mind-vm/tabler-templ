// Screenshots the page body of a rebuilt page next to the Tabler original, to check a rebuild by eye.
// Usage: node compare.mjs <demo base URL> <name> [out dir] [original base URL]
// The original defaults to https://preview.tabler.io; a local build of tabler/preview works as well.
import { chromium } from "playwright";
import { join } from "node:path";

const [base, name, out = "/tmp", orig = "https://preview.tabler.io"] = process.argv.slice(2);
const browser = await chromium.launch(process.env.CHROMIUM ? { executablePath: process.env.CHROMIUM } : {});
// Wide enough that .container-xl reaches its full width next to the original's sidebar as well.
const page = await browser.newPage({ viewport: { width: 1700, height: 900 } });
for (const [label, url] of [["ours", `${base}/${name}`], ["tabler", `${orig}/${name}.html`]]) {
  await page.goto(url, { waitUntil: "networkidle" });
  await page.evaluate(() => document.fonts.ready);
  const path = join(out, `${name.replaceAll("/", "-")}.${label}.png`);
  await page.locator(".page-body .container-xl, .page-center .container").first().screenshot({ path });
  console.log(path);
}
await browser.close();
