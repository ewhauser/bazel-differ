package internal

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
)

//go:generate mockgen -destination=../mocks/rule_provider_mock.go -package=mocks github.com/ewhauser/bazel-differ/internal RuleProvider
type RuleProvider interface {
	GetRule(rule *Rule) BazelRule
}

type ruleProvider struct {
}

func (r ruleProvider) GetRule(rule *Rule) BazelRule {
	return &bazelRule{rule: rule}
}

func NewRuleProvider() RuleProvider {
	return &ruleProvider{}
}

//go:generate mockgen -destination=../mocks/bazel_rule_mock.go -package=mocks github.com/ewhauser/bazel-differ/internal BazelRule
type BazelRule interface {
	GetDigest() ([]byte, error)
	GetRuleInputList() []string
	GetName() string
}

type bazelRule struct {
	rule *Rule
}

func (b *bazelRule) GetDigest() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if b.rule.RuleClass != nil {
		if _, err := buffer.Write([]byte(*b.rule.RuleClass)); err != nil {
			return nil, err
		}
	}
	if b.rule.Name != nil {
		if _, err := buffer.Write([]byte(*b.rule.Name)); err != nil {
			return nil, err
		}
	}
	if b.rule.SkylarkEnvironmentHashCode != nil {
		if _, err := buffer.Write([]byte(*b.rule.SkylarkEnvironmentHashCode)); err != nil {
			return nil, err
		}
	}
	if b.rule.Attribute != nil {
		for _, attribute := range b.rule.Attribute {
			if _, err := buffer.Write([]byte(attribute.String())); err != nil {
				return nil, err
			}
		}
	}
	checksum := sha256.Sum256(buffer.Bytes())
	return checksum[:], nil
}

func (b *bazelRule) GetRuleInputList() []string {
	var ruleInputList []string
	for _, ruleInput := range b.rule.RuleInput {
		if strings.HasPrefix(ruleInput, "@") {
			splitRule := strings.Split(ruleInput, "//")
			if len(splitRule) == 2 {
				externalRule := splitRule[0]
				externalRule = strings.ReplaceAll(externalRule, "@", "")
				ruleInputList = append(ruleInputList, fmt.Sprintf("//external:%s", externalRule))
			} else {
				ruleInputList = append(ruleInputList, ruleInput)
			}
		} else {
			ruleInputList = append(ruleInputList, ruleInput)
		}
	}
	return ruleInputList
}

func (b *bazelRule) GetName() string {
	return b.rule.GetName()
}

func (b *bazelRule) GetRuleClass() string {
	return b.rule.GetRuleClass()
}
