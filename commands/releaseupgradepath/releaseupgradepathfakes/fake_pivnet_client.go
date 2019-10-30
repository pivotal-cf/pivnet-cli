// Code generated by counterfeiter. DO NOT EDIT.
package releaseupgradepathfakes

import (
	"sync"

	pivnet "github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/pivnet-cli/commands/releaseupgradepath"
)

type FakePivnetClient struct {
	AddReleaseUpgradePathStub        func(string, int, int) error
	addReleaseUpgradePathMutex       sync.RWMutex
	addReleaseUpgradePathArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
	}
	addReleaseUpgradePathReturns struct {
		result1 error
	}
	addReleaseUpgradePathReturnsOnCall map[int]struct {
		result1 error
	}
	ReleaseForVersionStub        func(string, string) (pivnet.Release, error)
	releaseForVersionMutex       sync.RWMutex
	releaseForVersionArgsForCall []struct {
		arg1 string
		arg2 string
	}
	releaseForVersionReturns struct {
		result1 pivnet.Release
		result2 error
	}
	releaseForVersionReturnsOnCall map[int]struct {
		result1 pivnet.Release
		result2 error
	}
	ReleaseUpgradePathsStub        func(string, int) ([]pivnet.ReleaseUpgradePath, error)
	releaseUpgradePathsMutex       sync.RWMutex
	releaseUpgradePathsArgsForCall []struct {
		arg1 string
		arg2 int
	}
	releaseUpgradePathsReturns struct {
		result1 []pivnet.ReleaseUpgradePath
		result2 error
	}
	releaseUpgradePathsReturnsOnCall map[int]struct {
		result1 []pivnet.ReleaseUpgradePath
		result2 error
	}
	ReleasesForProductSlugStub        func(string, ...pivnet.QueryParameter) ([]pivnet.Release, error)
	releasesForProductSlugMutex       sync.RWMutex
	releasesForProductSlugArgsForCall []struct {
		arg1 string
		arg2 []pivnet.QueryParameter
	}
	releasesForProductSlugReturns struct {
		result1 []pivnet.Release
		result2 error
	}
	releasesForProductSlugReturnsOnCall map[int]struct {
		result1 []pivnet.Release
		result2 error
	}
	RemoveReleaseUpgradePathStub        func(string, int, int) error
	removeReleaseUpgradePathMutex       sync.RWMutex
	removeReleaseUpgradePathArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
	}
	removeReleaseUpgradePathReturns struct {
		result1 error
	}
	removeReleaseUpgradePathReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePivnetClient) AddReleaseUpgradePath(arg1 string, arg2 int, arg3 int) error {
	fake.addReleaseUpgradePathMutex.Lock()
	ret, specificReturn := fake.addReleaseUpgradePathReturnsOnCall[len(fake.addReleaseUpgradePathArgsForCall)]
	fake.addReleaseUpgradePathArgsForCall = append(fake.addReleaseUpgradePathArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	fake.recordInvocation("AddReleaseUpgradePath", []interface{}{arg1, arg2, arg3})
	fake.addReleaseUpgradePathMutex.Unlock()
	if fake.AddReleaseUpgradePathStub != nil {
		return fake.AddReleaseUpgradePathStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addReleaseUpgradePathReturns
	return fakeReturns.result1
}

func (fake *FakePivnetClient) AddReleaseUpgradePathCallCount() int {
	fake.addReleaseUpgradePathMutex.RLock()
	defer fake.addReleaseUpgradePathMutex.RUnlock()
	return len(fake.addReleaseUpgradePathArgsForCall)
}

func (fake *FakePivnetClient) AddReleaseUpgradePathCalls(stub func(string, int, int) error) {
	fake.addReleaseUpgradePathMutex.Lock()
	defer fake.addReleaseUpgradePathMutex.Unlock()
	fake.AddReleaseUpgradePathStub = stub
}

func (fake *FakePivnetClient) AddReleaseUpgradePathArgsForCall(i int) (string, int, int) {
	fake.addReleaseUpgradePathMutex.RLock()
	defer fake.addReleaseUpgradePathMutex.RUnlock()
	argsForCall := fake.addReleaseUpgradePathArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) AddReleaseUpgradePathReturns(result1 error) {
	fake.addReleaseUpgradePathMutex.Lock()
	defer fake.addReleaseUpgradePathMutex.Unlock()
	fake.AddReleaseUpgradePathStub = nil
	fake.addReleaseUpgradePathReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) AddReleaseUpgradePathReturnsOnCall(i int, result1 error) {
	fake.addReleaseUpgradePathMutex.Lock()
	defer fake.addReleaseUpgradePathMutex.Unlock()
	fake.AddReleaseUpgradePathStub = nil
	if fake.addReleaseUpgradePathReturnsOnCall == nil {
		fake.addReleaseUpgradePathReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addReleaseUpgradePathReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) ReleaseForVersion(arg1 string, arg2 string) (pivnet.Release, error) {
	fake.releaseForVersionMutex.Lock()
	ret, specificReturn := fake.releaseForVersionReturnsOnCall[len(fake.releaseForVersionArgsForCall)]
	fake.releaseForVersionArgsForCall = append(fake.releaseForVersionArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ReleaseForVersion", []interface{}{arg1, arg2})
	fake.releaseForVersionMutex.Unlock()
	if fake.ReleaseForVersionStub != nil {
		return fake.ReleaseForVersionStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.releaseForVersionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ReleaseForVersionCallCount() int {
	fake.releaseForVersionMutex.RLock()
	defer fake.releaseForVersionMutex.RUnlock()
	return len(fake.releaseForVersionArgsForCall)
}

func (fake *FakePivnetClient) ReleaseForVersionCalls(stub func(string, string) (pivnet.Release, error)) {
	fake.releaseForVersionMutex.Lock()
	defer fake.releaseForVersionMutex.Unlock()
	fake.ReleaseForVersionStub = stub
}

func (fake *FakePivnetClient) ReleaseForVersionArgsForCall(i int) (string, string) {
	fake.releaseForVersionMutex.RLock()
	defer fake.releaseForVersionMutex.RUnlock()
	argsForCall := fake.releaseForVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ReleaseForVersionReturns(result1 pivnet.Release, result2 error) {
	fake.releaseForVersionMutex.Lock()
	defer fake.releaseForVersionMutex.Unlock()
	fake.ReleaseForVersionStub = nil
	fake.releaseForVersionReturns = struct {
		result1 pivnet.Release
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ReleaseForVersionReturnsOnCall(i int, result1 pivnet.Release, result2 error) {
	fake.releaseForVersionMutex.Lock()
	defer fake.releaseForVersionMutex.Unlock()
	fake.ReleaseForVersionStub = nil
	if fake.releaseForVersionReturnsOnCall == nil {
		fake.releaseForVersionReturnsOnCall = make(map[int]struct {
			result1 pivnet.Release
			result2 error
		})
	}
	fake.releaseForVersionReturnsOnCall[i] = struct {
		result1 pivnet.Release
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ReleaseUpgradePaths(arg1 string, arg2 int) ([]pivnet.ReleaseUpgradePath, error) {
	fake.releaseUpgradePathsMutex.Lock()
	ret, specificReturn := fake.releaseUpgradePathsReturnsOnCall[len(fake.releaseUpgradePathsArgsForCall)]
	fake.releaseUpgradePathsArgsForCall = append(fake.releaseUpgradePathsArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("ReleaseUpgradePaths", []interface{}{arg1, arg2})
	fake.releaseUpgradePathsMutex.Unlock()
	if fake.ReleaseUpgradePathsStub != nil {
		return fake.ReleaseUpgradePathsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.releaseUpgradePathsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ReleaseUpgradePathsCallCount() int {
	fake.releaseUpgradePathsMutex.RLock()
	defer fake.releaseUpgradePathsMutex.RUnlock()
	return len(fake.releaseUpgradePathsArgsForCall)
}

func (fake *FakePivnetClient) ReleaseUpgradePathsCalls(stub func(string, int) ([]pivnet.ReleaseUpgradePath, error)) {
	fake.releaseUpgradePathsMutex.Lock()
	defer fake.releaseUpgradePathsMutex.Unlock()
	fake.ReleaseUpgradePathsStub = stub
}

func (fake *FakePivnetClient) ReleaseUpgradePathsArgsForCall(i int) (string, int) {
	fake.releaseUpgradePathsMutex.RLock()
	defer fake.releaseUpgradePathsMutex.RUnlock()
	argsForCall := fake.releaseUpgradePathsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ReleaseUpgradePathsReturns(result1 []pivnet.ReleaseUpgradePath, result2 error) {
	fake.releaseUpgradePathsMutex.Lock()
	defer fake.releaseUpgradePathsMutex.Unlock()
	fake.ReleaseUpgradePathsStub = nil
	fake.releaseUpgradePathsReturns = struct {
		result1 []pivnet.ReleaseUpgradePath
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ReleaseUpgradePathsReturnsOnCall(i int, result1 []pivnet.ReleaseUpgradePath, result2 error) {
	fake.releaseUpgradePathsMutex.Lock()
	defer fake.releaseUpgradePathsMutex.Unlock()
	fake.ReleaseUpgradePathsStub = nil
	if fake.releaseUpgradePathsReturnsOnCall == nil {
		fake.releaseUpgradePathsReturnsOnCall = make(map[int]struct {
			result1 []pivnet.ReleaseUpgradePath
			result2 error
		})
	}
	fake.releaseUpgradePathsReturnsOnCall[i] = struct {
		result1 []pivnet.ReleaseUpgradePath
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ReleasesForProductSlug(arg1 string, arg2 ...pivnet.QueryParameter) ([]pivnet.Release, error) {
	fake.releasesForProductSlugMutex.Lock()
	ret, specificReturn := fake.releasesForProductSlugReturnsOnCall[len(fake.releasesForProductSlugArgsForCall)]
	fake.releasesForProductSlugArgsForCall = append(fake.releasesForProductSlugArgsForCall, struct {
		arg1 string
		arg2 []pivnet.QueryParameter
	}{arg1, arg2})
	fake.recordInvocation("ReleasesForProductSlug", []interface{}{arg1, arg2})
	fake.releasesForProductSlugMutex.Unlock()
	if fake.ReleasesForProductSlugStub != nil {
		return fake.ReleasesForProductSlugStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.releasesForProductSlugReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ReleasesForProductSlugCallCount() int {
	fake.releasesForProductSlugMutex.RLock()
	defer fake.releasesForProductSlugMutex.RUnlock()
	return len(fake.releasesForProductSlugArgsForCall)
}

func (fake *FakePivnetClient) ReleasesForProductSlugCalls(stub func(string, ...pivnet.QueryParameter) ([]pivnet.Release, error)) {
	fake.releasesForProductSlugMutex.Lock()
	defer fake.releasesForProductSlugMutex.Unlock()
	fake.ReleasesForProductSlugStub = stub
}

func (fake *FakePivnetClient) ReleasesForProductSlugArgsForCall(i int) (string, []pivnet.QueryParameter) {
	fake.releasesForProductSlugMutex.RLock()
	defer fake.releasesForProductSlugMutex.RUnlock()
	argsForCall := fake.releasesForProductSlugArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ReleasesForProductSlugReturns(result1 []pivnet.Release, result2 error) {
	fake.releasesForProductSlugMutex.Lock()
	defer fake.releasesForProductSlugMutex.Unlock()
	fake.ReleasesForProductSlugStub = nil
	fake.releasesForProductSlugReturns = struct {
		result1 []pivnet.Release
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ReleasesForProductSlugReturnsOnCall(i int, result1 []pivnet.Release, result2 error) {
	fake.releasesForProductSlugMutex.Lock()
	defer fake.releasesForProductSlugMutex.Unlock()
	fake.ReleasesForProductSlugStub = nil
	if fake.releasesForProductSlugReturnsOnCall == nil {
		fake.releasesForProductSlugReturnsOnCall = make(map[int]struct {
			result1 []pivnet.Release
			result2 error
		})
	}
	fake.releasesForProductSlugReturnsOnCall[i] = struct {
		result1 []pivnet.Release
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePath(arg1 string, arg2 int, arg3 int) error {
	fake.removeReleaseUpgradePathMutex.Lock()
	ret, specificReturn := fake.removeReleaseUpgradePathReturnsOnCall[len(fake.removeReleaseUpgradePathArgsForCall)]
	fake.removeReleaseUpgradePathArgsForCall = append(fake.removeReleaseUpgradePathArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	fake.recordInvocation("RemoveReleaseUpgradePath", []interface{}{arg1, arg2, arg3})
	fake.removeReleaseUpgradePathMutex.Unlock()
	if fake.RemoveReleaseUpgradePathStub != nil {
		return fake.RemoveReleaseUpgradePathStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeReleaseUpgradePathReturns
	return fakeReturns.result1
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePathCallCount() int {
	fake.removeReleaseUpgradePathMutex.RLock()
	defer fake.removeReleaseUpgradePathMutex.RUnlock()
	return len(fake.removeReleaseUpgradePathArgsForCall)
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePathCalls(stub func(string, int, int) error) {
	fake.removeReleaseUpgradePathMutex.Lock()
	defer fake.removeReleaseUpgradePathMutex.Unlock()
	fake.RemoveReleaseUpgradePathStub = stub
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePathArgsForCall(i int) (string, int, int) {
	fake.removeReleaseUpgradePathMutex.RLock()
	defer fake.removeReleaseUpgradePathMutex.RUnlock()
	argsForCall := fake.removeReleaseUpgradePathArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePathReturns(result1 error) {
	fake.removeReleaseUpgradePathMutex.Lock()
	defer fake.removeReleaseUpgradePathMutex.Unlock()
	fake.RemoveReleaseUpgradePathStub = nil
	fake.removeReleaseUpgradePathReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) RemoveReleaseUpgradePathReturnsOnCall(i int, result1 error) {
	fake.removeReleaseUpgradePathMutex.Lock()
	defer fake.removeReleaseUpgradePathMutex.Unlock()
	fake.RemoveReleaseUpgradePathStub = nil
	if fake.removeReleaseUpgradePathReturnsOnCall == nil {
		fake.removeReleaseUpgradePathReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeReleaseUpgradePathReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addReleaseUpgradePathMutex.RLock()
	defer fake.addReleaseUpgradePathMutex.RUnlock()
	fake.releaseForVersionMutex.RLock()
	defer fake.releaseForVersionMutex.RUnlock()
	fake.releaseUpgradePathsMutex.RLock()
	defer fake.releaseUpgradePathsMutex.RUnlock()
	fake.releasesForProductSlugMutex.RLock()
	defer fake.releasesForProductSlugMutex.RUnlock()
	fake.removeReleaseUpgradePathMutex.RLock()
	defer fake.removeReleaseUpgradePathMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePivnetClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ releaseupgradepath.PivnetClient = new(FakePivnetClient)
