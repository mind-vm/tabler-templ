// Counts the pixels that differ between two screenshots (e.g. from compare.mjs) and the 50px bands they lie in.
// Usage: CHROMIUM=... node diff.mjs <a.png> <b.png>
import { chromium } from "playwright";
const [a,b,out]=process.argv.slice(2);
const browser = await chromium.launch({executablePath: process.env.CHROMIUM});
const page = await browser.newPage();
const fs = await import('node:fs');
const da = fs.readFileSync(a).toString('base64'), db = fs.readFileSync(b).toString('base64');
const r = await page.evaluate(async ([da,db]) => {
  const load = s => new Promise(r => { const i = new Image(); i.onload = () => r(i); i.src = 'data:image/png;base64,' + s; });
  const [ia, ib] = await Promise.all([load(da), load(db)]);
  const w = Math.min(ia.width, ib.width), h = Math.min(ia.height, ib.height);
  const c = document.createElement('canvas'); c.width = w; c.height = h; const x = c.getContext('2d');
  x.drawImage(ia,0,0); const pa = x.getImageData(0,0,w,h).data; x.clearRect(0,0,w,h); x.drawImage(ib,0,0); const pb = x.getImageData(0,0,w,h).data;
  let diff=0; const rows=new Set(); for (let i=0;i<pa.length;i+=4){ if (Math.abs(pa[i]-pb[i])+Math.abs(pa[i+1]-pb[i+1])+Math.abs(pa[i+2]-pb[i+2])>30){diff++; rows.add(Math.floor(i/4/w/50)*50);} }
  return {sizeA:[ia.width,ia.height], sizeB:[ib.width,ib.height], diff, rows:[...rows].slice(0,40)};
}, [da,db]);
console.log(JSON.stringify(r));
await browser.close();
