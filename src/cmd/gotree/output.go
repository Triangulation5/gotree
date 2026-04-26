package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
	"gotree/src/internal/model"
	"gotree/src/internal/printer"
)

type outputPayload struct {
	Root    outputNode      `json:"root" yaml:"root"`
	Summary printer.Summary `json:"summary" yaml:"summary"`
}

type outputNode struct {
	Name         string       `json:"name" yaml:"name"`
	Path         string       `json:"path" yaml:"path"`
	IsDir        bool         `json:"is_dir" yaml:"is_dir"`
	SizeBytes    int64        `json:"size_bytes" yaml:"size_bytes"`
	ModifiedUnix int64        `json:"modified_unix" yaml:"modified_unix"`
	SummaryMsg   string       `json:"summary_message,omitempty" yaml:"summary_message,omitempty"`
	Children     []outputNode `json:"children,omitempty" yaml:"children,omitempty"`
}

func printStructured(tree *model.Node, asJSON bool, asYAML bool) error {
	payload := outputPayload{
		Root:    buildOutputNode(tree),
		Summary: printer.BuildSummary(tree),
	}
	if asJSON {
		return encodeJSON(payload)
	}
	if asYAML {
		return encodeYAML(payload)
	}
	return nil
}

func buildOutputNode(n *model.Node) outputNode {
	node := outputNode{
		Name:         n.Name,
		Path:         n.Path,
		IsDir:        n.IsDir,
		SizeBytes:    n.SizeBytes,
		ModifiedUnix: n.ModifiedUnix,
		SummaryMsg:   n.SummaryMsg,
	}
	if len(n.Children) > 0 {
		node.Children = make([]outputNode, 0, len(n.Children))
		for _, child := range n.Children {
			node.Children = append(node.Children, buildOutputNode(child))
		}
	}
	return node
}

func encodeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func encodeYAML(v any) error {
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	return enc.Encode(v)
}
