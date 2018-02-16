package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	rcli "github.com/gomeet-examples/svc-profile/cmd/remotecli"
)

// consoleCmd represents the console command
var (
	consoleAddress string
	consoleCmd     = &cobra.Command{
		Use:   "console",
		Short: "Interactive console on svc-profile service",
		Long: `Use this command for prompt a interactive console
`,
		Run: console,
	}
)

func init() {
	RootCmd.AddCommand(consoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Port flag
	consoleCmd.PersistentFlags().StringVar(&consoleAddress, "address", "localhost:13000", "Address to server")

	// CA certificate
	consoleCmd.PersistentFlags().StringVar(&caCertificate, "ca", "", "X.509 certificate of the Certificate Authority (required for gRPC TLS support)")

	// gRPC client certificate
	consoleCmd.PersistentFlags().StringVar(&clientCertificate, "cert", "", "X.509 certificate of the gRPC client (required for gRPC TLS support)")

	// gRPC client private key
	consoleCmd.PersistentFlags().StringVar(&clientPrivateKey, "key", "", "RSA private key of the gRPC client (required for gRPC TLS support)")

	// gRPC timeout
	consoleCmd.PersistentFlags().IntVar(&timeoutSeconds, "timeout", 5, "gRPC timeout in seconds")

	// JSON Web Token
	consoleCmd.PersistentFlags().StringVar(&jwtToken, "jwt", "", "JSON Web Token")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func console(cmd *cobra.Command, args []string) {
	log.WithFields(log.Fields{
		"Name":     svc.Name,
		"Version":  svc.Version,
		"FullName": svcName,
	}).Info("console version")

	// Initialize remote cli
	c, err := rcli.NewRemoteCli(
		svc.Name,
		svc.Version,
		rcli.ConsoleCall,
		consoleAddress,
		timeoutSeconds,
		caCertificate,
		clientCertificate,
		clientPrivateKey,
		jwtToken,
	)

	if err != nil {
		log.Fatalf("Remote cli init fail - %v", err)
	}
	// Defer remote cli closing
	defer c.Close()

	// Get remote version for use in prompt
	var prompt string
	rVersion, err := c.RemoteVersion()
	if err != nil {
		log.Warnf("Get remote version fail : %v", err)
		prompt = svc.Name
	} else {
		log.WithFields(log.Fields{
			"Name":    rVersion.Name,
			"Version": rVersion.Version,
		}).Info("remote version")
		prompt = fmt.Sprintf("%s-%s@%s", rVersion.Name, rVersion.Version, consoleAddress)
	}

	var pfxCompl []readline.PrefixCompleterInterface

	for k, _ := range c.GetActionsMap() {
		pfxCompl = append(pfxCompl, readline.PcItem(k))
	}
	for _, v := range []string{"exit"} {
		pfxCompl = append(pfxCompl, readline.PcItem(v))
	}

	completer := readline.NewPrefixCompleter(pfxCompl...)

	// Set up interactive console
	cfg := &readline.Config{
		// Prompt definition
		Prompt:          fmt.Sprintf("└─┤%s├─$ ", prompt),
		HistoryFile:     fmt.Sprintf("/tmp/%s-%d.tmp", svc.Name, os.Getpid()),
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	}

	l, err := readline.NewEx(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.WithFields(log.Fields{
		"HistoryFile": cfg.HistoryFile,
		"Interrupt":   cfg.InterruptPrompt,
		"Exit":        cfg.EOFPrompt,
	}).Info("load console")

	// REPL

	log.SetOutput(l.Stderr())
	for {
		// read interaction
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			log.Info("")
			continue
		}

		// Eval/Print
		switch {
		case line == "exit":
			goto exit

		default:
			// Evaluate arguments string
			ok, err := c.Eval(line)
			if err != nil {
				log.Warn(err)
				break
			}
			log.Info(ok)
		} // loop
	}
exit:
}
