/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kafka_test

import (
	"github.com/RedHatInsights/insights-operator-utils/kafka"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	brokerCfg = kafka.BrokerConfiguration{
		Address: "localhost:1234",
		Topic:   "consumer-topic",
		Group:   "test-group",
		Enabled: true,
	}
	// Base UNIX time plus approximately 50 years (not long before year 2020).
	testTimestamp = time.Unix(50*365*24*60*60, 0)
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

// Test Producer creation with a non-accessible Kafka broker
func TestNewProducerBadBroker(t *testing.T) {
	const expectedErr = "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)"

	_, err := kafka.NewProducer(brokerCfg)
	assert.EqualError(t, err, expectedErr)
}

// TestProducerClose makes sure it's possible to close the producer.
func TestProducerClose(t *testing.T) {
	p := kafka.KafkaProducer{
		Producer: mocks.NewSyncProducer(t, nil),
	}

	err := p.Close()
	assert.NoError(t, err, "failed to close Kafka producer")
}

func TestNewProducer(t *testing.T) {
	mockBroker := sarama.NewMockBroker(t, 0)
	defer mockBroker.Close()

	mockBroker.SetHandlerByMap(helpers.GetHandlersMapForMockConsumer(t, mockBroker, brokerCfg.Topic))

	prod, err := kafka.NewProducer(
		kafka.BrokerConfiguration{
			Address: mockBroker.Addr(),
			Topic:   brokerCfg.Topic,
			Enabled: brokerCfg.Enabled,
		})
	helpers.FailOnError(t, err)

	helpers.FailOnError(t, prod.Close())
}
