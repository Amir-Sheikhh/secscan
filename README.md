# SecScan
Öğrenci: amir sheikh
Okul No: 24080410155
Ders: BMU1208 Web Tabanlı Programlama
Seri: Finale Doğru / Seri 3
Teslim Türü: Kısım 2 / `secscan`
Demo URL: SUBMIT_BEFORE_GITHUB
Screenshot Klasörü: `docs/screenshots/`

SecScan, kullanıcıdan hedef URL alıp 7 güvenlik modülünü paralel çalıştıran bir web güvenliği analiz platformudur. Backend katmanında Go tabanlı HTTP servis, frontend katmanında Next.js + TypeScript + TailwindCSS kullanılır. Tarama ilerlemesi SSE ile canlı izlenebilir, sonuçlar JSON ve PDF olarak alınabilir.

## Teslim Özeti

- Ayrı public GitHub repo olarak hazırlanmıştır.
- README ilk 10 satırında öğrenci bilgisi bulunur.
- Backend için çalışır API, SSRF guard ve testler vardır.
- Frontend için URL giriş ekranı, canlı durum ve sonuç dashboard'u vardır.
- CI ve security workflow dosyaları repo kökünde bulunur.

## Özellikler

- `POST /api/scan` ile asenkron tarama başlatma
- `GET /api/scan/:id` ile sonuç okuma
- `GET /api/scan/:id/stream` ile SSE akışı
- `GET /api/scan/:id/report.pdf` ile PDF raporu
- SSRF koruması: private, loopback ve link-local IP blokları reddedilir
- 7 modül: `ports`, `headers`, `tls`, `fuzz`, `xss`, `sqli`, `cve`

## Mimari

```text
secscan/
|-- README.md
|-- docker-compose.yml
|-- .github/workflows/
|-- docs/screenshots/
|-- backend/
|   |-- main.go
|   |-- internal/api
|   |-- internal/scanner
|   \-- internal/storage
\-- frontend/
    |-- app/
    |-- components/
    \-- lib/
```

## Çalıştırma

### Backend

```powershell
cd backend
go mod tidy
go test ./...
go run .
```

Varsayılan adres: `http://localhost:8080`

### Frontend

```powershell
cd frontend
npm install
npm run lint
npm run build
npm run dev
```

Varsayılan adres: `http://localhost:3000`

### Docker Compose

```powershell
docker compose up --build
```

Not: Bu ortamda `docker` kurulu olmadığı için compose komutu yerelde doğrulanamadı; repo dosyaları hazırlandı.

## Ortam Değişkenleri

### Backend

- `PORT=8080`
- `ALLOWED_ORIGINS=http://localhost:3000`
- `SECSCAN_SCAN_TIMEOUT=45s`
- `SECSCAN_ENABLE_ACTIVE_PROBES=false`

### Frontend

- `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080`

## API Uçları

- `GET /health`
- `POST /api/scan`
- `GET /api/scan/:id`
- `GET /api/scan/:id/stream`
- `GET /api/scan/:id/report.pdf`

Örnek istek:

```json
{
  "url": "https://example.com",
  "modules": ["ports", "headers", "tls", "fuzz", "xss", "sqli", "cve"]
}
```

## Modül Özeti

| Modül | Açıklama |
|---|---|
| `ports` | Sık kullanılan portları eşzamanlı TCP connect ile kontrol eder |
| `headers` | Güvenlik başlıklarını değerlendirir |
| `tls` | TLS sürümlerini ve sertifika durumunu inceler |
| `fuzz` | Yaygın path'ler üzerinde içerik keşfi yapar |
| `xss` | Reflected XSS için yansıtma heuristiği uygular |
| `sqli` | Hata ve cevap farkı temelli SQLi heuristiği uygular |
| `cve` | Teknoloji izi çıkarıp advisory eşleşmelerini arar |

## Doğrulama

Bu workspace içinde aşağıdaki kontroller başarıyla çalıştırıldı:

- `go test ./...`
- `npm run lint`
- `npm run build`
- `/health` smoke test
- `/api/scan` smoke test
- `/api/scan/:id/report.pdf` smoke test

## Demo ve Screenshot Kanıtları

Teslim öncesi aşağıdakileri doldur:

- `Demo URL: SUBMIT_BEFORE_GITHUB` satırını gerçek canlı link ile değiştir
- `docs/screenshots/` altına ana ekran, tarama ekranı ve rapor ekranı görsellerini ekle

Önerilen dosya adları:

- `docs/screenshots/01-home.png`
- `docs/screenshots/02-running-scan.png`
- `docs/screenshots/03-report-dashboard.png`

## AI Kullanımı Beyanı

Bu repo AI destekli geliştirme ile oluşturuldu. Kod ve README teslim öncesi gözden geçirildi; yine de canlı demo, ekran görüntüleri ve commit geçmişi öğrenci tarafından tamamlanmalıdır.

## Teslim Kontrol Listesi

- README ilk 10 satırda öğrenci bilgisi içeriyor
- Repo public
- Birden fazla anlamlı commit var
- Demo URL eklendi
- Screenshot dosyaları eklendi
- CI ve security workflow'ları GitHub'da çalışıyor
