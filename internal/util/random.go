// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

//RandomUserID generates a random user ID from 1->100,000,000 for use in mock events
func RandomUserID() string {
	uid, err := rand.Int(rand.Reader, big.NewInt(1*100*100*100*100))
	if err != nil {
		log.Fatal(err.Error())
	}
	return uid.String()
}

//RandomGUID generates a random GUID for use with creating IDs in the local store and for mock events
func RandomGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}
