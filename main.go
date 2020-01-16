// Copyright 2020 Jaime Soriano Pastor
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

var (
	flagSection string
	flagUpdate  bool
)

func init() {
	flag.StringVar(&flagSection, "section", "", "section to update")
	flag.BoolVar(&flagUpdate, "u", false, "update ~/.aws/credentials")
}

type StsOutput struct {
	Credentials struct {
		AccessKeyId     string
		SecretAccessKey string
		SessionToken    string
		Expiration      time.Time
	}
}

func main() {
	flag.Parse()
	if flagSection == "" {
		log.Fatal("-section flag needed")
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(usr.HomeDir, ".aws/credentials")

	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var sts StsOutput
	err = json.Unmarshal(d, &sts)
	if err != nil {
		log.Fatal(err)
	}

	creds := sts.Credentials

	section := cfg.Section(flagSection)
	section.Key("aws_access_key_id").SetValue(creds.AccessKeyId)
	section.Key("aws_secret_access_key").SetValue(creds.SecretAccessKey)
	section.Key("aws_session_token").SetValue(creds.SessionToken)

	if flagUpdate {
		err := cfg.SaveTo(path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cfg.WriteTo(os.Stdout)
	}
}
