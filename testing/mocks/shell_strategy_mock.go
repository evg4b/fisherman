package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i fisherman/pkg/shell.ShellStrategy -o ./testing/mocks/shell_strategy_mock.go -n ShellStrategyMock

import (
	"context"
	"os/exec"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ShellStrategyMock implements shell.ShellStrategy
type ShellStrategyMock struct {
	t minimock.Tester

	funcArgsWrapper          func(sa1 []string) (sa2 []string)
	inspectFuncArgsWrapper   func(sa1 []string)
	afterArgsWrapperCounter  uint64
	beforeArgsWrapperCounter uint64
	ArgsWrapperMock          mShellStrategyMockArgsWrapper

	funcEnvWrapper          func(sa1 []string) (sa2 []string)
	inspectFuncEnvWrapper   func(sa1 []string)
	afterEnvWrapperCounter  uint64
	beforeEnvWrapperCounter uint64
	EnvWrapperMock          mShellStrategyMockEnvWrapper

	funcGetCommand          func(ctx context.Context) (cp1 *exec.Cmd)
	inspectFuncGetCommand   func(ctx context.Context)
	afterGetCommandCounter  uint64
	beforeGetCommandCounter uint64
	GetCommandMock          mShellStrategyMockGetCommand

	funcGetName          func() (s1 string)
	inspectFuncGetName   func()
	afterGetNameCounter  uint64
	beforeGetNameCounter uint64
	GetNameMock          mShellStrategyMockGetName
}

// NewShellStrategyMock returns a mock for shell.ShellStrategy
func NewShellStrategyMock(t minimock.Tester) *ShellStrategyMock {
	m := &ShellStrategyMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ArgsWrapperMock = mShellStrategyMockArgsWrapper{mock: m}
	m.ArgsWrapperMock.callArgs = []*ShellStrategyMockArgsWrapperParams{}

	m.EnvWrapperMock = mShellStrategyMockEnvWrapper{mock: m}
	m.EnvWrapperMock.callArgs = []*ShellStrategyMockEnvWrapperParams{}

	m.GetCommandMock = mShellStrategyMockGetCommand{mock: m}
	m.GetCommandMock.callArgs = []*ShellStrategyMockGetCommandParams{}

	m.GetNameMock = mShellStrategyMockGetName{mock: m}

	return m
}

type mShellStrategyMockArgsWrapper struct {
	mock               *ShellStrategyMock
	defaultExpectation *ShellStrategyMockArgsWrapperExpectation
	expectations       []*ShellStrategyMockArgsWrapperExpectation

	callArgs []*ShellStrategyMockArgsWrapperParams
	mutex    sync.RWMutex
}

// ShellStrategyMockArgsWrapperExpectation specifies expectation struct of the ShellStrategy.ArgsWrapper
type ShellStrategyMockArgsWrapperExpectation struct {
	mock    *ShellStrategyMock
	params  *ShellStrategyMockArgsWrapperParams
	results *ShellStrategyMockArgsWrapperResults
	Counter uint64
}

// ShellStrategyMockArgsWrapperParams contains parameters of the ShellStrategy.ArgsWrapper
type ShellStrategyMockArgsWrapperParams struct {
	sa1 []string
}

// ShellStrategyMockArgsWrapperResults contains results of the ShellStrategy.ArgsWrapper
type ShellStrategyMockArgsWrapperResults struct {
	sa2 []string
}

// Expect sets up expected params for ShellStrategy.ArgsWrapper
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) Expect(sa1 []string) *mShellStrategyMockArgsWrapper {
	if mmArgsWrapper.mock.funcArgsWrapper != nil {
		mmArgsWrapper.mock.t.Fatalf("ShellStrategyMock.ArgsWrapper mock is already set by Set")
	}

	if mmArgsWrapper.defaultExpectation == nil {
		mmArgsWrapper.defaultExpectation = &ShellStrategyMockArgsWrapperExpectation{}
	}

	mmArgsWrapper.defaultExpectation.params = &ShellStrategyMockArgsWrapperParams{sa1}
	for _, e := range mmArgsWrapper.expectations {
		if minimock.Equal(e.params, mmArgsWrapper.defaultExpectation.params) {
			mmArgsWrapper.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmArgsWrapper.defaultExpectation.params)
		}
	}

	return mmArgsWrapper
}

