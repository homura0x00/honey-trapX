package docker

import (
	"context"
	"fmt"
	"io"
	"net/netip"
	"strconv"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// Manager 唯一接触 docker daemon 的封装。moby 类型不出此包。
type Manager struct {
	cli *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %w", err)
	}
	return &Manager{cli: cli}, nil
}

// EnsureImage 拉取镜像并消费进度流，确保 pull 真正完成。
func (m *Manager) EnsureImage(ctx context.Context, image string) error {
	rc, err := m.cli.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("拉取镜像 %s: %w", image, err)
	}
	defer rc.Close()
	if _, err := io.Copy(io.Discard, rc); err != nil {
		return fmt.Errorf("读取镜像拉取进度: %w", err)
	}
	return nil
}

// CreateContainer 创建容器（不自动删除，支持后续 stop→start），返回容器 ID。
func (m *Manager) CreateContainer(ctx context.Context, name, image string, hostPort, containerPort int) (string, error) {
	p, err := network.ParsePort(fmt.Sprintf("%d/tcp", containerPort))
	if err != nil {
		return "", fmt.Errorf("非法容器端口 %d: %w", containerPort, err)
	}

	cfg := &container.Config{
		Image:        image,
		ExposedPorts: network.PortSet{p: {}},
	}
	hostCfg := &container.HostConfig{
		AutoRemove: false,
		PortBindings: network.PortMap{
			p: {{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: strconv.Itoa(hostPort)}},
		},
	}

	resp, err := m.cli.ContainerCreate(ctx, cfg, hostCfg, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("创建容器失败: %w", err)
	}
	return resp.ID, nil
}

func (m *Manager) StartContainer(ctx context.Context, containerID string) error {
	if err := m.cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("启动容器失败: %w", err)
	}
	return nil
}

func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	if err := m.cli.ContainerStop(ctx, containerID, client.ContainerStopOptions{}); err != nil {
		return fmt.Errorf("停止容器失败: %w", err)
	}
	return nil
}

func (m *Manager) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	if err := m.cli.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{Force: force}); err != nil {
		return fmt.Errorf("删除容器失败: %w", err)
	}
	return nil
}
