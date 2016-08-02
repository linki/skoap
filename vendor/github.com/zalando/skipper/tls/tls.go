// Copyright 2016 Zalando SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tls

import "crypto/x509"

// CertPool is an interface for *x509.CertPool.
type CertPool interface {
	Set(**x509.CertPool)
}

// DefaultCertPool is the default cert pool.
type DefaultCertPool struct{}

// Set sets the cert pool to the provided pool pointer.
func (d *DefaultCertPool) Set(pool **x509.CertPool) {
	*pool = nil
}