// Inspect accepts an inspector function that has same arguments as the ShellStrategy.ArgsWrapper
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) Inspect(f func(sa1 []string)) *mShellStrategyMockArgsWrapper {
	if mmArgsWrapper.mock.inspectFuncArgsWrapper != nil {
		mmArgsWrapper.mock.t.Fatalf("Inspect function is already set for ShellStrategyMock.ArgsWrapper")
	}

	mmArgsWrapper.mock.inspectFuncArgsWrapper = f

	return mmArgsWrapper
}

// Return sets up results that will be returned by ShellStrategy.ArgsWrapper
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) Return(sa2 []string) *ShellStrategyMock {
	if mmArgsWrapper.mock.funcArgsWrapper != nil {
		mmArgsWrapper.mock.t.Fatalf("ShellStrategyMock.ArgsWrapper mock is already set by Set")
	}

	if mmArgsWrapper.defaultExpectation == nil {
		mmArgsWrapper.defaultExpectation = &ShellStrategyMockArgsWrapperExpectation{mock: mmArgsWrapper.mock}
	}
	mmArgsWrapper.defaultExpectation.results = &ShellStrategyMockArgsWrapperResults{sa2}
	return mmArgsWrapper.mock
}

//Set uses given function f to mock the ShellStrategy.ArgsWrapper method
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) Set(f func(sa1 []string) (sa2 []string)) *ShellStrategyMock {
	if mmArgsWrapper.defaultExpectation != nil {
		mmArgsWrapper.mock.t.Fatalf("Default expectation is already set for the ShellStrategy.ArgsWrapper method")
	}

	if len(mmArgsWrapper.expectations) > 0 {
		mmArgsWrapper.mock.t.Fatalf("Some expectations are already set for the ShellStrategy.ArgsWrapper method")
	}

	mmArgsWrapper.mock.funcArgsWrapper = f
	return mmArgsWrapper.mock
}

// When sets expectation for the ShellStrategy.ArgsWrapper which will trigger the result defined by the following
// Then helper
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) When(sa1 []string) *ShellStrategyMockArgsWrapperExpectation {
	if mmArgsWrapper.mock.funcArgsWrapper != nil {
		mmArgsWrapper.mock.t.Fatalf("ShellStrategyMock.ArgsWrapper mock is already set by Set")
	}

	expectation := &ShellStrategyMockArgsWrapperExpectation{
		mock:   mmArgsWrapper.mock,
		params: &ShellStrategyMockArgsWrapperParams{sa1},
	}
	mmArgsWrapper.expectations = append(mmArgsWrapper.expectations, expectation)
	return expectation
}

// Then sets up ShellStrategy.ArgsWrapper return parameters for the expectation previously defined by the When method
func (e *ShellStrategyMockArgsWrapperExpectation) Then(sa2 []string) *ShellStrategyMock {
	e.results = &ShellStrategyMockArgsWrapperResults{sa2}
	return e.mock
}

// ArgsWrapper implements shell.ShellStrategy
func (mmArgsWrapper *ShellStrategyMock) ArgsWrapper(sa1 []string) (sa2 []string) {
	mm_atomic.AddUint64(&mmArgsWrapper.beforeArgsWrapperCounter, 1)
	defer mm_atomic.AddUint64(&mmArgsWrapper.afterArgsWrapperCounter, 1)

	if mmArgsWrapper.inspectFuncArgsWrapper != nil {
		mmArgsWrapper.inspectFuncArgsWrapper(sa1)
	}

	mm_params := &ShellStrategyMockArgsWrapperParams{sa1}

	// Record call args
	mmArgsWrapper.ArgsWrapperMock.mutex.Lock()
	mmArgsWrapper.ArgsWrapperMock.callArgs = append(mmArgsWrapper.ArgsWrapperMock.callArgs, mm_params)
	mmArgsWrapper.ArgsWrapperMock.mutex.Unlock()

	for _, e := range mmArgsWrapper.ArgsWrapperMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.sa2
		}
	}

	if mmArgsWrapper.ArgsWrapperMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmArgsWrapper.ArgsWrapperMock.defaultExpectation.Counter, 1)
		mm_want := mmArgsWrapper.ArgsWrapperMock.defaultExpectation.params
		mm_got := ShellStrategyMockArgsWrapperParams{sa1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmArgsWrapper.t.Errorf("ShellStrategyMock.ArgsWrapper got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmArgsWrapper.ArgsWrapperMock.defaultExpectation.results
		if mm_results == nil {
			mmArgsWrapper.t.Fatal("No results are set for the ShellStrategyMock.ArgsWrapper")
		}
		return (*mm_results).sa2
	}
	if mmArgsWrapper.funcArgsWrapper != nil {
		return mmArgsWrapper.funcArgsWrapper(sa1)
	}
	mmArgsWrapper.t.Fatalf("Unexpected call to ShellStrategyMock.ArgsWrapper. %v", sa1)
	return
}

