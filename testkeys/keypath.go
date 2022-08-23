///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package testkeys

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/utils"
	"path/filepath"
	"runtime"
)

func getDirForFile() string {
	// Get the filename we're in
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

// These functions are used to cover TLS connection code in tests
func GetNodeCertPath() string {
	return filepath.Join(getDirForFile(), "cmix.rip.crt")
}

func GetNodeKeyPath() string {
	return filepath.Join(getDirForFile(), "cmix.rip.key")
}

func GetGatewayCertPath() string {
	return filepath.Join(getDirForFile(), "gateway.cmix.rip.crt")
}

func GetGatewayKeyPath() string {
	return filepath.Join(getDirForFile(), "gateway.cmix.rip.key")
}

func LoadFromPath(path string) []byte {
	// Check if the path is an embedded file
	switch filepath.Base(path) {
	case "cmix.rip.crt":
		return CmixCrt
	case "cmix.rip.key":
		return CmixKey
	case "gateway.cmix.rip.crt":
		return GatewayCrt
	case "gateway.cmix.rip.key":
		return GatewayKey
	}

	data, err := utils.ReadFile(path)
	if err != nil {
		jww.FATAL.Panicf("***Check your key!***\nFailed to read file at %s: %+v", path, err)
	}
	return data
}
