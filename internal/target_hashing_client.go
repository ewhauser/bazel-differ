package internal

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"io/fs"
)

var HexArray = []byte("0123456789abcdef")

type CacheOptions struct {
}

type TargetHashingClient interface {
	HashAllBazelTargetsAndSourcefiles(seedFilePaths map[string]bool) (map[string]string, error)
	GetImpactedTargets(startHashes map[string]string, endHashes map[string]string) (map[string]bool, error)
	GetNames(targets []*Target) map[string]bool
}

type targetHashingClient struct {
	bazelClient  BazelClient
	filesystem   fs.FS
	ruleProvider RuleProvider
}

func NewTargetHashingClient(client BazelClient, filesystem fs.FS, ruleProvider RuleProvider) TargetHashingClient {
	return &targetHashingClient{
		bazelClient:  client,
		filesystem:   filesystem,
		ruleProvider: ruleProvider,
	}
}

func (t targetHashingClient) HashAllBazelTargetsAndSourcefiles(seedFilePaths map[string]bool) (
	map[string]string, error) {
	bazelSourcefileTargets, err := t.bazelClient.QueryAllSourceFileTargets()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(seedFilePaths))
	for k := range seedFilePaths {
		keys = append(keys, k)
	}

	seedHashes, err := createSeedForFilepaths(t.filesystem, keys)
	if err != nil {
		return nil, err
	}

	return t.hashAllTargets(seedHashes, bazelSourcefileTargets)
}

func (t targetHashingClient) GetImpactedTargets(startHashes map[string]string, endHashes map[string]string) (map[string]bool, error) {
	impactedTargets := make(map[string]bool)
	for target, endHashValue := range endHashes {
		startHashValue, ok := startHashes[target]
		if !ok || startHashValue != endHashValue {
			impactedTargets[target] = true
		}
	}
	return impactedTargets, nil
}

func (t targetHashingClient) GetNames(targets []*Target) map[string]bool {
	targetNames := make(map[string]bool, len(targets))
	for _, target := range targets {
		targetNames[getNameForTarget(target)] = true
	}
	return targetNames
}

func createSeedForFilepaths(filesys fs.FS, seedFilepaths []string) ([]byte, error) {
	if len(seedFilepaths) == 0 {
		return []byte{}, nil
	}
	buffer := bytes.NewBuffer([]byte{})
	for _, path := range seedFilepaths {
		data, err := fs.ReadFile(filesys, path)
		if err != nil {
			return []byte{}, err
		}
		if _, err := buffer.Write(data); err != nil {
			return nil, err
		}
	}
	checksum := sha256.Sum256(buffer.Bytes())

	fmt.Println(b64.StdEncoding.EncodeToString(checksum[:]))

	return checksum[:], nil
}

func (t targetHashingClient) hashAllTargets(seedHash []byte,
	bazelSourcefileTargets map[string]*BazelSourceFileTarget) (
	map[string]string, error) {
	allTargets, err := t.bazelClient.QueryAllTargets()
	if err != nil {
		return nil, err
	}
	targetHashes := make(map[string]string)
	ruleHashes := make(map[string][]byte)
	allRulesMap := make(map[string]BazelRule)
	for _, target := range allTargets {
		targetName := getNameForTarget(target)
		if targetName == "" {
			continue
		}
		if target.Rule != nil {
			allRulesMap[targetName] = t.ruleProvider.GetRule(target.Rule)
		}

		if target.GeneratedFile != nil {
			allRulesMap[targetName] = allRulesMap[target.GetGeneratedFile().GetGeneratingRule()]
		}
	}

	for _, target := range allTargets {
		targetName := getNameForTarget(target)
		if targetName == "" {
			continue
		}
		targetDigest, err := t.createDigestForTarget(
			target,
			allRulesMap,
			bazelSourcefileTargets,
			ruleHashes,
			seedHash,
		)

		if err != nil {
			return nil, err
		}

		if targetDigest != nil {
			targetHashes[targetName] = convertByteArrayToString(targetDigest)
		}
	}
	return targetHashes, nil
}

