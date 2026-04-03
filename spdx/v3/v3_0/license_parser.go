package v3_0

import (
	"fmt"
	"strings"
)

// ParseLicenseExpression parses an SPDX license expression string into the
// corresponding model types. It handles AND, OR, WITH operators, the + (or-later)
// suffix, LicenseRef and DocumentRef references, and parenthesized sub-expressions.
//
// Operator precedence (lowest to highest): OR, AND, WITH, + (or-later)
//
// Examples:
//   - "MIT" → *ListedLicense
//   - "MIT OR Apache-2.0" → *DisjunctiveLicenseSet
//   - "MIT AND Apache-2.0" → *ConjunctiveLicenseSet
//   - "GPL-2.0-only WITH Classpath-exception-2.0" → *WithAdditionOperator
//   - "GPL-2.0-only+" → *OrLaterOperator wrapping *ListedLicense
//   - "LicenseRef-custom" → *CustomLicense
//   - "DocumentRef-ext:LicenseRef-custom" → *CustomLicense
func ParseLicenseExpression(expression string) (AnyLicenseInfo, error) {
	p := &licenseParser{input: strings.TrimSpace(expression)}
	if len(p.input) == 0 {
		return nil, fmt.Errorf("empty license expression")
	}
	result, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	p.skipWhitespace()
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected token at position %d: %q", p.pos, p.remaining())
	}
	return result, nil
}

type licenseParser struct {
	input string
	pos   int
}

func (p *licenseParser) remaining() string {
	if p.pos >= len(p.input) {
		return ""
	}
	return p.input[p.pos:]
}

// parseOr handles: and-expr ("OR" and-expr)*
// Flattens consecutive OR operands into a single DisjunctiveLicenseSet.
func (p *licenseParser) parseOr() (AnyLicenseInfo, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	members := LicenseInfoList{left}
	for p.peekKeyword("OR") {
		p.consumeKeyword("OR")
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		members = append(members, right)
	}

	if len(members) == 1 {
		return left, nil
	}
	return &DisjunctiveLicenseSet{Members: members}, nil
}

// parseAnd handles: with-expr ("AND" with-expr)*
// Flattens consecutive AND operands into a single ConjunctiveLicenseSet.
func (p *licenseParser) parseAnd() (AnyLicenseInfo, error) {
	left, err := p.parseWith()
	if err != nil {
		return nil, err
	}

	members := LicenseInfoList{left}
	for p.peekKeyword("AND") {
		p.consumeKeyword("AND")
		right, err := p.parseWith()
		if err != nil {
			return nil, err
		}
		members = append(members, right)
	}

	if len(members) == 1 {
		return left, nil
	}
	return &ConjunctiveLicenseSet{Members: members}, nil
}

// parseWith handles: simple-expr ("WITH" exception-id)?
func (p *licenseParser) parseWith() (AnyLicenseInfo, error) {
	license, err := p.parseSimple()
	if err != nil {
		return nil, err
	}

	if !p.peekKeyword("WITH") {
		return license, nil
	}
	p.consumeKeyword("WITH")

	extendable, ok := license.(AnyExtendableLicense)
	if !ok {
		return nil, fmt.Errorf("WITH operator requires a simple license or or-later expression, got %T", license)
	}

	exceptionID := p.scanIdent()
	if exceptionID == "" {
		return nil, fmt.Errorf("expected license exception identifier after WITH at position %d", p.pos)
	}

	return &WithAdditionOperator{
		SubjectExtendableLicense: extendable,
		SubjectAddition:          makeAddition(exceptionID),
	}, nil
}

// parseSimple handles: "(" expression ")" | license-id "+"? | license-ref
func (p *licenseParser) parseSimple() (AnyLicenseInfo, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	if p.input[p.pos] == '(' {
		p.pos++
		expr, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return nil, fmt.Errorf("expected ')' at position %d", p.pos)
		}
		p.pos++
		return expr, nil
	}

	ident := p.scanIdent()
	if ident == "" {
		return nil, fmt.Errorf("expected license identifier at position %d, got %q", p.pos, p.remaining())
	}

	// Check for or-later (+) suffix
	if p.pos < len(p.input) && p.input[p.pos] == '+' {
		p.pos++
		return &OrLaterOperator{SubjectLicense: makeLicense(ident)}, nil
	}

	return makeLicense(ident), nil
}

func (p *licenseParser) skipWhitespace() {
	for p.pos < len(p.input) && p.input[p.pos] == ' ' {
		p.pos++
	}
}

// peekKeyword checks if the next non-whitespace token is exactly the given keyword,
// ensuring it is not part of a longer identifier.
func (p *licenseParser) peekKeyword(keyword string) bool {
	p.skipWhitespace()
	end := p.pos + len(keyword)
	if end > len(p.input) {
		return false
	}
	if !strings.EqualFold(p.input[p.pos:end], keyword) {
		return false
	}
	if end < len(p.input) && isIdentChar(p.input[end]) {
		return false
	}
	return true
}

func (p *licenseParser) consumeKeyword(keyword string) {
	p.skipWhitespace()
	p.pos += len(keyword)
}

// scanIdent reads an SPDX identifier: [a-zA-Z0-9.-]+
// Also handles the DocumentRef-xxx:LicenseRef-xxx compound form.
func (p *licenseParser) scanIdent() string {
	p.skipWhitespace()
	start := p.pos
	for p.pos < len(p.input) && isIdentChar(p.input[p.pos]) {
		p.pos++
	}
	if p.pos == start {
		return ""
	}
	ident := p.input[start:p.pos]
	// Handle DocumentRef-xxx:LicenseRef-xxx
	if strings.HasPrefix(ident, "DocumentRef-") && p.pos < len(p.input) && p.input[p.pos] == ':' {
		p.pos++ // consume ':'
		for p.pos < len(p.input) && isIdentChar(p.input[p.pos]) {
			p.pos++
		}
		ident = p.input[start:p.pos]
	}
	return ident
}

func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '.'
}

// makeLicense creates the appropriate license type based on the identifier prefix.
// LicenseRef-* and DocumentRef-* identifiers produce CustomLicense, all others produce ListedLicense.
func makeLicense(ident string) AnyLicense {
	if strings.HasPrefix(ident, "LicenseRef-") || strings.HasPrefix(ident, "DocumentRef-") {
		return &CustomLicense{ID: ident}
	}
	return &ListedLicense{Name: ident}
}

// makeAddition creates the appropriate license addition type based on the identifier prefix.
// AdditionRef-*, LicenseRef-*, and DocumentRef-* identifiers produce CustomLicenseAddition,
// all others produce ListedLicenseException.
func makeAddition(ident string) AnyLicenseAddition {
	if strings.HasPrefix(ident, "AdditionRef-") || strings.HasPrefix(ident, "LicenseRef-") || strings.HasPrefix(ident, "DocumentRef-") {
		return &CustomLicenseAddition{ID: ident}
	}
	return &ListedLicenseException{Name: ident}
}
