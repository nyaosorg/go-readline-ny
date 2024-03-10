package auto

type Pilot struct {
	Text []string
}

func (*Pilot) Open(func(int)) error {
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

func (*Pilot) Size() (int, int, error) {
	return 80, 25, nil
}

func (*Pilot) Close() error {
	return nil
}
