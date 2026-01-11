package pixoo64

func (p *Pixoo64) PlayBuzzer(activeTimeInCycle int, offTimeInCycle int, playTotalTime int) error {
	data := map[string]any{
		"Command":           "Device/PlayBuzzer",
		"ActiveTimeInCycle": activeTimeInCycle,
		"OffTimeInCycle":    offTimeInCycle,
		"PlayTotalTime":     playTotalTime,
	}

	return p.callApi(data)
}
