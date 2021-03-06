// Code generated by counterfeiter. DO NOT EDIT.
package artifactreferencefakes

import (
	"sync"

	pivnet "github.com/pivotal-cf/go-pivnet/v7"
	"github.com/pivotal-cf/pivnet-cli/v3/commands/artifactreference"
)

type FakePivnetClient struct {
	AddArtifactReferenceToReleaseStub        func(string, int, int) error
	addArtifactReferenceToReleaseMutex       sync.RWMutex
	addArtifactReferenceToReleaseArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
	}
	addArtifactReferenceToReleaseReturns struct {
		result1 error
	}
	addArtifactReferenceToReleaseReturnsOnCall map[int]struct {
		result1 error
	}
	ArtifactReferenceStub        func(string, int) (pivnet.ArtifactReference, error)
	artifactReferenceMutex       sync.RWMutex
	artifactReferenceArgsForCall []struct {
		arg1 string
		arg2 int
	}
	artifactReferenceReturns struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	artifactReferenceReturnsOnCall map[int]struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	ArtifactReferenceForReleaseStub        func(string, int, int) (pivnet.ArtifactReference, error)
	artifactReferenceForReleaseMutex       sync.RWMutex
	artifactReferenceForReleaseArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
	}
	artifactReferenceForReleaseReturns struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	artifactReferenceForReleaseReturnsOnCall map[int]struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	ArtifactReferencesStub        func(string) ([]pivnet.ArtifactReference, error)
	artifactReferencesMutex       sync.RWMutex
	artifactReferencesArgsForCall []struct {
		arg1 string
	}
	artifactReferencesReturns struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	artifactReferencesReturnsOnCall map[int]struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	ArtifactReferencesForDigestStub        func(string, string) ([]pivnet.ArtifactReference, error)
	artifactReferencesForDigestMutex       sync.RWMutex
	artifactReferencesForDigestArgsForCall []struct {
		arg1 string
		arg2 string
	}
	artifactReferencesForDigestReturns struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	artifactReferencesForDigestReturnsOnCall map[int]struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	ArtifactReferencesForReleaseStub        func(string, int) ([]pivnet.ArtifactReference, error)
	artifactReferencesForReleaseMutex       sync.RWMutex
	artifactReferencesForReleaseArgsForCall []struct {
		arg1 string
		arg2 int
	}
	artifactReferencesForReleaseReturns struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	artifactReferencesForReleaseReturnsOnCall map[int]struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}
	CreateArtifactReferenceStub        func(pivnet.CreateArtifactReferenceConfig) (pivnet.ArtifactReference, error)
	createArtifactReferenceMutex       sync.RWMutex
	createArtifactReferenceArgsForCall []struct {
		arg1 pivnet.CreateArtifactReferenceConfig
	}
	createArtifactReferenceReturns struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	createArtifactReferenceReturnsOnCall map[int]struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	DeleteArtifactReferenceStub        func(string, int) (pivnet.ArtifactReference, error)
	deleteArtifactReferenceMutex       sync.RWMutex
	deleteArtifactReferenceArgsForCall []struct {
		arg1 string
		arg2 int
	}
	deleteArtifactReferenceReturns struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	deleteArtifactReferenceReturnsOnCall map[int]struct {
		result1 pivnet.ArtifactReference
		result2 error
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
	RemoveArtifactReferenceFromReleaseStub        func(string, int, int) error
	removeArtifactReferenceFromReleaseMutex       sync.RWMutex
	removeArtifactReferenceFromReleaseArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
	}
	removeArtifactReferenceFromReleaseReturns struct {
		result1 error
	}
	removeArtifactReferenceFromReleaseReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateArtifactReferenceStub        func(string, pivnet.ArtifactReference) (pivnet.ArtifactReference, error)
	updateArtifactReferenceMutex       sync.RWMutex
	updateArtifactReferenceArgsForCall []struct {
		arg1 string
		arg2 pivnet.ArtifactReference
	}
	updateArtifactReferenceReturns struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	updateArtifactReferenceReturnsOnCall map[int]struct {
		result1 pivnet.ArtifactReference
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePivnetClient) AddArtifactReferenceToRelease(arg1 string, arg2 int, arg3 int) error {
	fake.addArtifactReferenceToReleaseMutex.Lock()
	ret, specificReturn := fake.addArtifactReferenceToReleaseReturnsOnCall[len(fake.addArtifactReferenceToReleaseArgsForCall)]
	fake.addArtifactReferenceToReleaseArgsForCall = append(fake.addArtifactReferenceToReleaseArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	fake.recordInvocation("AddArtifactReferenceToRelease", []interface{}{arg1, arg2, arg3})
	fake.addArtifactReferenceToReleaseMutex.Unlock()
	if fake.AddArtifactReferenceToReleaseStub != nil {
		return fake.AddArtifactReferenceToReleaseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addArtifactReferenceToReleaseReturns
	return fakeReturns.result1
}

func (fake *FakePivnetClient) AddArtifactReferenceToReleaseCallCount() int {
	fake.addArtifactReferenceToReleaseMutex.RLock()
	defer fake.addArtifactReferenceToReleaseMutex.RUnlock()
	return len(fake.addArtifactReferenceToReleaseArgsForCall)
}

func (fake *FakePivnetClient) AddArtifactReferenceToReleaseCalls(stub func(string, int, int) error) {
	fake.addArtifactReferenceToReleaseMutex.Lock()
	defer fake.addArtifactReferenceToReleaseMutex.Unlock()
	fake.AddArtifactReferenceToReleaseStub = stub
}

func (fake *FakePivnetClient) AddArtifactReferenceToReleaseArgsForCall(i int) (string, int, int) {
	fake.addArtifactReferenceToReleaseMutex.RLock()
	defer fake.addArtifactReferenceToReleaseMutex.RUnlock()
	argsForCall := fake.addArtifactReferenceToReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) AddArtifactReferenceToReleaseReturns(result1 error) {
	fake.addArtifactReferenceToReleaseMutex.Lock()
	defer fake.addArtifactReferenceToReleaseMutex.Unlock()
	fake.AddArtifactReferenceToReleaseStub = nil
	fake.addArtifactReferenceToReleaseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) AddArtifactReferenceToReleaseReturnsOnCall(i int, result1 error) {
	fake.addArtifactReferenceToReleaseMutex.Lock()
	defer fake.addArtifactReferenceToReleaseMutex.Unlock()
	fake.AddArtifactReferenceToReleaseStub = nil
	if fake.addArtifactReferenceToReleaseReturnsOnCall == nil {
		fake.addArtifactReferenceToReleaseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addArtifactReferenceToReleaseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) ArtifactReference(arg1 string, arg2 int) (pivnet.ArtifactReference, error) {
	fake.artifactReferenceMutex.Lock()
	ret, specificReturn := fake.artifactReferenceReturnsOnCall[len(fake.artifactReferenceArgsForCall)]
	fake.artifactReferenceArgsForCall = append(fake.artifactReferenceArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("ArtifactReference", []interface{}{arg1, arg2})
	fake.artifactReferenceMutex.Unlock()
	if fake.ArtifactReferenceStub != nil {
		return fake.ArtifactReferenceStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.artifactReferenceReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ArtifactReferenceCallCount() int {
	fake.artifactReferenceMutex.RLock()
	defer fake.artifactReferenceMutex.RUnlock()
	return len(fake.artifactReferenceArgsForCall)
}

func (fake *FakePivnetClient) ArtifactReferenceCalls(stub func(string, int) (pivnet.ArtifactReference, error)) {
	fake.artifactReferenceMutex.Lock()
	defer fake.artifactReferenceMutex.Unlock()
	fake.ArtifactReferenceStub = stub
}

func (fake *FakePivnetClient) ArtifactReferenceArgsForCall(i int) (string, int) {
	fake.artifactReferenceMutex.RLock()
	defer fake.artifactReferenceMutex.RUnlock()
	argsForCall := fake.artifactReferenceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ArtifactReferenceReturns(result1 pivnet.ArtifactReference, result2 error) {
	fake.artifactReferenceMutex.Lock()
	defer fake.artifactReferenceMutex.Unlock()
	fake.ArtifactReferenceStub = nil
	fake.artifactReferenceReturns = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferenceReturnsOnCall(i int, result1 pivnet.ArtifactReference, result2 error) {
	fake.artifactReferenceMutex.Lock()
	defer fake.artifactReferenceMutex.Unlock()
	fake.ArtifactReferenceStub = nil
	if fake.artifactReferenceReturnsOnCall == nil {
		fake.artifactReferenceReturnsOnCall = make(map[int]struct {
			result1 pivnet.ArtifactReference
			result2 error
		})
	}
	fake.artifactReferenceReturnsOnCall[i] = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferenceForRelease(arg1 string, arg2 int, arg3 int) (pivnet.ArtifactReference, error) {
	fake.artifactReferenceForReleaseMutex.Lock()
	ret, specificReturn := fake.artifactReferenceForReleaseReturnsOnCall[len(fake.artifactReferenceForReleaseArgsForCall)]
	fake.artifactReferenceForReleaseArgsForCall = append(fake.artifactReferenceForReleaseArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	fake.recordInvocation("ArtifactReferenceForRelease", []interface{}{arg1, arg2, arg3})
	fake.artifactReferenceForReleaseMutex.Unlock()
	if fake.ArtifactReferenceForReleaseStub != nil {
		return fake.ArtifactReferenceForReleaseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.artifactReferenceForReleaseReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ArtifactReferenceForReleaseCallCount() int {
	fake.artifactReferenceForReleaseMutex.RLock()
	defer fake.artifactReferenceForReleaseMutex.RUnlock()
	return len(fake.artifactReferenceForReleaseArgsForCall)
}

func (fake *FakePivnetClient) ArtifactReferenceForReleaseCalls(stub func(string, int, int) (pivnet.ArtifactReference, error)) {
	fake.artifactReferenceForReleaseMutex.Lock()
	defer fake.artifactReferenceForReleaseMutex.Unlock()
	fake.ArtifactReferenceForReleaseStub = stub
}

func (fake *FakePivnetClient) ArtifactReferenceForReleaseArgsForCall(i int) (string, int, int) {
	fake.artifactReferenceForReleaseMutex.RLock()
	defer fake.artifactReferenceForReleaseMutex.RUnlock()
	argsForCall := fake.artifactReferenceForReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) ArtifactReferenceForReleaseReturns(result1 pivnet.ArtifactReference, result2 error) {
	fake.artifactReferenceForReleaseMutex.Lock()
	defer fake.artifactReferenceForReleaseMutex.Unlock()
	fake.ArtifactReferenceForReleaseStub = nil
	fake.artifactReferenceForReleaseReturns = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferenceForReleaseReturnsOnCall(i int, result1 pivnet.ArtifactReference, result2 error) {
	fake.artifactReferenceForReleaseMutex.Lock()
	defer fake.artifactReferenceForReleaseMutex.Unlock()
	fake.ArtifactReferenceForReleaseStub = nil
	if fake.artifactReferenceForReleaseReturnsOnCall == nil {
		fake.artifactReferenceForReleaseReturnsOnCall = make(map[int]struct {
			result1 pivnet.ArtifactReference
			result2 error
		})
	}
	fake.artifactReferenceForReleaseReturnsOnCall[i] = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferences(arg1 string) ([]pivnet.ArtifactReference, error) {
	fake.artifactReferencesMutex.Lock()
	ret, specificReturn := fake.artifactReferencesReturnsOnCall[len(fake.artifactReferencesArgsForCall)]
	fake.artifactReferencesArgsForCall = append(fake.artifactReferencesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ArtifactReferences", []interface{}{arg1})
	fake.artifactReferencesMutex.Unlock()
	if fake.ArtifactReferencesStub != nil {
		return fake.ArtifactReferencesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.artifactReferencesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ArtifactReferencesCallCount() int {
	fake.artifactReferencesMutex.RLock()
	defer fake.artifactReferencesMutex.RUnlock()
	return len(fake.artifactReferencesArgsForCall)
}

func (fake *FakePivnetClient) ArtifactReferencesCalls(stub func(string) ([]pivnet.ArtifactReference, error)) {
	fake.artifactReferencesMutex.Lock()
	defer fake.artifactReferencesMutex.Unlock()
	fake.ArtifactReferencesStub = stub
}

func (fake *FakePivnetClient) ArtifactReferencesArgsForCall(i int) string {
	fake.artifactReferencesMutex.RLock()
	defer fake.artifactReferencesMutex.RUnlock()
	argsForCall := fake.artifactReferencesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePivnetClient) ArtifactReferencesReturns(result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesMutex.Lock()
	defer fake.artifactReferencesMutex.Unlock()
	fake.ArtifactReferencesStub = nil
	fake.artifactReferencesReturns = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferencesReturnsOnCall(i int, result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesMutex.Lock()
	defer fake.artifactReferencesMutex.Unlock()
	fake.ArtifactReferencesStub = nil
	if fake.artifactReferencesReturnsOnCall == nil {
		fake.artifactReferencesReturnsOnCall = make(map[int]struct {
			result1 []pivnet.ArtifactReference
			result2 error
		})
	}
	fake.artifactReferencesReturnsOnCall[i] = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferencesForDigest(arg1 string, arg2 string) ([]pivnet.ArtifactReference, error) {
	fake.artifactReferencesForDigestMutex.Lock()
	ret, specificReturn := fake.artifactReferencesForDigestReturnsOnCall[len(fake.artifactReferencesForDigestArgsForCall)]
	fake.artifactReferencesForDigestArgsForCall = append(fake.artifactReferencesForDigestArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ArtifactReferencesForDigest", []interface{}{arg1, arg2})
	fake.artifactReferencesForDigestMutex.Unlock()
	if fake.ArtifactReferencesForDigestStub != nil {
		return fake.ArtifactReferencesForDigestStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.artifactReferencesForDigestReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ArtifactReferencesForDigestCallCount() int {
	fake.artifactReferencesForDigestMutex.RLock()
	defer fake.artifactReferencesForDigestMutex.RUnlock()
	return len(fake.artifactReferencesForDigestArgsForCall)
}

func (fake *FakePivnetClient) ArtifactReferencesForDigestCalls(stub func(string, string) ([]pivnet.ArtifactReference, error)) {
	fake.artifactReferencesForDigestMutex.Lock()
	defer fake.artifactReferencesForDigestMutex.Unlock()
	fake.ArtifactReferencesForDigestStub = stub
}

func (fake *FakePivnetClient) ArtifactReferencesForDigestArgsForCall(i int) (string, string) {
	fake.artifactReferencesForDigestMutex.RLock()
	defer fake.artifactReferencesForDigestMutex.RUnlock()
	argsForCall := fake.artifactReferencesForDigestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ArtifactReferencesForDigestReturns(result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesForDigestMutex.Lock()
	defer fake.artifactReferencesForDigestMutex.Unlock()
	fake.ArtifactReferencesForDigestStub = nil
	fake.artifactReferencesForDigestReturns = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferencesForDigestReturnsOnCall(i int, result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesForDigestMutex.Lock()
	defer fake.artifactReferencesForDigestMutex.Unlock()
	fake.ArtifactReferencesForDigestStub = nil
	if fake.artifactReferencesForDigestReturnsOnCall == nil {
		fake.artifactReferencesForDigestReturnsOnCall = make(map[int]struct {
			result1 []pivnet.ArtifactReference
			result2 error
		})
	}
	fake.artifactReferencesForDigestReturnsOnCall[i] = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferencesForRelease(arg1 string, arg2 int) ([]pivnet.ArtifactReference, error) {
	fake.artifactReferencesForReleaseMutex.Lock()
	ret, specificReturn := fake.artifactReferencesForReleaseReturnsOnCall[len(fake.artifactReferencesForReleaseArgsForCall)]
	fake.artifactReferencesForReleaseArgsForCall = append(fake.artifactReferencesForReleaseArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("ArtifactReferencesForRelease", []interface{}{arg1, arg2})
	fake.artifactReferencesForReleaseMutex.Unlock()
	if fake.ArtifactReferencesForReleaseStub != nil {
		return fake.ArtifactReferencesForReleaseStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.artifactReferencesForReleaseReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) ArtifactReferencesForReleaseCallCount() int {
	fake.artifactReferencesForReleaseMutex.RLock()
	defer fake.artifactReferencesForReleaseMutex.RUnlock()
	return len(fake.artifactReferencesForReleaseArgsForCall)
}

func (fake *FakePivnetClient) ArtifactReferencesForReleaseCalls(stub func(string, int) ([]pivnet.ArtifactReference, error)) {
	fake.artifactReferencesForReleaseMutex.Lock()
	defer fake.artifactReferencesForReleaseMutex.Unlock()
	fake.ArtifactReferencesForReleaseStub = stub
}

func (fake *FakePivnetClient) ArtifactReferencesForReleaseArgsForCall(i int) (string, int) {
	fake.artifactReferencesForReleaseMutex.RLock()
	defer fake.artifactReferencesForReleaseMutex.RUnlock()
	argsForCall := fake.artifactReferencesForReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) ArtifactReferencesForReleaseReturns(result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesForReleaseMutex.Lock()
	defer fake.artifactReferencesForReleaseMutex.Unlock()
	fake.ArtifactReferencesForReleaseStub = nil
	fake.artifactReferencesForReleaseReturns = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) ArtifactReferencesForReleaseReturnsOnCall(i int, result1 []pivnet.ArtifactReference, result2 error) {
	fake.artifactReferencesForReleaseMutex.Lock()
	defer fake.artifactReferencesForReleaseMutex.Unlock()
	fake.ArtifactReferencesForReleaseStub = nil
	if fake.artifactReferencesForReleaseReturnsOnCall == nil {
		fake.artifactReferencesForReleaseReturnsOnCall = make(map[int]struct {
			result1 []pivnet.ArtifactReference
			result2 error
		})
	}
	fake.artifactReferencesForReleaseReturnsOnCall[i] = struct {
		result1 []pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) CreateArtifactReference(arg1 pivnet.CreateArtifactReferenceConfig) (pivnet.ArtifactReference, error) {
	fake.createArtifactReferenceMutex.Lock()
	ret, specificReturn := fake.createArtifactReferenceReturnsOnCall[len(fake.createArtifactReferenceArgsForCall)]
	fake.createArtifactReferenceArgsForCall = append(fake.createArtifactReferenceArgsForCall, struct {
		arg1 pivnet.CreateArtifactReferenceConfig
	}{arg1})
	fake.recordInvocation("CreateArtifactReference", []interface{}{arg1})
	fake.createArtifactReferenceMutex.Unlock()
	if fake.CreateArtifactReferenceStub != nil {
		return fake.CreateArtifactReferenceStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createArtifactReferenceReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) CreateArtifactReferenceCallCount() int {
	fake.createArtifactReferenceMutex.RLock()
	defer fake.createArtifactReferenceMutex.RUnlock()
	return len(fake.createArtifactReferenceArgsForCall)
}

func (fake *FakePivnetClient) CreateArtifactReferenceCalls(stub func(pivnet.CreateArtifactReferenceConfig) (pivnet.ArtifactReference, error)) {
	fake.createArtifactReferenceMutex.Lock()
	defer fake.createArtifactReferenceMutex.Unlock()
	fake.CreateArtifactReferenceStub = stub
}

func (fake *FakePivnetClient) CreateArtifactReferenceArgsForCall(i int) pivnet.CreateArtifactReferenceConfig {
	fake.createArtifactReferenceMutex.RLock()
	defer fake.createArtifactReferenceMutex.RUnlock()
	argsForCall := fake.createArtifactReferenceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePivnetClient) CreateArtifactReferenceReturns(result1 pivnet.ArtifactReference, result2 error) {
	fake.createArtifactReferenceMutex.Lock()
	defer fake.createArtifactReferenceMutex.Unlock()
	fake.CreateArtifactReferenceStub = nil
	fake.createArtifactReferenceReturns = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) CreateArtifactReferenceReturnsOnCall(i int, result1 pivnet.ArtifactReference, result2 error) {
	fake.createArtifactReferenceMutex.Lock()
	defer fake.createArtifactReferenceMutex.Unlock()
	fake.CreateArtifactReferenceStub = nil
	if fake.createArtifactReferenceReturnsOnCall == nil {
		fake.createArtifactReferenceReturnsOnCall = make(map[int]struct {
			result1 pivnet.ArtifactReference
			result2 error
		})
	}
	fake.createArtifactReferenceReturnsOnCall[i] = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) DeleteArtifactReference(arg1 string, arg2 int) (pivnet.ArtifactReference, error) {
	fake.deleteArtifactReferenceMutex.Lock()
	ret, specificReturn := fake.deleteArtifactReferenceReturnsOnCall[len(fake.deleteArtifactReferenceArgsForCall)]
	fake.deleteArtifactReferenceArgsForCall = append(fake.deleteArtifactReferenceArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("DeleteArtifactReference", []interface{}{arg1, arg2})
	fake.deleteArtifactReferenceMutex.Unlock()
	if fake.DeleteArtifactReferenceStub != nil {
		return fake.DeleteArtifactReferenceStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteArtifactReferenceReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) DeleteArtifactReferenceCallCount() int {
	fake.deleteArtifactReferenceMutex.RLock()
	defer fake.deleteArtifactReferenceMutex.RUnlock()
	return len(fake.deleteArtifactReferenceArgsForCall)
}

func (fake *FakePivnetClient) DeleteArtifactReferenceCalls(stub func(string, int) (pivnet.ArtifactReference, error)) {
	fake.deleteArtifactReferenceMutex.Lock()
	defer fake.deleteArtifactReferenceMutex.Unlock()
	fake.DeleteArtifactReferenceStub = stub
}

func (fake *FakePivnetClient) DeleteArtifactReferenceArgsForCall(i int) (string, int) {
	fake.deleteArtifactReferenceMutex.RLock()
	defer fake.deleteArtifactReferenceMutex.RUnlock()
	argsForCall := fake.deleteArtifactReferenceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) DeleteArtifactReferenceReturns(result1 pivnet.ArtifactReference, result2 error) {
	fake.deleteArtifactReferenceMutex.Lock()
	defer fake.deleteArtifactReferenceMutex.Unlock()
	fake.DeleteArtifactReferenceStub = nil
	fake.deleteArtifactReferenceReturns = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) DeleteArtifactReferenceReturnsOnCall(i int, result1 pivnet.ArtifactReference, result2 error) {
	fake.deleteArtifactReferenceMutex.Lock()
	defer fake.deleteArtifactReferenceMutex.Unlock()
	fake.DeleteArtifactReferenceStub = nil
	if fake.deleteArtifactReferenceReturnsOnCall == nil {
		fake.deleteArtifactReferenceReturnsOnCall = make(map[int]struct {
			result1 pivnet.ArtifactReference
			result2 error
		})
	}
	fake.deleteArtifactReferenceReturnsOnCall[i] = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
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

func (fake *FakePivnetClient) RemoveArtifactReferenceFromRelease(arg1 string, arg2 int, arg3 int) error {
	fake.removeArtifactReferenceFromReleaseMutex.Lock()
	ret, specificReturn := fake.removeArtifactReferenceFromReleaseReturnsOnCall[len(fake.removeArtifactReferenceFromReleaseArgsForCall)]
	fake.removeArtifactReferenceFromReleaseArgsForCall = append(fake.removeArtifactReferenceFromReleaseArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	fake.recordInvocation("RemoveArtifactReferenceFromRelease", []interface{}{arg1, arg2, arg3})
	fake.removeArtifactReferenceFromReleaseMutex.Unlock()
	if fake.RemoveArtifactReferenceFromReleaseStub != nil {
		return fake.RemoveArtifactReferenceFromReleaseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeArtifactReferenceFromReleaseReturns
	return fakeReturns.result1
}

func (fake *FakePivnetClient) RemoveArtifactReferenceFromReleaseCallCount() int {
	fake.removeArtifactReferenceFromReleaseMutex.RLock()
	defer fake.removeArtifactReferenceFromReleaseMutex.RUnlock()
	return len(fake.removeArtifactReferenceFromReleaseArgsForCall)
}

func (fake *FakePivnetClient) RemoveArtifactReferenceFromReleaseCalls(stub func(string, int, int) error) {
	fake.removeArtifactReferenceFromReleaseMutex.Lock()
	defer fake.removeArtifactReferenceFromReleaseMutex.Unlock()
	fake.RemoveArtifactReferenceFromReleaseStub = stub
}

func (fake *FakePivnetClient) RemoveArtifactReferenceFromReleaseArgsForCall(i int) (string, int, int) {
	fake.removeArtifactReferenceFromReleaseMutex.RLock()
	defer fake.removeArtifactReferenceFromReleaseMutex.RUnlock()
	argsForCall := fake.removeArtifactReferenceFromReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePivnetClient) RemoveArtifactReferenceFromReleaseReturns(result1 error) {
	fake.removeArtifactReferenceFromReleaseMutex.Lock()
	defer fake.removeArtifactReferenceFromReleaseMutex.Unlock()
	fake.RemoveArtifactReferenceFromReleaseStub = nil
	fake.removeArtifactReferenceFromReleaseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) RemoveArtifactReferenceFromReleaseReturnsOnCall(i int, result1 error) {
	fake.removeArtifactReferenceFromReleaseMutex.Lock()
	defer fake.removeArtifactReferenceFromReleaseMutex.Unlock()
	fake.RemoveArtifactReferenceFromReleaseStub = nil
	if fake.removeArtifactReferenceFromReleaseReturnsOnCall == nil {
		fake.removeArtifactReferenceFromReleaseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeArtifactReferenceFromReleaseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePivnetClient) UpdateArtifactReference(arg1 string, arg2 pivnet.ArtifactReference) (pivnet.ArtifactReference, error) {
	fake.updateArtifactReferenceMutex.Lock()
	ret, specificReturn := fake.updateArtifactReferenceReturnsOnCall[len(fake.updateArtifactReferenceArgsForCall)]
	fake.updateArtifactReferenceArgsForCall = append(fake.updateArtifactReferenceArgsForCall, struct {
		arg1 string
		arg2 pivnet.ArtifactReference
	}{arg1, arg2})
	fake.recordInvocation("UpdateArtifactReference", []interface{}{arg1, arg2})
	fake.updateArtifactReferenceMutex.Unlock()
	if fake.UpdateArtifactReferenceStub != nil {
		return fake.UpdateArtifactReferenceStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.updateArtifactReferenceReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePivnetClient) UpdateArtifactReferenceCallCount() int {
	fake.updateArtifactReferenceMutex.RLock()
	defer fake.updateArtifactReferenceMutex.RUnlock()
	return len(fake.updateArtifactReferenceArgsForCall)
}

func (fake *FakePivnetClient) UpdateArtifactReferenceCalls(stub func(string, pivnet.ArtifactReference) (pivnet.ArtifactReference, error)) {
	fake.updateArtifactReferenceMutex.Lock()
	defer fake.updateArtifactReferenceMutex.Unlock()
	fake.UpdateArtifactReferenceStub = stub
}

func (fake *FakePivnetClient) UpdateArtifactReferenceArgsForCall(i int) (string, pivnet.ArtifactReference) {
	fake.updateArtifactReferenceMutex.RLock()
	defer fake.updateArtifactReferenceMutex.RUnlock()
	argsForCall := fake.updateArtifactReferenceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePivnetClient) UpdateArtifactReferenceReturns(result1 pivnet.ArtifactReference, result2 error) {
	fake.updateArtifactReferenceMutex.Lock()
	defer fake.updateArtifactReferenceMutex.Unlock()
	fake.UpdateArtifactReferenceStub = nil
	fake.updateArtifactReferenceReturns = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) UpdateArtifactReferenceReturnsOnCall(i int, result1 pivnet.ArtifactReference, result2 error) {
	fake.updateArtifactReferenceMutex.Lock()
	defer fake.updateArtifactReferenceMutex.Unlock()
	fake.UpdateArtifactReferenceStub = nil
	if fake.updateArtifactReferenceReturnsOnCall == nil {
		fake.updateArtifactReferenceReturnsOnCall = make(map[int]struct {
			result1 pivnet.ArtifactReference
			result2 error
		})
	}
	fake.updateArtifactReferenceReturnsOnCall[i] = struct {
		result1 pivnet.ArtifactReference
		result2 error
	}{result1, result2}
}

func (fake *FakePivnetClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addArtifactReferenceToReleaseMutex.RLock()
	defer fake.addArtifactReferenceToReleaseMutex.RUnlock()
	fake.artifactReferenceMutex.RLock()
	defer fake.artifactReferenceMutex.RUnlock()
	fake.artifactReferenceForReleaseMutex.RLock()
	defer fake.artifactReferenceForReleaseMutex.RUnlock()
	fake.artifactReferencesMutex.RLock()
	defer fake.artifactReferencesMutex.RUnlock()
	fake.artifactReferencesForDigestMutex.RLock()
	defer fake.artifactReferencesForDigestMutex.RUnlock()
	fake.artifactReferencesForReleaseMutex.RLock()
	defer fake.artifactReferencesForReleaseMutex.RUnlock()
	fake.createArtifactReferenceMutex.RLock()
	defer fake.createArtifactReferenceMutex.RUnlock()
	fake.deleteArtifactReferenceMutex.RLock()
	defer fake.deleteArtifactReferenceMutex.RUnlock()
	fake.releaseForVersionMutex.RLock()
	defer fake.releaseForVersionMutex.RUnlock()
	fake.removeArtifactReferenceFromReleaseMutex.RLock()
	defer fake.removeArtifactReferenceFromReleaseMutex.RUnlock()
	fake.updateArtifactReferenceMutex.RLock()
	defer fake.updateArtifactReferenceMutex.RUnlock()
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

var _ artifactreference.PivnetClient = new(FakePivnetClient)
