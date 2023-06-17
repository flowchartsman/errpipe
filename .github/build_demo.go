package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"
)

func buildDemos(action ActionFn) error {
	return RunMatrix(action, 0,
		Layer.WithName("errpipe")(
			vars{
				"DemoName":  "errpipe",
				"Framerate": 120,
			},
			Layer.WithName("full")(
				vars{
					"Width":  800,
					"Height": 400,
					"Env": vars{
						"ERRPIPE_MAX":   3,
						"ERRPIPE_STYLE": "Braille",
					},
					"FontSize": 28,
					"Pre":      []string{"Show"},
					"Post":     []string{"Sleep 500ms"},
				},
				Layer.WithName("showcase")(
					vars{
						"DemoLength": 25,
						"WindowBar":  "Colorful",
						"DemoArgs":   []string{"--demo delay"},
						"suffix":     "",
					},
				),
				Layer.WithName("with-warnings")(
					vars{
						"DemoLength": 17,
						"DemoArgs":   []string{"--demo main"},
						"Suffix":     "with-warnings",
						"EPArgs":     []string{"-w"},
					},
				),
			),
			Layer.WithName("styles")(
				vars{
					"Width":      800,
					"Height":     400,
					"TrimWidth":  277,
					"FontSize":   20,
					"DemoLength": 14,
					"Env": vars{
						"ERRPIPE_IDLE": 0,
						"ERRPIPE_MAX":  4,
					},
					"DemoArgs": []string{"--demo short"},
					"Pre": []string{
						"Set TypingSpeed 1ms",
						// "Show",
					},
					"Post": []string{
						"Show",
					},
					"PostProcess": []string{
						`{{$file := print ".github/" .DemoName "-" .Suffix ".gif" }}gifsicle --crop 12,23+{{.TrimWidth}}x25 -o {{$file}} {{$file}}`,
					},
				},
				Layer.WithName("single-width")(
					vars{},
					Layer.WithName("double-height")(
						vars{
							"DemoArgs": []string{"--demo tall"},
							"Env": vars{
								"ERRPIPE_MAX": 8,
							},
						},
						Layer.WithName("block")(
							vars{
								"Style":     "block",
								"Suffix":    "style-block",
								"TrimWidth": 279,
								"Env": vars{
									"ERRPIPE_STYLE": "block",
								},
							},
						),
					),
					Layer.WithName("single-height")(
						vars{},
						Layer.WithName("legacy-line")(
							vars{
								"Suffix": "style-legacy",
								"Env": vars{
									"ERRPIPE_STYLE": "legacy",
								},
							},
						),
						Layer.WithName("legacy-block")(
							vars{
								"Suffix": "style-legacy-block",
								"Env": vars{
									"ERRPIPE_STYLE": "legacy-block",
								},
							},
						),
						Layer.WithName("legacy-block-line")(
							vars{
								"Suffix": "style-legacy-block-line",
								"Env": vars{
									"ERRPIPE_STYLE": "legacy-block-line",
								},
							},
						),
					),
				),
				Layer.WithName("double-width")(
					vars{
						"DemoArgs": []string{"--demo double"},
						"NoCursor": true,
					},
					Layer.WithName("braille")(
						vars{
							"Suffix": "style-braille",
							"Env": vars{
								"ERRPIPE_STYLE": "braille",
							},
						},
					),
					Layer.WithName("braille-line")(
						vars{
							"Suffix": "style-braille-line",
							"Env": vars{
								"ERRPIPE_STYLE": "braille-line",
							},
						},
					),
					Layer.WithName("braille4")(
						vars{
							"Suffix": "style-braille4",
							"Env": vars{
								"ERRPIPE_STYLE": "braille4",
							},
						},
					),
					Layer.WithName("braille4-line")(
						vars{
							"Suffix": "style-braille4-line",
							"Env": vars{
								"ERRPIPE_STYLE": "braille4-line",
							},
						},
					),
				),
			),
		),
	)
}

func main() {
	log.SetFlags(0)
	tmpl := mustGet(template.New("").Funcs(template.FuncMap{"Args": argsJoin}).ParseFiles("errpipe.tape.tmpl"))
	os.Chdir("..")
	tmpdir := mustGet(os.MkdirTemp("", "errpipebuild"))
	defer os.RemoveAll(tmpdir) // really matters for local
	mustRun("go", "build", "-o", filepath.Join(tmpdir, "errpipe"))
	if err := buildDemos(runVHS(tmpl, tmpdir)); err != nil {
		log.Print(err)
		defer os.Exit(1)
		return
	}
}

func mustRun(cmdargs ...string) {
	if err := runCmd(cmdargs...); err != nil {
		log.Fatal(err) //nolint
	}
}

