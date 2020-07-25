package main

import (
    "strings"
    "regexp"
    "github.com/antchfx/xmlquery"
    "github.com/antchfx/jsonquery"
    "github.com/antchfx/htmlquery"
)

func jsonExtraction(config Config, varName string, content string, jpath string) bool {
	doc, err := jsonquery.Parse(strings.NewReader(content))
	if err != nil {
		LogError.Println("Json parsing error", err)
		return false
	}

	if node := jsonquery.FindOne(doc, jpath); node != nil {		
		config[varName] = node.InnerText()

		LogInfo.Println("Extracted parameter " + varName + " = " + config[varName])
		return true
	}
	return false
}

func xmlExtraction(config Config, varName string, content string, xpath string) bool {
	doc, err := xmlquery.Parse(strings.NewReader(content))
	if err != nil {
		LogError.Println("Xml parsing error", err)
		return false
	}

	if node := xmlquery.FindOne(doc, xpath); node != nil {
		config[varName] = node.InnerText()

		LogInfo.Println("Extracted parameter " + varName + " = " + config[varName])
		return true
	}
	return false
}

func htmlExtraction(config Config, varName string, content string, xpath string) bool {
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		LogError.Println("Html parsing error", err) 
		return false
	}
	
	if node := htmlquery.FindOne(doc, xpath); node != nil {
		config[varName] = htmlquery.InnerText(node)

		LogInfo.Println("Extracted parameter " + varName + " = " + config[varName])
		return true
	}
	return false
}

func textExtraction(config Config, varName string, content string, expression string) bool {
	re := regexp.MustCompile(expression)
	if val := re.FindString(content); val != "" {
		config[varName] = val

		LogInfo.Println("Extracted parameter " + varName + " = " + val)
		return true
	}
	return false
}

func doVariableExtraction(config Config, extract Extract) (bool, string)  {
	LogInfo.Println("Enter doVariableExtraction ", extract)

	status := true
	msg := ""

	if pVal, found := config[extract.From]; found {	
		switch extract.Type { 
			case "xml": 
				 status = xmlExtraction(config, extract.Name, pVal, extract.Value)
			case "json": 
				 status = jsonExtraction(config, extract.Name, pVal, extract.Value)
			case "html": 
				 status = htmlExtraction(config, extract.Name, pVal, extract.Value)
			case "text": 
				 status = textExtraction(config, extract.Name, pVal, extract.Value)
			case "param": 
				 config[extract.Name] = pVal
			default:
				status = false
				msg = "un supported extraction type "+extract.Type	
		}
			    
	} else {
		status = false
		msg = "parameter not found in the path"
	}

	LogInfo.Println("Extraction result", status, msg)

	return status, msg
}
