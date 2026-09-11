package worktree

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/gibbonmi/bench/internal/git"
)

func captureLayers(root, path string, keepWorking bool, tip string) (string, []string, string, error) {
	head, err := git.Output("-C", path, "rev-parse", "HEAD")
	if err != nil {
		return "", nil, "", fmt.Errorf("read recovery HEAD: %w", err)
	}
	headTree, err := git.Output("-C", path, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return "", nil, "", fmt.Errorf("read recovery base tree: %w", err)
	}
	admin, err := git.AdminDir(path)
	if err != nil {
		return "", nil, "", err
	}
	layerTrees := map[string]string{}
	workingTree, err := worktreeTree(path, admin)
	if err != nil {
		return "", nil, "", fmt.Errorf("capture working layer: %w", err)
	}
	if keepWorking || workingTree != headTree {
		layerTrees["working"] = workingTree
	}
	entries, conflicted, err := readIndexEntries(path)
	if err != nil {
		return "", nil, "", err
	}
	if conflicted {
		for stage, name := range map[int]string{1: "base", 2: "ours", 3: "theirs"} {
			tree, err := conflictTree(path, admin, entries, stage)
			if err != nil {
				return "", nil, "", fmt.Errorf("capture conflict %s layer: %w", name, err)
			}
			layerTrees[name] = tree
		}
	} else {
		stagedTree, err := realIndexTree(path, admin)
		if err != nil {
			return "", nil, "", fmt.Errorf("capture staged layer: %w", err)
		}
		if stagedTree != headTree {
			layerTrees["staged"] = stagedTree
		}
	}
	if len(layerTrees) == 0 {
		return "", nil, "", errors.New("recovery requested for a clean assignment")
	}
	treePayload := map[string]string{}
	layers := map[string]string{}
	names := make([]string, 0, len(layerTrees))
	for name := range layerTrees {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		tree := layerTrees[name]
		payload := treePayload[tree]
		if payload == "" {
			payload, err = commitTree(root, tree, []string{head}, "bench recovery payload: "+name+"\n")
			if err != nil {
				return "", nil, "", err
			}
			treePayload[tree] = payload
		}
		layers[name] = payload
	}
	payloads := make([]string, 0, len(treePayload))
	for _, payload := range treePayload {
		payloads = append(payloads, payload)
	}
	sort.Strings(payloads)
	manifest := recoveryManifest{Schema: recoverySchema, Base: head, Layers: layers}
	if tip != "" && tip != head {
		manifest.Tip = tip
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return "", nil, "", err
	}
	manifestBytes = append(manifestBytes, '\n')
	blob, err := gitInput(root, nil, manifestBytes, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", nil, "", err
	}
	rootTree, err := gitInput(root, nil, []byte("100644 blob "+blob+"\tmanifest.json\n"), "mktree")
	if err != nil {
		return "", nil, "", err
	}
	parents := append([]string(nil), payloads...)
	if manifest.Tip != "" {
		parents = append(parents, manifest.Tip)
	}
	rootOID, err := commitTree(root, rootTree, parents, "bench recovery root\n")
	if err != nil {
		return "", nil, "", err
	}
	return rootOID, payloads, head, nil
}