func (t targetHashingClient) createDigestForTarget(
	target *Target,
	allRulesMap map[string]BazelRule,
	bazelSourcefileTargets map[string]*BazelSourceFileTarget,
	ruleHashes map[string][]byte,
	seedHash []byte,
) ([]byte, error) {
	if target.SourceFile != nil {
		var sourceFileName = getNameForTarget(target)
		if sourceFileName != "" {
			buffer := bytes.NewBuffer([]byte{})
			sourceTargetDigestBytes, err := getDigestForSourceTargetName(sourceFileName, bazelSourcefileTargets)
			if err != nil {
				return nil, err
			}
			if sourceTargetDigestBytes != nil {
				if _, err := buffer.Write(sourceTargetDigestBytes); err != nil {
					return nil, err
				}
			}
			if seedHash != nil {
				if _, err := buffer.Write(seedHash); err != nil {
					return nil, err
				}
			}
			checksum := sha256.Sum256(buffer.Bytes())
			return checksum[:], nil
		}
	}
	if target.GeneratedFile != nil {
		var generatingRuleDigest = ruleHashes[target.GeneratedFile.GetGeneratingRule()]
		if generatingRuleDigest != nil {
			return createDigestForRule(allRulesMap[target.GeneratedFile.GetGeneratingRule()], allRulesMap, ruleHashes,
				bazelSourcefileTargets,
				seedHash)
		}
		return generatingRuleDigest, nil
	}
	return createDigestForRule(t.ruleProvider.GetRule(target.Rule), allRulesMap, ruleHashes, bazelSourcefileTargets, seedHash)
}

func createDigestForRule(
	rule BazelRule,
	allRulesMap map[string]BazelRule,
	ruleHashes map[string][]byte,
	bazelSourcefileTargets map[string]*BazelSourceFileTarget,
	seedHash []byte,
) ([]byte, error) {
	existingByteArray := ruleHashes[rule.Name()]
	if existingByteArray != nil {
		return existingByteArray, nil
	}
	buffer := bytes.NewBuffer([]byte{})
	ruleDigest, err := rule.Digest()
	if err != nil {
		return nil, err
	}
	if _, err := buffer.Write(ruleDigest); err != nil {
		return nil, err
	}
	if seedHash != nil {
		if _, err := buffer.Write(seedHash); err != nil {
			return nil, err
		}
	}
	for _, ruleInput := range rule.RuleInputList() {
		if _, err := buffer.Write([]byte(ruleInput)); err != nil {
			return nil, err
		}
		inputRule := allRulesMap[ruleInput]
		sourceFileDigest, err := getDigestForSourceTargetName(ruleInput, bazelSourcefileTargets)
		if err != nil {
			return nil, err
		}

		if inputRule != nil && inputRule.Name() != "" && !(inputRule.Name() == rule.Name()) {
			inputRuleHash, err := createDigestForRule(
				inputRule,
				allRulesMap,
				ruleHashes,
				bazelSourcefileTargets,
				seedHash,
			)
			if err != nil {
				return nil, err
			}
			if inputRuleHash != nil {
				if _, err := buffer.Write(inputRuleHash); err != nil {
					return nil, err
				}
			}
		} else if sourceFileDigest != nil {
			if _, err := buffer.Write(sourceFileDigest); err != nil {
				return nil, err
			}
		}
	}
	checksum := sha256.Sum256(buffer.Bytes())
	finalHashValue := checksum[:]
	ruleHashes[rule.Name()] = finalHashValue
	return finalHashValue, nil
}

func getNameForTarget(target *Target) string {
	if target.Rule != nil {
		return target.Rule.GetName()
	}
	if target.SourceFile != nil {
		return *target.SourceFile.Name
	}
	if target.GeneratedFile != nil {
		return *target.GeneratedFile.Name
	}
	return ""
}

func getDigestForSourceTargetName(sourceTargetName string, bazelSourcefileTargets map[string]*BazelSourceFileTarget) ([]byte, error) {
	target, ok := bazelSourcefileTargets[sourceTargetName]
	if !ok {
		return nil, nil
	}
	return (*target).Digest(), nil
}

func convertByteArrayToString(bytes []byte) string {
	hexChars := make([]byte, len(bytes)*2)
	for i := 0; i < len(bytes); i++ {
		v := bytes[i] & 0xFF
		hexChars[i*2] = HexArray[v>>4]
		hexChars[i*2+1] = HexArray[v&0x0F]
	}
	return string(hexChars)
}
