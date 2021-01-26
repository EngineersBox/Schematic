package parser

import (
	"fmt"
	"strings"
)

const providerReferenceDelimiter = "::"

type ProviderReference struct {
	Provider string
	Kind     string
}

func newProviderReference(provRefStr string) (*ProviderReference, error) {
	if !strings.Contains(provRefStr, providerReferenceDelimiter) {
		return nil, fmt.Errorf("invalid provider reference format, must be of the form <PROVIDER>::<KIND>")
	}
	provRef := strings.Split(provRefStr, providerReferenceDelimiter)
	if len(provRef) != 2 {
		return nil, fmt.Errorf("invalid provider reference format, must be of the form <PROVIDER>::<KIND>")
	}
	return &ProviderReference{
		Provider: provRef[0],
		Kind:     provRef[1],
	}, nil
}

func (p *ProviderReference) AsString() string {
	return fmt.Sprintf("%s::%s", p.Provider, p.Kind)
}
