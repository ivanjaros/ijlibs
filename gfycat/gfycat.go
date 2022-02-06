package gfycat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type GrantRequest struct {
	Type         string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type GrantResponse struct {
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type SearchResponse struct {
	Gifs    []ResponseGif `json:"gfycats"`
	Cursor  string        `json:"cursor"`
	Related []string      `json:"related"` // related words or suffixes(include forward space)
	Found   int           `json:"found"`
}

// the API returns inconsistent data types(int/string, 0, "0", "") and cannot be properly unmarshaled.
// Luckily we just want the ID and nothing more because we can construct URLs as needed and we do not need dimensions.
type ResponseGif struct {
	//Max2MbGif string `json:"max2mbGif"`
	//UserData  struct {
	//	Name            string `json:"name"`
	//	ProfileImageURL string `json:"profileImageUrl"`
	//	URL             string `json:"url"`
	//	Username        string `json:"username"`
	//	Followers       int    `json:"followers"`
	//	Subscription    int    `json:"subscription"`
	//	Following       int    `json:"following"`
	//	ProfileURL      string `json:"profileUrl"`
	//	Views           int    `json:"views"`
	//	Verified        bool   `json:"verified"`
	//} `json:"userData"`
	//Rating              string        `json:"rating"`
	//Source              int           `json:"source"`
	//FrameRate           int           `json:"frameRate"`
	//Sitename            string        `json:"sitename"`
	//Likes               string        `json:"likes"`
	//Height              int           `json:"height"`
	//UserProfileImageURL string        `json:"userProfileImageUrl"`
	//AvgColor            string        `json:"avgColor"`
	//Dislikes            int           `json:"dislikes"`
	//Published           int           `json:"published"`
	//Gif100Px            string        `json:"gif100px"`
	//Thumb100PosterURL   string        `json:"thumb100PosterUrl"`
	//Tags                []string      `json:"tags"`
	//GifURL              string        `json:"gifUrl"`
	//GfyNumber           string        `json:"gfyNumber"`
	//Mp4Size             int           `json:"mp4Size"`
	//LanguageCategories  []string      `json:"languageCategories"`
	//Max5MbGif           string        `json:"max5mbGif"`
	//GfySlug             string        `json:"gfySlug"`
	//Description         string        `json:"description"`
	//WebpURL             string        `json:"webpUrl"`
	//Title               string        `json:"title"`
	//DomainWhitelist     []string      `json:"domainWhitelist"`
	//Gatekeeper          int           `json:"gatekeeper"`
	//HasTransparency     bool          `json:"hasTransparency"`
	//PosterURL           string        `json:"posterUrl"`
	//MobilePosterURL     string        `json:"mobilePosterUrl"`
	//WebmSize            int           `json:"webmSize"`
	//MobileURL           string        `json:"mobileUrl"`
	GfyName string `json:"gfyName"`
	//Views               int           `json:"views"`
	//CreateDate          int           `json:"createDate"`
	//WebmURL             string        `json:"webmUrl"`
	//HasAudio            bool          `json:"hasAudio"`
	//ExtraLemmas         string        `json:"extraLemmas"`
	//Nsfw                string        `json:"nsfw"`
	//LanguageText2       string        `json:"languageText2"`
	//UserDisplayName     string        `json:"userDisplayName"`
	//MiniURL             string        `json:"miniUrl"`
	//UserName            string        `json:"userName"`
	//Max1MbGif           string        `json:"max1mbGif"`
	//GfyID               string        `json:"gfyId"`
	//NumFrames           int           `json:"numFrames"`
	//Curated             int           `json:"curated"`
	//MiniPosterURL       string        `json:"miniPosterUrl"`
	//Width               int           `json:"width"`
	//Mp4URL              string        `json:"mp4Url"`
	//Md5                 string        `json:"md5"`
	//ContentUrls         struct {
	//	Max2MbGif struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"max2mbGif"`
	//	Webp struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"webp"`
	//	Max1MbGif struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"max1mbGif"`
	//	One00PxGif struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"100pxGif"`
	//	MobilePoster struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"mobilePoster"`
	//	Mp4 struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"mp4"`
	//	Webm struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"webm"`
	//	Max5MbGif struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"max5mbGif"`
	//	LargeGif struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"largeGif"`
	//	Mobile struct {
	//		URL    string `json:"url"`
	//		Size   int    `json:"size"`
	//		Height int    `json:"height"`
	//		Width  int    `json:"width"`
	//	} `json:"mobile"`
	//} `json:"content_urls"`
}

func GetToken(clientId, secret string) (token string, expires_in int, err error) {
	req := GrantRequest{
		Type:         "client_credentials",
		ClientId:     clientId,
		ClientSecret: secret,
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(req); err != nil {
		return "", 0, err
	}

	// https://developers.gfycat.com/api/#authentication
	r, err := http.Post("https://api.gfycat.com/v1/oauth/token", "application/json", buf)
	if err != nil {
		return "", 0, err
	}

	var res GrantResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return "", 0, err
	}

	if res.TokenType != "bearer" {
		return "", 0, errors.New("unexpected token type " + res.TokenType)
	}

	return res.AccessToken, res.ExpiresIn, nil
}

func Search(token, query string, limit int, cursor string) ([]string, string, error) {
	u, _ := url.Parse("https://api.gfycat.com/v1/gfycats/search")
	q := u.Query()
	q.Set("search_text", query)
	q.Set("count", strconv.Itoa(limit))
	q.Set("cursor", cursor)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("authorization", "Bearer "+token)

	// this might not error. for example when search_text was missing, returned response was:
	// {"errorMessage":{"code":"Bad Request","description":"Missing parameter: search_text"}}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	if r.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		return nil, "", fmt.Errorf("gfycat service returned %s and response: %s", r.Status, string(msg))
	}

	defer r.Body.Close()

	var res SearchResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, "", err
	}

	gifs := make([]string, len(res.Gifs))
	for k, v := range res.Gifs {
		gifs[k] = v.GfyName
	}

	return gifs, res.Cursor, nil
}
