package response

// JSON ///
type JSON struct {
	ServiceDescription string `json:"service_description"`
	ServiceVersion     string `json:"service_version"`
	ServiceRespone     data   `json:"service_response"`
}

// Data ///
type data struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

//Respone constants
const (
	OK     = "OK"
	FAILED = "Failed"

	description = "simple auth api"
	version     = "v1.1.0"
)

// NewJSON returns response data
func NewJSON(status, inputData string) JSON {
	return JSON{description, version, data{status, inputData}}
}
