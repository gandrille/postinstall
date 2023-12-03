package functions

import (
	"os/exec"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// SSHFunction structure
type SSHFunction struct {
}

// Infos function
func (f SSHFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "SSH initialization",
		ShortDescription: "Initializes .ssh folder",
		LongDescription:  "Creates a " + keysize() + " rsa public/private key pair if needed",
	}
}

// Run function
func (f SSHFunction) Run() result.Result {
	return execute(f.Infos().Title, createKeyPair)
}

func createKeyPair() result.Result {
	if exists, _ := filesystem.RegularFileExists("~/.ssh/id_rsa"); exists {
		return result.NewUnchanged("SSH key pair already created")
	}

	if err := exec.Command("ssh-keygen", "-b", keysize(), "-t", "rsa", "-f", filesystem.HomeDir()+"/.ssh/id_rsa", "-P", "").Run(); err != nil {
		return result.NewError("ssh key pair creation failed: " + err.Error())
	}

	return result.NewCreated("ssh key pair created")
}

func keysize() string {
	return "4096"
}
