package deye8_decoder

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
