package html

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/BeamFoundry/rays/pkg/lang"

	core "github.com/gofiber/template/v2"
	"github.com/gofiber/utils"
)

type templateKey struct {
	key string
}

func (this templateKey) New(name string, layouts []string) templateKey {
	templates := []string{name}
	templates = append(templates, layouts...)
	slices.Sort(templates)
	this.key = strings.Join(templates, "_")
	return this
}

// Engine struct
type Engine struct {
	core.Engine
	// templates
	Templates   *template.Template
	TemplateMap map[templateKey]*template.Template
}

// New returns a HTML render engine for Fiber
func New(directory, extension string) *Engine {
	return newEngine(directory, extension, nil)
}

// NewFileSystem returns a HTML render engine for Fiber with file system
func NewFileSystem(fs http.FileSystem, extension string) *Engine {
	return newEngine("/", extension, fs)
}

// newEngine creates a new Engine instance with common initialization logic.
func newEngine(directory, extension string, fs http.FileSystem) *Engine {
	engine := &Engine{
		Engine: core.Engine{
			Left:       "{{",
			Right:      "}}",
			Directory:  directory,
			FileSystem: fs,
			Extension:  extension,
			Funcmap:    make(map[string]any),
		},
		TemplateMap: map[templateKey]*template.Template{},
	}
	addfunc := func(name string) {
		f := func() error {
			return errors.New(name + " called unexpectedly")
		}
		if err := f(); err != nil {
			engine.AddFunc(name, f)
		}
	}
	addfunc("include")
	addfunc("toHTML")
	return engine
}

func (e *Engine) parseTemplate(templates *template.Template, render string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		return Try(func() {
			// Return error if exist
			RaiseErr(err)

			// Get the relative file path
			// ./views/html/index.tmpl -> index.tmpl
			rel := RaiseVal(filepath.Rel(e.Directory, path))

			// Skip file if it's a directory or has no file info
			if info == nil || info.IsDir() {
				return
			}

			// Skip file if it does not equal the given template Extension
			if len(e.Extension) >= len(path) || path[len(path)-len(e.Extension):] != e.Extension {
				return
			}

			// Reverse slashes '\' -> '/' and
			// partials\footer.tmpl -> partials/footer.tmpl
			name := filepath.ToSlash(rel)
			// Remove ext from name 'index.tmpl' -> 'index'
			name = strings.TrimSuffix(name, e.Extension)

			if render != "" && render != name {
				return
			}

			// Read the file
			// #gosec G304
			buf := RaiseVal(utils.ReadFile(path, e.FileSystem))

			// Create new template associated with the current one
			// This enable use to invoke other templates {{ template .. }}
			RaiseVal(templates.New(name).Parse(string(buf)))

			// Debugging
			var err Error = NewError("Debug", "views: parsed template: "+name)
			if !e.Verbose {
				err = &NilError{}
			}
			err.Log(slog.LevelInfo)
		})
	}

	return Try(func() {
		if e.FileSystem != nil {
			RaiseErr(utils.Walk(e.FileSystem, e.Directory, walkFn))
		} else {
			RaiseErr(filepath.Walk(e.Directory, walkFn))
		}
	})
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	if e.Loaded {
		return nil
	}

	// race safe
	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	return Try(func() {
		e.Templates = template.New(e.Directory)
		e.Templates.Delims(e.Left, e.Right)
		e.Templates.Funcs(e.Funcmap)
		RaiseErr(e.parseTemplate(e.Templates, ""))

		// notify Engine that we parsed all templates
		e.Loaded = true
	})
}

func (e *Engine) getTemplates(name string, layouts []string) (*template.Template, error) {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	index := templateKey{}.New(name, layouts)

	if t, ok := e.TemplateMap[index]; ok {
		return t, nil
	}

	templates := template.New(e.Directory)
	templates.Delims(e.Left, e.Right)
	templates.Funcs(e.Funcmap)
	var funcMap template.FuncMap = map[string]any{}
	// Returning an HTML object, to be piped to "toString" to be handled
	funcMap["include"] = func(name string, data any) (template.HTML, error) {
		tmpl := template.HTML("")
		return tmpl, Try(func() {
			buf := bytes.NewBuffer(nil)
			RaiseErr(templates.ExecuteTemplate(buf, name, data))
			//nolint:gosec // Server side template
			tmpl = template.HTML(buf.String())
		})
	}
	// Returning an HTML object piped into "toString" back as an HTML object
	funcMap["toHTML"] = func(s string) template.HTML {
		//nolint:gosec // Server side template
		return template.HTML(s)
	}
	templates.Funcs(funcMap)

	return templates, Try(func() {
		for _, l := range layouts {
			RaiseErr(e.parseTemplate(templates, l))
		}

		RaiseErr(e.parseTemplate(templates, name))

		e.TemplateMap[index] = templates
	})
}

func (e *Engine) Render(out io.Writer, name string, binding any, layout ...string) error {
	valErr := func(val any, err error) error {
		ret := err
		if val != nil {
			ret = nil
		}
		return ret
	}
	return Try(func() {
		// Check if templates need to be loaded/reloaded
		if e.PreRenderCheck() {
			RaiseErr(e.Load())
		}

		e.Mutex.RLock()
		tmpl := e.Templates.Lookup(name)
		e.Mutex.RUnlock()

		RaiseErr(valErr(tmpl, fmt.Errorf("render: template %s does not exist", name)))

		layouts := []string{}

		for _, layName := range layout {
			if layName == "" {
				break
			}
			e.Mutex.RLock()
			lay := e.Templates.Lookup(layName)
			e.Mutex.RUnlock()
			RaiseErr(valErr(lay, fmt.Errorf("render: LayoutName %s does not exist", layName)))
			layouts = append(layouts, layName)
		}

		templates := RaiseVal(e.getTemplates(name, layouts))
		RaiseErr(templates.Lookup(name).Execute(out, binding))
	})
}
