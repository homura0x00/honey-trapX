package honey

import "context"

// ContainerOps honey 需要的容器操作窄接口，定义在消费者侧；
// *docker.Manager 隐式实现，测试时注入 fake。
type ContainerOps interface {
	EnsureImage(ctx context.Context, image string) error
	CreateContainer(ctx context.Context, name, image string, hostPort, containerPort int) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string, force bool) error
}
