package ge_go_sdl2

import (
	"reflect"
)

var (
	updateChannels        map[string]chan propertyChange
	PropertyChangeChannel chan propertyChange
)

type propertyChange struct {
	Id    string
	Name  string
	Value interface{}
}

func init() {
	updateChannels = make(map[string]chan propertyChange)
	PropertyChangeChannel = make(chan propertyChange)
}

func handlePropertyChanges() {
	for {
		pc := <-PropertyChangeChannel
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
	PropertyChangeChannel <- propertyChange{Id: id, Name: propertyName, Value: newValue}
}

func NotifyPropertyChangeAsync(id string, propertyName string, newValue interface{}) {
	go NotifyPropertyChange(id, propertyName, newValue)
}
