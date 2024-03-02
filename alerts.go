package main

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type alertRule struct {
	rule    string
	Program *vm.Program
}

func compileAlerts(alertRules *[]string) ([]alertRule, error) {
	alerts := make([]alertRule, len(*alertRules))
	for i, ruleStr := range *alertRules {
		program, err := expr.Compile(ruleStr, expr.Env(disk{}))
		if err != nil {
			return nil, err
		}
		alerts[i] = alertRule{
			rule:    ruleStr,
			Program: program,
		}
	}
	return alerts, nil
}
