package main

import (
    "os"
    "sort"
    "strings"
    "regexp"
    "path/filepath"    
    "io/ioutil"
    "encoding/xml"
    "crypto/rand"
)

/*
* File handlers to list or update
*
*/
func checkPathExist(fileOrFolder string) bool {
	if _, err := os.Stat(fileOrFolder); os.IsNotExist(err) {
		LogError.Println("File or folder "+fileOrFolder+" not exist")
		return false
	}
	return true
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func createFolders(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func getFoldersInPath(outFolders *[]string, path string, fullpath bool) bool {
	if checkPathExist(path) {	
		files, err := ioutil.ReadDir(path)
		if err != nil {
			LogError.Println("Failed to read directory", err)
			return false
		}
		for _, fInfo := range files {
			if (fInfo.IsDir() && fInfo.Name() != "__reports__") {
				if(fullpath) {
					*outFolders = append(*outFolders, filepath.Join(path, fInfo.Name()))   					
				} else {
					*outFolders = append(*outFolders, fInfo.Name())   
				}
			}
		}		
		return true
	}
	return false
}

func getFilesInPath(outFiles *[]string, path string, fullpath bool, extension string, removeExtension bool) bool {
	if checkPathExist(path) {	
		files, err := ioutil.ReadDir(path)
		if err != nil {
			LogError.Println("Failed to read directory", err)
			return false
		}
		for _, fInfo := range files {
			if !fInfo.IsDir() {
				r, err := regexp.MatchString(extension, fInfo.Name())
				if err == nil && r {			
					if(fullpath) {
						*outFiles = append(*outFiles, filepath.Join(path, fInfo.Name()))   					
					} else {
						if(removeExtension) {
							*outFiles = append(*outFiles, fileNameWithoutExtension(fInfo.Name()))   
						} else {
							*outFiles = append(*outFiles, fInfo.Name())   
						}
					}
				}  				
			}
		}		
		return true
	}
	return false
}

func getFoldersInPathOrderByTime(outFolders *[]ReportFolder, path string) bool {
	if checkPathExist(path) {	
		files, err := ioutil.ReadDir(path)
		if err != nil {
			LogError.Println("Failed to read directory", err)
			return false
		}

		sort.Slice(files, func(i,j int) bool{
		    return files[i].ModTime().After(files[j].ModTime())
		})

		for _, fInfo := range files {
			if fInfo.IsDir() {
				folder := ReportFolder{
					Name: fInfo.Name(),
					Path: filepath.Join(path, fInfo.Name()),
					ReportTime: fInfo.ModTime(),
				}
				*outFolders = append(*outFolders, folder)   
			}
		}		
		return true
	}
	return false
}

func getFileContentString(fileName string) string {    
    textContent := ""
    content, err := ioutil.ReadFile(fileName)  
    if err != nil {
        LogError.Println("Failed to read file", err)
	return textContent
    }
    textContent = string(content)
    return textContent
}

func parseTestCaseFile(fileName string) (testCase TestCase) {    
    xmlFile, err := os.Open(fileName)    
    if err != nil {
	LogError.Println("Failed to open test case file", err)
	return
    }
    defer xmlFile.Close()

    byteValue, _ := ioutil.ReadAll(xmlFile)   
    xml.Unmarshal(byteValue, &testCase)  
    return
}


/*
* Unique Id generator with different lenght of inputs
*
*/
const (
	StdLen = 16
	UUIDLen = 20
)
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func uniuriNew() string {
	return uniuriNewLenChars(StdLen, StdChars)
}

func uniuriNewLen(length int) string {
	return uniuriNewLenChars(length, StdChars)
}

func uniuriNewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		LogError.Println("uniuri: wrong charset length for NewLenChars")
		return ""
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			LogError.Println("uniuri: error reading random bytes", err)
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}


/*
* Conditional checks
* returns val1 if condition, otherwise val2
*	if conditionalString(true, "ok", "not ok") != "ok" {
*		t.Fatal("Error string")
*	}
*/
func conditionalDefaultString(value, defaultValue string) string {
	if len(strings.TrimSpace(value)) > 0 {	
		return value
	}
	return defaultValue
}

func conditionalString(condition bool, val1, val2 string) string {
	if condition {
		return val1
	}
	return val2
}

func conditionalInterface(condition bool, val1, val2 interface{}) interface{} {
	if condition {
		return val1
	}
	return val2
}

func conditionalInt(condition bool, val1, val2 int) int {
	if condition {
		return val1
	}
	return val2
}


/**
*  Regex expression to replace all the parameters in the config map and return the string
*/
func replace_vars_data(config Config, data string) string {
	outData := data
	re := regexp.MustCompile(`\$\{(.*?)\}`)	
	matchVariables:= re.FindAllStringSubmatch(outData, -1) 
	for _, name := range matchVariables {
		//fmt.Println(i[0])  ${myvar} //fmt.Println(i[1])  myvar
		if val, found := config[name[1]]; found {			
			outData = strings.Replace(outData, name[0], val, -1)
		}
	}
	return outData
}