package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// NotificationDevice represents a device for message notification
type NotificationDevice struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	Active      string `json:"Active"`
}

// listNotificationDevicesResponse handles generic notification device list responses
type listNotificationDevicesResponse struct {
	Total       string                  `json:"@total"`
	PhoneDevices OneOrMany[NotificationDevice]   `json:"PhoneDevice"`
	PagerDevices OneOrMany[NotificationDevice]   `json:"PagerDevice"`
	SmtpDevices  OneOrMany[NotificationDevice]   `json:"SmtpDevice"`
	HtmlDevices  OneOrMany[NotificationDevice]   `json:"HtmlDevice"`
}

// getDeviceKey returns the JSON key for the given device type
func getDeviceKey(deviceType string) string {
	switch deviceType {
	case "phonedevices":
		return "PhoneDevice"
	case "pagerdevices":
		return "PagerDevice"
	case "smtpdevices":
		return "SmtpDevice"
	case "htmldevices":
		return "HtmlDevice"
	default:
		return deviceType
	}
}

// ListNotificationDevices returns notification devices for a user by type
func ListNotificationDevices(host string, port int, user, pass, userObjectId, deviceType string) ([]NotificationDevice, error) {
	path := fmt.Sprintf("/users/%s/notificationdevices/%s", url.PathEscape(userObjectId), deviceType)
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification devices: %w", err)
	}

	key := getDeviceKey(deviceType)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse notification devices response: %w", err)
	}

	var devices []NotificationDevice
	if items, ok := data[key]; ok {
		itemsJSON, _ := json.Marshal(items)
		if err := json.Unmarshal(itemsJSON, &devices); err != nil {
			return nil, fmt.Errorf("failed to parse devices: %w", err)
		}
	}

	return devices, nil
}

// GetNotificationDevice retrieves a specific notification device
func GetNotificationDevice(host string, port int, user, pass, userObjectId, deviceType, deviceObjectId string) (*NotificationDevice, error) {
	path := fmt.Sprintf("/users/%s/notificationdevices/%s/%s", url.PathEscape(userObjectId), deviceType, url.PathEscape(deviceObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification device: %w", err)
	}

	var d NotificationDevice
	if err := json.Unmarshal(body, &d); err != nil {
		return nil, fmt.Errorf("failed to parse notification device: %w", err)
	}

	return &d, nil
}

// UpdateNotificationDevice updates a notification device
func UpdateNotificationDevice(host string, port int, user, pass, userObjectId, deviceType, deviceObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/notificationdevices/%s/%s", url.PathEscape(userObjectId), deviceType, url.PathEscape(deviceObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update notification device: %w", err)
	}
	return nil
}

// CreateNotificationDevice creates a notification device (only for phone and pager devices)
func CreateNotificationDevice(host string, port int, user, pass, userObjectId, deviceType string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/notificationdevices/%s", url.PathEscape(userObjectId), deviceType)
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create notification device: %w", err)
	}
	return nil
}

// DeleteNotificationDevice deletes a notification device (only for phone and pager devices)
func DeleteNotificationDevice(host string, port int, user, pass, userObjectId, deviceType, deviceObjectId string) error {
	path := fmt.Sprintf("/users/%s/notificationdevices/%s/%s", url.PathEscape(userObjectId), deviceType, url.PathEscape(deviceObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete notification device: %w", err)
	}
	return nil
}
