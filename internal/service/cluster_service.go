package service

import (
	"context"
	"fmt"
	"time"

	"ClusterGuard/internal/domain"

	"github.com/google/uuid"
)

// ClusterService handles cluster CRUD (Single Responsibility).
type ClusterService struct {
	repo domain.ClusterRepository
}

func NewClusterService(repo domain.ClusterRepository) *ClusterService {
	return &ClusterService{repo: repo}
}

type ClusterInput struct {
	ID          string
	Name        string
	Description string
	ServerIDs   []string
}

func (s *ClusterService) Create(ctx context.Context, in ClusterInput) (*domain.Cluster, error) {
	now := time.Now().Unix()
	cluster := &domain.Cluster{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: in.Description,
		ServerIDs:   in.ServerIDs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if cluster.ServerIDs == nil {
		cluster.ServerIDs = []string{}
	}
	if err := s.repo.Save(ctx, cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *ClusterService) Update(ctx context.Context, in ClusterInput) (*domain.Cluster, error) {
	if in.ID == "" {
		return nil, fmt.Errorf("cluster id is required")
	}
	existing, err := s.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = in.Name
	existing.Description = in.Description
	if in.ServerIDs != nil {
		existing.ServerIDs = in.ServerIDs
	}
	existing.UpdatedAt = time.Now().Unix()

	if err := s.repo.Save(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ClusterService) Get(ctx context.Context, id string) (*domain.Cluster, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ClusterService) List(ctx context.Context) ([]*domain.Cluster, error) {
	return s.repo.List(ctx)
}

func (s *ClusterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
