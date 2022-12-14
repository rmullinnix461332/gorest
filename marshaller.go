//Copyright 2011 Siyabonga Dlamini (siyabonga.dlamini@gmail.com). All rights reserved.
//
//Redistribution and use in source and binary forms, with or without
//modification, are permitted provided that the following conditions
//are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer
//     in the documentation and/or other materials provided with the
//     distribution.
//
//THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
//IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
//OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
//IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
//PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
//OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
//WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
//OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
//ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gorest

import (
	"bytes"
	"io"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"github.com/ajg/form"
	"strings"
)

//A Marshaller represents the two functions used to marshal/unmarshal interfaces back and forth.
type Marshaller struct {
	Marshal   func(v interface{}) (io.ReadCloser, error)
	Unmarshal func(data []byte, v interface{}) error
}

var marshallers map[string]*Marshaller

//Register a Marshaller. These registered Marshallers are shared by the client or servers side usage of gorest.
func RegisterMarshaller(mime string, m *Marshaller) {
	if marshallers == nil {
		marshallers = make(map[string]*Marshaller, 0)
	}
	if _, found := marshallers[mime]; !found {
		marshallers[mime] = m
	}
}

//Get an already registered Marshaller
func GetMarshallerByMime(mime string) (m *Marshaller) {
	if marshallers == nil {
		marshallers = make(map[string]*Marshaller, 0)
	}
	m, _ = marshallers[mime]
	return
}

//Predefined Marshallers

//JSON: This makes the JSON Marshaller. The Marshaller uses pkg: json
func NewJSONMarshaller() *Marshaller {
	m := Marshaller{jsonMarshal, jsonUnMarshal}
	return &m
}
func jsonMarshal(v interface{}) (io.ReadCloser, error) {
	j, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	return ioutil.NopCloser(bytes.NewBuffer(j)), nil
}
func jsonUnMarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

//XML
func NewXMLMarshaller() *Marshaller {
	m := Marshaller{xmlMarshal, xmlUnMarshal}
	return &m
}
func xmlMarshal(v interface{}) (io.ReadCloser, error) {
	x, e := xml.Marshal(v)
	if e != nil {
		return nil, e
	}
	return ioutil.NopCloser(bytes.NewBuffer(x)), nil
}
func xmlUnMarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

//application/x-www-form-urlencoded
func NewFormMarshaller() * Marshaller {
	m := Marshaller{formMarshal, formUnMarshal}
	return &m
}
func formMarshal(v interface{}) (io.ReadCloser, error) {
	s, err := form.EncodeToString(v)
	return ioutil.NopCloser(strings.NewReader(s)), err
}
func formUnMarshal(data []byte, v interface{}) error {
	err := form.DecodeString(v, string(data))
	return err
}
