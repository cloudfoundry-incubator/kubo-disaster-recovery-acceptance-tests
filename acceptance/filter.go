package acceptance

import "fmt"
import (
	"encoding/json"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

type ConfigTestCaseFilter map[string]interface{}

func NewConfigTestCaseFilter(configPath string) ConfigTestCaseFilter {
	rawConfig, err := ioutil.ReadFile(configPath)
	Expect(err).NotTo(HaveOccurred())

	filter := ConfigTestCaseFilter{}
	err = json.Unmarshal(rawConfig, &filter)
	Expect(err).NotTo(HaveOccurred())

	return filter
}

func (f ConfigTestCaseFilter) Filter(testCases []TestCase) []TestCase {
	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		flagValue := f.getFlagValue(testCase.Name())
		if flagValue == true {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}
	Expect(filteredTestCases).NotTo(BeEmpty())

	return filteredTestCases
}

func (f ConfigTestCaseFilter) getFlagValue(testCaseName string) bool {
	flagName := fmt.Sprintf("include_%s", testCaseName)

	flagValue, isDefined := f[flagName]
	if !isDefined {
		flagValue = false
	}

	boolFlagValue, isBool := flagValue.(bool)
	Expect(isBool).To(BeTrue())

	return boolFlagValue
}
