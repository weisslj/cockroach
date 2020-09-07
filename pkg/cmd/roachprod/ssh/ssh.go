// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.

package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/weisslj/cockroach/pkg/cmd/roachprod/config"
	"github.com/weisslj/cockroach/pkg/util/syncutil"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

var knownHosts ssh.HostKeyCallback
var knownHostsOnce sync.Once

// InsecureIgnoreHostKey TODO(peter): document
var InsecureIgnoreHostKey bool

func getKnownHosts() ssh.HostKeyCallback {
	knownHostsOnce.Do(func() {
		var err error
		if InsecureIgnoreHostKey {
			knownHosts = ssh.InsecureIgnoreHostKey()
		} else {
			knownHosts, err = knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	return knownHosts
}

func getSSHAgentSigners() []ssh.Signer {
	const authSockEnv = "SSH_AUTH_SOCK"
	agentSocket := os.Getenv(authSockEnv)
	if agentSocket == "" {
		return nil
	}
	sock, err := net.Dial("unix", agentSocket)
	if err != nil {
		log.Printf("SSH_AUTH_SOCK set but unable to connect to agent: %s", err)
		return nil
	}
	agent := agent.NewClient(sock)
	signers, err := agent.Signers()
	if err != nil {
		log.Printf("unable to retrieve keys from agent: %s", err)
		return nil
	}
	return signers
}

func getSSHKeySigner(path string, haveAgent bool) ssh.Signer {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("unable to read SSH key %q: %s", path, err)
		}
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		if strings.Contains(err.Error(), "cannot decode encrypted private key") {
			if !haveAgent {
				log.Printf(
					"skipping encrypted SSH key %q; if necessary, add the key to your SSH agent", path)
			}
		} else {
			log.Printf("unable to parse SSH key %q: %s", path, err)
		}
		return nil
	}
	return signer
}

func getDefaultSSHKeySigners(haveAgent bool) []ssh.Signer {
	var signers []ssh.Signer
	for _, name := range []string{"id_rsa", "google_compute_engine"} {
		s := getSSHKeySigner(filepath.Join(config.OSUser.HomeDir, ".ssh", name), haveAgent)
		if s != nil {
			signers = append(signers, s)
		}
	}
	return signers
}

func newSSHClient(user, host string) (*ssh.Client, net.Conn, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(sshState.signers...)},
		HostKeyCallback: getKnownHosts(),
	}
	config.SetDefaults()

	addr := fmt.Sprintf("%s:22", host)
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return nil, nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, nil, err
	}
	return ssh.NewClient(c, chans, reqs), conn, nil
}

type sshClient struct {
	syncutil.Mutex
	*ssh.Client
}

var sshState = struct {
	signers     []ssh.Signer
	signersInit sync.Once

	clients  map[string]*sshClient
	clientMu syncutil.Mutex
}{
	clients: map[string]*sshClient{},
}

// NewSSHSession TODO(peter): document
func NewSSHSession(user, host string) (*ssh.Session, error) {
	if host == "127.0.0.1" || host == "localhost" {
		return nil, errors.New("unable to ssh to localhost; file a bug")
	}

	sshState.clientMu.Lock()
	target := fmt.Sprintf("%s@%s", user, host)
	client := sshState.clients[target]
	if client == nil {
		client = &sshClient{}
		sshState.clients[target] = client
	}
	sshState.clientMu.Unlock()

	sshState.signersInit.Do(func() {
		sshState.signers = append(sshState.signers, getSSHAgentSigners()...)
		haveAgentSigner := len(sshState.signers) > 0
		sshState.signers = append(sshState.signers, getDefaultSSHKeySigners(haveAgentSigner)...)
	})

	client.Lock()
	defer client.Unlock()
	if client.Client == nil {
		var err error
		client.Client, _, err = newSSHClient(user, host)
		if err != nil {
			return nil, err
		}
	}
	return client.NewSession()
}

// IsSigKill TODO(peter): document
func IsSigKill(err error) bool {
	switch t := err.(type) {
	case *ssh.ExitError:
		return t.Signal() == string(ssh.SIGKILL)
	}
	return false
}

// ProgressWriter TODO(peter): document
type ProgressWriter struct {
	Writer   io.Writer
	Done     int64
	Total    int64
	Progress func(float64)
}

func (p *ProgressWriter) Write(b []byte) (int, error) {
	n, err := p.Writer.Write(b)
	if err == nil {
		p.Done += int64(n)
		p.Progress(float64(p.Done) / float64(p.Total))
	}
	return n, err
}
