package pcs250

import (
	"fmt"
	"math"
	"reflect"

	coreutils "github.com/johandrevandeventer/lora-worker/utils"
)

type PCS250Data struct {
	BatV       *float64 `json:"BatV,omitempty" mapstructure:"modbus_chn_1"`
	BatI       *float64 `json:"BatI,omitempty" mapstructure:"modbus_chn_2"`
	BatSOC     *float64 `json:"BatSOC,omitempty" mapstructure:"modbus_chn_3"`
	BatP       *float64 `json:"BatP,omitempty" mapstructure:"modbus_chn_4"`
	GridVa     *float64 `json:"GridVa,omitempty" mapstructure:"modbus_chn_5"`
	GridVb     *float64 `json:"GridVb,omitempty" mapstructure:"modbus_chn_6"`
	GridVc     *float64 `json:"GridVc,omitempty" mapstructure:"modbus_chn_7"`
	GridFreq   *float64 `json:"GridFreq,omitempty" mapstructure:"modbus_chn_8"`
	GridIa     *float64 `json:"GridIa,omitempty" mapstructure:"modbus_chn_9"`
	GridIb     *float64 `json:"GridIb,omitempty" mapstructure:"modbus_chn_10"`
	GridIc     *float64 `json:"GridIc,omitempty" mapstructure:"modbus_chn_11"`
	LoadVa     *float64 `json:"LoadVa,omitempty" mapstructure:"modbus_chn_12"`
	LoadVb     *float64 `json:"LoadVb,omitempty" mapstructure:"modbus_chn_13"`
	LoadVc     *float64 `json:"LoadVc,omitempty" mapstructure:"modbus_chn_14"`
	LoadIa     *float64 `json:"LoadIa,omitempty" mapstructure:"modbus_chn_15"`
	LoadIb     *float64 `json:"LoadIb,omitempty" mapstructure:"modbus_chn_16"`
	LoadIc     *float64 `json:"LoadIc,omitempty" mapstructure:"modbus_chn_17"`
	LoadtotP   *float64 `json:"LoadtotP,omitempty" mapstructure:"modbus_chn_18"`
	LoadE      *float64 `json:"LoadE,omitempty" mapstructure:"modbus_chn_19"`
	GridImpE   *float64 `json:"GridImpE,omitempty" mapstructure:"modbus_chn_20"`
	GridExpE   *float64 `json:"GridExpE,omitempty" mapstructure:"modbus_chn_21"`
	BatChargeE *float64 `json:"BatChargeE,omitempty" mapstructure:"modbus_chn_22"`
	BatDisE    *float64 `json:"BatDisE,omitempty" mapstructure:"modbus_chn_23"`
	GridtotP   *float64 `json:"GridtotP,omitempty" mapstructure:"modbus_chn_24"`
}

var ratioMap = map[string]float64{
	"BatV":       0.1,
	"BatI":       0.1,
	"BatSOC":     1.0,
	"BatP":       0.1,
	"GridVa":     0.1,
	"GridVb":     0.1,
	"GridVc":     0.1,
	"GridFreq":   0.01,
	"GridIa":     0.1,
	"GridIb":     0.1,
	"GridIc":     0.1,
	"LoadVa":     0.1,
	"LoadVb":     0.1,
	"LoadVc":     0.1,
	"LoadIa":     0.1,
	"LoadIb":     0.1,
	"LoadIc":     0.1,
	"LoadtotP":   0.1,
	"LoadE":      0.1,
	"GridImpE":   0.1,
	"GridExpE":   0.1,
	"BatChargeE": 0.1,
	"BatDisE":    0.1,
	"GridtotP":   0.1,
}

func Decoder(payload map[string]any) (rawData, processedData map[string]any, err error) {
	var pcs250Data PCS250Data

	err = coreutils.DecodeMapToStruct(payload, &pcs250Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error decoding PCS250 data: %w", err)
	}

	rawData, err = coreutils.StructToMap(pcs250Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error converting PCS250 data to map: %w", err)
	}

	applyRatios(&pcs250Data)

	processedData, err = coreutils.StructToMap(pcs250Data)
	if err != nil {
		return rawData, processedData, fmt.Errorf("error converting PCS250 data to map: %w", err)
	}

	return rawData, processedData, nil
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
