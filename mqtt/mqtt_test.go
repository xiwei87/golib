package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"testing"
)

func TestInitMqttClient(t *testing.T) {
	var config MqttConfig
	config.SecretKey = "123abc0099"
	config.ClientId = "box_210001"
	config.StationId = "111"
	config.ServerUri = "tls://139.224.238.81:8883"
	config.WillTopic = "box_disconnection"
	err := InitClient(&config)
	if err != nil {
		t.Error(err)
	}
	go func() {
		Subscribe("test/#", func(client mqtt.Client, msg mqtt.Message) {
			if string(msg.Payload()) != "mymessage" {
				t.Fatalf("want mymessage, got %s", msg.Payload())
			}
		})
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
