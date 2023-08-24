package keyboard

func ListenFunc(fn func(event KeyEvent)) error {
	keysEvents, err := GetKeys(10)
	if err != nil {
		return err
	}
	defer Close()
	for {
		event := <-keysEvents
		if event.Err != nil {
			continue
		}
		fn(event)
	}
}

func ListenKey(key Key) <-chan struct{} {
	c := make(chan struct{})
	go ListenFunc(func(event KeyEvent) {
		if event.Key == key {
			c <- struct{}{}
		}
	})
	return c
}
