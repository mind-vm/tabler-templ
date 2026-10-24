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
// A line of pages.txt may give a viewport width after the name, e.g. "themes 1600".
const widths = Object.fromEntries(names.map((n) => n.split(/\s+/)).filter((p) => p[1]).map(([n, w]) => [n, +w]));
await mkdir(out, { recursive: true });

const browser = await chromium.launch(process.env.CHROMIUM ? { executablePath: process.env.CHROMIUM } : {});
const page = await browser.newPage({ viewport: { width: 1280, height: 800 }, deviceScaleFactor: 1 });
for (const line of ["index", ...names]) {
  const name = line.split(/\s+/)[0];
  await page.setViewportSize({ width: widths[name] ?? 1280, height: 800 });
  await page.goto(name === "index" ? base : `${base}/${name}`, { waitUntil: "networkidle" });
  await page.evaluate(() => document.fonts.ready);
  // Charts are drawn as they come into view; draw all of them for the picture.
  await page.evaluate(() => window.tablerCharts?.flush?.());
  await page.screenshot({ path: join(out, `${name}.png`), fullPage: true });
  console.log(name);
}
await browser.close();
