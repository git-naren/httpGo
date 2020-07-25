package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    LogInfo.Println("-------------------------------------------------- "+PROJECT_NAME+" "+PROJECT_VERSION+" --------------------------------------------------")
    LogInfo.Println("Start executing the application with program arguments", os.Args)

    inputOk := true

    argsWithProg := os.Args

    if(len(os.Args) > 1 && len(os.Args)%2 == 0) {
	argsProgMode := os.Args[1]
	argsWithoutProg := os.Args[2:]
	
	projectName	:= ""
	groupName	:= ""
	testCaseName  := "*"

	for i, s := range argsWithoutProg {    
		if( ((i+1) < len(argsWithoutProg)) && (!strings.HasPrefix(argsWithoutProg[i+1], "-")) ) {
		       if(s == "-p") {
				projectName = argsWithoutProg[i+1]
		       } else if(s == "-g") {
				groupName = argsWithoutProg[i+1]
		       } else if(s == "-t") {
				testCaseName = argsWithoutProg[i+1]
		       } 
	       }
	}

	if (argsProgMode == "--run") {

		if(len(strings.TrimSpace(projectName)) == 0) {
			inputOk = false
		}

		if(len(strings.TrimSpace(groupName)) == 0) {
			inputOk = false
		}

		if(strings.TrimSpace(groupName) == "*" && strings.TrimSpace(testCaseName) != "*") {
			inputOk = false
		}

		if inputOk {			
			LogInfo.Println("Command Run mode,  projectName: " + projectName + ", groupName: " + groupName + ", testCaseName: " + testCaseName)

			response := processTestRequest(projectName, groupName, testCaseName)   

			LogInfo.Println("Test request status", response)
		}
 
	} else if (argsProgMode == "--list") {

		LogInfo.Println("Command List mode,  projectName: " + projectName + ", groupName: " + groupName + ", testCaseName: " + testCaseName)

		if(strings.TrimSpace(projectName) == "*") {
			LogInfo.Println("List all projects from the root path "+PROJECT_ROOT_PATH)

			var projects []string
			status := getFoldersInPath(&projects, PROJECT_ROOT_PATH, false)
			conditionalLogger(status, "Listing projects operation success" , "Listing projects failed")
			if(status) {
				LogInfo.Println("Projects:", projects)
				fmt.Println("================ LIST PROJECTS ================") 
				for i := 0; i < len(projects); i++ {						
					fmt.Println("\t"+projects[i]) 
				}
			} 
		} else if(len(strings.TrimSpace(projectName)) > 0) {

			if(strings.TrimSpace(groupName) == "*") {
				rootPath := PROJECT_ROOT_PATH + projectName	
				LogInfo.Println("List all groups from the project path "+rootPath)

				var groups []string
				status := getFoldersInPath(&groups, rootPath, false)
				conditionalLogger(status, "Listing test groups operation success" , "Listing test groups failed for project "+projectName)
				if(status) {
					LogInfo.Println("Test groups for project "+projectName, groups)
					fmt.Println("================ LIST GROUPS ["+projectName+"]================") 					
					for i := 0; i < len(groups); i++ {						
						fmt.Println("\t"+groups[i]) 
					}
				}
				
			} else if(len(strings.TrimSpace(groupName)) > 0) {

				if(strings.TrimSpace(testCaseName) == "*") {
					rootPath := PROJECT_ROOT_PATH + projectName + "/" + groupName	
					LogInfo.Println("List all test cases from the project group path "+rootPath)

					var tcNames []string
					status := getFilesInPath(&tcNames, rootPath, false, TC_FILE_EXT, true)
					conditionalLogger(status, "Listing test cases operation success" , "Listing test cases failed for project "+projectName+", group "+groupName)
					if(status) {
						LogInfo.Println("Test cases for project "+projectName+", group "+groupName, tcNames)
						fmt.Println("================ LIST TEST CASES ["+projectName+" @ "+groupName+"]================") 	
						for i := 0; i < len(tcNames); i++ {						
							fmt.Println("\t"+tcNames[i]) 
						}
					} 

				} else if(len(strings.TrimSpace(testCaseName)) > 0) {
					//get test case details
					fmt.Println("  Get one testcase details") 
				} else {
					inputOk = false
					fmt.Println("  Get one group details") 
				}
				
			} else {
				inputOk = false
				fmt.Println("  Get one project details") 
			}
			
		} else {
			inputOk = false
		}
	} else if (argsProgMode == "--h" || argsProgMode == "-h" || argsProgMode == "h" || argsProgMode == "--help" || argsProgMode == "-help" || argsProgMode == "help") {
		inputOk = false
	} else if (argsProgMode == "--v" || argsProgMode == "-v" || argsProgMode == "v" || argsProgMode == "--version" || argsProgMode == "-version" || argsProgMode == "version") {
		fmt.Println(PROJECT_NAME + " version " + PROJECT_VERSION)
	} else {
		inputOk = false
	}
	
    } else {
	inputOk = false
    }
    
    if !inputOk {
	LogError.Println("Invalid program options given, please check usage")

	fmt.Println("Usage:")	
	fmt.Println("\t"+ argsWithProg[0] +" --run -p projectName -g *")
	fmt.Println("\t"+ argsWithProg[0] +" --run -p projectName -g groupName -t *")
	fmt.Println("\t"+ argsWithProg[0] +" --run -p projectName -g groupName -t testCaseName")
	fmt.Println("\t\t\t[OR]")
	fmt.Println("\t"+ argsWithProg[0] +" --list -p * ")
	fmt.Println("\t"+ argsWithProg[0] +" --list -p projectName -g * ")
	fmt.Println("\t"+ argsWithProg[0] +" --list -p projectName -g groupName -t *")	
	fmt.Println("\n")
	fmt.Println("\t"+ argsWithProg[0] +" --help")
	fmt.Println("\t"+ argsWithProg[0] +" --version")
	os.Exit(3)
    }

}