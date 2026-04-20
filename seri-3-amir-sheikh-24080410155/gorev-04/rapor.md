# G4 - CSRF Token Sistemi

## 1. Uyguladigim adimlar

1. Basit transfer formu olan bir Express uygulamasi kuruldu.
2. Once CSRF korumasi olmadan form gonderimi test edildi.
3. Baska bir sayfadan otomatik POST atan `evil-form.html` ile deneme yapildi.
4. Sonra `csurf` middleware eklendi.
5. Form icine hidden CSRF token basildi.
6. Ayni saldiri tekrar denendi ve bloklandigi gozlemlendi.

## 2. Karsilastigim hatalar / debug

1. Cookie parser olmadan `csurf` beklenen sekilde calismadi.
2. Token form icinde dogru isimle gonderilmezse 403 alindi.
3. Test ederken ayni tarayici oturumu kullanmak gerekti.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-transfer-form.png`
- `ekran-goruntuleri/02-evil-form.png`
- `ekran-goruntuleri/03-403-forbidden.png`

## 4. Ogrendigim 3 sey

- CSRF token session veya cookie baglami ile calisir.
- Sadece POST olmak koruma saglamaz.
- SameSite cookie ve CSRF token birlikte daha gucludur.
