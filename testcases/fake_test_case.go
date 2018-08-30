package testcases

type FakeTestCase struct{}

func (t FakeTestCase) Name() string {
	return "fake_test_case"
}

func (t FakeTestCase) BeforeBackup(config Config) {}

func (t FakeTestCase) AfterBackup(config Config) {}

func (t FakeTestCase) AfterRestore(config Config) {}

func (t FakeTestCase) Cleanup(config Config) {}
