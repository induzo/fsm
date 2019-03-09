package fsm

import (
	"testing"
)

func TestGraph_AddNode(t *testing.T) {
	m := NewGraph()
	if err := m.AddNode(&Node{Status: "new"}); err != nil {
		t.Errorf("AddNode: %v", err)
		return
	}
	tests := []struct {
		name    string
		node    *Node
		wantErr bool
	}{
		{
			name: "linking the node accepted to node new",
			node: &Node{
				Sources: []StatusAction{
					StatusAction{
						Status: "new",
						Action: "accept",
					},
				},
				Status: "accepted",
			},
			wantErr: false,
		},
		{
			name: "linking the node accepted to node new",
			node: &Node{
				Sources: []StatusAction{
					StatusAction{
						Status: "new",
						Action: "bury",
					},
				},
				Status: "buried",
				Outcomes: []StatusAction{
					StatusAction{
						Status: "accepted",
						Action: "accept",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "linking the node canceled to node x, which doesn't exist",
			node: &Node{
				Sources: []StatusAction{
					StatusAction{
						Status: "x",
						Action: "cancel",
					},
				},
				Status: "canceled",
			},
			wantErr: true,
		},
		{
			name: "linking the node x to node new",
			node: &Node{
				Sources: []StatusAction{
					StatusAction{
						Status: "new",
						Action: "lose",
					},
				},
				Status: "lost",
				Outcomes: []StatusAction{
					StatusAction{
						Status: "x",
						Action: "accept",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "adding an existing node, which doesn't work",
			node: &Node{
				Status: "new",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.AddNode(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Graph.AddNode() error = %v, wantErr %v",
					err, tt.wantErr,
				)
			}
		})
	}
}

func TestGraph_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "working unmarshal",
			args: args{
				data: []byte(`
				[
					{
						"status": "New"
					}
				]`,
				),
			},
		},
		{
			name: "working unmarshal",
			args: args{
				data: getTestGraphJSON(),
			},
		},
		{
			name: "non working unmarshal",
			args: args{
				data: []byte(`
				[
						"status_id": "New"
					}
				]`,
				),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewGraph()
			if err := m.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Graph.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGraph_GetOutcomeStatus(t *testing.T) {
	type args struct {
		status Status
		action Action
	}
	tests := []struct {
		name    string
		args    args
		want    Status
		wantErr bool
	}{
		{
			name: "existing status and action",
			args: args{
				status: "new",
				action: "accept",
			},
			want: "accepted",
		},
		{
			name: "missing status",
			args: args{
				status: "old",
				action: "accept",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing action",
			args: args{
				status: "new",
				action: "miss",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := getTestGraph()
			got, err := g.GetOutcomeStatus(tt.args.status, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graph.GetOutcome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Graph.GetOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_Viz(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Working viz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := getTestGraph()
			_ = g.Viz()
		})
	}
}

func getTestGraphJSON() []byte {
	return []byte(`
	[
		{
			"status": "new"
		},
		{
			"status": "accepted",
			"sources": [
				{
					"status": "new",
					"action": "accept"
				}
			]
		},
		{
			"status": "canceled",
			"sources": [
				{
					"status": "new",
					"action": "cancel"
				},
				{
					"status": "accepted",
					"action": "cancel"
				}
			]
		}
	]`)
}

func getTestGraph() Graph {
	m := NewGraph()
	_ = m.AddNode(&Node{Status: "new"})
	_ = m.AddNode(
		&Node{
			Sources: []StatusAction{
				StatusAction{
					Status: "new",
					Action: "accept",
				},
			},
			Status: "accepted",
		},
	)
	_ = m.AddNode(
		&Node{
			Sources: []StatusAction{
				StatusAction{
					Status: "new",
					Action: "cancel",
				},
				StatusAction{
					Status: "accepted",
					Action: "cancel",
				},
			},
			Status: "canceled",
		},
	)

	return m
}
