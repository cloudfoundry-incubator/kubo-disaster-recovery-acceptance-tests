package acceptance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/gomega"
)

type TestCaseFilter map[string]interface{}

func NewTestCaseFilter(path string) TestCaseFilter {
	rawConfig, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	filter := TestCaseFilter{}
	err = json.Unmarshal(rawConfig, &filter)
	Expect(err).NotTo(HaveOccurred())

	return filter
}

func (f TestCaseFilter) Filter(testCases []TestCase) []TestCase {
	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		flagValue := f.getFlagValue(testCase.Name())
		if flagValue == true {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}
	Expect(filteredTestCases).NotTo(BeEmpty(), "must run at least one test case")

	return filteredTestCases
}

func (f TestCaseFilter) getFlagValue(testCaseName string) bool {
	flagName := fmt.Sprintf("include_%s", testCaseName)

	flagValue, isDefined := f[flagName]
	if !isDefined {
		flagValue = false
	}

	boolFlagValue, isBool := flagValue.(bool)
	Expect(isBool).To(BeTrue())

	return boolFlagValue
}
