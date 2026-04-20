# G5 - JWT Guvenligi Audit

## 1. Uyguladigim adimlar

1. Kendi JWT kullanan uygulamam secildi.
2. Secret, exp, storage ve logout mantigi kontrol edildi.
3. Checklist ile guvenlik acisindan eksikler not edildi.
4. Riskli alanlar icin duzeltme onerisi yazildi.

## 2. Audit checklist sonucu

| Kontrol | Sonuc | Not |
|---|---|---|
| Secret 256-bit random mi | Kontrol edildi | `.env` icinde tutulmali |
| Hard-coded secret var mi | Hayir | Kod icinde olmamali |
| `exp` claim var mi | Evet | Kisa omurlu olmali |
| `alg: none` reddi | Kontrol edildi | Zorunlu |
| Storage guvenli mi | Kismen | HttpOnly cookie tercih edilmeli |
| Logout / revoke mantigi | Eksik olabilir | Gerekiyorsa denylist eklenmeli |

## 3. Karsilastigim hatalar / debug

1. Bazi riskler kodda degil deployment ayarinda cikti.
2. LocalStorage kullanimini fark etmek kolay oldu ama cikarmak daha zahmetli olabilir.
3. Logout mantigi cogu basit projede yetersiz kaliyor.

## 4. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-jwt-config.png`
- `ekran-goruntuleri/02-env-secret.png`
- `ekran-goruntuleri/03-cookie-or-storage.png`

## 5. Ogrendigim 3 sey

- JWT dogru imzalanmis olsa bile kotu saklanirsa risk devam eder.
- Kisa `exp` suresi hasar alanini azaltir.
- Logout problemi stateless token kullaniminda ayrica dusunulmelidir.
