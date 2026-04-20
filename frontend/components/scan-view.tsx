"use client";

import Link from "next/link";
import { startTransition, useDeferredValue, useEffect, useState } from "react";
import { getScan, pdfUrl, streamScan } from "@/lib/api";
import { ModuleCard } from "@/components/module-card";
import { RadarChart } from "@/components/radar-chart";
import { defaultModules, type Scan } from "@/lib/types";

function formatDate(value?: string) {
  if (!value) {
    return "—";
  }
  return new Intl.DateTimeFormat("tr-TR", {
    dateStyle: "short",
    timeStyle: "medium"
  }).format(new Date(value));
}

export function ScanView({ scanId }: { scanId: string }) {
  const [scan, setScan] = useState<Scan | null>(null);
  const [streamState, setStreamState] = useState("connecting");
  const [error, setError] = useState<string | null>(null);
  const deferredScan = useDeferredValue(scan);

  useEffect(() => {
    let mounted = true;
    let eventSource: EventSource | null = null;
    let pollTimer: ReturnType<typeof setInterval> | undefined;

    async function load() {
      try {
        const initial = await getScan(scanId);
        if (!mounted) {
          return;
        }
        setScan(initial);
      } catch (cause) {
        if (!mounted) {
          return;
        }
        setError(cause instanceof Error ? cause.message : "Tarama okunamadı.");
      }
    }

    load().catch(() => undefined);

    eventSource = streamScan(
      scanId,
      (nextScan) => {
        startTransition(() => {
          setScan(nextScan);
        });
        setStreamState("live");
      },
      (message) => {
        setStreamState("fallback");
        setError(message);
      }
    );

    pollTimer = setInterval(async () => {
      try {
        const refreshed = await getScan(scanId);
        if (!mounted) {
          return;
        }
        setScan(refreshed);
      } catch {
        return;
      }
    }, 12000);

    return () => {
      mounted = false;
      if (pollTimer) {
        clearInterval(pollTimer);
      }
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [scanId]);

  const modules = deferredScan
    ? defaultModules.map((name) => [name, deferredScan.modules[name]] as const)
    : defaultModules.map((name) => [name, undefined] as const);

  return (
    <main className="page-shell">
      <section className="mx-auto flex min-h-screen w-full max-w-7xl flex-col gap-6 px-5 py-8 sm:px-8 lg:px-10">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.2em] text-signal">Scan Session</p>
            <h1 className="mt-2 text-3xl font-semibold text-white sm:text-5xl">{scan?.hostname ?? scanId}</h1>
            <p className="mt-3 max-w-3xl text-sm leading-6 text-slate-300">
              Hedef: <span className="text-white">{scan?.url ?? "yükleniyor"}</span>
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <Link href="/" className="rounded-full border border-white/10 px-4 py-2 text-sm text-slate-200 transition hover:border-white/20">
              Yeni tarama
            </Link>
            <a
              href={pdfUrl(scanId)}
              className="rounded-full bg-glow px-4 py-2 text-sm font-semibold text-slate-950 transition hover:brightness-105"
              target="_blank"
              rel="noreferrer"
            >
              PDF rapor
            </a>
          </div>
        </div>

        {error ? <p className="rounded-2xl border border-danger/40 bg-danger/10 px-4 py-3 text-sm text-danger">{error}</p> : null}

        <section className="grid gap-5 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="glass-panel rounded-[30px] p-6 sm:p-8">
            <div className="grid gap-4 sm:grid-cols-4">
              <MetricBox label="Durum" value={scan?.status ?? "queued"} tone="signal" />
              <MetricBox label="Skor" value={String(scan?.summary.score ?? 0)} tone="glow" />
              <MetricBox label="Not" value={scan?.summary.grade ?? "—"} tone="signal" />
              <MetricBox label="Risk" value={scan?.summary.riskLevel ?? "—"} tone="danger" />
            </div>

            <div className="mt-6 grid gap-4 sm:grid-cols-3">
              <div className="rounded-[24px] border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">Started</p>
                <p className="mt-2 text-sm text-white">{formatDate(scan?.startedAt)}</p>
              </div>
              <div className="rounded-[24px] border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">Completed</p>
                <p className="mt-2 text-sm text-white">{formatDate(scan?.completedAt)}</p>
              </div>
              <div className="rounded-[24px] border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">SSE</p>
                <p className="mt-2 text-sm text-white">{streamState}</p>
              </div>
            </div>

            <div className="mt-6 rounded-[28px] border border-white/10 bg-black/15 p-5">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">Event Feed</p>
              <div className="mt-4 max-h-64 space-y-3 overflow-auto pr-1">
                {(scan?.events ?? []).slice().reverse().map((event) => (
                  <div key={`${event.at}-${event.message}`} className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                      <span className="text-sm font-medium text-white">{event.message}</span>
                      <span className="text-xs uppercase tracking-[0.16em] text-slate-400">{event.status}</span>
                    </div>
                    <p className="mt-2 text-xs text-slate-400">{formatDate(event.at)}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {scan ? <RadarChart scan={scan} /> : <div className="glass-panel rounded-[30px] p-6">Yükleniyor...</div>}
        </section>

        <section className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          {modules.map(([name, result]) => (
            <ModuleCard key={name} moduleName={name} result={result} />
          ))}
        </section>
      </section>
    </main>
  );
}

function MetricBox({ label, value, tone }: { label: string; value: string; tone: "signal" | "glow" | "danger" }) {
  const toneClass =
    tone === "signal"
      ? "text-signal"
      : tone === "glow"
        ? "text-glow"
        : "text-danger";

  return (
    <div className="rounded-[24px] border border-white/10 bg-white/5 p-4">
      <p className="text-xs uppercase tracking-[0.18em] text-slate-400">{label}</p>
      <p className={`mt-3 text-3xl font-semibold ${toneClass}`}>{value}</p>
    </div>
  );
}
