package main

import (
    "fmt"
    "time"
    "strings"  
)

func processTestRequest(project string, group string, tcName string) (response TestRunResponse) {
	LogInfo.Println("Enter processTestRequest for project: " + project + ", group: "+ group + ", testcase: "+ tcName)

	response.Status =  SUCCESS
	response.Message =  SUCCESS_MSG
	response.Scheduled = false
    
	projectRoot := PROJECT_ROOT_PATH + project

	if (len(strings.TrimSpace(project)) > 0 && checkPathExist(projectRoot)) {	
		pCfgMap := Config{} 
 		eConfig, err := ReadConfig(PROJECT_ROOT_PATH + ENV_PROPERTIES_FILE)
 		if err != nil {
			LogWarn.Println("Failed to read environment properties file", err)
 		}
		for k, v := range eConfig {         
			pCfgMap[k] = v 		
		} 
		config, err := ReadConfig(projectRoot + "/" + PROJECT_PROPERTIES_FILE)
 		if err != nil {
			LogWarn.Println("Failed to read project properties file", err)
 		}
		for k, v := range config {         
			pCfgMap[k] = v 		
		} 

		identifier := uniuriNew() 		
		response.Identifier = identifier
		pCfgMap[TEST_RUN_ID] = identifier

		var testProject TestProject
		testProject.Name = project
		testProject.Path = projectRoot
		testProject.Identifier = identifier
		testProject.ReportPath = projectRoot + REPORTS_LOCATION + identifier
		
		LogInfo.Println("Process test request identifier is "+identifier)

		if(len(strings.TrimSpace(group)) != 0 && group == "*") {	
			prepareTestGroups(&testProject, projectRoot)	
			
			executeTestProject(pCfgMap, &testProject)

			response.Scheduled = true
			LogInfo.Println(identifier, "Test project execution is scheduled", response)			
			return response
		} else if(len(strings.TrimSpace(group)) != 0 && group != "*") {
			// SINGLE GROUP
			projectRoot = projectRoot + "/" + group
			if checkPathExist(projectRoot) {
				pCfgMap[GROUP_PATH] = projectRoot

				var testGroup TestGroup
				testGroup.Name = group
				testGroup.Path = projectRoot
				
				if(len(strings.TrimSpace(tcName)) != 0 && tcName == "*") {
					prepareTestCases(&testGroup, projectRoot)
					testProject.TestGroups = append(testProject.TestGroups, testGroup)		
					
					executeTestProject(pCfgMap, &testProject)

					response.Scheduled = true
					LogInfo.Println(identifier, "Test groups execution is scheduled", response)					
					return response
				} else if(len(strings.TrimSpace(tcName)) != 0 && tcName != "*") {
					// SINGLE TESTCASE
					projectRoot = projectRoot + "/" + tcName + TC_FILE_EXT
					if checkPathExist(projectRoot) {						
						testCase := parseTestCaseFile(projectRoot)
						testCase.Path = projectRoot
						testGroup.TestCases = append(testGroup.TestCases, testCase)
						testProject.TestGroups = append(testProject.TestGroups, testGroup)						
						executeTestProject(pCfgMap, &testProject)

						LogInfo.Println(identifier, "Test case executed", response)
						return response
					}
				}			
			}
		}	
	}    
    
	LogError.Println("Failed to process test request project or group or testcase not found")

	response.Status =  ERROR
	response.Message =  "operation failed, project (or) group (or) testcase not found"    

	LogInfo.Println("ProcessTestRequest response", response)

	return response	
}

func prepareTestGroups(testProject *TestProject, path string) {  
	LogInfo.Println("Enter prepareTestGroups "+path)

	var groupFolders []string	
	status := getFoldersInPath(&groupFolders, path, false)
	if !status {
		LogWarn.Println("Get test groups failed from "+path)
	}

	for _, folder := range groupFolders {
		groupPath := path+ "/" + folder
		LogInfo.Println("Loading test group "+groupPath)
		var testGroup TestGroup
		testGroup.Name = folder
		testGroup.Path = groupPath				
		prepareTestCases(&testGroup, groupPath)
		testProject.TestGroups = append(testProject.TestGroups, testGroup)
	}
}

