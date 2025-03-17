package deye8_decoder

import (
	"fmt"
	"math"
	"reflect"

	coreutils "github.com/johandrevandeventer/lora-worker/utils"
)

var ratioMap = map[string]float64{
	"FreqIn":      0.01,
	"Vin":         0.1,
	"Iin":         1,
	"Pin":         0.001,
	"Vout":        0.1,
	"Iout":        1,
	"Pout":        0.001,
	"BatV":        0.01,
	"BatSOC":      1,
	"BatP":        0.001,
	"BatI":        0.01,
	"BatChargeE":  0.001,
	"BatDisE":     0.001,
	"LoadE":       0.01,
	"GridImportE": 0.01,
	"GridExpE":    0.01,
}

// Deye8DecoderRaw decodes the raw Deye8 data
func Deye8DecoderRaw(payload map[string]any) (map[string]any, error) {
	var deye8Data Deye8Data

	// Decode map into struct
	err := coreutils.DecodeMapToStruct(payload, &deye8Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding Deye8 data: %w", err)
	}

	if deye8Data.Vin != nil && deye8Data.Pin != nil {
		if *deye8Data.Vin != 0 && *deye8Data.Pin != 0 {
			value := math.Round((*deye8Data.Pin*1000) / *deye8Data.Vin) / 100
			deye8Data.Iin = &value
		}
	}

	deye8DataMap, err := coreutils.StructToMap(deye8Data)
	if err != nil {
		return nil, fmt.Errorf("error converting Deye8 data to map: %w", err)
	}

	return deye8DataMap, nil
}

// Deye8DecoderProcessed decodes and processes the Deye8 data
func Deye8DecoderProcessed(payload map[string]any) (map[string]any, error) {
	var deye8Data Deye8Data

	// Decode map into struct
	err := coreutils.DecodeMapToStruct(payload, &deye8Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding Deye8 data: %w", err)
	}

	applyRatios(&deye8Data)

	if deye8Data.Vin != nil && deye8Data.Pin != nil {
		if *deye8Data.Vin != 0 && *deye8Data.Pin != 0 {
			value := math.Round((*deye8Data.Pin*1000) / *deye8Data.Vin) / 100
			deye8Data.Iin = &value
		}
	}

	deye8DataMap, err := coreutils.StructToMap(deye8Data)
	if err != nil {
		return nil, fmt.Errorf("error converting Deye8 data to map: %w", err)
	}

	return deye8DataMap, nil
}

func applyRatios(dm *Deye8Data) {
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
