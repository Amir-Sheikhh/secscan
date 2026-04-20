import { ModuleResult, moduleLabels } from "@/lib/types";

const toneMap: Record<string, string> = {
  critical: "border-danger/60 bg-danger/10 text-danger",
  high: "border-danger/40 bg-danger/10 text-danger",
  medium: "border-glow/40 bg-glow/10 text-glow",
  low: "border-signal/30 bg-signal/10 text-signal",
  info: "border-white/10 bg-white/5 text-slate-300"
};

const statusMap: Record<string, string> = {
  completed: "Tamamlandı",
  running: "Çalışıyor",
  failed: "Hata",
  pending: "Bekliyor",
  skipped: "Atlandı"
};

export function ModuleCard({ moduleName, result }: { moduleName: string; result?: ModuleResult }) {
  const tone = toneMap[result?.severity ?? "info"] ?? toneMap.info;
  const findings = result?.findings ?? [];

  return (
    <article className="glass-panel rounded-[26px] p-5">
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-slate-400">{moduleName}</p>
          <h3 className="mt-1 text-xl font-semibold text-white">
            {moduleLabels[moduleName as keyof typeof moduleLabels] ?? moduleName}
          </h3>
        </div>
        <span className={`rounded-full border px-3 py-1 text-xs uppercase tracking-[0.16em] ${tone}`}>
          {statusMap[result?.status ?? "pending"]}
        </span>
      </div>

      <div className="mt-5 flex items-end justify-between">
        <div>
          <p className="text-sm text-slate-400">Skor</p>
          <p className="text-4xl font-semibold text-white">{result?.score ?? 0}</p>
        </div>
        <p className="text-sm text-slate-400">{result?.durationMs ?? 0} ms</p>
      </div>

      <p className="mt-4 text-sm leading-6 text-slate-300">{result?.summary ?? "Modül beklemede."}</p>

      <div className="mt-5 space-y-3">
        {findings.length === 0 ? (
          <p className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-300">
            Öne çıkan bulgu yok.
          </p>
        ) : (
          findings.slice(0, 3).map((finding) => (
            <div key={`${finding.title}-${finding.evidence}`} className="rounded-2xl border border-white/10 bg-black/10 px-4 py-3">
              <div className="flex items-center justify-between gap-3">
                <p className="text-sm font-medium text-white">{finding.title}</p>
                <span className="text-xs uppercase tracking-[0.18em] text-slate-400">{finding.severity}</span>
              </div>
              <p className="mt-2 text-sm leading-6 text-slate-300">{finding.recommendation}</p>
            </div>
          ))
        )}
      </div>
    </article>
  );
}
