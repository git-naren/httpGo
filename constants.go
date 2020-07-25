package main

const (

	PROJECT_NAME = "httpGo"
	PROJECT_VERSION = "v1.0.0"
 	PROJECT_ROOT_PATH = "./data/"
	PROJECT_PROPERTIES_FILE = "project.properties"
	ENV_PROPERTIES_FILE = "env.properties"
	CONTENT_LOCATION = "/content/"
	REPORTS_LOCATION = "/__reports__/"
	REPORT_FILE_NAME = "/testStatusReport.html"	

	API_BASE_PATH = "/httpGo"
	WEB_UI_DEFAULT_PATH = API_BASE_PATH+"/web/index.html"

	TC_FILE_EXT = ".xml"
	REPORT_FILE_EXT = ".html"
	
	SERVER_DEFAULT_PORT = "9999"

	TEST_RUN_ID = "p_test_run_id"
	GROUP_PATH = "p_group_path"
	
	SMTP_HOST = "smtp.host"
	SMTP_PORT = "smtp.port"
	SMTP_USER = "smtp.user"
	SMTP_PASS = "smtp.pass"
	EMAIL_TO = "report.to.email"
	EMAIL_FROM = "report.from.email"	
	EMAIL_FLAG = "report.email.flag"

	TEST_STEP_TYPE_HTTP = "HTTP"
	TEST_STEP_TYPE_SSH   = "SSH"
	TEST_STEP_TYPE_SFTP = "SFTP"
	TEST_STEP_TYPE_TIMER = "TIMER"

	SUCCESS = "success"
	ERROR   = "error"
	SUCCESS_MSG = "Operation successfull"
	ERROR_MSG = "Operation failed"
	
	PARAM_HTTP_STATUS	= "http.status"
	PARAM_HTTP_CODE		= "http.code"
	PARAM_HTTP_BODY		= "http.body"
	PARAM_HTTP_HEAD		= "http.header"
	PARAM_SSH_STATUS	= "ssh.status"
	PARAM_SSH_OUTPUT	= "ssh.output"
	PARAM_SFTP_STATUS	= "sftp.status"

)

