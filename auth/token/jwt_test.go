package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAg84U9d4QdLHiXadtqiz6XL408jEgnsbuG+bRrYPCEU12woKn
/S9rn0U/4TSKp8KRdOhHVa0ZKKwFgHnYMlAuXGvpjxZA6VzNr8+r7lvXmySAl07d
xuUKqfAD09IzerNXYr2xNuET40uh55GZwjT+lGjlizf2JmutE4FCGh74E/ogmhJV
D78z8z0xcPgJLYzVCvxdQDj4yGfbYsQ+1Ni+KThhGZT/EXYdx5Z8iewCQV1ULhrK
+Rq0WjemrQ/I4prwRLNsmKvxLcI0iGCVpyuRWqkfuwV0YfyImkRCvZIHSKe4dO8Y
Y7mEJVliOX2WL0RxLLOkrV5cvOZhu/9bduKe1wIDAQABAoIBAEtR2gKCqu60tIIG
apHD8DJNc54vWs/BKFKDfbDlSWJv3PzcgzkY/yxd+1orl0y00EB5eEJKj8UBQIeh
mV1vGn8wH6Dn+6IfqV7dkwe8LiJ3IpDUvcGqI4TnJpjGVyq4D6jac8nDp4TNNLFL
MA2eERkoNHxxN1XPVHF//RFGP0fSZr57Sel7NuvJHcGxk5sz82oTGJfFpVpW4dZX
RvS6xjx3hlYwuXkEfuN2dgBGqL3nImcveWFRViAmPS3TW2VzLfA2CnCa7kxnG+lI
v3XvGRDpl5OYret59eJOPH6xG158yoGuCdqNPtwpzAqo2oAZ4mcoMjp3MSQOM74Z
y1uUTRECgYEAwteIGZZnIAOE64OZjPQpC507dPoRhX2B6/KyItma6Rh8Os8g5vMZ
y6dUpei4zxB8DIGcAuQL4sQ/2daE9LXq4ej+y32enQb3SvNsJL5WQPe4jGMTDWsc
U7GhCV4r4WQ0CRejYpYD5KOra93xsL6i3TkRYahhwvjsIZDbkHKuxz8CgYEArS08
kgVl1A7H0EfyAJOYaq4pckdiBFSVJqd6tKJKrvZL509AbR07o/SKiCUoHozolAu3
SaAHDOidXOlZpukZqBaSCyfTRI6oIXlOo3X/i1msSprKWxHBEyxZp2fgvLq/LyqZ
ZmIvUH23nOBiJLCoV3vUO/XrDm1QQC4cEM/pmmkCgYEAnqYGwObxc2TKL0aJmfcZ
EMbnKdmQuMQ4LNoB6FSNSW1RgkUzgjnB8EyApVL4YEoI59oFIWl0sCGh6As/WU5j
Qa2JAkJ4C14nr9TDYqvE6cOLdmwZkFx9xTwmZs1SJ4WCxUCFHfoOk3YdV4hxiru/
OyiDmaQUbkBnbPFZhqWK4NsCgYA9ybEdzG07jxZ92t2elQrBrWg+TPfM4bzhsMnY
HzuUV25XlnA/PjnkUsEGuHMrC02EXPXFgCJj2a8j0mJZajvsPDlZX5lCkb+tSdHk
Aprtxk3xxG7EtX308FMAptCJpfvGwWVAIXIOPvy/LVP3EUzPAfCEgEagvCHw9EKx
QO1xWQKBgQCvFHbTt4iy3UjBIfd/19+CPH0NhhA4zAsXVx26hdGvQ4ORamUt4n2v
heC7b07WqwwLfsMzbFSGuroliGTqGj70G9cNPFiqT333R+GPXzB9xpV6TtwaWgU7
fGGH9bylPyXW2WXrt4ThiUTpzLt6z8Fd1tmqxVhNv/MviY4ilQ87Bw==
-----END RSA PRIVATE KEY-----`

const PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAg84U9d4QdLHiXadtqiz6
XL408jEgnsbuG+bRrYPCEU12woKn/S9rn0U/4TSKp8KRdOhHVa0ZKKwFgHnYMlAu
XGvpjxZA6VzNr8+r7lvXmySAl07dxuUKqfAD09IzerNXYr2xNuET40uh55GZwjT+
lGjlizf2JmutE4FCGh74E/ogmhJVD78z8z0xcPgJLYzVCvxdQDj4yGfbYsQ+1Ni+
KThhGZT/EXYdx5Z8iewCQV1ULhrK+Rq0WjemrQ/I4prwRLNsmKvxLcI0iGCVpyuR
WqkfuwV0YfyImkRCvZIHSKe4dO8YY7mEJVliOX2WL0RxLLOkrV5cvOZhu/9bduKe
1wIDAQAB
-----END PUBLIC KEY-----`

func TestGenerateToken(t *testing.T) {
	// 解析 private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PRIVATE_KEY))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", err)
	}

	// 工厂模式获取 jwt generator
	g := NewJWTGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	// 生成 token
	token, err := g.GenerateToken("61b1e4caf6d536ccefdae779", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate jwt: %v", err)
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFiMWU0Y2FmNmQ1MzZjY2VmZGFlNzc5In0.BB7QiMOvoeJZTBLHYuEYaay2pHI1Yp4UKiGerdXom7Fs4YY00dta-JnJTirBU94sMITeUItaQXqGiVjhmYuQIh-34t7Itb7BFiIrDKCQuTGftiYm6KuVPmX_8JrL1qVK9ni8nTTZC6m2zgYUR5-mPh4LKNyp9XZEPfCn91TC8iXlY2MDkpoYO-RVAmGOHlVTXhmzWdStLyiQdRNEynDNQkfnf1eMqKe4EHk6EfwHC0XvJRzAKVeGKYJb5OHJD-XbNrHy_q-WncmXr6i-cUFqruB1kYoFn8tckQhQtDb-46cjNn894TI2q1PCHHqoYHr8cDwUAPcjX6990wYzbkWoYw"
	if token != want {
		t.Errorf("wrong token generate, \nwant: %q, \ngot: %q\n", want, token)
	}
}
