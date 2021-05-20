package apierrors

type ErrorStructure struct {
	Errors []IndividualError `json:"errors"`
}

type IndividualError struct {
	SpecificError string `json:"error"`
	Message       string `json:"message"`
	Source        Source `json:"source"`
}

type Source struct {
	Field string `json:"field"`
	Param string `json:"param"`
}
