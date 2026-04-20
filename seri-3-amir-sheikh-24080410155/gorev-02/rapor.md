# G2 - SQL Injection Lab

Bu calisma yerel ve egitsel ortamda yapilan SQL Injection laboratuvarini ozetler.

## 1. Uyguladigim adimlar

1. Juice Shop lokal ortamda calistirildi.
2. Login formunda farkli SQLi payloadlari denendi.
3. Davranis farklari not edildi.
4. Sonra parameterized query kullanan savunma kodu yazildi.
5. Before/after mantigi raporlandi.

## 2. Denenen payloadlar

1. `' OR 1=1--`
2. `admin'--`
3. `' OR '1'='1`

## 3. Karsilastigim hatalar / debug

1. Her payload her uygulamada calismadi.
2. Modern framework kullanan yerlerde dogrudan SQLi tetiklemek zorlasti.
3. Basarili denemenin neden basarili oldugunu anlamak icin cevap boyutu ve hata mesaji izlendi.

## 4. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-login-form.png`
- `ekran-goruntuleri/02-payload-test.png`
- `ekran-goruntuleri/03-savunma-kodu.png`

## 5. Ogrendigim 3 sey

- String birlestirme ile sorgu yazmak ciddi risk olusturur.
- Prepared statement SQLi riskini ciddi oranda azaltir.
- "Calisti" demek yetmez; neden calistigini yorumlamak gerekir.
