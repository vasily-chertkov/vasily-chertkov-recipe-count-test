package main

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {

	am1, _ := time.Parse("3PM", "1AM")
	pm7, _ := time.Parse("3PM", "7PM")
	pm9, _ := time.Parse("3PM", "9PM")
	am12, _ := time.Parse("3PM", "12AM")
	pm12, _ := time.Parse("3PM", "12PM")

	var tests = []struct {
		name   string
		input  string
		expect Time
	}{
		{
			"Correct Times",
			strconv.Quote("Wednesday 1AM - 7PM"),
			Time{
				From: am1,
				To:   pm7,
			},
		},
		{
			"Second half",
			strconv.Quote("Thursday 12PM - 9PM"),
			Time{
				From: pm12,
				To:   pm9,
			},
		},
		{
			"Night",
			strconv.Quote("Friday 12AM - 1AM"),
			Time{
				From: am12,
				To:   am1,
			},
		},
	}

	for _, test := range tests {
		var tm Time
		err := jsonfast.Unmarshal([]byte(test.input), &tm)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.expect, tm, test.name)
	}
}

func TestStats(t *testing.T) {
	b, err := ioutil.ReadFile("fixture.json")
	assert.Nil(t, err, "Could not open fixture file")

	am2, _ := time.Parse("3PM", "2AM")
	pm3, _ := time.Parse("3PM", "3PM")

	CajunSpicedPulledPork := &RecipeStats{"Cajun-Spiced Pulled Pork", 15}
	CheesyChickenEnchiladaBake := &RecipeStats{"Cheesy Chicken Enchilada Bake", 5}
	CherryBalsamicPorkChops := &RecipeStats{"Cherry Balsamic Pork Chops", 7}
	ChickenPineappleQuesadillas := &RecipeStats{"Chicken Pineapple Quesadillas", 9}
	ChickenSausagePizzas := &RecipeStats{"Chicken Sausage Pizzas", 9}
	CreamyDillChicken := &RecipeStats{"Creamy Dill Chicken", 10}
	CreamyShrimpTagliatelle := &RecipeStats{"Creamy Shrimp Tagliatelle", 7}
	CrispyCheddarFricoCheeseburgers := &RecipeStats{"Crispy Cheddar Frico Cheeseburgers", 2}
	GardenQuesadillas := &RecipeStats{"Garden Quesadillas", 5}
	GarlicHerbButterSteak := &RecipeStats{"Garlic Herb Butter Steak", 4}
	GrilledCheeseAndVeggieJumble := &RecipeStats{"Grilled Cheese and Veggie Jumble", 7}
	HeartyPorkChili := &RecipeStats{"Hearty Pork Chili", 4}
	HoneySesameChicken := &RecipeStats{"Honey Sesame Chicken", 6}
	HotHoneyBarbecueChickenLegs := &RecipeStats{"Hot Honey Barbecue Chicken Legs", 5}
	KoreanStyleChickenThighs := &RecipeStats{"Korean-Style Chicken Thighs", 7}
	MeatloafALaMom := &RecipeStats{"Meatloaf à La Mom", 7}
	MediterraneanBakedVeggies := &RecipeStats{"Mediterranean Baked Veggies", 5}
	MeltyMontereyJackBurgers := &RecipeStats{"Melty Monterey Jack Burgers", 7}
	MoleSpicedBeefTacos := &RecipeStats{"Mole-Spiced Beef Tacos", 4}
	OnePanOrzoItaliano := &RecipeStats{"One-Pan Orzo Italiano", 6}
	ParmesanCrustedPorkTenderloin := &RecipeStats{"Parmesan-Crusted Pork Tenderloin", 8}
	SpanishOnePanChicken := &RecipeStats{"Spanish One-Pan Chicken", 6}
	SpeedySteakFajitas := &RecipeStats{"Speedy Steak Fajitas", 13}
	SpinachArtichokePastaBake := &RecipeStats{"Spinach Artichoke Pasta Bake", 9}
	SteakhouseStyleNewYorkStrip := &RecipeStats{"Steakhouse-Style New York Strip", 6}
	StovetopMacNCheese := &RecipeStats{"Stovetop Mac 'N' Cheese", 3}
	SweetApplePorkTenderloin := &RecipeStats{"Sweet Apple Pork Tenderloin", 10}
	TexMexTilapia := &RecipeStats{"Tex-Mex Tilapia", 4}
	YellowSquashFlatbreads := &RecipeStats{"Yellow Squash Flatbreads", 9}

	countPerRecipe := CountPerRecipe{
		RecipeStatistics: map[string]*RecipeStats{
			"Cajun-Spiced Pulled Pork":           CajunSpicedPulledPork,
			"Cheesy Chicken Enchilada Bake":      CheesyChickenEnchiladaBake,
			"Cherry Balsamic Pork Chops":         CherryBalsamicPorkChops,
			"Chicken Pineapple Quesadillas":      ChickenPineappleQuesadillas,
			"Chicken Sausage Pizzas":             ChickenSausagePizzas,
			"Creamy Dill Chicken":                CreamyDillChicken,
			"Creamy Shrimp Tagliatelle":          CreamyShrimpTagliatelle,
			"Crispy Cheddar Frico Cheeseburgers": CrispyCheddarFricoCheeseburgers,
			"Garden Quesadillas":                 GardenQuesadillas,
			"Garlic Herb Butter Steak":           GarlicHerbButterSteak,
			"Grilled Cheese and Veggie Jumble":   GrilledCheeseAndVeggieJumble,
			"Hearty Pork Chili":                  HeartyPorkChili,
			"Honey Sesame Chicken":               HoneySesameChicken,
			"Hot Honey Barbecue Chicken Legs":    HotHoneyBarbecueChickenLegs,
			"Korean-Style Chicken Thighs":        KoreanStyleChickenThighs,
			"Meatloaf à La Mom":                  MeatloafALaMom,
			"Mediterranean Baked Veggies":        MediterraneanBakedVeggies,
			"Melty Monterey Jack Burgers":        MeltyMontereyJackBurgers,
			"Mole-Spiced Beef Tacos":             MoleSpicedBeefTacos,
			"One-Pan Orzo Italiano":              OnePanOrzoItaliano,
			"Parmesan-Crusted Pork Tenderloin":   ParmesanCrustedPorkTenderloin,
			"Spanish One-Pan Chicken":            SpanishOnePanChicken,
			"Speedy Steak Fajitas":               SpeedySteakFajitas,
			"Spinach Artichoke Pasta Bake":       SpinachArtichokePastaBake,
			"Steakhouse-Style New York Strip":    SteakhouseStyleNewYorkStrip,
			"Stovetop Mac 'N' Cheese":            StovetopMacNCheese,
			"Sweet Apple Pork Tenderloin":        SweetApplePorkTenderloin,
			"Tex-Mex Tilapia":                    TexMexTilapia,
			"Yellow Squash Flatbreads":           YellowSquashFlatbreads,
		},
		UniqueRecipes: []*RecipeStats{
			CajunSpicedPulledPork, CheesyChickenEnchiladaBake, CherryBalsamicPorkChops, ChickenPineappleQuesadillas,
			ChickenSausagePizzas, CreamyDillChicken, CreamyShrimpTagliatelle, CrispyCheddarFricoCheeseburgers,
			GardenQuesadillas, GarlicHerbButterSteak, GrilledCheeseAndVeggieJumble, HeartyPorkChili, HoneySesameChicken,
			HotHoneyBarbecueChickenLegs, KoreanStyleChickenThighs, MeatloafALaMom, MediterraneanBakedVeggies,
			MeltyMontereyJackBurgers, MoleSpicedBeefTacos, OnePanOrzoItaliano, ParmesanCrustedPorkTenderloin,
			SpanishOnePanChicken, SpeedySteakFajitas, SpinachArtichokePastaBake, SteakhouseStyleNewYorkStrip,
			StovetopMacNCheese, SweetApplePorkTenderloin, TexMexTilapia, YellowSquashFlatbreads,
		},
	}

	busiestPostcode := BusiestPostcode{
		maxDeliveryCount: 8,
		maxPostcode:      "10186",
		PostcodeStatistics: map[string]*PostcodeStats{
			"10115": {"10115", 1},
			"10116": {"10116", 2},
			"10117": {"10117", 1},
			"10118": {"10118", 2},
			"10119": {"10119", 2},
			"10120": {"10120", 4},
			"10121": {"10121", 2},
			"10122": {"10122", 1},
			"10123": {"10123", 2},
			"10124": {"10124", 1},
			"10126": {"10126", 1},
			"10127": {"10127", 1},
			"10128": {"10128", 1},
			"10129": {"10129", 4},
			"10130": {"10130", 2},
			"10131": {"10131", 1},
			"10133": {"10133", 2},
			"10134": {"10134", 1},
			"10135": {"10135", 2},
			"10136": {"10136", 3},
			"10137": {"10137", 2},
			"10138": {"10138", 1},
			"10139": {"10139", 3},
			"10140": {"10140", 1},
			"10141": {"10141", 2},
			"10142": {"10142", 1},
			"10143": {"10143", 1},
			"10145": {"10145", 1},
			"10146": {"10146", 1},
			"10147": {"10147", 4},
			"10148": {"10148", 1},
			"10149": {"10149", 4},
			"10150": {"10150", 1},
			"10152": {"10152", 2},
			"10153": {"10153", 3},
			"10154": {"10154", 1},
			"10158": {"10158", 2},
			"10159": {"10159", 6},
			"10162": {"10162", 2},
			"10163": {"10163", 2},
			"10164": {"10164", 1},
			"10166": {"10166", 2},
			"10167": {"10167", 1},
			"10168": {"10168", 3},
			"10170": {"10170", 5},
			"10171": {"10171", 2},
			"10172": {"10172", 1},
			"10173": {"10173", 4},
			"10175": {"10175", 1},
			"10176": {"10176", 1},
			"10178": {"10178", 2},
			"10179": {"10179", 2},
			"10180": {"10180", 2},
			"10182": {"10182", 2},
			"10183": {"10183", 1},
			"10184": {"10184", 6},
			"10185": {"10185", 2},
			"10186": {"10186", 8},
			"10187": {"10187", 2},
			"10189": {"10189", 3},
			"10190": {"10190", 4},
			"10192": {"10192", 3},
			"10193": {"10193", 1},
			"10194": {"10194", 5},
			"10196": {"10196", 1},
			"10197": {"10197", 4},
			"10198": {"10198", 3},
			"10200": {"10200", 1},
			"10201": {"10201", 1},
			"10202": {"10202", 3},
			"10203": {"10203", 1},
			"10204": {"10204", 2},
			"10205": {"10205", 3},
			"10206": {"10206", 2},
			"10208": {"10208", 3},
			"10209": {"10209", 2},
			"10210": {"10210", 2},
			"10211": {"10211", 2},
			"10213": {"10213", 5},
			"10214": {"10214", 1},
			"10215": {"10215", 3},
			"10216": {"10216", 5},
			"10217": {"10217", 3},
			"10218": {"10218", 4},
			"10222": {"10222", 3},
			"10223": {"10223", 3},
			"10224": {"10224", 1},
		},
	}

	type input struct {
		postcode string
		from     time.Time
		to       time.Time
		filters  []string
	}

	var tests = []struct {
		name   string
		input  input
		expect Stats
	}{
		{
			"Checking stats",
			input{
				postcode: "10120",
				from:     am2,
				to:       pm3,
				filters:  []string{"Yellow", "Mac"},
			},
			Stats{
				UniqueRecipeCount: 29,
				CountPerRecipe:    countPerRecipe,
				BusiestPostcode:   busiestPostcode,
				CountPerPostcodeAndTime: PostcodeTimeStats{
					Postcode:      "10120",
					From:          am2,
					To:            pm3,
					DeliveryCount: 1,
				},
			},
		},
	}

	for _, test := range tests {
		r := bytes.NewReader(b)
		stats := NewStatsIter(r, test.input.postcode, test.input.from, test.input.to, test.input.filters)
		assert.EqualValues(t, test.expect.UniqueRecipeCount, stats.UniqueRecipeCount, test.name+": Unique recipe count")

		sort.SliceStable(test.expect.CountPerRecipe.UniqueRecipes, func(i int, j int) bool {
			return test.expect.CountPerRecipe.UniqueRecipes[i].Recipe < test.expect.CountPerRecipe.UniqueRecipes[j].Recipe
		})
		sort.SliceStable(stats.CountPerRecipe.UniqueRecipes, func(i int, j int) bool {
			return stats.CountPerRecipe.UniqueRecipes[i].Recipe < stats.CountPerRecipe.UniqueRecipes[j].Recipe
		})
		assert.EqualValues(t, test.expect.CountPerRecipe.UniqueRecipes, stats.CountPerRecipe.UniqueRecipes, test.name+": Count per recipe")

		assert.EqualValues(t, test.expect.BusiestPostcode, stats.BusiestPostcode, test.name+": Busiest postcode")

		assert.EqualValues(t, test.expect.CountPerPostcodeAndTime.DeliveryCount, stats.CountPerPostcodeAndTime.DeliveryCount, test.name+": Count per postcode and time")

	}
}
