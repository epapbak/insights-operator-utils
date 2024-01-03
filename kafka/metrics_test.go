// Copyright 2020, 2021, 2022, 2023 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka_test

import (
	"fmt"
	"github.com/RedHatInsights/insights-operator-utils/kafka"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/Shopify/sarama/mocks"
	"github.com/prometheus/client_golang/prometheus"
	prommodels "github.com/prometheus/client_model/go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func assertCounterValue(tb testing.TB, expected int64, counter prometheus.Counter, initValue int64) {
	assert.Equal(tb, float64(expected+initValue), getCounterValue(counter))
}

func getCounterValue(counter prometheus.Counter) float64 {
	pb := &prommodels.Metric{}
	err := counter.Write(pb)
	if err != nil {
		panic(fmt.Sprintf("Unable to get counter from counter %v", err))
	}

	return pb.GetCounter().GetValue()
}

func TestProducedMessagesMetric(t *testing.T) {
	// other tests may run at the same process
	initValue := int64(getCounterValue(kafka.ProducedMessages))

	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	p := kafka.KafkaProducer{Producer: mockProducer}
	defer func() {
		helpers.FailOnError(t, p.Close())
	}()

	_, _, err := p.ProduceMessage([]byte("test"), "producer_topic")
	helpers.FailOnError(t, err)

	assertCounterValue(t, 1, kafka.ProducedMessages, initValue)
}
