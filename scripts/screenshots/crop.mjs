// Shows the same band of two screenshots side by side, to look at what diff.mjs found.
// Usage: node crop.mjs <a.png> <b.png> <y> <height> <out.png>
import { chromium } from "playwright";
import fs from 'node:fs';
const [a,b,y,h,out]=process.argv.slice(2);
const browser = await chromium.launch({executablePath: process.env.CHROMIUM});
const page = await browser.newPage({viewport:{width:2660,height:+h}});
const da = fs.readFileSync(a).toString('base64'), db = fs.readFileSync(b).toString('base64');
await page.setContent(`<body style="margin:0;display:flex;gap:20px"><div style="height:${h}px;width:1320px;overflow:hidden"><img style="margin-top:-${y}px" src="data:image/png;base64,${da}"></div><div style="height:${h}px;width:1320px;overflow:hidden"><img style="margin-top:-${y}px" src="data:image/png;base64,${db}"></div></body>`);
await page.screenshot({path: out});
await browser.close();
