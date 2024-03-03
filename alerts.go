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

func testAlertRules(response response, alertRules []alertRule) []string {

	alerts := []string{}
	for _, alertRule := range alertRules {
		fail, err := expr.Run(alertRule.Program, response)
		if err != nil {
			alerts = append(alerts, fmt.Sprintf("Failed to evaluate alert rule: %v", err))
		}
		if fail == true {
			alerts = append(alerts, alertRule.Rule)
		}
	}
	return alerts
}
