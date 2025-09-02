package aiclient

import (
	"fmt"
	"strings"
)

type prompt struct {
	role         string
	context      string
	examples     string
	task         string
	instructions string
}

func PromptBuilder() *prompt {
	p := &prompt{}

	return p
}

func (p *prompt) WithRole(role string) *prompt {
	p.role = role
	return p
}

func (p *prompt) WithContext(context string) *prompt {
	p.context = context
	return p
}

func (p *prompt) WithExamples(examples string) *prompt {
	p.examples = examples
	return p
}

func (p *prompt) WithTask(task string) *prompt {
	p.task = task
	return p
}

func (p *prompt) WithInstructions(instructions string) *prompt {
	p.instructions = instructions
	return p
}

func (p *prompt) Get() string {
	var builder strings.Builder
	if p.role != "" {
		builder.WriteString(fmt.Sprintf("[ROLE]\n%s\n\n", p.role))
	}
	if p.context != "" {
		builder.WriteString(fmt.Sprintf("[CONTEXT]\n%s\n\n", p.context))
	}
	if p.examples != "" {
		builder.WriteString(fmt.Sprintf("[EXAMPLES]\n%s\n\n", p.examples))
	}
	if p.task != "" {
		builder.WriteString(fmt.Sprintf("[TASK]\n%s\n\n", p.task))
	}
	if p.instructions != "" {
		builder.WriteString(fmt.Sprintf("[INSTRUCTIONS]\n%s\n\n", p.instructions))
	}
	return builder.String()
}
