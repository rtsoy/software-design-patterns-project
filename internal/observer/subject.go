package observer

type Subject interface {
	Subscribe(id uint) error
	Unsubscribe(id uint) error
	NotifyAll() ([]uint, error)
}
