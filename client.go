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


// Notice: This code has been modified from its original source.
// Modifications are licensed as specified below.
//
// Copyright (c) 2014, fromkeith
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice, this
//   list of conditions and the following disclaimer in the documentation and/or
//   other materials provided with the distribution.
//
// * Neither the name of the fromkeith nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
///



package gorest

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

var sharedClient *http.Client

//Use this if you have a *http.Client instance that you specifically want to use. 
//Otherwise just use NewRequestBuilder(), which uses the http.Client maintained by GoRest.
func NewRequestBuilderFromClient(client *http.Client, url string) (*RequestBuilder, error) {
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	rb := RequestBuilder{client, Application_Json, req}
	return &rb, nil
}

//This creates a new RequestBuilder, backed by GoRest's internally managed http.Client.
//Although http.Client is useable concurrently, an instance of RequestBuilder is not safe for this. 
//Because of http.Client's persistent(cached TCP connections) and concurrent nature, 
//this can be used safely multiple times from different go routines. 
func NewRequestBuilder(url string) (*RequestBuilder, error) {
	if sharedClient == nil {
		sharedClient = new(http.Client) //DefaultClient
	}
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	rb := RequestBuilder{sharedClient, Application_Json, req}
	return &rb, nil
}

type RequestBuilder struct {
	client             *http.Client
	defaultContentType string
	_req               *http.Request
}

func (this *RequestBuilder) Request() *http.Request {
	return this._req
}
func (this *RequestBuilder) UseContentType(mime string) *RequestBuilder {
	this.defaultContentType = mime
	return this
}

func (this *RequestBuilder) AddCookie(cookie *http.Cookie) *RequestBuilder {
	this._req.AddCookie(cookie)
	return this
}

func (this *RequestBuilder) Delete() (*http.Response, error) {
	//	u, err := url.Parse(url_)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	//this._req.URL = u
	this._req.Method = DELETE

	return this.client.Do(this._req)
}

func (this *RequestBuilder) Head() (*http.Response, error) {
	this._req.Method = HEAD
	return this.client.Do(this._req)
}

func (this *RequestBuilder) Options(opts *[]string) (*http.Response, error) {
	this._req.Method = OPTIONS

	res, err := this.client.Do(this._req)
	if err != nil {
		return res, err
	}

	for _, str := range res.Header["Allow"] {
		*opts = append(*opts, strings.Trim(str, " "))
	}
	return res, err
}

func (this *RequestBuilder) Get(i interface{}, expecting int) (*http.Response, error) {
	//this._req.URL = u
	this._req.Method = GET

	res, err := this.client.Do(this._req)
	if err != nil {
		return res, err
	}

	if res.StatusCode == expecting {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		res.Body.Close()
		err = bytesToInterface(buf, i, this.defaultContentType)
		return res, nil
	}

	return res, errors.New(res.Status)
}
func (this *RequestBuilder) Post(i interface{}) (*http.Response, error) {
	return this.PostWithResponse(i, nil)

}

func (this *RequestBuilder) PostWithResponse(i interface{}, output interface{}) (*http.Response, error) {
	this._req.Method = POST
	bb, err := interfaceToBytes(i, this.defaultContentType)
	if err != nil {
		return nil, err
	}
	this._req.Body = bb

	res, err := this.client.Do(this._req)
	if err != nil {
		return res, err
	}


	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	res.Body.Close()
	if buf.Len() > 0 && output != nil {
		err = bytesToInterface(buf, output, this.defaultContentType)
		return res, err
	}
	return res, nil

}
