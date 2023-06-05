package httputil

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ShutdownGracefully(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}

	return nil
}
