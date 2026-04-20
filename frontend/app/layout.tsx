import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SecScan",
  description: "Web security scanner dashboard for Finale Dogru / Seri 3"
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="tr">
      <body>{children}</body>
    </html>
  );
}
