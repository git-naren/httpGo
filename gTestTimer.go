package main

import ( 
    "time"
)

func timerStepExecutor(cm Config, step *TestStep) { 
    TID := cm[TEST_RUN_ID]
    LogInfo.Println(TID, "Enter timerStepExecutor "+ step.Description, "StartTime", time.Now().Unix())
    
    LogInfo.Println(TID, "WaitTimeInSec:", step.WaitTimeInSec)

    time.Sleep(time.Duration(step.WaitTimeInSec) * time.Second)    
    step.Status = true
    
    LogInfo.Println(TID, "timerStepExecutor EndTime:", time.Now().Unix())
}