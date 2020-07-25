package main

import ( 
    "strings"   
    "io/ioutil"   
    "net/http"
    "time"
    "strconv"
)

func httpStepExecutor(cm Config, step *TestStep) { 
    TID := cm[TEST_RUN_ID]
    LogInfo.Println(TID, "Enter httpStepExecutor "+ step.Description)
    
    httpURL := step.URL
    httpBody := strings.TrimSpace(step.Body)

    isBodyPresent := false
    if len(httpBody) > 0 {
	isBodyPresent = true
	if strings.HasPrefix(httpBody, "@file=") {
		sParts := strings.Split(httpBody, "=")
		if basepath, ok := cm[GROUP_PATH]; ok {
		    contentFile := basepath + CONTENT_LOCATION + strings.TrimSpace(sParts[1])
		    LogInfo.Println(TID, "Fetch http body from content file: "+ contentFile)
		    httpBody = getFileContentString(contentFile)
		}		
	}	
    }

    httpURL = replace_vars_data(cm, httpURL)
        
    var httpReq *http.Request
    if isBodyPresent {
	httpBody = replace_vars_data(cm, httpBody)
	httpReq, _ = http.NewRequest(step.Method, httpURL, strings.NewReader(httpBody))
    }	else {
	httpReq, _ = http.NewRequest(step.Method, httpURL, nil)
    }
    
    query := httpReq.URL.Query()
    for i := 0; i < len(step.Parameters); i++ {	
	query.Add( step.Parameters[i].Name, replace_vars_data(cm, step.Parameters[i].Value) ) 
    }

    for j := 0; j < len(step.Headers); j++ {	
	httpReq.Header.Add( step.Headers[j].Name, replace_vars_data(cm, step.Headers[j].Value) ) 
    }
    
    LogInfo.Println(TID, "HTTP REQUEST: ", httpReq, "\nBODY:", httpBody)

    step.LogData = "REQUEST: "+httpReq.URL.String()+"\n"+httpBody

    httpClient := http.Client {
	    Timeout: time.Duration(30 * time.Second),
    }
    httpRes, err := httpClient.Do(httpReq)

    if err != nil {
	LogError.Println(TID, "Failed to send http request", err)
	step.Status = false
	return
    }

    httpResData, _ := ioutil.ReadAll( httpRes.Body )
    httpRes.Body.Close()
	
    LogInfo.Println(TID, "HTTP RESPONSE: ", httpRes, "\nBODY:", httpResData)

    cm[PARAM_HTTP_STATUS] = httpRes.Status
    cm[PARAM_HTTP_CODE]   = strconv.Itoa(httpRes.StatusCode)
    cm[PARAM_HTTP_BODY]   = string(httpResData)
    
    step.Status = true

    for i := 0; i < len(step.Extracts); i++ {	
	if strings.EqualFold(step.Extracts[i].From, PARAM_HTTP_HEAD) {		
		cm[step.Extracts[i].Name] = httpRes.Header.Get(step.Extracts[i].Value)
	} else {
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
				LogWarn.Println(TID, "Failed to validate prameter "+step.Expects[j].Name, err)
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
    
    LogInfo.Println(TID, "httpStepExecutor result", step.Status)

    step.LogData = step.LogData + "RESPONSE: " + cm[PARAM_HTTP_STATUS] + "\n status code: " + cm[PARAM_HTTP_CODE] + "\n" + cm[PARAM_HTTP_BODY]

    delete(cm, PARAM_HTTP_STATUS)
    delete(cm, PARAM_HTTP_CODE)
    delete(cm, PARAM_HTTP_BODY)
}