package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQBwEAjRw383hs3ioocpqVYWESnaOVAneNJkIB9RCal+Ax0s0JDa
Pb4eR2lT3fsuRXA2LOKyXYVr2F5e2g9aw0VeIawpuMMUNnpNLKYtr9xo5ASpkgZe
2Vw2qaKSAZDceM4uss6woveg4HdBStf9akKcMFxXm0JQAQQ8SfEv4YL8QZpRYYOt
V6Bfr5FsES1vN6wvYjz7be0gAuU2dZ44zk5tRZRZBD0rPNnKLGyPX0C5vhMzmb74
Ygdd3liBgNnZq1dM2qX26CyMfvMcTyHFaKDVRZ624xQTuLNTst8QmuU+GmEU5Owh
CDUzZmVye1EHoI8y18pvcjY4+ddH5Pj+NiCbAgMBAAECggEAKJ/HcrOaVLjLvlMJ
CLNYf9ts/q00P3rvWljCwY/tH79lPi/FszgaycM9jae1PQrzKNCQxfRdP8FKRLxs
QIz1grBWRzGxloxBPeTnwTw1fNjcT9Sd9ph5c19/jcp/h01MAWXUw/fwKSP2V3cY
Yce5qgFmGGQyX+pebLSpeKAZ6fQqlg1kqpUwD+gQzOD2jwdDtDMw2e+VDjM5CgBw
c/Sp0ULZdJg52CWAPP/+oRDYA1SazMP/UN9BvJHbCGFJU8yMeWpq1vZdAGvw0xEi
yVBvulNSuy84NbudlIS3r/SHoiZ5gBYoarxvQVySPyFUz5qi4NSa1aexLQOIfp9F
31nrIQKBgQC/N6JYkPvsTrwLkdyhknjkvup+Kk9Er9xLXdKltA74zZ3gDYIKw6zW
WWsRRz32kBOOCLRJ81v13/Hp/VPblk0fOyrDXMOiGrx7PlJoBAcisssRwlbyMlsj
9aVB87Y6w2NthfZdqeP1sYy251ujSjhiwC5JkWxDPnXkFZ8pYk9DuQKBgQCWB0cw
Mm15yBDNdHeeM7Q2G4xxt9YQkrcKoQqlXzc5HOv4Dl9z2Y0JCciqjhxSSRSMQ7bM
kgt+FmCNAgcszQNm1ZgiNGkojYzAKJ4TqOhd7gjhlEJ6ko6CSOMinJ8a28ySCSCt
iCLIYKSFRvawo53k59EjfE6WXi3yWq7B2luY8wKBgQCAmD+b5FUnxIN+HGgOWNzo
qOQ2WOnaPZKfNMe3HHJ9LHqh9TOhXZrcxGHWytyG9rnP6Val7P3JSwP74pwErv3K
pRXLMJV4bRaYo37z6NGLqkYa+WFTnD3i3w52XhoNRt2VCGR49QoTb4Y8czcTh5SI
BMqcOypdxqGt9ESHpk+8MQKBgBBTDGjvI8Bbh5CDeryJDbPRfY+mHDHDdQ1ABqOF
38FkJzzfe31kj5cLLH16ZxSBBS5Z89sAjSWUE79kGBnE2Pjohtw57bNO66DBCqpC
9724t9irxCwIJF+CxMwwVJOclP7S5W5cYMSC14P0+8SumPCeVMBSt05fQwqa2WcY
dl0HAoGAZM0K+cmJiHi7AbRD/ooNUYeZ5+H9MYUTeCoADJzHnBHv0FI9ex1jfItl
wp4kbSF1j8ZMjU5cw/OELYiB2CxLzYXQr8zGnIuN5285iQZa6HW5E9K2oeLCw8qT
H3+dq7ZQ/dRAFDrYcWK2YcIp8s90OAggb39KvfIM9qyn+kKHwQA=
-----END RSA PRIVATE KEY-----`

func TestJWTTokenGen_GenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("ParseRSAPrivateKeyFromPEM error:%v\n", err) //会终止
	}
	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	token, err := g.GenerateToken("625905883b4ae1a7857692aa", 2*time.Hour)
	if err != nil {
		t.Errorf("GenerateToken error:%v\n", err) //不终止
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI1OTA1ODgzYjRhZTFhNzg1NzY5MmFhIn0.YvwmhMjACiaS-FA5C5nD2rb1I21IiVDZ2xRUlSbFP03QX--Ba_uWiHxaP-uNsep34fjgSeEVVlcPOnrcbJy_7rx-PNwFD_DYnHR230nJO0zo38nXHg5cNKZ5thT2YNDnmQ1jjRIQoIyUe6ubRmcHq8aNNkdJNrY_DDJgwHX5i-TIs8mzzmUIW_WAs4CjW8weLKA7shAjYI2cEt7hANovnXCZ7xuYk80hk_oI8ISdpCjWigwOw9ZiewF7XYeKWpmKc51--ZPRnXRIUI201xLw2TONuzofVbLztzVQ11gNe7dnv4zyEwCIyx9IKDT1VWWXjUh1YmxHdvE59n7vgCZ2hA"
	if want != token {
		fmt.Println(token)
		t.Fatalf("not this token\n")

	}
	fmt.Println(token)
}
