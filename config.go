package main

import (
 	"bufio"
 	"io"
 	"os"
 	"strings"
 )

 type Config map[string]string

 func ReadConfig(filename string) (Config, error) {
	LogInfo.Println("Read config file "+filename)

 	config := Config{}
 	if len(filename) == 0 {
 		return config, nil
 	}
 	file, err := os.Open(filename)
 	if err != nil {
		LogError.Println("Failed to open config file "+filename, err)
 		return nil, err
 	}
 	defer file.Close()
 	
 	reader := bufio.NewReader(file)
 	
 	for {
 		line, err := reader.ReadString('\n') 		
 		// check if the line has = sign
                // and process the line. Ignore the rest.
 		if equal := strings.Index(line, "="); equal >= 0 {
 			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
 				value := ""
 				if len(line) > equal {
 					value = strings.TrimSpace(line[equal+1:])
 				}
				// assign the config map
 				config[key] = value
 			}
 		}
 		if err == io.EOF {
 			break
 		}
 		if err != nil {
 			return nil, err
 		}
 	}
 	return config, nil
 }