// ArgsWrapperAfterCounter returns a count of finished ShellStrategyMock.ArgsWrapper invocations
func (mmArgsWrapper *ShellStrategyMock) ArgsWrapperAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmArgsWrapper.afterArgsWrapperCounter)
}

// ArgsWrapperBeforeCounter returns a count of ShellStrategyMock.ArgsWrapper invocations
func (mmArgsWrapper *ShellStrategyMock) ArgsWrapperBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmArgsWrapper.beforeArgsWrapperCounter)
}

// Calls returns a list of arguments used in each call to ShellStrategyMock.ArgsWrapper.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmArgsWrapper *mShellStrategyMockArgsWrapper) Calls() []*ShellStrategyMockArgsWrapperParams {
	mmArgsWrapper.mutex.RLock()

	argCopy := make([]*ShellStrategyMockArgsWrapperParams, len(mmArgsWrapper.callArgs))
	copy(argCopy, mmArgsWrapper.callArgs)

	mmArgsWrapper.mutex.RUnlock()

	return argCopy
}

// MinimockArgsWrapperDone returns true if the count of the ArgsWrapper invocations corresponds
// the number of defined expectations
func (m *ShellStrategyMock) MinimockArgsWrapperDone() bool {
	for _, e := range m.ArgsWrapperMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ArgsWrapperMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterArgsWrapperCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcArgsWrapper != nil && mm_atomic.LoadUint64(&m.afterArgsWrapperCounter) < 1 {
		return false
	}
	return true
}

// MinimockArgsWrapperInspect logs each unmet expectation
func (m *ShellStrategyMock) MinimockArgsWrapperInspect() {
	for _, e := range m.ArgsWrapperMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ShellStrategyMock.ArgsWrapper with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ArgsWrapperMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterArgsWrapperCounter) < 1 {
		if m.ArgsWrapperMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ShellStrategyMock.ArgsWrapper")
		} else {
			m.t.Errorf("Expected call to ShellStrategyMock.ArgsWrapper with params: %#v", *m.ArgsWrapperMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcArgsWrapper != nil && mm_atomic.LoadUint64(&m.afterArgsWrapperCounter) < 1 {
		m.t.Error("Expected call to ShellStrategyMock.ArgsWrapper")
	}
}

type mShellStrategyMockEnvWrapper struct {
	mock               *ShellStrategyMock
	defaultExpectation *ShellStrategyMockEnvWrapperExpectation
	expectations       []*ShellStrategyMockEnvWrapperExpectation

	callArgs []*ShellStrategyMockEnvWrapperParams
	mutex    sync.RWMutex
}

// ShellStrategyMockEnvWrapperExpectation specifies expectation struct of the ShellStrategy.EnvWrapper
type ShellStrategyMockEnvWrapperExpectation struct {
	mock    *ShellStrategyMock
	params  *ShellStrategyMockEnvWrapperParams
	results *ShellStrategyMockEnvWrapperResults
	Counter uint64
}

// ShellStrategyMockEnvWrapperParams contains parameters of the ShellStrategy.EnvWrapper
type ShellStrategyMockEnvWrapperParams struct {
	sa1 []string
}

// ShellStrategyMockEnvWrapperResults contains results of the ShellStrategy.EnvWrapper
type ShellStrategyMockEnvWrapperResults struct {
	sa2 []string
}

// Expect sets up expected params for ShellStrategy.EnvWrapper
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) Expect(sa1 []string) *mShellStrategyMockEnvWrapper {
	if mmEnvWrapper.mock.funcEnvWrapper != nil {
		mmEnvWrapper.mock.t.Fatalf("ShellStrategyMock.EnvWrapper mock is already set by Set")
	}

	if mmEnvWrapper.defaultExpectation == nil {
		mmEnvWrapper.defaultExpectation = &ShellStrategyMockEnvWrapperExpectation{}
	}

	mmEnvWrapper.defaultExpectation.params = &ShellStrategyMockEnvWrapperParams{sa1}
	for _, e := range mmEnvWrapper.expectations {
		if minimock.Equal(e.params, mmEnvWrapper.defaultExpectation.params) {
			mmEnvWrapper.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmEnvWrapper.defaultExpectation.params)
		}
	}

	return mmEnvWrapper
}

