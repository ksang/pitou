package board

import (
	"os"
	"testing"

	"github.com/ajstarks/svgo"
	"github.com/ksang/gotopo"
)

var (
	// node1 <-> node2 <-> node3 <-> node4
	simpleTopo = gotopo.Topology{
		Nodes: []gotopo.Node{
			gotopo.Node{
				ID:   1,
				Name: "node1",
			},
			gotopo.Node{
				ID:   2,
				Name: "node2",
			},
			gotopo.Node{
				ID:   3,
				Name: "node3",
			},
			gotopo.Node{
				ID:   4,
				Name: "node4",
			},
		},
		Graph: [][]*gotopo.Link{
			[]*gotopo.Link{
				nil, &gotopo.Link{}, nil, nil,
			},
			[]*gotopo.Link{
				&gotopo.Link{}, nil, &gotopo.Link{}, nil,
			},
			[]*gotopo.Link{
				nil, &gotopo.Link{}, nil, &gotopo.Link{},
			},
			[]*gotopo.Link{
				nil, nil, &gotopo.Link{}, nil,
			},
		},
	}
)

func TestSVGTopo(t *testing.T) {
	width := 500
	height := 500
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 100)
	canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}
