import express from "express";
import helmet from "helmet";

const app = express();

app.use(
  helmet({
    hsts: {
      maxAge: 31536000,
      includeSubDomains: true,
      preload: true
    },
    contentSecurityPolicy: {
      directives: {
        defaultSrc: ["'self'"],
        scriptSrc: ["'self'"],
        objectSrc: ["'none'"],
        upgradeInsecureRequests: []
      }
    }
  })
);

app.get("/", (_req, res) => {
  res.send("hello world");
});

app.listen(3003, () => {
  console.log("Helmet demo listening on http://localhost:3003");
});
