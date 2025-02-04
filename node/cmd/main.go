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

package main

import (
	"flag"
	"fmt"
	"os"

	driver "github.com/ibm/ibm-block-csi-driver/node/pkg/driver"
	"k8s.io/klog"
)

func main() {
	var (
		endpoint = flag.String("csi-endpoint", "unix://csi/csi.sock", "CSI Endpoint")
		version  = flag.Bool("version", false, "Print the version and exit.")
		configFile  = flag.String("config-file-path",  "./common/config.yaml", "Shared config file.")
		hostname  = flag.String("hostname",  "host-dns-name", "The name of the host the node is running on.")
	)

	klog.InitFlags(nil)
	flag.Parse()
	

	if *version {
		info, err := driver.GetVersionJSON(*configFile)
		if err != nil {
			klog.Fatalln(err)
		}
		fmt.Println(info)
		os.Exit(0)
	}

	drv, err := driver.NewDriver(*endpoint, *configFile, *hostname)
	if err != nil {
		klog.Fatalln(err)
	}
	if err := drv.Run(); err != nil {
		klog.Fatalln(err)
	}
}
