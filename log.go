package traceFall

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

const (
	EnvironmentDev  = `dev`
	EnvironmentProd = `prod`
	EnvironmentTest = `test`
)

type Log struct {
	Id          uuid.UUID
	Thread      uuid.UUID
	Name        string
	Data        ExtraData
	App         string
	Notes       NoteGroups
	Tags        Tags
	Error       error
	Environment string
	step        uint16

	Result bool
	Finish bool

	Time    time.Time
	TimeEnd *time.Time
	Parent  *Log
}

func (l *Log) SetName(name string) *Log {
	l.Name = name
	return l
}

func (l *Log) AddData(key string, val interface{}) *Log {
	l.Data.Set(key, val)
	return l
}

func (l *Log) AddNote(step, note string) *Log {
	l.Notes.Add(step, note)
	return l
}

// set finish time of the log
func (l *Log) FinishTimeEnd() *Log {
	n := time.Now()
	l.TimeEnd = &n
	return l
}

// finish thread line
func (l *Log) ThreadFinish() *Log {
	l.Finish = true
	return l
}

// set result of the log: success
func (l *Log) Success() *Log {
	l.FinishTimeEnd().Result = true
	return l
}

// set result of the log: error
func (l *Log) Fail(err error) *Log {
	l.Result = false
	l.Error = err
	return l.FinishTimeEnd()
}

func (l *Log) SetEnvironment(env string) *Log {
	l.Environment = env
	return l
}

var ErrorParentFinish = errors.New(`the Parent does not have to be the finish point`)
var ErrorParentThreadDiff = errors.New(`the Parent Thread is different from the Thread of own log`)

func (l *Log) SetParent(parent *Log) error {
	if parent.Finish {
		return ErrorParentFinish
	}

	if parent.Thread.String() != l.Thread.String() {
		return ErrorParentThreadDiff
	}

	l.Parent = parent

	return nil
}

func (l *Log) SetParentId(id uuid.UUID) *Log {
	l.Parent = &Log{Id: id, Thread: l.Thread}
	return l
}

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

func (l Log) String() string {
	return fmt.Sprintf("[%s] %s", l.Time, l.Name, )
}

func (l *Log) SetDefaults() *Log {
	l.App = `App`
	l.Environment = EnvironmentDev
	l.Result = false
	return l
}

func NewLog(name string) *Log {
	id := generateUUID()
	return (&Log{
		Id:     id,
		Thread: id,
		Name:   name,
		Data:   NewExtraData(),
		Notes:  NewNotesGroups(),
		Result: false,
		Tags:   Tags{},
		Time:   time.Now(),
	}).SetDefaults()
}

/*



type LogJson struct {
	Id          uuid.UUID      `json:"id"`
	Thread      uuid.UUID      `json:"thread"`
	Name        string         `json:"name"`
	App         string         `json:"app"`
	Module      string         `json:"module"`
	Time        int64          `json:"time"`
	TimeEnd     *int64         `json:"timeEnd"`
	Result      bool           `json:"result"`
	Finish      bool           `json:"finish"`
	Environment string         `json:"env"`
	Error       *string        `json:"error"`
	Data        LogExtraParams `json:"data"`
	Notes       LogNoteGroups  `json:"notes"`
	Tags        []string       `json:"tags"`
	Parent      *string        `json:"parent"`
}




func (l *Log) CreateChildParams(name string, params LogExtraParams) Log {
	return Log{
		Id:          generateUUID(),
		thread:      l.thread,
		Name:        name,
		App:         l.App,
		notes:       NewNotesGroups(),
		Data:        params,
		Result:      false,
		Time:        time.Now(),
		Parent:      l,
		Tags:        removeDuplicates(append(l.Tags, l.Id.String(), l.thread.String())),
		Environment: l.Environment,
	}
}

func (l Log) ToJsonByte() []byte {
	b, _ := l.MarshalJSON()
	return b
}

func (l *Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.ToLogJson())
}


func (l *Log) ToLogJson() LogJson {
	var (
		parentId, er *string
		te           *int64
	)
	if l.Parent != nil {
		pid := l.Parent.Id.String()
		parentId = &pid
	} else {
		parentId = nil
	}

	if l.TimeEnd != nil {
		teInt := l.TimeEnd.UnixNano()
		te = &teInt
	}

	if l.Error != nil {
		e1 := l.Error.Error()
		er = &e1
	}

	return LogJson{
		Id:          l.Id,
		Thread:      l.thread,
		Name:        l.Name,
		App:         l.App,
		Time:        l.Time.UnixNano(),
		TimeEnd:     te,
		Result:      l.Result,
		Finish:      l.Finish,
		Environment: l.Environment,
		Error:       er,
		Data:        l.Data,
		Notes:       l.notes,
		Tags:        l.Tags,
		Parent:      parentId,
	}
}



type LogShadow struct {
	Id     uuid.UUID `json:"id"`
	Thread uuid.UUID `json:"thread"`
}

func (l Log) ToShadow() LogShadow {
	return LogShadow{l.Id, l.thread}
}

func (l *Log) SetParentFromShadow(shadow *LogShadow) *Log {
	if shadow != nil {
		l.Parent = &Log{Id: shadow.Id, thread: shadow.Thread}
		l.thread = shadow.Thread
	}
	return l
}
*/
/*

func (l Log) GetLevel() int {
	current := l
	level := 0
	for {
		if current.Parent == nil {
			break
		}
		level++
		current = *current.Parent
	}

	return level
}

func (l Log) BuildTree() string {

	current := l
	level := l.GetLevel()
	var text string

	for {
		offset := strings.Repeat("\t", level)

		var parentId string
		if current.Parent != nil {
			parentId = current.Parent.Id.String()
		}
		t := fmt.Sprintf("{{offset}}- id: %s\n{{offset}}- time: %s\n{{offset}}- action: %s\n{{offset}}- level: %d\n{{offset}}- parent: %s\n", current.Id, current.Time, current.Action, level, parentId)
		text += strings.Replace(t, `{{offset}}`, offset, -1)

		if current.Parent == nil {
			break
		}
		level--
		current = *current.Parent
	}

	return text
}
*/
