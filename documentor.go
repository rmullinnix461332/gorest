//Copyright 2014  (rmullinnix461332@gmail.com). All rights reserved.
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

var documentors map[string]*Documentor

//Signiture of functions to be used as Documentors
type Documentor struct {
	Document func(string, map[string]ServiceMetaData, map[string]EndPointStruct, map[string]SecurityStruct)(interface{})
}

//Registers an Documentor for the specified mime type
func RegisterDocumentor(mime string, dec *Documentor) {
	if documentors == nil {
		documentors = make(map[string]*Documentor, 0)
	}
	if _, found := documentors[mime]; !found {
		documentors[mime] = dec
	}
}

//Returns the registred documentor for the specified mime type
func GetDocumentor(mime string) (dec *Documentor) {
	if documentors == nil {
		documentors = make(map[string]*Documentor, 0)
	}
	dec, _ = documentors[mime]
	return
}
