# G6 - OAuth 2.0 + PKCE Demo

## 1. Uyguladigim adimlar

1. Google OAuth istemcisi icin lokal redirect URI tanimlandi.
2. `code_verifier` ve `code_challenge` mantigi hazirlandi.
3. Login endpoint kullaniciyi yetkilendirme ekranina yonlendirdi.
4. Callback tarafinda code exchange akisi not edildi.
5. Profil verisinin ekrana basildigi akis dokumante edildi.

## 2. Karsilastigim hatalar / debug

1. Redirect URI bire bir eslesmezse hata alindi.
2. PKCE challenge formatinda base64url kullanmak onemliydi.
3. Secret dosyasinin `.gitignore` icinde oldugundan emin olmak gerekti.

## 3. Sonuc / ekran goruntusu

- `ekran-goruntuleri/01-google-console.png`
- `ekran-goruntuleri/02-login-redirect.png`
- `ekran-goruntuleri/03-callback-profile.png`

## 4. Ogrendigim 3 sey

- PKCE public client senaryolarinda ciddi guvenlik katkisi saglar.
- `code_verifier` asla istemci tarafinda loglanmamali.
- Redirect URI hassas bir kontrol noktasi oldugu icin bire bir dogru ayarlanmalidir.
