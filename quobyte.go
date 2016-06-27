// Golang API for the Quobyte Storage System
package quobyte

import (
	"bytes"
	"github.com/gorilla/rpc/v2/json2"
	"net/http"
)

type QuobyteClient struct {
	client   *http.Client
	url      string
	username string
	password string
}

// Create a new Quobyte API client
func NewQuobyteClient(url string, username string, password string) *QuobyteClient {
	result := new(QuobyteClient)
	result.client = new(http.Client)
	result.url = url
	result.username = username
	result.password = password
	return result
}

type CreateVolumeRequest struct {
	Name        string `json:"name"`
	RootUserId  string `json:"root_user_id"`
	RootGroupId string `json:"root_group_id"`
}

type CreateVolumeResponse struct {
	VolumeUuid string `json:"volume_uuid"`
}

// Create a new Quobyte volume. Its root directory will be owned by given user and group
func (client QuobyteClient) CreateVolume(name string, rootUserName string, rootGroupName string) (volumeUuid string, err error) {
	request := &CreateVolumeRequest{
		Name:        name,
		RootUserId:  rootUserName,
		RootGroupId: rootGroupName,
	}
	var response CreateVolumeResponse
	err = client.sendRequest("createVolume", request, &response)
	if err != nil {
		return "", err
	}
	return response.VolumeUuid, nil
}

type DeleteVolumeRequest struct {
	VolumeUuid string `json:"volume_uuid"`
}

type DeleteVolumeResponse struct {
}

// Delete a Quobyte volume. Its root directory will be owned by given user and group and have access 700.
func (client QuobyteClient) DeleteVolume(volumeUuid string) error {
	request := &DeleteVolumeRequest{
		VolumeUuid: volumeUuid,
	}
	var response DeleteVolumeResponse
	return client.sendRequest("deleteVolume", request, &response)
}

func (client QuobyteClient) sendRequest(method string, request interface{}, response interface{}) error {
	message, err := json2.EncodeClientRequest(method, request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", client.url, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(client.username, client.password)
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json2.DecodeClientResponse(resp.Body, &response)
	if err != nil {
		return err
	}
	return nil
}
