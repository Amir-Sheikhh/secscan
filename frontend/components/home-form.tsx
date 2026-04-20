"use client";

import { FormEvent, startTransition, useState } from "react";
import { useRouter } from "next/navigation";
import { defaultModules, moduleLabels, type ModuleName } from "@/lib/types";
import { startScan } from "@/lib/api";

export function HomeForm() {
  const router = useRouter();
  const [targetUrl, setTargetUrl] = useState("https://example.com");
  const [selected, setSelected] = useState<ModuleName[]>([...defaultModules]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const scan = await startScan({ url: targetUrl, modules: selected });
      startTransition(() => {
        router.push(`/scan/${scan.id}`);
      });
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Tarama başlatılamadı.");
      setSubmitting(false);
    }
  }

  function toggleModule(name: ModuleName) {
    setSelected((current) =>
      current.includes(name) ? current.filter((item) => item !== name) : [...current, name].sort()
    );
  }

  return (
    <form className="space-y-6" onSubmit={onSubmit}>
      <label className="block">
        <span className="mb-2 block text-sm uppercase tracking-[0.18em] text-slate-400">Target URL</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3 text-sm text-white outline-none transition focus:border-signal"
          type="url"
          required
          value={targetUrl}
          onChange={(event) => setTargetUrl(event.target.value)}
          placeholder="https://example.com"
        />
      </label>

      <div>
        <p className="mb-3 text-sm uppercase tracking-[0.18em] text-slate-400">Modules</p>
        <div className="grid gap-3 sm:grid-cols-2">
          {defaultModules.map((moduleName) => {
            const active = selected.includes(moduleName);
            return (
              <button
                key={moduleName}
                type="button"
                onClick={() => toggleModule(moduleName)}
                className={`rounded-2xl border px-4 py-3 text-left text-sm transition ${
                  active
                    ? "border-signal bg-signal/12 text-white"
                    : "border-white/10 bg-white/5 text-slate-300 hover:border-white/20"
                }`}
              >
                <span className="block font-medium">{moduleLabels[moduleName]}</span>
                <span className="mt-1 block text-xs uppercase tracking-[0.16em] text-slate-400">{moduleName}</span>
              </button>
            );
          })}
        </div>
      </div>

      {error ? <p className="rounded-2xl border border-danger/40 bg-danger/10 px-4 py-3 text-sm text-danger">{error}</p> : null}

      <button
        type="submit"
        disabled={submitting || selected.length === 0}
        className="w-full rounded-2xl bg-glow px-5 py-3 text-sm font-semibold text-slate-950 transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {submitting ? "Tarama başlatılıyor..." : "Scan başlat"}
      </button>
    </form>
  );
}
