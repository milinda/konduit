package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

func connectMqtt(brokerUrl string, userName string, password string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(brokerUrl)

	if len(userName) > 0 && len(password) > 0 {
		opts.SetUsername(userName)
		opts.SetPassword(password)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}


func startMqtt(hub *Hub, brokerUrl string, userName string, password string) {
	var err error
	var mqttClient mqtt.Client

	mqttClient, err = connectMqtt(brokerUrl, userName, password)

	if err != nil {
		zap.S().Panicf("Cannot connect to MQTT broker at %s", brokerUrl)
	} else {
		if token := mqttClient.Subscribe("homeassistant/#", 0, func(client mqtt.Client, message mqtt.Message) {
			notification := make(map[string]interface{})
			notification["topic"] = message.Topic()
			notification["payload"] = message.Payload()
			hub.notifications <- notification
		}); token.Wait() && token.Error() != nil {
			zap.S().Panic(token.Error())
		}
	}
}
