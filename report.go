package main

import (	
	"os"
	"time"
	"bytes"
	"strings"
	"net/smtp"
	"text/template"
)

const report_tmpl_html = `<!DOCTYPE html><html><head><title>{{.Name}} Test Report</title><style type="text/css">
*{margin:0;padding:0}#wrapper,body,html{height:100%}header{height:50px}.accordion{background-color:#eee;color:#444;cursor:pointer;padding:10px;width:100%;border:none;text-align:left;outline:0;font-size:13px;font-weight:700;transition:.4s;font-family:Cambria;color:#044b6e}.accordion:hover,.active{background-color:#cceef4}.accordion:before{content:'\002B';color:#777;font-weight:700;float:right;margin-left:5px}.active:before{content:"\2212"}.panel{padding:0 18px;background-color:#fff;max-height:0;overflow:hidden;transition:max-height .2s ease-out}blockquote{font:14px/22px normal helvetica,Cambria;margin-top:10px;margin-bottom:10px;margin-left:50px;padding-left:5px;border-left:1px solid #cceef4}input,select,textarea{font-family:Cambria;color:#044b6e}.custom-font{font-family:Cambria;color:#044b6e}.bg-default,.bg-error,.bg-primary,.bg-secondary,.bg-success,.bg-warning{color:#fff;border-radius:2px;text-shadow:0 1px 1px rgba(0,0,0,.2);font-family:Cambria;font-weight:700;padding-left:5px;padding-right:5px}.bg-success{background:#2ecc71}.bg-error{background:#ca3c3c}.bg-warning{background:#df7514}.bg-primary{background:#044b6e}.bg-secondary{background:#42b8dd}.i-xxsmall{font-size:60%}.i-xsmall{font-size:75%}.i-small{font-size:80%}.i-large{font-size:110%}.i-xlarge{font-size:125%}.bg-lg-primary{background-image:linear-gradient(to bottom,#26759e,#133d5b)}.bg-lg-info{background-image:linear-gradient(to bottom,#1ab0ec,#1a92c2)}.tc_details,.tc_hide:target,.tc_show{display:none}.tc_hide:target+.tc_show,.tc_hide:target~.tc_details{display:block}.testgroup_dummy{margin-top:50px;margin-left:50px;padding:12px;border-bottom:1px solid #cceef4}.testgroup_first{margin-top:50px;margin-left:50px;padding:12px;border:1px solid #cceef4}.testgroup{margin-left:50px;padding:12px;border-left:1px solid #cceef4;border-right:1px solid #cceef4;border-bottom:1px solid #cceef4}.testgroup_head{font-family:Cambria;font-size:14px;font-weight:700;color:#044b6e;padding-bottom:12px;cursor:pointer}.testgroup_head_left{width:70%;float:left;font-size:13px}.testgroup_head_right{width:30%;float:right;display:flex;justify-content:flex-end;font-size:12px;font-weight:400}.tc_details_align{margin-left:1%;height:100%;width:96%}.teststep{margin-left:30px;margin-top:5px;font-size:12px}.teststep_desc{font-weight:700}.teststep_content{padding-left:20px}.teststep_fail_content{padding-left:20px;font-weight:700;color:#ec7063}::-webkit-scrollbar{width:8px;height:10px;border-radius:0}::-webkit-scrollbar-track{background:#cceef4;border-radius:0}::-webkit-scrollbar-thumb{width:12px;height:10px;background:#13486d;border:2px solid #cceef4;border-radius:0}
</style></head><body><div id="wrapper" style="font-family: 'Cambria';color: #044B6E"> <header class="bg-lg-info"><div style="padding-top: 20px; padding-left: 42%; font-weight:bolder; font-family:'Courier New'; font-size:15px; color: #F4F6F7; ">{{.Name}}&nbsp;TEST REPORT</div></header><div style="width:80%; height:100%; float:right;margin-left:10%;margin-right:10%;"> <div style="height: 100%; background: white; border-bottom:2px solid #cceef4;margin-top:40px"><div class="testgroup_dummy"></div>{{range $idx,$group :=.TestGroups}}<div class="testgroup"> <div id="hide_tg{{$idx}}" class="testgroup_head tc_hide"><div class="testgroup_head_left">+&nbsp;{{if gt $group.FailCount 0}}<label class="bg-error">TEST GROUP</label>{{else}}<label class="bg-success">TEST GROUP</label>{{end}}&nbsp;{{$group.Name}}</div><div class="testgroup_head_right">Total&nbsp;<strong>{{$group.TotalCount}}</strong>&nbsp;test cases (<strong>{{$group.PassCount}}PASS/{{$group.FailCount}}FAIL</strong>)&nbsp;Execution time&nbsp;{{$group.ExecutionTime}}</div></div><div id="show_tg{{$idx}}" class="testgroup_head tc_show"><div class="testgroup_head_left">-&nbsp;{{if gt $group.FailCount 0}}<label class="bg-error">TEST GROUP</label>{{else}}<label class="bg-success">TEST GROUP</label>{{end}}&nbsp;{{$group.Name}}</div><div class="testgroup_head_right">Total&nbsp;<strong>{{$group.TotalCount}}</strong>&nbsp;test cases (<strong>{{$group.PassCount}}PASS/{{$group.FailCount}}FAIL</strong>)&nbsp;Execution time&nbsp;{{$group.ExecutionTime}}</div></div><div class="tc_details tc_details_align">{{range $group.TestCases}}<blockquote><div class="accordion">{{if .Status}}<label class="bg-success">TEST CASE</label>{{else}}<label class="bg-error">TEST CASE</label>{{end}}&nbsp;{{.Name}}<span style="font-weight: normal">&nbsp;{{.Description}}</span><span style="float:right; font-weight:normal">Execution time&nbsp;{{.ExecutionTime}}</span></div><div class="panel">{{range .TestSteps}}<div class="teststep">{{if .Status}}<span class="bg-secondary">STEP&nbsp;{{.Type}}</span>{{else}}<span class="bg-error">STEP&nbsp;{{.Type}}</span>{{end}}<span class="teststep_desc">&nbsp;{{.Description}}</span><p class="teststep_content">{{.LogData}}</p>{{if not .Status}}<span class="teststep_fail_content">{{.Message}}</span>{{end}}</div>{{end}}</div></blockquote>{{end}}</div></div>{{end}}</div></div></div>
<script>var i,acc=document.getElementsByClassName("accordion");for(i=0;i<acc.length;i++)acc[i].addEventListener("click",function(){this.classList.toggle("active");var e=this.nextElementSibling;e.style.maxHeight?e.style.maxHeight=null:e.style.maxHeight=e.scrollHeight+"px"});var j,testcases=document.getElementsByClassName("testgroup_head");for(j=0;j<testcases.length;j++)testcases[j].addEventListener("click",function(){window.location.hash=this.id});
</script></body></html>
`

