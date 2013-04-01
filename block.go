package hpcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Volume struct {
	Status           string            `json:"status"`
	CreatedAt        string            `json:"createdAt"`
	Size             int64             `json:"size"`
	DisplayName      string            `json:"display_name"`
	DisplayDesc      string            `json:"display_description"`
	SnapshotID       int64             `json:"snapshot_id"`
	ImageRef         int64             `json:"imageRef"`
	Metadata         map[string]string `json:"metadata"`
	AvailabilityZone string            `json:"availability_zone"`
	VolumeType       string            `json:"volume_type"`
	Attachments      []Attachment      `json:"attachments"`
}

type Attachment struct {
	ID       int64  `json:"id"`
	Device   string `json:"device"`
	ServerID int64  `json:"serverId"`
	VolumeID int64  `json:"volumeId"`
}

func (a Access) ListVolumes() ([]Volume, error) {
	resp, err := a.baseRequest(
		fmt.Sprintf("%s%s/os-volumes", COMPUTE_URL, a.TenantID),
		"GET", nil,
	)
	if err != nil {
		return nil, err
	}
	type Volumes struct {
		V []Volume `json:"volumes"`
	}
	vs := &Volumes{}
	json.Unmarshal(resp, vs)
	return vs.V, nil
}

func (a Access) ListVolumesForServer(server_id string) ([]Attachment, error) {
	resp, err := a.baseRequest(
		fmt.Sprintf("%s%s/servers/%s/os-volume_attachments", COMPUTE_URL, a.TenantID, server_id),
		"GET", nil,
	)
	if err != nil {
		return nil, err
	}
	type Attachments struct {
		V []Attachment `json:"volumeAttachments"`
	}
	vs := &Attachments{}
	json.Unmarshal(resp, vs)
	return vs.V, nil
}

func (a Access) ListSnapshots() ([]Volume, error) {
	resp, err := a.baseRequest(
		fmt.Sprintf("%s%s/os-snapshots", COMPUTE_URL, a.TenantID),
		"GET", nil,
	)
	if err != nil {
		return nil, err
	}
	type Volumes struct {
		V []Volume `json:"snapshots"`
	}
	vs := &Volumes{}
	json.Unmarshal(resp, vs)
	return vs.V, nil
}

func (a Access) NewVolume(v *Volume) error {
	b, err := v.MarshalJSON()
	if err != nil {
		return err
	}
	resp, err := a.baseRequest(
		fmt.Sprintf("%s%s/os-volumes", COMPUTE_URL, a.TenantID),
		"POST", strings.NewReader(string(b)),
	)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, v)
}

func (a Access) DetachVolume(at Attachment) error {
	_, err := a.baseRequest(
		fmt.Sprintf(
			"%s%s/servers/%d/os-volume_attachments/%d",
			COMPUTE_URL, a.TenantID, at.ServerID, at.VolumeID,
		), "DELETE", nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (v Volume) MarshalJSON() ([]byte, error) {
	if v.Size <= 0 {
		return nil, errors.New("Size cannot be <= 0")
	}
	b := bytes.NewBufferString(
		fmt.Sprintf(`{"volume":{"size":%d`, v.Size),
	)

	val := reflect.ValueOf(&v).Elem()
	var i64 int64
	var str string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		T := reflect.TypeOf(v).Field(i)
		if field.Type() == reflect.ValueOf(i64).Type() && field.Int() != i64 {
			b.WriteString(
				fmt.Sprintf(`,"%s": %d`, T.Tag.Get("json"), field.Int()),
			)
		}
		if field.Type() == reflect.ValueOf(str).Type() && field.String() != str {
			b.WriteString(
				fmt.Sprintf(`,"%s": "%s"`, T.Tag.Get("json"), field.String()),
			)
		}
	}
	if len(v.Metadata) > 0 {
		b.WriteString(`,"metadata":{`)
		metadata := make([]string, 0, len(v.Metadata))
		for k, v := range v.Metadata {
			metadata = append(metadata, fmt.Sprintf(`"%s":"%s"`, k, v))
		}
		b.WriteString(strings.Join(metadata, ",") + "}")
	}
	b.WriteString("}}")
	return b.Bytes(), nil
}
