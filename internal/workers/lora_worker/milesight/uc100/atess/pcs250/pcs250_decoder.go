package pcs250_decoder

import (
	"fmt"
	"math"
	"reflect"

	coreutils "github.com/johandrevandeventer/lora-worker/utils"
)

var ratioMap = map[string]float64{
	"BatV":       0.1,
	"BatI":       1.0,
	"BatSOC":     1.0,
	"BatP":       0.001,
	"GridVa":     0.1,
	"GridVb":     0.1,
	"GridVc":     0.1,
	"GridFreq":   1.0,
	"GridIa":     1.0,
	"GridIb":     1.0,
	"GridIc":     1.0,
	"LoadVa":     0.1,
	"LoadVb":     0.1,
	"LoadVc":     0.1,
	"LoadIa":     1.0,
	"LoadIb":     1.0,
	"LoadIc":     1.0,
	"LoadtotP":   0.001,
	"LoadE":      0.001,
	"GridImpE":   0.001,
	"GridExpE":   0.001,
	"BatChargeE": 0.001,
	"BatDisE":    0.001,
	"GridtotP":   0.001,
}

// PCS250DecoderRaw decodes the raw PCS250 data
func PCS250DecoderRaw(payload map[string]any) (map[string]any, error) {
	var pcs250250Data PCS250Data

	// Decode map into struct
	err := coreutils.DecodeMapToStruct(payload, &pcs250250Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding PCS250 data: %w", err)
	}

	pcs250250DataMap, err := coreutils.StructToMap(pcs250250Data)
	if err != nil {
		return nil, fmt.Errorf("error converting PCS250 data to map: %w", err)
	}

	return pcs250250DataMap, nil
}

// PCS250DecoderProcessed decodes and processes the PCS250 data
func PCS250DecoderProcessed(payload map[string]any) (map[string]any, error) {
	var pcs250Data PCS250Data

	// Decode map into struct
	err := coreutils.DecodeMapToStruct(payload, &pcs250Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding PCS250 data: %w", err)
	}

	applyRatios(&pcs250Data)

	pcs250DataMap, err := coreutils.StructToMap(pcs250Data)
	if err != nil {
		return nil, fmt.Errorf("error converting PCS250 data to map: %w", err)
	}

	return pcs250DataMap, nil
}

func applyRatios(dm *PCS250Data) {
	// Get the reflection value of the Deye8Data struct
	dmValue := reflect.ValueOf(dm).Elem()

	// Iterate over the fields in the struct
	for i := 0; i < dmValue.NumField(); i++ {
		// Get the field name and value
		field := dmValue.Type().Field(i)
		fieldValue := dmValue.Field(i)

		// Look for a ratio in ratioMap for the current field's name
		if ratio, ok := ratioMap[field.Name]; ok {
			// Ensure the field is a pointer to a float64 and not nil
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Float64 {
				if !fieldValue.IsNil() {
					// Apply the ratio
					original := fieldValue.Elem().Float()
					newValue := math.Round(original*ratio*100) / 100
					fieldValue.Elem().SetFloat(newValue)
				}
			}
		}
	}
}
