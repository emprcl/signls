package ui

import (
	"cykl/core"
)

func renderNode(node core.Node) string {
	switch node.(type) {
	case *core.OnceEmitter:
		switch node.(*core.OnceEmitter).Direction() {
		case 0:
			return "▀▀"
		case 1:
			return " █"
		case 2:
			return "▄▄"
		case 3:
			return "█ "
		}
	case *core.Signal:
		return "!!"
	}
	return ".."
}
