package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"
	"strings"
	"golang.org/x/crypto/ssh"
)

type ServerConnInfo struct {
	Host		string
	Port		string
	User		string
	Pass		string
	KeyFile     string
}

func (c *ServerConnInfo) Socket() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		LogError.Println("Failed to read public key file "+file, err)
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		LogError.Println("Failed to parse public key file "+file, err)
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func generateSession(s ServerConnInfo) (*ssh.Session, ssh.Conn, error) {
	//var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod {
			ssh.Password(s.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		Timeout: 15 * time.Second,
	}

	//IF key mode enabled overwrite config
	if (len(strings.TrimSpace(s.KeyFile)) > 0) {
		publicKey, err := publicKeyFile(s.KeyFile)
		if err != nil {
			return nil, nil, err
		}		
		config.Auth = []ssh.AuthMethod { publicKey, }		
	}

	conn, err := ssh.Dial("tcp", s.Socket(), config)
	if err != nil {
		LogError.Println("Failed to connect ssh server", err)
		return nil, nil, err
	}

	// Each ClientConn can support multiple interactive sessions, represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		LogError.Println("Failed to create ssh session", err)
		return nil, conn, err
	}

	return session, conn, nil
}

func SSHCommand(command string, sci ServerConnInfo) (bool, error, string) {
	LogInfo.Println("Run SSHCommand =>", command)

	session, conn, err := generateSession(sci)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return false, err, ""
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(command)

	session.Close()
	conn.Close()

	if err != nil {
		LogError.Println("Failed to run ssh command", err)
		return false, err, ""
	}
	return true, nil, strings.TrimSuffix(stdoutBuf.String(), "\n")
}

func sshStepExecutor(cm Config, step *TestStep) { 
	TID := cm[TEST_RUN_ID]
	LogInfo.Println(TID, "Enter sshStepExecutor "+ step.Description)

	serverInfo := ServerConnInfo {
		step.Host,
		step.Port,
		step.Username,
		step.Password,
		step.Keyfile,
	}

	if len(step.HostNode) > 0 {		
		serverInfo.Host	= cm[step.HostNode+".ssh.host"]
		serverInfo.Port	= cm[step.HostNode+".ssh.port"]
		serverInfo.User	= cm[step.HostNode+".ssh.user"]
		serverInfo.Pass	= cm[step.HostNode+".ssh.pass"]
		serverInfo.KeyFile = cm[step.HostNode+".ssh.key"]
	} else {
		serverInfo.Host	= replace_vars_data(cm, serverInfo.Host)
		serverInfo.Port	= replace_vars_data(cm, serverInfo.Port)
		serverInfo.User	= replace_vars_data(cm, serverInfo.User)
		serverInfo.Pass	= replace_vars_data(cm, serverInfo.Pass)
		serverInfo.KeyFile = replace_vars_data(cm, serverInfo.KeyFile)
	}
	
	LogInfo.Println(TID, "Connecting to ssh server", step.HostNode , serverInfo)

	step.Status = true

	for i := 0; i < len(step.SSHCommands); i++ {	
		status, exitError, output := SSHCommand(replace_vars_data(cm, step.SSHCommands[i].Command), serverInfo)
		LogInfo.Println(TID, "SSH RESPONSE <=", status, step.SSHCommands[i].Command, output)
		cm[PARAM_SSH_STATUS] += conditionalString(status, "true , ", "false , ")
		cm[PARAM_SSH_OUTPUT] += ("\n"+output)	
		if !status { 			
			LogWarn.Println(TID, "Run ssh command failed", exitError)
			if step.SSHCommands[i].Mandatory {
				step.Status = false
				step.Message = "Mandatory ssh command execution failed"
				LogError.Println(TID, step.Message)	
				break;
			}
		}		
	}	

	if step.Status {
		for i := 0; i < len(step.Extracts); i++ {	
			status, err := doVariableExtraction(cm, step.Extracts[i])
			if !status {
				LogWarn.Println(TID, "Variable extraction failed", err)
				if step.Extracts[i].Mandatory {
					step.Status = false
					step.Message = "Failed to extract mandatory variable "+step.Extracts[i].Name
					LogError.Println(TID, step.Message)						
					break;
				}			
			}			
		}
	}

	if step.Status {
		for j := 0; j < len(step.Expects); j++ {	
			if pNameVal, found := cm[step.Expects[j].Name]; found {
				status, err := doValidation(pNameVal, step.Expects[j].Check, replace_vars_data(cm, step.Expects[j].Value), step.Expects[j].ErrMsg)
				if !status {
					LogWarn.Println(TID, "Failed to validate parameter "+step.Expects[j].Name, err)
					if !step.Expects[j].Optional {
						step.Status = false
						step.Message = "Failed to validate parameter "+step.Expects[j].Name
						LogError.Println(TID, step.Message)						
						break;
					}
				}	    
			} else if !step.Expects[j].Optional {							 				
				step.Status = false
				step.Message = "Failed to validate parameter "+step.Expects[j].Name+", as it's not found"
				LogError.Println(TID, step.Message)	
				break;
			}
		}
	}

	LogInfo.Println(TID, "sshStepExecutor result", step.Status)

	step.LogData = cm[PARAM_SSH_OUTPUT]

	delete(cm, PARAM_SSH_STATUS)
	delete(cm, PARAM_SSH_OUTPUT)	
}
