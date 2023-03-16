// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// mkconsts generates const definition files for each GOEXPERIMENT.
package main

import (
	"bytes"
	"fmt"
	"github.com/ryan-berger/golang_tv/internal/src/internal/goexperiment"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	// Delete existing experiment constant files.
	ents, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, ent := range ents {
		name := ent.Name()
		if !strings.HasPrefix(name, "exp_") {
			continue
		}
		// Check that this is definitely a generated file.
		data, err := os.ReadFile(name)
		if err != nil {
			log.Fatalf("reading %s: %v", name, err)
		}
		if !bytes.Contains(data, []byte("Code generated by mkconsts")) {
			log.Fatalf("%s: expected generated file", name)
		}
		if err := os.Remove(name); err != nil {
			log.Fatal(err)
		}
	}

	// Generate new experiment constant files.
	rt := reflect.TypeOf(&goexperiment.Flags{}).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i).Name
		buildTag := "goexperiment." + strings.ToLower(f)
		for _, val := range []bool{false, true} {
			name := fmt.Sprintf("exp_%s_%s.go", strings.ToLower(f), pick(val, "off", "on"))
			data := fmt.Sprintf(`// Code generated by mkconsts.go. DO NOT EDIT.

//go:build %s%s
// +build %s%s

package goexperiment

const %s = %v
const %sInt = %s
`, pick(val, "!", ""), buildTag, pick(val, "!", ""), buildTag, f, val, f, pick(val, "0", "1"))
			if err := os.WriteFile(name, []byte(data), 0666); err != nil {
				log.Fatalf("writing %s: %v", name, err)
			}
		}
	}
}

func pick(v bool, f, t string) string {
	if v {
		return t
	}
	return f
}
