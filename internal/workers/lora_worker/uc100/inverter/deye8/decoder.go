package deye8

import (
	"fmt"
	"math"
	"reflect"

	coreutils "github.com/johandrevandeventer/lora-worker/utils"
)

type Deye8Data struct {
	FreqIn      *float64 `json:"FreqIn,omitempty" mapstructure:"modbus_chn_1"`
	Vin         *float64 `json:"Vin,omitempty" mapstructure:"modbus_chn_2"`
	Iin         *float64 `json:"Iin,omitempty" mapstructure:"modbus_chn_3"`
	Pin         *float64 `json:"Pin,omitempty" mapstructure:"modbus_chn_4"`
	Vout        *float64 `json:"Vout,omitempty" mapstructure:"modbus_chn_5"`
	Iout        *float64 `json:"Iout,omitempty" mapstructure:"modbus_chn_6"`
	Pout        *float64 `json:"Pout,omitempty" mapstructure:"modbus_chn_7"`
	BatV        *float64 `json:"BatV,omitempty" mapstructure:"modbus_chn_8"`
	BatSOC      *float64 `json:"BatSOC,omitempty" mapstructure:"modbus_chn_9"`
	BatP        *float64 `json:"BatP,omitempty" mapstructure:"modbus_chn_10"`
	BatI        *float64 `json:"BatI,omitempty" mapstructure:"modbus_chn_11"`
	BatChargeE  *float64 `json:"BatChargeE,omitempty" mapstructure:"modbus_chn_12"`
	BatDisE     *float64 `json:"BatDisE,omitempty" mapstructure:"modbus_chn_13"`
	LoadE       *float64 `json:"LoadE,omitempty" mapstructure:"modbus_chn_14"`
	GridImportE *float64 `json:"GridImportE,omitempty" mapstructure:"modbus_chn_15"`
	GridExpE    *float64 `json:"GridExpE,omitempty" mapstructure:"modbus_chn_16"`
}

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

func Decoder(payload map[string]any) (rawData, processedData map[string]any, err error) {
	var deye8Data Deye8Data

	err = coreutils.DecodeMapToStruct(payload, &deye8Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error decoding Deye8 data: %w", err)
	}

	rawData, err = coreutils.StructToMap(deye8Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error converting Deye8 data to map: %w", err)
	}

	applyRatios(&deye8Data)

	processedData, err = coreutils.StructToMap(deye8Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error converting Deye8 data to map: %w", err)
	}

	return rawData, processedData, nil
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
