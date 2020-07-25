package main

import (
    "strings"
    "strconv"
    "regexp"
)

func arrayContains(arr []string, str string) bool {
	for _, a := range arr {
	      if a == str {
		 return true
	      }
	}
	return false
}

func doValidation(param string, compare string, value string, message string) (bool, string) {
	LogInfo.Println("Enter doValidation ["+param+" ("+compare+") "+value+"]", message)
	
	result := true
	msg := ""
	
	//add length check
	switch compare { 
		case "empty", "em", "0", "null": 		
			if len(strings.TrimSpace(param)) > 0 {
				result = false
				msg = conditionalDefaultString(message, param+" is not empty")
			}
		case "not_empty", "not empty", "!em", "!0", "not null", "not_null", "!null", "!empty": 
			if len(strings.TrimSpace(param)) == 0 {
				result = false
				msg = conditionalDefaultString(message, param+" is empty")
			}
		case "eq", "equal", "equals", "=": 
			if !strings.EqualFold(strings.TrimSpace(param), strings.TrimSpace(value)) {
				result = false
				msg = conditionalDefaultString(message, param+" is not equals to "+value)
			}
		case "eq_case", "eq case", "equals case", "==", "equal case", "equal_case", "equals_case": 
			if strings.TrimSpace(param) != strings.TrimSpace(value) {
				result = false
				msg = conditionalDefaultString(message, param+" is not equals(with case) to "+value)
			}
		case "not_eq", "not eq", "not equals", "!=", "not equal", "not_equal", "not_equals", "!eq", "!equal", "!equals": 
			if strings.EqualFold(strings.TrimSpace(param), strings.TrimSpace(value)) {
				result = false
				msg = conditionalDefaultString(message, param+" is equals to "+value)
			}
		case "contain", "contains", "found": 
			if !strings.Contains(strings.ToLower(param), strings.ToLower(value)) {
				result = false
				msg = conditionalDefaultString(message, param+" does not contain "+value)
			}
		case "not_contain", "not_contains", "not contain", "not contains", "not_found", "not found", "!found", "!contain", "!contains": 
			if strings.Contains(strings.ToLower(param), strings.ToLower(value)) {
				result = false
				msg = conditionalDefaultString(message, param+" contains "+value)
			}
		case "match", "re": 
			re := regexp.MustCompile(value)
			if !re.MatchString(param) {
				result = false
				msg = conditionalDefaultString(message, param+" not matches expression "+value)
			}
		case "not_match", "not match", "!match", "!re": 
			re := regexp.MustCompile(value)
			if re.MatchString(param) {
				result = false
				msg = conditionalDefaultString(message, param+" matches expression "+value)
			}
		case "lt", "less_than", "less than", "<": 
			iP, errP := strconv.Atoi(param)
			iV, errV := strconv.Atoi(value)
			if (errP != nil || errV != nil){
				result = false
				msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
			} else {
				if !(iP < iV)  {
					result = false
					msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
				}
			}
		case "lt_eq", "less_than_eq", "less_than_equal", "less_than_equals", "lt eq", "less than eq", "less than equal", "less than equals", "<=": 
			iP, errP := strconv.Atoi(param)
			iV, errV := strconv.Atoi(value)
			if (errP != nil || errV != nil){
				result = false
				msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
			} else {
				if !(iP <= iV)  {
					result = false
					msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
				}
			}
		case "gt", "greater_than", "greater than", ">": 
			iP, errP := strconv.Atoi(param)
			iV, errV := strconv.Atoi(value)
			if (errP != nil || errV != nil){
				result = false
				msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
			} else {
				if !(iP > iV)  {
					result = false
					msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
				}
			}
		case "gt_eq", "greater_than_eq", "greater_than_equal", "greater_than_equals", "gt eq", "greater than eq", "greater than equal", "greater than equals", ">=": 
			iP, errP := strconv.Atoi(param)
			iV, errV := strconv.Atoi(value)
			if (errP != nil || errV != nil){
				result = false
				msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
			} else {
				if !(iP >= iV)  {
					result = false
					msg = conditionalDefaultString(message, "failed to compare "+value+" or "+param)
				}
			}
		default:
			result = false
			msg = "un supported comparator check "+compare			
	}  

	LogInfo.Println("Validation result", result, msg)

	return result, msg
}