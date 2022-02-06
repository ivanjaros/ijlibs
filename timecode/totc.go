package timecode

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"net/url"
	"time"
)

const QRCodeSize = 400
const DefaultSecretLength = 24

// this is how the totp.Generate() produces secrets,
// using simple randomly generated strings is not working, so we have to do it this way
func GenerateSecret(length ...int) (string, error) {
	ln := DefaultSecretLength
	if len(length) > 0 {
		ln = length[0]
	}
	secret := make([]byte, ln)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// this is a copy of the totp.Generate() with the exception
// of being able to directly provide secret and to use base32 encoding without padding, instead of trimming.
func Generate(issuer, secret, userId string) (*otp.Key, error) {
	v := url.Values{}
	v.Set("secret", secret)
	v.Set("issuer", issuer)
	v.Set("period", "30")
	v.Set("algorithm", otp.AlgorithmSHA1.String()) // most compatible setting
	v.Set("digits", otp.DigitsEight.String())      // 6 feels too little since we are using this instead of password so security matters

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + issuer + ":" + userId,
		RawQuery: v.Encode(),
	}

	return otp.NewKeyFromURL(u.String())
}

// this is a copy of totp.Validate with settings we are using in GenerateTOTP
func Validate(passcode string, secret string) bool {
	rv, _ := totp.ValidateCustom(
		passcode,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsEight,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	return rv
}

// image is base64 encoded byte slice
func GetImageData(key *otp.Key) (url string, image []byte, mimeType string, err error) {
	img, err := key.Image(QRCodeSize, QRCodeSize)
	if err != nil {
		return "", nil, "", err
	}

	buff := bytes.NewBuffer(nil)
	b64 := base64.NewEncoder(base64.StdEncoding, buff)
	if err := png.Encode(b64, img); err != nil {
		return "", nil, "", err
	}

	return key.URL(), buff.Bytes(), "image/png", nil
}
