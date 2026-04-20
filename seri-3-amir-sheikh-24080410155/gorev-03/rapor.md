# G3 - XSS + CSP Korumasi

## 1. Uyguladigim adimlar

1. Reflected XSS mantigini gormek icin basit bir Express ornegi yazildi.
2. `req.query.name` degeri HTML icine yazdirildi.
3. Payload ile yansima test edildi.
4. Sonrasinda CSP basligi eklendi.
5. Tarayici konsolunda CSP ihlali kontrol edildi.

## 2. Test mantigi

- Unsafe endpoint: ciktiyi dogrudan HTML'e basar.
- Safe endpoint: CSP header ekler.
- Gerekirse ek olarak HTML escaping uygulanir.

## 3. Karsilastigim hatalar / debug

1. Bazi payloadlar tarayicida filtreye takildi.
2. CSP dogru yazilmazsa beklenen bloklama gorulmedi.
3. Sadece CSP yetmez; output encoding de gerekir.

## 4. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-unsafe.png`
- `ekran-goruntuleri/02-devtools-csp.png`
- `ekran-goruntuleri/03-safe.png`

## 5. Ogrendigim 3 sey

- XSS yalnizca script etiketi ile sinirli degildir.
- CSP ikinci savunma hattidir.
- Kullanici girdisini HTML'e basmadan once encode etmek gerekir.
