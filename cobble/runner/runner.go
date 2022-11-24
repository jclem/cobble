package runner

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jclem/cobble/cobble/task"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	ScaffoldsDir string
	WorkingDir   string
	taskMap      map[string]*task.Task
	mutexMap     map[string]*sync.Mutex
}

type taskFile struct {
	Run   string
	Needs []string
	Touch []string
}

func (r *Runner) Run(taskNames ...string) error {
	var eg errgroup.Group

taskFor:
	for _, taskName := range taskNames {
		if strings.HasSuffix(taskName, "*") {
			tl, err := task.List()
			if err != nil {
				return err
			}

			prefix := strings.TrimSuffix(taskName, "*")

			for _, scaffold := range tl {
				if strings.HasPrefix(scaffold, prefix) {
					if err := r.addTask(scaffold); err != nil {
						return err
					}
				}
			}
			continue taskFor
		}

		if err := r.addTask(taskName); err != nil {
			return err
		}
	}

	for _, t := range r.taskMap {
		tt := t
		eg.Go(func() error {
			return tt.Run()
		})
	}

	return eg.Wait()
}

func (r *Runner) addTask(taskName string, from ...*task.Task) error {
	if t, ok := r.taskMap[taskName]; ok {
		t.AddChildren(from...)
		return nil
	}

	f, err := os.ReadFile(filepath.Join(r.ScaffoldsDir, taskName, "scaffold.yaml"))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if f == nil {
		f, err = os.ReadFile(filepath.Join(r.ScaffoldsDir, taskName, "scaffold.yml"))
		if err != nil {
			return err
		}
	}

	var tf taskFile
	if err := yaml.Unmarshal(f, &tf); err != nil {
		return err
	}

	mutexes := make([]*sync.Mutex, len(tf.Touch))

	for i, touch := range tf.Touch {
		if _, ok := r.mutexMap[touch]; !ok {
			r.mutexMap[touch] = &sync.Mutex{}
		}

		mutexes[i] = r.mutexMap[touch]
	}

	t, err := task.NewWithOpts(
		task.WithName(taskName),
		task.WithWorkingDir(r.WorkingDir),
		task.WithRun(tf.Run),
		task.WithDeps(tf.Needs...),
		task.WithChildren(from...),
		task.WithTouch(mutexes...),
	)

	r.taskMap[taskName] = t

	for _, dep := range tf.Needs {
		if err := r.addTask(dep, t); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

type Opt func(*Runner) error

func NewWithOpts(opts ...Opt) (*Runner, error) {
	runner := &Runner{
		taskMap:  make(map[string]*task.Task),
		mutexMap: make(map[string]*sync.Mutex),
	}

	for _, opt := range opts {
		if err := opt(runner); err != nil {
			return nil, err
		}
	}
	return runner, nil
}

func WithScaffoldsDir(scaffoldsDir string) Opt {
	return func(runner *Runner) error {
		runner.ScaffoldsDir = scaffoldsDir
		return nil
	}
}

func WithWorkingDir(workingDir string) Opt {
	return func(runner *Runner) error {
		runner.WorkingDir = workingDir
		return nil
	}
}
