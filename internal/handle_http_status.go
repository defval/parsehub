package internal

import "errors"

// Check HTTP status code
func CheckHTTPStatusCode(statusCode int) (bool, error) {
	switch statusCode {
	case 400:
		return false, errors.New("Bad request. Not able to get data from parsehub.")
	case 401:
		return false, errors.New("Unauthorized access. Not able to get data from parsehub. Please check api key.")
	case 403:
		return false, errors.New("Forbidden. Not able to get data from parsehub. Please check api key.")
	}

	return true, nil
}
