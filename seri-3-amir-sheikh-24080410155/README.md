# Seri 3 Lab
Öğrenci: amir sheikh
Okul No: 24080410155
Ders: BMU1208 Web Tabanlı Programlama
Seri: Finale Doğru / Seri 3
Teslim Türü: Kısım 1 / Rapor Çıktıları
Repo Adı Önerisi: `seri-3-lab`
Branch Alternatifi: `seri-3-lab`

Bu repo Seri 3 Kısım 1 için 10 görevin raporlarını ve görev bazlı yardımcı dosyaları içerir. Her görev için `rapor.md` ve `ekran-goruntuleri/` klasörü bulunur. Bazı görevlerde kurala uygun olarak ek dosyalar da yer alır.

## Görev Özeti

| Görev | Başlık | İçerik |
|---|---|---|
| G1 | OWASP Top 10 Haritalama | `rapor.md` + ekran görüntüleri |
| G2 | SQL Injection Lab | `rapor.md` + `saldiri.txt` + `savunma.js` |
| G3 | XSS + CSP Koruması | `rapor.md` + `xss-demo.js` |
| G4 | CSRF Token Sistemi | `rapor.md` + `evil-form.html` + `transfer-app.js` |
| G5 | JWT Güvenliği Audit | `rapor.md` + checklist |
| G6 | OAuth 2.0 + PKCE Demo | `rapor.md` + `pkce-example.js` |
| G7 | Nmap + ZAP Taraması | `rapor.md` + komut dosyası |
| G8 | SAST Pipeline (Semgrep) | `rapor.md` + workflow + örnek zafiyetli dosya |
| G9 | SBOM + Trivy Scan | `rapor.md` + komut listesi + CVE tablosu |
| G10 | Security Headers A+ | `rapor.md` + `helmet-example.js` |

## Repo Yapısı

```text
seri-3-amir-sheikh-24080410155/
|-- README.md
|-- .gitignore
|-- .github/workflows/semgrep.yml
|-- gorev-01/
|-- gorev-02/
|-- ...
\-- gorev-10/
```

## Rapor Kontrol Listesi

Her `rapor.md` şu 5 bölümü içerir:

1. Görev numarası ve başlığı
2. Uyguladığım adımlar
3. Karşılaştığım hatalar
4. Sonuç / ekran görüntüsü
5. Öğrendiğim 3 şey

## Kendi Tahminim

- Yapısal teslim gereksinimleri: hazır
- Görev klasörleri ve rapor başlıkları: hazır
- GitHub Actions için Semgrep workflow yolu: hazır
- Tam puan için kalan manuel işler: gerçek screenshot eklemek, raporları son kişisel anlatımla netleştirmek, GitHub public repo/branch ve düzenli commit geçmişi oluşturmak

## Teslim Kontrol Listesi

- README ilk 10 satırda öğrenci bilgisi içeriyor
- Her görevde `rapor.md` var
- Her görevde `ekran-goruntuleri/` klasörü var
- G2 için `saldiri.txt` ve `savunma.js` var
- G8 için repo kökünde `.github/workflows/semgrep.yml` var
- Gerçek ekran görüntüleri eklendi
- Repo public veya öğretim üyesine erişilebilir durumda
- Düzenli commit geçmişi oluştu
