package internal_test

import (
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"github.com/ewhauser/bazel-differ/mocks"
	"github.com/golang/mock/gomock"
	"testing"
	"testing/fstest"
)

var defaultTargets []*internal.Target

func createDefaultTargets(ctrl *gomock.Controller, ruleProvider *mocks.MockRuleProvider) {
	defaultTargets = []*internal.Target{
		createRuleTarget("rule1", []string{}, "rule1Digest", ctrl, ruleProvider),
		createRuleTarget("rule2", []string{}, "rule2Digest", ctrl, ruleProvider),
	}
}

func TestHashAllBazelTargets_ruleTargets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ruleProvider := mocks.NewMockRuleProvider(ctrl)

	createDefaultTargets(ctrl, ruleProvider)

	bazelClient := mocks.NewMockBazelClient(ctrl)

	bazelClient.EXPECT().QueryAllSourceFileTargets().Return(make(map[string]*internal.BazelSourceFileTarget), nil)
	bazelClient.EXPECT().QueryAllTargets().Return(defaultTargets, nil).AnyTimes()

	targetHashingClient := internal.NewTargetHashingClient(bazelClient, fstest.MapFS{}, ruleProvider)

	var seedFilePaths = map[string]bool{}

	hash, err := targetHashingClient.HashAllBazelTargetsAndSourcefiles(seedFilePaths)
	if err != nil {
		t.Error(err)
	}

	if len(hash) != 2 {
		t.Errorf("Expected 2 entries in hash map, got %d", len(hash))
	}

	if hash["rule1"] != "2c963f7c06bc1cead7e3b4759e1472383d4469fc3238dc42f8848190887b4775" {
		t.Errorf("Expected hash of rule1 to be 2c963f7c06bc1cead7e3b4759e1472383d4469fc3238dc42f8848190887b4775, got %s", hash["rule1"])
	}
	if hash["rule2"] != "bdc1abd0a07103cea34199a9c0d1020619136ff90fb88dcc3a8f873c811c1fe9" {
		t.Errorf("Expected hash of rule2 to be bdc1abd0a07103cea34199a9c0d1020619136ff90fb88dcc3a8f873c811c1fe9, got %s", hash["rule2"])
	}
}

func TestHashAllBazelTargets_ruleTargets_seedFilepaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ruleProvider := mocks.NewMockRuleProvider(ctrl)

	createDefaultTargets(ctrl, ruleProvider)

	bazelClient := mocks.NewMockBazelClient(ctrl)

	bazelClient.EXPECT().QueryAllSourceFileTargets().Return(make(map[string]*internal.BazelSourceFileTarget), nil)
	bazelClient.EXPECT().QueryAllTargets().Return(defaultTargets, nil).AnyTimes()

	m := fstest.MapFS{
		"somefile.txt": {
			Data: []byte("somecontents"),
		},
	}

	targetHashingClient := internal.NewTargetHashingClient(bazelClient, m, ruleProvider)

	var seedFilePaths = map[string]bool{"somefile.txt": true}

	hash, err := targetHashingClient.HashAllBazelTargetsAndSourcefiles(seedFilePaths)
	if err != nil {
		t.Error(err)
	}

	if len(hash) != 2 {
		t.Errorf("Expected 2 entries in hash map, got %d", len(hash))
	}

	if hash["rule1"] != "866578a84bf5c0ed4786469c98876be5995518598c8ddb69b4abb0ef5c50a8d5" {
		t.Errorf("Expected hash of rule1 to be 866578a84bf5c0ed4786469c98876be5995518598c8ddb69b4abb0ef5c50a8d5, got %s", hash["rule1"])
	}
	if hash["rule2"] != "3789f108d047c8f73a03af00804bad56e9f764a83fc6be4ccc599ee48efe7826" {
		t.Errorf("Expected hash of rule2 to be 3789f108d047c8f73a03af00804bad56e9f764a83fc6be4ccc599ee48efe7826, got %s", hash["rule2"])
	}
}

