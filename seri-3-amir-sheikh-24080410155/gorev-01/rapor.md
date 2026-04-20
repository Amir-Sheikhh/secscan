# G1 - OWASP Top 10 Haritalama

Bu raporda OWASP Top 10 kategorilerini secilen acik kaynak proje uzerinden haritaladim.

## 1. Uyguladigim adimlar

1. Orta olcekli acik kaynak web projesi olarak Supabase secildi.
2. README ve security ile ilgili dokumanlar incelendi.
3. OWASP Top 10 listesinden 5 kategori secildi.
4. Her kategori icin projedeki koruma mekanizmasi veya eksiklik not edildi.
5. Sonuclar tabloya yazildi.
6. En dikkat cekici eksik nokta kisa yorumla aciklandi.

## 2. Inceleme tablosu

| OWASP | Gozlem | Durum |
|---|---|---|
| A01 Broken Access Control | Policy ve rol bazli erisim yaklasimi mevcut | Kismen iyi |
| A02 Cryptographic Failures | Secret yonetimi ve TLS kullanimina vurgu var | Iyi |
| A03 Injection | ORM ve parametreli sorgu yaklasimi tercih ediliyor | Iyi |
| A05 Security Misconfiguration | Varsayilan ayarlar dikkatli incelenmeli | Risk var |
| A06 Vulnerable Components | Dependabot ve patch sureci onemli | Surekli takip gerekli |

## 3. Karsilastigim hatalar / debug

1. Security dokumanlari farkli klasorlerde daginikti.
2. Bazi korumalar koddan degil README ve config dosyalarindan anlasildi.
3. Guvenlik tabindaki bilgi ile kod tabani her zaman bire bir eslesmiyor.

## 4. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-repo-overview.png`
- `ekran-goruntuleri/02-security-notes.png`
- `ekran-goruntuleri/03-owasp-table.png`

## 5. Ogrendigim 3 sey

- OWASP kategorilerini gercek bir projede haritalamak teoriden daha ogretici.
- Bazi korumalar dogrudan kodda degil altyapi ve pipeline katmaninda yer aliyor.
- "Injection yok" demek yerine hangi mekanizma ile azaltildigini gostermek gerekiyor.
