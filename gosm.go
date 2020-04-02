package gosm

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FrankSantoso/gosm/fetch"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type OsmData struct {
	OsmType           string      `json:"osm_type"`
	OsmID             int64       `json:"osm_id"`
	Class             string      `json:"class"`
	Type              string      `json:"type"`
	AdminLevel        int64       `json:"admin_level"`
	RankSearch        int64       `json:"rank_search"`
	RankAddress       int64       `json:"rank_address"`
	PlaceID           int64       `json:"place_id"`
	ParentPlaceID     int64       `json:"parent_place_id"`
	Housenumber       int64       `json:"housenumber"`
	CountryCode       string      `json:"country_code"`
	Langaddress       string      `json:"langaddress"`
	Placename         string      `json:"placename"`
	Ref               string      `json:"ref"`
	Lon               string      `json:"lon"`
	Lat               string      `json:"lat"`
	Importance        float64     `json:"importance"`
	Addressimportance string      `json:"addressimportance"`
	ExtraPlace        interface{} `json:"extra_place"`
	Addresstype       string      `json:"addresstype"`
	ABoundingBox      []string    `json:"aBoundingBox"`
	Label             string      `json:"label"`
	Name              string      `json:"name"`
	Foundorder        float64     `json:"foundorder"`
}

func GetGeo(s string) ([]OsmData, error) {
	ctx := context.Background()
	body, err := fetch.Fetch(ctx, s)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var osmDatas []OsmData
	errchan := make(chan error, 1)
	doc.Find("script[type='text/javascript']").Each(func(i int, s *goquery.Selection) {
		sel, err := removesIrrelevantStrings(s)
		if err != nil {
			errchan <- err
		}
		if err = json.Unmarshal(sel, &osmDatas); err != nil {
			errchan <- err
		}
		close(errchan)
	})
	if err = <-errchan; err != nil {
		return nil, err
	}
	return osmDatas, nil
}

func removesIrrelevantStrings(s *goquery.Selection) ([]byte, error) {
	selString := strings.Split(s.Text(), ";")
	if len(selString) == 0 {
		return nil, errors.New("No location found")
	}
	// removes irrelevant bits
	return []byte(strings.ReplaceAll(
		selString[1],
		"var nominatim_results = ",
		"",
	)), nil

}
