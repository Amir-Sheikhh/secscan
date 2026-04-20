import express from "express";

const app = express();

app.get("/unsafe", (req, res) => {
  const name = req.query.name || "guest";
  res.send(`<h1>Hello ${name}</h1>`);
});

app.get("/safe", (req, res) => {
  const name = String(req.query.name || "guest")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");

  res.setHeader("Content-Security-Policy", "default-src 'self'; script-src 'self'");
  res.send(`<h1>Hello ${name}</h1>`);
});

app.listen(3001, () => {
  console.log("XSS demo listening on http://localhost:3001");
});
