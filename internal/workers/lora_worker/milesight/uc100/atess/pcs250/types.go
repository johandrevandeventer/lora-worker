package pcs250_decoder

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
