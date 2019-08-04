/**
 * Copyright 2019 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package driver

import (
	"fmt"
)

type ConfigYmlEmptyAttribute struct {
	Attr string
}

func (e *ConfigYmlEmptyAttribute) Error() string {
	return fmt.Sprintf("Missing attribute [%s] in driver config yaml file", e.Attr)
}

type RequestValidationError struct {
	Msg string
}

func (e *RequestValidationError) Error() string {
	return fmt.Sprintf("Request Validation Error: %s", e.Msg)
}


type ConnectivityIscsiStorageTargetNotFoundError struct {
	StorageTargetName string
	DirectoryPath string
}

func (e *ConnectivityIscsiStorageTargetNotFoundError) Error() string {
	return fmt.Sprintf("Connectivity Error: Storage target name [%s] was not found on the host, under directory %s", e.StorageTargetName, e.DirectoryPath)
}


type VolumeAlreadyProcessingError struct {
	volId string
}

func (e *VolumeAlreadyProcessingError) Error() string {
	return fmt.Sprintf("Volume %s is already processing. request cannot be completed.", e.volId)
}