func TestHashAllBazelTargets_ruleTargets_ruleInputs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ruleProvider := mocks.NewMockRuleProvider(ctrl)

	createDefaultTargets(ctrl, ruleProvider)

	bazelClient := mocks.NewMockBazelClient(ctrl)

	ruleInputs := []string{"rule1"}
	rule3 := createRuleTarget("rule3", ruleInputs, "digest", ctrl, ruleProvider)
	defaultTargets = append(defaultTargets, rule3)
	rule4 := createRuleTarget("rule4", ruleInputs, "digest2", ctrl, ruleProvider)
	defaultTargets = append(defaultTargets, rule4)

	bazelClient.EXPECT().QueryAllSourceFileTargets().Return(make(map[string]*internal.BazelSourceFileTarget), nil)
	bazelClient.EXPECT().QueryAllTargets().Return(defaultTargets, nil).AnyTimes()

	targetHashingClient := internal.NewTargetHashingClient(bazelClient, fstest.MapFS{}, ruleProvider)

	var seedFilePaths = map[string]bool{}

	hash, err := targetHashingClient.HashAllBazelTargetsAndSourcefiles(seedFilePaths)

	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 4 {
		t.Fatal(fmt.Sprintf("expected 4 hash values, got %d", len(hash)))
	}
	if hash["rule1"] != "2c963f7c06bc1cead7e3b4759e1472383d4469fc3238dc42f8848190887b4775" {
		t.Fatal(fmt.Sprintf("expected hash value for rule1 to be 2c963f7c06bc1cead7e3b4759e1472383d4469fc3238dc42f8848190887b4775, got %s", hash["rule1"]))
	}
	if hash["rule2"] != "bdc1abd0a07103cea34199a9c0d1020619136ff90fb88dcc3a8f873c811c1fe9" {
		t.Fatal(fmt.Sprintf("expected hash value for rule2 to be bdc1abd0a07103cea34199a9c0d1020619136ff90fb88dcc3a8f873c811c1fe9, got %s", hash["rule2"]))
	}
	if hash["rule3"] != "87dd050f1ca0f684f37970092ff6a02677d995718b5a05461706c0f41ffd4915" {
		t.Fatal(fmt.Sprintf("expected hash value for rule3 to be 87dd050f1ca0f684f37970092ff6a02677d995718b5a05461706c0f41ffd4915, got %s", hash["rule3"]))
	}
	if hash["rule4"] != "a7bc5d23cd98c4942dc879c649eb9646e38eddd773f9c7996fa0d96048cf63dc" {
		t.Fatal(fmt.Sprintf("expected hash value for rule4 to be a7bc5d23cd98c4942dc879c649eb9646e38eddd773f9c7996fa0d96048cf63dc, got %s", hash["rule4"]))
	}
}

func Test_hashAllBazelTargets_ruleTargets_ruleInputsWithSelfInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ruleProvider := mocks.NewMockRuleProvider(ctrl)

	createDefaultTargets(ctrl, ruleProvider)

	bazelClient := mocks.NewMockBazelClient(ctrl)

	ruleInputs := []string{"rule1", "rule4"}
	rule3 := createRuleTarget("rule3", ruleInputs, "digest", ctrl, ruleProvider)
	defaultTargets = append(defaultTargets, rule3)
	rule4 := createRuleTarget("rule4", ruleInputs, "digest2", ctrl, ruleProvider)
	defaultTargets = append(defaultTargets, rule4)
	bazelClient.EXPECT().QueryAllSourceFileTargets().Return(make(map[string]*internal.BazelSourceFileTarget), nil)
	bazelClient.EXPECT().QueryAllTargets().Return(defaultTargets, nil).AnyTimes()

	targetHashingClient := internal.NewTargetHashingClient(bazelClient, fstest.MapFS{}, ruleProvider)

	var seedFilePaths = map[string]bool{}

	hash, err := targetHashingClient.HashAllBazelTargetsAndSourcefiles(seedFilePaths)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(hash) != 4 {
		t.Errorf("expected 4 keys, got %d", len(hash))
	}
	if hash["rule4"] != "bf15e616e870aaacb02493ea0b8e90c6c750c266fa26375e22b30b78954ee523" {
		t.Errorf("expected %s, got %s", "bf15e616e870aaacb02493ea0b8e90c6c750c266fa26375e22b30b78954ee523", hash["rule4"])
	}
}

