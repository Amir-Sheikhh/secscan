import express from "express";
import cookieParser from "cookie-parser";
import csurf from "csurf";

const app = express();
const enableCSRF = process.env.ENABLE_CSRF === "true";

app.use(express.urlencoded({ extended: true }));
app.use(cookieParser());

if (enableCSRF) {
  app.use(csurf({ cookie: true }));
}

app.get("/transfer", (req, res) => {
  const tokenField = enableCSRF ? `<input type="hidden" name="_csrf" value="${req.csrfToken()}">` : "";
  res.send(`
    <form action="/transfer" method="POST">
      ${tokenField}
      <input name="iban" placeholder="IBAN">
      <input name="amount" placeholder="Amount">
      <button type="submit">Transfer</button>
    </form>
  `);
});

app.post("/transfer", (req, res) => {
  res.send(`Transfer simulated: ${req.body.amount} -> ${req.body.iban}`);
});

app.listen(3002, () => {
  console.log("CSRF demo listening on http://localhost:3002");
});