func runCmd(cmdargs ...string) error {
	var buf bytes.Buffer
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running %q:\n%s", cmd.String(), buf.String())
	}
	return nil
}

func mustGet[OUT any](o OUT, err error) OUT {
	if err != nil {
		log.Fatal(err) //nolint
	}
	return o
}

func runVHS(tmpl *template.Template, tmpdir string) ActionFn {
	return func(ctx context.Context, vars vars) error {
		vars["EpLoc"] = tmpdir
		cmd := exec.CommandContext(ctx, "vhs")
		pr, pw := io.Pipe()
		var errBuf bytes.Buffer
		cmd.Stdin = pr
		cmd.Stdout = &errBuf
		var tmplBuf bytes.Buffer
		err := tmpl.ExecuteTemplate(&tmplBuf, "errpipe.tape.tmpl", vars)
		if err != nil {
			return fmt.Errorf("vhs tmpl execute: %v", err)
		}
		go func() {
			io.Copy(pw, &tmplBuf)
			pw.Close()
		}()
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("vhs: %v\n%s", err, errBuf.String())
		}
		if pp, found := vars["PostProcess"]; found {
			for _, cmdStr := range pp.([]string) {
				cmdTmpl, err := tmpl.New("").Parse(cmdStr)
				if err != nil {
					return fmt.Errorf("post-process tmpl parse: %v", err)
				}
				var buf bytes.Buffer
				if err := cmdTmpl.Execute(&buf, vars); err != nil {
					return fmt.Errorf("post-process tmpl execute: %v", err)
				}
				cmd := strings.Split(buf.String(), " ")
				mustRun(cmd...)
			}
		}
		return nil
	}
}

func argsJoin(args []string) string {
	return " " + strings.Join(args, " ")
}

// MATRIX Code

func RunMatrix(action ActionFn, maxJobs int, root *layer) error {
	if root.name == "" {
		root.name = "[plan]"
	}
	p := plan{}
	buildPlan(path{root}, nil, &p)

	ctx, cf := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cf()
	var wg sync.WaitGroup
	var sem chan struct{}
	if maxJobs > 0 {
		sem = make(chan struct{}, maxJobs)
		for i := 0; i < maxJobs; i++ {
			sem <- struct{}{}
		}
	}
	var runErr error
	var errOnce sync.Once
	var failed atomic.Bool
	for i := range p {
		wg.Add(1)
		go func(run layer) {
			defer wg.Done()
			if sem != nil {
				<-sem
				defer func() {
					sem <- struct{}{}
				}()
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			if failed.Load() {
				return
			}
			log.Printf("RUN - %s", run.name)
			if err := action(ctx, run.vars); err != nil {
				errOnce.Do(func() {
					runErr = fmt.Errorf("%s: %w", run.name, err)
					failed.Store(true)
					cf()
				})
			}
		}(p[i])
	}
	wg.Wait()
	return runErr
}

type ActionFn func(ctx context.Context, vars vars) error

func buildPlan(path path, vars vars, plan *plan) {
	vars = merge(vars, path.current().vars)
	if path.bottom() {
		plan.Append(path.current().name, vars)
		return
	}

	for i, child := range path.current().children {
		if child.name == "" {
			child.name = path.current().name + fmt.Sprintf("[%d]", i)
		} else {
			child.name = path.current().name + fmt.Sprintf("[%q]", child.name)
		}
		buildPlan(append(path, child), vars, plan)
	}
}

type plan []layer

func (p *plan) Append(name string, vars vars) {
	*p = append(*p, layer{name: name, vars: vars})
}

type layer struct {
	name     string
	vars     vars
	children []*layer
}

type path []*layer

func (p path) current() *layer {
	return p[len(p)-1]
}

func (p path) bottom() bool {
	return len(p.current().children) == 0
}

type layerFn func(vars vars, subplans ...*layer) *layer

func (pf layerFn) WithName(name string) layerFn {
	return func(vars vars, subplans ...*layer) *layer {
		plan := pf(vars, subplans...)
		plan.name = name
		return plan
	}
}

var Layer layerFn = func(vars vars, layers ...*layer) *layer {
	return &layer{
		vars:     vars,
		children: layers,
	}
}

type vars map[string]any

func merge(a vars, b vars) vars {
	m := make(vars, len(a)+len(b))
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		switch bval := v.(type) {
		case vars:
			if m[k] == nil {
				m[k] = bval
				continue
			}
			aMap, ok := m[k].(vars)
			if !ok {
				panic(fmt.Sprintf("key %q: cannot merge vars into type %T", k, m[k]))
			}
			m[k] = merge(aMap, bval)
		default:
			m[k] = bval
		}
	}
	return m
}
