# JWT Audit Checklist

- Secret `.env` icinde mi
- En az 32 byte random deger mi
- Access token omru 30 dakikanin altinda mi
- Refresh token stratejisi var mi
- `alg: none` kabul edilmiyor mu
- Token localStorage yerine HttpOnly cookie icinde mi
- Logout sonrasi iptal mekanizmasi var mi
