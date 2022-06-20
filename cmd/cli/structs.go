package main

import (
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Datum struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery Time   `json:"delivery"`
}

type Time struct {
	From time.Time
	To   time.Time
}

// expected format: {weekday} {h}AM - {h}PM
func (t *Time) UnmarshalJSON(data []byte) error {
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	parts := strings.Split(str, " ")

	from, err := time.Parse("3PM", parts[1])
	if err != nil {
		return err
	}
	t.From = from

	to, err := time.Parse("3PM", parts[3])
	if err != nil {
		return err
	}
	t.To = to

	return nil
}

func NewStatsIter(reader io.Reader, postcode string, from time.Time, to time.Time, filters []string) Stats {
	d := jsonfast.NewDecoder(reader)
	var data []*Datum
	err := d.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	stats := Stats{
		UniqueRecipeCount:       0,
		CountPerRecipe:          CountPerRecipe{make(map[string]*RecipeStats, 2e3), make([]*RecipeStats, 0, 2e3)},
		BusiestPostcode:         BusiestPostcode{make(map[string]*PostcodeStats, 1e6), 0, ""},
		CountPerPostcodeAndTime: PostcodeTimeStats{Postcode: postcode, From: from, To: to},
		MatchByName:             MatchByName{Filter: filters},
	}

	for _, datum := range data {
		stats.CountPerRecipe.AddRecipe(datum.Recipe)
		stats.BusiestPostcode.AddPostcode(datum.Postcode)
		stats.CountPerPostcodeAndTime.AddData(datum)
	}

	stats.UniqueRecipeCount = len(stats.CountPerRecipe.UniqueRecipes)
	for _, v := range stats.CountPerRecipe.RecipeStatistics {
		stats.MatchByName.AddRecipe(v.Recipe)
	}

	return stats
}

/* Didn't work so fast
func NewStats(reader io.Reader, postcode string, from time.Time, to time.Time, filterRecipes []string) Stats {
	decoder := json.NewDecoder(reader)

	// expect first byte is square bracket
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	stats := Stats{
		UniqueRecipeCount:       0,
		CountPerRecipe:          CountPerRecipe{make(map[string]*RecipeStats, 2e3), make([]*RecipeStats, 0, 2e3)},
		BusiestPostcode:         BusiestPostcode{make(map[string]*PostcodeStats, 1e6), 0, ""},
		CountPerPostcodeAndTime: PostcodeTimeStats{Postcode: postcode, From: from, To: to},
		MatchByName:             MatchByName{Filter: filterRecipes},
	}

	for decoder.More() {
		datum := &Datum{}
		err := decoder.Decode(datum)
		if err != nil {
			log.Fatal(err)
		}

		stats.CountPerRecipe.AddRecipe(datum.Recipe)
		stats.BusiestPostcode.AddPostcode(datum.Postcode)
		stats.CountPerPostcodeAndTime.AddData(datum)
	}

	stats.UniqueRecipeCount = len(stats.CountPerRecipe.UniqueRecipes)
	for _, v := range stats.CountPerRecipe.RecipeStatistics {
		stats.MatchByName.AddRecipe(v.Recipe)
	}

	return stats
}
*/

type Stats struct {
	UniqueRecipeCount       int               `json:"unique_recipe_count"`
	CountPerRecipe          CountPerRecipe    `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode   `json:"busiest_postcode"`
	CountPerPostcodeAndTime PostcodeTimeStats `json:"count_per_postcode_and_time"`
	MatchByName             MatchByName       `json:"match_by_name"`
}

type CountPerRecipe struct {
	RecipeStatistics map[string]*RecipeStats
	UniqueRecipes    []*RecipeStats
}

func (cpr *CountPerRecipe) AddRecipe(recipe string) {
	if r, ok := cpr.RecipeStatistics[recipe]; ok {
		r.Count += 1
	} else {
		rs := &RecipeStats{recipe, 1}
		cpr.RecipeStatistics[recipe] = rs
		cpr.UniqueRecipes = append(cpr.UniqueRecipes, rs)
	}
}

func (cpr CountPerRecipe) MarshalJSON() ([]byte, error) {
	if cpr.RecipeStatistics == nil {
		return []byte("null"), nil
	}
	sort.SliceStable(cpr.UniqueRecipes, func(i int, j int) bool {
		return cpr.UniqueRecipes[i].Recipe < cpr.UniqueRecipes[j].Recipe
	})
	return jsonfast.MarshalIndent(cpr.UniqueRecipes, "", "    ")
}

type RecipeStats struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type BusiestPostcode struct {
	PostcodeStatistics map[string]*PostcodeStats
	maxDeliveryCount   int
	maxPostcode        string
}

func (bp *BusiestPostcode) AddPostcode(postcode string) {
	if pc, ok := bp.PostcodeStatistics[postcode]; ok {
		pc.DeliveryCount += 1
	} else {
		bp.PostcodeStatistics[postcode] = &PostcodeStats{postcode, 1}
	}
	if bp.PostcodeStatistics[postcode].DeliveryCount > bp.maxDeliveryCount {
		bp.maxDeliveryCount = bp.PostcodeStatistics[postcode].DeliveryCount
		bp.maxPostcode = postcode
	}
}

func (bp BusiestPostcode) MarshalJSON() ([]byte, error) {
	bpstats, ok := bp.PostcodeStatistics[bp.maxPostcode]
	if bp.PostcodeStatistics == nil || !ok {
		return []byte("null"), nil
	}
	return jsonfast.MarshalIndent(bpstats, "", "    ")
}

type PostcodeStats struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type PostcodeTimeStats struct {
	Postcode      string    `json:"postcode"`
	From          time.Time `json:"from"`
	To            time.Time `json:"to"`
	DeliveryCount int       `json:"delivery_count"`
}

func (pts *PostcodeTimeStats) AddData(datum *Datum) {
	if pts.Postcode == datum.Postcode &&
		(datum.Delivery.From == pts.From || datum.Delivery.From.After(pts.From)) &&
		(datum.Delivery.To == pts.To || datum.Delivery.To.Before(pts.To)) {
		pts.DeliveryCount += 1
	}
}

func (pts *PostcodeTimeStats) MarshalJSON() ([]byte, error) {
	type Alias PostcodeTimeStats
	return jsonfast.MarshalIndent(&struct {
		*Alias
		From string `json:"from"`
		To   string `json:"to"`
	}{
		Alias: (*Alias)(pts),
		From:  pts.From.Format("3PM"),
		To:    pts.To.Format("3PM"),
	}, "", "    ")
}

type MatchByName struct {
	Filter []string
	Match  []string
}

func (mbn *MatchByName) AddRecipe(recipe string) {
	for _, f := range mbn.Filter {
		if strings.Contains(recipe, f) {
			mbn.Match = append(mbn.Match, recipe)
			break
		}
	}
}

func (mbn MatchByName) MarshalJSON() ([]byte, error) {
	if mbn.Match == nil {
		return []byte("null"), nil
	}
	sort.SliceStable(mbn.Match, func(i int, j int) bool {
		return mbn.Match[i] < mbn.Match[j]
	})
	return jsonfast.MarshalIndent(mbn.Match, "", "    ")
}
