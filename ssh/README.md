# SSH

## SSH tunnel example

```go
import (
	"fmt"
	"github.com/manifoldco/promptui"
	netssh "github.com/pete911/go-net/ssh"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"strings"
)

func GetRequest() {

	tunnel := NewTunnel("root", "bastion.com", 22, "site.internal", 80)
	http.Get(fmt.Sprintf("%s:%d", tunnel.Host, tunnel.Port))
}

func NewTunnel(bastionUsername, bastionHost string, bastionPort int, host string, port int) netssh.Endpoint {

	bastionEndpoint := netssh.NewEndpoint(bastionUsername, bastionHost, bastionPort)
	ldapEndpoint := netssh.NewEndpoint("", host, port)

	tunnel, err := netssh.NewTunnel(bastionEndpoint, ldapEndpoint, SSHKeyboardAuth())
	if err != nil {
		log.Printf("ssh tunnel: %v", err)
	}
	go tunnel.Start()
	return tunnel.LocalEndpoint
}

func SSHKeyboardAuth() ssh.AuthMethod {

	return ssh.KeyboardInteractive(func(_, _ string, questions []string, _ []bool) ([]string, error) {

		if len(questions) == 0 {
			log.Printf("prompt ssh auth: keyboard interactive: no questions")
			return nil, nil
		}

		label := getLabel(questions)
		prompt := promptui.Prompt{Label: label}
		if strings.Contains(strings.ToLower(label), "password") {
			prompt.Mask = '*'
		}

		answer, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		return []string{answer}, nil
	})
}

func getLabel(questions []string) string {

	// labels do not handle multilines, split on new line
	var questionLines []string
	for _, question := range questions {
		questionLines = append(questionLines, strings.Split(question, "\n")...)
	}

	// if there are multiple questions (or lines) print all except last line
	for i := 0; i < len(questionLines)-1; i++ {
		log.Println(questionLines[i])
	}

	// use last line as prompt and remove ':', promptui prints this by default provides
	lastQuestion := questionLines[len(questionLines)-1]
	return strings.TrimSuffix(strings.TrimSpace(lastQuestion), ":")
}
```