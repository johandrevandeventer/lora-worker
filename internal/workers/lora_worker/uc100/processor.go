package uc100

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/johandrevandeventer/kafkaclient/payload"
	"github.com/johandrevandeventer/lora-worker/internal/workers"
	"github.com/johandrevandeventer/lora-worker/internal/workers/lora_worker/uc100/atess/pcs250"
	"github.com/johandrevandeventer/lora-worker/internal/workers/lora_worker/uc100/inverter/deye8"
	"github.com/johandrevandeventer/lora-worker/internal/workers/types"
	"go.uber.org/zap"
)

const (
	DeviceTypePCS250 = "pcs250"
	DeviceTypeDeye8  = "deye8"
)

func Processor(msg payload.Payload, logger *zap.Logger) (MessageInfo *types.MessageInfo, err error) {
	var uc100Data UC100
	if err := json.Unmarshal(msg.Message, &uc100Data); err != nil {
		return MessageInfo, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var controllerID string
	var deviceID string

	controllerID = uc100Data.EndDeviceIDs.DevEUI

	logger.Debug("Processing controller", zap.String("controllerID", controllerID))

	ignoredControllers, err := workers.GetIgnoredControllers()
	if err != nil {
		return MessageInfo, fmt.Errorf("error getting ignored controllers: %w", err)
	}

	if slices.Contains(ignoredControllers, controllerID) {
		return MessageInfo, fmt.Errorf("controller is ignored: %s", controllerID)
	}

	devicesResult, err := workers.GetDevicesByControllerIdentifier(controllerID)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return MessageInfo, fmt.Errorf("controller not found: %s", controllerID)
		}
		return MessageInfo, fmt.Errorf("error getting devices by controller ID: %w", err)
	}

	if len(devicesResult) == 0 {
		logger.Warn("No devices found for controller", zap.String("controller_id", controllerID))
	}

	var devices []types.Device

	for _, device := range devicesResult {
		deviceID = device.DeviceIdentifier

		logger.Debug("Processing device", zap.String("deviceID", deviceID))

		ignoredDevices, err := workers.GetIgnoredDevices()
		if err != nil {
			return MessageInfo, fmt.Errorf("error getting ignored devices: %w", err)
		}

		if slices.Contains(ignoredDevices, deviceID) {
			return MessageInfo, fmt.Errorf("device is ignored: %s", deviceID)
		}

		device, err := workers.GetDevicesByDeviceIdentifier(deviceID)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return MessageInfo, fmt.Errorf("device not found: %s", deviceID)
			}

			return MessageInfo, fmt.Errorf("error getting device by device ID - %s: %w", deviceID, err)
		}

		deviceType := device.DeviceType
		deviceTypeLower := strings.ToLower(deviceType)
		timestamp := msg.MessageTimestamp

		data := uc100Data.UplinkMessage.DecodedPayload

		var rawData map[string]any
		var processedData map[string]any

		logger.Debug(fmt.Sprintf("%s :: %s", device.Controller, device.DeviceType))

		switch deviceTypeLower {
		// Process Genset devices
		case DeviceTypePCS250:
			rawData, processedData, err = pcs250.Decoder(data)
			if err != nil {
				return MessageInfo, fmt.Errorf("error decoding pcs250 data: %w", err)
			}
		case DeviceTypeDeye8:
			rawData, processedData, err = deye8.Decoder(data)
			if err != nil {
				return MessageInfo, fmt.Errorf("error decoding deye8 data: %w", err)
			}
		}

		rawData["SerialNo1"] = device.ControllerIdentifier
		processedData["SerialNo1"] = device.ControllerIdentifier

		deviceStruct := &types.Device{
			CustomerID:           device.Site.Customer.ID,
			CustomerName:         device.Site.Customer.Name,
			SiteID:               device.Site.ID,
			SiteName:             device.Site.Name,
			Controller:           device.Controller,
			DeviceType:           device.DeviceType,
			ControllerIdentifier: device.ControllerIdentifier,
			DeviceName:           device.DeviceName,
			DeviceIdentifier:     device.DeviceIdentifier,
			RawData:              rawData,
			ProcessedData:        processedData,
			Timestamp:            timestamp,
		}

		devices = append(devices, *deviceStruct)
	}

	return &types.MessageInfo{
		MessageID: msg.ID.String(),
		Devices:   devices,
	}, nil

}
