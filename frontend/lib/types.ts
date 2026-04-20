export const defaultModules = ["ports", "headers", "tls", "fuzz", "xss", "sqli", "cve"] as const;

export type ModuleName = (typeof defaultModules)[number];

export const moduleLabels: Record<ModuleName, string> = {
  ports: "Port Scan",
  headers: "Security Headers",
  tls: "TLS Audit",
  fuzz: "Directory Fuzzer",
  xss: "XSS Reflection",
  sqli: "SQLi Heuristics",
  cve: "CVE Intel"
};

export interface Finding {
  title: string;
  severity: string;
  category: string;
  description: string;
  recommendation: string;
  evidence?: string;
}

export interface ModuleResult {
  name: string;
  status: "pending" | "running" | "completed" | "failed" | "skipped";
  startedAt?: string;
  completedAt?: string;
  durationMs: number;
  score: number;
  severity: string;
  summary: string;
  findings: Finding[];
  details?: Record<string, unknown>;
  error?: string;
}

export interface ScanEvent {
  type: string;
  module?: string;
  status: string;
  message: string;
  at: string;
}

export interface ScanSummary {
  score: number;
  grade: string;
  riskLevel: string;
  findings: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  passed: number;
  failed: number;
  moduleRuns: number;
  durationMs: number;
}

export interface Scan {
  id: string;
  url: string;
  hostname: string;
  resolvedIp: string;
  status: "queued" | "running" | "completed" | "failed";
  createdAt: string;
  startedAt?: string;
  completedAt?: string;
  modules: Record<string, ModuleResult>;
  events: ScanEvent[];
  summary: ScanSummary;
}
