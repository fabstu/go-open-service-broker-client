/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2

import (
	"fmt"
	"net/http"
)

func (c *client) GetInstance(r *GetInstanceRequest) (*GetInstanceResponse, error) {
	if err := c.validateAlphaAPIMethodsAllowed(); err != nil {
		return nil, GetInstanceNotAllowedError{
			reason: err.Error(),
		}
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)

	response, err := c.prepareAndDo(http.MethodGet, fullURL, nil /* params */, nil /* request body */, nil /* originating identity */)
	if err != nil {
		return nil, err
	}

	defer func() {
		drainReader(response.Body)
		response.Body.Close()
	}()

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &GetInstanceResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, HTTPStatusCodeError{StatusCode: response.StatusCode, ResponseError: err}
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}
