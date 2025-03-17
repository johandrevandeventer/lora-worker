package milesightworker

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/johandrevandeventer/kafkaclient/payload"
	"github.com/johandrevandeventer/lora-worker/internal/workers"
	uc100_decoder "github.com/johandrevandeventer/lora-worker/internal/workers/lora_worker/milesight/uc100"
	"go.uber.org/zap"
)

const (
	ControllerTypeUC100 = "uc100"
)

func MileSightWorker(msg payload.Payload, logger *zap.Logger) (*workers.DataStruct, *workers.DataStruct, error) {
	logger.Debug("Decoding milesight data")
	mileSightData := &MileSightData{}
	rawDataStruct := &workers.DataStruct{}
	processedDataStruct := &workers.DataStruct{}

	// Parse the JSON string into a map
	err := json.Unmarshal([]byte(msg.Message), &mileSightData)
	if err != nil {
		return rawDataStruct, processedDataStruct, fmt.Errorf("error unmarshalling milesight data: %w", err)
	}

	// Get the controller ID and payload
	controllerId := mileSightData.EndDeviceIDs.DevEUI

	logger.Debug("Processing controller", zap.String("controllerId", controllerId))

	ignoredControllers, err := workers.GetIgnoredControllers()
	if err != nil {
		return rawDataStruct, processedDataStruct, fmt.Errorf("error getting ignored controllers: %w", err)
	}

	if slices.Contains(ignoredControllers, controllerId) {
		logger.Warn("Controller is ignored", zap.String("controller_id", controllerId))
		return rawDataStruct, processedDataStruct, nil
	}

	logger.Debug("Fetching devices by controller ID", zap.String("controllerId", controllerId))

	devices, err := workers.GetDevicesByControllerSerialNumber(controllerId)
	if err != nil {
		return rawDataStruct, processedDataStruct, fmt.Errorf("error getting devices by controller ID: %w", err)
	}

	if len(devices) == 0 {
		logger.Warn("No devices found for controller", zap.String("controller_id", controllerId))
	}

	for _, device := range devices {
		if device.ControllerSerialNumber != controllerId {
			continue
		}

		timestamp := msg.MessageTimestamp
		controller := device.Controller
		controllerLower := strings.ToLower(controller)

		ignoredDevices, err := workers.GetIgnoredDevices()
		if err != nil {
			return rawDataStruct, processedDataStruct, fmt.Errorf("error getting ignored devices: %w", err)
		}

		for _, ignoredDevice := range ignoredDevices {
			if device.DeviceSerialNumber == ignoredDevice {
				logger.Warn("Device is ignored", zap.String("device_serial_number", device.DeviceSerialNumber))
				continue
			}
		}

		logger.Info(fmt.Sprintf("%s :: %s :: %s :: %s :: %s :: %s", device.Gateway, device.Site.Customer.Name, device.Site.Name, device.DeviceType, device.ControllerSerialNumber, device.DeviceSerialNumber))

		switch controllerLower {
		case ControllerTypeUC100:
			rawData, processedData, err := uc100_decoder.UC100Decoder(mileSightData.UplinkMessage.DecodedPayload, device.DeviceType, logger)
			if err != nil {
				return rawDataStruct, processedDataStruct, fmt.Errorf("error decoding UC100 data: %w", err)
			}

			if len(rawData) == 0 {
				logger.Warn("No raw data found", zap.String("controller", controllerId))
				continue
			}

			if len(processedData) == 0 {
				logger.Warn("No processed data found", zap.String("controller", controllerId))
				continue
			}

			rawData["SerialNo1"] = device.ControllerSerialNumber
			processedData["SerialNo1"] = device.ControllerSerialNumber

			rawDataStruct = &workers.DataStruct{
				State:                  "Pre",
				CustomerID:             device.Site.Customer.ID,
				CustomerName:           device.Site.Customer.Name,
				SiteID:                 device.Site.ID,
				SiteName:               device.Site.Name,
				Gateway:                device.Gateway,
				Controller:             device.Controller,
				DeviceType:             device.DeviceType,
				ControllerSerialNumber: device.ControllerSerialNumber,
				DeviceName:             device.DeviceName,
				DeviceSerialNumber:     device.DeviceSerialNumber,
				Data:                   rawData,
				Timestamp:              timestamp,
			}

			processedDataStruct = &workers.DataStruct{
				State:                  "Post",
				CustomerID:             device.Site.Customer.ID,
				CustomerName:           device.Site.Customer.Name,
				SiteID:                 device.Site.ID,
				SiteName:               device.Site.Name,
				Gateway:                device.Gateway,
				Controller:             device.Controller,
				DeviceType:             device.DeviceType,
				ControllerSerialNumber: device.ControllerSerialNumber,
				DeviceName:             device.DeviceName,
				DeviceSerialNumber:     device.DeviceSerialNumber,
				Data:                   processedData,
				Timestamp:              timestamp,
			}

			return rawDataStruct, processedDataStruct, nil
		default:
			logger.Warn("Controller not supported", zap.String("controller", controller))
		}
	}

	return rawDataStruct, processedDataStruct, nil
}