// Inspect accepts an inspector function that has same arguments as the ShellStrategy.EnvWrapper
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) Inspect(f func(sa1 []string)) *mShellStrategyMockEnvWrapper {
	if mmEnvWrapper.mock.inspectFuncEnvWrapper != nil {
		mmEnvWrapper.mock.t.Fatalf("Inspect function is already set for ShellStrategyMock.EnvWrapper")
	}

	mmEnvWrapper.mock.inspectFuncEnvWrapper = f

	return mmEnvWrapper
}

// Return sets up results that will be returned by ShellStrategy.EnvWrapper
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) Return(sa2 []string) *ShellStrategyMock {
	if mmEnvWrapper.mock.funcEnvWrapper != nil {
		mmEnvWrapper.mock.t.Fatalf("ShellStrategyMock.EnvWrapper mock is already set by Set")
	}

	if mmEnvWrapper.defaultExpectation == nil {
		mmEnvWrapper.defaultExpectation = &ShellStrategyMockEnvWrapperExpectation{mock: mmEnvWrapper.mock}
	}
	mmEnvWrapper.defaultExpectation.results = &ShellStrategyMockEnvWrapperResults{sa2}
	return mmEnvWrapper.mock
}

//Set uses given function f to mock the ShellStrategy.EnvWrapper method
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) Set(f func(sa1 []string) (sa2 []string)) *ShellStrategyMock {
	if mmEnvWrapper.defaultExpectation != nil {
		mmEnvWrapper.mock.t.Fatalf("Default expectation is already set for the ShellStrategy.EnvWrapper method")
	}

	if len(mmEnvWrapper.expectations) > 0 {
		mmEnvWrapper.mock.t.Fatalf("Some expectations are already set for the ShellStrategy.EnvWrapper method")
	}

	mmEnvWrapper.mock.funcEnvWrapper = f
	return mmEnvWrapper.mock
}

// When sets expectation for the ShellStrategy.EnvWrapper which will trigger the result defined by the following
// Then helper
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) When(sa1 []string) *ShellStrategyMockEnvWrapperExpectation {
	if mmEnvWrapper.mock.funcEnvWrapper != nil {
		mmEnvWrapper.mock.t.Fatalf("ShellStrategyMock.EnvWrapper mock is already set by Set")
	}

	expectation := &ShellStrategyMockEnvWrapperExpectation{
		mock:   mmEnvWrapper.mock,
		params: &ShellStrategyMockEnvWrapperParams{sa1},
	}
	mmEnvWrapper.expectations = append(mmEnvWrapper.expectations, expectation)
	return expectation
}

// Then sets up ShellStrategy.EnvWrapper return parameters for the expectation previously defined by the When method
func (e *ShellStrategyMockEnvWrapperExpectation) Then(sa2 []string) *ShellStrategyMock {
	e.results = &ShellStrategyMockEnvWrapperResults{sa2}
	return e.mock
}

