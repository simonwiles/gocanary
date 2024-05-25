package main

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type alertRule struct {
	Rule    string
	Program *vm.Program
	Label   string
}

func compileAlerts(alertExprs *[]string) ([]alertRule, error) {
	alertRules := make([]alertRule, len(*alertExprs))
	for i, ruleDef := range *alertExprs {

		ruleStr, label, hasLabel := strings.Cut(ruleDef, "::")
		ruleStr = strings.TrimSpace(ruleStr)
		label = strings.TrimSpace(label)
		if !hasLabel {
			label = ruleStr
		}

		program, err := expr.Compile(ruleStr, expr.Env(response{}))
		if err != nil {
			return nil, err
		}

		alertRules[i] = alertRule{
			Rule:    ruleStr,
			Program: program,
			Label:   label,
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
		alerts[alertRule.Label] = fail
		fmt.Printf("Alert: %s\n\t%v\n", alertRule.Rule, fail)
	}
	return alerts
}
