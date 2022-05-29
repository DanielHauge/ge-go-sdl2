package ge_go_sdl2

type PropertyChange struct {
	Id    string
	Name  string
	Value interface{}
}

func handlePropertyChanges(pcChan <-chan PropertyChange, quit <-chan int) {
	for {
		select {
		case <-quit:
			return
		case pc := <-pcChan:
			changeProperty(pc)
		}
	}
}

func changeProperty(pc PropertyChange) {

}
