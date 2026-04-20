# G9 - SBOM + Trivy Scan

## 1. Uyguladigim adimlar

1. Proje kokunde SBOM olusturmak icin `syft` kullanildi.
2. Dosya sistemi taramasi `trivy fs` ile yapildi.
3. Ciddi CVE'ler tablo halinde not edildi.
4. En az bir paket icin versiyon guncelleme plani yazildi.

## 2. Karsilastigim hatalar / debug

1. Arac kurulu degilse komutlar calismadi.
2. Bazi paket adlari lock dosyasindan okunuyor.
3. CVE listesi ciktiginda hangi paketin dogrudan veya dolayli bagimlilik oldugunu ayirmak gerekti.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-sbom.png`
- `ekran-goruntuleri/02-trivy.png`
- `ekran-goruntuleri/03-fix-plan.png`

## 4. Ogrendigim 3 sey

- SBOM sadece guvenlik degil envanter acisindan da faydalidir.
- Trivy sonucu gormek kadar onceliklendirmek de onemlidir.
- Versiyon guncellemesi bazen en hizli savunmadir.
