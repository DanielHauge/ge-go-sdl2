package ge_go_sdl2

import (
	"reflect"
)

var (
	updateChannels        map[string]chan propertyChange
	propertyChangeChannel chan propertyChange
)

type propertyChange struct {
	Id     string
	Name   string
	Value  interface{}
	DoneCB chan int
}

func init() {
	updateChannels = make(map[string]chan propertyChange)
	propertyChangeChannel = make(chan propertyChange)
}

func handlePropertyChanges() {
	for {
		pc := <-propertyChangeChannel
		channel, ok := updateChannels[pc.Id]
		if ok {
			channel <- pc
		}
	}
}

func handleElementPropertyChanges[T uiElement](pcChan <-chan propertyChange, ph *T, redrawFunc func(t *T)) {
	for {
		select {
		case pc := <-pcChan:
			changeElementProperty(pc, ph)
			redrawFunc(ph)
			window.UpdateSurface()
			if pc.DoneCB != nil {
				go func() { pc.DoneCB <- 1 }()
			}
		}
	}
}

func changeElementProperty(pc propertyChange, ph interface{}) {
	value := reflect.ValueOf(ph)
	pcValue := reflect.ValueOf(pc.Value)
	if value.Elem().Kind() == reflect.Struct {
		structValue := value.Elem()
		fieldValue := structValue.FieldByName(pc.Name)
		fieldValue.Set(pcValue)
		ref := reflect.Indirect(value)
		ref.Set(structValue)
	}
}

func NotifyPropertyChange(id string, propertyName string, newValue interface{}) {
	propertyChangeChannel <- propertyChange{Id: id, Name: propertyName, Value: newValue}
}

func NotifyPropertyChangeCb(id string, propertyName string, newValue interface{}, cb chan int) {
	propertyChangeChannel <- propertyChange{Id: id, Name: propertyName, Value: newValue, DoneCB: cb}
}

func NotifyPropertyChangeCbAsync(id string, propertyName string, newValue interface{}, cb chan int) {
	go NotifyPropertyChangeCb(id, propertyName, newValue, cb)
}

func NotifyPropertyChangeAsync(id string, propertyName string, newValue interface{}) {
	go NotifyPropertyChange(id, propertyName, newValue)
}
