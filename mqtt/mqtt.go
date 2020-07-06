package mqtt

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xiwei87/golib/utils"
)

var client mqtt.Client

const mqtt_tls_ca = `-----BEGIN CERTIFICATE-----
MIIDfTCCAmWgAwIBAgIJAOfJPbFi+b+tMA0GCSqGSIb3DQEBCwUAMFUxCzAJBgNV
BAYTAkNOMREwDwYDVQQIDAhaaGVqaWFuZzEQMA4GA1UEBwwHSGFuZ2hvdTEQMA4G
A1UECgwHNjZJRlVFTDEPMA0GA1UEAwwGUm9vdENBMB4XDTIwMDMwOTA1NTM0MFoX
DTMwMDMwNzA1NTM0MFowVTELMAkGA1UEBhMCQ04xETAPBgNVBAgMCFpoZWppYW5n
MRAwDgYDVQQHDAdIYW5naG91MRAwDgYDVQQKDAc2NklGVUVMMQ8wDQYDVQQDDAZS
b290Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCfDeMcIARLesQL
cWwvhRC2N4Lg+iPjpmkIEW/aAnkUTWF28AMGsnu4dP/ZJUnVlba82xi7KO3butND
xxQgEcNaghXiStg7G9sQo5kLFAsM49G3JBF76LirNhvCvxgH7/HNKC9PMkUzhRnN
MnWn51MO1fJ8eIcpsi9qJ7mvk4654WbN+LWEOKFLwy51YGVR0gxBy3H1lkKZI3IE
f39/HJqJYuj71Q7BmFqAc9zBVRwNKbqDGOEn3Xs1EmMwCplhMpLV1mFGSPZdslav
WGz6xeij5XSVOTQkkMMu0bAd2scJzppt3TIis1GM4ZOTnZkHDbYjLUZEw1OxB3w4
o9HHJq9VAgMBAAGjUDBOMB0GA1UdDgQWBBQppfCCWTXwai+yNUZZkbfL6OyITzAf
BgNVHSMEGDAWgBQppfCCWTXwai+yNUZZkbfL6OyITzAMBgNVHRMEBTADAQH/MA0G
CSqGSIb3DQEBCwUAA4IBAQB8QLwyaYRuWCKyHEJd9WpTPNZbPkCIPgBqz+wilLzX
Px0pHM/E32gcxtKrVhh6cCKuSUIKLmgPvGPBx+pEK0y8z0f2YzAzO6UVDEKgy6p4
we0vGxWBjqGo1wU8WSXCX4LXmPXw+qL67jGUaQEe33Zn6ffjR9A7v+U0zvlWLupr
LMgQ1svW2L79epeo/VvvAVmbmya04t61swbdbqdT+mOAp3wLceDA0eT5j7s+PQyP
OHaLngx4pkRnLDSu5BoavFZ0reShyzpVlbAY3vuEQ6K6rz6asDeUpNZINYMsiTOv
17owR7lHYkOWNXIEXAHzFVlSsR32GudDnuzXpjJPiUN5
-----END CERTIFICATE-----`

type MqttConfig struct {
	StationId string `yaml:"station_id"`
	ClientId  string `yaml:"client_id"`
	ServerUri string `yaml:"server_uri"`
	SecretKey string `yaml:"secret_key"`
	WillTopic string `yaml:"will_topic"`
}

//初始化连接服务器
func Connect(config *MqttConfig) error {
	var (
		err      error
		mqtt_opt *mqtt.ClientOptions
	)
	if mqtt_opt, err = initOptions(config); err != nil {
		return err
	}
	client = mqtt.NewClient(mqtt_opt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

//关闭连接
func Disconnect() {
	client.Disconnect(250)
}

//订阅主题
func Subscribe(topic string, callback mqtt.MessageHandler) error {
	if token := client.Subscribe(topic, 1, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

//发送消息
func Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if token := client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

//设置MQTT参数
func initOptions(config *MqttConfig) (*mqtt.ClientOptions, error) {
	//加载ROOT CA
	root_ca := x509.NewCertPool()
	load_ca := root_ca.AppendCertsFromPEM([]byte(mqtt_tls_ca))
	if !load_ca {
		return nil, errors.New("failed to parse root certificate")
	}
	tlsConfig := &tls.Config{RootCAs: root_ca, InsecureSkipVerify: true}
	//用户名
	user_name := fmt.Sprintf("%s&%d", config.ClientId, time.Now().UnixNano()/1e6)
	//密码
	password_hash := hmac.New(sha1.New, []byte(config.SecretKey))
	password_hash.Write([]byte(user_name))
	password := base64.StdEncoding.EncodeToString(password_hash.Sum(nil))
	//设置参数
	opts := mqtt.NewClientOptions()
	opts.SetTLSConfig(tlsConfig)
	opts.AddBroker(config.ServerUri)
	opts.SetClientID(utils.NewUUID())
	opts.SetCleanSession(true)
	opts.SetUsername(user_name)
	opts.SetPassword(password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetWill(config.WillTopic, config.ClientId, 1, false)

	return opts, nil
}
