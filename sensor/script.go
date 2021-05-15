package sensor

import (
	"bufio"
	"bytes"
	"context"
	"hacompanion/entity"
	"log"
	"os/exec"
	"strings"
)

type Script struct {
	cfg entity.ScriptConfig
}

func NewScriptRunner(cfg entity.ScriptConfig) *Script {
	return &Script{
		cfg: cfg,
	}
}

func (s Script) Run(ctx context.Context) (*entity.Payload, error) {
	var err error
	var out bytes.Buffer

	// Call the custom script.
	cmd := exec.CommandContext(ctx, s.cfg.Path)
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	n := 0
	p := entity.NewPayload()
	sc := bufio.NewScanner(strings.NewReader(out.String()))
	for sc.Scan() {
		n++
		line := strings.TrimSpace(sc.Text())
		// First line has to contain state.
		if n == 1 {
			p.State = line
			continue
		}
		// Other lines are attributes.
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			log.Printf("ignoring custom script line with less than two parts: %s\n", line)
			continue
		}
		p.Attributes[strings.TrimSpace(parts[0])] = strings.TrimSpace(strings.Join(parts[1:], ":"))
	}
	return p, nil
}
