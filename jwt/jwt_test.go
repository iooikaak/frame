package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestGetJwtToken(t *testing.T) {
	sign := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----
`)
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(sign)
	t.Log(err)
	claims := &CustomClaims{
		UserInfo: &UserInfo{
			UserId:   111,
			RoleId:   1,
			UserName: "czl",
		},
		RegisteredClaims: RegisteredClaims{
			Issuer:    "test_issuer",
			Subject:   "test_subject",
			Audience:  []string{"test-look"},
			ExpiresAt: NewNumericDate(time.Now().Add(300 * time.Second)),
		},
	}
	a, b := GetJwtToken(SigningMethodRS256, parsedKey, claims)
	t.Log(a, b)
}

func TestParseJwtToken(t *testing.T) {
	sign := []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`)
	key, err := ParseRSAPublicKeyFromPEM(sign)
	t.Log(err)
	a, b := ParseJwtToken(key, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2luZm8iOnsidXNlcl9pZCI6MTExLCJyb2xlX2lkIjoxLCJ1c2VyX25hbWUiOiJjemwifSwiaXNzIjoidGVzdF9pc3N1ZXIiLCJzdWIiOiJ0ZXN0X3N1YmplY3QiLCJhdWQiOlsidGVzdC1sb29rIl0sImV4cCI6MTYzNzE0MDk4MCwibmJmIjoxNjM3MTQwNjgwLCJpYXQiOjE2MzcxNDA2ODAsImp0aSI6IjIwN0VGNThDOEZBODRCNTU5NzJGRThDNkZCNzc4RjNGIn0.0mtevFoqGM8dA5Arp8LAZMt6TyL6mSRqwJQH8p-WpLZM9NJ_u1xhdY6C84AyTEGrCww_0txhqrfbuRYshLltZ_6t3TmyKvnIXq3iITEA7nzk1LpQTUTmd9kjvvTjlAEmuluEzYis5KN85QmmEFWD4BkrmslfGpHbYt0ct5dIF2a2PlehZ2G7ET94fvd3v1tBDfbbTGDl-_u6zPMlue4mqztjBpkb2_q_C22V1V-7R1BmNiqpIwHPH2FTxrKmYmToWnbJiOSXaRn1eSoE58gzQJS3xzf9VVRk5MgQqLXShgcDJD16hczr_RhITR3b5U62oVFPU9dcaNZ6q6kO3LuEYw")
	t.Log(a, b)
}

func TestGetJwtTokenByHS256(t *testing.T) {
	parsedKey := []byte(`test`)
	claims := &CustomClaims{
		UserInfo: &UserInfo{
			UserId:   111,
			RoleId:   1,
			UserName: "czl",
		},
		RegisteredClaims: RegisteredClaims{
			Issuer:    "test_issuer",
			Subject:   "test_subject",
			Audience:  []string{"test-look"},
			ExpiresAt: NewNumericDate(time.Now().Add(30000 * time.Second)),
		},
	}
	a, b := GetJwtToken(SigningMethodHS256, parsedKey, claims)
	t.Log(a, b)
}

func TestParseJwtTokenByHS256(t *testing.T) {
	key := []byte(`test`)
	a, b := ParseJwtToken(key, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2luZm8iOnsidXNlcl9pZCI6MTExLCJyb2xlX2lkIjoxLCJ1c2VyX25hbWUiOiJjemwifSwiaXNzIjoidGVzdF9pc3N1ZXIiLCJzdWIiOiJ0ZXN0X3N1YmplY3QiLCJhdWQiOlsidGVzdC1sb29rIl0sImV4cCI6MTYzNzE4NDUyOSwibmJmIjoxNjM3MTU0NTI5LCJpYXQiOjE2MzcxNTQ1MjksImp0aSI6IjgyMkY4RkM1QTI0QjQ5Qjk4NzBDMjE5RjE3OThBQjI4In0.FioBXLH8WIXCEVKQd-C70ihjF_sm0DuBOoaA1QIb6fA")
	t.Log(a, b)
}