func TestHashAllBazelTargets_generatedTargets(t *testing.T) {
	t.Skip("Work in progress...")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ruleProvider := mocks.NewMockRuleProvider(ctrl)

	generator := createRuleTarget("rule1", []string{"rule0"}, "rule1Digest", ctrl, ruleProvider)
	target := createGeneratedTarget("rule0", "rule1")

	ruleInputs := []string{"rule0"}
	rule3 := createRuleTarget("rule3", ruleInputs, "digest", ctrl, ruleProvider)

	oldHash := ""
	newHash := ""

	bazelClient := mocks.NewMockBazelClient(ctrl)

	bazelClient.EXPECT().QueryAllSourceFileTargets().Return(make(map[string]*internal.BazelSourceFileTarget), nil)
	bazelClient.EXPECT().QueryAllTargets().Return([]*internal.Target{rule3, target, generator}, nil)
	targetHashingClient := internal.NewTargetHashingClient(bazelClient, fstest.MapFS{}, ruleProvider)

	hash, err := targetHashingClient.HashAllBazelTargetsAndSourcefiles(make(map[string]bool))
	if err != nil {
		t.Fatal(err)
	}

	if len(hash) != 3 {
		t.Error("Expecting 3 results")
	}
	oldHash = hash["rule3"]

	ruleProvider.EXPECT().GetRule(NameMatcher("rule1")).DoAndReturn(func(*internal.Rule) internal.BazelRule {
		rule := mocks.NewMockBazelRule(ctrl)
		rule.EXPECT().Name().Return("rule1").AnyTimes()
		rule.EXPECT().RuleInputList().Return([]string{}).AnyTimes()
		rule.EXPECT().Digest().Return([]byte("newDigest"), nil).AnyTimes()
		return rule
	})

	hash, err = targetHashingClient.HashAllBazelTargetsAndSourcefiles(make(map[string]bool))
	if err != nil {
		t.Fatal(err)
	}

	if len(hash) != 3 {
		t.Error("Expecting 3 results")
	}
	newHash = hash["rule3"]

	if oldHash != newHash {
		t.Errorf("Expected %s and %s to not be equal", oldHash, newHash)
	}
}

type nameMatcher struct {
	name string
}

func (m nameMatcher) Matches(arg interface{}) bool {
	sarg := arg.(*internal.Rule)
	return *sarg.Name == m.name
}

func (m nameMatcher) String() string {
	return "Matches on name: " + m.name
}

func NameMatcher(name string) gomock.Matcher {
	return &nameMatcher{name: name}
}

func createRuleTarget(name string, inputs []string, digest string, ctrl *gomock.Controller, provider *mocks.MockRuleProvider) *internal.Target {
	target := &internal.Target{
		Rule: &internal.Rule{
			Name:      &name,
			RuleInput: inputs,
		},
	}

	provider.EXPECT().GetRule(NameMatcher(name)).DoAndReturn(func(*internal.Rule) internal.BazelRule {
		rule := mocks.NewMockBazelRule(ctrl)
		rule.EXPECT().Name().Return(name).AnyTimes()
		rule.EXPECT().RuleInputList().Return(inputs).AnyTimes()
		rule.EXPECT().Digest().Return([]byte(digest), nil).AnyTimes()
		return rule
	}).AnyTimes()

	return target
}

//nolint
func createGeneratedTarget(name string, generatingRuleName string) *internal.Target {
	generatingFileName := generatingRuleName + "_file"
	target := &internal.Target{
		Rule: &internal.Rule{
			Name: &name,
		},
		GeneratedFile: &internal.GeneratedFile{
			Name:           &generatingFileName,
			GeneratingRule: &generatingRuleName,
		},
	}
	return target
}