func prepareTestCases(testGroup *TestGroup, path string) {  
	LogInfo.Println("Enter prepareTestCases "+path)

	var testFiles []string
	status := getFilesInPath(&testFiles, path, true, TC_FILE_EXT, false)
	if !status {
		LogWarn.Println("Get test case files failed from "+path)
	}
	for _, file := range testFiles {
		LogInfo.Println("Loading test case file "+file)
		var testCase TestCase
		testCase = parseTestCaseFile(file)
		testCase.Path = file
		testGroup.TestCases = append(testGroup.TestCases, testCase)
	}
}

func executeTestProject(pcm Config, testProject *TestProject) {  
	TID := pcm[TEST_RUN_ID]
	LogInfo.Println(TID, "Enter executeTestProject "+testProject.Name)

	//Create default project dirs	
	createFolders(testProject.ReportPath)

	sTime := time.Now()
	for i := 0; i < len(testProject.TestGroups); i++ {
		executeTestGroup(pcm, &testProject.TestGroups[i])
		
		for _, testCase := range testProject.TestGroups[i].TestCases {
			testProject.TestGroups[i].TotalCount += 1
			if testCase.Status {
				testProject.TestGroups[i].PassCount += 1
			} else {
				testProject.TestGroups[i].FailCount += 1
			}
		}
	}
	testProject.ExecutionTime = time.Since(sTime)

	//Generate report
	generateTestReport(testProject)

	//Send email report
	sendEmailReport(pcm, testProject)
}

func executeTestGroup(pcm Config, testGroup *TestGroup) {  
	TID := pcm[TEST_RUN_ID]
	LogInfo.Println(TID, "Enter executeTestGroup "+testGroup.Name)
	fmt.Println("  Running Test Group ("+testGroup.Name+")")

	sTime := time.Now()
	for i := 0; i < len(testGroup.TestCases); i++ {
		executeTestCase(pcm, &testGroup.TestCases[i])
	}
	testGroup.ExecutionTime = time.Since(sTime)
}

func executeTestCase(pcm Config, testCase *TestCase) { 
	TID := pcm[TEST_RUN_ID]
	LogInfo.Println(TID, "Enter executeTestCase", testCase.Name, testCase)
	fmt.Print("\tExecuting Test Case ["+testCase.Name+"] ")

	sTime := time.Now()

	tCfgMap := Config{} 
	for k, v := range pcm {         
		tCfgMap[k] = v 
	}     
	for i := 0; i < len(testCase.Variables); i++ {
		tCfgMap[testCase.Variables[i].Name] = testCase.Variables[i].Value
	}

	fmt.Print(">")

	tcStatus := true
	for i := 0; i < len(testCase.TestSteps); i++ {	
		//Check for must execute case even failed
		if (!tcStatus && !testCase.TestSteps[i].RunAlways) {
			LogInfo.Println(TID, "No need to run the step", testCase.TestSteps[i].Description)
			continue;
		}

		LogInfo.Println(TID, "Execute "+testCase.TestSteps[i].Type+" Step", testCase.TestSteps[i])
		fmt.Print(".")
	
		stepType := testCase.TestSteps[i].Type	  	
		if stepType == TEST_STEP_TYPE_HTTP {
			httpStepExecutor(tCfgMap, &testCase.TestSteps[i])
		} else if stepType == TEST_STEP_TYPE_SSH {
			sshStepExecutor(tCfgMap, &testCase.TestSteps[i])
		} else if stepType == TEST_STEP_TYPE_SFTP {
			sftpStepExecutor(tCfgMap, &testCase.TestSteps[i])
		} else if stepType == TEST_STEP_TYPE_TIMER {
			timerStepExecutor(tCfgMap, &testCase.TestSteps[i])
		} else {
			LogWarn.Println(TID, "Unsupported step type "+stepType)
		}	

		LogInfo.Println(TID, testCase.TestSteps[i].Type, "Step Status:", conditionalString(testCase.TestSteps[i].Status, "success", "failed"))

		if (tcStatus && !testCase.TestSteps[i].Status) {	
			tcStatus = false
			//break;
			LogWarn.Println(TID, "Test step execution failed, check for any mandatory steps to be run after failure")
		}
	}
	
	testCase.Status = tcStatus
	testCase.Result = conditionalString(tcStatus, "PASS", "FAIL")
	testCase.ExecutionTime = time.Since(sTime)

	LogInfo.Println(TID, testCase.Name, "Test Case Result: ", testCase.Result, testCase.ExecutionTime)
	fmt.Println(testCase.Result, testCase.ExecutionTime)
}