package main

import (
	"github.com/mfsyahrz/image_feed_api/internal/interface/ioc"
	"github.com/mfsyahrz/image_feed_api/internal/interface/server/rest"
)

func main() {
	rest.StartRestServer(ioc.Setup())
}
