package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"
)

// NewSpring creates new transient kafka pipe
func NewSpring(topic string, tail int, commitAfter int, config *kafka.ConfigMap, formatter UnmarshalFunction) (spring kepler.Spring, err error) {

	log.Debug("creating KafkaSpring")
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Errorf("Failed to create consumer: %s\n", err)
		return
	}

	log.Debug("Kafka consumer created")
	log.Debug("Subscribing...")
	err = c.Subscribe(topic, nil)

	if err != nil {
		log.Error("Failed to Subscribe: %s\n", err)
		c.Close()
	}

	counter := 0
	//var delta Delta

	log.Info("Subscribed to kafka")
	spring = kepler.NewSpring(func(ctx context.Context, out chan<- kepler.Message) {
		for {
			select {
			case ev := <-c.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					log.Debugln("Assigning...")
					err = c.Assign(e.Partitions)
					log.Debugln("Current partitions")
					for _, ap := range e.Partitions {
						log.Debugf("%s[%d]@%v \n", *ap.Topic, ap.Partition, ap.Offset)
					}

					if tail > 0 {
						log.Println("ReAssigning tail partition offsets...")
						parts := make([]kafka.TopicPartition, len(e.Partitions))
						for i, tp := range e.Partitions {
							parts[i] = kafka.TopicPartition{Topic: tp.Topic, Partition: tp.Partition, Offset: kafka.OffsetTail(50)}
						}

						log.Println("New partitions")
						for _, ap := range parts {
							log.Printf("%s[%d]@%v \n", *ap.Topic, ap.Partition, ap.Offset)
						}

						log.Println("Assigning...")
						err = c.Assign(parts)
						if err != nil {
							panic(err)
						}
					}
					log.Printf("[Done]")
				case kafka.RevokedPartitions:
					c.Unassign()
				case *kafka.Message:
					log.Infof("Partition loaded: %v\n", e.TopicPartition)
					var msg kepler.Message
					if msg, err = formatter(e); err != nil {
						log.Error("Error: %v", err)
						log.Error("Skip invalid delta: %v", string(e.Value))
						continue
					}

					out <- msg
					counter++
					if commitAfter > 0 && counter > commitAfter {
						c.Commit()
						counter = 0
					}
				case kafka.PartitionEOF:
					log.Printf("EOF Reached %v\n", e)
				case kafka.Error:
					log.Printf("Error: %v\n", e)
					break
				}
			case <-ctx.Done():
				log.Info("kafka spring done")
				close(out)
			}
		}
	})

	return
}
