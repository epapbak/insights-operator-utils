/*
Copyright © 2020, 2021, 2022 Red Hat, Inc.

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

// This file contains the Kafka-related metrics that needs to be exposed to Prometheus
// Currently, the following metrics are exposed:
//
//   consumed_messages - total number of messages consumed from selected broker
//
//   consuming_errors - total number of errors during consuming messages from selected broker
//
//   produced_messages - total number of produced messages sent to Payload Tracker's Kafka topic

package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ProducedMessages shows number of messages produced by producer
var ProducedMessages = promauto.NewCounter(prometheus.CounterOpts{
	Name: "produced_messages",
	Help: "The total number of produced messages sent to Payload Tracker's Kafka topic",
})
