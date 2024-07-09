package gitconfig

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var (
		configContent []byte
		err           error
	)
	testcases := []struct {
		name       string
		wantErr    error
		configPath string
	}{
		{
			name:       "bad section name",
			configPath: "configsamples/badSectionName.gitconfig",
			wantErr: &ParseError{
				Err:        ErrInvalidSection,
				LineNumber: 1,
				Line:       "[fo!o] # section name can only contain alphanumeric characters and '-'",
			},
		},
		{
			name:       "bad variable name",
			configPath: "configsamples/badVariableName.gitconfig",
			wantErr: &ParseError{
				Err:        ErrInvalidLine,
				LineNumber: 2,
				Line:       "	1bar = baz  # variable name must start with an alphabetic character",
			},
		},
		{
			name:       "bad variable value",
			configPath: "configsamples/badVariableValue.gitconfig",
			wantErr: &ParseError{
				Err:        ErrInvalidVariableValue,
				LineNumber: 2,
				Line:       `	bar = baz  \ # '\' indicates that the value continues on the next line, there MUST NOT be any character after '\'`,
			},
		},
		{
			name:       "good config",
			configPath: "configsamples/good.gitconfig",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			configContent, err = os.ReadFile(tt.configPath)
			if err != nil {
				t.Fatalf("os.ReadFile(%s) error = %v, want = %v", tt.configPath, err, nil)
			}
			_, err = Parse(configContent)
			if tt.wantErr != err {
				if tt.wantErr == nil {
					t.Fatalf("Parse() error = %v, want = %T", err, tt.wantErr)
				}
				parseErr := new(ParseError)
				if !errors.As(err, &parseErr) {
					t.Errorf("Parse() error = %v, want = %T", err, tt.wantErr)
				}

				ttwantErr := new(ParseError)
				errors.As(tt.wantErr, &ttwantErr)

				// check error values
				if parseErr.LineNumber != ttwantErr.LineNumber {
					t.Errorf("expected error to be at line number %d, got number %d", ttwantErr.LineNumber, parseErr.LineNumber)
				}
				if parseErr.Line != ttwantErr.Line {
					t.Errorf("expected error to be at line %s, got %s", ttwantErr, parseErr.Line)
				}
				if !errors.Is(parseErr.Err, ttwantErr.Err) {
					t.Errorf("expected error to be %s, got %s", ttwantErr, parseErr.Err)
				}
			}
		})
	}

}

func TestParsedValue(t *testing.T) {
	configContent, err := os.ReadFile("configsamples/good.gitconfig")
	if err != nil {
		t.Fatalf("os.ReadFile error = %v, want %v", err, nil)
	}

	parsed, err := Parse(configContent)
	if err != nil {
		t.Fatalf("Parse error = %v, want %v", err, nil)
	}

	for section := range parsed.data {
		keyval := parsed.data[section]
		for key := range keyval {
			keyStr := fmt.Sprintf("%s.%s", section.DottedString(), key)
			originalValue := getOriginalValue(keyStr)
			parsedValue, err := parsed.GetAll(keyStr)
			if err != nil {
				t.Errorf("Parsed.GetAll(%s) error = %v, want %v", keyStr, err, nil)
			}
			parsedValueStr := join(parsedValue)
			if originalValue != parsedValueStr {
				t.Errorf("%s = %s, want %s", keyStr, parsedValueStr, originalValue)
			}
		}
	}
}

func getOriginalValue(key string) string {
	cmd := exec.Command("git", "config", "--file", "configsamples/good.gitconfig", "--get-all", key)
	gitOutput, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return string(bytes.TrimSpace(gitOutput))
}

func join(vals []Value) string {
	var sb strings.Builder

	for i := range vals {
		sb.WriteString(vals[i].String())
		if len(vals) > 1 && i < len(vals)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
