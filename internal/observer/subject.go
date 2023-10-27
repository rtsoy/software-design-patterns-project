package observer

type Subject interface {
	Subscribe(id int64) error
	Unsubscribe(id int64) error
	NotifyAll() ([]int64, error)
}
