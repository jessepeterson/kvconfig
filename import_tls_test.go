package kvconfig

import "testing"

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
)

func TestTLSImport(t *testing.T) {
	type TestStruct struct {
		TLSCert *tls.Certificate `kvconfig:"tls_cert_test"`
	}

	rsaCertB64 := "MIIDtTCCAp2gAwIBAgIJAMMMdz7/T5DIMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMTcwMzAzMDY1OTQ5WhcNMTgwMzAzMDY1OTQ5WjBFMQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqMbsFXOgm6B/sOu6Ueo22wMv+DISa372v6/7OOBFRhHZE3ZHD5gucZAwLQ2Rr4lG6D1w4AlgL43k1Xc6VHTfJ9BkbVUg+iW7ihKpGNbFhIC0ve0sdl0eB02cuhXdzkLMfvLliM+ad6uZ4LLpygOVt5+Tux2W0ok+8Us+H7Ghu14CdqHwTsssFycv3PP3ySEfV3NHbJIBjcOTYClCVfKkvgdOPbXimfSwJ/QHrFUdszlk33Vq9KCNclcbWcVPN64a4ou6PmlX9j6SPclD16pC2vMZ+PnXuKbTJiyAIzEFVhNcjpcSufuwBcxuF3bdO0NKgluKEJe2ke0bJWuZm9CzdQIDAQABo4GnMIGkMB0GA1UdDgQWBBRxYa+7s0rTxlT0H0abmgE17qejSjB1BgNVHSMEbjBsgBRxYa+7s0rTxlT0H0abmgE17qejSqFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJAMMMdz7/T5DIMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAAVFhV41ug8tCPLMembVXu4szz3MUftEYu4AvxV2x9/CHs8ueUC+eJSFCnQgPD953f1stmYQ/ZdJvoP20OJK6V7Pq1AoUGES9/hdAc5hFpg99+w5ySEeT9MkuKniT8a1jojV+kA2aJYuPY19KkMXFWnaZxSYcw3sT03iAYZHkG/yMEcpZau0GuB9qYrYPpptiG7gQT4IdRWJaKG4WtO4P6HJAnsfsmcl8janfBIjvuKURTutLcw/6/idAfHN1MFUqguITsJYqYPNCFRqvJ772GLcOZ1SPyNMe6NPcrYV3Wm5ExThiHR0xYWB7OHn/yoclffY1E2j6WLSriVh/sTt3FA="
	rsaKeyB64 := "MIIEpAIBAAKCAQEAqMbsFXOgm6B/sOu6Ueo22wMv+DISa372v6/7OOBFRhHZE3ZHD5gucZAwLQ2Rr4lG6D1w4AlgL43k1Xc6VHTfJ9BkbVUg+iW7ihKpGNbFhIC0ve0sdl0eB02cuhXdzkLMfvLliM+ad6uZ4LLpygOVt5+Tux2W0ok+8Us+H7Ghu14CdqHwTsssFycv3PP3ySEfV3NHbJIBjcOTYClCVfKkvgdOPbXimfSwJ/QHrFUdszlk33Vq9KCNclcbWcVPN64a4ou6PmlX9j6SPclD16pC2vMZ+PnXuKbTJiyAIzEFVhNcjpcSufuwBcxuF3bdO0NKgluKEJe2ke0bJWuZm9CzdQIDAQABAoIBAQCWY7JoJwD8y5YccuAyL64zl3J+CTgKmzaJdek4M/bmSe8Q/XqydZskzCNxcb7YGE2bkWvr5c7UcO5wG+5Y5U8XbgSeu5VH8KlcjeYpYO7dc8YZ2qWczrp8LXczBVsAeNs5X3ySXNK6Qak65JGX1XvqBAKiX+pNrcftQGuZ2DFR/y6dS4dVLPXpEtwDvs4/38ALCS9yYUDgKImrmN2XuApPeQ+EIHJxvU5xJA74zQCWVN9pzHsQyvz1DM6SJBZHDlePQxPt8Q9pTfxfEUSmC64gsH5QcUlE4gwPYibzje3dMZYfnPawsZ4icft5lzTQf4H7IZqT8CEXpZZSvqj41Bz9AoGBAN3p+Mt1/nGEqtY2bcWYaSP9JownJ3UknND8746gj3wP1pncyqQ27MHnYJt1qnDlN54pSglRCryEmtgPKX7B17EGLtC9Y5rbrqyzaCyf7K974TnI7VsKKBbIBM7Y5EOgPyqOcRa5lGa26D10+cTVYZwHV2MLRVMavO9llcoSMHb7AoGBAMKzhJoxgR3q3SgvgVNggZ+Ng7capRjjLBu5ZShaqp5veoKVE88KzvDOWGhdfVhPZZMBmysm1warWS6ioLjo8/z8ntSNxUBrImU8UrjNgJ0tIINf2xDD5BSYEF6MK8+JRAY6Q4wcflCtUm8NjiadqxD1oe9ZRKvpu8kuznOZQzRPAoGAO+c71NhuLgCNCTQ6H5vLzf45GJ49JX8TocqVdB/de7Tezjvuq7Nz58foqS5zKvSFNfmZVbh9uHPnRKmbHu9+pPexTYHCUHw6w73OQjWNc7VyD+IwSGIOfk/SFHAx9htc0cUPu/2ulKeNO4HHJp4fMjo9GaxiM1PFaq42aAzO7l8CgYBrZ86ZpP9+Oobf2Tz1esJm+xETHF7BGOjHLoHQPhvrJMIncQepamP4UUxR3mj8I2h8LSGlL1rlMfcTk+EnwFKAV/dieAa9X5xszlcv3SW7Dx7leiaF3BphBfXZwmeUqDtfWBrVGw7PgJ1957NoOAgbZfV77PnGAD14YRrAiGabXwKBgQCWOkL1LI6C3gVXC+TLjvlsaOdir9ZYnw0TSo57fJ5Nv4GzHc4YuLDnGiSk2r3r5NmtN1YAdNtl/niKmFL8XNzErH/yGUae11cNTdQuTh8zXv44uW0qGGxt0S+KqaszizDxjgN+sA2xQshYgPSnFwjmOZDjp9xtoPVMnVrQwzF3Ng=="

	keyBytes, _ := base64.StdEncoding.DecodeString(rsaKeyB64)
	crtBytes, _ := base64.StdEncoding.DecodeString(rsaCertB64)

	kv := &MapStrStr{
		"tls_cert_test_cert_0": rsaCertB64,
		"tls_cert_test_pk_0":   rsaKeyB64,
	}

	ts := TestStruct{}

	Import(kv, &ts)

	if ts.TLSCert == nil {
		t.Error("TestStruct.TLSCert == nil")
	} else {
		tC := *(ts.TLSCert)
		if len(tC.Certificate) != 1 {
			t.Errorf("len(*(TestStruct.TLSCert).Certificate) = %d; wanted 1", len(tC.Certificate))
		} else {
			cmp := bytes.Compare(tC.Certificate[0], crtBytes)
			if cmp != 0 {
				t.Errorf("*(TestStruct.TLSCert).Certificate[0] does not match wanted (%d)", cmp)
			}
			keyValBytes := x509.MarshalPKCS1PrivateKey(tC.PrivateKey.(*rsa.PrivateKey))
			cmp = bytes.Compare(keyValBytes, keyBytes)
			if cmp != 0 {
				t.Errorf("*(TestStruct.TLSCert).PrivateKey does not match wanted (%d)", cmp)
			}
		}
	}
}
