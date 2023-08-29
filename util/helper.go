package util

import (
	"douyin/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

func GetUser(c *gin.Context) *model.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	} else {
		return user.(*model.User)
	}
}

var allCategories = []string{"dongman", "fengjing", "biying", "meinv"}
var pxTypes = [...]string{"", "m", "pc"}

const nPxTypes = len(pxTypes)

func RandomImageURL(category []string) string {
	if len(category) == 0 {
		category = allCategories
	}
	pxType := pxTypes[rand.Intn(nPxTypes)]

	randomImageURL := url.URL{
		Scheme: "https",
		Host:   "tuapi.eees.cc",
		Path:   "/api.php",
	}
	values := url.Values{
		"type":     {"json"},
		"category": {fmt.Sprintf("{%s}", strings.Join(category, ","))},
	}
	if pxType != "" {
		values.Set("px", pxType)
	}
	randomImageURL.RawQuery = values.Encode()

	u := randomImageURL.String()
	response, err := http.Get(u)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	var randomImageResponse struct {
		Error  string `json:"error"`
		Result string `json:"result"`
		Width  string `json:"width"`
		Height string `json:"height"`
		Format string `json:"format"`
		URL    string `json:"img"`
	}

	err = json.Unmarshal(body, &randomImageResponse)
	if err != nil {
		return ""
	}

	return randomImageResponse.URL
}
