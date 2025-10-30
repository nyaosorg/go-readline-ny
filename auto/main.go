package auto

type Pilot struct {
	Text   []string
	Width  int
	Height int
}

func (p *Pilot) Open(func(int)) error {
	if p.Width <= 0 {
		p.Width = 80
	}
	if p.Height <= 0 {
		p.Height = 24
	}
	return nil
}

func (ap *Pilot) GetKey() (string, error) {
	if len(ap.Text) <= 0 {
		return "\r", nil
	}
	result := ap.Text[0]
	ap.Text = ap.Text[1:]
	return result, nil
}

func (p *Pilot) Size() (int, int, error) {
	return p.Width, p.Height, nil
}

func (*Pilot) Close() error {
	return nil
}
