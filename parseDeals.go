package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type gameDeal struct {
	Title              string  `json:"title"`
	OriginalPrice      float64 `json:"originalPrice"`
	DiscountPrice      float64 `json:"discountPrice"`
	ThumbnailURL       string  `json:"thumbnailURL"`
	SteamRatingPercent string  `json:"steamRating"`
}

func parseGameDeals(c chan []gameDeal, maxPrice string, minPrice string, title string) {
	requestString := fmt.Sprintf("https://www.cheapshark.com/api/1.0/deals?storeID=1&lowerPrice=%s&upperPrice=%s&onSale=1&title=%s", minPrice, maxPrice, title)
	fmt.Println(requestString)
	resp, err := http.Get(requestString)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}
	var respArray []interface{}
	err = json.Unmarshal(bodyBytes, &respArray)
	if err != nil {
		return
	}
	var dealsSlice []gameDeal
	for _, d := range respArray {
		dealMap := d.(map[string]interface{})
		normalPrice, _ := strconv.ParseFloat(dealMap["normalPrice"].(string), 64)
		salePrice, _ := strconv.ParseFloat(dealMap["salePrice"].(string), 64)
		thumb := dealMap["thumb"].(string)
		splitRes := strings.Split(thumb, "/")
		thumbnailLink := strings.Replace(splitRes[6], "capsule_sm_120", "header", 1)
		steamRating := dealMap["steamRatingPercent"].(string)
		if steamRating == "0" {
			steamRating = "N/A"
		}
		thumbnailLink = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", splitRes[0], splitRes[1], splitRes[2], splitRes[3], splitRes[4], splitRes[5], thumbnailLink)

		gd := gameDeal{
			dealMap["title"].(string),
			normalPrice,
			salePrice,
			thumbnailLink,
			steamRating,
		}
		dealsSlice = append(dealsSlice, gd)
	}
	c <- dealsSlice
}
