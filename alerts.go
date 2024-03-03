package main

import (
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

func testAlertRules(response response, alertRules []alertRule) (bool, *string) {

	for _, alertRule := range alertRules {
		fail, err := expr.Run(alertRule.Program, response)
		if err != nil {
			return true, &alertRule.Rule
		}
		if fail == true {
			return true, &alertRule.Rule
		}
	}
	return false, nil
}
