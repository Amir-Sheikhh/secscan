# G8 - SAST Pipeline (Semgrep)

## 1. Uyguladigim adimlar

1. GitHub Actions altina Semgrep workflow eklendi.
2. Ornek zafiyetli kod dosyasi olusturuldu.
3. CI calistiginda bulgu uretmesi beklendi.
4. Sonra kod duzeltilip temiz sonuc hedeflendi.

## 2. Karsilastigim hatalar / debug

1. Workflow dosyasi yanlis path'te ise Actions tetiklenmedi.
2. Kasti bug yeterince belirgin degilse Semgrep yakalamayabilir.
3. Sadece CI eklemek yetmez; raporu okumak gerekir.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-actions-run.png`
- `ekran-goruntuleri/02-semgrep-alert.png`
- `ekran-goruntuleri/03-fixed-run.png`

## 4. Ogrendigim 3 sey

- SAST en verimli halini pipeline icinde verir.
- Kodu PR seviyesinde taramak geri bildirim suresini kisaltir.
- False positive ile gercek bug'u ayirt etmeyi ogrenmek gerekir.
