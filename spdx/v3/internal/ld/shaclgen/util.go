package shaclgen

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/deiu/rdf2go"
)

type usedFunc func(...*rdf2go.Triple) []*rdf2go.Triple

func oneOptional(used usedFunc, all []*rdf2go.Triple) *rdf2go.Triple {
	if len(all) > 1 {
		panic(fmt.Errorf("too many results for: %#v", all))
	}
	if len(all) > 0 {
		used(all[0])
		return all[0]
	}
	return nil
}

func oneRequired(used usedFunc, all []*rdf2go.Triple) *rdf2go.Triple {
	if len(all) != 1 {
		panic(fmt.Errorf("required exactly 1 result for: %#v", all))
	}
	used(all[0])
	return all[0]
}

func cleanIRI(iri string) string {
	return strings.Trim(iri, "<>")
}

func cleanText(iri string) string {
	return strings.Trim(iri, "\"")
}

var logEnabled = false

func log(msg ...any) {
	if !logEnabled {
		return
	}
	for _, m := range msg {
		switch m := m.(type) {
		case *rdf2go.Triple:
			_, _ = fmt.Fprint(os.Stderr, nodeDisplay(m))
		default:
			_, _ = fmt.Fprint(os.Stderr, m)
		}
		_, _ = fmt.Fprint(os.Stderr, " ")
	}
	_, _ = fmt.Fprintln(os.Stderr)
}

func nodeDisplay(triple *rdf2go.Triple) string {
	return join("Subject: ", triple.Subject.String(), ", Predicate: ", triple.Predicate.String(), ", Object: ", triple.Object.String())
}

func join(parts ...string) string {
	return strings.Join(parts, "")
}

func bySubject(t *rdf2go.Triple) string {
	return t.Subject.String()
}

func byObject(t *rdf2go.Triple) string {
	return t.Object.String()
}

func sorted(values []*rdf2go.Triple, by func(triple *rdf2go.Triple) string) []*rdf2go.Triple {
	slices.SortFunc(values, func(a, b *rdf2go.Triple) int {
		if a == nil && b == nil {
			return 0
		}
		if a == nil {
			return 1
		}
		if b == nil {
			return -1
		}
		return strings.Compare(by(a), by(b))
	})
	return values
}

func fetch(definitions string) []byte {
	spdxTTLRes := get(http.Get(definitions))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = os.Stderr.Write([]byte(fmt.Sprint(err)))
		}
	}(spdxTTLRes.Body)
	return get(io.ReadAll(spdxTTLRes.Body))
}

func get[T any](t T, err error) T {
	must(err)
	return t
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func in[K comparable, V any](values map[K]V, value K) bool {
	_, ok := values[value]
	return ok
}

func keys[K cmp.Ordered, V any](values map[K]V) []K {
	out := make([]K, 0, len(values))
	for v := range maps.Keys(values) {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}
