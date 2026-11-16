package upload

import "github.com/mind-vm/tabler-templ/kit"

func iconOf(p ButtonProps) string {
	if p.Camera {
		return "ti-camera"
	}
	return "ti-upload"
}

func labelOf(m kit.Messages, p ButtonProps) string {
	switch {
	case p.Label != "":
		return p.Label
	case p.Camera:
		return m.TakePhoto
	}
	return m.ChooseFile
}
