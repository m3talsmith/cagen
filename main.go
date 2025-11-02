package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path"
	"time"
)

var (
	caOutputPath         string
	caCommonName         string
	caOrganization       string
	caOrganizationalUnit string
	caAddress            string
	caLocality           string
	caProvince           string
	caPostalCode         string
	caCountry            string
	caKeySize            int
	caDaysToExpire       int
)

func init() {
	flag.StringVar(&caOutputPath, "output-path", "./", "The path to output the CA certificate")
	flag.StringVar(&caCommonName, "common-name", "", "The common name of the CA")
	flag.StringVar(&caOrganization, "organization", "", "The organization of the CA")
	flag.StringVar(&caOrganizationalUnit, "organizational-unit", "", "The organizational unit of the CA")
	flag.StringVar(&caAddress, "address", "", "The address of the CA")
	flag.StringVar(&caLocality, "locality", "", "The locality of the CA")
	flag.StringVar(&caProvince, "province", "", "The province of the CA")
	flag.StringVar(&caPostalCode, "postal-code", "", "The postal code of the CA")
	flag.StringVar(&caCountry, "country", "", "The country of the CA")
	flag.IntVar(&caKeySize, "key-size", 2048, "The key size of the CA")
	flag.IntVar(&caDaysToExpire, "days-to-expire", 365, "The number of days to expire the CA")
	flag.Parse()

	if caCommonName == "" {
		fmt.Println("Error: ca-common-name is required")
		os.Exit(1)
	}

	if caOrganization == "" {
		fmt.Println("Error: ca-organization is required")
		os.Exit(1)
	}

	if caOrganizationalUnit == "" {
		fmt.Println("Error: ca-organizational-unit is required")
		os.Exit(1)
	}

	if caAddress == "" {
		fmt.Println("Error: ca-address is required")
		os.Exit(1)
	}

	if caLocality == "" {
		fmt.Println("Error: ca-locality is required")
		os.Exit(1)
	}

	if caProvince == "" {
		fmt.Println("Error: ca-province is required")
		os.Exit(1)
	}

	if caPostalCode == "" {
		fmt.Println("Error: ca-postal-code is required")
		os.Exit(1)
	}

	if caCountry == "" {
		fmt.Println("Error: ca-country is required")
		os.Exit(1)
	}

	if caKeySize <= 1024 {
		fmt.Println("Error: ca-key-size must be greater than 1024")
		os.Exit(1)
	}

	if caKeySize%1024 != 0 {
		fmt.Println("Error: ca-key-size must be a multiple of 1024")
		os.Exit(1)
	}

	if caDaysToExpire <= 90 {
		fmt.Println("Error: ca-days-to-expire must be greater than 90")
		os.Exit(1)
	}
}

func main() {
	t := time.Now()
	ca := x509.Certificate{
		SerialNumber: big.NewInt(int64(t.Year())),
		Subject: pkix.Name{
			CommonName:         caCommonName,
			Organization:       []string{caOrganization},
			OrganizationalUnit: []string{caOrganizationalUnit},
			Locality:           []string{caLocality},
			Province:           []string{caProvince},
			PostalCode:         []string{caPostalCode},
			Country:            []string{caCountry},
		},
		NotBefore:             t,
		NotAfter:              t.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, caKeySize)
	if err != nil {
		fmt.Println("Error: failed to generate CA private key")
		os.Exit(1)
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &ca, &ca, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Println("Error: failed to create CA certificate")
		os.Exit(1)
	}

	caPEM := bytes.NewBuffer(nil)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate,
	})

	caPrivateKeyPEM := bytes.NewBuffer(nil)
	pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	pem, err := tls.X509KeyPair(caPEM.Bytes(), caPrivateKeyPEM.Bytes())
	if err != nil {
		fmt.Println("Error: failed to create CA pem")
		os.Exit(1)
	}

	os.WriteFile(path.Join(caOutputPath, "ca.pem"), []byte(fmt.Sprintf("%s\n%s", caPEM.String(), caPrivateKeyPEM.String())), 0644)
	os.WriteFile(path.Join(caOutputPath, "ca.crt"), pem.Certificate[0], 0644)
}
