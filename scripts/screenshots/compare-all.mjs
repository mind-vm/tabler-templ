// Compares every rebuilt page with the Tabler original and lists the pages by how many pixels differ, most
// first: after a Tabler update, the list of pages to look at. The page body is compared as compare.mjs does;
// pages with a different height are listed as such. Screenshots of the pages that differ stay in the out dir
// for crop.mjs.
// Usage: CHROMIUM=... node compare-all.mjs <demo base URL> <original base URL> [out dir] [min pixels]
import { chromium } from "playwright";
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";

const [base, orig, out = "/tmp/tabler-compare", min = "0"] = process.argv.slice(2);
if (!base || !orig) {
  console.error("usage: node compare-all.mjs <demo base URL> <original base URL> [out dir] [min pixels]");
  process.exit(2);
}
await mkdir(out, { recursive: true });
const names = await (await fetch(`${base}/pages.json`)).json();
const browser = await chromium.launch(process.env.CHROMIUM ? { executablePath: process.env.CHROMIUM } : {});
const page = await browser.newPage({ viewport: { width: 1700, height: 900 } });
// Charts and count-ups animate; reduced motion stops what honours it.
await page.emulateMedia({ reducedMotion: "reduce" });

async function shoot(url) {
  try {
    const res = await page.goto(url, { waitUntil: "networkidle", timeout: 30000 });
    if (!res || !res.ok()) return null;
    await page.evaluate(() => document.fonts.ready);
    // Charts are drawn as they come into view; draw all of them for the picture.
    await page.evaluate(() => window.tablerCharts?.flush?.());
    await page.waitForTimeout(300);
    const body = page.locator(".page-body .container-xl");
    const el = (await body.count()) ? body.first() : page.locator("main").first();
    if (await el.isVisible()) return await el.screenshot({ timeout: 10000 });
    return await page.screenshot({ fullPage: true }); // e.g. a page that is only a loader or a dialog
  } catch (e) {
    console.error(`\n${url}: ${e.message.split("\n")[0]}`);
    return null;
  }
}

// diff counts the pixels that differ by more than a little, in a page of its own.
const differ = await browser.newPage();
async function diff(a, b) {
  return differ.evaluate(async ([a, b]) => {
    const load = (s) => new Promise((r) => { const i = new Image(); i.onload = () => r(i); i.src = "data:image/png;base64," + s; });
    const [ia, ib] = await Promise.all([load(a), load(b)]);
    if (ia.width !== ib.width || ia.height !== ib.height) return { size: `${ia.width}x${ia.height} vs ${ib.width}x${ib.height}` };
    const c = Object.assign(document.createElement("canvas"), { width: ia.width, height: ia.height }).getContext("2d");
    c.drawImage(ia, 0, 0); const pa = c.getImageData(0, 0, ia.width, ia.height).data;
    c.drawImage(ib, 0, 0); const pb = c.getImageData(0, 0, ia.width, ia.height).data;
    let n = 0;
    for (let i = 0; i < pa.length; i += 4) if (Math.abs(pa[i] - pb[i]) + Math.abs(pa[i + 1] - pb[i + 1]) + Math.abs(pa[i + 2] - pb[i + 2]) > 30) n++;
    return { pixels: n };
  }, [a.toString("base64"), b.toString("base64")]);
}

const rows = [];
for (const name of names) {
  const file = name.replaceAll("/", "-");
  const ours = await shoot(`${base}/${name}`);
  const theirs = await shoot(`${orig}/${name}.html`);
  let r;
  if (!ours || !theirs) r = { name, note: !ours ? "ours failed" : "original failed" };
  else r = { name, ...(await diff(ours, theirs)) };
  if (r.size || r.pixels > +min) {
    await writeFile(join(out, `${file}.ours.png`), ours);
    await writeFile(join(out, `${file}.tabler.png`), theirs);
  }
  rows.push(r);
  process.stderr.write(".");
}
await browser.close();
process.stderr.write("\n");

const rank = (r) => (r.note ? 3e9 : r.size ? 2e9 : r.pixels);
rows.sort((a, b) => rank(b) - rank(a));
const lines = rows.filter((r) => r.note || r.size || r.pixels > +min)
  .map((r) => `${r.name}\t${r.note ?? r.size ?? r.pixels}`);
await writeFile(join(out, "report.tsv"), lines.join("\n") + "\n");
console.log(lines.join("\n"));
console.log(`\n${lines.length} of ${rows.length} pages differ; screenshots and report.tsv in ${out}`);
