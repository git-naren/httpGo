package main

import (	
	"io"	
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"	
	"path/filepath"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpServerInfo struct {	
	Host		string
	Port		string
	User		string
	Pass		string
	KeyFile     string
}

func (c *SftpServerInfo) Socket() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		LogError.Println("Failed to read public key for host", host, err)
	}
	defer file.Close()
 
	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				//log.Fatalf("error parsing %q: %v", fields[2], err)
				LogError.Println("Failed to parse key file ", err)
			}
			break
		}
	}
 
	if hostKey == nil {
		LogWarn.Println("No hostkey found for %s", host)
	}
 
	return hostKey
}

func ftpFileProcess(s SftpServerInfo, sfName string, dfName string, ftpMode string) (bool, string) {
	LogInfo.Println("Enter ftpFileProcess sfName: "+sfName+" , dfName: "+dfName+", ftpMode: "+ftpMode)
	// get host public key
	//hostKey := getHostKey(s.Host)
 
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Pass),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 15 * time.Second,
	}
 
	// connect
	conn, err := ssh.Dial("tcp", s.Socket(), config)
	if err != nil {
		LogError.Println("Failed to connect sftp server", err)
		return false, "Failed to connect sftp server"
	}
	defer conn.Close()
 
	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		LogError.Println("Failed to create ftp session", err)
		return false, "Failed to create ftp session"
	}
	defer client.Close()
 
	if (ftpMode == "upload") {
		// create destination file
		dstFile, err := client.Create(dfName)
		if err != nil {
			LogError.Println("Failed to create destination file", err)
			return false, "Failed to create destination file"
		}
		defer dstFile.Close()
	 
		// create source file
		srcFile, err := os.Open(sfName)
		if err != nil {
			LogError.Println("Failed to open source file", err)
			return false, "Failed to open source file"
		}
	 
		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			LogError.Println("Failed to copy from source to destination file", err)
			return false, "Failed to copy from source to destination file"
		}
		LogInfo.Println("%d bytes copied", bytes)

		LogInfo.Println("File uploaded successfully to remote server")

		return true, "File uploaded successfully"
	} else {
		// create destination file
		dstFile, err := os.Create(dfName)
		if err != nil {
			LogError.Println("Failed to create destination file", err)
			return false, "Failed to create destination file"
		}
		defer dstFile.Close()
	 
		// open source file
		srcFile, err := client.Open(sfName)
		if err != nil {
			LogError.Println("Failed to open source file", err)
			return false, "Failed to open source file"
		}
	 
		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			LogError.Println("Failed to copy from source to destination file", err)
			return false, "Failed to copy from source to destination file"
		}
		LogInfo.Println("%d bytes copied", bytes)
	 
		// flush in-memory copy
		err = dstFile.Sync()
		if err != nil {
			LogError.Println("Failed to flush data from memory to file", err)
			return false, "Failed to flush data from memory to file"
		}

		LogInfo.Println("File downloaded successfully from remote server")

		return true, "File downloaded successfully"
	}
}

func sftpStepExecutor(cm Config, step *TestStep) { 
	TID := cm[TEST_RUN_ID]
	LogInfo.Println(TID, "Enter sftpStepExecutor "+ step.Description)
	
	serverInfo := SftpServerInfo {
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
	
	LogInfo.Println(TID, "Connecting to sftp server", step.HostNode , serverInfo)

	step.Status = true

	for i := 0; i < len(step.FTPFiles); i++ {	
		status, error := ftpFileProcess(serverInfo, step.FTPFiles[i].Source, step.FTPFiles[i].Destination, step.FTPFiles[i].Type) 
		LogInfo.Println(TID, "SFTP RESPONSE <=", status, step.FTPFiles[i], error)
		cm[PARAM_SFTP_STATUS] += conditionalString(status, "true , ", "false , ")
		if !status {
			LogWarn.Println(TID, "Run sftp command failed", error)
			if step.FTPFiles[i].Mandatory {
				step.Status = false
				step.Message = "Mandatory sftp operation failed"
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

	
	LogInfo.Println(TID, "sftpStepExecutor result", step.Status)

	//step.LogData = cm[PARAM_SSH_OUTPUT]

	delete(cm, PARAM_SFTP_STATUS)
}
