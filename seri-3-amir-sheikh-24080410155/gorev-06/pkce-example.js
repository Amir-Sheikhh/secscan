import crypto from "crypto";

function base64url(buffer) {
  return buffer
    .toString("base64")
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replaceAll("=", "");
}

export function createPKCEPair() {
  const verifier = base64url(crypto.randomBytes(32));
  const challenge = base64url(crypto.createHash("sha256").update(verifier).digest());

  return {
    verifier,
    challenge,
    method: "S256"
  };
}

console.log(createPKCEPair());
