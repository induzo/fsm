package fsm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/emicklei/dot"
)

// Action represents an action
type Action string

// Status represent a status
type Status string

// StatusAction represents either
// an action and its origin state (Source)
// an action and its state resulting of the action applied (Outcome)
type StatusAction struct {
	Action Action `json:"action"`
	Status Status `json:"status"`
}

// Node represent a Node in the state graph
type Node struct {
	Sources  []StatusAction `json:"sources,omitempty"`
	Status   Status         `json:"status"`
	Outcomes []StatusAction `json:"outcomes,omitempty"`
}

// Graph represent the container of all nodes
type Graph map[Status]*Node

// NewGraph will create a new Graph
func NewGraph() Graph {
	return Graph(make(map[Status]*Node))
}

// AddNode will add a node to the state graph
func (g Graph) AddNode(n *Node) error {
	// Make sure it doesn't already exists
	if _, ok := g[n.Status]; ok {
		return fmt.Errorf("AddNode(%s): node already exists", n.Status)
	}

	// Make sure all srcs and outcomes nodes exists
	for _, nso := range append(n.Sources, n.Outcomes...) {
		// Make sure it doesn't already exists
		if _, ok := g[nso.Status]; !ok {
			return fmt.Errorf(
				"AddNode(%s): node %s doesn't exist",
				n.Status, nso.Status,
			)
		}
	}

	// make sure the source node has this one as an outcome
	if err := g.appendOutcomeToSourceNodes(n.Sources, n.Status); err != nil {
		return fmt.Errorf("AddNode(%s): %v", n.Status, err)
	}
	// make sure the outcome node has this one as a source
	if err := g.appendSourceToOutcomeNodes(n.Status, n.Outcomes); err != nil {
		return fmt.Errorf("AddNode(%s): %v", n.Status, err)
	}

	g[n.Status] = n

	return nil
}

func (g Graph) appendOutcomeToSourceNodes(
	srcs []StatusAction,
	st Status,
) error {
	for _, src := range srcs {
		if _, ok := g[src.Status]; !ok {
			return fmt.Errorf(
				"appendOutcomeToSourceNodes: source node %s doesn't exist",
				src.Status,
			)
		}

		// Add the outcome to the source node, if it doesn't already exist
		outcomeExists := false
		for _, srcO := range g[src.Status].Outcomes {
			if srcO.Action == src.Action && srcO.Status == st {
				outcomeExists = true
			}
		}
		if !outcomeExists {
			g[src.Status].Outcomes = append(
				g[src.Status].Outcomes,
				StatusAction{
					Action: src.Action,
					Status: st,
				},
			)
		}
	}
	return nil
}

func (g Graph) appendSourceToOutcomeNodes(
	st Status,
	os []StatusAction,
) error {
	for _, o := range os {
		if _, ok := g[o.Status]; !ok {
			return fmt.Errorf(
				"appendSourceToOutcomeNodes: outcome node %s doesn't exist",
				o.Status,
			)
		}

		// Add the outcome to the source node, if it doesn't already exist
		srcExists := false
		for _, oSrc := range g[o.Status].Sources {
			if oSrc.Action == o.Action && oSrc.Status == st {
				srcExists = true
			}
		}
		if !srcExists {
			g[o.Status].Outcomes = append(
				g[o.Status].Outcomes,
				StatusAction{
					Action: o.Action,
					Status: st,
				},
			)
		}
	}
	return nil
}

// Viz will create the graphviz graph for the graph
func (g Graph) Viz() *dot.Graph {
	gv := dot.NewGraph(dot.Directed)

	ns := make(map[string]dot.Node)
	for _, n := range g {
		ns[string(n.Status)] = gv.Node(string(n.Status))
	}

	for _, n := range g {
		for _, o := range n.Outcomes {
			gv.Edge(
				ns[string(n.Status)],
				ns[string(o.Status)],
				string(o.Action),
			)
		}
	}

	return gv
}

// GetOutcomeStatus will return the potential outcomes for a curr status
func (g Graph) GetOutcomeStatus(st Status, ac Action) (Status, error) {
	n, ok := g[st]
	if !ok {
		return "",
			fmt.Errorf(
				"GetOutcomeStatus(%s, %s): node %s not found",
				st, ac, st,
			)
	}
	for _, o := range n.Outcomes {
		if o.Action == ac {
			return o.Status, nil
		}
	}
	return "", fmt.Errorf(
		"GetOutcomeStatus(%s, %s): action %s not found",
		st, ac, ac,
	)
}

// GeneratePNG will generate a png from the graph
// dependency on dot
func (g Graph) GeneratePNG(f io.Writer) error {
	path, err := exec.LookPath("dot")
	if err != nil {
		return fmt.Errorf("GeneratePNG: you must have dot installed")
	}
	// Generate PNG graph
	v := g.Viz()

	// nolint: gosec
	dotgraph := exec.Command(path, "-Tpng")
	dotgraph.Stdin = strings.NewReader(v.String())

	dotgraph.Stdout = f
	if err := dotgraph.Run(); err != nil {
		return fmt.Errorf("GeneratePNG: %v", err)
	}
	return nil
}

// UnmarshalJSON will parse a json into a tree of nodes
func (g Graph) UnmarshalJSON(data []byte) error {
	r := bytes.NewReader(data)
	var aux []*Node

	if err := json.NewDecoder(r).Decode(&aux); err != nil {
		return fmt.Errorf("UnmarshalJSON: %v", err)
	}

	for _, a := range aux {
		if err := g.AddNode(a); err != nil {
			return fmt.Errorf("UnmarshalJSON: %v", err)
		}
	}

	return nil
}

// MarshalJSON will return json
func (g Graph) MarshalJSON() ([]byte, error) {
	var ns []*Node
	for _, n := range g {
		ns = append(ns, n)
	}

	return json.Marshal(ns)
}
