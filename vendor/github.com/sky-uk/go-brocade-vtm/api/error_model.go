package api

// VTMError : Generic error object for Brocade vTM
type VTMError struct {
	VTMError VTMErrorInfo `json:"error"`
}

// VTMErrorInfo : Generic error info object for Brocade vTM
type VTMErrorInfo struct {
	ErrorID   string                 `json:"error_id"`
	ErrorText string                 `json:"error_text"`
	ErrorInfo map[string]interface{} `json:"error_info,omitempty"`
}
