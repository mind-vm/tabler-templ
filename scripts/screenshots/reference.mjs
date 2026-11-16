// Screenshots of the reference app (examples/reference) at the desktop (1400 × 950) and on a phone
// (390 × 844), with its dialogs and the panel open, into docs/screenshots/reference/. It fails when a page
// scrolls sideways (GS-14), when the browser logs an error, or when a page loads anything from outside.
// Usage: node reference.mjs http://localhost:8080 ../../docs/screenshots/reference
// scripts/reference.sh starts the demo and runs it.
import { chromium } from "playwright";
import { mkdir } from "node:fs/promises";
import { join } from "node:path";

const [base = "http://localhost:8080", out = "../../docs/screenshots/reference"] = process.argv.slice(2);
const ref = base + "/reference";

// name, path, and what to do before the picture (open a dialog or the panel).
const shots = [
  ["overview", ""],
  ["b01-rahmen", "/rahmen"],
  ["b01-panel", "/rahmen", async (p) => { await p.click('[data-kit-sidepanel-tab="assistent"] >> visible=true'); await p.waitForTimeout(400); }],
  ["b02-start", "/start"],
  ["b03-liste", "/kunden"],
  ["b03-leere-suche", "/kunden?q=Zeppelin"],
  ["b04-detail", "/kunden/K-1001"],
  ["b04-dialog", "/kunden/K-1001", async (p) => { await p.click('[data-bs-target="#dialog-stammdaten"]'); await p.waitForSelector("#dialog-stammdaten.show"); await p.waitForTimeout(300); }],
  ["b04-dialog-fehler", "/kunden/K-1001", async (p) => {
    await p.click('[data-bs-target="#dialog-stammdaten"]'); await p.waitForSelector("#dialog-stammdaten.show");
    await p.fill("#stammdaten-name", ""); await p.click('#dialog-stammdaten .modal-footer .btn-primary');
    await p.waitForSelector("#dialog-stammdaten.show .kit-dialog-error"); await p.waitForTimeout(300);
  }],
  ["b05-angebot", "/angebote/A-2026-001"],
  ["b05-liste-bearbeiten", "/angebote/A-2026-001", async (p) => { await p.click('[data-bs-target="#dialog-positionen"]'); await p.waitForSelector("#dialog-positionen.show"); await p.waitForTimeout(300); }],
  ["b05-leere-liste", "/angebote/A-2026-041", async (p) => {
    await p.click('[data-bs-target="#dialog-positionen"]'); await p.waitForSelector("#dialog-positionen.show");
    // With only the new row the add button stays in the row, not in a dialog footer.
    if (!(await p.isVisible("#dialog-positionen form.kit-rowedit-new .kit-rowedit-add")) || (await p.locator("#dialog-positionen .kit-dialog-footer").count()))
      throw new Error("empty list: the add button left its row");
    await p.waitForTimeout(300);
  }],
  ["b05-schreibgeschuetzt", "/angebote/A-2026-002"],
  ["b06-anlegen", "/kunden/neu"],
  ["b07-zustaende", "/zustaende?meldung=1", async (p) => { await p.waitForSelector("[data-ref-load]:not([aria-busy])"); }],
  ["b07-grund", "/zustaende", async (p) => { await p.click('[data-bs-target="#dialog-zuruecknehmen"]'); await p.waitForSelector("#dialog-zuruecknehmen.show"); await p.waitForTimeout(300); }],
  ["b07-fehlerseite", "/gibt-es-nicht", null, 404],
  ["b08-marke", "/marke"],
  ["b09-assistent", "/assistent"],
  // The concepts are long: only the visible part (a "before" step means no full-page picture).
  ["konzept", "/konzept", async () => {}],
  ["konzept-assistent", "/konzept/assistent", async () => {}],
  ["konzept-regel", "/konzept#gs-41", async (p) => { await p.waitForTimeout(1000); }],
  // Other themes last: the choice is kept in the browser for the following pages of the context.
  ["theme-tabler-b05", "/angebote/A-2026-001?theme=tabler"],
  ["theme-rose-b02", "/start?theme=rose"],
  ["theme-midnight-b04", "/kunden/K-1001?theme=midnight"],
];

const sizes = [
  ["desktop", { viewport: { width: 1400, height: 950 }, deviceScaleFactor: 1 }],
  ["phone", { viewport: { width: 390, height: 844 }, deviceScaleFactor: 2, isMobile: true, hasTouch: true }],
];

await mkdir(out, { recursive: true });
const browser = await chromium.launch(process.env.CHROMIUM ? { executablePath: process.env.CHROMIUM } : {});
const problems = [];
for (const [size, options] of sizes) {
  const context = await browser.newContext(options);
  for (const [name, path, before, status = 200] of shots) {
    const page = await context.newPage();
    const origin = new URL(base).origin;
    page.on("console", (m) => { if (m.type() === "error" && !/status of (404|403|422)/.test(m.text())) problems.push(`${size} ${name}: console: ${m.text()}`); });
    page.on("pageerror", (e) => problems.push(`${size} ${name}: ${e.message}`));
    page.on("request", (r) => { if (!r.url().startsWith(origin) && !r.url().startsWith("data:")) problems.push(`${size} ${name}: loads from outside: ${r.url()}`); });
    // A new context per shot keeps the panel's cookie from one shot out of the next.
    await context.clearCookies();
    const res = await page.goto(ref + path, { waitUntil: "networkidle" });
    if (res.status() !== status) problems.push(`${size} ${name}: status ${res.status()}, want ${status}`);
    await page.evaluate(() => document.fonts.ready);
    if (before) await before(page);
    // GS-14: the page never scrolls sideways; wide tables scroll inside their island.
    const wide = await page.evaluate(() => {
      const d = document.documentElement;
      return d.scrollWidth > d.clientWidth + 1 ? `${d.scrollWidth} > ${d.clientWidth}` : "";
    });
    if (wide) problems.push(`${size} ${name}: scrolls sideways (${wide})`);
    await page.screenshot({ path: join(out, `${name}-${size}.png`), fullPage: !before });
    console.log(`${size} ${name}${wide ? " WIDE" : ""}`);
    await page.close();
  }
  await context.close();
}
await browser.close();
if (problems.length) {
  console.error("\n" + problems.join("\n"));
  process.exit(1);
}
