package polling

import (
	"context"
	"fmt"

	"github.com/flux-iac/tofu-controller/internal/config"
)

const DefaultConfigMapName = "branch-planner"

func (s *Server) readConfig(ctx context.Context) (*config.Config, error) {
	configMap, err := config.ReadConfig(ctx, s.clusterClient, s.configMapRef)
	if err != nil {
		return nil, fmt.Errorf("unable to read ConfigMap: %w", err)
	}

	return &configMap, nil
}