// EnvWrapper implements shell.ShellStrategy
func (mmEnvWrapper *ShellStrategyMock) EnvWrapper(sa1 []string) (sa2 []string) {
	mm_atomic.AddUint64(&mmEnvWrapper.beforeEnvWrapperCounter, 1)
	defer mm_atomic.AddUint64(&mmEnvWrapper.afterEnvWrapperCounter, 1)

	if mmEnvWrapper.inspectFuncEnvWrapper != nil {
		mmEnvWrapper.inspectFuncEnvWrapper(sa1)
	}

	mm_params := &ShellStrategyMockEnvWrapperParams{sa1}

	// Record call args
	mmEnvWrapper.EnvWrapperMock.mutex.Lock()
	mmEnvWrapper.EnvWrapperMock.callArgs = append(mmEnvWrapper.EnvWrapperMock.callArgs, mm_params)
	mmEnvWrapper.EnvWrapperMock.mutex.Unlock()

	for _, e := range mmEnvWrapper.EnvWrapperMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.sa2
		}
	}

	if mmEnvWrapper.EnvWrapperMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmEnvWrapper.EnvWrapperMock.defaultExpectation.Counter, 1)
		mm_want := mmEnvWrapper.EnvWrapperMock.defaultExpectation.params
		mm_got := ShellStrategyMockEnvWrapperParams{sa1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmEnvWrapper.t.Errorf("ShellStrategyMock.EnvWrapper got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmEnvWrapper.EnvWrapperMock.defaultExpectation.results
		if mm_results == nil {
			mmEnvWrapper.t.Fatal("No results are set for the ShellStrategyMock.EnvWrapper")
		}
		return (*mm_results).sa2
	}
	if mmEnvWrapper.funcEnvWrapper != nil {
		return mmEnvWrapper.funcEnvWrapper(sa1)
	}
	mmEnvWrapper.t.Fatalf("Unexpected call to ShellStrategyMock.EnvWrapper. %v", sa1)
	return
}

// EnvWrapperAfterCounter returns a count of finished ShellStrategyMock.EnvWrapper invocations
func (mmEnvWrapper *ShellStrategyMock) EnvWrapperAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmEnvWrapper.afterEnvWrapperCounter)
}

// EnvWrapperBeforeCounter returns a count of ShellStrategyMock.EnvWrapper invocations
func (mmEnvWrapper *ShellStrategyMock) EnvWrapperBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmEnvWrapper.beforeEnvWrapperCounter)
}

// Calls returns a list of arguments used in each call to ShellStrategyMock.EnvWrapper.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmEnvWrapper *mShellStrategyMockEnvWrapper) Calls() []*ShellStrategyMockEnvWrapperParams {
	mmEnvWrapper.mutex.RLock()

	argCopy := make([]*ShellStrategyMockEnvWrapperParams, len(mmEnvWrapper.callArgs))
	copy(argCopy, mmEnvWrapper.callArgs)

	mmEnvWrapper.mutex.RUnlock()

	return argCopy
}

// MinimockEnvWrapperDone returns true if the count of the EnvWrapper invocations corresponds
// the number of defined expectations
func (m *ShellStrategyMock) MinimockEnvWrapperDone() bool {
	for _, e := range m.EnvWrapperMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.EnvWrapperMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterEnvWrapperCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcEnvWrapper != nil && mm_atomic.LoadUint64(&m.afterEnvWrapperCounter) < 1 {
		return false
	}
	return true
}

// MinimockEnvWrapperInspect logs each unmet expectation
func (m *ShellStrategyMock) MinimockEnvWrapperInspect() {
	for _, e := range m.EnvWrapperMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ShellStrategyMock.EnvWrapper with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.EnvWrapperMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterEnvWrapperCounter) < 1 {
		if m.EnvWrapperMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ShellStrategyMock.EnvWrapper")
		} else {
			m.t.Errorf("Expected call to ShellStrategyMock.EnvWrapper with params: %#v", *m.EnvWrapperMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcEnvWrapper != nil && mm_atomic.LoadUint64(&m.afterEnvWrapperCounter) < 1 {
		m.t.Error("Expected call to ShellStrategyMock.EnvWrapper")
	}
}

type mShellStrategyMockGetCommand struct {
	mock               *ShellStrategyMock
	defaultExpectation *ShellStrategyMockGetCommandExpectation
	expectations       []*ShellStrategyMockGetCommandExpectation

	callArgs []*ShellStrategyMockGetCommandParams
	mutex    sync.RWMutex
}

// ShellStrategyMockGetCommandExpectation specifies expectation struct of the ShellStrategy.GetCommand
type ShellStrategyMockGetCommandExpectation struct {
	mock    *ShellStrategyMock
	params  *ShellStrategyMockGetCommandParams
	results *ShellStrategyMockGetCommandResults
	Counter uint64
}

// ShellStrategyMockGetCommandParams contains parameters of the ShellStrategy.GetCommand
type ShellStrategyMockGetCommandParams struct {
	ctx context.Context
}

// ShellStrategyMockGetCommandResults contains results of the ShellStrategy.GetCommand
type ShellStrategyMockGetCommandResults struct {
	cp1 *exec.Cmd
}

// Expect sets up expected params for ShellStrategy.GetCommand
func (mmGetCommand *mShellStrategyMockGetCommand) Expect(ctx context.Context) *mShellStrategyMockGetCommand {
	if mmGetCommand.mock.funcGetCommand != nil {
		mmGetCommand.mock.t.Fatalf("ShellStrategyMock.GetCommand mock is already set by Set")
	}

	if mmGetCommand.defaultExpectation == nil {
		mmGetCommand.defaultExpectation = &ShellStrategyMockGetCommandExpectation{}
	}

	mmGetCommand.defaultExpectation.params = &ShellStrategyMockGetCommandParams{ctx}
	for _, e := range mmGetCommand.expectations {
		if minimock.Equal(e.params, mmGetCommand.defaultExpectation.params) {
			mmGetCommand.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetCommand.defaultExpectation.params)
		}
	}

	return mmGetCommand
}

