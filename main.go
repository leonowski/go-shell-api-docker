package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var allowedCommands string

func generateSelfSignedCert() (tls.Certificate, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

	cert, err := tls.X509KeyPair(certPEM, privKeyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

func main() {
	binary := flag.String("b", "", "Path to the executable binary")
	port := flag.Int("p", 8080, "HTTP port to listen on")
	username := flag.String("u", "", "Basic authentication username")
	password := flag.String("pw", "", "Basic authentication password")
	useHTTPS := flag.Bool("https", false, "Serve via HTTPS using a self-signed certificate or an optional custom certificate")
	certFile := flag.String("cert", "", "Path to the server certificate file (requires -https flag)")
	keyFile := flag.String("key", "", "Path to the server key file (requires -https flag)")
	listenAddr := flag.String("l", "0.0.0.0", "Address to listen on")
	flag.Parse()

	if *binary == "" {
		fmt.Println("Path to binary not specified.")
		return
	}

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *username != "" && *password != "" {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			authDecoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
			if err != nil || string(authDecoded) != fmt.Sprintf("%s:%s", *username, *password) {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
		}

		var argString string
		if r.Body != nil {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				l.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			argString = string(data)
		}

		fields := strings.Fields(*binary)
		args := append(fields[1:], strings.Fields(argString)...)

		command := fmt.Sprintf("%s %s", fields[0], strings.Join(args, " "))
		l.Printf("Command: [%s]", command)

		if !isValidCommand(command) {
			http.Error(w, "Invalid command.", http.StatusBadRequest)
			return
		}

		if !isCommandAllowed(command) {
			http.Error(w, "Command not allowed.", http.StatusForbidden)
			return
		}

		output, err := exec.Command(fields[0], args...).Output()

		l.Printf("Command: [%s]", command)

		output, err = exec.Command(fields[0], args...).Output()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(output)
	})

	if *useHTTPS {
		var tlsConfig *tls.Config
		if *certFile != "" && *keyFile != "" {
			cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
			if err != nil {
				l.Fatal(err)
			}
			tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		} else {
			cert, err := generateSelfSignedCert()
			if err != nil {
				l.Fatal(err)
			}
			tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		}

		server := &http.Server{
			Addr:      fmt.Sprintf("%s:%d", *listenAddr, *port),
			TLSConfig: tlsConfig,
		}

		l.Printf("Listening on %s:%d using HTTPS...", *listenAddr, *port)
		l.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		l.Printf("Listening on %s:%d using HTTP...", *listenAddr, *port)
		l.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddr, *port), nil))
	}
}

func isCommandAllowed(command string) bool {
	if allowedCommands == "" {
		return true
	}

	allowedCommandsSlice := strings.Split(allowedCommands, ",")
	for _, allowedCommand := range allowedCommandsSlice {
		if strings.HasPrefix(command, allowedCommand) {
			return true
		}
	}
	return false
}

// this is an attempt to sanitize shell commands that may cause a sort of fork bomb.
func isValidCommand(command string) bool {
	dangerousSequences := []string{"|", ";", "&", "$(", "`", ">/dev/null", ">/dev/random"}

	for _, seq := range dangerousSequences {
		if strings.Contains(command, seq) {
			return false
		}
	}
	return true
}
