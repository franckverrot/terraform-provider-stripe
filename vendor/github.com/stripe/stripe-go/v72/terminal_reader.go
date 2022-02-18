//
//
// File generated from our OpenAPI spec
//
//

package stripe

// Type of reader, one of `bbpos_chipper2x`, `bbpos_wisepos_e`, or `verifone_P400`.
type TerminalReaderDeviceType string

// List of values that TerminalReaderDeviceType can take
const (
	TerminalReaderDeviceTypeBBPOSChipper2X TerminalReaderDeviceType = "bbpos_chipper2x"
	TerminalReaderDeviceTypeBBPOSWisePOSE  TerminalReaderDeviceType = "bbpos_wisepos_e"
	TerminalReaderDeviceTypeVerifoneP400   TerminalReaderDeviceType = "verifone_P400"
)

// Updates a Reader object by setting the values of the parameters passed. Any parameters not provided will be left unchanged.
type TerminalReaderParams struct {
	Params `form:"*"`
	// Custom label given to the reader for easier identification. If no label is specified, the registration code will be used.
	Label *string `form:"label"`
	// The location to assign the reader to.
	Location *string `form:"location"`
	// A code generated by the reader used for registering to an account.
	RegistrationCode *string `form:"registration_code"`
}

// TerminalReaderGetParams is the set of parameters that can be used to get a terminal reader.
type TerminalReaderGetParams struct {
	Params `form:"*"`
}

// Returns a list of Reader objects.
type TerminalReaderListParams struct {
	ListParams `form:"*"`
	// Filters readers by device type
	DeviceType *string `form:"device_type"`
	// A location ID to filter the response list to only readers at the specific location
	Location *string `form:"location"`
	// A status filter to filter readers to only offline or online readers
	Status *string `form:"status"`
}

// A Reader represents a physical device for accepting payment details.
//
// Related guide: [Connecting to a Reader](https://stripe.com/docs/terminal/payments/connect-reader).
type TerminalReader struct {
	APIResource
	Deleted bool `json:"deleted"`
	// The current software version of the reader.
	DeviceSwVersion string `json:"device_sw_version"`
	// Type of reader, one of `bbpos_chipper2x`, `bbpos_wisepos_e`, or `verifone_P400`.
	DeviceType TerminalReaderDeviceType `json:"device_type"`
	// Unique identifier for the object.
	ID string `json:"id"`
	// The local IP address of the reader.
	IPAddress string `json:"ip_address"`
	// Custom label given to the reader for easier identification.
	Label string `json:"label"`
	// Has the value `true` if the object exists in live mode or the value `false` if the object exists in test mode.
	Livemode bool `json:"livemode"`
	// The location identifier of the reader.
	Location string `json:"location"`
	// Set of [key-value pairs](https://stripe.com/docs/api/metadata) that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
	Metadata map[string]string `json:"metadata"`
	// String representing the object's type. Objects of the same type share the same value.
	Object string `json:"object"`
	// Serial number of the reader.
	SerialNumber string `json:"serial_number"`
	// The networking status of the reader.
	Status string `json:"status"`
}

// TerminalReaderList is a list of Readers as retrieved from a list endpoint.
type TerminalReaderList struct {
	APIResource
	ListMeta
	Data     []*TerminalReader `json:"data"`
	Location *string           `json:"location"`
	Status   *string           `json:"status"`
}