package hpcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Volume encapsulates the volumes available in the OpenStack/HP Cloud
// system *as well as* the volumes you are actually running.
//
// Typically you would create a value of this type and pass it to the
// CreateVolume function. This type will be returned from any query
// endpoints which return containers of info about the volumes you may
// have.
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

// ListVolumes returns a slice of volumes which are currently
// associated with the token_id you provide.
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

// ListSnapshots will return a slice of Volumes for which are in-fact
// snapshots of your systems.
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

// NewVolume takes a volume instance and will create that in the
// cloud. This function will return *before* the instance is
// created. In order to know when the instance has been created you
// will need to check the status using the provided methods.
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

// DetachVolume will remove a volume from whatever server it is
// attached to.
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

// We override MarshalJSON because we want to provide additional
// marshaling logic when creating new compute nodes. This is because
// the zero values of Volumes are not valid parameters for the compute
// API.
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

	// Iterate through via reflect on the remaining fields. These
	// fields require no special treatment and therefore can be simply
	// interpolated into the request.
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
	// We ignore the metadata completely if there are no additional
	// fields.
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
