// Code generated by counterfeiter. DO NOT EDIT.
package subscriptiongroupfakes

import (
	"sync"

	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/subscriptiongroup"
)

type FakePivnetClient struct {
	AddSubscriptionGroupMemberStub        func(int, string, string) (pivnet.SubscriptionGroup, error)
	addSubscriptionGroupMemberMutex       sync.RWMutex
	addSubscriptionGroupMemberArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 string
	}
	addSubscriptionGroupMemberReturns struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	addSubscriptionGroupMemberReturnsOnCall map[int]struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	RemoveSubscriptionGroupMemberStub        func(int, string) (pivnet.SubscriptionGroup, error)
	removeSubscriptionGroupMemberMutex       sync.RWMutex
	removeSubscriptionGroupMemberArgsForCall []struct {
		arg1 int
		arg2 string
	}
	removeSubscriptionGroupMemberReturns struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	removeSubscriptionGroupMemberReturnsOnCall map[int]struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	SubscriptionGroupStub        func(int) (pivnet.SubscriptionGroup, error)
	subscriptionGroupMutex       sync.RWMutex
	subscriptionGroupArgsForCall []struct {
		arg1 int
	}
	subscriptionGroupReturns struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	subscriptionGroupReturnsOnCall map[int]struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}
	SubscriptionGroupsStub        func() ([]pivnet.SubscriptionGroup, error)
	subscriptionGroupsMutex       sync.RWMutex
	subscriptionGroupsArgsForCall []struct {
	}
	subscriptionGroupsReturns struct {
		result1 []pivnet.SubscriptionGroup
		result2 error
	}
	subscriptionGroupsReturnsOnCall map[int]struct {
		result1 []pivnet.SubscriptionGroup
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePivnetClient) AddSubscriptionGroupMember(arg1 int, arg2 string, arg3 string) (pivnet.SubscriptionGroup, error) {
	fake.addSubscriptionGroupMemberMutex.Lock()
	ret, specificReturn := fake.addSubscriptionGroupMemberReturnsOnCall[len(fake.addSubscriptionGroupMemberArgsForCall)]
	fake.addSubscriptionGroupMemberArgsForCall = append(fake.addSubscriptionGroupMemberArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("AddSubscriptionGroupMember", []interface{}{arg1, arg2, arg3})
	fake.addSubscriptionGroupMemberMutex.Unlock()
	if fake.AddSubscriptionGroupMemberStub != nil {
		return fake.AddSubscriptionGroupMemberStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.addSubscriptionGroupMemberReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) AddSubscriptionGroupMemberCallCount() int {
	fake.addSubscriptionGroupMemberMutex.RLock()
	defer fake.addSubscriptionGroupMemberMutex.RUnlock()
	return len(fake.addSubscriptionGroupMemberArgsForCall)
}

func (fake *FakePivnetClient) AddSubscriptionGroupMemberCalls(stub func(int, string, string) (pivnet.SubscriptionGroup, error)) {
	fake.addSubscriptionGroupMemberMutex.Lock()
	defer fake.addSubscriptionGroupMemberMutex.Unlock()
	fake.AddSubscriptionGroupMemberStub = stub
}

func (fake *FakePivnetClient) AddSubscriptionGroupMemberArgsForCall(i int) (int, string, string) {
	fake.addSubscriptionGroupMemberMutex.RLock()
	defer fake.addSubscriptionGroupMemberMutex.RUnlock()
	argsForCall := fake.addSubscriptionGroupMemberArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) AddSubscriptionGroupMemberReturns(result1 pivnet.SubscriptionGroup, result2 error) {
	fake.addSubscriptionGroupMemberMutex.Lock()
	defer fake.addSubscriptionGroupMemberMutex.Unlock()
	fake.AddSubscriptionGroupMemberStub = nil
	fake.addSubscriptionGroupMemberReturns = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) AddSubscriptionGroupMemberReturnsOnCall(i int, result1 pivnet.SubscriptionGroup, result2 error) {
	fake.addSubscriptionGroupMemberMutex.Lock()
	defer fake.addSubscriptionGroupMemberMutex.Unlock()
	fake.AddSubscriptionGroupMemberStub = nil
	if fake.addSubscriptionGroupMemberReturnsOnCall == nil {
		fake.addSubscriptionGroupMemberReturnsOnCall = make(map[int]struct {
			result1 pivnet.SubscriptionGroup
			result2 error
		})
	}
	fake.addSubscriptionGroupMemberReturnsOnCall[i] = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMember(arg1 int, arg2 string) (pivnet.SubscriptionGroup, error) {
	fake.removeSubscriptionGroupMemberMutex.Lock()
	ret, specificReturn := fake.removeSubscriptionGroupMemberReturnsOnCall[len(fake.removeSubscriptionGroupMemberArgsForCall)]
	fake.removeSubscriptionGroupMemberArgsForCall = append(fake.removeSubscriptionGroupMemberArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("RemoveSubscriptionGroupMember", []interface{}{arg1, arg2})
	fake.removeSubscriptionGroupMemberMutex.Unlock()
	if fake.RemoveSubscriptionGroupMemberStub != nil {
		return fake.RemoveSubscriptionGroupMemberStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.removeSubscriptionGroupMemberReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMemberCallCount() int {
	fake.removeSubscriptionGroupMemberMutex.RLock()
	defer fake.removeSubscriptionGroupMemberMutex.RUnlock()
	return len(fake.removeSubscriptionGroupMemberArgsForCall)
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMemberCalls(stub func(int, string) (pivnet.SubscriptionGroup, error)) {
	fake.removeSubscriptionGroupMemberMutex.Lock()
	defer fake.removeSubscriptionGroupMemberMutex.Unlock()
	fake.RemoveSubscriptionGroupMemberStub = stub
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMemberArgsForCall(i int) (int, string) {
	fake.removeSubscriptionGroupMemberMutex.RLock()
	defer fake.removeSubscriptionGroupMemberMutex.RUnlock()
	argsForCall := fake.removeSubscriptionGroupMemberArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMemberReturns(result1 pivnet.SubscriptionGroup, result2 error) {
	fake.removeSubscriptionGroupMemberMutex.Lock()
	defer fake.removeSubscriptionGroupMemberMutex.Unlock()
	fake.RemoveSubscriptionGroupMemberStub = nil
	fake.removeSubscriptionGroupMemberReturns = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) RemoveSubscriptionGroupMemberReturnsOnCall(i int, result1 pivnet.SubscriptionGroup, result2 error) {
	fake.removeSubscriptionGroupMemberMutex.Lock()
	defer fake.removeSubscriptionGroupMemberMutex.Unlock()
	fake.RemoveSubscriptionGroupMemberStub = nil
	if fake.removeSubscriptionGroupMemberReturnsOnCall == nil {
		fake.removeSubscriptionGroupMemberReturnsOnCall = make(map[int]struct {
			result1 pivnet.SubscriptionGroup
			result2 error
		})
	}
	fake.removeSubscriptionGroupMemberReturnsOnCall[i] = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) SubscriptionGroup(arg1 int) (pivnet.SubscriptionGroup, error) {
	fake.subscriptionGroupMutex.Lock()
	ret, specificReturn := fake.subscriptionGroupReturnsOnCall[len(fake.subscriptionGroupArgsForCall)]
	fake.subscriptionGroupArgsForCall = append(fake.subscriptionGroupArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("SubscriptionGroup", []interface{}{arg1})
	fake.subscriptionGroupMutex.Unlock()
	if fake.SubscriptionGroupStub != nil {
		return fake.SubscriptionGroupStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.subscriptionGroupReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) SubscriptionGroupCallCount() int {
	fake.subscriptionGroupMutex.RLock()
	defer fake.subscriptionGroupMutex.RUnlock()
	return len(fake.subscriptionGroupArgsForCall)
}

func (fake *FakePivnetClient) SubscriptionGroupCalls(stub func(int) (pivnet.SubscriptionGroup, error)) {
	fake.subscriptionGroupMutex.Lock()
	defer fake.subscriptionGroupMutex.Unlock()
	fake.SubscriptionGroupStub = stub
}

func (fake *FakePivnetClient) SubscriptionGroupArgsForCall(i int) int {
	fake.subscriptionGroupMutex.RLock()
	defer fake.subscriptionGroupMutex.RUnlock()
	argsForCall := fake.subscriptionGroupArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePivnetClient) SubscriptionGroupReturns(result1 pivnet.SubscriptionGroup, result2 error) {
	fake.subscriptionGroupMutex.Lock()
	defer fake.subscriptionGroupMutex.Unlock()
	fake.SubscriptionGroupStub = nil
	fake.subscriptionGroupReturns = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) SubscriptionGroupReturnsOnCall(i int, result1 pivnet.SubscriptionGroup, result2 error) {
	fake.subscriptionGroupMutex.Lock()
	defer fake.subscriptionGroupMutex.Unlock()
	fake.SubscriptionGroupStub = nil
	if fake.subscriptionGroupReturnsOnCall == nil {
		fake.subscriptionGroupReturnsOnCall = make(map[int]struct {
			result1 pivnet.SubscriptionGroup
			result2 error
		})
	}
	fake.subscriptionGroupReturnsOnCall[i] = struct {
		result1 pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) SubscriptionGroups() ([]pivnet.SubscriptionGroup, error) {
	fake.subscriptionGroupsMutex.Lock()
	ret, specificReturn := fake.subscriptionGroupsReturnsOnCall[len(fake.subscriptionGroupsArgsForCall)]
	fake.subscriptionGroupsArgsForCall = append(fake.subscriptionGroupsArgsForCall, struct {
	}{})
	fake.recordInvocation("SubscriptionGroups", []interface{}{})
	fake.subscriptionGroupsMutex.Unlock()
	if fake.SubscriptionGroupsStub != nil {
		return fake.SubscriptionGroupsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.subscriptionGroupsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) SubscriptionGroupsCallCount() int {
	fake.subscriptionGroupsMutex.RLock()
	defer fake.subscriptionGroupsMutex.RUnlock()
	return len(fake.subscriptionGroupsArgsForCall)
}

func (fake *FakePivnetClient) SubscriptionGroupsCalls(stub func() ([]pivnet.SubscriptionGroup, error)) {
	fake.subscriptionGroupsMutex.Lock()
	defer fake.subscriptionGroupsMutex.Unlock()
	fake.SubscriptionGroupsStub = stub
}

func (fake *FakePivnetClient) SubscriptionGroupsReturns(result1 []pivnet.SubscriptionGroup, result2 error) {
	fake.subscriptionGroupsMutex.Lock()
	defer fake.subscriptionGroupsMutex.Unlock()
	fake.SubscriptionGroupsStub = nil
	fake.subscriptionGroupsReturns = struct {
		result1 []pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) SubscriptionGroupsReturnsOnCall(i int, result1 []pivnet.SubscriptionGroup, result2 error) {
	fake.subscriptionGroupsMutex.Lock()
	defer fake.subscriptionGroupsMutex.Unlock()
	fake.SubscriptionGroupsStub = nil
	if fake.subscriptionGroupsReturnsOnCall == nil {
		fake.subscriptionGroupsReturnsOnCall = make(map[int]struct {
			result1 []pivnet.SubscriptionGroup
			result2 error
		})
	}
	fake.subscriptionGroupsReturnsOnCall[i] = struct {
		result1 []pivnet.SubscriptionGroup
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addSubscriptionGroupMemberMutex.RLock()
	defer fake.addSubscriptionGroupMemberMutex.RUnlock()
	fake.removeSubscriptionGroupMemberMutex.RLock()
	defer fake.removeSubscriptionGroupMemberMutex.RUnlock()
	fake.subscriptionGroupMutex.RLock()
	defer fake.subscriptionGroupMutex.RUnlock()
	fake.subscriptionGroupsMutex.RLock()
	defer fake.subscriptionGroupsMutex.RUnlock()
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

var _ subscriptiongroup.PivnetClient = new(FakePivnetClient)
