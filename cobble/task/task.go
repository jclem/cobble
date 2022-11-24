package task

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var notifStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#99ff99"))
var cmdStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9999ff"))

type Task struct {
	Name  string
	Ready chan struct{}

	workingDir string
	runCmd     string
	depsCount  int
	children   []*Task
	touch      []*sync.Mutex
}

func (t *Task) Run() error {
	readyCount := 0

	for {
		select {
		case <-t.Ready:
			readyCount++
			if readyCount == t.depsCount {
				return t.run()
			}
		default:
			if readyCount == t.depsCount {
				return t.run()
			}
		}
	}
}

func (t *Task) AddChildren(children ...*Task) {
	t.children = append(t.children, children...)
}

func (t *Task) run() error {
	for _, touch := range t.touch {
		touch.Lock()
		defer touch.Unlock()
	}

	fmt.Println(notifStyle.Render("Running task:"), cmdStyle.Render(t.Name))

	c := exec.Command("/bin/sh", "-c", t.runCmd)
	c.Stdout = os.Stdout
	c.Stderr = c.Stdout
	c.Env = os.Environ()
	err := c.Run()

	for _, child := range t.children {
		child.Ready <- struct{}{}
	}

	return err
}

type Opt func(*Task) error

func NewWithOpts(opts ...Opt) (*Task, error) {
	task := &Task{Ready: make(chan struct{})}
	for _, opt := range opts {
		if err := opt(task); err != nil {
			return nil, err
		}
	}
	return task, nil
}

func WithChildren(children ...*Task) Opt {
	return func(task *Task) error {
		task.children = children
		return nil
	}
}

func WithDeps(deps ...string) Opt {
	return func(task *Task) error {
		task.depsCount = len(deps)
		return nil
	}
}

func WithName(name string) Opt {
	return func(task *Task) error {
		task.Name = name
		return nil
	}
}

func WithRun(run string) Opt {
	return func(task *Task) error {
		task.runCmd = run
		return nil
	}
}

func WithTouch(touch ...*sync.Mutex) Opt {
	return func(task *Task) error {
		task.touch = touch
		return nil
	}
}

func WithWorkingDir(workingDir string) Opt {
	return func(task *Task) error {
		task.workingDir = workingDir
		return nil
	}
}
