package hcaptcha

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func Verify(secret, response, siteId, ip string) (Response, error) {
	values := url.Values{
		"secret":   []string{secret},
		"response": []string{response},
	}

	if ip != "" {
		values.Set("remoteip", ip)
	}

	if siteId != "" {
		values.Set("sitekey", siteId)
	}

	res, err := http.PostForm("https://hcaptcha.com/siteverify", values)
	if err != nil {
		return Response{}, err
	}

	var result Response
	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		return Response{}, err
	}

	return result, nil
}

type Response struct {
	Success     bool        `json:"success"`               // is the passcode valid, and does it meet security criteria you specified, e.g. sitekey?
	ChallengeTS string      `json:"challenge_ts"`          // timestamp of the challenge (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string      `json:"hostname"`              // the hostname of the site where the challenge was solved
	Credit      *bool       `json:"credit,omitempty"`      // optional: whether the response will be credited
	ErrorCodes  []ErrorCode `json:"error-codes,omitempty"` // optional: any error codes
	Score       float32     `json:"score"`                 // ENTERPRISE feature: a score denoting malicious activity.
	ScoreReason []string    `json:"score_reason,omitempty"`
}

func (r *Response) Error() error {
	if r.Success {
		return nil
	}
	if len(r.ErrorCodes) == 0 {
		return errors.New("unknown error")
	}
	return errors.New(r.ErrorCodes[0].Reason())
}

type ErrorCode string

func (e ErrorCode) Reason() string {
	switch e {
	case "missing-input-secret":
		return "Your secret key is missing."
	case "invalid-input-secret":
		return "Your secret key is invalid or malformed."
	case "missing-input-response":
		return "The response parameter (verification token) is missing."
	case "invalid-input-response":
		return "The response parameter (verification token) is invalid or malformed."
	case "bad-request":
		return "The request is invalid or malformed."
	case "invalid-or-already-seen-response":
		return "The response parameter has already been checked, or has another issue."
	case "not-using-dummy-passcode":
		return "You have used a testing sitekey but have not used its matching secret."
	case "sitekey-secret-mismatch":
		return "The sitekey is not registered with the provided secret."
	case "invalid-remoteip":
		return "Invalid client IP address." // @todo not documented in https://docs.hcaptcha.com/configuration#error-codes
	default:
		return "Unknown"
	}
}
