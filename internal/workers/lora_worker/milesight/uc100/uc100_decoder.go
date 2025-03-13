package uc100_decoder

import (
	"fmt"
	"strings"

	pcs250_decoder "github.com/johandrevandeventer/lora-worker/internal/workers/lora_worker/milesight/uc100/atess/pcs250"
	deye8_decoder "github.com/johandrevandeventer/lora-worker/internal/workers/lora_worker/milesight/uc100/inverters/deye8"
	"go.uber.org/zap"
)

const (
	deviceTypeDeye8  = "deye8"
	deviceTypePCS250 = "pcs250"
)

func UC100Decoder(payload map[string]any, deviceType string, logger *zap.Logger) (rawData, processedData map[string]any, err error) {
	deviceTypeLower := strings.ToLower(deviceType)

	var decodeRawFunc func(map[string]any) (map[string]any, error)
	var decodeProcessedFunc func(map[string]any) (map[string]any, error)

	logger.Debug("Decoding device data", zap.String("deviceType", deviceType))

	switch deviceTypeLower {
	case deviceTypeDeye8:
		decodeRawFunc = deye8_decoder.Deye8DecoderRaw
		decodeProcessedFunc = deye8_decoder.Deye8DecoderProcessed
	case deviceTypePCS250:
		decodeRawFunc = pcs250_decoder.PCS250DecoderRaw
		decodeProcessedFunc = pcs250_decoder.PCS250DecoderProcessed
	default:
		// return nil, nil, fmt.Errorf("unknown device type: %s", deviceType)
		logger.Warn("Unknown device type", zap.String("deviceType", deviceType))
		return nil, nil, nil
	}

	// Execute the decoding functions
	if rawData, err = decodeRawFunc(payload); err != nil {
		return nil, nil, fmt.Errorf("error decoding raw data for %s: %w", deviceType, err)
	}

	if processedData, err = decodeProcessedFunc(payload); err != nil {
		return nil, nil, fmt.Errorf("error processing data for %s: %w", deviceType, err)
	}

	return rawData, processedData, nil
}
