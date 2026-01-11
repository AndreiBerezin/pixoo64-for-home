package pixoo64

func (p *Pixoo64) SetBrightness(brightness int) error {
	data := map[string]any{
		"Command":    "Channel/SetBrightness",
		"Brightness": brightness,
	}

	return p.callApi(data)
}

func (p *Pixoo64) OnOffScreen(on bool) error {
	onOff := 0
	if on {
		onOff = 1
	}
	data := map[string]any{
		"Command": "Channel/OnOffScreen",
		"OnOff":   onOff,
	}

	return p.callApi(data)
}
