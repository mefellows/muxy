package command

import (
	"flag"
	"fmt"
	"github.com/mefellows/pkigo/pki"
	"github.com/mitchellh/cli"
	"os"
	"strings"
	"time"
)

type PkiCommand struct {
	meta             Meta
	caHost           string
	outputCA         bool
	importClientCert string
	importClientKey  string
	outputClientCert bool
	outputClientKey  bool
	importCA         string
	generateCert     bool
	configure        bool
	removePKI        bool
}

func (c *PkiCommand) Run(args []string) int {
	c.meta = Meta{
		Ui: &cli.ColoredUi{
			Ui:          &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin, ErrorWriter: os.Stderr},
			OutputColor: cli.UiColorNone,
			InfoColor:   cli.UiColorNone,
			ErrorColor:  cli.UiColorRed,
		},
	}
	cmdFlags := flag.NewFlagSet("pki", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.meta.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.caHost, "caHost", "localhost", "Specify the CAs custom hostname")
	cmdFlags.StringVar(&c.importCA, "importCA", "", "Path to CA Cert to import")
	cmdFlags.StringVar(&c.importClientCert, "importClientCert", "", "Path of client certificate to import and set as the default")
	cmdFlags.StringVar(&c.importClientKey, "importClientKey", "", "Path of client key to import and set as the default")
	cmdFlags.BoolVar(&c.configure, "configure", false, "Configures a default PKI infrastructure. Warning: This will clear any existing PKI files")
	cmdFlags.BoolVar(&c.removePKI, "removePKI", false, "Remove existing PKI keys and certs.")
	cmdFlags.BoolVar(&c.outputCA, "outputCA", false, "Output the CA Certificate of this node")
	cmdFlags.BoolVar(&c.outputClientCert, "outputClientCert", false, "Output the Client Certificate")
	cmdFlags.BoolVar(&c.outputClientKey, "outputClientKey", false, "Output the Client Key")
	cmdFlags.BoolVar(&c.generateCert, "generateCert", false, "Generate a custom cert from this nodes' CA")

	pki, err := pki.New()
	if err != nil {
		c.meta.Ui.Error(fmt.Sprintf("Unable to setup public key infrastructure: %s", err.Error()))
		return 1
	}

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if c.configure {
		c.meta.Ui.Output(fmt.Sprintf("Setting up PKI for %s...", c.caHost))
		pki.RemovePKI()
		err := pki.SetupPKI(c.caHost)
		if err != nil {
			c.meta.Ui.Error(err.Error())
		}
		c.meta.Ui.Output("PKI setup complete.")
	}

	if c.importCA != "" {
		c.meta.Ui.Output(fmt.Sprintf("Importing CA from %s", c.importCA))
		timestamp := time.Now().Unix()
		err := pki.ImportCA(fmt.Sprintf("%d", timestamp), c.importCA)
		if err != nil {
			c.meta.Ui.Error(fmt.Sprintf("Failed to import CA: %s", err.Error()))
		} else {
			c.meta.Ui.Info("CA successfully imported")
		}
	}

	if c.importClientCert != "" && c.importClientKey != "" {
		err := pki.ImportClientCertAndKey(c.importClientCert, c.importClientKey)
		if err != nil {
			c.meta.Ui.Error(fmt.Sprintf("Failed to import client keys: %s", err.Error()))
		} else {
			c.meta.Ui.Info("Client keys successfully imported")
		}
	}
	if c.outputCA {
		cert, _ := pki.OutputCACert()
		c.meta.Ui.Output(cert)
	}

	if c.outputClientCert {
		cert, _ := pki.OutputClientCert()
		c.meta.Ui.Output(cert)
	}

	if c.outputClientKey {
		cert, _ := pki.OutputClientKey()
		c.meta.Ui.Output(cert)
	}

	if c.removePKI {
		c.meta.Ui.Output("Removing existing PKI")
		err := pki.RemovePKI()
		if err != nil {
			c.meta.Ui.Error(err.Error())
		}
		c.meta.Ui.Output("PKI removal complete.")
	}

	if c.generateCert {
		c.meta.Ui.Output("Generating a new client cert")
		err := pki.GenerateClientCertificate([]string{"localhost"})
		if err != nil {
			c.meta.Ui.Error(err.Error())
		}
		c.meta.Ui.Output("Cert generation complete")
	}

	return 0
}

func (c *PkiCommand) Help() string {
	helpText := `
Usage: <application> pki [options] 

  Sets up the PKI infrastructure for secure communication.
  
Options:

  --configure                 (Re-)configure PKI infrastructure on this node. This is generally only required if something strange happens. 
  --caHost                    Specify a custom CA Host when generating the PKI.
  --importCA                  Trust the provided CA.
  --outputCA                  Output the CA Certificate for this node.
  --importClientCert          Import the current Client Certificate (.crt). Must be accompanied by --importClientKey.
  --importClientKey           Import the current Client Key (.pem) file. Must be accompanied by --importClientCert.
  --outputClientCert          Output the current Client Certificate (.crt).
  --outputClientKey           Output the current Client Key (.pem) file.
  --generateCert              Generate a client cert trusted by this nodes CA.
  --removePKI                 Removes existing PKI.
`

	return strings.TrimSpace(helpText)
}

func (c *PkiCommand) Synopsis() string {
	return "Setup the PKI infrastructure for secure communication"
}
