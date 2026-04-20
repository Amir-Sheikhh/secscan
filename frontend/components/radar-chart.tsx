import { moduleLabels, type Scan } from "@/lib/types";

function pointFor(angle: number, radius: number) {
  const radians = (angle - 90) * (Math.PI / 180);
  const x = 150 + Math.cos(radians) * radius;
  const y = 150 + Math.sin(radians) * radius;
  return `${x},${y}`;
}

export function RadarChart({ scan }: { scan: Scan }) {
  const modules = Object.entries(scan.modules);
  const axisCount = Math.max(modules.length, 1);
  const points = modules.map(([name, result], index) =>
    pointFor((360 / axisCount) * index, 100 * ((result.score ?? 0) / 100))
  );

  return (
    <div className="glass-panel rounded-[28px] p-6">
      <p className="text-xs uppercase tracking-[0.22em] text-signal">Module Profile</p>
      <h3 className="mt-2 text-2xl font-semibold">Radar görünümü</h3>
      <div className="mt-6 flex flex-col items-center gap-6 lg:flex-row">
        <svg viewBox="0 0 300 300" className="h-72 w-72 shrink-0">
          {[40, 70, 100].map((radius) => (
            <polygon
              key={radius}
              points={modules.map((_, index) => pointFor((360 / axisCount) * index, radius)).join(" ")}
              fill="none"
              stroke="rgba(255,255,255,0.12)"
            />
          ))}

          {modules.map(([name], index) => (
            <line
              key={name}
              x1="150"
              y1="150"
              x2={pointFor((360 / axisCount) * index, 110).split(",")[0]}
              y2={pointFor((360 / axisCount) * index, 110).split(",")[1]}
              stroke="rgba(255,255,255,0.18)"
            />
          ))}

          <polygon points={points.join(" ")} fill="rgba(106,242,212,0.24)" stroke="#6af2d4" strokeWidth="2" />

          {modules.map(([name, result], index) => {
            const [x, y] = pointFor((360 / axisCount) * index, 124).split(",");
            return (
              <text key={`${name}-label`} x={x} y={y} fill="#d9e8e6" fontSize="11" textAnchor="middle">
                {moduleLabels[name as keyof typeof moduleLabels] ?? name} {result.score}
              </text>
            );
          })}
        </svg>

        <div className="grid w-full gap-3">
          {modules.map(([name, result]) => (
            <div key={name} className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3">
              <div className="flex items-center justify-between gap-3">
                <span className="text-sm font-medium text-white">{moduleLabels[name as keyof typeof moduleLabels] ?? name}</span>
                <span className="text-sm text-signal">{result.score}</span>
              </div>
              <div className="mt-3 h-2 rounded-full bg-white/10">
                <div className="h-2 rounded-full bg-signal" style={{ width: `${result.score}%` }} />
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
