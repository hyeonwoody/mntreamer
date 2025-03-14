package lib

import (
	"fmt"
	"mntreamer/monolithic/cmd/configuration"
)

func Start(ctnr *configuration.MonolithicContainer) error {

	errChan := make(chan error, 1)
	go func() {
		if err := ctnr.RunRouter(); err != nil {
			errChan <- fmt.Errorf("Router failed to run :%w", err)
		}
	}()

	return <-errChan
}
