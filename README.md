# httpGo<sup>Test</sup>
 **Learn Testing, without coding!!**
 
Open source API testing and automation framework written in Go language. 

## Overview
httpGo<sup>Test</sup> enables to write the REST API test cases and automate them with minimal effort. Even non-programmers can also write and manage the test cases wit ease. This framework is written completely in Go language and can be extensible other than HTTP as per use case. For running the tool no external dependencies are required which is super simple.


## Test Project Structure
Test project will consolidate the test groups that contains the list of test cases. Each test case will contain the variables and list of api test steps that needs to be performed.
In httpGo test case is nothing but a simple xml file, if you can understand xml you are ready to write test cases in minutes. 

![N|Solid](https://github.com/git-naren/httpGo/blob/master/images/test_project_structure.png)

All the projects will be located in **~/data/** path of httoGo. Creating adding a new project is simply creating new folder in project path, follwed by group and testcase

**Sample Test Case XML**  [TAG **&lt;testcase&gt;**]
Following show the sample test case xml to fetch the user details via HTTP API with user name and validate the response status
```sh
  <?xml version="1.0" encoding="UTF-8"?>
  <testcase name="TC0001_Get_api_sample" desc="Test case to show the sample GET http api test format">
    <variable name="user.id">testUser123</variable>
    <variable name="app.source">API</variable>
    <teststep type="HTTP" desc="GET user details based on userId">		
      <method>GET</method>
      <url>http://localhost:9090/app/fetch/profile/${user.id}</url>	
      <header name="source">${app.source}</header>
      <header name="user-agent">httpGo Server v1.0.0</header>
      <extract variable="api.res.status" from="http.body" type="json" value="//status">status parameter not found in the response body</extract>	
      <expect name="http.code" check="eq" value="200">api response status code is other than 200 OK</expect>	
    </teststep>
  </testcase>
```
  
### Test Steps [TAG **&lt;teststep&gt;** Attributes *type, desc, runalways*]
Test steps are the core functionality that has to be achieved by the test case. Currently httpGo supports the following test step types

  - **HTTP STEP** for achieving the http protocol based API transactions Ex:- Sending API request
  - **SSH STEP** to execute the remote commands that may be required during test case execution Ex:- Starting or stopping the servers before running the test
  - **SFTP STEP** for processing the upload and download file operations between remote machines and httpGo server. Ex:- Downloading the log files
  - **TIMER STEP** to keep some delay between the test steps execution when needed. Ex:- after sending api request if acoount sync takes some time to reflect in main database
  
 About the detailed tags and supported options about each type of step please refer [DempTestProject Samples!](https://github.com/git-naren/httpGo/tree/master/data/DemoTestProject) 


### Variables [TAG **&lt;variable&gt;** Attribute *name*]
Variables are the parameters that can be used during test case execution either to replace the data or to represent the servers. Variables are classified into four types Environmental, Global, Local type and Fixed type.

  - Environmental variables are shared across all projects defined in **env.properties** file in root data path. Ex- server nodes config.
  - Global variables are shared specific to project that are defined in **project.properties** located in project path.
  - Local variables are limited to the test case and defined in the test case xml it self.
  - Fixed variables are prepared internally to store the response information based on step type and scope is limited to test step
      <br>HTTP step results "http.status", "http.code", "http.header" and "http.body" variables
      <br>SSH step results "ssh.status", and "ssh.output" variables
      <br>SFTP step results "sftp.status" variable

All varibales needs to be use in the notaion **${varibale.name}** in test case files for replacement or reference
   
   
### Variable Extractions [TAG **&lt;extract&gt;** Attributes *variable, from, type, value*]
Variable extractions are to create a new variable from the test step response data. The scope of these new variable extracted is limited to the test case for using in next steps or for validations. Currently httpGo can extract the variables from JSON/XML/HTML from API responses using jpath/xpat notation or by string regular expression.


### Expects (or) Validations  [TAG **&lt;expect&gt;** Attributes *name, check, value*]
Expects are the validation that needs to be performed after executing the test step. for example http status code 200 for the api response. httpGo supports most of the validation types in simple keywords that are user friendly. Below are the list of validations that can be performed on a variable

|"empty", "em", "0", "null", "not_empty", "not empty", "!em", "!0", "not null", "not_null", "!null", "!empty", "eq", "equal", "equals", "=", "eq_case", "eq case", "equals case", "==", "equal case", "equal_case", "equals_case","not_eq", "not eq", "not equals", "!=", "not equal", "not_equal", "not_equals", "!eq", "!equal", "!equals", "contain", "contains", "found", "not_contain", "not_contains", "not contain", "not contains", "not_found", "not found", "!found", "!contain", "!contains“ "match", "re“ "not_match", "not match", "!match", "!re","lt", "less_than", "less than", "<“ "lt_eq", "less_than_eq", "less_than_equal", "less_than_equals", "lt eq", "less than eq", "less than equal", "less than equals", "<=", "gt", "greater_than", "greater than", ">“"gt_eq", "greater_than_eq", "greater_than_equal", "greater_than_equals", "gt eq", "greater than eq", "greater than equal", "greater than equals", ">="


## Installation and Usage
No complex installation is needed for httpGo, you just need to extract the archive and start using it by executing **httpGo.exe** from command line.
Package can be downloaded from [Release Link](https://github.com/git-naren/httpGo/releases/tag/v.1.0.0) OR build the matching OS binary from source code.

![N|Solid](https://github.com/git-naren/httpGo/blob/master/images/http_root_files.png)

        Usage:
                httpGo.exe --list -p *
                httpGo.exe --list -p projectName -g *
                httpGo.exe --list -p projectName -g groupName -t *
                                [OR]
                httpGo.exe --run -p projectName -g *
                httpGo.exe --run -p projectName -g groupName -t *
                httpGo.exe --run -p projectName -g groupName -t testCaseName

                httpGo.exe --help
                httpGo.exe --version
        
        --list option is used for viewing the projects, groups or test cases
        --run option is used for running the project, specific group, all test cases in a group or single test case.
        -p name of the test project
        -g name of the group
        -t name of the test case
        
 Once the test operation is performed it will generate the test report in $project_path/__reports__/last_test_run_Id/ path. Also the status report email alerts can be enabled in environmental variables.

 
 ### Test Case Automation
 Once the test projects are created with test cases, test execution can be triggered from the continuous integration tools CI/Jenkins/DevOps when ever build is changed via command mode. 
 
 
 ### TODO's
  - Parallel execution of test cases
  - Performance optimization and memory issues
  - Enhanced reporting format
  - Web UI for test project management
  - Support Import & Export of test cases
  - Support of Security & Performance Testing
  
  
  
  **Happy Testing, without coding!!**
 