func generateTestReport(testProject *TestProject) bool {	
	LogInfo.Println("Enter generateTestReport for "+testProject.Name)

	report_tmpl := template.New("project report template")
        report_tmpl, err := report_tmpl.Parse(report_tmpl_html)
        if err != nil {
		LogError.Println("Parse report_tmpl failed", err)
                return false
        }

	filePath := testProject.ReportPath + REPORT_FILE_NAME
	outfile, err := os.Create(filePath)
	if err != nil {
		LogError.Println("Failed to create report file failed", err)
		return false
	}

	current_time := time.Now() 	
	testProject.ReportTime = current_time.Format(time.RFC3339Nano)

        err = report_tmpl.Execute(outfile, testProject)
        if err != nil {
                LogError.Println("Execute report_tmpl failed", err)
                return false
        }

	LogInfo.Println("Test report generated successfully", filePath)
	return true
}

func sendEmailReport(pcm Config, testProject *TestProject) bool {	
	
	if pcm[EMAIL_FLAG] != "enable" {
		LogError.Println("Email reporting is not enabled!")
                return false
	}

	//TODO:: template preparation
	email_tmpl := template.New("project email template")
        email_tmpl, err := email_tmpl.Parse(report_tmpl_html)
        if err != nil {
		LogError.Println("Parse email_tmpl failed", err)
                return false
        }

	var outBuffer bytes.Buffer
	err = email_tmpl.Execute(&outBuffer, testProject)
        if err != nil {
                LogError.Println("Execute email_tmpl failed", err)
                return false
        }

	bodyHtml := outBuffer.String()
	subject := ""
	toEmails := strings.Split(pcm[EMAIL_TO], ",")
	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	
	//"From: " + from + "\n" +
	body := "To: " + toEmails[0] + "\r\nSubject: " + subject + "\r\n" + MIME + "\r\n" + bodyHtml
	SMTP := pcm[SMTP_HOST]+":"+pcm[SMTP_PORT]
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", pcm[EMAIL_FROM], pcm[SMTP_PORT], pcm[SMTP_HOST]), pcm[EMAIL_FROM], toEmails, []byte(body)); err != nil {
		LogError.Println("Failed to send email report", err)
		return false
	}

	LogInfo.Println("Test report mailed successfully")
	return true
}