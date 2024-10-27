package rest

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo"
	"github.com/mfsyahrz/image_feed_api/internal/interface/ioc"
)

func StartRestServer(container *ioc.IOC) {

	var (
		httpAddr = fmt.Sprintf(":%s", container.Config.Service.Port.REST)
		errChan  = make(chan error)
		e        = echo.New()
	)

	SetupMiddleware(e, container)
	SetupRouter(e, SetupHandler(container))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errChan <- e.Start(httpAddr)
	}()

	err := <-errChan
	panic(err.Error())

}