// Inspect accepts an inspector function that has same arguments as the ShellStrategy.GetCommand
func (mmGetCommand *mShellStrategyMockGetCommand) Inspect(f func(ctx context.Context)) *mShellStrategyMockGetCommand {
	if mmGetCommand.mock.inspectFuncGetCommand != nil {
		mmGetCommand.mock.t.Fatalf("Inspect function is already set for ShellStrategyMock.GetCommand")
	}

	mmGetCommand.mock.inspectFuncGetCommand = f

	return mmGetCommand
}

// Return sets up results that will be returned by ShellStrategy.GetCommand
func (mmGetCommand *mShellStrategyMockGetCommand) Return(cp1 *exec.Cmd) *ShellStrategyMock {
	if mmGetCommand.mock.funcGetCommand != nil {
		mmGetCommand.mock.t.Fatalf("ShellStrategyMock.GetCommand mock is already set by Set")
	}

	if mmGetCommand.defaultExpectation == nil {
		mmGetCommand.defaultExpectation = &ShellStrategyMockGetCommandExpectation{mock: mmGetCommand.mock}
	}
	mmGetCommand.defaultExpectation.results = &ShellStrategyMockGetCommandResults{cp1}
	return mmGetCommand.mock
}

//Set uses given function f to mock the ShellStrategy.GetCommand method
func (mmGetCommand *mShellStrategyMockGetCommand) Set(f func(ctx context.Context) (cp1 *exec.Cmd)) *ShellStrategyMock {
	if mmGetCommand.defaultExpectation != nil {
		mmGetCommand.mock.t.Fatalf("Default expectation is already set for the ShellStrategy.GetCommand method")
	}

	if len(mmGetCommand.expectations) > 0 {
		mmGetCommand.mock.t.Fatalf("Some expectations are already set for the ShellStrategy.GetCommand method")
	}

	mmGetCommand.mock.funcGetCommand = f
	return mmGetCommand.mock
}

// When sets expectation for the ShellStrategy.GetCommand which will trigger the result defined by the following
// Then helper
func (mmGetCommand *mShellStrategyMockGetCommand) When(ctx context.Context) *ShellStrategyMockGetCommandExpectation {
	if mmGetCommand.mock.funcGetCommand != nil {
		mmGetCommand.mock.t.Fatalf("ShellStrategyMock.GetCommand mock is already set by Set")
	}

	expectation := &ShellStrategyMockGetCommandExpectation{
		mock:   mmGetCommand.mock,
		params: &ShellStrategyMockGetCommandParams{ctx},
	}
	mmGetCommand.expectations = append(mmGetCommand.expectations, expectation)
	return expectation
}

// Then sets up ShellStrategy.GetCommand return parameters for the expectation previously defined by the When method
func (e *ShellStrategyMockGetCommandExpectation) Then(cp1 *exec.Cmd) *ShellStrategyMock {
	e.results = &ShellStrategyMockGetCommandResults{cp1}
	return e.mock
}

