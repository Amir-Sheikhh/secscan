import mysql from "mysql2/promise";

const pool = mysql.createPool({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME
});

export async function secureLogin(email, password) {
  const query = "SELECT id, email FROM users WHERE email = ? AND password = SHA2(?, 256) LIMIT 1";
  const [rows] = await pool.execute(query, [email, password]);
  return rows;
}
