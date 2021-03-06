package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testSecureBindingParameters = SecureBindingParameters{
		"foo": fooValue,
	}
	testSecureBindingDetails = SecureBindingDetails{
		"foo": fooValue,
	}
	testBinding     Binding
	testBindingJSON []byte
)

func init() {
	bindingID := "test-binding-id"
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	bindingParameters := BindingParameters{
		"foo": "bar",
	}
	bindingParametersJSONStr := []byte(`{"foo":"bar"}`)
	secureBindingParameters := SecureBindingParameters{
		"foo": "bar",
	}
	encryptedSecureBindingParameters := []byte(`{"foo":"bar"}`)
	statusReason := "in-progress"
	bindingDetails := BindingDetails{
		"foo": "bat",
	}
	bindingDetailsJSONStr := `{"foo":"bat"}`
	secureBindingDetails := SecureBindingDetails{
		"foo": "bat",
	}
	encryptedSecureBindingDetails := []byte(`{"foo":"bat"}`)
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testBinding = Binding{
		BindingID:                        bindingID,
		InstanceID:                       instanceID,
		ServiceID:                        serviceID,
		BindingParameters:                bindingParameters,
		EncryptedSecureBindingParameters: encryptedSecureBindingParameters,
		SecureBindingParameters:          secureBindingParameters,
		Status:                           BindingStateBound,
		StatusReason:                     statusReason,
		Details:                          bindingDetails,
		EncryptedSecureDetails:           encryptedSecureBindingDetails,
		SecureDetails:                    secureBindingDetails,
		Created:                          created,
	}

	b64EncryptedSecureBindingParameters := base64.StdEncoding.EncodeToString(
		encryptedSecureBindingParameters,
	)
	b64EncryptedSecureBindingDetails := base64.StdEncoding.EncodeToString(
		encryptedSecureBindingDetails,
	)

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":%s,
			"secureBindingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"details":%s,
			"secureDetails":"%s",
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		bindingParametersJSONStr,
		b64EncryptedSecureBindingParameters,
		BindingStateBound,
		statusReason,
		bindingDetailsJSONStr,
		b64EncryptedSecureBindingDetails,
		created.Format(time.RFC3339),
	)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, " ", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\n", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\t", "", -1)
	testBindingJSON = []byte(testBindingJSONStr)
}

func TestNewBindingFromJSON(t *testing.T) {
	binding, err := NewBindingFromJSON(testBindingJSON, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	json, err := testBinding.ToJSON(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, json)
}

func TestEncryptSecureBindingParameters(t *testing.T) {
	binding := Binding{
		SecureBindingParameters: testSecureBindingParameters,
	}
	var err error
	binding, err = binding.encryptSecureBindingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedSecureBindingParameters,
	)
}

func TestEncryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		SecureDetails: testSecureBindingDetails,
	}
	var err error
	binding, err = binding.encryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedSecureDetails,
	)
}

func TestDecryptSecureBindingParameters(t *testing.T) {
	binding := Binding{
		EncryptedSecureBindingParameters: testArbitraryObjectJSON,
		SecureBindingParameters:          SecureBindingParameters{},
	}
	var err error
	binding, err = binding.decryptSecureBindingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testSecureBindingParameters, binding.SecureBindingParameters)
}

func TestDecryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		EncryptedSecureDetails: testArbitraryObjectJSON,
		SecureDetails:          SecureBindingDetails{},
	}
	var err error
	binding, err = binding.decryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testSecureBindingDetails, binding.SecureDetails)
}
