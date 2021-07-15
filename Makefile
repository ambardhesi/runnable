.PHONY: binaries clean test vendor ca_key ca_cert svr_key svr_csr svr_cert alice_key alice_csr alice_cert \
	bob_key bob_csr bob_cert bad_key bad_cert

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
ROOT=$(shell pwd)
COMMANDS=$(shell go list ./... | grep cmd/)

OPENSSL?=/usr/bin/openssl
CERT_DIR=$(ROOT)/certs
CA_CERT="$(CERT_DIR)/ca-cert.pem"
CA_KEY="$(CERT_DIR)/ca-key.pem"
SVR_CERT="$(CERT_DIR)/svr-cert.pem"
SVR_KEY="$(CERT_DIR)/svr-key.pem"
SVR_CSR="$(CERT_DIR)/svr-csr.csr"
ALICE_CERT="$(CERT_DIR)/alice-cert.pem"
ALICE_KEY="$(CERT_DIR)/alice-key.pem"
ALICE_CSR="$(CERT_DIR)/alice-csr.csr"
BOB_CERT="$(CERT_DIR)/bob-cert.pem"
BOB_KEY="$(CERT_DIR)/bob-key.pem"
BOB_CSR="$(CERT_DIR)/bob-csr.csr"
BAD_CERT="$(CERT_DIR)/bad-cert.pem"
BAD_KEY="$(CERT_DIR)/bad-key.pem"

binaries:
	$(GOBUILD) -o "." $(COMMANDS) 	

clean:
	rm -rf vendor
	rm certs/*.pem
	rm certs/*.csr
	rm certs/*.srl
	$(GOCLEAN)

test: vendor
	$(GOTEST) -v -race ./...

vendor:
	$(GOCMD) mod vendor

certs: ca_cert svr_cert alice_cert bob_cert bad_key bad_cert

ca_key: 
	$(OPENSSL) genpkey -algorithm ed25519 -out $(CA_KEY) -outform PEM

ca_cert: ca_key
	$(OPENSSL) req -x509 -newkey rsa:4096 -key $(CA_KEY) -out $(CA_CERT) -subj "/C=CA/ST=BC/L=Vancouver/OU=SignedCA/CN=localhost/emailAddress=foo@foo.com"

svr_key:
	$(OPENSSL) genpkey -algorithm ed25519 -out $(SVR_KEY) -outform PEM

svr_csr: svr_key
	$(OPENSSL) req -new -key $(SVR_KEY) -out $(SVR_CSR) -config "$(CERT_DIR)/svr.conf" 

svr_cert: svr_csr
	$(OPENSSL) x509 -req -in $(SVR_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SVR_CERT) -extfile $(CERT_DIR)/svr.conf -extensions my_extensions

alice_key:
	$(OPENSSL) genpkey -algorithm ed25519 -out $(ALICE_KEY) -outform PEM

alice_csr: alice_key
	$(OPENSSL) req -new -key $(ALICE_KEY) -out $(ALICE_CSR) -subj "/CN=alice"

alice_cert: alice_csr
	$(OPENSSL) x509 -req -in $(ALICE_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(ALICE_CERT) 

bob_key: 
	$(OPENSSL) genpkey -algorithm ed25519 -out $(BOB_KEY) -outform PEM

bob_csr: bob_key
	$(OPENSSL) req -new -key $(BOB_KEY) -out $(BOB_CSR) -subj "/CN=bob"

bob_cert: bob_csr
	$(OPENSSL) x509 -req -in $(BOB_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(BOB_CERT) 

bad_key:
	$(OPENSSL) genpkey -algorithm ed25519 -out $(BAD_KEY) -outform PEM

bad_cert: bad_key
	$(OPENSSL) req -new -key $(BAD_KEY) -out $(BAD_CERT) -subj "/CN=peter"
