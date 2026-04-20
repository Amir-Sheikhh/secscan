# G7 - Nmap + ZAP Taramasi

Yerel uygulama veya yazili izinli hedef disinda kullanma.

## 1. Uyguladigim adimlar

1. Lokal servis ayaga kaldirildi.
2. `nmap -sV -p 3000 localhost` ile versiyon tespiti yapildi.
3. ZAP baseline scan ile pasif tarama raporu alindi.
4. Bulgular OWASP kategorilerine eslendi.
5. En az bir bulgu icin duzeltme notu yazildi.

## 2. Karsilastigim hatalar / debug

1. Docker yoksa ZAP komutu calismayabilir.
2. Host adresi Docker icinde farkli tanimlanabilir.
3. Baseline scan aktif exploitation yapmadigi icin daha sinirli bulgu verir.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-nmap.png`
- `ekran-goruntuleri/02-zap-summary.png`
- `ekran-goruntuleri/03-fix-note.png`

## 4. Ogrendigim 3 sey

- Nmap servis gorunurlugunu hizli sekilde ozetler.
- ZAP pasif taramada bile faydali konfig bulgulari verebilir.
- Tarama sonucu kadar sonucu yorumlamak da onemlidir.
