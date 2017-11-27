package board

import (
	"io"
	"net/http"

	"github.com/ksang/gotopo"
)

func (b *Board) SVGHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	genTopoSVG(w, b.Topo)
}

func genTopoSVG(w io.Writer, topo *gotopo.Topology) {

}
