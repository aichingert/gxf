package test

import (
	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/drawing"
	"github.com/aichingert/dxf/pkg/entity"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/aichingert/dxf"
)

var (
	plan *drawing.Dxf
	err  error
)

func TestMain(m *testing.M) {
	plan, err = dxf.Open("test.dxf")

	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestCorrectNumberOfInserts(t *testing.T) {
	const expected = 3
	ok := assert.Len(t, plan.Inserts, expected)

	if !ok {
		t.Fatal("expected ", expected, " inserts in test.dxf file: actual", len(plan.Inserts))
	}
}

func TestParsedInsertsCorrectly(t *testing.T) {
	expected := []*entity.Insert{
		{BlockName: "STANDARD_RZ_PR"},
		{BlockName: "STANDARD_SL", Attributes: []*entity.Attrib{{
			Entity: &entity.EntityData{Handle: 65932, Owner: 65924, LayerName: "0"},
			Text: &entity.Text{
				Entity:      &entity.EntityData{},
				Style:       "STANDARD",
				Flags:       1,
				Thickness:   1,
				XScale:      1,
				Rotation:    270,
				Height:      0.0025,
				Vector:      [3]float64{0, 0, 0},
				Coordinates: [3]float64{37.312, 33.3717, 0},
			},
			Tag:   "System",
			Flags: 1,
		}}},
		{BlockName: "QRCode"},
	}

	for i, insert := range plan.Inserts {
		assert.Equal(t, expected[i].BlockName, insert.BlockName)

		if expected[i].Attributes != nil {
			assert.Equal(t, expected[i].Attributes, insert.Attributes)
		}
	}
}

func TestCorrectNumberOfLayoutTableParsed(t *testing.T) {
	const expected = 27
	ok := assert.Equal(t, len(plan.Layers), expected)

	if !ok {
		t.Fatal("expected ", expected, " layers in test.dxf file: actual", len(plan.Layers))
	}
}

func TestCorrectLayoutValuesParsed(t *testing.T) {
	expected := map[string][]int64{
		"Defpoints":        {7, 0},
		"XRef":             {8, 0},
		"REMARKS":          {190, 0},
		"0":                {7, 0},
		"10_":              {114, 2140502},
		"10_yellow":        {40, 0},
		"10_white":         {7, 0},
		"11_system":        {114, 2140502},
		"12_modules":       {114, 2140502},
		"30_color_green":   {114, 2140502},
		"31_color_white":   {255, 0},
		"32_color_black":   {250, 0},
		"33_color_yellow":  {40, 0},
		"34_color_red":     {10, 0},
		"34_color_sperr":   {10, 0},
		"40_UMRISS":        {250, 0},
		"50_type_name":     {114, 2140502},
		"51_typenumber":    {114, 2140502},
		"52_info":          {114, 2140502},
		"53_additional_1":  {114, 2140502},
		"54_additional_2":  {114, 2140502},
		"61_labeling":      {114, 2140502},
		"63_luminaire_ID":  {114, 2140502},
		"70_legend_green":  {114, 2140502},
		"70_legend_white":  {7, 0},
		"70_legend_yellow": {40, 0},
		"99_general":       {114, 2140502},
	}

	for name, layer := range plan.Layers {
		actual := []int64{layer.Color, layer.TrueColor}
		assert.Equal(t, expected[name], actual)
		delete(expected, name)
	}

	if len(expected) > 0 {
		t.Fatal("did not find all layers in test.dxf remaining: ", len(expected))
	}
}

func TestCorrectAmountOfBlocksPresent(t *testing.T) {
	const expected = 2
	ok := assert.Equal(t, len(plan.Blocks), expected)

	if !ok {
		t.Fatal("expected ", expected, " blocks in test.dxf file: actual", len(plan.Blocks))
	}
}

func TestCorrectBlockValuesParsed(t *testing.T) {
	expected := map[string]any{
		"STANDARD_SL": blocks.Block{EntitiesData: &entity.EntitiesData{
			Lines: []*entity.Line{
				{
					Entity: &entity.EntityData{Handle: 240, Owner: 183, LayerName: "32_color_black"},
					Src:    [3]float64{-326.418307086611, 0.000000000001, 0},
					Dst:    [3]float64{326.418307086616, -329.982283388572, 0},
				},
				{
					Entity: &entity.EntityData{Handle: 241, Owner: 183, LayerName: "32_color_black"},
					Src:    [3]float64{326.418307086616, -0.000000000001, 0},
					Dst:    [3]float64{-326.418307086626, -329.982283388572, 0},
				},
			},
			Polylines: []*entity.Polyline{
				{
					Entity: &entity.EntityData{Handle: 237, Owner: 183, LayerName: "32_color_black"},
					Flag:   1, Coordinates: []entity.PLine{
						{X: -326.418307086626, Y: -329.982283388572},
						{X: -326.418307086611, Y: 0.000000000001},
						{X: 326.418307086616, Y: -0.000000000001},
						{X: 326.418307086616, Y: -329.982283388572},
					},
				},
				{
					Entity: &entity.EntityData{Handle: 242, Owner: 183, LayerName: "40_UMRISS"},
					Flag:   1, Coordinates: []entity.PLine{
						{X: -326.418307086611, Y: 0.000000000001},
						{X: 326.418307086616, Y: 0.000000000001},
						{X: 326.418307086616, Y: -329.982283388572},
						{X: -326.418307086611, Y: -329.982283388572},
					},
				},
			},
			Hatches: []*entity.Hatch{
				{
					Entity:      &entity.EntityData{Handle: 232, Owner: 183, LayerName: "0"},
					PatternName: "SOLID",
					SolidFill:   1,
					Style:       1,
					Pattern:     1,
					SeedPoint:   [3]float64{118.6771251745169, -182.3154832900236, 0},
					BoundaryPaths: []*entity.BoundaryPath{
						{
							Flag: 1,
							Lines: []*entity.Line{
								{
									Src: [3]float64{-0.000000000000739, -164.9911416942843, 0},
									Dst: [3]float64{-326.4183070866112, 0.00000000000108, 0},
								},
								{
									Src: [3]float64{-326.4183070866112, 0.0000000000011084, 0},
									Dst: [3]float64{-326.4183070866255, -329.9822833885714, 0},
								},
								{
									Src: [3]float64{-326.4183070866256, -329.9822833885715, 0},
									Dst: [3]float64{-0.000000000000739, -164.9911416942844, 0},
								},
							},
						},
						{
							Flag: 1,
							Lines: []*entity.Line{
								{
									Src: [3]float64{326.4183070866164, -329.9822833885715, 0},
									Dst: [3]float64{326.4183070866163, -0.0000000000012506, 0},
								},
								{
									Src: [3]float64{326.4183070866163, -0.0000000000011653, 0},
									Dst: [3]float64{-0.000000000000739, -164.9911416942844, 0},
								},
								{
									Src: [3]float64{-0.0000000000006821, -164.9911416942844, 0},
									Dst: [3]float64{326.4183070866163, -329.9822833885715, 0},
								},
							},
						},
					},
				},
				{
					Entity:      &entity.EntityData{Handle: 235, Owner: 183, LayerName: "31_color_white"},
					SolidFill:   1,
					PatternName: "SOLID",
					Style:       1,
					Pattern:     1,
					SeedPoint:   [3]float64{12.60739704665474, -302.8154203338559, 0},
					BoundaryPaths: []*entity.BoundaryPath{
						{
							Flag: 1,
							Lines: []*entity.Line{
								{
									Src: [3]float64{-326.4183070866112, 0.0000000000011084, 0},
									Dst: [3]float64{0.0000000000000568, -164.9911416942839, 0},
								},
								{
									Src: [3]float64{0, -164.9911416942839, 0},
									Dst: [3]float64{326.4183070866163, -0.0000000000011369, 0},
								},
								{
									Src: [3]float64{326.4183070866163, -0.0000000000011653, 0},
									Dst: [3]float64{-326.4183070866114, 0.0000000000010775, 0},
								},
							},
						},
						{
							Flag: 1,
							Lines: []*entity.Line{
								{
									Src: [3]float64{326.4183070866163, -329.9822833885715, 0},
									Dst: [3]float64{0.0000000000000568, -164.991141694284, 0},
								},
								{
									Src: [3]float64{0.0000000000000568, -164.9911416942839, 0},
									Dst: [3]float64{-326.4183070866256, -329.9822833885714, 0},
								},
								{
									Src: [3]float64{-326.4183070866256, -329.9822833885715, 0},
									Dst: [3]float64{326.4183070866163, -329.9822833885715, 0},
								},
							},
						},
					},
				},
			},
		}},
		"STANDARD_RZ_PR": blocks.Block{EntitiesData: &entity.EntitiesData{
			Hatches: []*entity.Hatch{
				{
					Entity:      &entity.EntityData{Handle: 285, Owner: 254, LayerName: "0"},
					PatternName: "SOLID",
					SolidFill:   1,
					Style:       1,
					Pattern:     1,
					PixelSize:   0.551072650145,
					SeedPoint:   [3]float64{5.87295002037331, -73.53701621187813, 0},
					BoundaryPaths: []*entity.BoundaryPath{
						{
							Flag: 7,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -326.4183070866112, Y: 0.0000000000011369},
									{X: -326.4183070866112, Y: -329.9779836250452},
									{X: 328.0590986456044, Y: -329.9779836250452},
									{X: 328.0590986455917, Y: 0.0000000000038654},
								},
							},
						},
						{
							Flag: 22,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: 63.44523592602116, Y: -60.49372612885668},
									{X: 150.5379367708664, Y: -60.49372612885668},
									{X: 257.0903297166986, Y: -161.7184994274128},
									{X: 150.5379367708665, Y: -262.9432727259716},
									{X: 63.44523592602116, Y: -262.9432727259679},
									{X: 138.4265494805021, Y: -191.7110248492088},
									{X: 18.45644779332895, Y: -191.7110248492088},
									{X: 18.45644779332895, Y: -131.7259740056195},
									{X: 138.4265494805094, Y: -131.7259740056195},
								},
							},
						},
						{
							Flag: 22,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -76.17227789834044, Y: -136.8709992295896},
									{X: -104.25568688113, Y: -136.8709992295896, Bulge: -0.2399524395792842},
									{X: -107.287691577128, Y: -135.3270278762274},
									{X: -130.542887456965, Y: -103.3511335414381, Bulge: 0.2399524395792833},
									{X: -145.7029109369615, Y: -95.63127677462376},
									{X: -212.2326049179741, Y: -95.63127677462376, Bulge: 0.2840790438404092},
									{X: -232.352202615976, Y: -108.0658719919629},
									{X: -248.629299073652, Y: -140.6200649073141, Bulge: 1},
									{X: -231.8629676586471, Y: -149.0032306148166},
									{X: -215.5858712009779, Y: -116.449037699469, Bulge: -0.2840790438404418},
									{X: -212.2326049179772, Y: -114.3766051632454},
									{X: -189.7382108516359, Y: -114.3766051632454},
									{X: -203.5597649909132, Y: -152.0673701953252, Bulge: 0.0687858006254981},
									{X: -206.8321992824213, Y: -167.1157506353802},
									{X: -210.7281762729564, Y: -217.7634515122316, Bulge: -0.3919022311992417},
									{X: -214.466199023372, Y: -221.2249769783855},
									{X: -263.6255617845537, Y: -221.2249769783855, Bulge: 0.3631588356135515},
									{X: -294.9924956785818, Y: -247.4684367224543},
									{X: -263.625561784551, Y: -247.4684367224561},
									{X: -214.4661990233739, Y: -247.4684367224561, Bulge: 0.3919022311992351},
									{X: -184.5620170200571, Y: -219.776232993233},
									{X: -181.2234173657872, Y: -176.3744374878268},
									{X: -141.261013849475, Y: -273.7803472622229, Bulge: 0.3042679099072648},
									{X: -127.3869924453993, Y: -283.0845606608361},
									{X: -109.0775649966, Y: -283.0845606608361},
									{X: -156.3728768002399, Y: -167.8051362784292, Bulge: -0.1873489039227227},
									{X: -156.4242319135346, Y: -165.0913583761178},
									{X: -142.8859513795958, Y: -128.1730679001891, Bulge: -0.5998762433504838},
									{X: -136.3340862192784, Y: -127.2587391500678},
									{X: -122.4477150571095, Y: -146.3524994980441, Bulge: 0.2399524395791912},
									{X: -104.2556868811259, Y: -155.6163276182112},
									{X: -76.17227789834044, Y: -155.6163276182112},
									{X: -76.17227789834044, Y: -290.5826920162851},
									{X: -94.91760628696476, Y: -309.328020404904},
									{X: -128.6591973864762, Y: -309.328020404904},
									{X: -109.9138689978591, Y: -290.5826920162851},
									{X: -129.9456733421635, Y: -290.5826920162851, Bulge: 0.1989123673795757},
									{X: -140.5496323974874, Y: -294.9749956718473},
									{X: -154.9026571305532, Y: -309.328020404904},
									{X: -282.3708901731735, Y: -309.328020404904},
									{X: -263.625561784551, Y: -290.5826920162851},
									{X: -263.6255617845509, Y: -253.0920352390419},
									{X: -301.1162185617978, Y: -253.0920352390419, Bulge: -0.414213562373095},
									{X: -263.625561784551, Y: -215.6013784617996},
									{X: -263.625561784551, Y: -20.64996322013826},
									{X: -76.17227789834588, Y: -20.64996322013827},
								},
							},
						},
						{
							Flag: 6,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -168.0243870026028, Y: -65.63875135283047, Bulge: -1.0},
									{X: -123.035598869912, Y: -65.63875135283047, Bulge: -1.0},
								},
							},
						},
					},
				},
				{
					Entity:      &entity.EntityData{Handle: 286, Owner: 254, LayerName: "31_color_white"},
					SolidFill:   1,
					Associative: 1,
					PatternName: "SOLID",
					Pattern:     1,
					Style:       1,
					PixelSize:   0.551072650145,
					SeedPoint:   [3]float64{-197.445021060176, -73.56187378603633, 0},
					BoundaryPaths: []*entity.BoundaryPath{
						{
							Flag: 7,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: 150.5379367708664, Y: -60.49372612885668},
									{X: 63.44523592602116, Y: -60.49372612885668},
									{X: 138.4265494805094, Y: -131.7259740056195},
									{X: 18.45644779332895, Y: -131.7259740056195},
									{X: 18.45644779332895, Y: -191.7110248492088},
									{X: 138.4265494805021, Y: -191.7110248492088},
									{X: 63.44523592602116, Y: -262.943272725968},
									{X: 150.5379367708664, Y: -262.9432727259716},
									{X: 257.0903297166986, Y: -161.7184994274128},
								},
							},
						},
						{
							Flag: 7,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -145.7029109369615, Y: -95.63127677462376, Bulge: -0.2399524395792819},
									{X: -130.542887456965, Y: -103.3511335414381},
									{X: -107.287691577128, Y: -135.3270278762274, Bulge: 0.2399524395792833},
									{X: -104.2556868811299, Y: -136.8709992295896},
									{X: -76.17227789834044, Y: -136.8709992295896},
									{X: -76.1722778983459, Y: -20.64996322013826},
									{X: -263.6255617845509, Y: -20.64996322013826},
									{X: -263.625561784551, Y: -215.6013784617996, Bulge: 0.414213562373095},
									{X: -301.1162185617978, Y: -253.0920352390419},
									{X: -263.625561784551, Y: -253.0920352390419},
									{X: -263.625561784551, Y: -290.5826920162851},
									{X: -282.3708901731735, Y: -309.328020404904},
									{X: -154.9026571305532, Y: -309.328020404904},
									{X: -140.5496323974874, Y: -294.9749956718473, Bulge: -0.198912367379576},
									{X: -129.9456733421634, Y: -290.5826920162851},
									{X: -109.9138689978591, Y: -290.5826920162851},
									{X: -128.6591973864762, Y: -309.328020404904},
									{X: -94.91760628696474, Y: -309.328020404904},
									{X: -76.17227789834044, Y: -290.5826920162851},
									{X: -76.17227789834047, Y: -155.6163276182112},
									{X: -104.2556868811258, Y: -155.6163276182112, Bulge: -0.2399524395791912},
									{X: -122.4477150571095, Y: -146.3524994980441},
									{X: -136.3340862192784, Y: -127.2587391500678, Bulge: 0.5998762433504935},
									{X: -142.8859513795958, Y: -128.173067900189},
									{X: -156.4242319135346, Y: -165.0913583761178, Bulge: 0.1873489039227227},
									{X: -156.3728768002399, Y: -167.8051362784292},
									{X: -109.0775649966, Y: -283.0845606608361},
									{X: -127.3869924453994, Y: -283.0845606608361, Bulge: -0.3042679099072649},
									{X: -141.261013849475, Y: -273.780347262223},
									{X: -181.2234173657872, Y: -176.3744374878268},
									{X: -184.5620170200571, Y: -219.776232993233, Bulge: -0.3919022311992351},
									{X: -214.4661990233739, Y: -247.4684367224561},
									{X: -263.625561784551, Y: -247.4684367224561},
									{X: -294.9924956785818, Y: -247.4684367224543, Bulge: -0.3631588356135517},
									{X: -263.6255617845537, Y: -221.2249769783855},
									{X: -214.466199023372, Y: -221.2249769783855, Bulge: 0.3919022311992419},
									{X: -210.7281762729564, Y: -217.7634515122316},
									{X: -206.8321992824213, Y: -167.1157506353802, Bulge: -0.0687858006254979},
									{X: -203.5597649909132, Y: -152.0673701953252},
									{X: -189.738210851636, Y: -114.3766051632454},
									{X: -212.2326049179772, Y: -114.3766051632454, Bulge: 0.2840790438404417},
									{X: -215.5858712009779, Y: -116.449037699469},
									{X: -231.8629676586467, Y: -149.0032306148157, Bulge: -1.0},
									{X: -248.6292990736522, Y: -140.6200649073145},
									{X: -232.352202615976, Y: -108.0658719919629, Bulge: -0.2840790438404093},
									{X: -212.2326049179741, Y: -95.63127677462376},
								},
							},
						},
						{
							Flag: 22,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -168.0243870026028, Y: -65.63875135283047, Bulge: -1},
									{X: -123.035598869912, Y: -65.63875135283047, Bulge: -1},
								},
							},
						},
					},
				},
				{
					Entity:      &entity.EntityData{Handle: 291, Owner: 254, LayerName: "0"},
					SolidFill:   1,
					Associative: 1,
					Style:       1,
					Pattern:     1,
					PatternName: "SOLID",
					PixelSize:   0.551072650145,
					SeedPoint:   [3]float64{-145.6046038420661, -59.78505755621447, 0},
					BoundaryPaths: []*entity.BoundaryPath{
						{
							Flag: 7,
							Polyline: &entity.Polyline{
								Flag: 1,
								Coordinates: []entity.PLine{
									{X: -168.0243870026028, Y: -65.63875135283047, Bulge: 1.0},
									{X: -123.035598869912, Y: -65.63875135283047, Bulge: 1.0},
								},
							},
						},
					},
				},
			},
			Polylines: []*entity.Polyline{
				{
					Entity: &entity.EntityData{Handle: 287, Owner: 254, LayerName: "32_color_black"},
					Flag:   1,
					Coordinates: []entity.PLine{
						{X: 150.537936770866, Y: -60.493726128857},
						{X: 63.445235926021, Y: -60.493726128857},
						{X: 138.426549480509, Y: -131.725974005619},
						{X: 18.456447793329, Y: -131.725974005619},
						{X: 18.456447793329, Y: -191.711024849209},
						{X: 138.426549480502, Y: -191.711024849209},
						{X: 63.445235926021, Y: -262.943272725968},
						{X: 150.537936770866, Y: -262.943272725972},
						{X: 257.090329716699, Y: -161.718499427413},
					},
				},
				{
					Entity: &entity.EntityData{Handle: 288, Owner: 254, LayerName: "32_color_black"},
					Flag:   1,
					Coordinates: []entity.PLine{
						{X: 328.059098645604, Y: -329.977983625045},
						{X: 328.059098645592, Y: 0.000000000004},
						{X: -326.418307086611, Y: 0.000000000001},
						{X: -326.418307086611, Y: -329.977983625045},
					},
				},
				{
					Entity: &entity.EntityData{Handle: 289, Owner: 254, LayerName: "32_color_black"},
					Coordinates: []entity.PLine{
						{X: -76.17227789834, Y: -136.87099922959},
						{X: -76.172277898346, Y: -20.649963220138},
						{X: -263.625561784551, Y: -20.649963220138},
						{X: -263.625561784551, Y: -215.6013784618, Bulge: 0.414213562373},
						{X: -301.116218561798, Y: -253.092035239042},
						{X: -263.625561784551, Y: -253.092035239042},
						{X: -263.625561784551, Y: -290.582692016285},
						{X: -282.370890173174, Y: -309.328020404904},
						{X: -154.902657130553, Y: -309.328020404904},
						{X: -140.549632397493, Y: -294.974995671853, Bulge: -0.19891236738},
						{X: -129.945673342163, Y: -290.582692016285},
						{X: -109.913868997859, Y: -290.582692016285},
						{X: -128.659197386476, Y: -309.328020404904},
						{X: -94.917606286965, Y: -309.328020404904},
						{X: -76.17227789834, Y: -290.582692016285},
						{X: -76.17227789834, Y: -155.616327618211},
					},
				},
				{
					Entity: &entity.EntityData{Handle: 290, Owner: 254, LayerName: "32_color_black"},
					Coordinates: []entity.PLine{
						{X: -76.17227789834, Y: -136.87099922959},
						{X: -104.255686881128, Y: -136.87099922959, Bulge: -0.239952439579},
						{X: -107.287691577127, Y: -135.327027876229},
						{X: -130.542887456968, Y: -103.351133541434, Bulge: 0.239952439579},
						{X: -145.702910936963, Y: -95.631276774624},
						{X: -212.232604917977, Y: -95.631276774624, Bulge: 0.28407904384},
						{X: -232.352202615977, Y: -108.065871991965},
						{X: -248.629299073653, Y: -140.620064907316, Bulge: 1.0},
						{X: -231.862967658648, Y: -149.003230614818},
						{X: -215.585871200978, Y: -116.449037699469, Bulge: -0.28407904384},
						{X: -212.232604917977, Y: -114.376605163245},
						{X: -189.738210851636, Y: -114.376605163245},
						{X: -203.559764990898, Y: -152.067370195284, Bulge: 0.068785800626},
						{X: -206.832199282419, Y: -167.115750635345},
						{X: -210.728176272956, Y: -217.763451512232, Bulge: -0.391902231199},
						{X: -214.466199023373, Y: -221.224976978386},
						{X: -263.625561784551, Y: -221.224976978386, Bulge: 0.363158835614},
						{X: -294.992495678582, Y: -247.468436722454},
						{X: -263.625561784551, Y: -247.468436722456},
						{X: -214.466199023373, Y: -247.468436722456, Bulge: 0.391902231199},
						{X: -184.562017020057, Y: -219.776232993225},
						{X: -181.223417365787, Y: -176.374437487827},
						{X: -141.261013849475, Y: -273.780347262223, Bulge: 0.304267909907},
						{X: -127.386992445399, Y: -283.084560660836},
						{X: -109.0775649966, Y: -283.084560660836},
						{X: -156.372876800237, Y: -167.805136278437, Bulge: -0.187348903923},
						{X: -156.424231913537, Y: -165.091358376126},
						{X: -142.885951379596, Y: -128.17306790019, Bulge: -0.599876243351},
						{X: -136.334086219278, Y: -127.258739150068},
						{X: -122.447715057115, Y: -146.352499498037, Bulge: 0.239952439579},
						{X: -104.255686881128, Y: -155.616327618211},
						{X: -76.17227789834, Y: -155.616327618211},
					},
				},
				{
					Entity: &entity.EntityData{Handle: 293, Owner: 254, LayerName: "40_UMRISS"},
					Flag:   1,
					Coordinates: []entity.PLine{
						{X: -326.418307086611, Y: 0.000000000001},
						{X: 328.059098645604, Y: 0.000000000001},
						{X: 328.059098645604, Y: -329.977983625045},
						{X: -326.418307086611, Y: -329.977983625045},
					},
				},
			},
			Circles: []*entity.Circle{
				{
					Entity:      &entity.EntityData{Handle: 292, Owner: 254, LayerName: "32_color_black"},
					Coordinates: [3]float64{-145.529992936257, -65.63875135283, 0},
					Radius:      22.49439406634541,
				},
			},
		}},
	}

	for name, block := range plan.Blocks {
		actual := blocks.Block{EntitiesData: block.EntitiesData}
		assert.Equal(t, expected[name], actual)
		delete(expected, name)
	}

	if len(expected) > 0 {
		t.Fatal("Block did not parse all entities missing: ", expected)
	}
}
