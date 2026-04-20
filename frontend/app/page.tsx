import { HomeForm } from "@/components/home-form";

const readinessPoints = [
  "README ilk 10 satırında ad soyad ve okul no bulunmalı.",
  "Repo public olmalı ve tek commit yerine süreç görünmeli.",
  "Backend `/health` ve `/api/scan` uçları çalışır durumda olmalı.",
  "Frontend canlı ilerleme, sonuç özeti ve PDF bağlantısı göstermeli."
];

export default function HomePage() {
  return (
    <main className="page-shell">
      <section className="mx-auto flex min-h-screen w-full max-w-7xl flex-col gap-10 px-5 py-8 sm:px-8 lg:px-10">
        <header className="grid gap-8 lg:grid-cols-[1.2fr_0.8fr]">
          <div className="glass-panel rounded-[32px] p-8 sm:p-10">
            <p className="mb-4 inline-flex rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs uppercase tracking-[0.3em] text-signal">
              Finale Doğru • SecScan
            </p>
            <h1 className="max-w-3xl text-4xl font-semibold leading-tight sm:text-6xl">
              Hedef URL gir, 7 modülü paralel çalıştır, raporu canlı izle.
            </h1>
            <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300 sm:text-lg">
              Bu arayüz Seri 3 final projesi için hazırlandı. Backend tarafında SSRF guard, async tarama ve SSE
              bulunur; burada ise başlatma, ilerleme, özet skor ve rapor akışı tek ekranda toplanır.
            </p>
            <div className="mt-8 grid gap-3 sm:grid-cols-2">
              {readinessPoints.map((point) => (
                <div
                  key={point}
                  className="rounded-2xl border border-white/10 bg-white/5 px-4 py-4 text-sm leading-6 text-slate-200"
                >
                  {point}
                </div>
              ))}
            </div>
          </div>

          <aside className="glass-panel rounded-[32px] border border-glow/20 p-8 sm:p-10">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <p className="text-sm uppercase tracking-[0.22em] text-glow">Scan Control</p>
                <h2 className="mt-2 text-2xl font-semibold">Yeni Tarama</h2>
              </div>
              <div className="metric-ring flex h-16 w-16 items-center justify-center rounded-full p-[6px]" style={{ ["--score" as string]: 92 }}>
                <div className="flex h-full w-full items-center justify-center rounded-full bg-base text-lg font-semibold">A</div>
              </div>
            </div>
            <HomeForm />
          </aside>
        </header>

        <section className="grid gap-5 md:grid-cols-3">
          <div className="glass-panel rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.22em] text-signal">Rule Fit</p>
            <h3 className="mt-3 text-xl font-semibold">AI değerlendirmeye uygun yapı</h3>
            <p className="mt-3 text-sm leading-6 text-slate-300">
              Çalışır klasör yapısı, CI dosyaları, testler ve README teslim kurallarıyla birlikte gelir.
            </p>
          </div>
          <div className="glass-panel rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.22em] text-signal">Safety</p>
            <h3 className="mt-3 text-xl font-semibold">SSRF koruması varsayılan açık</h3>
            <p className="mt-3 text-sm leading-6 text-slate-300">
              Private, loopback ve link-local ağlara istek atılmaz. Aktif SQLi testleri varsayılan olarak kapalıdır.
            </p>
          </div>
          <div className="glass-panel rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.22em] text-signal">Deliverable</p>
            <h3 className="mt-3 text-xl font-semibold">Dashboard + PDF + canlı durum</h3>
            <p className="mt-3 text-sm leading-6 text-slate-300">
              Tarama bittiğinde detay kartları, olay akışı ve backend PDF çıktısı aynı akış içinde görünür.
            </p>
          </div>
        </section>
      </section>
    </main>
  );
}
