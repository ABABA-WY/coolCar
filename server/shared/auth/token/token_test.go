package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBITANBgkqhkiG9w0BAQEFAAOCAQ4AMIIBCQKCAQBwEAjRw383hs3ioocpqVYW
ESnaOVAneNJkIB9RCal+Ax0s0JDaPb4eR2lT3fsuRXA2LOKyXYVr2F5e2g9aw0Ve
IawpuMMUNnpNLKYtr9xo5ASpkgZe2Vw2qaKSAZDceM4uss6woveg4HdBStf9akKc
MFxXm0JQAQQ8SfEv4YL8QZpRYYOtV6Bfr5FsES1vN6wvYjz7be0gAuU2dZ44zk5t
RZRZBD0rPNnKLGyPX0C5vhMzmb74Ygdd3liBgNnZq1dM2qX26CyMfvMcTyHFaKDV
RZ624xQTuLNTst8QmuU+GmEU5OwhCDUzZmVye1EHoI8y18pvcjY4+ddH5Pj+NiCb
AgMBAAE=
-----END PUBLIC KEY-----`

func TestJWTTokenVerifier_Verify(t *testing.T) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))

	if err != nil {
		t.Fatalf("cannot ParseRSAPublicKeyFromPEM:%v\n", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pk,
	}
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI1OTA1ODgzYjRhZTFhNzg1NzY5MmFhIn0.YvwmhMjACiaS-FA5C5nD2rb1I21IiVDZ2xRUlSbFP03QX--Ba_uWiHxaP-uNsep34fjgSeEVVlcPOnrcbJy_7rx-PNwFD_DYnHR230nJO0zo38nXHg5cNKZ5thT2YNDnmQ1jjRIQoIyUe6ubRmcHq8aNNkdJNrY_DDJgwHX5i-TIs8mzzmUIW_WAs4CjW8weLKA7shAjYI2cEt7hANovnXCZ7xuYk80hk_oI8ISdpCjWigwOw9ZiewF7XYeKWpmKc51--ZPRnXRIUI201xLw2TONuzofVbLztzVQ11gNe7dnv4zyEwCIyx9IKDT1VWWXjUh1YmxHdvE59n7vgCZ2hA"
	jwt.TimeFunc = func() time.Time { return time.Unix(1516239122, 0) }
	verify, err := v.Verify(token)
	if err != nil {
		t.Fatalf("cannot verify:%v\n", err)
	}
	fmt.Println(verify)
}
