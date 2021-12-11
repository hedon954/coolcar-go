package shared_token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAg84U9d4QdLHiXadtqiz6
XL408jEgnsbuG+bRrYPCEU12woKn/S9rn0U/4TSKp8KRdOhHVa0ZKKwFgHnYMlAu
XGvpjxZA6VzNr8+r7lvXmySAl07dxuUKqfAD09IzerNXYr2xNuET40uh55GZwjT+
lGjlizf2JmutE4FCGh74E/ogmhJVD78z8z0xcPgJLYzVCvxdQDj4yGfbYsQ+1Ni+
KThhGZT/EXYdx5Z8iewCQV1ULhrK+Rq0WjemrQ/I4prwRLNsmKvxLcI0iGCVpyuR
WqkfuwV0YfyImkRCvZIHSKe4dO8YY7mEJVliOX2WL0RxLLOkrV5cvOZhu/9bduKe
1wIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {

	// parse public key
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PUBLIC_KEY))
	if err != nil {
		t.Errorf("cannot parse public key: %v", err)
	}

	// get JWT verifer
	v := &JWTVerifier{
		PublicKey: pk,
	}

	// set standard time, for test stability
	jwt.TimeFunc = func() time.Time {
		return time.Unix(1639231380, 0)
	}

	// test tables
	cases := []struct {
		name    string
		token   string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid_token",
			token:   "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMzg1ODAsImlhdCI6MTYzOTIzMTM4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFiMWU0Y2FmNmQ1MzZjY2VmZGFlNzc5In0.gEri-huXgqdr5Vb-69deUjOKAdNCI5NAjoyfluiT3Ai4MgtmNg8K03uawkfG7Z66Duxh2nz0TBzPKRp61aKUajQITXB3tCGbvKENmVEx3C__rTDSUTkU-EPJGIn_bb1lQhwqcKBOQIcxlAsfKnZIHPpUCchgaAKr79Bpq0KZtx_Yt27mDEFRGHsE4zZaLYOAAGewqWa8oJgdlTzhGeHXJtqX5fvdL3t_kp0fEFML1eDkL1fqYB5XGpSfg6CtLOzJWzbshYWJMcHBmsTyq9GpzuJEr9Fvu6PkMEHtT3s4LICKkLWKFIllXxHSN6G8i-TTSDfiv-rQCEZBefoG0hET3w",
			now:     time.Unix(1639231480, 0),
			want:    "61b1e4caf6d536ccefdae779",
			wantErr: false,
		},
		{
			name:    "valid_expired",
			token:   "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMzg1ODAsImlhdCI6MTYzOTIzMTM4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFiMWU0Y2FmNmQ1MzZjY2VmZGFlNzc5In0.gEri-huXgqdr5Vb-69deUjOKAdNCI5NAjoyfluiT3Ai4MgtmNg8K03uawkfG7Z66Duxh2nz0TBzPKRp61aKUajQITXB3tCGbvKENmVEx3C__rTDSUTkU-EPJGIn_bb1lQhwqcKBOQIcxlAsfKnZIHPpUCchgaAKr79Bpq0KZtx_Yt27mDEFRGHsE4zZaLYOAAGewqWa8oJgdlTzhGeHXJtqX5fvdL3t_kp0fEFML1eDkL1fqYB5XGpSfg6CtLOzJWzbshYWJMcHBmsTyq9GpzuJEr9Fvu6PkMEHtT3s4LICKkLWKFIllXxHSN6G8i-TTSDfiv-rQCEZBefoG0hET3w",
			now:     time.Unix(1539231480, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			token:   "bad_token",
			now:     time.Unix(1639231480, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			token:   "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMzg1ODAsImlhdCI6MTYzOTIzMTM4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFiMWU0Y2FmNmQ1MzZjY2VmZGFlNzc0In0.Dt4KvsA_KfyyGPoQhBmZ4ckBi8tUeaOE54P00XTaFN8uk_pvxtsJq-Q9GBpyEmY33HFqtGLinpxZYFz-Zd8UXGIGavuYGQWZYKvk2qk8koLv8UBGMkXbPTmXha_CYC_YYDRcyvt0TmluQCj3zHx3_GpFv2fwHnZ1PsuC1vb9S4mfCfXY9zXvvV1_oeu9i5HXwEjeuEmDwBx9vpi8FMrpFjHPOy1U6YcL8GuYYfTEm3yjuUbL_JiDfmgw2tOjpbgkudJlo7taRRpi7Wwanbnd6Neas5Uf_u7XUgZKD4du2GDi-TO5SKp_RCT3EqD4XFhgjx59zXYqpKj12xHXw21Ygw",
			now:     time.Unix(1639231480, 0),
			wantErr: true,
		},
	}

	// test
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return testCase.now
			}
			accountID, err := v.Verify(testCase.token)
			if !testCase.wantErr && err != nil {
				t.Errorf("verify token failed, error: %v", err)
			}
			if testCase.wantErr && err == nil {
				t.Errorf("verification error, should have error, but got no error")
			}
			if accountID != testCase.want {
				t.Errorf("verify token failed, want: %q, got %q\n", testCase.want, accountID)
			}
		})
	}

	// verify token
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMzg1ODAsImlhdCI6MTYzOTIzMTM4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFiMWU0Y2FmNmQ1MzZjY2VmZGFlNzc5In0.gEri-huXgqdr5Vb-69deUjOKAdNCI5NAjoyfluiT3Ai4MgtmNg8K03uawkfG7Z66Duxh2nz0TBzPKRp61aKUajQITXB3tCGbvKENmVEx3C__rTDSUTkU-EPJGIn_bb1lQhwqcKBOQIcxlAsfKnZIHPpUCchgaAKr79Bpq0KZtx_Yt27mDEFRGHsE4zZaLYOAAGewqWa8oJgdlTzhGeHXJtqX5fvdL3t_kp0fEFML1eDkL1fqYB5XGpSfg6CtLOzJWzbshYWJMcHBmsTyq9GpzuJEr9Fvu6PkMEHtT3s4LICKkLWKFIllXxHSN6G8i-TTSDfiv-rQCEZBefoG0hET3w"
	accountID, err := v.Verify(token)
	if err != nil {
		t.Errorf("verify token error: %v", err)
	}

	// compare
	want := "61b1e4caf6d536ccefdae779"
	if want != accountID {
		t.Errorf("verify token failed, want: %q, got: %q\n", want, accountID)
	}

}
