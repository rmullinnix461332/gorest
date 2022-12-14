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
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//Marshals the data in interface i into a byte slice, using the Marhaller/Unmarshaller specified in mime.
//The Marhaller/Unmarshaller must have been registered before using gorest.RegisterMarshaller
func Marshal(i interface{}, mime string) (io.ReadCloser, error) {
	return interfaceToBytes(i, mime)
}

//Marshals the data in interface i into a byte slice, using the Marhaller/Unmarshaller specified in mime.
//The Marhaller/Unmarshaller must have been registered before using gorest.RegisterMarshaller
func interfaceToBytes(i interface{}, mime string) (io.ReadCloser, error) {
	marshalType := mime

	if strings.Contains(mime, "json") {
		marshalType = "json"
	} else if strings.Contains(mime, "xml") {
		marshalType = "xml"
	}

	m := GetMarshallerByMime(marshalType)
	if m != nil {
		return m.Marshal(i)
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		x := v.Bool()
		if x {
			return ioutil.NopCloser(bytes.NewBuffer([]byte("true"))), nil
		}
		return ioutil.NopCloser(bytes.NewBuffer([]byte("false"))), nil
	case reflect.String:
		return ioutil.NopCloser(bytes.NewBuffer([]byte(v.String()))), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ioutil.NopCloser(bytes.NewBuffer([]byte(strconv.FormatInt(v.Int(), 10)))), nil
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		m := GetMarshallerByMime(marshalType)
		return m.Marshal(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return ioutil.NopCloser(bytes.NewBuffer([]byte(strconv.FormatUint(v.Uint(), 10)))), nil
	case reflect.Float32, reflect.Float64:
		return ioutil.NopCloser(bytes.NewBuffer([]byte(strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits())))), nil
	case reflect.Invalid:
		return nil, errors.New("Type is invalid!")
	default:
		return nil, errors.New("Type " + v.Type().Name() + " is not handled by GoRest.")
	}
}

//Unmarshals the data in buf into interface i, using the Marhaller/Unmarshaller specified in mime.
//The Marhaller/Unmarshaller must have been registered before using gorest.RegisterMarshaller
func Unmarshal(buf *bytes.Buffer, i interface{}, mime string) error {
	return bytesToInterface(buf, i, mime)
}

func bytesToInterface(buf *bytes.Buffer, i interface{}, mime string) error {
	marshalType := mime

	if strings.Contains(mime, "json") {
		marshalType = "json"
	} else if strings.Contains(mime, "xml") {
		marshalType = "xml"
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:

		n, err := strconv.ParseBool(buf.String())
		if err != nil {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetBool(n)
		break
	case reflect.String:
		reflect.ValueOf(i).Elem().SetString(buf.String())
		break
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		m := GetMarshallerByMime(marshalType)
		return m.Unmarshal(buf.Bytes(), i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		n, err := strconv.ParseInt(buf.String(), 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetInt(n)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(buf.String(), 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errors.New("Invalid value. " + err.Error())
		}

		reflect.ValueOf(i).Elem().SetUint(n)
		break
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(buf.String(), v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetFloat(n)
	default:
		return errors.New("Type " + v.Type().Name() + " is not handled by GoRest.")
	}
	return nil

}
