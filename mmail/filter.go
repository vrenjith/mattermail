package mmail

import (
	"fmt"
	"strings"
)

// Rule for filter
type Rule struct {
	To      string
	From    string
	Subject string
	Channel string
}

// Filter has an array of rules
type Filter []*Rule

// Fix remove spaces and convert to lower case the rules
func (r *Rule) Fix() {
	r.To = strings.TrimSpace(strings.ToLower(r.To))
	r.From = strings.TrimSpace(strings.ToLower(r.From))
	r.Subject = strings.TrimSpace(strings.ToLower(r.Subject))
	r.Channel = strings.TrimSpace(strings.ToLower(r.Channel))
}

// IsValid check if this rule is valid
func (r *Rule) IsValid() error {
	if len(r.From) == 0 && len(r.Subject) == 0 && len(r.To) == 0 {
		return fmt.Errorf("Need to set From or Subject or To")
	}

	if len(r.Channel) == 0 {
		return fmt.Errorf("Need to set a Channel")
	}

	if !strings.HasPrefix(r.Channel, "#") && !strings.HasPrefix(r.Channel, "@") {
		return fmt.Errorf("Need to set a #channel or @user")
	}
	return nil
}

func (r *Rule) meetsFrom(from string) bool {
	from = strings.ToLower(from)
	if len(r.From) == 0 {
		return true
	}
	return strings.Contains(from, r.From)
}

func (r *Rule) meetsSubject(subject string) bool {
	subject = strings.ToLower(subject)
	if len(r.Subject) == 0 {
		return true
	}
	return strings.Contains(subject, r.Subject)
}

func (r *Rule) meetsTo(to string) bool {
	to = strings.ToLower(to)
	if len(r.To) == 0 {
		return true
	}
	return strings.Contains(to, r.To)
}

// MeetsRule check if from and subject meets this rule
func (r *Rule) MeetsRule(from, to, subject string) bool {
	return r.meetsFrom(from) && r.meetsSubject(subject) && r.meetsTo(to)
}

// GetChannel return the first channel with attempt the rules
func (f *Filter) GetChannel(from, to, subject string) string {
	for _, r := range *f {
		if r.MeetsRule(from, to, subject) {
			return r.Channel
		}
	}
	return ""
}

// Valid check if all rules is valid and fix
func (f *Filter) Valid() error {
	for _, r := range *f {
		r.Fix()
		if err := r.IsValid(); err != nil {
			return err
		}
	}
	return nil
}
