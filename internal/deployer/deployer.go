package deployer

// Deployer - интерфейс взаимодействия с кибернетисом
type Deployer interface {
	CreatePod(name string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}
