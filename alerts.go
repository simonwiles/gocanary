package main

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type alertRule struct {
	Rule    string
	Program *vm.Program
}

func compileAlerts(alertExprs *[]string) ([]alertRule, error) {
	alertRules := make([]alertRule, len(*alertExprs))
	for i, ruleStr := range *alertExprs {
		program, err := expr.Compile(ruleStr, expr.Env(response{}))
		if err != nil {
			return nil, err
		}
		alertRules[i] = alertRule{
			Rule:    ruleStr,
			Program: program,
		}
	}
	return alertRules, nil
}

func testAlertRules(response response, alertRules []alertRule) map[string]interface{} {

	alerts := make(map[string]interface{})
	for _, alertRule := range alertRules {
		fail, err := expr.Run(alertRule.Program, response)
		if err != nil {
			alerts[alertRule.Rule] = fmt.Sprintf("Failed to evaluate alert rule: %v", err)
			continue
		}
		if fail == false || len(fail.([]interface{})) == 0 {
			continue
		}
		alerts[alertRule.Rule] = fail
		fmt.Printf("Alert: %s\n\t%v\n", alertRule.Rule, fail)
	}
	return alerts
}
