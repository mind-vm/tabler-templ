// Takes screenshots of a running demo server (go run ./cmd/demo) and writes docs/screenshots/<name>.png:
// the pages listed in docs/screenshots/pages.txt (the README's gallery), or every rebuilt page with --all.
// Usage: node shoot.mjs http://localhost:8080 ../../docs/screenshots [--all]
import { chromium } from "playwright";
import { mkdir } from "node:fs/promises";
import { join } from "node:path";

const [base = "http://localhost:8080", out = "../../docs/screenshots"] = process.argv.slice(2);
import { readFile } from "node:fs/promises";
const all = process.argv.includes("--all");
const names = all
  ? await (await fetch(`${base}/pages.json`)).json()
  : (await readFile(join(out, "pages.txt"), "utf8")).split("\n").map((l) => l.trim()).filter((l) => l && !l.startsWith("#") && l !== "index");
await mkdir(out, { recursive: true });

const browser = await chromium.launch(process.env.CHROMIUM ? { executablePath: process.env.CHROMIUM } : {});
const page = await browser.newPage({ viewport: { width: 1280, height: 800 }, deviceScaleFactor: 1 });
for (const name of ["index", ...names]) {
  await page.goto(name === "index" ? base : `${base}/${name}`, { waitUntil: "networkidle" });
  await page.evaluate(() => document.fonts.ready);
  await page.screenshot({ path: join(out, `${name}.png`), fullPage: true });
  console.log(name);
}
await browser.close();