// GetCommand implements shell.ShellStrategy
func (mmGetCommand *ShellStrategyMock) GetCommand(ctx context.Context) (cp1 *exec.Cmd) {
	mm_atomic.AddUint64(&mmGetCommand.beforeGetCommandCounter, 1)
	defer mm_atomic.AddUint64(&mmGetCommand.afterGetCommandCounter, 1)

	if mmGetCommand.inspectFuncGetCommand != nil {
		mmGetCommand.inspectFuncGetCommand(ctx)
	}

	mm_params := &ShellStrategyMockGetCommandParams{ctx}

	// Record call args
	mmGetCommand.GetCommandMock.mutex.Lock()
	mmGetCommand.GetCommandMock.callArgs = append(mmGetCommand.GetCommandMock.callArgs, mm_params)
	mmGetCommand.GetCommandMock.mutex.Unlock()

	for _, e := range mmGetCommand.GetCommandMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.cp1
		}
	}

	if mmGetCommand.GetCommandMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetCommand.GetCommandMock.defaultExpectation.Counter, 1)
		mm_want := mmGetCommand.GetCommandMock.defaultExpectation.params
		mm_got := ShellStrategyMockGetCommandParams{ctx}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetCommand.t.Errorf("ShellStrategyMock.GetCommand got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetCommand.GetCommandMock.defaultExpectation.results
		if mm_results == nil {
			mmGetCommand.t.Fatal("No results are set for the ShellStrategyMock.GetCommand")
		}
		return (*mm_results).cp1
	}
	if mmGetCommand.funcGetCommand != nil {
		return mmGetCommand.funcGetCommand(ctx)
	}
	mmGetCommand.t.Fatalf("Unexpected call to ShellStrategyMock.GetCommand. %v", ctx)
	return
}

// GetCommandAfterCounter returns a count of finished ShellStrategyMock.GetCommand invocations
func (mmGetCommand *ShellStrategyMock) GetCommandAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCommand.afterGetCommandCounter)
}

// GetCommandBeforeCounter returns a count of ShellStrategyMock.GetCommand invocations
func (mmGetCommand *ShellStrategyMock) GetCommandBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCommand.beforeGetCommandCounter)
}

// Calls returns a list of arguments used in each call to ShellStrategyMock.GetCommand.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetCommand *mShellStrategyMockGetCommand) Calls() []*ShellStrategyMockGetCommandParams {
	mmGetCommand.mutex.RLock()

	argCopy := make([]*ShellStrategyMockGetCommandParams, len(mmGetCommand.callArgs))
	copy(argCopy, mmGetCommand.callArgs)

	mmGetCommand.mutex.RUnlock()

	return argCopy
}

