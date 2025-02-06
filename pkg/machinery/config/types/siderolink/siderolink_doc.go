// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Code generated by hack/docgen tool. DO NOT EDIT.

package siderolink

import (
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
)

func (ConfigV1Alpha1) Doc() *encoder.Doc {
	doc := &encoder.Doc{
		Type:        "SideroLinkConfig",
		Comments:    [3]string{"" /* encoder.HeadComment */, "SideroLinkConfig is a SideroLink connection machine configuration document." /* encoder.LineComment */, "" /* encoder.FootComment */},
		Description: "SideroLinkConfig is a SideroLink connection machine configuration document.",
		Fields: []encoder.Doc{
			{}, {
				Name:        "apiUrl",
				Type:        "URL",
				Note:        "",
				Description: "SideroLink API URL to connect to.",
				Comments:    [3]string{"" /* encoder.HeadComment */, "SideroLink API URL to connect to." /* encoder.LineComment */, "" /* encoder.FootComment */},
			},
		},
	}

	doc.AddExample("", exampleConfigV1Alpha1())

	doc.Fields[1].AddExample("", "https://siderolink.api/?jointoken=secret")

	return doc
}

// GetFileDoc returns documentation for the file ./siderolink_doc.go.
func GetFileDoc() *encoder.FileDoc {
	return &encoder.FileDoc{
		Name:        "siderolink",
		Description: "Package siderolink provides SideroLink machine configuration documents.\n",
		Structs: []*encoder.Doc{
			ConfigV1Alpha1{}.Doc(),
		},
	}
}
