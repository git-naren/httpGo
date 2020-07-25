package main

import "time"
import "encoding/xml"

type Variable struct {
	XMLName	xml.Name   `xml:"variable"`
	Name	string   `xml:"name,attr"`
	Value	string   `xml:",chardata"`
}

type Expect struct {
	XMLName   xml.Name   `xml:"expect"`
	Name	string   `xml:"name,attr"`
	Check   string   `xml:"check,attr"`
	Value    string   `xml:"value,attr"`
	Optional  bool   `xml:"optional,attr"`
	ErrMsg  string   `xml:",chardata"`
}

type Extract struct {
	XMLName   xml.Name   `xml:"extract"`
	Name	string   `xml:"variable,attr"`
	Type    string   `xml:"type,attr"`
	From    string   `xml:"from,attr"`
	Value    string   `xml:"value,attr"`
	Mandatory  bool   `xml:"mandatory,attr"`
	ErrMsg  string   `xml:",chardata"`
}

type Header struct {
	XMLName   xml.Name   `xml:"header"`
	Name   string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

type QueryParam struct {
	XMLName   xml.Name   `xml:"param"`
	Name   string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

type SSHCmd struct {
	XMLName   xml.Name   `xml:"command"`
	Command  	string   `xml:",chardata"`
	Mandatory    bool   `xml:"mandatory,attr"`
}

type FTPFile struct {
	XMLName   xml.Name   `xml:"ftp"`
	Type		string   `xml:"type,attr"`
	Source		string   `xml:"source"`
	Destination	string   `xml:"destination"`
	Mandatory      bool   `xml:"mandatory,attr"`
}

type TestStep struct {
	XMLName   xml.Name   `xml:"teststep"`    
	Type			string   `xml:"type,attr"`
	Description	string   `xml:"desc,attr"`
	RunAlways	bool	   `xml:"runalways,attr"`
	
	//HTTP
	Method	string   `xml:"method"`
	URL	string   `xml:"url"`
	Parameters  []QueryParam   `xml:"param"`
	Headers       []Header   `xml:"header"`
	Body		string   `xml:"body"`
    
	//SSH & SFTP    
	HostNode	string   `xml:"nodeid"`
	Host		string   `xml:"host"`
	Port		string   `xml:"port"`
	Username	string   `xml:"username"`
	Password	string   `xml:"password"`
	Keyfile		string   `xml:"keyfile"`
	SSHCommands  []SSHCmd   `xml:"command"`
	FTPFiles      []FTPFile   `xml:"ftp"`

	//TIMER
	WaitTimeInSec	int32  `xml:"seconds,attr"`

	//Response
	Expects	       []Expect  `xml:"expect"`
	Extracts       []Extract  `xml:"extract"`

	//Report
	Status	   bool
	Message  string
	LogData   string
}

type TestCase struct {
	XMLName   xml.Name   `xml:"testcase"`
	Name         string   `xml:"name,attr"`
	Tags           string   `xml:"tag,attr"`
	Description string   `xml:"desc,attr"`
	Variables    []Variable   `xml:"variable"`
	TestSteps   []TestStep   `xml:"teststep"`
    
	Path      string
	Status   bool
	Result   string 
	ExecutionTime time.Duration
}

type TestGroup struct {
	Name string
	Path   string
	TestCases []TestCase		
	ExecutionTime time.Duration
	TotalCount int
	PassCount int
	FailCount int
}

type TestProject struct {
	Name string
	Path   string
	TestGroups []TestGroup	
	Identifier string
	ExecutionTime time.Duration
	ReportPath  string
	ReportTime string
}

type TestRunResponse struct {
	Status        string `json:"status"`
	Message    string `json:"message"`
	Identifier    string `json:"identifier"`
	Scheduled  bool `json:"scheduled"`
}

type ReportFolder struct {
	Name string
	Path  string
	ReportTime time.Time	
}