// MinimockGetCommandDone returns true if the count of the GetCommand invocations corresponds
// the number of defined expectations
func (m *ShellStrategyMock) MinimockGetCommandDone() bool {
	for _, e := range m.GetCommandMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCommandMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCommandCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCommand != nil && mm_atomic.LoadUint64(&m.afterGetCommandCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetCommandInspect logs each unmet expectation
func (m *ShellStrategyMock) MinimockGetCommandInspect() {
	for _, e := range m.GetCommandMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ShellStrategyMock.GetCommand with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCommandMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCommandCounter) < 1 {
		if m.GetCommandMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ShellStrategyMock.GetCommand")
		} else {
			m.t.Errorf("Expected call to ShellStrategyMock.GetCommand with params: %#v", *m.GetCommandMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCommand != nil && mm_atomic.LoadUint64(&m.afterGetCommandCounter) < 1 {
		m.t.Error("Expected call to ShellStrategyMock.GetCommand")
	}
}

type mShellStrategyMockGetName struct {
	mock               *ShellStrategyMock
	defaultExpectation *ShellStrategyMockGetNameExpectation
	expectations       []*ShellStrategyMockGetNameExpectation
}

// ShellStrategyMockGetNameExpectation specifies expectation struct of the ShellStrategy.GetName
type ShellStrategyMockGetNameExpectation struct {
	mock *ShellStrategyMock

	results *ShellStrategyMockGetNameResults
	Counter uint64
}

// ShellStrategyMockGetNameResults contains results of the ShellStrategy.GetName
type ShellStrategyMockGetNameResults struct {
	s1 string
}

// Expect sets up expected params for ShellStrategy.GetName
func (mmGetName *mShellStrategyMockGetName) Expect() *mShellStrategyMockGetName {
	if mmGetName.mock.funcGetName != nil {
		mmGetName.mock.t.Fatalf("ShellStrategyMock.GetName mock is already set by Set")
	}

	if mmGetName.defaultExpectation == nil {
		mmGetName.defaultExpectation = &ShellStrategyMockGetNameExpectation{}
	}

	return mmGetName
}

// Inspect accepts an inspector function that has same arguments as the ShellStrategy.GetName
func (mmGetName *mShellStrategyMockGetName) Inspect(f func()) *mShellStrategyMockGetName {
	if mmGetName.mock.inspectFuncGetName != nil {
		mmGetName.mock.t.Fatalf("Inspect function is already set for ShellStrategyMock.GetName")
	}

	mmGetName.mock.inspectFuncGetName = f

	return mmGetName
}

// Return sets up results that will be returned by ShellStrategy.GetName
func (mmGetName *mShellStrategyMockGetName) Return(s1 string) *ShellStrategyMock {
	if mmGetName.mock.funcGetName != nil {
		mmGetName.mock.t.Fatalf("ShellStrategyMock.GetName mock is already set by Set")
	}

	if mmGetName.defaultExpectation == nil {
		mmGetName.defaultExpectation = &ShellStrategyMockGetNameExpectation{mock: mmGetName.mock}
	}
	mmGetName.defaultExpectation.results = &ShellStrategyMockGetNameResults{s1}
	return mmGetName.mock
}

//Set uses given function f to mock the ShellStrategy.GetName method
func (mmGetName *mShellStrategyMockGetName) Set(f func() (s1 string)) *ShellStrategyMock {
	if mmGetName.defaultExpectation != nil {
		mmGetName.mock.t.Fatalf("Default expectation is already set for the ShellStrategy.GetName method")
	}

	if len(mmGetName.expectations) > 0 {
		mmGetName.mock.t.Fatalf("Some expectations are already set for the ShellStrategy.GetName method")
	}

	mmGetName.mock.funcGetName = f
	return mmGetName.mock
}

// GetName implements shell.ShellStrategy
func (mmGetName *ShellStrategyMock) GetName() (s1 string) {
	mm_atomic.AddUint64(&mmGetName.beforeGetNameCounter, 1)
	defer mm_atomic.AddUint64(&mmGetName.afterGetNameCounter, 1)

	if mmGetName.inspectFuncGetName != nil {
		mmGetName.inspectFuncGetName()
	}

	if mmGetName.GetNameMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetName.GetNameMock.defaultExpectation.Counter, 1)

		mm_results := mmGetName.GetNameMock.defaultExpectation.results
		if mm_results == nil {
			mmGetName.t.Fatal("No results are set for the ShellStrategyMock.GetName")
		}
		return (*mm_results).s1
	}
	if mmGetName.funcGetName != nil {
		return mmGetName.funcGetName()
	}
	mmGetName.t.Fatalf("Unexpected call to ShellStrategyMock.GetName.")
	return
}

// GetNameAfterCounter returns a count of finished ShellStrategyMock.GetName invocations
func (mmGetName *ShellStrategyMock) GetNameAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetName.afterGetNameCounter)
}

// GetNameBeforeCounter returns a count of ShellStrategyMock.GetName invocations
func (mmGetName *ShellStrategyMock) GetNameBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetName.beforeGetNameCounter)
}

// MinimockGetNameDone returns true if the count of the GetName invocations corresponds
// the number of defined expectations
func (m *ShellStrategyMock) MinimockGetNameDone() bool {
	for _, e := range m.GetNameMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNameMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNameCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetName != nil && mm_atomic.LoadUint64(&m.afterGetNameCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetNameInspect logs each unmet expectation
func (m *ShellStrategyMock) MinimockGetNameInspect() {
	for _, e := range m.GetNameMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ShellStrategyMock.GetName")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNameMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNameCounter) < 1 {
		m.t.Error("Expected call to ShellStrategyMock.GetName")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetName != nil && mm_atomic.LoadUint64(&m.afterGetNameCounter) < 1 {
		m.t.Error("Expected call to ShellStrategyMock.GetName")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ShellStrategyMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockArgsWrapperInspect()

		m.MinimockEnvWrapperInspect()

		m.MinimockGetCommandInspect()

		m.MinimockGetNameInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ShellStrategyMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ShellStrategyMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockArgsWrapperDone() &&
		m.MinimockEnvWrapperDone() &&
		m.MinimockGetCommandDone() &&
		m.MinimockGetNameDone()
}