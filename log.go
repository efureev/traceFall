package tracefall

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Environments
const (
	EnvironmentDev  = `dev`
	EnvironmentProd = `prod`
	EnvironmentTest = `test`
)

// Logable interface
type Logable interface {
	Success() Logable
	Fail(err error) Logable
	SetParentID(id uuid.UUID) Logable
	SetParent(parent Logable) error
	CreateChild(name string) (Logable, error)
	ToJSON() []byte
	ToLogJSON() LogJSON
}

// LogJSON struct
type LogJSON struct {
	ID          uuid.UUID     `json:"id"`
	Thread      uuid.UUID     `json:"thread"`
	Name        string        `json:"name"`
	App         string        `json:"app"`
	Time        int64         `json:"time"`
	TimeEnd     *int64        `json:"timeEnd"`
	Result      bool          `json:"result"`
	Finish      bool          `json:"finish"`
	Environment string        `json:"env"`
	Error       *string       `json:"error"`
	Data        ExtraData     `json:"data"`
	Notes       NoteGroupList `json:"notes"`
	Tags        []string      `json:"tags"`
	Parent      *string       `json:"parent"`
	//Step        uint16       `json:"step"`
}

// Log struct
type Log struct {
	ID          uuid.UUID
	Thread      uuid.UUID
	Name        string
	Data        ExtraData
	App         string
	Notes       NoteGroups
	Tags        Tags
	Error       error
	Environment string
	//Step        uint16

	Result bool
	Finish bool

	Time    time.Time
	TimeEnd *time.Time
	Parent  *Log
	//items   []*Log
}

// SetName set log name
func (l *Log) SetName(name string) *Log {
	l.Name = name
	return l
}

// FinishTimeEnd set finish time of the log
func (l *Log) FinishTimeEnd() *Log {
	n := time.Now()
	l.TimeEnd = &n
	return l
}

// ThreadFinish finish thread line
func (l *Log) ThreadFinish() *Log {
	l.Finish = true
	return l
}

// Success set result of the log: success
func (l *Log) Success() *Log {
	l.FinishTimeEnd().Result = true
	return l
}

// Fail set result of the log: error
func (l *Log) Fail(err error) *Log {
	l.Result = false
	l.Error = err
	return l.FinishTimeEnd()
}

// SetEnvironment set environment name of log
func (l *Log) SetEnvironment(env string) *Log {
	l.Environment = env
	return l
}

// SetApplication set application name of log
func (l *Log) SetApplication(str string) *Log {
	l.App = str
	return l
}

// ErrorParentFinish error
var ErrorParentFinish = errors.New(`the Parent does not have to be the finish point`)

// ErrorParentThreadDiff error
var ErrorParentThreadDiff = errors.New(`the Parent Thread is different from the Thread of own log`)

// SetParent set parent to log for for create Thread
func (l *Log) SetParent(parent *Log) error {
	if parent.Finish {
		return ErrorParentFinish
	}

	if parent.Thread.String() != l.Thread.String() {
		return ErrorParentThreadDiff
	}

	if parent != nil {
		l.Parent = parent
		//parent.items = append(parent.items, l)
	}

	return nil
}

// SetParentID set parent ID to log
func (l *Log) SetParentID(id uuid.UUID) *Log {
	l.Parent = &Log{ID: id, Thread: l.Thread}
	return l
}

// CreateChild make new log and attach it to current log as child
func (l *Log) CreateChild(name string) (*Log, error) {
	if l.Finish {
		return nil, ErrorParentFinish
	}
	child := NewLog(name)
	child.Thread = l.Thread
	child.App = l.App
	child.Environment = l.Environment
	child.Parent = l

	return child, nil
}

// ToJSON create json bytes from Log data
func (l Log) ToJSON() []byte {
	b, _ := l.MarshalJSON()
	return b
}

// MarshalJSON marshal json
func (l *Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.ToLogJSON())
}

// ToLogJSON return JsonLog Struct
func (l Log) ToLogJSON() *LogJSON {
	var (
		parentID, er *string
		te           *int64
	)
	if l.Parent != nil {
		pid := l.Parent.ID.String()
		parentID = &pid
	} else {
		parentID = nil
	}

	if l.TimeEnd != nil {
		teInt := l.TimeEnd.UnixNano()
		te = &teInt
	}

	if l.Error != nil {
		e1 := l.Error.Error()
		er = &e1
	}

	return &LogJSON{
		ID:          l.ID,
		Thread:      l.Thread,
		Name:        l.Name,
		App:         l.App,
		Time:        l.Time.UnixNano(),
		TimeEnd:     te,
		Result:      l.Result,
		Finish:      l.Finish,
		Environment: l.Environment,
		Error:       er,
		Data:        l.Data,
		Notes:       l.Notes.prepareToJSON(),
		Tags:        l.Tags,
		Parent:      parentID,
		//Step : l.Step,
	}
}

// String return string representation of log
func (l Log) String() string {
	return fmt.Sprintf("[%s] %s", l.Time, l.Name)
}

// SetDefaults set values for Log by default
func (l *Log) SetDefaults() *Log {
	l.App = `App`
	l.Environment = EnvironmentDev
	l.Result = false
	return l
}

// NewLog create new Log
func NewLog(name string) *Log {
	id := generateUUID()
	return (&Log{
		ID:     id,
		Thread: id,
		Name:   name,
		Data:   NewExtraData(),
		Notes:  NewNotesGroups(),
		Result: false,
		Tags:   Tags{},
		Time:   time.Now(),
	}).SetDefaults()
}

// LogParentShadow struct
type LogParentShadow struct {
	ID     uuid.UUID `json:"id"`
	Thread uuid.UUID `json:"thread"`
}

// ToShadow create new shadow struct of log
func (l Log) ToShadow() *LogParentShadow {
	return &LogParentShadow{l.ID, l.Thread}
}

// ParentFromShadow return Parent's ID from LogShadow
func (l *Log) ParentFromShadow(shadow *LogParentShadow) *Log {
	if shadow != nil {
		l.Parent = &Log{ID: shadow.ID, Thread: shadow.Thread}
		l.Thread = shadow.Thread
	}
	return l
}

const rootLevel int = 0

// GetLevel return depth level of Log in thread
func (l *Log) GetLevel() int {
	current := l
	level := rootLevel
	for {
		if current.Parent == nil {
			break
		}
		level++
		current = current.Parent
	}

	return level
}
