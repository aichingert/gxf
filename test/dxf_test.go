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
	// TODO: finish adding remaining fields
	expected := map[string]any{
		"STANDARD_SL": blocks.Block{EntitiesData: &entity.EntitiesData{
			Hatches: []*entity.Hatch{
				&entity.Hatch{
					PatternName: "SOLID",
					SolidFill:   1,
					Associative: 0,
				},
			},
		}},
		"STANDARD_RZ_PR": blocks.Block{},
	}

	for name, block := range plan.Blocks {
		actual := blocks.Block{EntitiesData: block.EntitiesData}

		assert.Equal(t, expected[name], actual)
	}
}
