# G10 - Security Headers A+

## 1. Uyguladigim adimlar

1. Basit bir Express uygulamasi deploy veya lokal test icin hazirlandi.
2. Ilk durumda response header'lari kontrol edildi.
3. `helmet` kullanilarak guvenlik basliklari eklendi.
4. Son durumda skorun nasil iyilestigi not edildi.

## 2. Karsilastigim hatalar / debug

1. Bazi basliklar reverse proxy uzerinden ezilebilir.
2. HSTS yalnizca HTTPS senaryosunda anlamlidir.
3. Securityheaders.com sonucunu yorumlarken deployment ortami da onemlidir.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-before.png`
- `ekran-goruntuleri/02-after.png`
- `ekran-goruntuleri/03-report-link.png`

## 4. Ogrendigim 3 sey

- Header sertlestirmesi hizli ama etkili bir savunma katmanidir.
- CSP, HSTS ve nosniff gibi basliklar birlikte anlam kazanir.
- Public test sonucu ile lokal header kontrolu birbirini tamamlar.
