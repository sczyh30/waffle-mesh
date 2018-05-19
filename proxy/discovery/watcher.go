package discovery

import (
	"time"
	"errors"
	"context"
	"log"
)

type ResourceWatcher struct {
	XdsConsumer *XdsConsumer
}

func (w *ResourceWatcher) StartWatching(stop chan struct{}) error {
	defer w.XdsConsumer.conn.Close()

	if w.XdsConsumer == nil {
		return errors.New("unexpected state for resource discovery service watcher")
	}
	ticker := time.NewTicker(time.Second * 20)
	ctx := context.Background()
	select {
	case <-stop:
		log.Println("Stopping the xDS watcher")
		break
	case <-ticker.C:
		go func() {
			err := w.XdsConsumer.RetrieveAndUpdate(ctx)
			if err != nil {
				log.Println("Error when retrieving xDS resources: " + err.Error())
			}
		}()
	}
	return nil
}
