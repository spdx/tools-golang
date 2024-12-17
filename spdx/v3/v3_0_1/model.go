//
//
//

package v3_0_1

import (
    "encoding/json"
    "fmt"
    "reflect"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/ncruces/go-strftime"
)

// Validation Error
type ValidationError struct {
    Property string
    Err string
}

func (e *ValidationError) Error() string { return e.Property + ": " + e.Err }

// Conversion Error
type ConversionError struct {
    From string
    To string
}

func (e *ConversionError) Error() string {
    return "Unable to convert from " + e.From + " to " + e.To
}

// Decode Error
type DecodeError struct {
    Path Path
    Err string
}

func (e *DecodeError) Error() string {
    return e.Path.ToString() + ": " + e.Err
}

type EncodeError struct {
    Path Path
    Err string
}

func (e *EncodeError) Error() string {
    return e.Path.ToString() + ": " + e.Err
}

// Path
type Path struct {
    Path []string
}

func (p *Path) PushPath(s string) Path {
    new_p := *p
    new_p.Path = append(new_p.Path, s)
    return new_p
}

func (p *Path) PushIndex(idx int) Path {
    return p.PushPath("[" + strconv.Itoa(idx) + "]")
}

func (p *Path) ToString() string {
    return "." + strings.Join(p.Path, ".")
}

// Error Handler
type ErrorHandler interface {
    HandleError(error, Path)
}

// Reference
type Ref[T SHACLObject] interface {
    GetIRI() string
    GetObj() T
    IsSet() bool
    IsObj() bool
    IsIRI() bool
}

type ref[T SHACLObject] struct {
    obj *T
    iri string
}

func (r ref[T]) GetIRI() string {
    if r.iri != "" {
        return r.iri
    }
    if r.obj != nil {
        o := *r.obj
        if o.ID().IsSet() {
            return o.ID().Get()
        }
    }
    return ""
}

func (r ref[T]) GetObj() T {
    return *r.obj
}

func (r ref[T]) IsSet() bool { return r.IsIRI() || r.IsObj() }
func (r ref[T]) IsObj() bool { return r.obj != nil }
func (r ref[T]) IsIRI() bool { return r.iri != "" }

func MakeObjectRef[T SHACLObject](obj T) Ref[T] {
    return ref[T]{&obj, ""}
}

func MakeIRIRef[T SHACLObject](iri string) Ref[T] {
    return ref[T]{nil, iri}
}

// Convert one reference to another. Note that the output type is first so it
// can be specified, while the input type is generally inferred from the argument
func ConvertRef[TO SHACLObject, FROM SHACLObject](in Ref[FROM]) (Ref[TO], error) {
    if in.IsObj() {
        out_obj, ok := any(in.GetObj()).(TO)
        if !ok {
            return nil, &ConversionError{reflect.TypeOf(ref[FROM]{}).Name(), reflect.TypeOf(ref[TO]{}).Name()}
        }
        return ref[TO]{&out_obj, in.GetIRI()}, nil
    }
    return ref[TO]{nil, in.GetIRI()}, nil
}

type Visit func(Path, any)

// Base SHACL Object
type SHACLObjectBase struct {
    // Object ID
    id Property[string]
    typ SHACLType
    typeIRI string
}

func (self *SHACLObjectBase) ID() PropertyInterface[string] { return &self.id }

func (self *SHACLObjectBase) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true

    switch self.typ.GetNodeKind() {
    case NodeKindBlankNode:
        if self.ID().IsSet() && ! IsBlankNode(self.ID().Get()) {
            handler.HandleError(&ValidationError{
                "id",
                "ID must by be blank node"},
                path.PushPath("id"))
            valid = false
        }
    case NodeKindIRI:
        if ! self.ID().IsSet() || ! IsIRI(self.ID().Get()) {
            handler.HandleError(&ValidationError{
                "id",
                "ID must be an IRI"},
                path.PushPath("id"))
            valid = false
        }
    case NodeKindBlankNodeOrIRI:
        if self.ID().IsSet() && ! IsBlankNode(self.ID().Get()) && ! IsIRI(self.ID().Get()) {
            handler.HandleError(&ValidationError{
                "id",
                "ID must be a blank node or IRI"},
                path.PushPath("id"))
            valid = false
        }
    default:
        panic("Unknown node kind")
    }

    return valid
}

func (self *SHACLObjectBase) Walk(path Path, visit Visit) {
    self.id.Walk(path, visit)
}

func (self *SHACLObjectBase) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if self.typeIRI != "" {
        data["type"] = self.typeIRI
    } else {
        data["type"] = self.typ.GetCompactTypeIRI().GetDefault(self.typ.GetTypeIRI())
    }

    id_prop := self.typ.GetIDAlias().GetDefault("@id")

    if self.id.IsSet() {
        val, err := EncodeIRI(self.id.Get(), path.PushPath(id_prop), map[string]string{}, state)
        if err != nil {
            return err
        }
        data[id_prop] = val
    }

    return nil
}

func (self *SHACLObjectBase) GetType() SHACLType {
    return self.typ
}

func (self *SHACLObjectBase) setTypeIRI(iri string) {
    self.typeIRI = iri
}

func ConstructSHACLObjectBase(o *SHACLObjectBase, typ SHACLType) *SHACLObjectBase {
    o.id = NewProperty[string]("id", []Validator[string]{ IDValidator{}, })
    o.typ = typ
    return o
}

type SHACLObject interface {
    ID() PropertyInterface[string]
    Validate(path Path, handler ErrorHandler) bool
    Walk(path Path, visit Visit)
    EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error
    GetType() SHACLType
    setTypeIRI(iri string)
}

func EncodeSHACLObject(o SHACLObject, path Path, state *EncodeState) (any, error) {
    if state != nil {
        if state.Written[o] {
            if o.ID().IsSet() {
                return o.ID().Get(), nil
            }

            return nil, &EncodeError{
                path,
                "Object referenced multiple times, but does not have an ID assigned",
            }
        }

        state.Written[o] = true
    }

    d := make(map[string]interface{})
    return d, o.EncodeProperties(d, path, state)
}

// Extensible Object

type SHACLExtensibleBase struct {
    properties map[string][]any
}

func (self *SHACLExtensibleBase) GetExtProperty(name string) []any {
    return self.properties[name]
}

func (self *SHACLExtensibleBase) SetExtProperty(name string, value []any) {
    if self.properties == nil {
        self.properties = make(map[string][]any)
    }
    self.properties[name] = value
}

func (self *SHACLExtensibleBase) DeleteExtProperty(name string) {
    delete(self.properties, name)
}

func (self *SHACLExtensibleBase) EncodeExtProperties(data map[string]any, path Path) error {
    for k, values := range self.properties {
        if len(values) == 0 {
            continue
        }

        lst := []any{}
        for _, v := range values {
            lst = append(lst, v)
        }
        data[k] = lst
    }
    return nil
}

type SHACLExtensibleObject interface {
    GetExtProperty(string) []any
    SetExtProperty(string, []any)
    DeleteExtProperty(string)
}

// Type Metadata
const NodeKindBlankNode = 0
const NodeKindIRI = 1
const NodeKindBlankNodeOrIRI = 2

type SHACLType interface {
    GetTypeIRI() string
    GetCompactTypeIRI() Optional[string]
    GetNodeKind() int
    GetIDAlias() Optional[string]
    DecodeProperty(SHACLObject, string, interface{}, Path) (bool, error)
    Create() SHACLObject
    IsAbstract() bool
    IsExtensible() bool
    IsSubClassOf(SHACLType) bool
}

type SHACLTypeBase struct {
    typeIRI string
    compactTypeIRI Optional[string]
    idAlias Optional[string]
    isExtensible Optional[bool]
    isAbstract bool
    parentIRIs []string
    nodeKind Optional[int]
}

func (self SHACLTypeBase) GetTypeIRI() string {
    return self.typeIRI
}

func (self SHACLTypeBase) GetCompactTypeIRI() Optional[string] {
    return self.compactTypeIRI
}

func (self SHACLTypeBase) GetNodeKind() int {
    if self.nodeKind.IsSet() {
        return self.nodeKind.Get()
    }

    for _, parent_id := range(self.parentIRIs) {
        p := objectTypes[parent_id]
        return p.GetNodeKind()
    }

    return NodeKindBlankNodeOrIRI
}

func (self SHACLTypeBase) GetIDAlias() Optional[string] {
    if self.idAlias.IsSet() {
        return self.idAlias
    }

    for _, parent_id := range(self.parentIRIs) {
        p := objectTypes[parent_id]
        a := p.GetIDAlias()
        if a.IsSet() {
            return a
        }
    }

    return self.idAlias
}

func (self SHACLTypeBase) IsAbstract() bool {
    return self.isAbstract
}

func (self SHACLTypeBase) IsExtensible() bool {
    if self.isExtensible.IsSet() {
        return self.isExtensible.Get()
    }

    for _, parent_id := range(self.parentIRIs) {
        p := objectTypes[parent_id]
        if p.IsExtensible() {
            return true
        }
    }

    return false
}

func (self SHACLTypeBase) IsSubClassOf(other SHACLType) bool {
    if other.GetTypeIRI() == self.typeIRI {
        return true
    }

    for _, parent_id := range(self.parentIRIs) {
        p := objectTypes[parent_id]
        if p.IsSubClassOf(other) {
            return true
        }
    }

    return false
}

type EncodeState struct {
    Written map[SHACLObject]bool
}

func (self SHACLTypeBase) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    id_alias := self.GetIDAlias()
    if id_alias.IsSet() {
        switch name {
        case id_alias.Get():
            val, err := DecodeString(value, path.PushPath(name), map[string]string{})
            if err != nil {
                return false, err
            }
            err = o.ID().Set(val)
            if err != nil {
                return false, err
            }
            return true, nil
        case "@id":
            return true, &DecodeError{
                path.PushPath(name),
                "'@id' is not allowed for " + self.GetTypeIRI() + " which has an ID alias",
            }
        }
    } else if name == "@id" {
        val, err := DecodeString(value, path.PushPath(name), map[string]string{})
        if err != nil {
            return false, err
        }
        err = o.ID().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    }

    for _, parent_id := range(self.parentIRIs) {
        p := objectTypes[parent_id]
        found, err := p.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    if self.isExtensible.GetDefault(false) {
        obj := o.(SHACLExtensibleObject)
        v, err := DecodeAny(value, path, map[string]string{})
        if err != nil {
            return false, err
        }

        lst, is_list := v.([]interface{})
        if is_list {
            obj.SetExtProperty(name, lst)
        } else {
            obj.SetExtProperty(name, []interface{}{v})
        }
        return true, nil
    }
    return false, nil
}


var objectTypes map[string] SHACLType

func RegisterType(typ SHACLType) {
    objectTypes[typ.GetTypeIRI()] = typ
    compact := typ.GetCompactTypeIRI()
    if compact.IsSet() {
        objectTypes[compact.Get()] = typ
    }
}

// SHACLObjectSet
type SHACLObjectSet interface {
    AddObject(r SHACLObject)
    Decode(decoder *json.Decoder) error
    Encode(encoder *json.Encoder) error
    Walk(visit Visit)
    Validate(handler ErrorHandler) bool
}

type SHACLObjectSetObject struct {
    objects []SHACLObject
}

func (self *SHACLObjectSetObject) AddObject(r SHACLObject) {
    self.objects = append(self.objects, r)
}

func (self *SHACLObjectSetObject) Decode(decoder *json.Decoder) error {
    path := Path{}

    var data map[string]interface{}
    if err := decoder.Decode(&data); err != nil {
        return err
    }

    {
        v, ok := data["@context"]
        if ! ok {
            return &DecodeError{path, "@context missing"}
        }

        sub_path := path.PushPath("@context")
        value, ok := v.(string)
        if ! ok {
            return &DecodeError{sub_path, "@context must be a string, or list of string"}
        }
        if value != "https://spdx.org/rdf/3.0.1/spdx-context.jsonld" {
            return &DecodeError{sub_path, "Wrong context URL '" + value + "'"}
        }
    }

    delete(data, "@context")

    decodeProxy := func (data any, path Path, context map[string]string) (SHACLObject, error) {
        return DecodeSHACLObject[SHACLObject](data, path, context, nil)
    }

    _, has_graph := data["@graph"]
    if has_graph {
        for k, v := range data {
            switch k {
            case "@graph": {
                objs, err := DecodeList[SHACLObject](
                    v,
                    path.PushPath("@graph"),
                    map[string]string{},
                    decodeProxy,
                )

                if err != nil {
                    return err
                }

                for _, obj := range objs {
                    self.AddObject(obj)
                }
            }

            default:
                return &DecodeError{path, "Unknown property '" + k + "'"}
            }
        }
    } else {
        obj, err := decodeProxy(data, path, map[string]string{})
        if err != nil {
            return err
        }

        self.AddObject(obj)
    }

    return nil
}

func (self *SHACLObjectSetObject) Encode(encoder *json.Encoder) error {
    data := make(map[string]interface{})
    data["@context"] = "https://spdx.org/rdf/3.0.1/spdx-context.jsonld"
    path := Path{}
    state := EncodeState{
        Written: make(map[SHACLObject]bool),
    }

    ref_counts := make(map[SHACLObject]int)

    visit := func (path Path, v any) {
        r, ok := v.(Ref[SHACLObject])
        if ! ok {
            return
        }

        if ! r.IsObj() {
            return
        }

        o := r.GetObj()

        // Remove blank nodes for reassignment
        if o.ID().IsSet() && IsBlankNode(o.ID().Get()) {
            o.ID().Delete()
        }

        ref_counts[o] = ref_counts[o] + 1
    }

    self.Walk(visit)

    blank_count := 0
    for o, count := range ref_counts {
        if count <= 1 {
            continue
        }

        if o.ID().IsSet() {
            continue
        }

        o.ID().Set(fmt.Sprintf("_:%d", blank_count))
        blank_count += 1
    }

    if len(self.objects) == 1 {
        err := self.objects[0].EncodeProperties(data, path, &state)
        if err != nil {
            return err
        }
    } else if len(self.objects) > 1 {
        // All objects directly added to the object set should be written as
        // top level objects, so mark then as written until they are ready to
        // be serialized, which will force them to be referenced by IRI until
        // we are ready
        for _, o := range self.objects {
            state.Written[o] = true
        }

        graph_path := path.PushPath("@graph")
        lst := []interface{}{}
        for idx, o := range self.objects {
            // Remove this object from the written set now so it gets serialized
            delete(state.Written, o)

            d, err := EncodeSHACLObject(o, graph_path.PushIndex(idx), &state)
            if err != nil {
                return err
            }
            lst = append(lst, d)
        }

        data["@graph"] = lst
    }

    return encoder.Encode(data)
}

func (self *SHACLObjectSetObject) Walk(visit Visit) {
    path := Path{}
    visited := map[SHACLObject]bool{}

    visit_proxy := func (path Path, v any) {
        switch v.(type) {
        case Ref[SHACLObject]:
            r := v.(Ref[SHACLObject])
            if ! r.IsObj() {
                visit(path, v)
                return
            }

            o := r.GetObj()
            _, ok := visited[o]
            if ok {
                return
            }
            visited[o] = true
            visit(path, v)
            o.Walk(path, visit)
            return

        default:
            visit(path, v)
            return
        }
    }

    for idx, o := range(self.objects) {
        sub_path := path.PushIndex(idx)
        visit_proxy(sub_path, MakeObjectRef(o))
    }
}

func (self *SHACLObjectSetObject) Validate(handler ErrorHandler) bool {
    valid := true

    visit_proxy := func (path Path, v any) {
        r, ok := v.(Ref[SHACLObject])
        if ! ok {
            return
        }

        if ! r.IsObj() {
            return
        }

        if ! r.GetObj().Validate(path, handler) {
            valid = false
        }
    }

    self.Walk(visit_proxy)

    return valid
}

func NewSHACLObjectSet() SHACLObjectSet {
    os := SHACLObjectSetObject{}
    return &os
}

func DecodeAny(data any, path Path, context map[string]string) (any, error) {
    switch data.(type) {
    case map[string]interface{}:
        return DecodeRef[SHACLObject](data, path, context, nil)
    case string:
        return DecodeString(data, path, context)
    case int:
        return DecodeInteger(data, path, context)
    case float64:
        return DecodeFloat(data, path, context)
    case bool:
        return DecodeBoolean(data, path, context)
    case []interface{}:
        return DecodeList[any](data, path, context, DecodeAny)
    default:
        return nil, &DecodeError{path, "Unknown type "+ reflect.TypeOf(data).Name()}
    }
}

func DecodeSHACLObject[T SHACLObject](data any, path Path, context map[string]string, targetType SHACLType) (T, error) {
    dict, ok := data.(map[string]interface{})
    if ! ok {
        return *new(T), &DecodeError{path, "Expected dictionary or string. Got " + reflect.TypeOf(data).Name()}
    }

    var v interface{}
    v, ok = dict["@type"]
    if ! ok {
        v, ok = dict["type"]
        if ! ok {
            return *new(T), &DecodeError{path, "type missing"}
        }
    }

    var type_iri string
    var create_type SHACLType

    type_iri, ok = v.(string)
    if ! ok {
        return *new(T), &DecodeError{path, "Wrong type for @type. Got " + reflect.TypeOf(v).Name()}
    }

    iri_typ, ok := objectTypes[type_iri]
    if ok {
        if targetType != nil && !iri_typ.IsSubClassOf(targetType) {
            return *new(T), &DecodeError{path, "Type " + type_iri + " is not valid where " +
                    targetType.GetTypeIRI() + " is expected"}
        }

        if iri_typ.IsAbstract() {
            return *new(T), &DecodeError{path, "Unable to create abstract type '" + type_iri + "'"}
        }

        create_type = iri_typ
    } else if targetType != nil && targetType.IsExtensible() {
        // An extensible type is expected, so make one of the correct type
        //
        // Note: An abstract extensible class is actually allowed to be created
        // here
        create_type = targetType
    } else {
        if IsIRI(type_iri)  {
            // It's not clear exactly which type should be created. Search through
            // all types and collect a list of possible Extensible types that are
            // valid in this location.
            possible := []SHACLType{}
            for _, v := range objectTypes {
                if ! v.IsExtensible() {
                    continue
                }

                if v.IsAbstract() {
                    continue
                }

                // If a type was specified, only subclasses of that type are
                // allowed
                if targetType != nil && ! v.IsSubClassOf(targetType) {
                    continue
                }

                possible = append(possible, v)
            }

            // Sort for determinism
            sort.Slice(possible, func(i, j int) bool {
                return possible[i].GetTypeIRI() < possible[j].GetTypeIRI()
            })

            for _, t := range(possible) {
                // Ignore errors
                o, err := DecodeSHACLObject[T](data, path, context, t)
                if err == nil {
                    o.setTypeIRI(type_iri)
                    return o, nil
                }
            }
        }
        return *new(T), &DecodeError{path, "Unable to create object of type '" + type_iri + "' (no matching extensible object)"}
    }

    obj, ok := create_type.Create().(T)
    if ! ok {
        return *new(T), &DecodeError{path, "Unable to create object of type '" + type_iri + "'"}
    }
    obj.setTypeIRI(type_iri)

    for k, v := range dict {
        if k == "@type" {
            continue
        }
        if k == "type" {
            continue
        }

        sub_path := path.PushPath(k)
        found, err := create_type.DecodeProperty(obj, k, v, sub_path)
        if err != nil {
            return *new(T), err
        }
        if ! found {
            return *new(T), &DecodeError{path, "Unknown property '" + k + "'"}
        }
    }

    return obj, nil
}

func DecodeRef[T SHACLObject](data any, path Path, context map[string]string, typ SHACLType) (Ref[T], error) {
    switch data.(type) {
    case string:
        s, err := DecodeIRI(data, path, context)
        if err != nil {
            return nil, err
        }
        return MakeIRIRef[T](s), nil
    }

    obj, err := DecodeSHACLObject[T](data, path, context, typ)
    if err != nil {
        return nil, err
    }

    return MakeObjectRef[T](obj), nil
}

func EncodeRef[T SHACLObject](value Ref[T], path Path, context map[string]string, state *EncodeState) (any, error) {
    if value.IsIRI() {
        v := value.GetIRI()
        compact, ok := context[v]
        if ok {
            return compact, nil
        }
        return v, nil
    }
    return EncodeSHACLObject(value.GetObj(), path, state)
}

func DecodeString(data any, path Path, context map[string]string) (string, error) {
    v, ok := data.(string)
    if ! ok {
        return v, &DecodeError{path, "String expected. Got " + reflect.TypeOf(data).Name()}
    }
    return v, nil
}

func EncodeString(value string, path Path, context map[string]string, state *EncodeState) (any, error) {
    return value, nil
}

func DecodeIRI(data any, path Path, context map[string]string) (string, error) {
    s, err := DecodeString(data, path, context)
    if err != nil {
        return s, err
    }

    for k, v := range context {
        if s == v {
            s = k
            break
        }
    }

    if ! IsBlankNode(s) && ! IsIRI(s) {
        return s, &DecodeError{path, "Must be blank node or IRI. Got '" + s + "'"}
    }

    return s, nil
}

func EncodeIRI(value string, path Path, context map[string]string, state *EncodeState) (any, error) {
    compact, ok := context[value]
    if ok {
        return compact, nil
    }
    return value, nil
}

func DecodeBoolean(data any, path Path, context map[string]string) (bool, error) {
    v, ok := data.(bool)
    if ! ok {
        return v, &DecodeError{path, "Boolean expected. Got " + reflect.TypeOf(data).Name()}
    }
    return v, nil
}

func EncodeBoolean(value bool, path Path, context map[string]string, state *EncodeState) (any, error) {
    return value, nil
}

func DecodeInteger(data any, path Path, context map[string]string) (int, error) {
    switch data.(type) {
    case int:
        return data.(int), nil
    case float64:
        v := data.(float64)
        if v == float64(int64(v)) {
            return int(v), nil
        }
        return 0, &DecodeError{path, "Value must be an integer. Got " + fmt.Sprintf("%f", v)}
    default:
        return 0, &DecodeError{path, "Integer expected. Got " + reflect.TypeOf(data).Name()}
    }
}

func EncodeInteger(value int, path Path, context map[string]string, state *EncodeState) (any, error) {
    return value, nil
}

func DecodeFloat(data any, path Path, context map[string]string) (float64, error) {
    switch data.(type) {
    case float64:
        return data.(float64), nil
    case string:
        v, err := strconv.ParseFloat(data.(string), 64)
        if err != nil {
            return 0, err
        }
        return v, nil
    default:
        return 0, &DecodeError{path, "Float expected. Got " + reflect.TypeOf(data).Name()}
    }
}

func EncodeFloat(value float64, path Path, context map[string]string, state *EncodeState) (any, error) {
    return strconv.FormatFloat(value, 'f', -1, 64), nil
}

const UtcFormatStr = "%Y-%m-%dT%H:%M:%SZ"
const TzFormatStr = "%Y-%m-%dT%H:%M:%S%:z"

func decodeDateTimeString(data any, path Path, re *regexp.Regexp) (time.Time, error) {
    v, ok := data.(string)
    if ! ok {
        return time.Time{}, &DecodeError{path, "String expected. Got " + reflect.TypeOf(data).Name()}
    }

    match := re.FindStringSubmatch(v)

    if match == nil {
        return time.Time{}, &DecodeError{path, "Invalid date time string '" + v + "'"}
    }

    var format string
    s := match[1]
    tzstr := match[2]

    switch tzstr {
    case "Z":
        s += "+00:00"
        format = "%Y-%m-%dT%H:%M:%S%:z"
    case "":
        format = "%Y-%m-%dT%H:%M:%S"
    default:
        s += tzstr
        format = "%Y-%m-%dT%H:%M:%S%:z"
    }

    t, err := strftime.Parse(format, v)
    if err != nil {
        return time.Time{}, &DecodeError{path, "Invalid date time string '" + v + "': " + err.Error()}
    }
    return t, nil
}

var dateTimeRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(Z|[+-]\d{2}:\d{2})?$`)
func DecodeDateTime(data any, path Path, context map[string]string) (time.Time, error) {
    return decodeDateTimeString(data, path, dateTimeRegex)
}

var dateTimeStampRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(Z|[+-]\d{2}:\d{2})$`)
func DecodeDateTimeStamp(data any, path Path, context map[string]string) (time.Time, error) {
    return decodeDateTimeString(data, path, dateTimeStampRegex)
}

func EncodeDateTime(value time.Time, path Path, context map[string]string, state *EncodeState) (any, error) {
    if value.Location() == time.UTC {
        return strftime.Format(UtcFormatStr, value), nil
    }
    return strftime.Format(TzFormatStr, value), nil
}

func DecodeList[T any](data any, path Path, context map[string]string, f func (any, Path, map[string]string) (T, error)) ([]T, error) {
    lst, ok := data.([]interface{})
    if ! ok {
        return nil, &DecodeError{path, "Must be a list"}
    }

    var result []T
    for idx, v := range lst {
        sub_path := path.PushIndex(idx)
        item, err := f(v, sub_path, context)
        if err != nil {
            return nil, err
        }
        result = append(result, item)
    }

    return result, nil
}

func EncodeList[T any](value []T, path Path, context map[string]string, state *EncodeState, f func (T, Path, map[string]string, *EncodeState) (any, error)) (any, error) {
    lst := []any{}
    for idx, v := range value {
        val, err := f(v, path.PushIndex(idx), context, state)
        if err != nil {
            return lst, err
        }

        lst = append(lst, val)
    }
    return lst, nil
}

// IRI Validation
func IsIRI(iri string) bool {
    if strings.HasPrefix(iri, "_:") {
        return false
    }
    if strings.Contains(iri, ":") {
        return true
    }
    return false
}

func IsBlankNode(iri string) bool {
    return strings.HasPrefix(iri, "_:")
}

// Optional
type Optional[T any] struct {
    value *T
}

func (self Optional[T]) Get() T {
    return *self.value
}

func (self Optional[T]) GetDefault(val T) T {
    if ! self.IsSet() {
        return val
    }
    return *self.value
}

func (self Optional[T]) IsSet() bool {
    return self.value != nil
}

func NewOptional[T any](value T) Optional[T] {
    return Optional[T]{&value}
}

func NewEmptyOptional[T any]() Optional[T] {
    return Optional[T]{nil}
}

// Validator
type Validator[T any] interface {
    Check(T, string) error
}

func ValueToString(val any) string {
    switch val.(type) {
    case string:
        return val.(string)
    case int:
        return strconv.Itoa(val.(int))
    case time.Time:
        t := val.(time.Time)
        if t.Location() == time.UTC {
            return strftime.Format(UtcFormatStr, t)
        }
        return strftime.Format(TzFormatStr, t)
    }
    panic("Unsupported Type " + reflect.TypeOf(val).Name())
}


// ID Validator
type IDValidator struct {}

func (self IDValidator) Check(val string, name string) error {
    if ! IsIRI(val) && ! IsBlankNode(val) {
        return &ValidationError{name, "Must be an IRI or a Blank Node"}
    }
    return nil
}


// Regex Validator
type RegexValidator[T int | time.Time | string] struct {
    Regex string
}

func (self RegexValidator[T]) Check(val T, name string) error {
    s := ValueToString(val)

    m, err := regexp.MatchString(self.Regex, s)
    if err != nil {
        return err
    }
    if ! m {
        return &ValidationError{name, "Value '" + s + "' does not match pattern"}
    }
    return nil
}

// Integer Min Validator
type IntegerMinValidator struct {
    Min int
}

func (self IntegerMinValidator) Check(val int, name string) error {
    if val < self.Min {
        return &ValidationError{name, "Value " + strconv.Itoa(val) + " is less than minimum " + strconv.Itoa(self.Min)}
    }
    return nil
}

// Integer Max Validator
type IntegerMaxValidator struct {
    Max int
}

func (self IntegerMaxValidator) Check(val int, name string) error {
    if val > self.Max {
        return &ValidationError{name, "Value " + strconv.Itoa(val) + " is greater than maximum" + strconv.Itoa(self.Max)}
    }
    return nil
}

// Enum Validator
type EnumValidator struct {
    Values []string
}

func (self EnumValidator) Check(val string, name string) error {
    for _, v := range self.Values {
        if val == v {
            return nil
        }
    }
    return &ValidationError{name, "Value '" + val + "' is not a valid enumerated value" }
}

// Property
type PropertyInterface[T any] interface {
    Get() T
    Set(val T) error
    Delete()
    IsSet() bool
    Walk(path Path, visit Visit)
}

type Property[T any] struct {
    value Optional[T]
    name string
    validators []Validator[T]
}

func NewProperty[T any](name string, validators []Validator[T]) Property[T] {
    return Property[T]{
        value: NewEmptyOptional[T](),
        name: name,
        validators: validators,
    }
}

func (self *Property[T]) Get() T {
    return self.value.Get()
}

func (self *Property[T]) Set(val T) error {
    for _, validator := range self.validators {
        err := validator.Check(val, self.name)
        if err != nil {
            return err
        }
    }

    self.value = NewOptional(val)
    return nil
}

func (self *Property[T]) Delete() {
    self.value = NewEmptyOptional[T]()
}

func (self *Property[T]) IsSet() bool {
    return self.value.IsSet()
}

func (self *Property[T]) Check(path Path, handler ErrorHandler) bool {
    if ! self.value.IsSet() {
        return true
    }

    var valid bool
    valid = true

    for _, validator := range self.validators {
        err := validator.Check(self.value.Get(), self.name)
        if err != nil {
            if handler != nil {
                handler.HandleError(err, path)
            }
            valid = false
        }
    }
    return valid
}

func (self *Property[T]) Walk(path Path, visit Visit) {
    if ! self.value.IsSet() {
        return
    }

    visit(path.PushPath(self.name), self.value.Get())
}

// Ref Property
type RefPropertyInterface[T SHACLObject] interface {
    PropertyInterface[Ref[T]]

    GetIRI() string
    GetObj() T
    IsObj() bool
    IsIRI() bool
}

type RefProperty[T SHACLObject] struct {
    Property[Ref[T]]
}

func NewRefProperty[T SHACLObject](name string, validators []Validator[Ref[T]]) RefProperty[T] {
    return RefProperty[T]{
        Property: Property[Ref[T]]{
            value: NewEmptyOptional[Ref[T]](),
            name: name,
            validators: validators,
        },
    }
}

func (self *RefProperty[T]) GetIRI() string {
    return self.Get().GetIRI()
}

func (self *RefProperty[T]) GetObj() T {
    return self.Get().GetObj()
}

func (self *RefProperty[T]) IsSet() bool {
    return self.Property.IsSet() && self.Get().IsSet()
}

func (self *RefProperty[T]) IsObj() bool {
    return self.Property.IsSet() && self.Get().IsObj()
}

func (self *RefProperty[T]) IsIRI() bool {
    return self.Property.IsSet() && self.Get().IsIRI()
}

func (self *RefProperty[T]) Walk(path Path, visit Visit) {
    if ! self.IsSet() {
        return
    }

    r, err := ConvertRef[SHACLObject](self.value.Get())
    if err != nil {
        return
    }

    visit(path.PushPath(self.name), r)
}

// List Property
type ListPropertyInterface[T any] interface {
    Get() []T
    Set(val []T) error
    Append(val T) error
    Delete()
    Walk(path Path, visit Visit)
    IsSet() bool
}

type ListProperty[T any] struct {
    value []T
    name string
    validators []Validator[T]
}

func NewListProperty[T any](name string, validators []Validator[T]) ListProperty[T] {
    return ListProperty[T]{
        value: []T{},
        name: name,
        validators: validators,
    }
}

func (self *ListProperty[T]) Get() []T {
    return self.value
}

func (self *ListProperty[T]) Set(val []T) error {
    for _, v := range val {
        for _, validator := range self.validators {
            err := validator.Check(v, self.name)
            if err != nil {
                return err
            }
        }
    }

    self.value = val
    return nil
}

func (self *ListProperty[T]) Append(val T) error {
    for _, validator := range self.validators {
        err := validator.Check(val, self.name)
        if err != nil {
            return err
        }
    }

    self.value = append(self.value, val)
    return nil
}

func (self *ListProperty[T]) Delete() {
    self.value = []T{}
}

func (self *ListProperty[T]) IsSet() bool {
    return self.value != nil && len(self.value) > 0
}

func (self *ListProperty[T]) Check(path Path, handler ErrorHandler) bool {
    var valid bool
    valid = true

    for idx, v := range self.value {
        for _, validator := range self.validators {
            err := validator.Check(v, self.name)
            if err != nil {
                if handler != nil {
                    handler.HandleError(err, path.PushIndex(idx))
                }
                valid = false
            }
        }
    }
    return valid
}

func (self *ListProperty[T]) Walk(path Path, visit Visit) {
    sub_path := path.PushPath(self.name)

    for idx, v := range self.value {
        visit(sub_path.PushIndex(idx), v)
    }
}

type RefListProperty[T SHACLObject] struct {
    ListProperty[Ref[T]]
}

func NewRefListProperty[T SHACLObject](name string, validators []Validator[Ref[T]]) RefListProperty[T] {
    return RefListProperty[T]{
        ListProperty: ListProperty[Ref[T]]{
            value: []Ref[T]{},
            name: name,
            validators: validators,
        },
    }
}

func (self *RefListProperty[T]) Walk(path Path, visit Visit) {
    sub_path := path.PushPath(self.name)

    for idx, v := range self.value {
        r, err := ConvertRef[SHACLObject](v)
        if err != nil {
            visit(sub_path.PushIndex(idx), r)
        }
    }
}


// A class for describing the energy consumption incurred by an AI model in

// different stages of its lifecycle.
type AiEnergyConsumptionObject struct {
    SHACLObjectBase

    // Specifies the amount of energy consumed when finetuning the AI model that is
    // being used in the AI system.
    aiFinetuningEnergyConsumption RefListProperty[AiEnergyConsumptionDescription]
    // Specifies the amount of energy consumed during inference time by an AI model
    // that is being used in the AI system.
    aiInferenceEnergyConsumption RefListProperty[AiEnergyConsumptionDescription]
    // Specifies the amount of energy consumed when training the AI model that is
    // being used in the AI system.
    aiTrainingEnergyConsumption RefListProperty[AiEnergyConsumptionDescription]
}


type AiEnergyConsumptionObjectType struct {
    SHACLTypeBase
}
var aiEnergyConsumptionType AiEnergyConsumptionObjectType
var aiEnergyConsumptionAiFinetuningEnergyConsumptionContext = map[string]string{}
var aiEnergyConsumptionAiInferenceEnergyConsumptionContext = map[string]string{}
var aiEnergyConsumptionAiTrainingEnergyConsumptionContext = map[string]string{}

func DecodeAiEnergyConsumption (data any, path Path, context map[string]string) (Ref[AiEnergyConsumption], error) {
    return DecodeRef[AiEnergyConsumption](data, path, context, aiEnergyConsumptionType)
}

func (self AiEnergyConsumptionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AiEnergyConsumption)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/AI/finetuningEnergyConsumption", "ai_finetuningEnergyConsumption":
        val, err := DecodeList[Ref[AiEnergyConsumptionDescription]](value, path, aiEnergyConsumptionAiFinetuningEnergyConsumptionContext, DecodeAiEnergyConsumptionDescription)
        if err != nil {
            return false, err
        }
        err = obj.AiFinetuningEnergyConsumption().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/inferenceEnergyConsumption", "ai_inferenceEnergyConsumption":
        val, err := DecodeList[Ref[AiEnergyConsumptionDescription]](value, path, aiEnergyConsumptionAiInferenceEnergyConsumptionContext, DecodeAiEnergyConsumptionDescription)
        if err != nil {
            return false, err
        }
        err = obj.AiInferenceEnergyConsumption().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/trainingEnergyConsumption", "ai_trainingEnergyConsumption":
        val, err := DecodeList[Ref[AiEnergyConsumptionDescription]](value, path, aiEnergyConsumptionAiTrainingEnergyConsumptionContext, DecodeAiEnergyConsumptionDescription)
        if err != nil {
            return false, err
        }
        err = obj.AiTrainingEnergyConsumption().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AiEnergyConsumptionObjectType) Create() SHACLObject {
    return ConstructAiEnergyConsumptionObject(&AiEnergyConsumptionObject{}, self)
}

func ConstructAiEnergyConsumptionObject(o *AiEnergyConsumptionObject, typ SHACLType) *AiEnergyConsumptionObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[Ref[AiEnergyConsumptionDescription]]{}
        o.aiFinetuningEnergyConsumption = NewRefListProperty[AiEnergyConsumptionDescription]("aiFinetuningEnergyConsumption", validators)
    }
    {
        validators := []Validator[Ref[AiEnergyConsumptionDescription]]{}
        o.aiInferenceEnergyConsumption = NewRefListProperty[AiEnergyConsumptionDescription]("aiInferenceEnergyConsumption", validators)
    }
    {
        validators := []Validator[Ref[AiEnergyConsumptionDescription]]{}
        o.aiTrainingEnergyConsumption = NewRefListProperty[AiEnergyConsumptionDescription]("aiTrainingEnergyConsumption", validators)
    }
    return o
}

type AiEnergyConsumption interface {
    SHACLObject
    AiFinetuningEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]]
    AiInferenceEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]]
    AiTrainingEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]]
}


func MakeAiEnergyConsumption() AiEnergyConsumption {
    return ConstructAiEnergyConsumptionObject(&AiEnergyConsumptionObject{}, aiEnergyConsumptionType)
}

func MakeAiEnergyConsumptionRef() Ref[AiEnergyConsumption] {
    o := MakeAiEnergyConsumption()
    return MakeObjectRef[AiEnergyConsumption](o)
}

func (self *AiEnergyConsumptionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("aiFinetuningEnergyConsumption")
        if ! self.aiFinetuningEnergyConsumption.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiInferenceEnergyConsumption")
        if ! self.aiInferenceEnergyConsumption.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiTrainingEnergyConsumption")
        if ! self.aiTrainingEnergyConsumption.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *AiEnergyConsumptionObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.aiFinetuningEnergyConsumption.Walk(path, visit)
    self.aiInferenceEnergyConsumption.Walk(path, visit)
    self.aiTrainingEnergyConsumption.Walk(path, visit)
}


func (self *AiEnergyConsumptionObject) AiFinetuningEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]] { return &self.aiFinetuningEnergyConsumption }
func (self *AiEnergyConsumptionObject) AiInferenceEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]] { return &self.aiInferenceEnergyConsumption }
func (self *AiEnergyConsumptionObject) AiTrainingEnergyConsumption() ListPropertyInterface[Ref[AiEnergyConsumptionDescription]] { return &self.aiTrainingEnergyConsumption }

func (self *AiEnergyConsumptionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.aiFinetuningEnergyConsumption.IsSet() {
        val, err := EncodeList[Ref[AiEnergyConsumptionDescription]](self.aiFinetuningEnergyConsumption.Get(), path.PushPath("aiFinetuningEnergyConsumption"), aiEnergyConsumptionAiFinetuningEnergyConsumptionContext, state, EncodeRef[AiEnergyConsumptionDescription])
        if err != nil {
            return err
        }
        data["ai_finetuningEnergyConsumption"] = val
    }
    if self.aiInferenceEnergyConsumption.IsSet() {
        val, err := EncodeList[Ref[AiEnergyConsumptionDescription]](self.aiInferenceEnergyConsumption.Get(), path.PushPath("aiInferenceEnergyConsumption"), aiEnergyConsumptionAiInferenceEnergyConsumptionContext, state, EncodeRef[AiEnergyConsumptionDescription])
        if err != nil {
            return err
        }
        data["ai_inferenceEnergyConsumption"] = val
    }
    if self.aiTrainingEnergyConsumption.IsSet() {
        val, err := EncodeList[Ref[AiEnergyConsumptionDescription]](self.aiTrainingEnergyConsumption.Get(), path.PushPath("aiTrainingEnergyConsumption"), aiEnergyConsumptionAiTrainingEnergyConsumptionContext, state, EncodeRef[AiEnergyConsumptionDescription])
        if err != nil {
            return err
        }
        data["ai_trainingEnergyConsumption"] = val
    }
    return nil
}

// The class that helps note down the quantity of energy consumption and the unit

// used for measurement.
type AiEnergyConsumptionDescriptionObject struct {
    SHACLObjectBase

    // Represents the energy quantity.
    aiEnergyQuantity Property[float64]
    // Specifies the unit in which energy is measured.
    aiEnergyUnit Property[string]
}


type AiEnergyConsumptionDescriptionObjectType struct {
    SHACLTypeBase
}
var aiEnergyConsumptionDescriptionType AiEnergyConsumptionDescriptionObjectType
var aiEnergyConsumptionDescriptionAiEnergyQuantityContext = map[string]string{}
var aiEnergyConsumptionDescriptionAiEnergyUnitContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/kilowattHour": "kilowattHour",
    "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/megajoule": "megajoule",
    "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/other": "other",}

func DecodeAiEnergyConsumptionDescription (data any, path Path, context map[string]string) (Ref[AiEnergyConsumptionDescription], error) {
    return DecodeRef[AiEnergyConsumptionDescription](data, path, context, aiEnergyConsumptionDescriptionType)
}

func (self AiEnergyConsumptionDescriptionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AiEnergyConsumptionDescription)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/AI/energyQuantity", "ai_energyQuantity":
        val, err := DecodeFloat(value, path, aiEnergyConsumptionDescriptionAiEnergyQuantityContext)
        if err != nil {
            return false, err
        }
        err = obj.AiEnergyQuantity().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/energyUnit", "ai_energyUnit":
        val, err := DecodeIRI(value, path, aiEnergyConsumptionDescriptionAiEnergyUnitContext)
        if err != nil {
            return false, err
        }
        err = obj.AiEnergyUnit().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AiEnergyConsumptionDescriptionObjectType) Create() SHACLObject {
    return ConstructAiEnergyConsumptionDescriptionObject(&AiEnergyConsumptionDescriptionObject{}, self)
}

func ConstructAiEnergyConsumptionDescriptionObject(o *AiEnergyConsumptionDescriptionObject, typ SHACLType) *AiEnergyConsumptionDescriptionObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[float64]{}
        o.aiEnergyQuantity = NewProperty[float64]("aiEnergyQuantity", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/kilowattHour",
                "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/megajoule",
                "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/other",
        }})
        o.aiEnergyUnit = NewProperty[string]("aiEnergyUnit", validators)
    }
    return o
}

type AiEnergyConsumptionDescription interface {
    SHACLObject
    AiEnergyQuantity() PropertyInterface[float64]
    AiEnergyUnit() PropertyInterface[string]
}


func MakeAiEnergyConsumptionDescription() AiEnergyConsumptionDescription {
    return ConstructAiEnergyConsumptionDescriptionObject(&AiEnergyConsumptionDescriptionObject{}, aiEnergyConsumptionDescriptionType)
}

func MakeAiEnergyConsumptionDescriptionRef() Ref[AiEnergyConsumptionDescription] {
    o := MakeAiEnergyConsumptionDescription()
    return MakeObjectRef[AiEnergyConsumptionDescription](o)
}

func (self *AiEnergyConsumptionDescriptionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("aiEnergyQuantity")
        if ! self.aiEnergyQuantity.Check(prop_path, handler) {
            valid = false
        }
        if ! self.aiEnergyQuantity.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"aiEnergyQuantity", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiEnergyUnit")
        if ! self.aiEnergyUnit.Check(prop_path, handler) {
            valid = false
        }
        if ! self.aiEnergyUnit.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"aiEnergyUnit", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *AiEnergyConsumptionDescriptionObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.aiEnergyQuantity.Walk(path, visit)
    self.aiEnergyUnit.Walk(path, visit)
}


func (self *AiEnergyConsumptionDescriptionObject) AiEnergyQuantity() PropertyInterface[float64] { return &self.aiEnergyQuantity }
func (self *AiEnergyConsumptionDescriptionObject) AiEnergyUnit() PropertyInterface[string] { return &self.aiEnergyUnit }

func (self *AiEnergyConsumptionDescriptionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.aiEnergyQuantity.IsSet() {
        val, err := EncodeFloat(self.aiEnergyQuantity.Get(), path.PushPath("aiEnergyQuantity"), aiEnergyConsumptionDescriptionAiEnergyQuantityContext, state)
        if err != nil {
            return err
        }
        data["ai_energyQuantity"] = val
    }
    if self.aiEnergyUnit.IsSet() {
        val, err := EncodeIRI(self.aiEnergyUnit.Get(), path.PushPath("aiEnergyUnit"), aiEnergyConsumptionDescriptionAiEnergyUnitContext, state)
        if err != nil {
            return err
        }
        data["ai_energyUnit"] = val
    }
    return nil
}

// Specifies the unit of energy consumption.
type AiEnergyUnitTypeObject struct {
    SHACLObjectBase

}

// Kilowatt-hour.
const AiEnergyUnitTypeKilowattHour = "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/kilowattHour"
// Megajoule.
const AiEnergyUnitTypeMegajoule = "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/megajoule"
// Any other units of energy measurement.
const AiEnergyUnitTypeOther = "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType/other"

type AiEnergyUnitTypeObjectType struct {
    SHACLTypeBase
}
var aiEnergyUnitTypeType AiEnergyUnitTypeObjectType

func DecodeAiEnergyUnitType (data any, path Path, context map[string]string) (Ref[AiEnergyUnitType], error) {
    return DecodeRef[AiEnergyUnitType](data, path, context, aiEnergyUnitTypeType)
}

func (self AiEnergyUnitTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AiEnergyUnitType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AiEnergyUnitTypeObjectType) Create() SHACLObject {
    return ConstructAiEnergyUnitTypeObject(&AiEnergyUnitTypeObject{}, self)
}

func ConstructAiEnergyUnitTypeObject(o *AiEnergyUnitTypeObject, typ SHACLType) *AiEnergyUnitTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type AiEnergyUnitType interface {
    SHACLObject
}


func MakeAiEnergyUnitType() AiEnergyUnitType {
    return ConstructAiEnergyUnitTypeObject(&AiEnergyUnitTypeObject{}, aiEnergyUnitTypeType)
}

func MakeAiEnergyUnitTypeRef() Ref[AiEnergyUnitType] {
    o := MakeAiEnergyUnitType()
    return MakeObjectRef[AiEnergyUnitType](o)
}

func (self *AiEnergyUnitTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *AiEnergyUnitTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *AiEnergyUnitTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Specifies the safety risk level.
type AiSafetyRiskAssessmentTypeObject struct {
    SHACLObjectBase

}

// The second-highest level of risk posed by an AI system.
const AiSafetyRiskAssessmentTypeHigh = "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/high"
// Low/no risk is posed by an AI system.
const AiSafetyRiskAssessmentTypeLow = "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/low"
// The third-highest level of risk posed by an AI system.
const AiSafetyRiskAssessmentTypeMedium = "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/medium"
// The highest level of risk posed by an AI system.
const AiSafetyRiskAssessmentTypeSerious = "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/serious"

type AiSafetyRiskAssessmentTypeObjectType struct {
    SHACLTypeBase
}
var aiSafetyRiskAssessmentTypeType AiSafetyRiskAssessmentTypeObjectType

func DecodeAiSafetyRiskAssessmentType (data any, path Path, context map[string]string) (Ref[AiSafetyRiskAssessmentType], error) {
    return DecodeRef[AiSafetyRiskAssessmentType](data, path, context, aiSafetyRiskAssessmentTypeType)
}

func (self AiSafetyRiskAssessmentTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AiSafetyRiskAssessmentType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AiSafetyRiskAssessmentTypeObjectType) Create() SHACLObject {
    return ConstructAiSafetyRiskAssessmentTypeObject(&AiSafetyRiskAssessmentTypeObject{}, self)
}

func ConstructAiSafetyRiskAssessmentTypeObject(o *AiSafetyRiskAssessmentTypeObject, typ SHACLType) *AiSafetyRiskAssessmentTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type AiSafetyRiskAssessmentType interface {
    SHACLObject
}


func MakeAiSafetyRiskAssessmentType() AiSafetyRiskAssessmentType {
    return ConstructAiSafetyRiskAssessmentTypeObject(&AiSafetyRiskAssessmentTypeObject{}, aiSafetyRiskAssessmentTypeType)
}

func MakeAiSafetyRiskAssessmentTypeRef() Ref[AiSafetyRiskAssessmentType] {
    o := MakeAiSafetyRiskAssessmentType()
    return MakeObjectRef[AiSafetyRiskAssessmentType](o)
}

func (self *AiSafetyRiskAssessmentTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *AiSafetyRiskAssessmentTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *AiSafetyRiskAssessmentTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Specifies the type of an annotation.
type AnnotationTypeObject struct {
    SHACLObjectBase

}

// Used to store extra information about an Element which is not part of a review (e.g. extra information provided during the creation of the Element).
const AnnotationTypeOther = "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/other"
// Used when someone reviews the Element.
const AnnotationTypeReview = "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/review"

type AnnotationTypeObjectType struct {
    SHACLTypeBase
}
var annotationTypeType AnnotationTypeObjectType

func DecodeAnnotationType (data any, path Path, context map[string]string) (Ref[AnnotationType], error) {
    return DecodeRef[AnnotationType](data, path, context, annotationTypeType)
}

func (self AnnotationTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AnnotationType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AnnotationTypeObjectType) Create() SHACLObject {
    return ConstructAnnotationTypeObject(&AnnotationTypeObject{}, self)
}

func ConstructAnnotationTypeObject(o *AnnotationTypeObject, typ SHACLType) *AnnotationTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type AnnotationType interface {
    SHACLObject
}


func MakeAnnotationType() AnnotationType {
    return ConstructAnnotationTypeObject(&AnnotationTypeObject{}, annotationTypeType)
}

func MakeAnnotationTypeRef() Ref[AnnotationType] {
    o := MakeAnnotationType()
    return MakeObjectRef[AnnotationType](o)
}

func (self *AnnotationTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *AnnotationTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *AnnotationTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Provides information about the creation of the Element.
type CreationInfoObject struct {
    SHACLObjectBase

    // Provide consumers with comments by the creator of the Element about the
    // Element.
    comment Property[string]
    // Identifies when the Element was originally created.
    created Property[time.Time]
    // Identifies who or what created the Element.
    createdBy RefListProperty[Agent]
    // Identifies the tooling that was used during the creation of the Element.
    createdUsing RefListProperty[Tool]
    // Provides a reference number that can be used to understand how to parse and
    // interpret an Element.
    specVersion Property[string]
}


type CreationInfoObjectType struct {
    SHACLTypeBase
}
var creationInfoType CreationInfoObjectType
var creationInfoCommentContext = map[string]string{}
var creationInfoCreatedContext = map[string]string{}
var creationInfoCreatedByContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",}
var creationInfoCreatedUsingContext = map[string]string{}
var creationInfoSpecVersionContext = map[string]string{}

func DecodeCreationInfo (data any, path Path, context map[string]string) (Ref[CreationInfo], error) {
    return DecodeRef[CreationInfo](data, path, context, creationInfoType)
}

func (self CreationInfoObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(CreationInfo)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/comment", "comment":
        val, err := DecodeString(value, path, creationInfoCommentContext)
        if err != nil {
            return false, err
        }
        err = obj.Comment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/created", "created":
        val, err := DecodeDateTimeStamp(value, path, creationInfoCreatedContext)
        if err != nil {
            return false, err
        }
        err = obj.Created().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/createdBy", "createdBy":
        val, err := DecodeList[Ref[Agent]](value, path, creationInfoCreatedByContext, DecodeAgent)
        if err != nil {
            return false, err
        }
        err = obj.CreatedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/createdUsing", "createdUsing":
        val, err := DecodeList[Ref[Tool]](value, path, creationInfoCreatedUsingContext, DecodeTool)
        if err != nil {
            return false, err
        }
        err = obj.CreatedUsing().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/specVersion", "specVersion":
        val, err := DecodeString(value, path, creationInfoSpecVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.SpecVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self CreationInfoObjectType) Create() SHACLObject {
    return ConstructCreationInfoObject(&CreationInfoObject{}, self)
}

func ConstructCreationInfoObject(o *CreationInfoObject, typ SHACLType) *CreationInfoObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.comment = NewProperty[string]("comment", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.created = NewProperty[time.Time]("created", validators)
    }
    {
        validators := []Validator[Ref[Agent]]{}
        o.createdBy = NewRefListProperty[Agent]("createdBy", validators)
    }
    {
        validators := []Validator[Ref[Tool]]{}
        o.createdUsing = NewRefListProperty[Tool]("createdUsing", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators, RegexValidator[string]{`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`})
        o.specVersion = NewProperty[string]("specVersion", validators)
    }
    return o
}

type CreationInfo interface {
    SHACLObject
    Comment() PropertyInterface[string]
    Created() PropertyInterface[time.Time]
    CreatedBy() ListPropertyInterface[Ref[Agent]]
    CreatedUsing() ListPropertyInterface[Ref[Tool]]
    SpecVersion() PropertyInterface[string]
}


func MakeCreationInfo() CreationInfo {
    return ConstructCreationInfoObject(&CreationInfoObject{}, creationInfoType)
}

func MakeCreationInfoRef() Ref[CreationInfo] {
    o := MakeCreationInfo()
    return MakeObjectRef[CreationInfo](o)
}

func (self *CreationInfoObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("comment")
        if ! self.comment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("created")
        if ! self.created.Check(prop_path, handler) {
            valid = false
        }
        if ! self.created.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"created", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("createdBy")
        if ! self.createdBy.Check(prop_path, handler) {
            valid = false
        }
        if len(self.createdBy.Get()) < 1 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "createdBy",
                    "Too few elements. Minimum of 1 required"},
                    prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("createdUsing")
        if ! self.createdUsing.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("specVersion")
        if ! self.specVersion.Check(prop_path, handler) {
            valid = false
        }
        if ! self.specVersion.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"specVersion", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *CreationInfoObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.comment.Walk(path, visit)
    self.created.Walk(path, visit)
    self.createdBy.Walk(path, visit)
    self.createdUsing.Walk(path, visit)
    self.specVersion.Walk(path, visit)
}


func (self *CreationInfoObject) Comment() PropertyInterface[string] { return &self.comment }
func (self *CreationInfoObject) Created() PropertyInterface[time.Time] { return &self.created }
func (self *CreationInfoObject) CreatedBy() ListPropertyInterface[Ref[Agent]] { return &self.createdBy }
func (self *CreationInfoObject) CreatedUsing() ListPropertyInterface[Ref[Tool]] { return &self.createdUsing }
func (self *CreationInfoObject) SpecVersion() PropertyInterface[string] { return &self.specVersion }

func (self *CreationInfoObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.comment.IsSet() {
        val, err := EncodeString(self.comment.Get(), path.PushPath("comment"), creationInfoCommentContext, state)
        if err != nil {
            return err
        }
        data["comment"] = val
    }
    if self.created.IsSet() {
        val, err := EncodeDateTime(self.created.Get(), path.PushPath("created"), creationInfoCreatedContext, state)
        if err != nil {
            return err
        }
        data["created"] = val
    }
    if self.createdBy.IsSet() {
        val, err := EncodeList[Ref[Agent]](self.createdBy.Get(), path.PushPath("createdBy"), creationInfoCreatedByContext, state, EncodeRef[Agent])
        if err != nil {
            return err
        }
        data["createdBy"] = val
    }
    if self.createdUsing.IsSet() {
        val, err := EncodeList[Ref[Tool]](self.createdUsing.Get(), path.PushPath("createdUsing"), creationInfoCreatedUsingContext, state, EncodeRef[Tool])
        if err != nil {
            return err
        }
        data["createdUsing"] = val
    }
    if self.specVersion.IsSet() {
        val, err := EncodeString(self.specVersion.Get(), path.PushPath("specVersion"), creationInfoSpecVersionContext, state)
        if err != nil {
            return err
        }
        data["specVersion"] = val
    }
    return nil
}

// A key with an associated value.
type DictionaryEntryObject struct {
    SHACLObjectBase

    // A key used in a generic key-value pair.
    key Property[string]
    // A value used in a generic key-value pair.
    value Property[string]
}


type DictionaryEntryObjectType struct {
    SHACLTypeBase
}
var dictionaryEntryType DictionaryEntryObjectType
var dictionaryEntryKeyContext = map[string]string{}
var dictionaryEntryValueContext = map[string]string{}

func DecodeDictionaryEntry (data any, path Path, context map[string]string) (Ref[DictionaryEntry], error) {
    return DecodeRef[DictionaryEntry](data, path, context, dictionaryEntryType)
}

func (self DictionaryEntryObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(DictionaryEntry)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/key", "key":
        val, err := DecodeString(value, path, dictionaryEntryKeyContext)
        if err != nil {
            return false, err
        }
        err = obj.Key().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/value", "value":
        val, err := DecodeString(value, path, dictionaryEntryValueContext)
        if err != nil {
            return false, err
        }
        err = obj.Value().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self DictionaryEntryObjectType) Create() SHACLObject {
    return ConstructDictionaryEntryObject(&DictionaryEntryObject{}, self)
}

func ConstructDictionaryEntryObject(o *DictionaryEntryObject, typ SHACLType) *DictionaryEntryObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.key = NewProperty[string]("key", validators)
    }
    {
        validators := []Validator[string]{}
        o.value = NewProperty[string]("value", validators)
    }
    return o
}

type DictionaryEntry interface {
    SHACLObject
    Key() PropertyInterface[string]
    Value() PropertyInterface[string]
}


func MakeDictionaryEntry() DictionaryEntry {
    return ConstructDictionaryEntryObject(&DictionaryEntryObject{}, dictionaryEntryType)
}

func MakeDictionaryEntryRef() Ref[DictionaryEntry] {
    o := MakeDictionaryEntry()
    return MakeObjectRef[DictionaryEntry](o)
}

func (self *DictionaryEntryObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("key")
        if ! self.key.Check(prop_path, handler) {
            valid = false
        }
        if ! self.key.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"key", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("value")
        if ! self.value.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *DictionaryEntryObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.key.Walk(path, visit)
    self.value.Walk(path, visit)
}


func (self *DictionaryEntryObject) Key() PropertyInterface[string] { return &self.key }
func (self *DictionaryEntryObject) Value() PropertyInterface[string] { return &self.value }

func (self *DictionaryEntryObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.key.IsSet() {
        val, err := EncodeString(self.key.Get(), path.PushPath("key"), dictionaryEntryKeyContext, state)
        if err != nil {
            return err
        }
        data["key"] = val
    }
    if self.value.IsSet() {
        val, err := EncodeString(self.value.Get(), path.PushPath("value"), dictionaryEntryValueContext, state)
        if err != nil {
            return err
        }
        data["value"] = val
    }
    return nil
}

// Base domain class from which all other SPDX-3.0 domain classes derive.
type ElementObject struct {
    SHACLObjectBase

    // Provide consumers with comments by the creator of the Element about the
    // Element.
    comment Property[string]
    // Provides information about the creation of the Element.
    creationInfo RefProperty[CreationInfo]
    // Provides a detailed description of the Element.
    description Property[string]
    // Specifies an Extension characterization of some aspect of an Element.
    extension RefListProperty[ExtensionExtension]
    // Provides a reference to a resource outside the scope of SPDX-3.0 content
    // that uniquely identifies an Element.
    externalIdentifier RefListProperty[ExternalIdentifier]
    // Points to a resource outside the scope of the SPDX-3.0 content
    // that provides additional characteristics of an Element.
    externalRef RefListProperty[ExternalRef]
    // Identifies the name of an Element as designated by the creator.
    name Property[string]
    // A short description of an Element.
    summary Property[string]
    // Provides an IntegrityMethod with which the integrity of an Element can be
    // asserted.
    verifiedUsing RefListProperty[IntegrityMethod]
}


type ElementObjectType struct {
    SHACLTypeBase
}
var elementType ElementObjectType
var elementCommentContext = map[string]string{}
var elementCreationInfoContext = map[string]string{}
var elementDescriptionContext = map[string]string{}
var elementExtensionContext = map[string]string{}
var elementExternalIdentifierContext = map[string]string{}
var elementExternalRefContext = map[string]string{}
var elementNameContext = map[string]string{}
var elementSummaryContext = map[string]string{}
var elementVerifiedUsingContext = map[string]string{}

func DecodeElement (data any, path Path, context map[string]string) (Ref[Element], error) {
    return DecodeRef[Element](data, path, context, elementType)
}

func (self ElementObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Element)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/comment", "comment":
        val, err := DecodeString(value, path, elementCommentContext)
        if err != nil {
            return false, err
        }
        err = obj.Comment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/creationInfo", "creationInfo":
        val, err := DecodeCreationInfo(value, path, elementCreationInfoContext)
        if err != nil {
            return false, err
        }
        err = obj.CreationInfo().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/description", "description":
        val, err := DecodeString(value, path, elementDescriptionContext)
        if err != nil {
            return false, err
        }
        err = obj.Description().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/extension", "extension":
        val, err := DecodeList[Ref[ExtensionExtension]](value, path, elementExtensionContext, DecodeExtensionExtension)
        if err != nil {
            return false, err
        }
        err = obj.Extension().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/externalIdentifier", "externalIdentifier":
        val, err := DecodeList[Ref[ExternalIdentifier]](value, path, elementExternalIdentifierContext, DecodeExternalIdentifier)
        if err != nil {
            return false, err
        }
        err = obj.ExternalIdentifier().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/externalRef", "externalRef":
        val, err := DecodeList[Ref[ExternalRef]](value, path, elementExternalRefContext, DecodeExternalRef)
        if err != nil {
            return false, err
        }
        err = obj.ExternalRef().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/name", "name":
        val, err := DecodeString(value, path, elementNameContext)
        if err != nil {
            return false, err
        }
        err = obj.Name().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/summary", "summary":
        val, err := DecodeString(value, path, elementSummaryContext)
        if err != nil {
            return false, err
        }
        err = obj.Summary().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/verifiedUsing", "verifiedUsing":
        val, err := DecodeList[Ref[IntegrityMethod]](value, path, elementVerifiedUsingContext, DecodeIntegrityMethod)
        if err != nil {
            return false, err
        }
        err = obj.VerifiedUsing().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ElementObjectType) Create() SHACLObject {
    return ConstructElementObject(&ElementObject{}, self)
}

func ConstructElementObject(o *ElementObject, typ SHACLType) *ElementObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.comment = NewProperty[string]("comment", validators)
    }
    {
        validators := []Validator[Ref[CreationInfo]]{}
        o.creationInfo = NewRefProperty[CreationInfo]("creationInfo", validators)
    }
    {
        validators := []Validator[string]{}
        o.description = NewProperty[string]("description", validators)
    }
    {
        validators := []Validator[Ref[ExtensionExtension]]{}
        o.extension = NewRefListProperty[ExtensionExtension]("extension", validators)
    }
    {
        validators := []Validator[Ref[ExternalIdentifier]]{}
        o.externalIdentifier = NewRefListProperty[ExternalIdentifier]("externalIdentifier", validators)
    }
    {
        validators := []Validator[Ref[ExternalRef]]{}
        o.externalRef = NewRefListProperty[ExternalRef]("externalRef", validators)
    }
    {
        validators := []Validator[string]{}
        o.name = NewProperty[string]("name", validators)
    }
    {
        validators := []Validator[string]{}
        o.summary = NewProperty[string]("summary", validators)
    }
    {
        validators := []Validator[Ref[IntegrityMethod]]{}
        o.verifiedUsing = NewRefListProperty[IntegrityMethod]("verifiedUsing", validators)
    }
    return o
}

type Element interface {
    SHACLObject
    Comment() PropertyInterface[string]
    CreationInfo() RefPropertyInterface[CreationInfo]
    Description() PropertyInterface[string]
    Extension() ListPropertyInterface[Ref[ExtensionExtension]]
    ExternalIdentifier() ListPropertyInterface[Ref[ExternalIdentifier]]
    ExternalRef() ListPropertyInterface[Ref[ExternalRef]]
    Name() PropertyInterface[string]
    Summary() PropertyInterface[string]
    VerifiedUsing() ListPropertyInterface[Ref[IntegrityMethod]]
}



func (self *ElementObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("comment")
        if ! self.comment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("creationInfo")
        if ! self.creationInfo.Check(prop_path, handler) {
            valid = false
        }
        if ! self.creationInfo.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"creationInfo", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("description")
        if ! self.description.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("extension")
        if ! self.extension.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("externalIdentifier")
        if ! self.externalIdentifier.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("externalRef")
        if ! self.externalRef.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("name")
        if ! self.name.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("summary")
        if ! self.summary.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("verifiedUsing")
        if ! self.verifiedUsing.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ElementObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.comment.Walk(path, visit)
    self.creationInfo.Walk(path, visit)
    self.description.Walk(path, visit)
    self.extension.Walk(path, visit)
    self.externalIdentifier.Walk(path, visit)
    self.externalRef.Walk(path, visit)
    self.name.Walk(path, visit)
    self.summary.Walk(path, visit)
    self.verifiedUsing.Walk(path, visit)
}


func (self *ElementObject) Comment() PropertyInterface[string] { return &self.comment }
func (self *ElementObject) CreationInfo() RefPropertyInterface[CreationInfo] { return &self.creationInfo }
func (self *ElementObject) Description() PropertyInterface[string] { return &self.description }
func (self *ElementObject) Extension() ListPropertyInterface[Ref[ExtensionExtension]] { return &self.extension }
func (self *ElementObject) ExternalIdentifier() ListPropertyInterface[Ref[ExternalIdentifier]] { return &self.externalIdentifier }
func (self *ElementObject) ExternalRef() ListPropertyInterface[Ref[ExternalRef]] { return &self.externalRef }
func (self *ElementObject) Name() PropertyInterface[string] { return &self.name }
func (self *ElementObject) Summary() PropertyInterface[string] { return &self.summary }
func (self *ElementObject) VerifiedUsing() ListPropertyInterface[Ref[IntegrityMethod]] { return &self.verifiedUsing }

func (self *ElementObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.comment.IsSet() {
        val, err := EncodeString(self.comment.Get(), path.PushPath("comment"), elementCommentContext, state)
        if err != nil {
            return err
        }
        data["comment"] = val
    }
    if self.creationInfo.IsSet() {
        val, err := EncodeRef[CreationInfo](self.creationInfo.Get(), path.PushPath("creationInfo"), elementCreationInfoContext, state)
        if err != nil {
            return err
        }
        data["creationInfo"] = val
    }
    if self.description.IsSet() {
        val, err := EncodeString(self.description.Get(), path.PushPath("description"), elementDescriptionContext, state)
        if err != nil {
            return err
        }
        data["description"] = val
    }
    if self.extension.IsSet() {
        val, err := EncodeList[Ref[ExtensionExtension]](self.extension.Get(), path.PushPath("extension"), elementExtensionContext, state, EncodeRef[ExtensionExtension])
        if err != nil {
            return err
        }
        data["extension"] = val
    }
    if self.externalIdentifier.IsSet() {
        val, err := EncodeList[Ref[ExternalIdentifier]](self.externalIdentifier.Get(), path.PushPath("externalIdentifier"), elementExternalIdentifierContext, state, EncodeRef[ExternalIdentifier])
        if err != nil {
            return err
        }
        data["externalIdentifier"] = val
    }
    if self.externalRef.IsSet() {
        val, err := EncodeList[Ref[ExternalRef]](self.externalRef.Get(), path.PushPath("externalRef"), elementExternalRefContext, state, EncodeRef[ExternalRef])
        if err != nil {
            return err
        }
        data["externalRef"] = val
    }
    if self.name.IsSet() {
        val, err := EncodeString(self.name.Get(), path.PushPath("name"), elementNameContext, state)
        if err != nil {
            return err
        }
        data["name"] = val
    }
    if self.summary.IsSet() {
        val, err := EncodeString(self.summary.Get(), path.PushPath("summary"), elementSummaryContext, state)
        if err != nil {
            return err
        }
        data["summary"] = val
    }
    if self.verifiedUsing.IsSet() {
        val, err := EncodeList[Ref[IntegrityMethod]](self.verifiedUsing.Get(), path.PushPath("verifiedUsing"), elementVerifiedUsingContext, state, EncodeRef[IntegrityMethod])
        if err != nil {
            return err
        }
        data["verifiedUsing"] = val
    }
    return nil
}

// A collection of Elements, not necessarily with unifying context.
type ElementCollectionObject struct {
    ElementObject

    // Refers to one or more Elements that are part of an ElementCollection.
    element RefListProperty[Element]
    // Describes one a profile which the creator of this ElementCollection intends to
    // conform to.
    profileConformance ListProperty[string]
    // This property is used to denote the root Element(s) of a tree of elements contained in a BOM.
    rootElement RefListProperty[Element]
}


type ElementCollectionObjectType struct {
    SHACLTypeBase
}
var elementCollectionType ElementCollectionObjectType
var elementCollectionElementContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement": "NoAssertionElement",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement": "NoneElement",
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}
var elementCollectionProfileConformanceContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/ai": "ai",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/build": "build",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/core": "core",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/dataset": "dataset",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/expandedLicensing": "expandedLicensing",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/extension": "extension",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/lite": "lite",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/security": "security",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/simpleLicensing": "simpleLicensing",
    "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/software": "software",}
var elementCollectionRootElementContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement": "NoAssertionElement",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement": "NoneElement",
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}

func DecodeElementCollection (data any, path Path, context map[string]string) (Ref[ElementCollection], error) {
    return DecodeRef[ElementCollection](data, path, context, elementCollectionType)
}

func (self ElementCollectionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ElementCollection)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/element", "element":
        val, err := DecodeList[Ref[Element]](value, path, elementCollectionElementContext, DecodeElement)
        if err != nil {
            return false, err
        }
        err = obj.Element().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/profileConformance", "profileConformance":
        val, err := DecodeList[string](value, path, elementCollectionProfileConformanceContext, DecodeIRI)
        if err != nil {
            return false, err
        }
        err = obj.ProfileConformance().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/rootElement", "rootElement":
        val, err := DecodeList[Ref[Element]](value, path, elementCollectionRootElementContext, DecodeElement)
        if err != nil {
            return false, err
        }
        err = obj.RootElement().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ElementCollectionObjectType) Create() SHACLObject {
    return ConstructElementCollectionObject(&ElementCollectionObject{}, self)
}

func ConstructElementCollectionObject(o *ElementCollectionObject, typ SHACLType) *ElementCollectionObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[Ref[Element]]{}
        o.element = NewRefListProperty[Element]("element", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/ai",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/build",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/core",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/dataset",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/expandedLicensing",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/extension",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/lite",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/security",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/simpleLicensing",
                "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/software",
        }})
        o.profileConformance = NewListProperty[string]("profileConformance", validators)
    }
    {
        validators := []Validator[Ref[Element]]{}
        o.rootElement = NewRefListProperty[Element]("rootElement", validators)
    }
    return o
}

type ElementCollection interface {
    Element
    Element() ListPropertyInterface[Ref[Element]]
    ProfileConformance() ListPropertyInterface[string]
    RootElement() ListPropertyInterface[Ref[Element]]
}



func (self *ElementCollectionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("element")
        if ! self.element.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("profileConformance")
        if ! self.profileConformance.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("rootElement")
        if ! self.rootElement.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ElementCollectionObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.element.Walk(path, visit)
    self.profileConformance.Walk(path, visit)
    self.rootElement.Walk(path, visit)
}


func (self *ElementCollectionObject) Element() ListPropertyInterface[Ref[Element]] { return &self.element }
func (self *ElementCollectionObject) ProfileConformance() ListPropertyInterface[string] { return &self.profileConformance }
func (self *ElementCollectionObject) RootElement() ListPropertyInterface[Ref[Element]] { return &self.rootElement }

func (self *ElementCollectionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.element.IsSet() {
        val, err := EncodeList[Ref[Element]](self.element.Get(), path.PushPath("element"), elementCollectionElementContext, state, EncodeRef[Element])
        if err != nil {
            return err
        }
        data["element"] = val
    }
    if self.profileConformance.IsSet() {
        val, err := EncodeList[string](self.profileConformance.Get(), path.PushPath("profileConformance"), elementCollectionProfileConformanceContext, state, EncodeIRI)
        if err != nil {
            return err
        }
        data["profileConformance"] = val
    }
    if self.rootElement.IsSet() {
        val, err := EncodeList[Ref[Element]](self.rootElement.Get(), path.PushPath("rootElement"), elementCollectionRootElementContext, state, EncodeRef[Element])
        if err != nil {
            return err
        }
        data["rootElement"] = val
    }
    return nil
}

// A reference to a resource identifier defined outside the scope of SPDX-3.0 content that uniquely identifies an Element.
type ExternalIdentifierObject struct {
    SHACLObjectBase

    // Provide consumers with comments by the creator of the Element about the
    // Element.
    comment Property[string]
    // Specifies the type of the external identifier.
    externalIdentifierType Property[string]
    // Uniquely identifies an external element.
    identifier Property[string]
    // Provides the location for more information regarding an external identifier.
    identifierLocator ListProperty[string]
    // An entity that is authorized to issue identification credentials.
    issuingAuthority Property[string]
}


type ExternalIdentifierObjectType struct {
    SHACLTypeBase
}
var externalIdentifierType ExternalIdentifierObjectType
var externalIdentifierCommentContext = map[string]string{}
var externalIdentifierExternalIdentifierTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe22": "cpe22",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe23": "cpe23",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cve": "cve",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/email": "email",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/gitoid": "gitoid",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/packageUrl": "packageUrl",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/securityOther": "securityOther",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swhid": "swhid",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swid": "swid",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/urlScheme": "urlScheme",}
var externalIdentifierIdentifierContext = map[string]string{}
var externalIdentifierIdentifierLocatorContext = map[string]string{}
var externalIdentifierIssuingAuthorityContext = map[string]string{}

func DecodeExternalIdentifier (data any, path Path, context map[string]string) (Ref[ExternalIdentifier], error) {
    return DecodeRef[ExternalIdentifier](data, path, context, externalIdentifierType)
}

func (self ExternalIdentifierObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExternalIdentifier)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/comment", "comment":
        val, err := DecodeString(value, path, externalIdentifierCommentContext)
        if err != nil {
            return false, err
        }
        err = obj.Comment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/externalIdentifierType", "externalIdentifierType":
        val, err := DecodeIRI(value, path, externalIdentifierExternalIdentifierTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.ExternalIdentifierType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/identifier", "identifier":
        val, err := DecodeString(value, path, externalIdentifierIdentifierContext)
        if err != nil {
            return false, err
        }
        err = obj.Identifier().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/identifierLocator", "identifierLocator":
        val, err := DecodeList[string](value, path, externalIdentifierIdentifierLocatorContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.IdentifierLocator().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/issuingAuthority", "issuingAuthority":
        val, err := DecodeString(value, path, externalIdentifierIssuingAuthorityContext)
        if err != nil {
            return false, err
        }
        err = obj.IssuingAuthority().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExternalIdentifierObjectType) Create() SHACLObject {
    return ConstructExternalIdentifierObject(&ExternalIdentifierObject{}, self)
}

func ConstructExternalIdentifierObject(o *ExternalIdentifierObject, typ SHACLType) *ExternalIdentifierObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.comment = NewProperty[string]("comment", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe22",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe23",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cve",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/email",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/gitoid",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/packageUrl",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/securityOther",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swhid",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swid",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/urlScheme",
        }})
        o.externalIdentifierType = NewProperty[string]("externalIdentifierType", validators)
    }
    {
        validators := []Validator[string]{}
        o.identifier = NewProperty[string]("identifier", validators)
    }
    {
        validators := []Validator[string]{}
        o.identifierLocator = NewListProperty[string]("identifierLocator", validators)
    }
    {
        validators := []Validator[string]{}
        o.issuingAuthority = NewProperty[string]("issuingAuthority", validators)
    }
    return o
}

type ExternalIdentifier interface {
    SHACLObject
    Comment() PropertyInterface[string]
    ExternalIdentifierType() PropertyInterface[string]
    Identifier() PropertyInterface[string]
    IdentifierLocator() ListPropertyInterface[string]
    IssuingAuthority() PropertyInterface[string]
}


func MakeExternalIdentifier() ExternalIdentifier {
    return ConstructExternalIdentifierObject(&ExternalIdentifierObject{}, externalIdentifierType)
}

func MakeExternalIdentifierRef() Ref[ExternalIdentifier] {
    o := MakeExternalIdentifier()
    return MakeObjectRef[ExternalIdentifier](o)
}

func (self *ExternalIdentifierObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("comment")
        if ! self.comment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("externalIdentifierType")
        if ! self.externalIdentifierType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.externalIdentifierType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"externalIdentifierType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("identifier")
        if ! self.identifier.Check(prop_path, handler) {
            valid = false
        }
        if ! self.identifier.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"identifier", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("identifierLocator")
        if ! self.identifierLocator.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("issuingAuthority")
        if ! self.issuingAuthority.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExternalIdentifierObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.comment.Walk(path, visit)
    self.externalIdentifierType.Walk(path, visit)
    self.identifier.Walk(path, visit)
    self.identifierLocator.Walk(path, visit)
    self.issuingAuthority.Walk(path, visit)
}


func (self *ExternalIdentifierObject) Comment() PropertyInterface[string] { return &self.comment }
func (self *ExternalIdentifierObject) ExternalIdentifierType() PropertyInterface[string] { return &self.externalIdentifierType }
func (self *ExternalIdentifierObject) Identifier() PropertyInterface[string] { return &self.identifier }
func (self *ExternalIdentifierObject) IdentifierLocator() ListPropertyInterface[string] { return &self.identifierLocator }
func (self *ExternalIdentifierObject) IssuingAuthority() PropertyInterface[string] { return &self.issuingAuthority }

func (self *ExternalIdentifierObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.comment.IsSet() {
        val, err := EncodeString(self.comment.Get(), path.PushPath("comment"), externalIdentifierCommentContext, state)
        if err != nil {
            return err
        }
        data["comment"] = val
    }
    if self.externalIdentifierType.IsSet() {
        val, err := EncodeIRI(self.externalIdentifierType.Get(), path.PushPath("externalIdentifierType"), externalIdentifierExternalIdentifierTypeContext, state)
        if err != nil {
            return err
        }
        data["externalIdentifierType"] = val
    }
    if self.identifier.IsSet() {
        val, err := EncodeString(self.identifier.Get(), path.PushPath("identifier"), externalIdentifierIdentifierContext, state)
        if err != nil {
            return err
        }
        data["identifier"] = val
    }
    if self.identifierLocator.IsSet() {
        val, err := EncodeList[string](self.identifierLocator.Get(), path.PushPath("identifierLocator"), externalIdentifierIdentifierLocatorContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["identifierLocator"] = val
    }
    if self.issuingAuthority.IsSet() {
        val, err := EncodeString(self.issuingAuthority.Get(), path.PushPath("issuingAuthority"), externalIdentifierIssuingAuthorityContext, state)
        if err != nil {
            return err
        }
        data["issuingAuthority"] = val
    }
    return nil
}

// Specifies the type of an external identifier.
type ExternalIdentifierTypeObject struct {
    SHACLObjectBase

}

// [Common Platform Enumeration Specification 2.2](https://cpe.mitre.org/files/cpe-specification_2.2.pdf)
const ExternalIdentifierTypeCpe22 = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe22"
// [Common Platform Enumeration: Naming Specification Version 2.3](https://csrc.nist.gov/publications/detail/nistir/7695/final)
const ExternalIdentifierTypeCpe23 = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cpe23"
// Common Vulnerabilities and Exposures identifiers, an identifier for a specific software flaw defined within the official CVE Dictionary and that conforms to the [CVE specification](https://csrc.nist.gov/glossary/term/cve_id).
const ExternalIdentifierTypeCve = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/cve"
// Email address, as defined in [RFC 3696](https://datatracker.ietf.org/doc/rfc3986/) Section 3.
const ExternalIdentifierTypeEmail = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/email"
// [Gitoid](https://www.iana.org/assignments/uri-schemes/prov/gitoid), stands for [Git Object ID](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects). A gitoid of type blob is a unique hash of a binary artifact. A gitoid may represent either an [Artifact Identifier](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-identifier-types) for the software artifact or an [Input Manifest Identifier](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#input-manifest-identifier) for the software artifact's associated [Artifact Input Manifest](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-input-manifest); this ambiguity exists because the Artifact Input Manifest is itself an artifact, and the gitoid of that artifact is its valid identifier. Gitoids calculated on software artifacts (Snippet, File, or Package Elements) should be recorded in the SPDX 3.0 SoftwareArtifact's contentIdentifier property. Gitoids calculated on the Artifact Input Manifest (Input Manifest Identifier) should be recorded in the SPDX 3.0 Element's externalIdentifier property. See [OmniBOR Specification](https://github.com/omnibor/spec/), a minimalistic specification for describing software [Artifact Dependency Graphs](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-dependency-graph-adg).
const ExternalIdentifierTypeGitoid = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/gitoid"
// Used when the type does not match any of the other options.
const ExternalIdentifierTypeOther = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/other"
// Package URL, as defined in the corresponding [Annex](../../../annexes/pkg-url-specification.md) of this specification.
const ExternalIdentifierTypePackageUrl = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/packageUrl"
// Used when there is a security related identifier of unspecified type.
const ExternalIdentifierTypeSecurityOther = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/securityOther"
// SoftWare Hash IDentifier, a persistent intrinsic identifier for digital artifacts, such as files, trees (also known as directories or folders), commits, and other objects typically found in version control systems. The format of the identifiers is defined in the [SWHID specification](https://www.swhid.org/specification/v1.1/4.Syntax) (ISO/IEC DIS 18670). They typically look like `swh:1:cnt:94a9ed024d3859793618152ea559a168bbcbb5e2`.
const ExternalIdentifierTypeSwhid = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swhid"
// Concise Software Identification (CoSWID) tag, as defined in [RFC 9393](https://datatracker.ietf.org/doc/rfc9393/) Section 2.3.
const ExternalIdentifierTypeSwid = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/swid"
// [Uniform Resource Identifier (URI) Schemes](https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml). The scheme used in order to locate a resource.
const ExternalIdentifierTypeUrlScheme = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType/urlScheme"

type ExternalIdentifierTypeObjectType struct {
    SHACLTypeBase
}
var externalIdentifierTypeType ExternalIdentifierTypeObjectType

func DecodeExternalIdentifierType (data any, path Path, context map[string]string) (Ref[ExternalIdentifierType], error) {
    return DecodeRef[ExternalIdentifierType](data, path, context, externalIdentifierTypeType)
}

func (self ExternalIdentifierTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExternalIdentifierType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExternalIdentifierTypeObjectType) Create() SHACLObject {
    return ConstructExternalIdentifierTypeObject(&ExternalIdentifierTypeObject{}, self)
}

func ConstructExternalIdentifierTypeObject(o *ExternalIdentifierTypeObject, typ SHACLType) *ExternalIdentifierTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type ExternalIdentifierType interface {
    SHACLObject
}


func MakeExternalIdentifierType() ExternalIdentifierType {
    return ConstructExternalIdentifierTypeObject(&ExternalIdentifierTypeObject{}, externalIdentifierTypeType)
}

func MakeExternalIdentifierTypeRef() Ref[ExternalIdentifierType] {
    o := MakeExternalIdentifierType()
    return MakeObjectRef[ExternalIdentifierType](o)
}

func (self *ExternalIdentifierTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExternalIdentifierTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *ExternalIdentifierTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A map of Element identifiers that are used within an SpdxDocument but defined

// external to that SpdxDocument.
type ExternalMapObject struct {
    SHACLObjectBase

    // Artifact representing a serialization instance of SPDX data containing the
    // definition of a particular Element.
    definingArtifact RefProperty[Artifact]
    // Identifies an external Element used within an SpdxDocument but defined
    // external to that SpdxDocument.
    externalSpdxId Property[string]
    // Provides an indication of where to retrieve an external Element.
    locationHint Property[string]
    // Provides an IntegrityMethod with which the integrity of an Element can be
    // asserted.
    verifiedUsing RefListProperty[IntegrityMethod]
}


type ExternalMapObjectType struct {
    SHACLTypeBase
}
var externalMapType ExternalMapObjectType
var externalMapDefiningArtifactContext = map[string]string{}
var externalMapExternalSpdxIdContext = map[string]string{}
var externalMapLocationHintContext = map[string]string{}
var externalMapVerifiedUsingContext = map[string]string{}

func DecodeExternalMap (data any, path Path, context map[string]string) (Ref[ExternalMap], error) {
    return DecodeRef[ExternalMap](data, path, context, externalMapType)
}

func (self ExternalMapObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExternalMap)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/definingArtifact", "definingArtifact":
        val, err := DecodeArtifact(value, path, externalMapDefiningArtifactContext)
        if err != nil {
            return false, err
        }
        err = obj.DefiningArtifact().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/externalSpdxId", "externalSpdxId":
        val, err := DecodeString(value, path, externalMapExternalSpdxIdContext)
        if err != nil {
            return false, err
        }
        err = obj.ExternalSpdxId().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/locationHint", "locationHint":
        val, err := DecodeString(value, path, externalMapLocationHintContext)
        if err != nil {
            return false, err
        }
        err = obj.LocationHint().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/verifiedUsing", "verifiedUsing":
        val, err := DecodeList[Ref[IntegrityMethod]](value, path, externalMapVerifiedUsingContext, DecodeIntegrityMethod)
        if err != nil {
            return false, err
        }
        err = obj.VerifiedUsing().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExternalMapObjectType) Create() SHACLObject {
    return ConstructExternalMapObject(&ExternalMapObject{}, self)
}

func ConstructExternalMapObject(o *ExternalMapObject, typ SHACLType) *ExternalMapObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[Ref[Artifact]]{}
        o.definingArtifact = NewRefProperty[Artifact]("definingArtifact", validators)
    }
    {
        validators := []Validator[string]{}
        o.externalSpdxId = NewProperty[string]("externalSpdxId", validators)
    }
    {
        validators := []Validator[string]{}
        o.locationHint = NewProperty[string]("locationHint", validators)
    }
    {
        validators := []Validator[Ref[IntegrityMethod]]{}
        o.verifiedUsing = NewRefListProperty[IntegrityMethod]("verifiedUsing", validators)
    }
    return o
}

type ExternalMap interface {
    SHACLObject
    DefiningArtifact() RefPropertyInterface[Artifact]
    ExternalSpdxId() PropertyInterface[string]
    LocationHint() PropertyInterface[string]
    VerifiedUsing() ListPropertyInterface[Ref[IntegrityMethod]]
}


func MakeExternalMap() ExternalMap {
    return ConstructExternalMapObject(&ExternalMapObject{}, externalMapType)
}

func MakeExternalMapRef() Ref[ExternalMap] {
    o := MakeExternalMap()
    return MakeObjectRef[ExternalMap](o)
}

func (self *ExternalMapObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("definingArtifact")
        if ! self.definingArtifact.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("externalSpdxId")
        if ! self.externalSpdxId.Check(prop_path, handler) {
            valid = false
        }
        if ! self.externalSpdxId.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"externalSpdxId", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("locationHint")
        if ! self.locationHint.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("verifiedUsing")
        if ! self.verifiedUsing.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExternalMapObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.definingArtifact.Walk(path, visit)
    self.externalSpdxId.Walk(path, visit)
    self.locationHint.Walk(path, visit)
    self.verifiedUsing.Walk(path, visit)
}


func (self *ExternalMapObject) DefiningArtifact() RefPropertyInterface[Artifact] { return &self.definingArtifact }
func (self *ExternalMapObject) ExternalSpdxId() PropertyInterface[string] { return &self.externalSpdxId }
func (self *ExternalMapObject) LocationHint() PropertyInterface[string] { return &self.locationHint }
func (self *ExternalMapObject) VerifiedUsing() ListPropertyInterface[Ref[IntegrityMethod]] { return &self.verifiedUsing }

func (self *ExternalMapObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.definingArtifact.IsSet() {
        val, err := EncodeRef[Artifact](self.definingArtifact.Get(), path.PushPath("definingArtifact"), externalMapDefiningArtifactContext, state)
        if err != nil {
            return err
        }
        data["definingArtifact"] = val
    }
    if self.externalSpdxId.IsSet() {
        val, err := EncodeString(self.externalSpdxId.Get(), path.PushPath("externalSpdxId"), externalMapExternalSpdxIdContext, state)
        if err != nil {
            return err
        }
        data["externalSpdxId"] = val
    }
    if self.locationHint.IsSet() {
        val, err := EncodeString(self.locationHint.Get(), path.PushPath("locationHint"), externalMapLocationHintContext, state)
        if err != nil {
            return err
        }
        data["locationHint"] = val
    }
    if self.verifiedUsing.IsSet() {
        val, err := EncodeList[Ref[IntegrityMethod]](self.verifiedUsing.Get(), path.PushPath("verifiedUsing"), externalMapVerifiedUsingContext, state, EncodeRef[IntegrityMethod])
        if err != nil {
            return err
        }
        data["verifiedUsing"] = val
    }
    return nil
}

// A reference to a resource outside the scope of SPDX-3.0 content related to an Element.
type ExternalRefObject struct {
    SHACLObjectBase

    // Provide consumers with comments by the creator of the Element about the
    // Element.
    comment Property[string]
    // Provides information about the content type of an Element or a Property.
    contentType Property[string]
    // Specifies the type of the external reference.
    externalRefType Property[string]
    // Provides the location of an external reference.
    locator ListProperty[string]
}


type ExternalRefObjectType struct {
    SHACLTypeBase
}
var externalRefType ExternalRefObjectType
var externalRefCommentContext = map[string]string{}
var externalRefContentTypeContext = map[string]string{}
var externalRefExternalRefTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altDownloadLocation": "altDownloadLocation",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altWebPage": "altWebPage",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/binaryArtifact": "binaryArtifact",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/bower": "bower",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildMeta": "buildMeta",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildSystem": "buildSystem",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/certificationReport": "certificationReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/chat": "chat",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/componentAnalysisReport": "componentAnalysisReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/cwe": "cwe",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/documentation": "documentation",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/dynamicAnalysisReport": "dynamicAnalysisReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/eolNotice": "eolNotice",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/exportControlAssessment": "exportControlAssessment",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/funding": "funding",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/issueTracker": "issueTracker",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/license": "license",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mailingList": "mailingList",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mavenCentral": "mavenCentral",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/metrics": "metrics",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/npm": "npm",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/nuget": "nuget",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/privacyAssessment": "privacyAssessment",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/productMetadata": "productMetadata",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/purchaseOrder": "purchaseOrder",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/qualityAssessmentReport": "qualityAssessmentReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseHistory": "releaseHistory",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseNotes": "releaseNotes",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/riskAssessment": "riskAssessment",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/runtimeAnalysisReport": "runtimeAnalysisReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/secureSoftwareAttestation": "secureSoftwareAttestation",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdversaryModel": "securityAdversaryModel",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdvisory": "securityAdvisory",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityFix": "securityFix",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityOther": "securityOther",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPenTestReport": "securityPenTestReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPolicy": "securityPolicy",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityThreatModel": "securityThreatModel",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/socialMedia": "socialMedia",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/sourceArtifact": "sourceArtifact",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/staticAnalysisReport": "staticAnalysisReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/support": "support",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vcs": "vcs",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityDisclosureReport": "vulnerabilityDisclosureReport",
    "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityExploitabilityAssessment": "vulnerabilityExploitabilityAssessment",}
var externalRefLocatorContext = map[string]string{}

func DecodeExternalRef (data any, path Path, context map[string]string) (Ref[ExternalRef], error) {
    return DecodeRef[ExternalRef](data, path, context, externalRefType)
}

func (self ExternalRefObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExternalRef)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/comment", "comment":
        val, err := DecodeString(value, path, externalRefCommentContext)
        if err != nil {
            return false, err
        }
        err = obj.Comment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/contentType", "contentType":
        val, err := DecodeString(value, path, externalRefContentTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.ContentType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/externalRefType", "externalRefType":
        val, err := DecodeIRI(value, path, externalRefExternalRefTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.ExternalRefType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/locator", "locator":
        val, err := DecodeList[string](value, path, externalRefLocatorContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.Locator().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExternalRefObjectType) Create() SHACLObject {
    return ConstructExternalRefObject(&ExternalRefObject{}, self)
}

func ConstructExternalRefObject(o *ExternalRefObject, typ SHACLType) *ExternalRefObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.comment = NewProperty[string]("comment", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators, RegexValidator[string]{`^[^\/]+\/[^\/]+$`})
        o.contentType = NewProperty[string]("contentType", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altDownloadLocation",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altWebPage",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/binaryArtifact",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/bower",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildMeta",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildSystem",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/certificationReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/chat",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/componentAnalysisReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/cwe",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/documentation",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/dynamicAnalysisReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/eolNotice",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/exportControlAssessment",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/funding",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/issueTracker",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/license",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mailingList",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mavenCentral",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/metrics",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/npm",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/nuget",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/privacyAssessment",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/productMetadata",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/purchaseOrder",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/qualityAssessmentReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseHistory",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseNotes",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/riskAssessment",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/runtimeAnalysisReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/secureSoftwareAttestation",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdversaryModel",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdvisory",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityFix",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityOther",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPenTestReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPolicy",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityThreatModel",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/socialMedia",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/sourceArtifact",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/staticAnalysisReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/support",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vcs",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityDisclosureReport",
                "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityExploitabilityAssessment",
        }})
        o.externalRefType = NewProperty[string]("externalRefType", validators)
    }
    {
        validators := []Validator[string]{}
        o.locator = NewListProperty[string]("locator", validators)
    }
    return o
}

type ExternalRef interface {
    SHACLObject
    Comment() PropertyInterface[string]
    ContentType() PropertyInterface[string]
    ExternalRefType() PropertyInterface[string]
    Locator() ListPropertyInterface[string]
}


func MakeExternalRef() ExternalRef {
    return ConstructExternalRefObject(&ExternalRefObject{}, externalRefType)
}

func MakeExternalRefRef() Ref[ExternalRef] {
    o := MakeExternalRef()
    return MakeObjectRef[ExternalRef](o)
}

func (self *ExternalRefObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("comment")
        if ! self.comment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("contentType")
        if ! self.contentType.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("externalRefType")
        if ! self.externalRefType.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("locator")
        if ! self.locator.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExternalRefObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.comment.Walk(path, visit)
    self.contentType.Walk(path, visit)
    self.externalRefType.Walk(path, visit)
    self.locator.Walk(path, visit)
}


func (self *ExternalRefObject) Comment() PropertyInterface[string] { return &self.comment }
func (self *ExternalRefObject) ContentType() PropertyInterface[string] { return &self.contentType }
func (self *ExternalRefObject) ExternalRefType() PropertyInterface[string] { return &self.externalRefType }
func (self *ExternalRefObject) Locator() ListPropertyInterface[string] { return &self.locator }

func (self *ExternalRefObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.comment.IsSet() {
        val, err := EncodeString(self.comment.Get(), path.PushPath("comment"), externalRefCommentContext, state)
        if err != nil {
            return err
        }
        data["comment"] = val
    }
    if self.contentType.IsSet() {
        val, err := EncodeString(self.contentType.Get(), path.PushPath("contentType"), externalRefContentTypeContext, state)
        if err != nil {
            return err
        }
        data["contentType"] = val
    }
    if self.externalRefType.IsSet() {
        val, err := EncodeIRI(self.externalRefType.Get(), path.PushPath("externalRefType"), externalRefExternalRefTypeContext, state)
        if err != nil {
            return err
        }
        data["externalRefType"] = val
    }
    if self.locator.IsSet() {
        val, err := EncodeList[string](self.locator.Get(), path.PushPath("locator"), externalRefLocatorContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["locator"] = val
    }
    return nil
}

// Specifies the type of an external reference.
type ExternalRefTypeObject struct {
    SHACLObjectBase

}

// A reference to an alternative download location.
const ExternalRefTypeAltDownloadLocation = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altDownloadLocation"
// A reference to an alternative web page.
const ExternalRefTypeAltWebPage = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/altWebPage"
// A reference to binary artifacts related to a package.
const ExternalRefTypeBinaryArtifact = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/binaryArtifact"
// A reference to a Bower package. The package locator format, looks like `package#version`, is defined in the "install" section of [Bower API documentation](https://bower.io/docs/api/#install).
const ExternalRefTypeBower = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/bower"
// A reference build metadata related to a published package.
const ExternalRefTypeBuildMeta = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildMeta"
// A reference build system used to create or publish the package.
const ExternalRefTypeBuildSystem = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/buildSystem"
// A reference to a certification report for a package from an accredited/independent body.
const ExternalRefTypeCertificationReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/certificationReport"
// A reference to the instant messaging system used by the maintainer for a package.
const ExternalRefTypeChat = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/chat"
// A reference to a Software Composition Analysis (SCA) report.
const ExternalRefTypeComponentAnalysisReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/componentAnalysisReport"
// [Common Weakness Enumeration](https://csrc.nist.gov/glossary/term/common_weakness_enumeration). A reference to a source of software flaw defined within the official [CWE List](https://cwe.mitre.org/data/) that conforms to the [CWE specification](https://cwe.mitre.org/).
const ExternalRefTypeCwe = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/cwe"
// A reference to the documentation for a package.
const ExternalRefTypeDocumentation = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/documentation"
// A reference to a dynamic analysis report for a package.
const ExternalRefTypeDynamicAnalysisReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/dynamicAnalysisReport"
// A reference to the End Of Sale (EOS) and/or End Of Life (EOL) information related to a package.
const ExternalRefTypeEolNotice = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/eolNotice"
// A reference to a export control assessment for a package.
const ExternalRefTypeExportControlAssessment = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/exportControlAssessment"
// A reference to funding information related to a package.
const ExternalRefTypeFunding = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/funding"
// A reference to the issue tracker for a package.
const ExternalRefTypeIssueTracker = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/issueTracker"
// A reference to additional license information related to an artifact.
const ExternalRefTypeLicense = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/license"
// A reference to the mailing list used by the maintainer for a package.
const ExternalRefTypeMailingList = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mailingList"
// A reference to a Maven repository artifact. The artifact locator format is defined in the [Maven documentation](https://maven.apache.org/guides/mini/guide-naming-conventions.html) and looks like `groupId:artifactId[:version]`.
const ExternalRefTypeMavenCentral = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/mavenCentral"
// A reference to metrics related to package such as OpenSSF scorecards.
const ExternalRefTypeMetrics = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/metrics"
// A reference to an npm package. The package locator format is defined in the [npm documentation](https://docs.npmjs.com/cli/v10/configuring-npm/package-json) and looks like `package@version`.
const ExternalRefTypeNpm = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/npm"
// A reference to a NuGet package. The package locator format is defined in the [NuGet documentation](https://docs.nuget.org) and looks like `package/version`.
const ExternalRefTypeNuget = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/nuget"
// Used when the type does not match any of the other options.
const ExternalRefTypeOther = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/other"
// A reference to a privacy assessment for a package.
const ExternalRefTypePrivacyAssessment = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/privacyAssessment"
// A reference to additional product metadata such as reference within organization's product catalog.
const ExternalRefTypeProductMetadata = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/productMetadata"
// A reference to a purchase order for a package.
const ExternalRefTypePurchaseOrder = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/purchaseOrder"
// A reference to a quality assessment for a package.
const ExternalRefTypeQualityAssessmentReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/qualityAssessmentReport"
// A reference to a published list of releases for a package.
const ExternalRefTypeReleaseHistory = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseHistory"
// A reference to the release notes for a package.
const ExternalRefTypeReleaseNotes = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/releaseNotes"
// A reference to a risk assessment for a package.
const ExternalRefTypeRiskAssessment = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/riskAssessment"
// A reference to a runtime analysis report for a package.
const ExternalRefTypeRuntimeAnalysisReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/runtimeAnalysisReport"
// A reference to information assuring that the software is developed using security practices as defined by [NIST SP 800-218 Secure Software Development Framework (SSDF) Version 1.1](https://csrc.nist.gov/pubs/sp/800/218/final) or [CISA Secure Software Development Attestation Form](https://www.cisa.gov/resources-tools/resources/secure-software-development-attestation-form).
const ExternalRefTypeSecureSoftwareAttestation = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/secureSoftwareAttestation"
// A reference to the security adversary model for a package.
const ExternalRefTypeSecurityAdversaryModel = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdversaryModel"
// A reference to a published security advisory (where advisory as defined per [ISO 29147:2018](https://www.iso.org/standard/72311.html)) that may affect one or more elements, e.g., vendor advisories or specific NVD entries.
const ExternalRefTypeSecurityAdvisory = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityAdvisory"
// A reference to the patch or source code that fixes a vulnerability.
const ExternalRefTypeSecurityFix = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityFix"
// A reference to related security information of unspecified type.
const ExternalRefTypeSecurityOther = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityOther"
// A reference to a [penetration test](https://en.wikipedia.org/wiki/Penetration_test) report for a package.
const ExternalRefTypeSecurityPenTestReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPenTestReport"
// A reference to instructions for reporting newly discovered security vulnerabilities for a package.
const ExternalRefTypeSecurityPolicy = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityPolicy"
// A reference the [security threat model](https://en.wikipedia.org/wiki/Threat_model) for a package.
const ExternalRefTypeSecurityThreatModel = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/securityThreatModel"
// A reference to a social media channel for a package.
const ExternalRefTypeSocialMedia = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/socialMedia"
// A reference to an artifact containing the sources for a package.
const ExternalRefTypeSourceArtifact = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/sourceArtifact"
// A reference to a static analysis report for a package.
const ExternalRefTypeStaticAnalysisReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/staticAnalysisReport"
// A reference to the software support channel or other support information for a package.
const ExternalRefTypeSupport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/support"
// A reference to a version control system related to a software artifact.
const ExternalRefTypeVcs = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vcs"
// A reference to a Vulnerability Disclosure Report (VDR) which provides the software supplier's analysis and findings describing the impact (or lack of impact) that reported vulnerabilities have on packages or products in the supplier's SBOM as defined in [NIST SP 800-161 Cybersecurity Supply Chain Risk Management Practices for Systems and Organizations](https://csrc.nist.gov/pubs/sp/800/161/r1/final).
const ExternalRefTypeVulnerabilityDisclosureReport = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityDisclosureReport"
// A reference to a Vulnerability Exploitability eXchange (VEX) statement which provides information on whether a product is impacted by a specific vulnerability in an included package and, if affected, whether there are actions recommended to remediate. See also [NTIA VEX one-page summary](https://ntia.gov/files/ntia/publications/vex_one-page_summary.pdf).
const ExternalRefTypeVulnerabilityExploitabilityAssessment = "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType/vulnerabilityExploitabilityAssessment"

type ExternalRefTypeObjectType struct {
    SHACLTypeBase
}
var externalRefTypeType ExternalRefTypeObjectType

func DecodeExternalRefType (data any, path Path, context map[string]string) (Ref[ExternalRefType], error) {
    return DecodeRef[ExternalRefType](data, path, context, externalRefTypeType)
}

func (self ExternalRefTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExternalRefType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExternalRefTypeObjectType) Create() SHACLObject {
    return ConstructExternalRefTypeObject(&ExternalRefTypeObject{}, self)
}

func ConstructExternalRefTypeObject(o *ExternalRefTypeObject, typ SHACLType) *ExternalRefTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type ExternalRefType interface {
    SHACLObject
}


func MakeExternalRefType() ExternalRefType {
    return ConstructExternalRefTypeObject(&ExternalRefTypeObject{}, externalRefTypeType)
}

func MakeExternalRefTypeRef() Ref[ExternalRefType] {
    o := MakeExternalRefType()
    return MakeObjectRef[ExternalRefType](o)
}

func (self *ExternalRefTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExternalRefTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *ExternalRefTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A mathematical algorithm that maps data of arbitrary size to a bit string.
type HashAlgorithmObject struct {
    SHACLObjectBase

}

// Adler-32 checksum is part of the widely used zlib compression library as defined in [RFC 1950](https://datatracker.ietf.org/doc/rfc1950/) Section 2.3.
const HashAlgorithmAdler32 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/adler32"
// BLAKE2b algorithm with a digest size of 256, as defined in [RFC 7693](https://datatracker.ietf.org/doc/rfc7693/) Section 4.
const HashAlgorithmBlake2b256 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b256"
// BLAKE2b algorithm with a digest size of 384, as defined in [RFC 7693](https://datatracker.ietf.org/doc/rfc7693/) Section 4.
const HashAlgorithmBlake2b384 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b384"
// BLAKE2b algorithm with a digest size of 512, as defined in [RFC 7693](https://datatracker.ietf.org/doc/rfc7693/) Section 4.
const HashAlgorithmBlake2b512 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b512"
// [BLAKE3](https://github.com/BLAKE3-team/BLAKE3-specs/blob/master/blake3.pdf)
const HashAlgorithmBlake3 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake3"
// [Dilithium](https://pq-crystals.org/dilithium/)
const HashAlgorithmCrystalsDilithium = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsDilithium"
// [Kyber](https://pq-crystals.org/kyber/)
const HashAlgorithmCrystalsKyber = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsKyber"
// [FALCON](https://falcon-sign.info/falcon.pdf)
const HashAlgorithmFalcon = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/falcon"
// MD2 message-digest algorithm, as defined in [RFC 1319](https://datatracker.ietf.org/doc/rfc1319/).
const HashAlgorithmMd2 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md2"
// MD4 message-digest algorithm, as defined in [RFC 1186](https://datatracker.ietf.org/doc/rfc1186/).
const HashAlgorithmMd4 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md4"
// MD5 message-digest algorithm, as defined in [RFC 1321](https://datatracker.ietf.org/doc/rfc1321/).
const HashAlgorithmMd5 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md5"
// [MD6 hash function](https://people.csail.mit.edu/rivest/pubs/RABCx08.pdf)
const HashAlgorithmMd6 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md6"
// any hashing algorithm that does not exist in this list of entries
const HashAlgorithmOther = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/other"
// SHA-1, a secure hashing algorithm, as defined in [RFC 3174](https://datatracker.ietf.org/doc/rfc3174/).
const HashAlgorithmSha1 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha1"
// SHA-2 with a digest length of 224, as defined in [RFC 3874](https://datatracker.ietf.org/doc/rfc3874/).
const HashAlgorithmSha224 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha224"
// SHA-2 with a digest length of 256, as defined in [RFC 6234](https://datatracker.ietf.org/doc/rfc6234/).
const HashAlgorithmSha256 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha256"
// SHA-2 with a digest length of 384, as defined in [RFC 6234](https://datatracker.ietf.org/doc/rfc6234/).
const HashAlgorithmSha384 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha384"
// SHA-3 with a digest length of 224, as defined in [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final).
const HashAlgorithmSha3224 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_224"
// SHA-3 with a digest length of 256, as defined in [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final).
const HashAlgorithmSha3256 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_256"
// SHA-3 with a digest length of 384, as defined in [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final).
const HashAlgorithmSha3384 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_384"
// SHA-3 with a digest length of 512, as defined in [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final).
const HashAlgorithmSha3512 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_512"
// SHA-2 with a digest length of 512, as defined in [RFC 6234](https://datatracker.ietf.org/doc/rfc6234/).
const HashAlgorithmSha512 = "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha512"

type HashAlgorithmObjectType struct {
    SHACLTypeBase
}
var hashAlgorithmType HashAlgorithmObjectType

func DecodeHashAlgorithm (data any, path Path, context map[string]string) (Ref[HashAlgorithm], error) {
    return DecodeRef[HashAlgorithm](data, path, context, hashAlgorithmType)
}

func (self HashAlgorithmObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(HashAlgorithm)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self HashAlgorithmObjectType) Create() SHACLObject {
    return ConstructHashAlgorithmObject(&HashAlgorithmObject{}, self)
}

func ConstructHashAlgorithmObject(o *HashAlgorithmObject, typ SHACLType) *HashAlgorithmObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type HashAlgorithm interface {
    SHACLObject
}


func MakeHashAlgorithm() HashAlgorithm {
    return ConstructHashAlgorithmObject(&HashAlgorithmObject{}, hashAlgorithmType)
}

func MakeHashAlgorithmRef() Ref[HashAlgorithm] {
    o := MakeHashAlgorithm()
    return MakeObjectRef[HashAlgorithm](o)
}

func (self *HashAlgorithmObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *HashAlgorithmObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *HashAlgorithmObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A concrete subclass of Element used by Individuals in the

// Core profile.
type IndividualElementObject struct {
    ElementObject

}

// An Individual Value for Element representing a set of Elements of unknown

// identify or cardinality (number).
const IndividualElementNoAssertionElement = "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement"
// An Individual Value for Element representing a set of Elements with

// cardinality (number/count) of zero.
const IndividualElementNoneElement = "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement"

type IndividualElementObjectType struct {
    SHACLTypeBase
}
var individualElementType IndividualElementObjectType

func DecodeIndividualElement (data any, path Path, context map[string]string) (Ref[IndividualElement], error) {
    return DecodeRef[IndividualElement](data, path, context, individualElementType)
}

func (self IndividualElementObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(IndividualElement)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self IndividualElementObjectType) Create() SHACLObject {
    return ConstructIndividualElementObject(&IndividualElementObject{}, self)
}

func ConstructIndividualElementObject(o *IndividualElementObject, typ SHACLType) *IndividualElementObject {
    ConstructElementObject(&o.ElementObject, typ)
    return o
}

type IndividualElement interface {
    Element
}


func MakeIndividualElement() IndividualElement {
    return ConstructIndividualElementObject(&IndividualElementObject{}, individualElementType)
}

func MakeIndividualElementRef() Ref[IndividualElement] {
    o := MakeIndividualElement()
    return MakeObjectRef[IndividualElement](o)
}

func (self *IndividualElementObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *IndividualElementObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
}



func (self *IndividualElementObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Provides an independently reproducible mechanism that permits verification of a specific Element.
type IntegrityMethodObject struct {
    SHACLObjectBase

    // Provide consumers with comments by the creator of the Element about the
    // Element.
    comment Property[string]
}


type IntegrityMethodObjectType struct {
    SHACLTypeBase
}
var integrityMethodType IntegrityMethodObjectType
var integrityMethodCommentContext = map[string]string{}

func DecodeIntegrityMethod (data any, path Path, context map[string]string) (Ref[IntegrityMethod], error) {
    return DecodeRef[IntegrityMethod](data, path, context, integrityMethodType)
}

func (self IntegrityMethodObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(IntegrityMethod)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/comment", "comment":
        val, err := DecodeString(value, path, integrityMethodCommentContext)
        if err != nil {
            return false, err
        }
        err = obj.Comment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self IntegrityMethodObjectType) Create() SHACLObject {
    return ConstructIntegrityMethodObject(&IntegrityMethodObject{}, self)
}

func ConstructIntegrityMethodObject(o *IntegrityMethodObject, typ SHACLType) *IntegrityMethodObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.comment = NewProperty[string]("comment", validators)
    }
    return o
}

type IntegrityMethod interface {
    SHACLObject
    Comment() PropertyInterface[string]
}



func (self *IntegrityMethodObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("comment")
        if ! self.comment.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *IntegrityMethodObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.comment.Walk(path, visit)
}


func (self *IntegrityMethodObject) Comment() PropertyInterface[string] { return &self.comment }

func (self *IntegrityMethodObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.comment.IsSet() {
        val, err := EncodeString(self.comment.Get(), path.PushPath("comment"), integrityMethodCommentContext, state)
        if err != nil {
            return err
        }
        data["comment"] = val
    }
    return nil
}

// Provide an enumerated set of lifecycle phases that can provide context to relationships.
type LifecycleScopeTypeObject struct {
    SHACLObjectBase

}

// A relationship has specific context implications during an element's build phase, during development.
const LifecycleScopeTypeBuild = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/build"
// A relationship has specific context implications during an element's design.
const LifecycleScopeTypeDesign = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/design"
// A relationship has specific context implications during development phase of an element.
const LifecycleScopeTypeDevelopment = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/development"
// A relationship has other specific context information necessary to capture that the above set of enumerations does not handle.
const LifecycleScopeTypeOther = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/other"
// A relationship has specific context implications during the execution phase of an element.
const LifecycleScopeTypeRuntime = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/runtime"
// A relationship has specific context implications during an element's testing phase, during development.
const LifecycleScopeTypeTest = "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/test"

type LifecycleScopeTypeObjectType struct {
    SHACLTypeBase
}
var lifecycleScopeTypeType LifecycleScopeTypeObjectType

func DecodeLifecycleScopeType (data any, path Path, context map[string]string) (Ref[LifecycleScopeType], error) {
    return DecodeRef[LifecycleScopeType](data, path, context, lifecycleScopeTypeType)
}

func (self LifecycleScopeTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(LifecycleScopeType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self LifecycleScopeTypeObjectType) Create() SHACLObject {
    return ConstructLifecycleScopeTypeObject(&LifecycleScopeTypeObject{}, self)
}

func ConstructLifecycleScopeTypeObject(o *LifecycleScopeTypeObject, typ SHACLType) *LifecycleScopeTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type LifecycleScopeType interface {
    SHACLObject
}


func MakeLifecycleScopeType() LifecycleScopeType {
    return ConstructLifecycleScopeTypeObject(&LifecycleScopeTypeObject{}, lifecycleScopeTypeType)
}

func MakeLifecycleScopeTypeRef() Ref[LifecycleScopeType] {
    o := MakeLifecycleScopeType()
    return MakeObjectRef[LifecycleScopeType](o)
}

func (self *LifecycleScopeTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *LifecycleScopeTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *LifecycleScopeTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A mapping between prefixes and namespace partial URIs.
type NamespaceMapObject struct {
    SHACLObjectBase

    // Provides an unambiguous mechanism for conveying a URI fragment portion of an
    // Element ID.
    namespace Property[string]
    // A substitute for a URI.
    prefix Property[string]
}


type NamespaceMapObjectType struct {
    SHACLTypeBase
}
var namespaceMapType NamespaceMapObjectType
var namespaceMapNamespaceContext = map[string]string{}
var namespaceMapPrefixContext = map[string]string{}

func DecodeNamespaceMap (data any, path Path, context map[string]string) (Ref[NamespaceMap], error) {
    return DecodeRef[NamespaceMap](data, path, context, namespaceMapType)
}

func (self NamespaceMapObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(NamespaceMap)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/namespace", "namespace":
        val, err := DecodeString(value, path, namespaceMapNamespaceContext)
        if err != nil {
            return false, err
        }
        err = obj.Namespace().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/prefix", "prefix":
        val, err := DecodeString(value, path, namespaceMapPrefixContext)
        if err != nil {
            return false, err
        }
        err = obj.Prefix().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self NamespaceMapObjectType) Create() SHACLObject {
    return ConstructNamespaceMapObject(&NamespaceMapObject{}, self)
}

func ConstructNamespaceMapObject(o *NamespaceMapObject, typ SHACLType) *NamespaceMapObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.namespace = NewProperty[string]("namespace", validators)
    }
    {
        validators := []Validator[string]{}
        o.prefix = NewProperty[string]("prefix", validators)
    }
    return o
}

type NamespaceMap interface {
    SHACLObject
    Namespace() PropertyInterface[string]
    Prefix() PropertyInterface[string]
}


func MakeNamespaceMap() NamespaceMap {
    return ConstructNamespaceMapObject(&NamespaceMapObject{}, namespaceMapType)
}

func MakeNamespaceMapRef() Ref[NamespaceMap] {
    o := MakeNamespaceMap()
    return MakeObjectRef[NamespaceMap](o)
}

func (self *NamespaceMapObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("namespace")
        if ! self.namespace.Check(prop_path, handler) {
            valid = false
        }
        if ! self.namespace.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"namespace", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("prefix")
        if ! self.prefix.Check(prop_path, handler) {
            valid = false
        }
        if ! self.prefix.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"prefix", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *NamespaceMapObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.namespace.Walk(path, visit)
    self.prefix.Walk(path, visit)
}


func (self *NamespaceMapObject) Namespace() PropertyInterface[string] { return &self.namespace }
func (self *NamespaceMapObject) Prefix() PropertyInterface[string] { return &self.prefix }

func (self *NamespaceMapObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.namespace.IsSet() {
        val, err := EncodeString(self.namespace.Get(), path.PushPath("namespace"), namespaceMapNamespaceContext, state)
        if err != nil {
            return err
        }
        data["namespace"] = val
    }
    if self.prefix.IsSet() {
        val, err := EncodeString(self.prefix.Get(), path.PushPath("prefix"), namespaceMapPrefixContext, state)
        if err != nil {
            return err
        }
        data["prefix"] = val
    }
    return nil
}

// An SPDX version 2.X compatible verification method for software packages.
type PackageVerificationCodeObject struct {
    IntegrityMethodObject

    // Specifies the algorithm used for calculating the hash value.
    algorithm Property[string]
    // The result of applying a hash algorithm to an Element.
    hashValue Property[string]
    // The relative file name of a file to be excluded from the
    // `PackageVerificationCode`.
    packageVerificationCodeExcludedFile ListProperty[string]
}


type PackageVerificationCodeObjectType struct {
    SHACLTypeBase
}
var packageVerificationCodeType PackageVerificationCodeObjectType
var packageVerificationCodeAlgorithmContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/adler32": "adler32",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b256": "blake2b256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b384": "blake2b384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b512": "blake2b512",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake3": "blake3",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsDilithium": "crystalsDilithium",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsKyber": "crystalsKyber",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/falcon": "falcon",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md2": "md2",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md4": "md4",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md5": "md5",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md6": "md6",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha1": "sha1",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha224": "sha224",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha256": "sha256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha384": "sha384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_224": "sha3_224",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_256": "sha3_256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_384": "sha3_384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_512": "sha3_512",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha512": "sha512",}
var packageVerificationCodeHashValueContext = map[string]string{}
var packageVerificationCodePackageVerificationCodeExcludedFileContext = map[string]string{}

func DecodePackageVerificationCode (data any, path Path, context map[string]string) (Ref[PackageVerificationCode], error) {
    return DecodeRef[PackageVerificationCode](data, path, context, packageVerificationCodeType)
}

func (self PackageVerificationCodeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(PackageVerificationCode)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/algorithm", "algorithm":
        val, err := DecodeIRI(value, path, packageVerificationCodeAlgorithmContext)
        if err != nil {
            return false, err
        }
        err = obj.Algorithm().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/hashValue", "hashValue":
        val, err := DecodeString(value, path, packageVerificationCodeHashValueContext)
        if err != nil {
            return false, err
        }
        err = obj.HashValue().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/packageVerificationCodeExcludedFile", "packageVerificationCodeExcludedFile":
        val, err := DecodeList[string](value, path, packageVerificationCodePackageVerificationCodeExcludedFileContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.PackageVerificationCodeExcludedFile().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self PackageVerificationCodeObjectType) Create() SHACLObject {
    return ConstructPackageVerificationCodeObject(&PackageVerificationCodeObject{}, self)
}

func ConstructPackageVerificationCodeObject(o *PackageVerificationCodeObject, typ SHACLType) *PackageVerificationCodeObject {
    ConstructIntegrityMethodObject(&o.IntegrityMethodObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/adler32",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b512",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake3",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsDilithium",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsKyber",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/falcon",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md2",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md4",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md5",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md6",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha1",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha224",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_224",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_512",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha512",
        }})
        o.algorithm = NewProperty[string]("algorithm", validators)
    }
    {
        validators := []Validator[string]{}
        o.hashValue = NewProperty[string]("hashValue", validators)
    }
    {
        validators := []Validator[string]{}
        o.packageVerificationCodeExcludedFile = NewListProperty[string]("packageVerificationCodeExcludedFile", validators)
    }
    return o
}

type PackageVerificationCode interface {
    IntegrityMethod
    Algorithm() PropertyInterface[string]
    HashValue() PropertyInterface[string]
    PackageVerificationCodeExcludedFile() ListPropertyInterface[string]
}


func MakePackageVerificationCode() PackageVerificationCode {
    return ConstructPackageVerificationCodeObject(&PackageVerificationCodeObject{}, packageVerificationCodeType)
}

func MakePackageVerificationCodeRef() Ref[PackageVerificationCode] {
    o := MakePackageVerificationCode()
    return MakeObjectRef[PackageVerificationCode](o)
}

func (self *PackageVerificationCodeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.IntegrityMethodObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("algorithm")
        if ! self.algorithm.Check(prop_path, handler) {
            valid = false
        }
        if ! self.algorithm.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"algorithm", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("hashValue")
        if ! self.hashValue.Check(prop_path, handler) {
            valid = false
        }
        if ! self.hashValue.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"hashValue", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("packageVerificationCodeExcludedFile")
        if ! self.packageVerificationCodeExcludedFile.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *PackageVerificationCodeObject) Walk(path Path, visit Visit) {
    self.IntegrityMethodObject.Walk(path, visit)
    self.algorithm.Walk(path, visit)
    self.hashValue.Walk(path, visit)
    self.packageVerificationCodeExcludedFile.Walk(path, visit)
}


func (self *PackageVerificationCodeObject) Algorithm() PropertyInterface[string] { return &self.algorithm }
func (self *PackageVerificationCodeObject) HashValue() PropertyInterface[string] { return &self.hashValue }
func (self *PackageVerificationCodeObject) PackageVerificationCodeExcludedFile() ListPropertyInterface[string] { return &self.packageVerificationCodeExcludedFile }

func (self *PackageVerificationCodeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.IntegrityMethodObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.algorithm.IsSet() {
        val, err := EncodeIRI(self.algorithm.Get(), path.PushPath("algorithm"), packageVerificationCodeAlgorithmContext, state)
        if err != nil {
            return err
        }
        data["algorithm"] = val
    }
    if self.hashValue.IsSet() {
        val, err := EncodeString(self.hashValue.Get(), path.PushPath("hashValue"), packageVerificationCodeHashValueContext, state)
        if err != nil {
            return err
        }
        data["hashValue"] = val
    }
    if self.packageVerificationCodeExcludedFile.IsSet() {
        val, err := EncodeList[string](self.packageVerificationCodeExcludedFile.Get(), path.PushPath("packageVerificationCodeExcludedFile"), packageVerificationCodePackageVerificationCodeExcludedFileContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["packageVerificationCodeExcludedFile"] = val
    }
    return nil
}

// A tuple of two positive integers that define a range.
type PositiveIntegerRangeObject struct {
    SHACLObjectBase

    // Defines the beginning of a range.
    beginIntegerRange Property[int]
    // Defines the end of a range.
    endIntegerRange Property[int]
}


type PositiveIntegerRangeObjectType struct {
    SHACLTypeBase
}
var positiveIntegerRangeType PositiveIntegerRangeObjectType
var positiveIntegerRangeBeginIntegerRangeContext = map[string]string{}
var positiveIntegerRangeEndIntegerRangeContext = map[string]string{}

func DecodePositiveIntegerRange (data any, path Path, context map[string]string) (Ref[PositiveIntegerRange], error) {
    return DecodeRef[PositiveIntegerRange](data, path, context, positiveIntegerRangeType)
}

func (self PositiveIntegerRangeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(PositiveIntegerRange)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/beginIntegerRange", "beginIntegerRange":
        val, err := DecodeInteger(value, path, positiveIntegerRangeBeginIntegerRangeContext)
        if err != nil {
            return false, err
        }
        err = obj.BeginIntegerRange().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/endIntegerRange", "endIntegerRange":
        val, err := DecodeInteger(value, path, positiveIntegerRangeEndIntegerRangeContext)
        if err != nil {
            return false, err
        }
        err = obj.EndIntegerRange().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self PositiveIntegerRangeObjectType) Create() SHACLObject {
    return ConstructPositiveIntegerRangeObject(&PositiveIntegerRangeObject{}, self)
}

func ConstructPositiveIntegerRangeObject(o *PositiveIntegerRangeObject, typ SHACLType) *PositiveIntegerRangeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[int]{}
        validators = append(validators, IntegerMinValidator{1})
        o.beginIntegerRange = NewProperty[int]("beginIntegerRange", validators)
    }
    {
        validators := []Validator[int]{}
        validators = append(validators, IntegerMinValidator{1})
        o.endIntegerRange = NewProperty[int]("endIntegerRange", validators)
    }
    return o
}

type PositiveIntegerRange interface {
    SHACLObject
    BeginIntegerRange() PropertyInterface[int]
    EndIntegerRange() PropertyInterface[int]
}


func MakePositiveIntegerRange() PositiveIntegerRange {
    return ConstructPositiveIntegerRangeObject(&PositiveIntegerRangeObject{}, positiveIntegerRangeType)
}

func MakePositiveIntegerRangeRef() Ref[PositiveIntegerRange] {
    o := MakePositiveIntegerRange()
    return MakeObjectRef[PositiveIntegerRange](o)
}

func (self *PositiveIntegerRangeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("beginIntegerRange")
        if ! self.beginIntegerRange.Check(prop_path, handler) {
            valid = false
        }
        if ! self.beginIntegerRange.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"beginIntegerRange", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("endIntegerRange")
        if ! self.endIntegerRange.Check(prop_path, handler) {
            valid = false
        }
        if ! self.endIntegerRange.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"endIntegerRange", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *PositiveIntegerRangeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.beginIntegerRange.Walk(path, visit)
    self.endIntegerRange.Walk(path, visit)
}


func (self *PositiveIntegerRangeObject) BeginIntegerRange() PropertyInterface[int] { return &self.beginIntegerRange }
func (self *PositiveIntegerRangeObject) EndIntegerRange() PropertyInterface[int] { return &self.endIntegerRange }

func (self *PositiveIntegerRangeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.beginIntegerRange.IsSet() {
        val, err := EncodeInteger(self.beginIntegerRange.Get(), path.PushPath("beginIntegerRange"), positiveIntegerRangeBeginIntegerRangeContext, state)
        if err != nil {
            return err
        }
        data["beginIntegerRange"] = val
    }
    if self.endIntegerRange.IsSet() {
        val, err := EncodeInteger(self.endIntegerRange.Get(), path.PushPath("endIntegerRange"), positiveIntegerRangeEndIntegerRangeContext, state)
        if err != nil {
            return err
        }
        data["endIntegerRange"] = val
    }
    return nil
}

// Categories of presence or absence.
type PresenceTypeObject struct {
    SHACLObjectBase

}

// Indicates absence of the field.
const PresenceTypeNo = "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no"
// Makes no assertion about the field.
const PresenceTypeNoAssertion = "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion"
// Indicates presence of the field.
const PresenceTypeYes = "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes"

type PresenceTypeObjectType struct {
    SHACLTypeBase
}
var presenceTypeType PresenceTypeObjectType

func DecodePresenceType (data any, path Path, context map[string]string) (Ref[PresenceType], error) {
    return DecodeRef[PresenceType](data, path, context, presenceTypeType)
}

func (self PresenceTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(PresenceType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self PresenceTypeObjectType) Create() SHACLObject {
    return ConstructPresenceTypeObject(&PresenceTypeObject{}, self)
}

func ConstructPresenceTypeObject(o *PresenceTypeObject, typ SHACLType) *PresenceTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type PresenceType interface {
    SHACLObject
}


func MakePresenceType() PresenceType {
    return ConstructPresenceTypeObject(&PresenceTypeObject{}, presenceTypeType)
}

func MakePresenceTypeRef() Ref[PresenceType] {
    o := MakePresenceType()
    return MakeObjectRef[PresenceType](o)
}

func (self *PresenceTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *PresenceTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *PresenceTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Enumeration of the valid profiles.
type ProfileIdentifierTypeObject struct {
    SHACLObjectBase

}

// the element follows the AI profile specification
const ProfileIdentifierTypeAi = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/ai"
// the element follows the Build profile specification
const ProfileIdentifierTypeBuild = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/build"
// the element follows the Core profile specification
const ProfileIdentifierTypeCore = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/core"
// the element follows the Dataset profile specification
const ProfileIdentifierTypeDataset = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/dataset"
// the element follows the ExpandedLicensing profile specification
const ProfileIdentifierTypeExpandedLicensing = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/expandedLicensing"
// the element follows the Extension profile specification
const ProfileIdentifierTypeExtension = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/extension"
// the element follows the Lite profile specification
const ProfileIdentifierTypeLite = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/lite"
// the element follows the Security profile specification
const ProfileIdentifierTypeSecurity = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/security"
// the element follows the SimpleLicensing profile specification
const ProfileIdentifierTypeSimpleLicensing = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/simpleLicensing"
// the element follows the Software profile specification
const ProfileIdentifierTypeSoftware = "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType/software"

type ProfileIdentifierTypeObjectType struct {
    SHACLTypeBase
}
var profileIdentifierTypeType ProfileIdentifierTypeObjectType

func DecodeProfileIdentifierType (data any, path Path, context map[string]string) (Ref[ProfileIdentifierType], error) {
    return DecodeRef[ProfileIdentifierType](data, path, context, profileIdentifierTypeType)
}

func (self ProfileIdentifierTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ProfileIdentifierType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ProfileIdentifierTypeObjectType) Create() SHACLObject {
    return ConstructProfileIdentifierTypeObject(&ProfileIdentifierTypeObject{}, self)
}

func ConstructProfileIdentifierTypeObject(o *ProfileIdentifierTypeObject, typ SHACLType) *ProfileIdentifierTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type ProfileIdentifierType interface {
    SHACLObject
}


func MakeProfileIdentifierType() ProfileIdentifierType {
    return ConstructProfileIdentifierTypeObject(&ProfileIdentifierTypeObject{}, profileIdentifierTypeType)
}

func MakeProfileIdentifierTypeRef() Ref[ProfileIdentifierType] {
    o := MakeProfileIdentifierType()
    return MakeObjectRef[ProfileIdentifierType](o)
}

func (self *ProfileIdentifierTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ProfileIdentifierTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *ProfileIdentifierTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Describes a relationship between one or more elements.
type RelationshipObject struct {
    ElementObject

    // Provides information about the completeness of relationships.
    completeness Property[string]
    // Specifies the time from which an element is no longer applicable / valid.
    endTime Property[time.Time]
    // References the Element on the left-hand side of a relationship.
    from RefProperty[Element]
    // Information about the relationship between two Elements.
    relationshipType Property[string]
    // Specifies the time from which an element is applicable / valid.
    startTime Property[time.Time]
    // References an Element on the right-hand side of a relationship.
    to RefListProperty[Element]
}


type RelationshipObjectType struct {
    SHACLTypeBase
}
var relationshipType RelationshipObjectType
var relationshipCompletenessContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/complete": "complete",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/incomplete": "incomplete",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/noAssertion": "noAssertion",}
var relationshipEndTimeContext = map[string]string{}
var relationshipFromContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement": "NoAssertionElement",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement": "NoneElement",
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}
var relationshipRelationshipTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/affects": "affects",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/amendedBy": "amendedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/ancestorOf": "ancestorOf",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/availableFrom": "availableFrom",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/configures": "configures",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/contains": "contains",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/coordinatedBy": "coordinatedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/copiedTo": "copiedTo",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/delegatedTo": "delegatedTo",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/dependsOn": "dependsOn",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/descendantOf": "descendantOf",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/describes": "describes",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/doesNotAffect": "doesNotAffect",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/expandsTo": "expandsTo",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/exploitCreatedBy": "exploitCreatedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedBy": "fixedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedIn": "fixedIn",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/foundBy": "foundBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/generates": "generates",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAddedFile": "hasAddedFile",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssessmentFor": "hasAssessmentFor",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssociatedVulnerability": "hasAssociatedVulnerability",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasConcludedLicense": "hasConcludedLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDataFile": "hasDataFile",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeclaredLicense": "hasDeclaredLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeletedFile": "hasDeletedFile",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDependencyManifest": "hasDependencyManifest",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDistributionArtifact": "hasDistributionArtifact",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDocumentation": "hasDocumentation",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDynamicLink": "hasDynamicLink",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasEvidence": "hasEvidence",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasExample": "hasExample",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasHost": "hasHost",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasInput": "hasInput",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasMetadata": "hasMetadata",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalComponent": "hasOptionalComponent",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalDependency": "hasOptionalDependency",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOutput": "hasOutput",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasPrerequisite": "hasPrerequisite",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasProvidedDependency": "hasProvidedDependency",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasRequirement": "hasRequirement",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasSpecification": "hasSpecification",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasStaticLink": "hasStaticLink",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTest": "hasTest",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTestCase": "hasTestCase",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasVariant": "hasVariant",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/invokedBy": "invokedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/modifiedBy": "modifiedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/packagedBy": "packagedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/patchedBy": "patchedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/publishedBy": "publishedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/reportedBy": "reportedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/republishedBy": "republishedBy",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/serializedInArtifact": "serializedInArtifact",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/testedOn": "testedOn",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/trainedOn": "trainedOn",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/underInvestigationFor": "underInvestigationFor",
    "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/usesTool": "usesTool",}
var relationshipStartTimeContext = map[string]string{}
var relationshipToContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement": "NoAssertionElement",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement": "NoneElement",
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}

func DecodeRelationship (data any, path Path, context map[string]string) (Ref[Relationship], error) {
    return DecodeRef[Relationship](data, path, context, relationshipType)
}

func (self RelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Relationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/completeness", "completeness":
        val, err := DecodeIRI(value, path, relationshipCompletenessContext)
        if err != nil {
            return false, err
        }
        err = obj.Completeness().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/endTime", "endTime":
        val, err := DecodeDateTimeStamp(value, path, relationshipEndTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.EndTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/from", "from":
        val, err := DecodeElement(value, path, relationshipFromContext)
        if err != nil {
            return false, err
        }
        err = obj.From().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/relationshipType", "relationshipType":
        val, err := DecodeIRI(value, path, relationshipRelationshipTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.RelationshipType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/startTime", "startTime":
        val, err := DecodeDateTimeStamp(value, path, relationshipStartTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.StartTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/to", "to":
        val, err := DecodeList[Ref[Element]](value, path, relationshipToContext, DecodeElement)
        if err != nil {
            return false, err
        }
        err = obj.To().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self RelationshipObjectType) Create() SHACLObject {
    return ConstructRelationshipObject(&RelationshipObject{}, self)
}

func ConstructRelationshipObject(o *RelationshipObject, typ SHACLType) *RelationshipObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/complete",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/incomplete",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/noAssertion",
        }})
        o.completeness = NewProperty[string]("completeness", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.endTime = NewProperty[time.Time]("endTime", validators)
    }
    {
        validators := []Validator[Ref[Element]]{}
        o.from = NewRefProperty[Element]("from", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/affects",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/amendedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/ancestorOf",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/availableFrom",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/configures",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/contains",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/coordinatedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/copiedTo",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/delegatedTo",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/dependsOn",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/descendantOf",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/describes",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/doesNotAffect",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/expandsTo",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/exploitCreatedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedIn",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/foundBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/generates",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAddedFile",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssessmentFor",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssociatedVulnerability",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasConcludedLicense",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDataFile",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeclaredLicense",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeletedFile",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDependencyManifest",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDistributionArtifact",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDocumentation",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDynamicLink",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasEvidence",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasExample",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasHost",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasInput",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasMetadata",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalComponent",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalDependency",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOutput",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasPrerequisite",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasProvidedDependency",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasRequirement",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasSpecification",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasStaticLink",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTest",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTestCase",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasVariant",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/invokedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/modifiedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/packagedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/patchedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/publishedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/reportedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/republishedBy",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/serializedInArtifact",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/testedOn",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/trainedOn",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/underInvestigationFor",
                "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/usesTool",
        }})
        o.relationshipType = NewProperty[string]("relationshipType", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.startTime = NewProperty[time.Time]("startTime", validators)
    }
    {
        validators := []Validator[Ref[Element]]{}
        o.to = NewRefListProperty[Element]("to", validators)
    }
    return o
}

type Relationship interface {
    Element
    Completeness() PropertyInterface[string]
    EndTime() PropertyInterface[time.Time]
    From() RefPropertyInterface[Element]
    RelationshipType() PropertyInterface[string]
    StartTime() PropertyInterface[time.Time]
    To() ListPropertyInterface[Ref[Element]]
}


func MakeRelationship() Relationship {
    return ConstructRelationshipObject(&RelationshipObject{}, relationshipType)
}

func MakeRelationshipRef() Ref[Relationship] {
    o := MakeRelationship()
    return MakeObjectRef[Relationship](o)
}

func (self *RelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("completeness")
        if ! self.completeness.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("endTime")
        if ! self.endTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("from")
        if ! self.from.Check(prop_path, handler) {
            valid = false
        }
        if ! self.from.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"from", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("relationshipType")
        if ! self.relationshipType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.relationshipType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"relationshipType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("startTime")
        if ! self.startTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("to")
        if ! self.to.Check(prop_path, handler) {
            valid = false
        }
        if len(self.to.Get()) < 1 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "to",
                    "Too few elements. Minimum of 1 required"},
                    prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *RelationshipObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.completeness.Walk(path, visit)
    self.endTime.Walk(path, visit)
    self.from.Walk(path, visit)
    self.relationshipType.Walk(path, visit)
    self.startTime.Walk(path, visit)
    self.to.Walk(path, visit)
}


func (self *RelationshipObject) Completeness() PropertyInterface[string] { return &self.completeness }
func (self *RelationshipObject) EndTime() PropertyInterface[time.Time] { return &self.endTime }
func (self *RelationshipObject) From() RefPropertyInterface[Element] { return &self.from }
func (self *RelationshipObject) RelationshipType() PropertyInterface[string] { return &self.relationshipType }
func (self *RelationshipObject) StartTime() PropertyInterface[time.Time] { return &self.startTime }
func (self *RelationshipObject) To() ListPropertyInterface[Ref[Element]] { return &self.to }

func (self *RelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.completeness.IsSet() {
        val, err := EncodeIRI(self.completeness.Get(), path.PushPath("completeness"), relationshipCompletenessContext, state)
        if err != nil {
            return err
        }
        data["completeness"] = val
    }
    if self.endTime.IsSet() {
        val, err := EncodeDateTime(self.endTime.Get(), path.PushPath("endTime"), relationshipEndTimeContext, state)
        if err != nil {
            return err
        }
        data["endTime"] = val
    }
    if self.from.IsSet() {
        val, err := EncodeRef[Element](self.from.Get(), path.PushPath("from"), relationshipFromContext, state)
        if err != nil {
            return err
        }
        data["from"] = val
    }
    if self.relationshipType.IsSet() {
        val, err := EncodeIRI(self.relationshipType.Get(), path.PushPath("relationshipType"), relationshipRelationshipTypeContext, state)
        if err != nil {
            return err
        }
        data["relationshipType"] = val
    }
    if self.startTime.IsSet() {
        val, err := EncodeDateTime(self.startTime.Get(), path.PushPath("startTime"), relationshipStartTimeContext, state)
        if err != nil {
            return err
        }
        data["startTime"] = val
    }
    if self.to.IsSet() {
        val, err := EncodeList[Ref[Element]](self.to.Get(), path.PushPath("to"), relationshipToContext, state, EncodeRef[Element])
        if err != nil {
            return err
        }
        data["to"] = val
    }
    return nil
}

// Indicates whether a relationship is known to be complete, incomplete, or if no assertion is made with respect to relationship completeness.
type RelationshipCompletenessObject struct {
    SHACLObjectBase

}

// The relationship is known to be exhaustive.
const RelationshipCompletenessComplete = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/complete"
// The relationship is known not to be exhaustive.
const RelationshipCompletenessIncomplete = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/incomplete"
// No assertion can be made about the completeness of the relationship.
const RelationshipCompletenessNoAssertion = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness/noAssertion"

type RelationshipCompletenessObjectType struct {
    SHACLTypeBase
}
var relationshipCompletenessType RelationshipCompletenessObjectType

func DecodeRelationshipCompleteness (data any, path Path, context map[string]string) (Ref[RelationshipCompleteness], error) {
    return DecodeRef[RelationshipCompleteness](data, path, context, relationshipCompletenessType)
}

func (self RelationshipCompletenessObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(RelationshipCompleteness)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self RelationshipCompletenessObjectType) Create() SHACLObject {
    return ConstructRelationshipCompletenessObject(&RelationshipCompletenessObject{}, self)
}

func ConstructRelationshipCompletenessObject(o *RelationshipCompletenessObject, typ SHACLType) *RelationshipCompletenessObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type RelationshipCompleteness interface {
    SHACLObject
}


func MakeRelationshipCompleteness() RelationshipCompleteness {
    return ConstructRelationshipCompletenessObject(&RelationshipCompletenessObject{}, relationshipCompletenessType)
}

func MakeRelationshipCompletenessRef() Ref[RelationshipCompleteness] {
    o := MakeRelationshipCompleteness()
    return MakeObjectRef[RelationshipCompleteness](o)
}

func (self *RelationshipCompletenessObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *RelationshipCompletenessObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *RelationshipCompletenessObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Information about the relationship between two Elements.
type RelationshipTypeObject struct {
    SHACLObjectBase

}

// The `from` Vulnerability affects each `to` Element. The use of the `affects` type is constrained to `VexAffectedVulnAssessmentRelationship` classed relationships.
const RelationshipTypeAffects = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/affects"
// The `from` Element is amended by each `to` Element.
const RelationshipTypeAmendedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/amendedBy"
// The `from` Element is an ancestor of each `to` Element.
const RelationshipTypeAncestorOf = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/ancestorOf"
// The `from` Element is available from the additional supplier described by each `to` Element.
const RelationshipTypeAvailableFrom = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/availableFrom"
// The `from` Element is a configuration applied to each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeConfigures = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/configures"
// The `from` Element contains each `to` Element.
const RelationshipTypeContains = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/contains"
// The `from` Vulnerability is coordinatedBy the `to` Agent(s) (vendor, researcher, or consumer agent).
const RelationshipTypeCoordinatedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/coordinatedBy"
// The `from` Element has been copied to each `to` Element.
const RelationshipTypeCopiedTo = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/copiedTo"
// The `from` Agent is delegating an action to the Agent of the `to` Relationship (which must be of type invokedBy), during a LifecycleScopeType (e.g. the `to` invokedBy Relationship is being done on behalf of `from`).
const RelationshipTypeDelegatedTo = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/delegatedTo"
// The `from` Element depends on each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeDependsOn = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/dependsOn"
// The `from` Element is a descendant of each `to` Element.
const RelationshipTypeDescendantOf = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/descendantOf"
// The `from` Element describes each `to` Element. To denote the root(s) of a tree of elements in a collection, the rootElement property should be used.
const RelationshipTypeDescribes = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/describes"
// The `from` Vulnerability has no impact on each `to` Element. The use of the `doesNotAffect` is constrained to `VexNotAffectedVulnAssessmentRelationship` classed relationships.
const RelationshipTypeDoesNotAffect = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/doesNotAffect"
// The `from` archive expands out as an artifact described by each `to` Element.
const RelationshipTypeExpandsTo = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/expandsTo"
// The `from` Vulnerability has had an exploit created against it by each `to` Agent.
const RelationshipTypeExploitCreatedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/exploitCreatedBy"
// Designates a `from` Vulnerability has been fixed by the `to` Agent(s).
const RelationshipTypeFixedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedBy"
// A `from` Vulnerability has been fixed in each `to` Element. The use of the `fixedIn` type is constrained to `VexFixedVulnAssessmentRelationship` classed relationships.
const RelationshipTypeFixedIn = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/fixedIn"
// Designates a `from` Vulnerability was originally discovered by the `to` Agent(s).
const RelationshipTypeFoundBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/foundBy"
// The `from` Element generates each `to` Element.
const RelationshipTypeGenerates = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/generates"
// Every `to` Element is a file added to the `from` Element (`from` hasAddedFile `to`).
const RelationshipTypeHasAddedFile = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAddedFile"
// Relates a `from` Vulnerability and each `to` Element with a security assessment. To be used with `VulnAssessmentRelationship` types.
const RelationshipTypeHasAssessmentFor = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssessmentFor"
// Used to associate a `from` Artifact with each `to` Vulnerability.
const RelationshipTypeHasAssociatedVulnerability = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasAssociatedVulnerability"
// The `from` SoftwareArtifact is concluded by the SPDX data creator to be governed by each `to` license.
const RelationshipTypeHasConcludedLicense = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasConcludedLicense"
// The `from` Element treats each `to` Element as a data file. A data file is an artifact that stores data required or optional for the `from` Element's functionality. A data file can be a database file, an index file, a log file, an AI model file, a calibration data file, a temporary file, a backup file, and more. For AI training dataset, test dataset, test artifact, configuration data, build input data, and build output data, please consider using the more specific relationship types: `trainedOn`, `testedOn`, `hasTest`, `configures`, `hasInput`, and `hasOutput`, respectively. This relationship does not imply dependency.
const RelationshipTypeHasDataFile = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDataFile"
// The `from` SoftwareArtifact was discovered to actually contain each `to` license, for example as detected by use of automated tooling.
const RelationshipTypeHasDeclaredLicense = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeclaredLicense"
// Every `to` Element is a file deleted from the `from` Element (`from` hasDeletedFile `to`).
const RelationshipTypeHasDeletedFile = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDeletedFile"
// The `from` Element has manifest files that contain dependency information in each `to` Element.
const RelationshipTypeHasDependencyManifest = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDependencyManifest"
// The `from` Element is distributed as an artifact in each `to` Element (e.g. an RPM or archive file).
const RelationshipTypeHasDistributionArtifact = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDistributionArtifact"
// The `from` Element is documented by each `to` Element.
const RelationshipTypeHasDocumentation = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDocumentation"
// The `from` Element dynamically links in each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeHasDynamicLink = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasDynamicLink"
// Every `to` Element is considered as evidence for the `from` Element (`from` hasEvidence `to`).
const RelationshipTypeHasEvidence = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasEvidence"
// Every `to` Element is an example for the `from` Element (`from` hasExample `to`).
const RelationshipTypeHasExample = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasExample"
// The `from` Build was run on the `to` Element during a LifecycleScopeType period (e.g. the host that the build runs on).
const RelationshipTypeHasHost = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasHost"
// The `from` Build has each `to` Element as an input, during a LifecycleScopeType period.
const RelationshipTypeHasInput = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasInput"
// Every `to` Element is metadata about the `from` Element (`from` hasMetadata `to`).
const RelationshipTypeHasMetadata = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasMetadata"
// Every `to` Element is an optional component of the `from` Element (`from` hasOptionalComponent `to`).
const RelationshipTypeHasOptionalComponent = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalComponent"
// The `from` Element optionally depends on each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeHasOptionalDependency = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOptionalDependency"
// The `from` Build element generates each `to` Element as an output, during a LifecycleScopeType period.
const RelationshipTypeHasOutput = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasOutput"
// The `from` Element has a prerequisite on each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeHasPrerequisite = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasPrerequisite"
// The `from` Element has a dependency on each `to` Element, dependency is not in the distributed artifact, but assumed to be provided, during a LifecycleScopeType period.
const RelationshipTypeHasProvidedDependency = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasProvidedDependency"
// The `from` Element has a requirement on each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeHasRequirement = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasRequirement"
// Every `to` Element is a specification for the `from` Element (`from` hasSpecification `to`), during a LifecycleScopeType period.
const RelationshipTypeHasSpecification = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasSpecification"
// The `from` Element statically links in each `to` Element, during a LifecycleScopeType period.
const RelationshipTypeHasStaticLink = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasStaticLink"
// Every `to` Element is a test artifact for the `from` Element (`from` hasTest `to`), during a LifecycleScopeType period.
const RelationshipTypeHasTest = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTest"
// Every `to` Element is a test case for the `from` Element (`from` hasTestCase `to`).
const RelationshipTypeHasTestCase = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasTestCase"
// Every `to` Element is a variant the `from` Element (`from` hasVariant `to`).
const RelationshipTypeHasVariant = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/hasVariant"
// The `from` Element was invoked by the `to` Agent, during a LifecycleScopeType period (for example, a Build element that describes a build step).
const RelationshipTypeInvokedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/invokedBy"
// The `from` Element is modified by each `to` Element.
const RelationshipTypeModifiedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/modifiedBy"
// Every `to` Element is related to the `from` Element where the relationship type is not described by any of the SPDX relationship types (this relationship is directionless).
const RelationshipTypeOther = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/other"
// Every `to` Element is a packaged instance of the `from` Element (`from` packagedBy `to`).
const RelationshipTypePackagedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/packagedBy"
// Every `to` Element is a patch for the `from` Element (`from` patchedBy `to`).
const RelationshipTypePatchedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/patchedBy"
// Designates a `from` Vulnerability was made available for public use or reference by each `to` Agent.
const RelationshipTypePublishedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/publishedBy"
// Designates a `from` Vulnerability was first reported to a project, vendor, or tracking database for formal identification by each `to` Agent.
const RelationshipTypeReportedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/reportedBy"
// Designates a `from` Vulnerability's details were tracked, aggregated, and/or enriched to improve context (i.e. NVD) by each `to` Agent.
const RelationshipTypeRepublishedBy = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/republishedBy"
// The `from` SpdxDocument can be found in a serialized form in each `to` Artifact.
const RelationshipTypeSerializedInArtifact = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/serializedInArtifact"
// The `from` Element has been tested on the `to` Element(s).
const RelationshipTypeTestedOn = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/testedOn"
// The `from` Element has been trained on the `to` Element(s).
const RelationshipTypeTrainedOn = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/trainedOn"
// The `from` Vulnerability impact is being investigated for each `to` Element. The use of the `underInvestigationFor` type is constrained to `VexUnderInvestigationVulnAssessmentRelationship` classed relationships.
const RelationshipTypeUnderInvestigationFor = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/underInvestigationFor"
// The `from` Element uses each `to` Element as a tool, during a LifecycleScopeType period.
const RelationshipTypeUsesTool = "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/usesTool"

type RelationshipTypeObjectType struct {
    SHACLTypeBase
}
var relationshipTypeType RelationshipTypeObjectType

func DecodeRelationshipType (data any, path Path, context map[string]string) (Ref[RelationshipType], error) {
    return DecodeRef[RelationshipType](data, path, context, relationshipTypeType)
}

func (self RelationshipTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(RelationshipType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self RelationshipTypeObjectType) Create() SHACLObject {
    return ConstructRelationshipTypeObject(&RelationshipTypeObject{}, self)
}

func ConstructRelationshipTypeObject(o *RelationshipTypeObject, typ SHACLType) *RelationshipTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type RelationshipType interface {
    SHACLObject
}


func MakeRelationshipType() RelationshipType {
    return ConstructRelationshipTypeObject(&RelationshipTypeObject{}, relationshipTypeType)
}

func MakeRelationshipTypeRef() Ref[RelationshipType] {
    o := MakeRelationshipType()
    return MakeObjectRef[RelationshipType](o)
}

func (self *RelationshipTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *RelationshipTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *RelationshipTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A collection of SPDX Elements that could potentially be serialized.
type SpdxDocumentObject struct {
    ElementCollectionObject

    // Provides the license under which the SPDX documentation of the Element can be
    // used.
    dataLicense RefProperty[SimplelicensingAnyLicenseInfo]
    // Provides an ExternalMap of Element identifiers.
    import_ RefListProperty[ExternalMap]
    // Provides a NamespaceMap of prefixes and associated namespace partial URIs applicable to an SpdxDocument and independent of any specific serialization format or instance.
    namespaceMap RefListProperty[NamespaceMap]
}


type SpdxDocumentObjectType struct {
    SHACLTypeBase
}
var spdxDocumentType SpdxDocumentObjectType
var spdxDocumentDataLicenseContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}
var spdxDocumentImportContext = map[string]string{}
var spdxDocumentNamespaceMapContext = map[string]string{}

func DecodeSpdxDocument (data any, path Path, context map[string]string) (Ref[SpdxDocument], error) {
    return DecodeRef[SpdxDocument](data, path, context, spdxDocumentType)
}

func (self SpdxDocumentObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SpdxDocument)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/dataLicense", "dataLicense":
        val, err := DecodeSimplelicensingAnyLicenseInfo(value, path, spdxDocumentDataLicenseContext)
        if err != nil {
            return false, err
        }
        err = obj.DataLicense().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/import", "import":
        val, err := DecodeList[Ref[ExternalMap]](value, path, spdxDocumentImportContext, DecodeExternalMap)
        if err != nil {
            return false, err
        }
        err = obj.Import().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/namespaceMap", "namespaceMap":
        val, err := DecodeList[Ref[NamespaceMap]](value, path, spdxDocumentNamespaceMapContext, DecodeNamespaceMap)
        if err != nil {
            return false, err
        }
        err = obj.NamespaceMap().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SpdxDocumentObjectType) Create() SHACLObject {
    return ConstructSpdxDocumentObject(&SpdxDocumentObject{}, self)
}

func ConstructSpdxDocumentObject(o *SpdxDocumentObject, typ SHACLType) *SpdxDocumentObject {
    ConstructElementCollectionObject(&o.ElementCollectionObject, typ)
    {
        validators := []Validator[Ref[SimplelicensingAnyLicenseInfo]]{}
        o.dataLicense = NewRefProperty[SimplelicensingAnyLicenseInfo]("dataLicense", validators)
    }
    {
        validators := []Validator[Ref[ExternalMap]]{}
        o.import_ = NewRefListProperty[ExternalMap]("import_", validators)
    }
    {
        validators := []Validator[Ref[NamespaceMap]]{}
        o.namespaceMap = NewRefListProperty[NamespaceMap]("namespaceMap", validators)
    }
    return o
}

type SpdxDocument interface {
    ElementCollection
    DataLicense() RefPropertyInterface[SimplelicensingAnyLicenseInfo]
    Import() ListPropertyInterface[Ref[ExternalMap]]
    NamespaceMap() ListPropertyInterface[Ref[NamespaceMap]]
}


func MakeSpdxDocument() SpdxDocument {
    return ConstructSpdxDocumentObject(&SpdxDocumentObject{}, spdxDocumentType)
}

func MakeSpdxDocumentRef() Ref[SpdxDocument] {
    o := MakeSpdxDocument()
    return MakeObjectRef[SpdxDocument](o)
}

func (self *SpdxDocumentObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementCollectionObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("dataLicense")
        if ! self.dataLicense.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("import_")
        if ! self.import_.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("namespaceMap")
        if ! self.namespaceMap.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SpdxDocumentObject) Walk(path Path, visit Visit) {
    self.ElementCollectionObject.Walk(path, visit)
    self.dataLicense.Walk(path, visit)
    self.import_.Walk(path, visit)
    self.namespaceMap.Walk(path, visit)
}


func (self *SpdxDocumentObject) DataLicense() RefPropertyInterface[SimplelicensingAnyLicenseInfo] { return &self.dataLicense }
func (self *SpdxDocumentObject) Import() ListPropertyInterface[Ref[ExternalMap]] { return &self.import_ }
func (self *SpdxDocumentObject) NamespaceMap() ListPropertyInterface[Ref[NamespaceMap]] { return &self.namespaceMap }

func (self *SpdxDocumentObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementCollectionObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.dataLicense.IsSet() {
        val, err := EncodeRef[SimplelicensingAnyLicenseInfo](self.dataLicense.Get(), path.PushPath("dataLicense"), spdxDocumentDataLicenseContext, state)
        if err != nil {
            return err
        }
        data["dataLicense"] = val
    }
    if self.import_.IsSet() {
        val, err := EncodeList[Ref[ExternalMap]](self.import_.Get(), path.PushPath("import_"), spdxDocumentImportContext, state, EncodeRef[ExternalMap])
        if err != nil {
            return err
        }
        data["import"] = val
    }
    if self.namespaceMap.IsSet() {
        val, err := EncodeList[Ref[NamespaceMap]](self.namespaceMap.Get(), path.PushPath("namespaceMap"), spdxDocumentNamespaceMapContext, state, EncodeRef[NamespaceMap])
        if err != nil {
            return err
        }
        data["namespaceMap"] = val
    }
    return nil
}

// Indicates the type of support that is associated with an artifact.
type SupportTypeObject struct {
    SHACLObjectBase

}

// in addition to being supported by the supplier, the software is known to have been deployed and is in use.  For a software as a service provider, this implies the software is now available as a service.
const SupportTypeDeployed = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/deployed"
// the artifact is in active development and is not considered ready for formal support from the supplier.
const SupportTypeDevelopment = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/development"
// there is a defined end of support for the artifact from the supplier.  This may also be referred to as end of life. There is a validUntilDate that can be used to signal when support ends for the artifact.
const SupportTypeEndOfSupport = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/endOfSupport"
// the artifact has been released, and there is limited support available from the supplier. There is a validUntilDate that can provide additional information about the duration of support.
const SupportTypeLimitedSupport = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/limitedSupport"
// no assertion about the type of support is made.   This is considered the default if no other support type is used.
const SupportTypeNoAssertion = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noAssertion"
// there is no support for the artifact from the supplier, consumer assumes any support obligations.
const SupportTypeNoSupport = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noSupport"
// the artifact has been released, and is supported from the supplier.   There is a validUntilDate that can provide additional information about the duration of support.
const SupportTypeSupport = "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/support"

type SupportTypeObjectType struct {
    SHACLTypeBase
}
var supportTypeType SupportTypeObjectType

func DecodeSupportType (data any, path Path, context map[string]string) (Ref[SupportType], error) {
    return DecodeRef[SupportType](data, path, context, supportTypeType)
}

func (self SupportTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SupportType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SupportTypeObjectType) Create() SHACLObject {
    return ConstructSupportTypeObject(&SupportTypeObject{}, self)
}

func ConstructSupportTypeObject(o *SupportTypeObject, typ SHACLType) *SupportTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SupportType interface {
    SHACLObject
}


func MakeSupportType() SupportType {
    return ConstructSupportTypeObject(&SupportTypeObject{}, supportTypeType)
}

func MakeSupportTypeRef() Ref[SupportType] {
    o := MakeSupportType()
    return MakeObjectRef[SupportType](o)
}

func (self *SupportTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SupportTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SupportTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// An element of hardware and/or software utilized to carry out a particular function.
type ToolObject struct {
    ElementObject

}


type ToolObjectType struct {
    SHACLTypeBase
}
var toolType ToolObjectType

func DecodeTool (data any, path Path, context map[string]string) (Ref[Tool], error) {
    return DecodeRef[Tool](data, path, context, toolType)
}

func (self ToolObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Tool)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ToolObjectType) Create() SHACLObject {
    return ConstructToolObject(&ToolObject{}, self)
}

func ConstructToolObject(o *ToolObject, typ SHACLType) *ToolObject {
    ConstructElementObject(&o.ElementObject, typ)
    return o
}

type Tool interface {
    Element
}


func MakeTool() Tool {
    return ConstructToolObject(&ToolObject{}, toolType)
}

func MakeToolRef() Ref[Tool] {
    o := MakeTool()
    return MakeObjectRef[Tool](o)
}

func (self *ToolObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ToolObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
}



func (self *ToolObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Categories of confidentiality level.
type DatasetConfidentialityLevelTypeObject struct {
    SHACLObjectBase

}

// Data points in the dataset can be shared only with specific organizations and their clients on a need to know basis.
const DatasetConfidentialityLevelTypeAmber = "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/amber"
// Dataset may be distributed freely, without restriction.
const DatasetConfidentialityLevelTypeClear = "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/clear"
// Dataset can be shared within a community of peers and partners.
const DatasetConfidentialityLevelTypeGreen = "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/green"
// Data points in the dataset are highly confidential and can only be shared with named recipients.
const DatasetConfidentialityLevelTypeRed = "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/red"

type DatasetConfidentialityLevelTypeObjectType struct {
    SHACLTypeBase
}
var datasetConfidentialityLevelTypeType DatasetConfidentialityLevelTypeObjectType

func DecodeDatasetConfidentialityLevelType (data any, path Path, context map[string]string) (Ref[DatasetConfidentialityLevelType], error) {
    return DecodeRef[DatasetConfidentialityLevelType](data, path, context, datasetConfidentialityLevelTypeType)
}

func (self DatasetConfidentialityLevelTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(DatasetConfidentialityLevelType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self DatasetConfidentialityLevelTypeObjectType) Create() SHACLObject {
    return ConstructDatasetConfidentialityLevelTypeObject(&DatasetConfidentialityLevelTypeObject{}, self)
}

func ConstructDatasetConfidentialityLevelTypeObject(o *DatasetConfidentialityLevelTypeObject, typ SHACLType) *DatasetConfidentialityLevelTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type DatasetConfidentialityLevelType interface {
    SHACLObject
}


func MakeDatasetConfidentialityLevelType() DatasetConfidentialityLevelType {
    return ConstructDatasetConfidentialityLevelTypeObject(&DatasetConfidentialityLevelTypeObject{}, datasetConfidentialityLevelTypeType)
}

func MakeDatasetConfidentialityLevelTypeRef() Ref[DatasetConfidentialityLevelType] {
    o := MakeDatasetConfidentialityLevelType()
    return MakeObjectRef[DatasetConfidentialityLevelType](o)
}

func (self *DatasetConfidentialityLevelTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *DatasetConfidentialityLevelTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *DatasetConfidentialityLevelTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Availability of dataset.
type DatasetDatasetAvailabilityTypeObject struct {
    SHACLObjectBase

}

// the dataset is not publicly available and can only be accessed after affirmatively accepting terms on a clickthrough webpage.
const DatasetDatasetAvailabilityTypeClickthrough = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/clickthrough"
// the dataset is publicly available and can be downloaded directly.
const DatasetDatasetAvailabilityTypeDirectDownload = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/directDownload"
// the dataset is publicly available, but not all at once, and can only be accessed through queries which return parts of the dataset.
const DatasetDatasetAvailabilityTypeQuery = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/query"
// the dataset is not publicly available and an email registration is required before accessing the dataset, although without an affirmative acceptance of terms.
const DatasetDatasetAvailabilityTypeRegistration = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/registration"
// the dataset provider is not making available the underlying data and the dataset must be reassembled, typically using the provided script for scraping the data.
const DatasetDatasetAvailabilityTypeScrapingScript = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/scrapingScript"

type DatasetDatasetAvailabilityTypeObjectType struct {
    SHACLTypeBase
}
var datasetDatasetAvailabilityTypeType DatasetDatasetAvailabilityTypeObjectType

func DecodeDatasetDatasetAvailabilityType (data any, path Path, context map[string]string) (Ref[DatasetDatasetAvailabilityType], error) {
    return DecodeRef[DatasetDatasetAvailabilityType](data, path, context, datasetDatasetAvailabilityTypeType)
}

func (self DatasetDatasetAvailabilityTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(DatasetDatasetAvailabilityType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self DatasetDatasetAvailabilityTypeObjectType) Create() SHACLObject {
    return ConstructDatasetDatasetAvailabilityTypeObject(&DatasetDatasetAvailabilityTypeObject{}, self)
}

func ConstructDatasetDatasetAvailabilityTypeObject(o *DatasetDatasetAvailabilityTypeObject, typ SHACLType) *DatasetDatasetAvailabilityTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type DatasetDatasetAvailabilityType interface {
    SHACLObject
}


func MakeDatasetDatasetAvailabilityType() DatasetDatasetAvailabilityType {
    return ConstructDatasetDatasetAvailabilityTypeObject(&DatasetDatasetAvailabilityTypeObject{}, datasetDatasetAvailabilityTypeType)
}

func MakeDatasetDatasetAvailabilityTypeRef() Ref[DatasetDatasetAvailabilityType] {
    o := MakeDatasetDatasetAvailabilityType()
    return MakeObjectRef[DatasetDatasetAvailabilityType](o)
}

func (self *DatasetDatasetAvailabilityTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *DatasetDatasetAvailabilityTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *DatasetDatasetAvailabilityTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Enumeration of dataset types.
type DatasetDatasetTypeObject struct {
    SHACLObjectBase

}

// data is audio based, such as a collection of music from the 80s.
const DatasetDatasetTypeAudio = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/audio"
// data that is classified into a discrete number of categories, such as the eye color of a population of people.
const DatasetDatasetTypeCategorical = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/categorical"
// data is in the form of a graph where entries are somehow related to each other through edges, such a social network of friends.
const DatasetDatasetTypeGraph = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/graph"
// data is a collection of images such as pictures of animals.
const DatasetDatasetTypeImage = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/image"
// data type is not known.
const DatasetDatasetTypeNoAssertion = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/noAssertion"
// data consists only of numeric entries.
const DatasetDatasetTypeNumeric = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/numeric"
// data is of a type not included in this list.
const DatasetDatasetTypeOther = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/other"
// data is recorded from a physical sensor, such as a thermometer reading or biometric device.
const DatasetDatasetTypeSensor = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/sensor"
// data is stored in tabular format or retrieved from a relational database.
const DatasetDatasetTypeStructured = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/structured"
// data describes the syntax or semantics of a language or text, such as a parse tree used for natural language processing.
const DatasetDatasetTypeSyntactic = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/syntactic"
// data consists of unstructured text, such as a book, Wikipedia article (without images), or transcript.
const DatasetDatasetTypeText = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/text"
// data is recorded in an ordered sequence of timestamped entries, such as the price of a stock over the course of a day.
const DatasetDatasetTypeTimeseries = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timeseries"
// data is recorded with a timestamp for each entry, but not necessarily ordered or at specific intervals, such as when a taxi ride starts and ends.
const DatasetDatasetTypeTimestamp = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timestamp"
// data is video based, such as a collection of movie clips featuring Tom Hanks.
const DatasetDatasetTypeVideo = "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/video"

type DatasetDatasetTypeObjectType struct {
    SHACLTypeBase
}
var datasetDatasetTypeType DatasetDatasetTypeObjectType

func DecodeDatasetDatasetType (data any, path Path, context map[string]string) (Ref[DatasetDatasetType], error) {
    return DecodeRef[DatasetDatasetType](data, path, context, datasetDatasetTypeType)
}

func (self DatasetDatasetTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(DatasetDatasetType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self DatasetDatasetTypeObjectType) Create() SHACLObject {
    return ConstructDatasetDatasetTypeObject(&DatasetDatasetTypeObject{}, self)
}

func ConstructDatasetDatasetTypeObject(o *DatasetDatasetTypeObject, typ SHACLType) *DatasetDatasetTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type DatasetDatasetType interface {
    SHACLObject
}


func MakeDatasetDatasetType() DatasetDatasetType {
    return ConstructDatasetDatasetTypeObject(&DatasetDatasetTypeObject{}, datasetDatasetTypeType)
}

func MakeDatasetDatasetTypeRef() Ref[DatasetDatasetType] {
    o := MakeDatasetDatasetType()
    return MakeObjectRef[DatasetDatasetType](o)
}

func (self *DatasetDatasetTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *DatasetDatasetTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *DatasetDatasetTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Abstract class for additional text intended to be added to a License, but

// which is not itself a standalone License.
type ExpandedlicensingLicenseAdditionObject struct {
    ElementObject

    // Identifies the full text of a LicenseAddition.
    expandedlicensingAdditionText Property[string]
    // Specifies whether an additional text identifier has been marked as deprecated.
    expandedlicensingIsDeprecatedAdditionId Property[bool]
    // Identifies all the text and metadata associated with a license in the license
    // XML format.
    expandedlicensingLicenseXml Property[string]
    // Specifies the licenseId that is preferred to be used in place of a deprecated
    // License or LicenseAddition.
    expandedlicensingObsoletedBy Property[string]
    // Contains a URL where the License or LicenseAddition can be found in use.
    expandedlicensingSeeAlso ListProperty[string]
    // Identifies the full text of a LicenseAddition, in SPDX templating format.
    expandedlicensingStandardAdditionTemplate Property[string]
}


type ExpandedlicensingLicenseAdditionObjectType struct {
    SHACLTypeBase
}
var expandedlicensingLicenseAdditionType ExpandedlicensingLicenseAdditionObjectType
var expandedlicensingLicenseAdditionExpandedlicensingAdditionTextContext = map[string]string{}
var expandedlicensingLicenseAdditionExpandedlicensingIsDeprecatedAdditionIdContext = map[string]string{}
var expandedlicensingLicenseAdditionExpandedlicensingLicenseXmlContext = map[string]string{}
var expandedlicensingLicenseAdditionExpandedlicensingObsoletedByContext = map[string]string{}
var expandedlicensingLicenseAdditionExpandedlicensingSeeAlsoContext = map[string]string{}
var expandedlicensingLicenseAdditionExpandedlicensingStandardAdditionTemplateContext = map[string]string{}

func DecodeExpandedlicensingLicenseAddition (data any, path Path, context map[string]string) (Ref[ExpandedlicensingLicenseAddition], error) {
    return DecodeRef[ExpandedlicensingLicenseAddition](data, path, context, expandedlicensingLicenseAdditionType)
}

func (self ExpandedlicensingLicenseAdditionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingLicenseAddition)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/additionText", "expandedlicensing_additionText":
        val, err := DecodeString(value, path, expandedlicensingLicenseAdditionExpandedlicensingAdditionTextContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingAdditionText().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/isDeprecatedAdditionId", "expandedlicensing_isDeprecatedAdditionId":
        val, err := DecodeBoolean(value, path, expandedlicensingLicenseAdditionExpandedlicensingIsDeprecatedAdditionIdContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingIsDeprecatedAdditionId().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/licenseXml", "expandedlicensing_licenseXml":
        val, err := DecodeString(value, path, expandedlicensingLicenseAdditionExpandedlicensingLicenseXmlContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingLicenseXml().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/obsoletedBy", "expandedlicensing_obsoletedBy":
        val, err := DecodeString(value, path, expandedlicensingLicenseAdditionExpandedlicensingObsoletedByContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingObsoletedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/seeAlso", "expandedlicensing_seeAlso":
        val, err := DecodeList[string](value, path, expandedlicensingLicenseAdditionExpandedlicensingSeeAlsoContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingSeeAlso().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/standardAdditionTemplate", "expandedlicensing_standardAdditionTemplate":
        val, err := DecodeString(value, path, expandedlicensingLicenseAdditionExpandedlicensingStandardAdditionTemplateContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingStandardAdditionTemplate().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingLicenseAdditionObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingLicenseAdditionObject(&ExpandedlicensingLicenseAdditionObject{}, self)
}

func ConstructExpandedlicensingLicenseAdditionObject(o *ExpandedlicensingLicenseAdditionObject, typ SHACLType) *ExpandedlicensingLicenseAdditionObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[string]{}
        o.expandedlicensingAdditionText = NewProperty[string]("expandedlicensingAdditionText", validators)
    }
    {
        validators := []Validator[bool]{}
        o.expandedlicensingIsDeprecatedAdditionId = NewProperty[bool]("expandedlicensingIsDeprecatedAdditionId", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingLicenseXml = NewProperty[string]("expandedlicensingLicenseXml", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingObsoletedBy = NewProperty[string]("expandedlicensingObsoletedBy", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingSeeAlso = NewListProperty[string]("expandedlicensingSeeAlso", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingStandardAdditionTemplate = NewProperty[string]("expandedlicensingStandardAdditionTemplate", validators)
    }
    return o
}

type ExpandedlicensingLicenseAddition interface {
    Element
    ExpandedlicensingAdditionText() PropertyInterface[string]
    ExpandedlicensingIsDeprecatedAdditionId() PropertyInterface[bool]
    ExpandedlicensingLicenseXml() PropertyInterface[string]
    ExpandedlicensingObsoletedBy() PropertyInterface[string]
    ExpandedlicensingSeeAlso() ListPropertyInterface[string]
    ExpandedlicensingStandardAdditionTemplate() PropertyInterface[string]
}



func (self *ExpandedlicensingLicenseAdditionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingAdditionText")
        if ! self.expandedlicensingAdditionText.Check(prop_path, handler) {
            valid = false
        }
        if ! self.expandedlicensingAdditionText.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"expandedlicensingAdditionText", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingIsDeprecatedAdditionId")
        if ! self.expandedlicensingIsDeprecatedAdditionId.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingLicenseXml")
        if ! self.expandedlicensingLicenseXml.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingObsoletedBy")
        if ! self.expandedlicensingObsoletedBy.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingSeeAlso")
        if ! self.expandedlicensingSeeAlso.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingStandardAdditionTemplate")
        if ! self.expandedlicensingStandardAdditionTemplate.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingLicenseAdditionObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.expandedlicensingAdditionText.Walk(path, visit)
    self.expandedlicensingIsDeprecatedAdditionId.Walk(path, visit)
    self.expandedlicensingLicenseXml.Walk(path, visit)
    self.expandedlicensingObsoletedBy.Walk(path, visit)
    self.expandedlicensingSeeAlso.Walk(path, visit)
    self.expandedlicensingStandardAdditionTemplate.Walk(path, visit)
}


func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingAdditionText() PropertyInterface[string] { return &self.expandedlicensingAdditionText }
func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingIsDeprecatedAdditionId() PropertyInterface[bool] { return &self.expandedlicensingIsDeprecatedAdditionId }
func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingLicenseXml() PropertyInterface[string] { return &self.expandedlicensingLicenseXml }
func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingObsoletedBy() PropertyInterface[string] { return &self.expandedlicensingObsoletedBy }
func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingSeeAlso() ListPropertyInterface[string] { return &self.expandedlicensingSeeAlso }
func (self *ExpandedlicensingLicenseAdditionObject) ExpandedlicensingStandardAdditionTemplate() PropertyInterface[string] { return &self.expandedlicensingStandardAdditionTemplate }

func (self *ExpandedlicensingLicenseAdditionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingAdditionText.IsSet() {
        val, err := EncodeString(self.expandedlicensingAdditionText.Get(), path.PushPath("expandedlicensingAdditionText"), expandedlicensingLicenseAdditionExpandedlicensingAdditionTextContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_additionText"] = val
    }
    if self.expandedlicensingIsDeprecatedAdditionId.IsSet() {
        val, err := EncodeBoolean(self.expandedlicensingIsDeprecatedAdditionId.Get(), path.PushPath("expandedlicensingIsDeprecatedAdditionId"), expandedlicensingLicenseAdditionExpandedlicensingIsDeprecatedAdditionIdContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_isDeprecatedAdditionId"] = val
    }
    if self.expandedlicensingLicenseXml.IsSet() {
        val, err := EncodeString(self.expandedlicensingLicenseXml.Get(), path.PushPath("expandedlicensingLicenseXml"), expandedlicensingLicenseAdditionExpandedlicensingLicenseXmlContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_licenseXml"] = val
    }
    if self.expandedlicensingObsoletedBy.IsSet() {
        val, err := EncodeString(self.expandedlicensingObsoletedBy.Get(), path.PushPath("expandedlicensingObsoletedBy"), expandedlicensingLicenseAdditionExpandedlicensingObsoletedByContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_obsoletedBy"] = val
    }
    if self.expandedlicensingSeeAlso.IsSet() {
        val, err := EncodeList[string](self.expandedlicensingSeeAlso.Get(), path.PushPath("expandedlicensingSeeAlso"), expandedlicensingLicenseAdditionExpandedlicensingSeeAlsoContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["expandedlicensing_seeAlso"] = val
    }
    if self.expandedlicensingStandardAdditionTemplate.IsSet() {
        val, err := EncodeString(self.expandedlicensingStandardAdditionTemplate.Get(), path.PushPath("expandedlicensingStandardAdditionTemplate"), expandedlicensingLicenseAdditionExpandedlicensingStandardAdditionTemplateContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_standardAdditionTemplate"] = val
    }
    return nil
}

// A license exception that is listed on the SPDX Exceptions list.
type ExpandedlicensingListedLicenseExceptionObject struct {
    ExpandedlicensingLicenseAdditionObject

    // Specifies the SPDX License List version in which this license or exception
    // identifier was deprecated.
    expandedlicensingDeprecatedVersion Property[string]
    // Specifies the SPDX License List version in which this ListedLicense or
    // ListedLicenseException identifier was first added.
    expandedlicensingListVersionAdded Property[string]
}


type ExpandedlicensingListedLicenseExceptionObjectType struct {
    SHACLTypeBase
}
var expandedlicensingListedLicenseExceptionType ExpandedlicensingListedLicenseExceptionObjectType
var expandedlicensingListedLicenseExceptionExpandedlicensingDeprecatedVersionContext = map[string]string{}
var expandedlicensingListedLicenseExceptionExpandedlicensingListVersionAddedContext = map[string]string{}

func DecodeExpandedlicensingListedLicenseException (data any, path Path, context map[string]string) (Ref[ExpandedlicensingListedLicenseException], error) {
    return DecodeRef[ExpandedlicensingListedLicenseException](data, path, context, expandedlicensingListedLicenseExceptionType)
}

func (self ExpandedlicensingListedLicenseExceptionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingListedLicenseException)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/deprecatedVersion", "expandedlicensing_deprecatedVersion":
        val, err := DecodeString(value, path, expandedlicensingListedLicenseExceptionExpandedlicensingDeprecatedVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingDeprecatedVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/listVersionAdded", "expandedlicensing_listVersionAdded":
        val, err := DecodeString(value, path, expandedlicensingListedLicenseExceptionExpandedlicensingListVersionAddedContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingListVersionAdded().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingListedLicenseExceptionObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingListedLicenseExceptionObject(&ExpandedlicensingListedLicenseExceptionObject{}, self)
}

func ConstructExpandedlicensingListedLicenseExceptionObject(o *ExpandedlicensingListedLicenseExceptionObject, typ SHACLType) *ExpandedlicensingListedLicenseExceptionObject {
    ConstructExpandedlicensingLicenseAdditionObject(&o.ExpandedlicensingLicenseAdditionObject, typ)
    {
        validators := []Validator[string]{}
        o.expandedlicensingDeprecatedVersion = NewProperty[string]("expandedlicensingDeprecatedVersion", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingListVersionAdded = NewProperty[string]("expandedlicensingListVersionAdded", validators)
    }
    return o
}

type ExpandedlicensingListedLicenseException interface {
    ExpandedlicensingLicenseAddition
    ExpandedlicensingDeprecatedVersion() PropertyInterface[string]
    ExpandedlicensingListVersionAdded() PropertyInterface[string]
}


func MakeExpandedlicensingListedLicenseException() ExpandedlicensingListedLicenseException {
    return ConstructExpandedlicensingListedLicenseExceptionObject(&ExpandedlicensingListedLicenseExceptionObject{}, expandedlicensingListedLicenseExceptionType)
}

func MakeExpandedlicensingListedLicenseExceptionRef() Ref[ExpandedlicensingListedLicenseException] {
    o := MakeExpandedlicensingListedLicenseException()
    return MakeObjectRef[ExpandedlicensingListedLicenseException](o)
}

func (self *ExpandedlicensingListedLicenseExceptionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingLicenseAdditionObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingDeprecatedVersion")
        if ! self.expandedlicensingDeprecatedVersion.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingListVersionAdded")
        if ! self.expandedlicensingListVersionAdded.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingListedLicenseExceptionObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingLicenseAdditionObject.Walk(path, visit)
    self.expandedlicensingDeprecatedVersion.Walk(path, visit)
    self.expandedlicensingListVersionAdded.Walk(path, visit)
}


func (self *ExpandedlicensingListedLicenseExceptionObject) ExpandedlicensingDeprecatedVersion() PropertyInterface[string] { return &self.expandedlicensingDeprecatedVersion }
func (self *ExpandedlicensingListedLicenseExceptionObject) ExpandedlicensingListVersionAdded() PropertyInterface[string] { return &self.expandedlicensingListVersionAdded }

func (self *ExpandedlicensingListedLicenseExceptionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingLicenseAdditionObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingDeprecatedVersion.IsSet() {
        val, err := EncodeString(self.expandedlicensingDeprecatedVersion.Get(), path.PushPath("expandedlicensingDeprecatedVersion"), expandedlicensingListedLicenseExceptionExpandedlicensingDeprecatedVersionContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_deprecatedVersion"] = val
    }
    if self.expandedlicensingListVersionAdded.IsSet() {
        val, err := EncodeString(self.expandedlicensingListVersionAdded.Get(), path.PushPath("expandedlicensingListVersionAdded"), expandedlicensingListedLicenseExceptionExpandedlicensingListVersionAddedContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_listVersionAdded"] = val
    }
    return nil
}

// A property name with an associated value.
type ExtensionCdxPropertyEntryObject struct {
    SHACLObjectBase

    // A name used in a CdxPropertyEntry name-value pair.
    extensionCdxPropName Property[string]
    // A value used in a CdxPropertyEntry name-value pair.
    extensionCdxPropValue Property[string]
}


type ExtensionCdxPropertyEntryObjectType struct {
    SHACLTypeBase
}
var extensionCdxPropertyEntryType ExtensionCdxPropertyEntryObjectType
var extensionCdxPropertyEntryExtensionCdxPropNameContext = map[string]string{}
var extensionCdxPropertyEntryExtensionCdxPropValueContext = map[string]string{}

func DecodeExtensionCdxPropertyEntry (data any, path Path, context map[string]string) (Ref[ExtensionCdxPropertyEntry], error) {
    return DecodeRef[ExtensionCdxPropertyEntry](data, path, context, extensionCdxPropertyEntryType)
}

func (self ExtensionCdxPropertyEntryObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExtensionCdxPropertyEntry)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Extension/cdxPropName", "extension_cdxPropName":
        val, err := DecodeString(value, path, extensionCdxPropertyEntryExtensionCdxPropNameContext)
        if err != nil {
            return false, err
        }
        err = obj.ExtensionCdxPropName().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Extension/cdxPropValue", "extension_cdxPropValue":
        val, err := DecodeString(value, path, extensionCdxPropertyEntryExtensionCdxPropValueContext)
        if err != nil {
            return false, err
        }
        err = obj.ExtensionCdxPropValue().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExtensionCdxPropertyEntryObjectType) Create() SHACLObject {
    return ConstructExtensionCdxPropertyEntryObject(&ExtensionCdxPropertyEntryObject{}, self)
}

func ConstructExtensionCdxPropertyEntryObject(o *ExtensionCdxPropertyEntryObject, typ SHACLType) *ExtensionCdxPropertyEntryObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    {
        validators := []Validator[string]{}
        o.extensionCdxPropName = NewProperty[string]("extensionCdxPropName", validators)
    }
    {
        validators := []Validator[string]{}
        o.extensionCdxPropValue = NewProperty[string]("extensionCdxPropValue", validators)
    }
    return o
}

type ExtensionCdxPropertyEntry interface {
    SHACLObject
    ExtensionCdxPropName() PropertyInterface[string]
    ExtensionCdxPropValue() PropertyInterface[string]
}


func MakeExtensionCdxPropertyEntry() ExtensionCdxPropertyEntry {
    return ConstructExtensionCdxPropertyEntryObject(&ExtensionCdxPropertyEntryObject{}, extensionCdxPropertyEntryType)
}

func MakeExtensionCdxPropertyEntryRef() Ref[ExtensionCdxPropertyEntry] {
    o := MakeExtensionCdxPropertyEntry()
    return MakeObjectRef[ExtensionCdxPropertyEntry](o)
}

func (self *ExtensionCdxPropertyEntryObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("extensionCdxPropName")
        if ! self.extensionCdxPropName.Check(prop_path, handler) {
            valid = false
        }
        if ! self.extensionCdxPropName.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"extensionCdxPropName", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("extensionCdxPropValue")
        if ! self.extensionCdxPropValue.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExtensionCdxPropertyEntryObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
    self.extensionCdxPropName.Walk(path, visit)
    self.extensionCdxPropValue.Walk(path, visit)
}


func (self *ExtensionCdxPropertyEntryObject) ExtensionCdxPropName() PropertyInterface[string] { return &self.extensionCdxPropName }
func (self *ExtensionCdxPropertyEntryObject) ExtensionCdxPropValue() PropertyInterface[string] { return &self.extensionCdxPropValue }

func (self *ExtensionCdxPropertyEntryObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.extensionCdxPropName.IsSet() {
        val, err := EncodeString(self.extensionCdxPropName.Get(), path.PushPath("extensionCdxPropName"), extensionCdxPropertyEntryExtensionCdxPropNameContext, state)
        if err != nil {
            return err
        }
        data["extension_cdxPropName"] = val
    }
    if self.extensionCdxPropValue.IsSet() {
        val, err := EncodeString(self.extensionCdxPropValue.Get(), path.PushPath("extensionCdxPropValue"), extensionCdxPropertyEntryExtensionCdxPropValueContext, state)
        if err != nil {
            return err
        }
        data["extension_cdxPropValue"] = val
    }
    return nil
}

// A characterization of some aspect of an Element that is associated with the Element in a generalized fashion.
type ExtensionExtensionObject struct {
    SHACLObjectBase
    SHACLExtensibleBase

}


type ExtensionExtensionObjectType struct {
    SHACLTypeBase
}
var extensionExtensionType ExtensionExtensionObjectType

func DecodeExtensionExtension (data any, path Path, context map[string]string) (Ref[ExtensionExtension], error) {
    return DecodeRef[ExtensionExtension](data, path, context, extensionExtensionType)
}

func (self ExtensionExtensionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExtensionExtension)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExtensionExtensionObjectType) Create() SHACLObject {
    return ConstructExtensionExtensionObject(&ExtensionExtensionObject{}, self)
}

func ConstructExtensionExtensionObject(o *ExtensionExtensionObject, typ SHACLType) *ExtensionExtensionObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type ExtensionExtension interface {
    SHACLObject
}



func (self *ExtensionExtensionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExtensionExtensionObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *ExtensionExtensionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    self.SHACLExtensibleBase.EncodeExtProperties(data, path)
    return nil
}

// Specifies the CVSS base, temporal, threat, or environmental severity type.
type SecurityCvssSeverityTypeObject struct {
    SHACLObjectBase

}

// When a CVSS score is between 9.0 - 10.0
const SecurityCvssSeverityTypeCritical = "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/critical"
// When a CVSS score is between 7.0 - 8.9
const SecurityCvssSeverityTypeHigh = "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/high"
// When a CVSS score is between 0.1 - 3.9
const SecurityCvssSeverityTypeLow = "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/low"
// When a CVSS score is between 4.0 - 6.9
const SecurityCvssSeverityTypeMedium = "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/medium"
// When a CVSS score is 0.0
const SecurityCvssSeverityTypeNone = "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/none"

type SecurityCvssSeverityTypeObjectType struct {
    SHACLTypeBase
}
var securityCvssSeverityTypeType SecurityCvssSeverityTypeObjectType

func DecodeSecurityCvssSeverityType (data any, path Path, context map[string]string) (Ref[SecurityCvssSeverityType], error) {
    return DecodeRef[SecurityCvssSeverityType](data, path, context, securityCvssSeverityTypeType)
}

func (self SecurityCvssSeverityTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityCvssSeverityType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityCvssSeverityTypeObjectType) Create() SHACLObject {
    return ConstructSecurityCvssSeverityTypeObject(&SecurityCvssSeverityTypeObject{}, self)
}

func ConstructSecurityCvssSeverityTypeObject(o *SecurityCvssSeverityTypeObject, typ SHACLType) *SecurityCvssSeverityTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SecurityCvssSeverityType interface {
    SHACLObject
}


func MakeSecurityCvssSeverityType() SecurityCvssSeverityType {
    return ConstructSecurityCvssSeverityTypeObject(&SecurityCvssSeverityTypeObject{}, securityCvssSeverityTypeType)
}

func MakeSecurityCvssSeverityTypeRef() Ref[SecurityCvssSeverityType] {
    o := MakeSecurityCvssSeverityType()
    return MakeObjectRef[SecurityCvssSeverityType](o)
}

func (self *SecurityCvssSeverityTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecurityCvssSeverityTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SecurityCvssSeverityTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Specifies the exploit catalog type.
type SecurityExploitCatalogTypeObject struct {
    SHACLObjectBase

}

// CISA's Known Exploited Vulnerability (KEV) Catalog
const SecurityExploitCatalogTypeKev = "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/kev"
// Other exploit catalogs
const SecurityExploitCatalogTypeOther = "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/other"

type SecurityExploitCatalogTypeObjectType struct {
    SHACLTypeBase
}
var securityExploitCatalogTypeType SecurityExploitCatalogTypeObjectType

func DecodeSecurityExploitCatalogType (data any, path Path, context map[string]string) (Ref[SecurityExploitCatalogType], error) {
    return DecodeRef[SecurityExploitCatalogType](data, path, context, securityExploitCatalogTypeType)
}

func (self SecurityExploitCatalogTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityExploitCatalogType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityExploitCatalogTypeObjectType) Create() SHACLObject {
    return ConstructSecurityExploitCatalogTypeObject(&SecurityExploitCatalogTypeObject{}, self)
}

func ConstructSecurityExploitCatalogTypeObject(o *SecurityExploitCatalogTypeObject, typ SHACLType) *SecurityExploitCatalogTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SecurityExploitCatalogType interface {
    SHACLObject
}


func MakeSecurityExploitCatalogType() SecurityExploitCatalogType {
    return ConstructSecurityExploitCatalogTypeObject(&SecurityExploitCatalogTypeObject{}, securityExploitCatalogTypeType)
}

func MakeSecurityExploitCatalogTypeRef() Ref[SecurityExploitCatalogType] {
    o := MakeSecurityExploitCatalogType()
    return MakeObjectRef[SecurityExploitCatalogType](o)
}

func (self *SecurityExploitCatalogTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecurityExploitCatalogTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SecurityExploitCatalogTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Specifies the SSVC decision type.
type SecuritySsvcDecisionTypeObject struct {
    SHACLObjectBase

}

// The vulnerability requires attention from the organization's internal, supervisory-level and leadership-level individuals. Necessary actions include requesting assistance or information about the vulnerability, as well as publishing a notification either internally and/or externally. Typically, internal groups would meet to determine the overall response and then execute agreed upon actions. CISA recommends remediating Act vulnerabilities as soon as possible.
const SecuritySsvcDecisionTypeAct = "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/act"
// The vulnerability requires attention from the organization's internal, supervisory-level individuals. Necessary actions include requesting assistance or information about the vulnerability, and may involve publishing a notification either internally and/or externally. CISA recommends remediating Attend vulnerabilities sooner than standard update timelines.
const SecuritySsvcDecisionTypeAttend = "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/attend"
// The vulnerability does not require action at this time. The organization would continue to track the vulnerability and reassess it if new information becomes available. CISA recommends remediating Track vulnerabilities within standard update timelines.
const SecuritySsvcDecisionTypeTrack = "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/track"
// ("Track\*" in the SSVC spec) The vulnerability contains specific characteristics that may require closer monitoring for changes. CISA recommends remediating Track\* vulnerabilities within standard update timelines.
const SecuritySsvcDecisionTypeTrackStar = "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/trackStar"

type SecuritySsvcDecisionTypeObjectType struct {
    SHACLTypeBase
}
var securitySsvcDecisionTypeType SecuritySsvcDecisionTypeObjectType

func DecodeSecuritySsvcDecisionType (data any, path Path, context map[string]string) (Ref[SecuritySsvcDecisionType], error) {
    return DecodeRef[SecuritySsvcDecisionType](data, path, context, securitySsvcDecisionTypeType)
}

func (self SecuritySsvcDecisionTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecuritySsvcDecisionType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecuritySsvcDecisionTypeObjectType) Create() SHACLObject {
    return ConstructSecuritySsvcDecisionTypeObject(&SecuritySsvcDecisionTypeObject{}, self)
}

func ConstructSecuritySsvcDecisionTypeObject(o *SecuritySsvcDecisionTypeObject, typ SHACLType) *SecuritySsvcDecisionTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SecuritySsvcDecisionType interface {
    SHACLObject
}


func MakeSecuritySsvcDecisionType() SecuritySsvcDecisionType {
    return ConstructSecuritySsvcDecisionTypeObject(&SecuritySsvcDecisionTypeObject{}, securitySsvcDecisionTypeType)
}

func MakeSecuritySsvcDecisionTypeRef() Ref[SecuritySsvcDecisionType] {
    o := MakeSecuritySsvcDecisionType()
    return MakeObjectRef[SecuritySsvcDecisionType](o)
}

func (self *SecuritySsvcDecisionTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecuritySsvcDecisionTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SecuritySsvcDecisionTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Specifies the VEX justification type.
type SecurityVexJustificationTypeObject struct {
    SHACLObjectBase

}

// The software is not affected because the vulnerable component is not in the product.
const SecurityVexJustificationTypeComponentNotPresent = "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/componentNotPresent"
// Built-in inline controls or mitigations prevent an adversary from leveraging the vulnerability.
const SecurityVexJustificationTypeInlineMitigationsAlreadyExist = "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/inlineMitigationsAlreadyExist"
// The vulnerable component is present, and the component contains the vulnerable code. However, vulnerable code is used in such a way that an attacker cannot mount any anticipated attack.
const SecurityVexJustificationTypeVulnerableCodeCannotBeControlledByAdversary = "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeCannotBeControlledByAdversary"
// The affected code is not reachable through the execution of the code, including non-anticipated states of the product.
const SecurityVexJustificationTypeVulnerableCodeNotInExecutePath = "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotInExecutePath"
// The product is not affected because the code underlying the vulnerability is not present in the product.
const SecurityVexJustificationTypeVulnerableCodeNotPresent = "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotPresent"

type SecurityVexJustificationTypeObjectType struct {
    SHACLTypeBase
}
var securityVexJustificationTypeType SecurityVexJustificationTypeObjectType

func DecodeSecurityVexJustificationType (data any, path Path, context map[string]string) (Ref[SecurityVexJustificationType], error) {
    return DecodeRef[SecurityVexJustificationType](data, path, context, securityVexJustificationTypeType)
}

func (self SecurityVexJustificationTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexJustificationType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexJustificationTypeObjectType) Create() SHACLObject {
    return ConstructSecurityVexJustificationTypeObject(&SecurityVexJustificationTypeObject{}, self)
}

func ConstructSecurityVexJustificationTypeObject(o *SecurityVexJustificationTypeObject, typ SHACLType) *SecurityVexJustificationTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SecurityVexJustificationType interface {
    SHACLObject
}


func MakeSecurityVexJustificationType() SecurityVexJustificationType {
    return ConstructSecurityVexJustificationTypeObject(&SecurityVexJustificationTypeObject{}, securityVexJustificationTypeType)
}

func MakeSecurityVexJustificationTypeRef() Ref[SecurityVexJustificationType] {
    o := MakeSecurityVexJustificationType()
    return MakeObjectRef[SecurityVexJustificationType](o)
}

func (self *SecurityVexJustificationTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecurityVexJustificationTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SecurityVexJustificationTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Abstract ancestor class for all vulnerability assessments
type SecurityVulnAssessmentRelationshipObject struct {
    RelationshipObject

    // Identifies who or what supplied the artifact or VulnAssessmentRelationship
    // referenced by the Element.
    suppliedBy RefProperty[Agent]
    // Specifies an Element contained in a piece of software where a vulnerability was
    // found.
    securityAssessedElement RefProperty[SoftwareSoftwareArtifact]
    // Specifies a time when a vulnerability assessment was modified
    securityModifiedTime Property[time.Time]
    // Specifies the time when a vulnerability was published.
    securityPublishedTime Property[time.Time]
    // Specified the time and date when a vulnerability was withdrawn.
    securityWithdrawnTime Property[time.Time]
}


type SecurityVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVulnAssessmentRelationshipType SecurityVulnAssessmentRelationshipObjectType
var securityVulnAssessmentRelationshipSuppliedByContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",}
var securityVulnAssessmentRelationshipSecurityAssessedElementContext = map[string]string{}
var securityVulnAssessmentRelationshipSecurityModifiedTimeContext = map[string]string{}
var securityVulnAssessmentRelationshipSecurityPublishedTimeContext = map[string]string{}
var securityVulnAssessmentRelationshipSecurityWithdrawnTimeContext = map[string]string{}

func DecodeSecurityVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVulnAssessmentRelationship](data, path, context, securityVulnAssessmentRelationshipType)
}

func (self SecurityVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/suppliedBy", "suppliedBy":
        val, err := DecodeAgent(value, path, securityVulnAssessmentRelationshipSuppliedByContext)
        if err != nil {
            return false, err
        }
        err = obj.SuppliedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/assessedElement", "security_assessedElement":
        val, err := DecodeSoftwareSoftwareArtifact(value, path, securityVulnAssessmentRelationshipSecurityAssessedElementContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityAssessedElement().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/modifiedTime", "security_modifiedTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnAssessmentRelationshipSecurityModifiedTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityModifiedTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/publishedTime", "security_publishedTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnAssessmentRelationshipSecurityPublishedTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityPublishedTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/withdrawnTime", "security_withdrawnTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnAssessmentRelationshipSecurityWithdrawnTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityWithdrawnTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVulnAssessmentRelationshipObject(&SecurityVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVulnAssessmentRelationshipObject(o *SecurityVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVulnAssessmentRelationshipObject {
    ConstructRelationshipObject(&o.RelationshipObject, typ)
    {
        validators := []Validator[Ref[Agent]]{}
        o.suppliedBy = NewRefProperty[Agent]("suppliedBy", validators)
    }
    {
        validators := []Validator[Ref[SoftwareSoftwareArtifact]]{}
        o.securityAssessedElement = NewRefProperty[SoftwareSoftwareArtifact]("securityAssessedElement", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityModifiedTime = NewProperty[time.Time]("securityModifiedTime", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityPublishedTime = NewProperty[time.Time]("securityPublishedTime", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityWithdrawnTime = NewProperty[time.Time]("securityWithdrawnTime", validators)
    }
    return o
}

type SecurityVulnAssessmentRelationship interface {
    Relationship
    SuppliedBy() RefPropertyInterface[Agent]
    SecurityAssessedElement() RefPropertyInterface[SoftwareSoftwareArtifact]
    SecurityModifiedTime() PropertyInterface[time.Time]
    SecurityPublishedTime() PropertyInterface[time.Time]
    SecurityWithdrawnTime() PropertyInterface[time.Time]
}



func (self *SecurityVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.RelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("suppliedBy")
        if ! self.suppliedBy.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityAssessedElement")
        if ! self.securityAssessedElement.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityModifiedTime")
        if ! self.securityModifiedTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityPublishedTime")
        if ! self.securityPublishedTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityWithdrawnTime")
        if ! self.securityWithdrawnTime.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SecurityVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.RelationshipObject.Walk(path, visit)
    self.suppliedBy.Walk(path, visit)
    self.securityAssessedElement.Walk(path, visit)
    self.securityModifiedTime.Walk(path, visit)
    self.securityPublishedTime.Walk(path, visit)
    self.securityWithdrawnTime.Walk(path, visit)
}


func (self *SecurityVulnAssessmentRelationshipObject) SuppliedBy() RefPropertyInterface[Agent] { return &self.suppliedBy }
func (self *SecurityVulnAssessmentRelationshipObject) SecurityAssessedElement() RefPropertyInterface[SoftwareSoftwareArtifact] { return &self.securityAssessedElement }
func (self *SecurityVulnAssessmentRelationshipObject) SecurityModifiedTime() PropertyInterface[time.Time] { return &self.securityModifiedTime }
func (self *SecurityVulnAssessmentRelationshipObject) SecurityPublishedTime() PropertyInterface[time.Time] { return &self.securityPublishedTime }
func (self *SecurityVulnAssessmentRelationshipObject) SecurityWithdrawnTime() PropertyInterface[time.Time] { return &self.securityWithdrawnTime }

func (self *SecurityVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.RelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.suppliedBy.IsSet() {
        val, err := EncodeRef[Agent](self.suppliedBy.Get(), path.PushPath("suppliedBy"), securityVulnAssessmentRelationshipSuppliedByContext, state)
        if err != nil {
            return err
        }
        data["suppliedBy"] = val
    }
    if self.securityAssessedElement.IsSet() {
        val, err := EncodeRef[SoftwareSoftwareArtifact](self.securityAssessedElement.Get(), path.PushPath("securityAssessedElement"), securityVulnAssessmentRelationshipSecurityAssessedElementContext, state)
        if err != nil {
            return err
        }
        data["security_assessedElement"] = val
    }
    if self.securityModifiedTime.IsSet() {
        val, err := EncodeDateTime(self.securityModifiedTime.Get(), path.PushPath("securityModifiedTime"), securityVulnAssessmentRelationshipSecurityModifiedTimeContext, state)
        if err != nil {
            return err
        }
        data["security_modifiedTime"] = val
    }
    if self.securityPublishedTime.IsSet() {
        val, err := EncodeDateTime(self.securityPublishedTime.Get(), path.PushPath("securityPublishedTime"), securityVulnAssessmentRelationshipSecurityPublishedTimeContext, state)
        if err != nil {
            return err
        }
        data["security_publishedTime"] = val
    }
    if self.securityWithdrawnTime.IsSet() {
        val, err := EncodeDateTime(self.securityWithdrawnTime.Get(), path.PushPath("securityWithdrawnTime"), securityVulnAssessmentRelationshipSecurityWithdrawnTimeContext, state)
        if err != nil {
            return err
        }
        data["security_withdrawnTime"] = val
    }
    return nil
}

// Abstract class representing a license combination consisting of one or more licenses.
type SimplelicensingAnyLicenseInfoObject struct {
    ElementObject

}


type SimplelicensingAnyLicenseInfoObjectType struct {
    SHACLTypeBase
}
var simplelicensingAnyLicenseInfoType SimplelicensingAnyLicenseInfoObjectType

func DecodeSimplelicensingAnyLicenseInfo (data any, path Path, context map[string]string) (Ref[SimplelicensingAnyLicenseInfo], error) {
    return DecodeRef[SimplelicensingAnyLicenseInfo](data, path, context, simplelicensingAnyLicenseInfoType)
}

func (self SimplelicensingAnyLicenseInfoObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SimplelicensingAnyLicenseInfo)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SimplelicensingAnyLicenseInfoObjectType) Create() SHACLObject {
    return ConstructSimplelicensingAnyLicenseInfoObject(&SimplelicensingAnyLicenseInfoObject{}, self)
}

func ConstructSimplelicensingAnyLicenseInfoObject(o *SimplelicensingAnyLicenseInfoObject, typ SHACLType) *SimplelicensingAnyLicenseInfoObject {
    ConstructElementObject(&o.ElementObject, typ)
    return o
}

type SimplelicensingAnyLicenseInfo interface {
    Element
}



func (self *SimplelicensingAnyLicenseInfoObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SimplelicensingAnyLicenseInfoObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
}



func (self *SimplelicensingAnyLicenseInfoObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// An SPDX Element containing an SPDX license expression string.
type SimplelicensingLicenseExpressionObject struct {
    SimplelicensingAnyLicenseInfoObject

    // Maps a LicenseRef or AdditionRef string for a Custom License or a Custom
    // License Addition to its URI ID.
    simplelicensingCustomIdToUri RefListProperty[DictionaryEntry]
    // A string in the license expression format.
    simplelicensingLicenseExpression Property[string]
    // The version of the SPDX License List used in the license expression.
    simplelicensingLicenseListVersion Property[string]
}


type SimplelicensingLicenseExpressionObjectType struct {
    SHACLTypeBase
}
var simplelicensingLicenseExpressionType SimplelicensingLicenseExpressionObjectType
var simplelicensingLicenseExpressionSimplelicensingCustomIdToUriContext = map[string]string{}
var simplelicensingLicenseExpressionSimplelicensingLicenseExpressionContext = map[string]string{}
var simplelicensingLicenseExpressionSimplelicensingLicenseListVersionContext = map[string]string{}

func DecodeSimplelicensingLicenseExpression (data any, path Path, context map[string]string) (Ref[SimplelicensingLicenseExpression], error) {
    return DecodeRef[SimplelicensingLicenseExpression](data, path, context, simplelicensingLicenseExpressionType)
}

func (self SimplelicensingLicenseExpressionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SimplelicensingLicenseExpression)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/customIdToUri", "simplelicensing_customIdToUri":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, simplelicensingLicenseExpressionSimplelicensingCustomIdToUriContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.SimplelicensingCustomIdToUri().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/licenseExpression", "simplelicensing_licenseExpression":
        val, err := DecodeString(value, path, simplelicensingLicenseExpressionSimplelicensingLicenseExpressionContext)
        if err != nil {
            return false, err
        }
        err = obj.SimplelicensingLicenseExpression().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/licenseListVersion", "simplelicensing_licenseListVersion":
        val, err := DecodeString(value, path, simplelicensingLicenseExpressionSimplelicensingLicenseListVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.SimplelicensingLicenseListVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SimplelicensingLicenseExpressionObjectType) Create() SHACLObject {
    return ConstructSimplelicensingLicenseExpressionObject(&SimplelicensingLicenseExpressionObject{}, self)
}

func ConstructSimplelicensingLicenseExpressionObject(o *SimplelicensingLicenseExpressionObject, typ SHACLType) *SimplelicensingLicenseExpressionObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.simplelicensingCustomIdToUri = NewRefListProperty[DictionaryEntry]("simplelicensingCustomIdToUri", validators)
    }
    {
        validators := []Validator[string]{}
        o.simplelicensingLicenseExpression = NewProperty[string]("simplelicensingLicenseExpression", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators, RegexValidator[string]{`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`})
        o.simplelicensingLicenseListVersion = NewProperty[string]("simplelicensingLicenseListVersion", validators)
    }
    return o
}

type SimplelicensingLicenseExpression interface {
    SimplelicensingAnyLicenseInfo
    SimplelicensingCustomIdToUri() ListPropertyInterface[Ref[DictionaryEntry]]
    SimplelicensingLicenseExpression() PropertyInterface[string]
    SimplelicensingLicenseListVersion() PropertyInterface[string]
}


func MakeSimplelicensingLicenseExpression() SimplelicensingLicenseExpression {
    return ConstructSimplelicensingLicenseExpressionObject(&SimplelicensingLicenseExpressionObject{}, simplelicensingLicenseExpressionType)
}

func MakeSimplelicensingLicenseExpressionRef() Ref[SimplelicensingLicenseExpression] {
    o := MakeSimplelicensingLicenseExpression()
    return MakeObjectRef[SimplelicensingLicenseExpression](o)
}

func (self *SimplelicensingLicenseExpressionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("simplelicensingCustomIdToUri")
        if ! self.simplelicensingCustomIdToUri.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("simplelicensingLicenseExpression")
        if ! self.simplelicensingLicenseExpression.Check(prop_path, handler) {
            valid = false
        }
        if ! self.simplelicensingLicenseExpression.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"simplelicensingLicenseExpression", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("simplelicensingLicenseListVersion")
        if ! self.simplelicensingLicenseListVersion.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SimplelicensingLicenseExpressionObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
    self.simplelicensingCustomIdToUri.Walk(path, visit)
    self.simplelicensingLicenseExpression.Walk(path, visit)
    self.simplelicensingLicenseListVersion.Walk(path, visit)
}


func (self *SimplelicensingLicenseExpressionObject) SimplelicensingCustomIdToUri() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.simplelicensingCustomIdToUri }
func (self *SimplelicensingLicenseExpressionObject) SimplelicensingLicenseExpression() PropertyInterface[string] { return &self.simplelicensingLicenseExpression }
func (self *SimplelicensingLicenseExpressionObject) SimplelicensingLicenseListVersion() PropertyInterface[string] { return &self.simplelicensingLicenseListVersion }

func (self *SimplelicensingLicenseExpressionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.simplelicensingCustomIdToUri.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.simplelicensingCustomIdToUri.Get(), path.PushPath("simplelicensingCustomIdToUri"), simplelicensingLicenseExpressionSimplelicensingCustomIdToUriContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["simplelicensing_customIdToUri"] = val
    }
    if self.simplelicensingLicenseExpression.IsSet() {
        val, err := EncodeString(self.simplelicensingLicenseExpression.Get(), path.PushPath("simplelicensingLicenseExpression"), simplelicensingLicenseExpressionSimplelicensingLicenseExpressionContext, state)
        if err != nil {
            return err
        }
        data["simplelicensing_licenseExpression"] = val
    }
    if self.simplelicensingLicenseListVersion.IsSet() {
        val, err := EncodeString(self.simplelicensingLicenseListVersion.Get(), path.PushPath("simplelicensingLicenseListVersion"), simplelicensingLicenseExpressionSimplelicensingLicenseListVersionContext, state)
        if err != nil {
            return err
        }
        data["simplelicensing_licenseListVersion"] = val
    }
    return nil
}

// A license or addition that is not listed on the SPDX License List.
type SimplelicensingSimpleLicensingTextObject struct {
    ElementObject

    // Identifies the full text of a License or Addition.
    simplelicensingLicenseText Property[string]
}


type SimplelicensingSimpleLicensingTextObjectType struct {
    SHACLTypeBase
}
var simplelicensingSimpleLicensingTextType SimplelicensingSimpleLicensingTextObjectType
var simplelicensingSimpleLicensingTextSimplelicensingLicenseTextContext = map[string]string{}

func DecodeSimplelicensingSimpleLicensingText (data any, path Path, context map[string]string) (Ref[SimplelicensingSimpleLicensingText], error) {
    return DecodeRef[SimplelicensingSimpleLicensingText](data, path, context, simplelicensingSimpleLicensingTextType)
}

func (self SimplelicensingSimpleLicensingTextObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SimplelicensingSimpleLicensingText)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/licenseText", "simplelicensing_licenseText":
        val, err := DecodeString(value, path, simplelicensingSimpleLicensingTextSimplelicensingLicenseTextContext)
        if err != nil {
            return false, err
        }
        err = obj.SimplelicensingLicenseText().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SimplelicensingSimpleLicensingTextObjectType) Create() SHACLObject {
    return ConstructSimplelicensingSimpleLicensingTextObject(&SimplelicensingSimpleLicensingTextObject{}, self)
}

func ConstructSimplelicensingSimpleLicensingTextObject(o *SimplelicensingSimpleLicensingTextObject, typ SHACLType) *SimplelicensingSimpleLicensingTextObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[string]{}
        o.simplelicensingLicenseText = NewProperty[string]("simplelicensingLicenseText", validators)
    }
    return o
}

type SimplelicensingSimpleLicensingText interface {
    Element
    SimplelicensingLicenseText() PropertyInterface[string]
}


func MakeSimplelicensingSimpleLicensingText() SimplelicensingSimpleLicensingText {
    return ConstructSimplelicensingSimpleLicensingTextObject(&SimplelicensingSimpleLicensingTextObject{}, simplelicensingSimpleLicensingTextType)
}

func MakeSimplelicensingSimpleLicensingTextRef() Ref[SimplelicensingSimpleLicensingText] {
    o := MakeSimplelicensingSimpleLicensingText()
    return MakeObjectRef[SimplelicensingSimpleLicensingText](o)
}

func (self *SimplelicensingSimpleLicensingTextObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("simplelicensingLicenseText")
        if ! self.simplelicensingLicenseText.Check(prop_path, handler) {
            valid = false
        }
        if ! self.simplelicensingLicenseText.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"simplelicensingLicenseText", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SimplelicensingSimpleLicensingTextObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.simplelicensingLicenseText.Walk(path, visit)
}


func (self *SimplelicensingSimpleLicensingTextObject) SimplelicensingLicenseText() PropertyInterface[string] { return &self.simplelicensingLicenseText }

func (self *SimplelicensingSimpleLicensingTextObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.simplelicensingLicenseText.IsSet() {
        val, err := EncodeString(self.simplelicensingLicenseText.Get(), path.PushPath("simplelicensingLicenseText"), simplelicensingSimpleLicensingTextSimplelicensingLicenseTextContext, state)
        if err != nil {
            return err
        }
        data["simplelicensing_licenseText"] = val
    }
    return nil
}

// A canonical, unique, immutable identifier
type SoftwareContentIdentifierObject struct {
    IntegrityMethodObject

    // Specifies the type of the content identifier.
    softwareContentIdentifierType Property[string]
    // Specifies the value of the content identifier.
    softwareContentIdentifierValue Property[string]
}


type SoftwareContentIdentifierObjectType struct {
    SHACLTypeBase
}
var softwareContentIdentifierType SoftwareContentIdentifierObjectType
var softwareContentIdentifierSoftwareContentIdentifierTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/gitoid": "gitoid",
    "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/swhid": "swhid",}
var softwareContentIdentifierSoftwareContentIdentifierValueContext = map[string]string{}

func DecodeSoftwareContentIdentifier (data any, path Path, context map[string]string) (Ref[SoftwareContentIdentifier], error) {
    return DecodeRef[SoftwareContentIdentifier](data, path, context, softwareContentIdentifierType)
}

func (self SoftwareContentIdentifierObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareContentIdentifier)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Software/contentIdentifierType", "software_contentIdentifierType":
        val, err := DecodeIRI(value, path, softwareContentIdentifierSoftwareContentIdentifierTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareContentIdentifierType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/contentIdentifierValue", "software_contentIdentifierValue":
        val, err := DecodeString(value, path, softwareContentIdentifierSoftwareContentIdentifierValueContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareContentIdentifierValue().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareContentIdentifierObjectType) Create() SHACLObject {
    return ConstructSoftwareContentIdentifierObject(&SoftwareContentIdentifierObject{}, self)
}

func ConstructSoftwareContentIdentifierObject(o *SoftwareContentIdentifierObject, typ SHACLType) *SoftwareContentIdentifierObject {
    ConstructIntegrityMethodObject(&o.IntegrityMethodObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/gitoid",
                "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/swhid",
        }})
        o.softwareContentIdentifierType = NewProperty[string]("softwareContentIdentifierType", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwareContentIdentifierValue = NewProperty[string]("softwareContentIdentifierValue", validators)
    }
    return o
}

type SoftwareContentIdentifier interface {
    IntegrityMethod
    SoftwareContentIdentifierType() PropertyInterface[string]
    SoftwareContentIdentifierValue() PropertyInterface[string]
}


func MakeSoftwareContentIdentifier() SoftwareContentIdentifier {
    return ConstructSoftwareContentIdentifierObject(&SoftwareContentIdentifierObject{}, softwareContentIdentifierType)
}

func MakeSoftwareContentIdentifierRef() Ref[SoftwareContentIdentifier] {
    o := MakeSoftwareContentIdentifier()
    return MakeObjectRef[SoftwareContentIdentifier](o)
}

func (self *SoftwareContentIdentifierObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.IntegrityMethodObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("softwareContentIdentifierType")
        if ! self.softwareContentIdentifierType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.softwareContentIdentifierType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"softwareContentIdentifierType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareContentIdentifierValue")
        if ! self.softwareContentIdentifierValue.Check(prop_path, handler) {
            valid = false
        }
        if ! self.softwareContentIdentifierValue.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"softwareContentIdentifierValue", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SoftwareContentIdentifierObject) Walk(path Path, visit Visit) {
    self.IntegrityMethodObject.Walk(path, visit)
    self.softwareContentIdentifierType.Walk(path, visit)
    self.softwareContentIdentifierValue.Walk(path, visit)
}


func (self *SoftwareContentIdentifierObject) SoftwareContentIdentifierType() PropertyInterface[string] { return &self.softwareContentIdentifierType }
func (self *SoftwareContentIdentifierObject) SoftwareContentIdentifierValue() PropertyInterface[string] { return &self.softwareContentIdentifierValue }

func (self *SoftwareContentIdentifierObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.IntegrityMethodObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.softwareContentIdentifierType.IsSet() {
        val, err := EncodeIRI(self.softwareContentIdentifierType.Get(), path.PushPath("softwareContentIdentifierType"), softwareContentIdentifierSoftwareContentIdentifierTypeContext, state)
        if err != nil {
            return err
        }
        data["software_contentIdentifierType"] = val
    }
    if self.softwareContentIdentifierValue.IsSet() {
        val, err := EncodeString(self.softwareContentIdentifierValue.Get(), path.PushPath("softwareContentIdentifierValue"), softwareContentIdentifierSoftwareContentIdentifierValueContext, state)
        if err != nil {
            return err
        }
        data["software_contentIdentifierValue"] = val
    }
    return nil
}

// Specifies the type of a content identifier.
type SoftwareContentIdentifierTypeObject struct {
    SHACLObjectBase

}

// [Gitoid](https://www.iana.org/assignments/uri-schemes/prov/gitoid), stands for [Git Object ID](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects). A gitoid of type blob is a unique hash of a binary artifact. A gitoid may represent either an [Artifact Identifier](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-identifier-types) for the software artifact or an [Input Manifest Identifier](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#input-manifest-identifier) for the software artifact's associated [Artifact Input Manifest](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-input-manifest); this ambiguity exists because the Artifact Input Manifest is itself an artifact, and the gitoid of that artifact is its valid identifier. Gitoids calculated on software artifacts (Snippet, File, or Package Elements) should be recorded in the SPDX 3.0 SoftwareArtifact's contentIdentifier property. Gitoids calculated on the Artifact Input Manifest (Input Manifest Identifier) should be recorded in the SPDX 3.0 Element's externalIdentifier property. See [OmniBOR Specification](https://github.com/omnibor/spec/), a minimalistic specification for describing software [Artifact Dependency Graphs](https://github.com/omnibor/spec/blob/eb1ee5c961c16215eb8709b2975d193a2007a35d/spec/SPEC.md#artifact-dependency-graph-adg).
const SoftwareContentIdentifierTypeGitoid = "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/gitoid"
// SoftWare Hash IDentifier, a persistent intrinsic identifier for digital artifacts, such as files, trees (also known as directories or folders), commits, and other objects typically found in version control systems. The format of the identifiers is defined in the [SWHID specification](https://www.swhid.org/specification/v1.1/4.Syntax) (ISO/IEC DIS 18670). They typically look like `swh:1:cnt:94a9ed024d3859793618152ea559a168bbcbb5e2`.
const SoftwareContentIdentifierTypeSwhid = "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType/swhid"

type SoftwareContentIdentifierTypeObjectType struct {
    SHACLTypeBase
}
var softwareContentIdentifierTypeType SoftwareContentIdentifierTypeObjectType

func DecodeSoftwareContentIdentifierType (data any, path Path, context map[string]string) (Ref[SoftwareContentIdentifierType], error) {
    return DecodeRef[SoftwareContentIdentifierType](data, path, context, softwareContentIdentifierTypeType)
}

func (self SoftwareContentIdentifierTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareContentIdentifierType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareContentIdentifierTypeObjectType) Create() SHACLObject {
    return ConstructSoftwareContentIdentifierTypeObject(&SoftwareContentIdentifierTypeObject{}, self)
}

func ConstructSoftwareContentIdentifierTypeObject(o *SoftwareContentIdentifierTypeObject, typ SHACLType) *SoftwareContentIdentifierTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SoftwareContentIdentifierType interface {
    SHACLObject
}


func MakeSoftwareContentIdentifierType() SoftwareContentIdentifierType {
    return ConstructSoftwareContentIdentifierTypeObject(&SoftwareContentIdentifierTypeObject{}, softwareContentIdentifierTypeType)
}

func MakeSoftwareContentIdentifierTypeRef() Ref[SoftwareContentIdentifierType] {
    o := MakeSoftwareContentIdentifierType()
    return MakeObjectRef[SoftwareContentIdentifierType](o)
}

func (self *SoftwareContentIdentifierTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SoftwareContentIdentifierTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SoftwareContentIdentifierTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Enumeration of the different kinds of SPDX file.
type SoftwareFileKindTypeObject struct {
    SHACLObjectBase

}

// The file represents a directory and all content stored in that directory.
const SoftwareFileKindTypeDirectory = "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/directory"
// The file represents a single file (default).
const SoftwareFileKindTypeFile = "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/file"

type SoftwareFileKindTypeObjectType struct {
    SHACLTypeBase
}
var softwareFileKindTypeType SoftwareFileKindTypeObjectType

func DecodeSoftwareFileKindType (data any, path Path, context map[string]string) (Ref[SoftwareFileKindType], error) {
    return DecodeRef[SoftwareFileKindType](data, path, context, softwareFileKindTypeType)
}

func (self SoftwareFileKindTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareFileKindType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareFileKindTypeObjectType) Create() SHACLObject {
    return ConstructSoftwareFileKindTypeObject(&SoftwareFileKindTypeObject{}, self)
}

func ConstructSoftwareFileKindTypeObject(o *SoftwareFileKindTypeObject, typ SHACLType) *SoftwareFileKindTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SoftwareFileKindType interface {
    SHACLObject
}


func MakeSoftwareFileKindType() SoftwareFileKindType {
    return ConstructSoftwareFileKindTypeObject(&SoftwareFileKindTypeObject{}, softwareFileKindTypeType)
}

func MakeSoftwareFileKindTypeRef() Ref[SoftwareFileKindType] {
    o := MakeSoftwareFileKindType()
    return MakeObjectRef[SoftwareFileKindType](o)
}

func (self *SoftwareFileKindTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SoftwareFileKindTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SoftwareFileKindTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Provides a set of values to be used to describe the common types of SBOMs that

// tools may create.
type SoftwareSbomTypeObject struct {
    SHACLObjectBase

}

// SBOM generated through analysis of artifacts (e.g., executables, packages, containers, and virtual machine images) after its build. Such analysis generally requires a variety of heuristics. In some contexts, this may also be referred to as a "3rd party" SBOM.
const SoftwareSbomTypeAnalyzed = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/analyzed"
// SBOM generated as part of the process of building the software to create a releasable artifact (e.g., executable or package) from data such as source files, dependencies, built components, build process ephemeral data, and other SBOMs.
const SoftwareSbomTypeBuild = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/build"
// SBOM provides an inventory of software that is present on a system. This may be an assembly of other SBOMs that combines analysis of configuration options, and examination of execution behavior in a (potentially simulated) deployment environment.
const SoftwareSbomTypeDeployed = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/deployed"
// SBOM of intended, planned software project or product with included components (some of which may not yet exist) for a new software artifact.
const SoftwareSbomTypeDesign = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/design"
// SBOM generated through instrumenting the system running the software, to capture only components present in the system, as well as external call-outs or dynamically loaded components. In some contexts, this may also be referred to as an "Instrumented" or "Dynamic" SBOM.
const SoftwareSbomTypeRuntime = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/runtime"
// SBOM created directly from the development environment, source files, and included dependencies used to build an product artifact.
const SoftwareSbomTypeSource = "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/source"

type SoftwareSbomTypeObjectType struct {
    SHACLTypeBase
}
var softwareSbomTypeType SoftwareSbomTypeObjectType

func DecodeSoftwareSbomType (data any, path Path, context map[string]string) (Ref[SoftwareSbomType], error) {
    return DecodeRef[SoftwareSbomType](data, path, context, softwareSbomTypeType)
}

func (self SoftwareSbomTypeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareSbomType)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareSbomTypeObjectType) Create() SHACLObject {
    return ConstructSoftwareSbomTypeObject(&SoftwareSbomTypeObject{}, self)
}

func ConstructSoftwareSbomTypeObject(o *SoftwareSbomTypeObject, typ SHACLType) *SoftwareSbomTypeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SoftwareSbomType interface {
    SHACLObject
}


func MakeSoftwareSbomType() SoftwareSbomType {
    return ConstructSoftwareSbomTypeObject(&SoftwareSbomTypeObject{}, softwareSbomTypeType)
}

func MakeSoftwareSbomTypeRef() Ref[SoftwareSbomType] {
    o := MakeSoftwareSbomType()
    return MakeObjectRef[SoftwareSbomType](o)
}

func (self *SoftwareSbomTypeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SoftwareSbomTypeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SoftwareSbomTypeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Provides information about the primary purpose of an Element.
type SoftwareSoftwarePurposeObject struct {
    SHACLObjectBase

}

// The Element is a software application.
const SoftwareSoftwarePurposeApplication = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/application"
// The Element is an archived collection of one or more files (.tar, .zip, etc.).
const SoftwareSoftwarePurposeArchive = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/archive"
// The Element is a bill of materials.
const SoftwareSoftwarePurposeBom = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/bom"
// The Element is configuration data.
const SoftwareSoftwarePurposeConfiguration = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/configuration"
// The Element is a container image which can be used by a container runtime application.
const SoftwareSoftwarePurposeContainer = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/container"
// The Element is data.
const SoftwareSoftwarePurposeData = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/data"
// The Element refers to a chipset, processor, or electronic board.
const SoftwareSoftwarePurposeDevice = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/device"
// The Element represents software that controls hardware devices.
const SoftwareSoftwarePurposeDeviceDriver = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/deviceDriver"
// The Element refers to a disk image that can be written to a disk, booted in a VM, etc. A disk image typically contains most or all of the components necessary to boot, such as bootloaders, kernels, firmware, userspace, etc.
const SoftwareSoftwarePurposeDiskImage = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/diskImage"
// The Element is documentation.
const SoftwareSoftwarePurposeDocumentation = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/documentation"
// The Element is the evidence that a specification or requirement has been fulfilled.
const SoftwareSoftwarePurposeEvidence = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/evidence"
// The Element is an Artifact that can be run on a computer.
const SoftwareSoftwarePurposeExecutable = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/executable"
// The Element is a single file which can be independently distributed (configuration file, statically linked binary, Kubernetes deployment, etc.).
const SoftwareSoftwarePurposeFile = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/file"
// The Element is a file system image that can be written to a disk (or virtual) partition.
const SoftwareSoftwarePurposeFilesystemImage = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/filesystemImage"
// The Element provides low level control over a device's hardware.
const SoftwareSoftwarePurposeFirmware = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/firmware"
// The Element is a software framework.
const SoftwareSoftwarePurposeFramework = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/framework"
// The Element is used to install software on disk.
const SoftwareSoftwarePurposeInstall = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/install"
// The Element is a software library.
const SoftwareSoftwarePurposeLibrary = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/library"
// The Element is a software manifest.
const SoftwareSoftwarePurposeManifest = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/manifest"
// The Element is a machine learning or artificial intelligence model.
const SoftwareSoftwarePurposeModel = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/model"
// The Element is a module of a piece of software.
const SoftwareSoftwarePurposeModule = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/module"
// The Element is an operating system.
const SoftwareSoftwarePurposeOperatingSystem = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/operatingSystem"
// The Element doesn't fit into any of the other categories.
const SoftwareSoftwarePurposeOther = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/other"
// The Element contains a set of changes to update, fix, or improve another Element.
const SoftwareSoftwarePurposePatch = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/patch"
// The Element represents a runtime environment.
const SoftwareSoftwarePurposePlatform = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/platform"
// The Element provides a requirement needed as input for another Element.
const SoftwareSoftwarePurposeRequirement = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/requirement"
// The Element is a single or a collection of source files.
const SoftwareSoftwarePurposeSource = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/source"
// The Element is a plan, guideline or strategy how to create, perform or analyze an application.
const SoftwareSoftwarePurposeSpecification = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/specification"
// The Element is a test used to verify functionality on an software element.
const SoftwareSoftwarePurposeTest = "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/test"

type SoftwareSoftwarePurposeObjectType struct {
    SHACLTypeBase
}
var softwareSoftwarePurposeType SoftwareSoftwarePurposeObjectType

func DecodeSoftwareSoftwarePurpose (data any, path Path, context map[string]string) (Ref[SoftwareSoftwarePurpose], error) {
    return DecodeRef[SoftwareSoftwarePurpose](data, path, context, softwareSoftwarePurposeType)
}

func (self SoftwareSoftwarePurposeObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareSoftwarePurpose)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareSoftwarePurposeObjectType) Create() SHACLObject {
    return ConstructSoftwareSoftwarePurposeObject(&SoftwareSoftwarePurposeObject{}, self)
}

func ConstructSoftwareSoftwarePurposeObject(o *SoftwareSoftwarePurposeObject, typ SHACLType) *SoftwareSoftwarePurposeObject {
    ConstructSHACLObjectBase(&o.SHACLObjectBase, typ)
    return o
}

type SoftwareSoftwarePurpose interface {
    SHACLObject
}


func MakeSoftwareSoftwarePurpose() SoftwareSoftwarePurpose {
    return ConstructSoftwareSoftwarePurposeObject(&SoftwareSoftwarePurposeObject{}, softwareSoftwarePurposeType)
}

func MakeSoftwareSoftwarePurposeRef() Ref[SoftwareSoftwarePurpose] {
    o := MakeSoftwareSoftwarePurpose()
    return MakeObjectRef[SoftwareSoftwarePurpose](o)
}

func (self *SoftwareSoftwarePurposeObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SHACLObjectBase.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SoftwareSoftwarePurposeObject) Walk(path Path, visit Visit) {
    self.SHACLObjectBase.Walk(path, visit)
}



func (self *SoftwareSoftwarePurposeObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SHACLObjectBase.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Class that describes a build instance of software/artifacts.
type BuildBuildObject struct {
    ElementObject

    // Property that describes the time at which a build stops.
    buildBuildEndTime Property[time.Time]
    // A buildId is a locally unique identifier used by a builder to identify a unique
    // instance of a build produced by it.
    buildBuildId Property[string]
    // Property describing the start time of a build.
    buildBuildStartTime Property[time.Time]
    // A buildType is a hint that is used to indicate the toolchain, platform, or
    // infrastructure that the build was invoked on.
    buildBuildType Property[string]
    // Property that describes the digest of the build configuration file used to
    // invoke a build.
    buildConfigSourceDigest RefListProperty[Hash]
    // Property describes the invocation entrypoint of a build.
    buildConfigSourceEntrypoint ListProperty[string]
    // Property that describes the URI of the build configuration source file.
    buildConfigSourceUri ListProperty[string]
    // Property describing the session in which a build is invoked.
    buildEnvironment RefListProperty[DictionaryEntry]
    // Property describing a parameter used in an instance of a build.
    buildParameter RefListProperty[DictionaryEntry]
}


type BuildBuildObjectType struct {
    SHACLTypeBase
}
var buildBuildType BuildBuildObjectType
var buildBuildBuildBuildEndTimeContext = map[string]string{}
var buildBuildBuildBuildIdContext = map[string]string{}
var buildBuildBuildBuildStartTimeContext = map[string]string{}
var buildBuildBuildBuildTypeContext = map[string]string{}
var buildBuildBuildConfigSourceDigestContext = map[string]string{}
var buildBuildBuildConfigSourceEntrypointContext = map[string]string{}
var buildBuildBuildConfigSourceUriContext = map[string]string{}
var buildBuildBuildEnvironmentContext = map[string]string{}
var buildBuildBuildParameterContext = map[string]string{}

func DecodeBuildBuild (data any, path Path, context map[string]string) (Ref[BuildBuild], error) {
    return DecodeRef[BuildBuild](data, path, context, buildBuildType)
}

func (self BuildBuildObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(BuildBuild)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Build/buildEndTime", "build_buildEndTime":
        val, err := DecodeDateTimeStamp(value, path, buildBuildBuildBuildEndTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.BuildBuildEndTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/buildId", "build_buildId":
        val, err := DecodeString(value, path, buildBuildBuildBuildIdContext)
        if err != nil {
            return false, err
        }
        err = obj.BuildBuildId().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/buildStartTime", "build_buildStartTime":
        val, err := DecodeDateTimeStamp(value, path, buildBuildBuildBuildStartTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.BuildBuildStartTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/buildType", "build_buildType":
        val, err := DecodeString(value, path, buildBuildBuildBuildTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.BuildBuildType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/configSourceDigest", "build_configSourceDigest":
        val, err := DecodeList[Ref[Hash]](value, path, buildBuildBuildConfigSourceDigestContext, DecodeHash)
        if err != nil {
            return false, err
        }
        err = obj.BuildConfigSourceDigest().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/configSourceEntrypoint", "build_configSourceEntrypoint":
        val, err := DecodeList[string](value, path, buildBuildBuildConfigSourceEntrypointContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.BuildConfigSourceEntrypoint().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/configSourceUri", "build_configSourceUri":
        val, err := DecodeList[string](value, path, buildBuildBuildConfigSourceUriContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.BuildConfigSourceUri().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/environment", "build_environment":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, buildBuildBuildEnvironmentContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.BuildEnvironment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Build/parameter", "build_parameter":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, buildBuildBuildParameterContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.BuildParameter().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self BuildBuildObjectType) Create() SHACLObject {
    return ConstructBuildBuildObject(&BuildBuildObject{}, self)
}

func ConstructBuildBuildObject(o *BuildBuildObject, typ SHACLType) *BuildBuildObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.buildBuildEndTime = NewProperty[time.Time]("buildBuildEndTime", validators)
    }
    {
        validators := []Validator[string]{}
        o.buildBuildId = NewProperty[string]("buildBuildId", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.buildBuildStartTime = NewProperty[time.Time]("buildBuildStartTime", validators)
    }
    {
        validators := []Validator[string]{}
        o.buildBuildType = NewProperty[string]("buildBuildType", validators)
    }
    {
        validators := []Validator[Ref[Hash]]{}
        o.buildConfigSourceDigest = NewRefListProperty[Hash]("buildConfigSourceDigest", validators)
    }
    {
        validators := []Validator[string]{}
        o.buildConfigSourceEntrypoint = NewListProperty[string]("buildConfigSourceEntrypoint", validators)
    }
    {
        validators := []Validator[string]{}
        o.buildConfigSourceUri = NewListProperty[string]("buildConfigSourceUri", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.buildEnvironment = NewRefListProperty[DictionaryEntry]("buildEnvironment", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.buildParameter = NewRefListProperty[DictionaryEntry]("buildParameter", validators)
    }
    return o
}

type BuildBuild interface {
    Element
    BuildBuildEndTime() PropertyInterface[time.Time]
    BuildBuildId() PropertyInterface[string]
    BuildBuildStartTime() PropertyInterface[time.Time]
    BuildBuildType() PropertyInterface[string]
    BuildConfigSourceDigest() ListPropertyInterface[Ref[Hash]]
    BuildConfigSourceEntrypoint() ListPropertyInterface[string]
    BuildConfigSourceUri() ListPropertyInterface[string]
    BuildEnvironment() ListPropertyInterface[Ref[DictionaryEntry]]
    BuildParameter() ListPropertyInterface[Ref[DictionaryEntry]]
}


func MakeBuildBuild() BuildBuild {
    return ConstructBuildBuildObject(&BuildBuildObject{}, buildBuildType)
}

func MakeBuildBuildRef() Ref[BuildBuild] {
    o := MakeBuildBuild()
    return MakeObjectRef[BuildBuild](o)
}

func (self *BuildBuildObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("buildBuildEndTime")
        if ! self.buildBuildEndTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildBuildId")
        if ! self.buildBuildId.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildBuildStartTime")
        if ! self.buildBuildStartTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildBuildType")
        if ! self.buildBuildType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.buildBuildType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"buildBuildType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildConfigSourceDigest")
        if ! self.buildConfigSourceDigest.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildConfigSourceEntrypoint")
        if ! self.buildConfigSourceEntrypoint.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildConfigSourceUri")
        if ! self.buildConfigSourceUri.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildEnvironment")
        if ! self.buildEnvironment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("buildParameter")
        if ! self.buildParameter.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *BuildBuildObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.buildBuildEndTime.Walk(path, visit)
    self.buildBuildId.Walk(path, visit)
    self.buildBuildStartTime.Walk(path, visit)
    self.buildBuildType.Walk(path, visit)
    self.buildConfigSourceDigest.Walk(path, visit)
    self.buildConfigSourceEntrypoint.Walk(path, visit)
    self.buildConfigSourceUri.Walk(path, visit)
    self.buildEnvironment.Walk(path, visit)
    self.buildParameter.Walk(path, visit)
}


func (self *BuildBuildObject) BuildBuildEndTime() PropertyInterface[time.Time] { return &self.buildBuildEndTime }
func (self *BuildBuildObject) BuildBuildId() PropertyInterface[string] { return &self.buildBuildId }
func (self *BuildBuildObject) BuildBuildStartTime() PropertyInterface[time.Time] { return &self.buildBuildStartTime }
func (self *BuildBuildObject) BuildBuildType() PropertyInterface[string] { return &self.buildBuildType }
func (self *BuildBuildObject) BuildConfigSourceDigest() ListPropertyInterface[Ref[Hash]] { return &self.buildConfigSourceDigest }
func (self *BuildBuildObject) BuildConfigSourceEntrypoint() ListPropertyInterface[string] { return &self.buildConfigSourceEntrypoint }
func (self *BuildBuildObject) BuildConfigSourceUri() ListPropertyInterface[string] { return &self.buildConfigSourceUri }
func (self *BuildBuildObject) BuildEnvironment() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.buildEnvironment }
func (self *BuildBuildObject) BuildParameter() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.buildParameter }

func (self *BuildBuildObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.buildBuildEndTime.IsSet() {
        val, err := EncodeDateTime(self.buildBuildEndTime.Get(), path.PushPath("buildBuildEndTime"), buildBuildBuildBuildEndTimeContext, state)
        if err != nil {
            return err
        }
        data["build_buildEndTime"] = val
    }
    if self.buildBuildId.IsSet() {
        val, err := EncodeString(self.buildBuildId.Get(), path.PushPath("buildBuildId"), buildBuildBuildBuildIdContext, state)
        if err != nil {
            return err
        }
        data["build_buildId"] = val
    }
    if self.buildBuildStartTime.IsSet() {
        val, err := EncodeDateTime(self.buildBuildStartTime.Get(), path.PushPath("buildBuildStartTime"), buildBuildBuildBuildStartTimeContext, state)
        if err != nil {
            return err
        }
        data["build_buildStartTime"] = val
    }
    if self.buildBuildType.IsSet() {
        val, err := EncodeString(self.buildBuildType.Get(), path.PushPath("buildBuildType"), buildBuildBuildBuildTypeContext, state)
        if err != nil {
            return err
        }
        data["build_buildType"] = val
    }
    if self.buildConfigSourceDigest.IsSet() {
        val, err := EncodeList[Ref[Hash]](self.buildConfigSourceDigest.Get(), path.PushPath("buildConfigSourceDigest"), buildBuildBuildConfigSourceDigestContext, state, EncodeRef[Hash])
        if err != nil {
            return err
        }
        data["build_configSourceDigest"] = val
    }
    if self.buildConfigSourceEntrypoint.IsSet() {
        val, err := EncodeList[string](self.buildConfigSourceEntrypoint.Get(), path.PushPath("buildConfigSourceEntrypoint"), buildBuildBuildConfigSourceEntrypointContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["build_configSourceEntrypoint"] = val
    }
    if self.buildConfigSourceUri.IsSet() {
        val, err := EncodeList[string](self.buildConfigSourceUri.Get(), path.PushPath("buildConfigSourceUri"), buildBuildBuildConfigSourceUriContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["build_configSourceUri"] = val
    }
    if self.buildEnvironment.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.buildEnvironment.Get(), path.PushPath("buildEnvironment"), buildBuildBuildEnvironmentContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["build_environment"] = val
    }
    if self.buildParameter.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.buildParameter.Get(), path.PushPath("buildParameter"), buildBuildBuildParameterContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["build_parameter"] = val
    }
    return nil
}

// Agent represents anything with the potential to act on a system.
type AgentObject struct {
    ElementObject

}


type AgentObjectType struct {
    SHACLTypeBase
}
var agentType AgentObjectType

func DecodeAgent (data any, path Path, context map[string]string) (Ref[Agent], error) {
    return DecodeRef[Agent](data, path, context, agentType)
}

func (self AgentObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Agent)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AgentObjectType) Create() SHACLObject {
    return ConstructAgentObject(&AgentObject{}, self)
}

func ConstructAgentObject(o *AgentObject, typ SHACLType) *AgentObject {
    ConstructElementObject(&o.ElementObject, typ)
    return o
}

type Agent interface {
    Element
}


func MakeAgent() Agent {
    return ConstructAgentObject(&AgentObject{}, agentType)
}

func MakeAgentRef() Ref[Agent] {
    o := MakeAgent()
    return MakeObjectRef[Agent](o)
}

func (self *AgentObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *AgentObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
}



func (self *AgentObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// An assertion made in relation to one or more elements.
type AnnotationObject struct {
    ElementObject

    // Describes the type of annotation.
    annotationType Property[string]
    // Provides information about the content type of an Element or a Property.
    contentType Property[string]
    // Commentary on an assertion that an annotator has made.
    statement Property[string]
    // An Element an annotator has made an assertion about.
    subject RefProperty[Element]
}


type AnnotationObjectType struct {
    SHACLTypeBase
}
var annotationType AnnotationObjectType
var annotationAnnotationTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/review": "review",}
var annotationContentTypeContext = map[string]string{}
var annotationStatementContext = map[string]string{}
var annotationSubjectContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement": "NoAssertionElement",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement": "NoneElement",
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}

func DecodeAnnotation (data any, path Path, context map[string]string) (Ref[Annotation], error) {
    return DecodeRef[Annotation](data, path, context, annotationType)
}

func (self AnnotationObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Annotation)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/annotationType", "annotationType":
        val, err := DecodeIRI(value, path, annotationAnnotationTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.AnnotationType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/contentType", "contentType":
        val, err := DecodeString(value, path, annotationContentTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.ContentType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/statement", "statement":
        val, err := DecodeString(value, path, annotationStatementContext)
        if err != nil {
            return false, err
        }
        err = obj.Statement().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/subject", "subject":
        val, err := DecodeElement(value, path, annotationSubjectContext)
        if err != nil {
            return false, err
        }
        err = obj.Subject().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AnnotationObjectType) Create() SHACLObject {
    return ConstructAnnotationObject(&AnnotationObject{}, self)
}

func ConstructAnnotationObject(o *AnnotationObject, typ SHACLType) *AnnotationObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType/review",
        }})
        o.annotationType = NewProperty[string]("annotationType", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators, RegexValidator[string]{`^[^\/]+\/[^\/]+$`})
        o.contentType = NewProperty[string]("contentType", validators)
    }
    {
        validators := []Validator[string]{}
        o.statement = NewProperty[string]("statement", validators)
    }
    {
        validators := []Validator[Ref[Element]]{}
        o.subject = NewRefProperty[Element]("subject", validators)
    }
    return o
}

type Annotation interface {
    Element
    AnnotationType() PropertyInterface[string]
    ContentType() PropertyInterface[string]
    Statement() PropertyInterface[string]
    Subject() RefPropertyInterface[Element]
}


func MakeAnnotation() Annotation {
    return ConstructAnnotationObject(&AnnotationObject{}, annotationType)
}

func MakeAnnotationRef() Ref[Annotation] {
    o := MakeAnnotation()
    return MakeObjectRef[Annotation](o)
}

func (self *AnnotationObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("annotationType")
        if ! self.annotationType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.annotationType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"annotationType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("contentType")
        if ! self.contentType.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("statement")
        if ! self.statement.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("subject")
        if ! self.subject.Check(prop_path, handler) {
            valid = false
        }
        if ! self.subject.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"subject", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *AnnotationObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.annotationType.Walk(path, visit)
    self.contentType.Walk(path, visit)
    self.statement.Walk(path, visit)
    self.subject.Walk(path, visit)
}


func (self *AnnotationObject) AnnotationType() PropertyInterface[string] { return &self.annotationType }
func (self *AnnotationObject) ContentType() PropertyInterface[string] { return &self.contentType }
func (self *AnnotationObject) Statement() PropertyInterface[string] { return &self.statement }
func (self *AnnotationObject) Subject() RefPropertyInterface[Element] { return &self.subject }

func (self *AnnotationObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.annotationType.IsSet() {
        val, err := EncodeIRI(self.annotationType.Get(), path.PushPath("annotationType"), annotationAnnotationTypeContext, state)
        if err != nil {
            return err
        }
        data["annotationType"] = val
    }
    if self.contentType.IsSet() {
        val, err := EncodeString(self.contentType.Get(), path.PushPath("contentType"), annotationContentTypeContext, state)
        if err != nil {
            return err
        }
        data["contentType"] = val
    }
    if self.statement.IsSet() {
        val, err := EncodeString(self.statement.Get(), path.PushPath("statement"), annotationStatementContext, state)
        if err != nil {
            return err
        }
        data["statement"] = val
    }
    if self.subject.IsSet() {
        val, err := EncodeRef[Element](self.subject.Get(), path.PushPath("subject"), annotationSubjectContext, state)
        if err != nil {
            return err
        }
        data["subject"] = val
    }
    return nil
}

// A distinct article or unit within the digital domain.
type ArtifactObject struct {
    ElementObject

    // Specifies the time an artifact was built.
    builtTime Property[time.Time]
    // Identifies from where or whom the Element originally came.
    originatedBy RefListProperty[Agent]
    // Specifies the time an artifact was released.
    releaseTime Property[time.Time]
    // The name of a relevant standard that may apply to an artifact.
    standardName ListProperty[string]
    // Identifies who or what supplied the artifact or VulnAssessmentRelationship
    // referenced by the Element.
    suppliedBy RefProperty[Agent]
    // Specifies the level of support associated with an artifact.
    supportLevel ListProperty[string]
    // Specifies until when the artifact can be used before its usage needs to be
    // reassessed.
    validUntilTime Property[time.Time]
}


type ArtifactObjectType struct {
    SHACLTypeBase
}
var artifactType ArtifactObjectType
var artifactBuiltTimeContext = map[string]string{}
var artifactOriginatedByContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",}
var artifactReleaseTimeContext = map[string]string{}
var artifactStandardNameContext = map[string]string{}
var artifactSuppliedByContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization": "SpdxOrganization",}
var artifactSupportLevelContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/deployed": "deployed",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/development": "development",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/endOfSupport": "endOfSupport",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/limitedSupport": "limitedSupport",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noAssertion": "noAssertion",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noSupport": "noSupport",
    "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/support": "support",}
var artifactValidUntilTimeContext = map[string]string{}

func DecodeArtifact (data any, path Path, context map[string]string) (Ref[Artifact], error) {
    return DecodeRef[Artifact](data, path, context, artifactType)
}

func (self ArtifactObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Artifact)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/builtTime", "builtTime":
        val, err := DecodeDateTimeStamp(value, path, artifactBuiltTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.BuiltTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/originatedBy", "originatedBy":
        val, err := DecodeList[Ref[Agent]](value, path, artifactOriginatedByContext, DecodeAgent)
        if err != nil {
            return false, err
        }
        err = obj.OriginatedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/releaseTime", "releaseTime":
        val, err := DecodeDateTimeStamp(value, path, artifactReleaseTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.ReleaseTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/standardName", "standardName":
        val, err := DecodeList[string](value, path, artifactStandardNameContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.StandardName().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/suppliedBy", "suppliedBy":
        val, err := DecodeAgent(value, path, artifactSuppliedByContext)
        if err != nil {
            return false, err
        }
        err = obj.SuppliedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/supportLevel", "supportLevel":
        val, err := DecodeList[string](value, path, artifactSupportLevelContext, DecodeIRI)
        if err != nil {
            return false, err
        }
        err = obj.SupportLevel().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/validUntilTime", "validUntilTime":
        val, err := DecodeDateTimeStamp(value, path, artifactValidUntilTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.ValidUntilTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ArtifactObjectType) Create() SHACLObject {
    return ConstructArtifactObject(&ArtifactObject{}, self)
}

func ConstructArtifactObject(o *ArtifactObject, typ SHACLType) *ArtifactObject {
    ConstructElementObject(&o.ElementObject, typ)
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.builtTime = NewProperty[time.Time]("builtTime", validators)
    }
    {
        validators := []Validator[Ref[Agent]]{}
        o.originatedBy = NewRefListProperty[Agent]("originatedBy", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.releaseTime = NewProperty[time.Time]("releaseTime", validators)
    }
    {
        validators := []Validator[string]{}
        o.standardName = NewListProperty[string]("standardName", validators)
    }
    {
        validators := []Validator[Ref[Agent]]{}
        o.suppliedBy = NewRefProperty[Agent]("suppliedBy", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/deployed",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/development",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/endOfSupport",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/limitedSupport",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noAssertion",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/noSupport",
                "https://spdx.org/rdf/3.0.1/terms/Core/SupportType/support",
        }})
        o.supportLevel = NewListProperty[string]("supportLevel", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.validUntilTime = NewProperty[time.Time]("validUntilTime", validators)
    }
    return o
}

type Artifact interface {
    Element
    BuiltTime() PropertyInterface[time.Time]
    OriginatedBy() ListPropertyInterface[Ref[Agent]]
    ReleaseTime() PropertyInterface[time.Time]
    StandardName() ListPropertyInterface[string]
    SuppliedBy() RefPropertyInterface[Agent]
    SupportLevel() ListPropertyInterface[string]
    ValidUntilTime() PropertyInterface[time.Time]
}



func (self *ArtifactObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("builtTime")
        if ! self.builtTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("originatedBy")
        if ! self.originatedBy.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("releaseTime")
        if ! self.releaseTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("standardName")
        if ! self.standardName.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("suppliedBy")
        if ! self.suppliedBy.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("supportLevel")
        if ! self.supportLevel.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("validUntilTime")
        if ! self.validUntilTime.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ArtifactObject) Walk(path Path, visit Visit) {
    self.ElementObject.Walk(path, visit)
    self.builtTime.Walk(path, visit)
    self.originatedBy.Walk(path, visit)
    self.releaseTime.Walk(path, visit)
    self.standardName.Walk(path, visit)
    self.suppliedBy.Walk(path, visit)
    self.supportLevel.Walk(path, visit)
    self.validUntilTime.Walk(path, visit)
}


func (self *ArtifactObject) BuiltTime() PropertyInterface[time.Time] { return &self.builtTime }
func (self *ArtifactObject) OriginatedBy() ListPropertyInterface[Ref[Agent]] { return &self.originatedBy }
func (self *ArtifactObject) ReleaseTime() PropertyInterface[time.Time] { return &self.releaseTime }
func (self *ArtifactObject) StandardName() ListPropertyInterface[string] { return &self.standardName }
func (self *ArtifactObject) SuppliedBy() RefPropertyInterface[Agent] { return &self.suppliedBy }
func (self *ArtifactObject) SupportLevel() ListPropertyInterface[string] { return &self.supportLevel }
func (self *ArtifactObject) ValidUntilTime() PropertyInterface[time.Time] { return &self.validUntilTime }

func (self *ArtifactObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.builtTime.IsSet() {
        val, err := EncodeDateTime(self.builtTime.Get(), path.PushPath("builtTime"), artifactBuiltTimeContext, state)
        if err != nil {
            return err
        }
        data["builtTime"] = val
    }
    if self.originatedBy.IsSet() {
        val, err := EncodeList[Ref[Agent]](self.originatedBy.Get(), path.PushPath("originatedBy"), artifactOriginatedByContext, state, EncodeRef[Agent])
        if err != nil {
            return err
        }
        data["originatedBy"] = val
    }
    if self.releaseTime.IsSet() {
        val, err := EncodeDateTime(self.releaseTime.Get(), path.PushPath("releaseTime"), artifactReleaseTimeContext, state)
        if err != nil {
            return err
        }
        data["releaseTime"] = val
    }
    if self.standardName.IsSet() {
        val, err := EncodeList[string](self.standardName.Get(), path.PushPath("standardName"), artifactStandardNameContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["standardName"] = val
    }
    if self.suppliedBy.IsSet() {
        val, err := EncodeRef[Agent](self.suppliedBy.Get(), path.PushPath("suppliedBy"), artifactSuppliedByContext, state)
        if err != nil {
            return err
        }
        data["suppliedBy"] = val
    }
    if self.supportLevel.IsSet() {
        val, err := EncodeList[string](self.supportLevel.Get(), path.PushPath("supportLevel"), artifactSupportLevelContext, state, EncodeIRI)
        if err != nil {
            return err
        }
        data["supportLevel"] = val
    }
    if self.validUntilTime.IsSet() {
        val, err := EncodeDateTime(self.validUntilTime.Get(), path.PushPath("validUntilTime"), artifactValidUntilTimeContext, state)
        if err != nil {
            return err
        }
        data["validUntilTime"] = val
    }
    return nil
}

// A collection of Elements that have a shared context.
type BundleObject struct {
    ElementCollectionObject

    // Gives information about the circumstances or unifying properties
    // that Elements of the bundle have been assembled under.
    context Property[string]
}


type BundleObjectType struct {
    SHACLTypeBase
}
var bundleType BundleObjectType
var bundleContextContext = map[string]string{}

func DecodeBundle (data any, path Path, context map[string]string) (Ref[Bundle], error) {
    return DecodeRef[Bundle](data, path, context, bundleType)
}

func (self BundleObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Bundle)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/context", "context":
        val, err := DecodeString(value, path, bundleContextContext)
        if err != nil {
            return false, err
        }
        err = obj.Context().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self BundleObjectType) Create() SHACLObject {
    return ConstructBundleObject(&BundleObject{}, self)
}

func ConstructBundleObject(o *BundleObject, typ SHACLType) *BundleObject {
    ConstructElementCollectionObject(&o.ElementCollectionObject, typ)
    {
        validators := []Validator[string]{}
        o.context = NewProperty[string]("context", validators)
    }
    return o
}

type Bundle interface {
    ElementCollection
    Context() PropertyInterface[string]
}


func MakeBundle() Bundle {
    return ConstructBundleObject(&BundleObject{}, bundleType)
}

func MakeBundleRef() Ref[Bundle] {
    o := MakeBundle()
    return MakeObjectRef[Bundle](o)
}

func (self *BundleObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ElementCollectionObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("context")
        if ! self.context.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *BundleObject) Walk(path Path, visit Visit) {
    self.ElementCollectionObject.Walk(path, visit)
    self.context.Walk(path, visit)
}


func (self *BundleObject) Context() PropertyInterface[string] { return &self.context }

func (self *BundleObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ElementCollectionObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.context.IsSet() {
        val, err := EncodeString(self.context.Get(), path.PushPath("context"), bundleContextContext, state)
        if err != nil {
            return err
        }
        data["context"] = val
    }
    return nil
}

// A mathematically calculated representation of a grouping of data.
type HashObject struct {
    IntegrityMethodObject

    // Specifies the algorithm used for calculating the hash value.
    algorithm Property[string]
    // The result of applying a hash algorithm to an Element.
    hashValue Property[string]
}


type HashObjectType struct {
    SHACLTypeBase
}
var hashType HashObjectType
var hashAlgorithmContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/adler32": "adler32",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b256": "blake2b256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b384": "blake2b384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b512": "blake2b512",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake3": "blake3",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsDilithium": "crystalsDilithium",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsKyber": "crystalsKyber",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/falcon": "falcon",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md2": "md2",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md4": "md4",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md5": "md5",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md6": "md6",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha1": "sha1",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha224": "sha224",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha256": "sha256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha384": "sha384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_224": "sha3_224",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_256": "sha3_256",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_384": "sha3_384",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_512": "sha3_512",
    "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha512": "sha512",}
var hashHashValueContext = map[string]string{}

func DecodeHash (data any, path Path, context map[string]string) (Ref[Hash], error) {
    return DecodeRef[Hash](data, path, context, hashType)
}

func (self HashObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Hash)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/algorithm", "algorithm":
        val, err := DecodeIRI(value, path, hashAlgorithmContext)
        if err != nil {
            return false, err
        }
        err = obj.Algorithm().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Core/hashValue", "hashValue":
        val, err := DecodeString(value, path, hashHashValueContext)
        if err != nil {
            return false, err
        }
        err = obj.HashValue().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self HashObjectType) Create() SHACLObject {
    return ConstructHashObject(&HashObject{}, self)
}

func ConstructHashObject(o *HashObject, typ SHACLType) *HashObject {
    ConstructIntegrityMethodObject(&o.IntegrityMethodObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/adler32",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake2b512",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/blake3",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsDilithium",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/crystalsKyber",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/falcon",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md2",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md4",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md5",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/md6",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha1",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha224",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_224",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_256",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_384",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha3_512",
                "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm/sha512",
        }})
        o.algorithm = NewProperty[string]("algorithm", validators)
    }
    {
        validators := []Validator[string]{}
        o.hashValue = NewProperty[string]("hashValue", validators)
    }
    return o
}

type Hash interface {
    IntegrityMethod
    Algorithm() PropertyInterface[string]
    HashValue() PropertyInterface[string]
}


func MakeHash() Hash {
    return ConstructHashObject(&HashObject{}, hashType)
}

func MakeHashRef() Ref[Hash] {
    o := MakeHash()
    return MakeObjectRef[Hash](o)
}

func (self *HashObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.IntegrityMethodObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("algorithm")
        if ! self.algorithm.Check(prop_path, handler) {
            valid = false
        }
        if ! self.algorithm.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"algorithm", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("hashValue")
        if ! self.hashValue.Check(prop_path, handler) {
            valid = false
        }
        if ! self.hashValue.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"hashValue", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *HashObject) Walk(path Path, visit Visit) {
    self.IntegrityMethodObject.Walk(path, visit)
    self.algorithm.Walk(path, visit)
    self.hashValue.Walk(path, visit)
}


func (self *HashObject) Algorithm() PropertyInterface[string] { return &self.algorithm }
func (self *HashObject) HashValue() PropertyInterface[string] { return &self.hashValue }

func (self *HashObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.IntegrityMethodObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.algorithm.IsSet() {
        val, err := EncodeIRI(self.algorithm.Get(), path.PushPath("algorithm"), hashAlgorithmContext, state)
        if err != nil {
            return err
        }
        data["algorithm"] = val
    }
    if self.hashValue.IsSet() {
        val, err := EncodeString(self.hashValue.Get(), path.PushPath("hashValue"), hashHashValueContext, state)
        if err != nil {
            return err
        }
        data["hashValue"] = val
    }
    return nil
}

// Provide context for a relationship that occurs in the lifecycle.
type LifecycleScopedRelationshipObject struct {
    RelationshipObject

    // Capture the scope of information about a specific relationship between elements.
    scope Property[string]
}


type LifecycleScopedRelationshipObjectType struct {
    SHACLTypeBase
}
var lifecycleScopedRelationshipType LifecycleScopedRelationshipObjectType
var lifecycleScopedRelationshipScopeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/build": "build",
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/design": "design",
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/development": "development",
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/runtime": "runtime",
    "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/test": "test",}

func DecodeLifecycleScopedRelationship (data any, path Path, context map[string]string) (Ref[LifecycleScopedRelationship], error) {
    return DecodeRef[LifecycleScopedRelationship](data, path, context, lifecycleScopedRelationshipType)
}

func (self LifecycleScopedRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(LifecycleScopedRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/scope", "scope":
        val, err := DecodeIRI(value, path, lifecycleScopedRelationshipScopeContext)
        if err != nil {
            return false, err
        }
        err = obj.Scope().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self LifecycleScopedRelationshipObjectType) Create() SHACLObject {
    return ConstructLifecycleScopedRelationshipObject(&LifecycleScopedRelationshipObject{}, self)
}

func ConstructLifecycleScopedRelationshipObject(o *LifecycleScopedRelationshipObject, typ SHACLType) *LifecycleScopedRelationshipObject {
    ConstructRelationshipObject(&o.RelationshipObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/build",
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/design",
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/development",
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/other",
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/runtime",
                "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType/test",
        }})
        o.scope = NewProperty[string]("scope", validators)
    }
    return o
}

type LifecycleScopedRelationship interface {
    Relationship
    Scope() PropertyInterface[string]
}


func MakeLifecycleScopedRelationship() LifecycleScopedRelationship {
    return ConstructLifecycleScopedRelationshipObject(&LifecycleScopedRelationshipObject{}, lifecycleScopedRelationshipType)
}

func MakeLifecycleScopedRelationshipRef() Ref[LifecycleScopedRelationship] {
    o := MakeLifecycleScopedRelationship()
    return MakeObjectRef[LifecycleScopedRelationship](o)
}

func (self *LifecycleScopedRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.RelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("scope")
        if ! self.scope.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *LifecycleScopedRelationshipObject) Walk(path Path, visit Visit) {
    self.RelationshipObject.Walk(path, visit)
    self.scope.Walk(path, visit)
}


func (self *LifecycleScopedRelationshipObject) Scope() PropertyInterface[string] { return &self.scope }

func (self *LifecycleScopedRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.RelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.scope.IsSet() {
        val, err := EncodeIRI(self.scope.Get(), path.PushPath("scope"), lifecycleScopedRelationshipScopeContext, state)
        if err != nil {
            return err
        }
        data["scope"] = val
    }
    return nil
}

// A group of people who work together in an organized way for a shared purpose.
type OrganizationObject struct {
    AgentObject

}

// An Organization representing the SPDX Project.
const OrganizationSpdxOrganization = "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization"

type OrganizationObjectType struct {
    SHACLTypeBase
}
var organizationType OrganizationObjectType

func DecodeOrganization (data any, path Path, context map[string]string) (Ref[Organization], error) {
    return DecodeRef[Organization](data, path, context, organizationType)
}

func (self OrganizationObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Organization)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self OrganizationObjectType) Create() SHACLObject {
    return ConstructOrganizationObject(&OrganizationObject{}, self)
}

func ConstructOrganizationObject(o *OrganizationObject, typ SHACLType) *OrganizationObject {
    ConstructAgentObject(&o.AgentObject, typ)
    return o
}

type Organization interface {
    Agent
}


func MakeOrganization() Organization {
    return ConstructOrganizationObject(&OrganizationObject{}, organizationType)
}

func MakeOrganizationRef() Ref[Organization] {
    o := MakeOrganization()
    return MakeObjectRef[Organization](o)
}

func (self *OrganizationObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.AgentObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *OrganizationObject) Walk(path Path, visit Visit) {
    self.AgentObject.Walk(path, visit)
}



func (self *OrganizationObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.AgentObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// An individual human being.
type PersonObject struct {
    AgentObject

}


type PersonObjectType struct {
    SHACLTypeBase
}
var personType PersonObjectType

func DecodePerson (data any, path Path, context map[string]string) (Ref[Person], error) {
    return DecodeRef[Person](data, path, context, personType)
}

func (self PersonObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Person)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self PersonObjectType) Create() SHACLObject {
    return ConstructPersonObject(&PersonObject{}, self)
}

func ConstructPersonObject(o *PersonObject, typ SHACLType) *PersonObject {
    ConstructAgentObject(&o.AgentObject, typ)
    return o
}

type Person interface {
    Agent
}


func MakePerson() Person {
    return ConstructPersonObject(&PersonObject{}, personType)
}

func MakePersonRef() Ref[Person] {
    o := MakePerson()
    return MakeObjectRef[Person](o)
}

func (self *PersonObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.AgentObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *PersonObject) Walk(path Path, visit Visit) {
    self.AgentObject.Walk(path, visit)
}



func (self *PersonObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.AgentObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A software agent.
type SoftwareAgentObject struct {
    AgentObject

}


type SoftwareAgentObjectType struct {
    SHACLTypeBase
}
var softwareAgentType SoftwareAgentObjectType

func DecodeSoftwareAgent (data any, path Path, context map[string]string) (Ref[SoftwareAgent], error) {
    return DecodeRef[SoftwareAgent](data, path, context, softwareAgentType)
}

func (self SoftwareAgentObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareAgent)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareAgentObjectType) Create() SHACLObject {
    return ConstructSoftwareAgentObject(&SoftwareAgentObject{}, self)
}

func ConstructSoftwareAgentObject(o *SoftwareAgentObject, typ SHACLType) *SoftwareAgentObject {
    ConstructAgentObject(&o.AgentObject, typ)
    return o
}

type SoftwareAgent interface {
    Agent
}


func MakeSoftwareAgent() SoftwareAgent {
    return ConstructSoftwareAgentObject(&SoftwareAgentObject{}, softwareAgentType)
}

func MakeSoftwareAgentRef() Ref[SoftwareAgent] {
    o := MakeSoftwareAgent()
    return MakeObjectRef[SoftwareAgent](o)
}

func (self *SoftwareAgentObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.AgentObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SoftwareAgentObject) Walk(path Path, visit Visit) {
    self.AgentObject.Walk(path, visit)
}



func (self *SoftwareAgentObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.AgentObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Portion of an AnyLicenseInfo representing a set of licensing information

// where all elements apply.
type ExpandedlicensingConjunctiveLicenseSetObject struct {
    SimplelicensingAnyLicenseInfoObject

    // A license expression participating in a license set.
    expandedlicensingMember RefListProperty[SimplelicensingAnyLicenseInfo]
}


type ExpandedlicensingConjunctiveLicenseSetObjectType struct {
    SHACLTypeBase
}
var expandedlicensingConjunctiveLicenseSetType ExpandedlicensingConjunctiveLicenseSetObjectType
var expandedlicensingConjunctiveLicenseSetExpandedlicensingMemberContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}

func DecodeExpandedlicensingConjunctiveLicenseSet (data any, path Path, context map[string]string) (Ref[ExpandedlicensingConjunctiveLicenseSet], error) {
    return DecodeRef[ExpandedlicensingConjunctiveLicenseSet](data, path, context, expandedlicensingConjunctiveLicenseSetType)
}

func (self ExpandedlicensingConjunctiveLicenseSetObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingConjunctiveLicenseSet)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/member", "expandedlicensing_member":
        val, err := DecodeList[Ref[SimplelicensingAnyLicenseInfo]](value, path, expandedlicensingConjunctiveLicenseSetExpandedlicensingMemberContext, DecodeSimplelicensingAnyLicenseInfo)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingMember().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingConjunctiveLicenseSetObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingConjunctiveLicenseSetObject(&ExpandedlicensingConjunctiveLicenseSetObject{}, self)
}

func ConstructExpandedlicensingConjunctiveLicenseSetObject(o *ExpandedlicensingConjunctiveLicenseSetObject, typ SHACLType) *ExpandedlicensingConjunctiveLicenseSetObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    {
        validators := []Validator[Ref[SimplelicensingAnyLicenseInfo]]{}
        o.expandedlicensingMember = NewRefListProperty[SimplelicensingAnyLicenseInfo]("expandedlicensingMember", validators)
    }
    return o
}

type ExpandedlicensingConjunctiveLicenseSet interface {
    SimplelicensingAnyLicenseInfo
    ExpandedlicensingMember() ListPropertyInterface[Ref[SimplelicensingAnyLicenseInfo]]
}


func MakeExpandedlicensingConjunctiveLicenseSet() ExpandedlicensingConjunctiveLicenseSet {
    return ConstructExpandedlicensingConjunctiveLicenseSetObject(&ExpandedlicensingConjunctiveLicenseSetObject{}, expandedlicensingConjunctiveLicenseSetType)
}

func MakeExpandedlicensingConjunctiveLicenseSetRef() Ref[ExpandedlicensingConjunctiveLicenseSet] {
    o := MakeExpandedlicensingConjunctiveLicenseSet()
    return MakeObjectRef[ExpandedlicensingConjunctiveLicenseSet](o)
}

func (self *ExpandedlicensingConjunctiveLicenseSetObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingMember")
        if ! self.expandedlicensingMember.Check(prop_path, handler) {
            valid = false
        }
        if len(self.expandedlicensingMember.Get()) < 2 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "expandedlicensingMember",
                    "Too few elements. Minimum of 2 required"},
                    prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingConjunctiveLicenseSetObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
    self.expandedlicensingMember.Walk(path, visit)
}


func (self *ExpandedlicensingConjunctiveLicenseSetObject) ExpandedlicensingMember() ListPropertyInterface[Ref[SimplelicensingAnyLicenseInfo]] { return &self.expandedlicensingMember }

func (self *ExpandedlicensingConjunctiveLicenseSetObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingMember.IsSet() {
        val, err := EncodeList[Ref[SimplelicensingAnyLicenseInfo]](self.expandedlicensingMember.Get(), path.PushPath("expandedlicensingMember"), expandedlicensingConjunctiveLicenseSetExpandedlicensingMemberContext, state, EncodeRef[SimplelicensingAnyLicenseInfo])
        if err != nil {
            return err
        }
        data["expandedlicensing_member"] = val
    }
    return nil
}

// A license addition that is not listed on the SPDX Exceptions List.
type ExpandedlicensingCustomLicenseAdditionObject struct {
    ExpandedlicensingLicenseAdditionObject

}


type ExpandedlicensingCustomLicenseAdditionObjectType struct {
    SHACLTypeBase
}
var expandedlicensingCustomLicenseAdditionType ExpandedlicensingCustomLicenseAdditionObjectType

func DecodeExpandedlicensingCustomLicenseAddition (data any, path Path, context map[string]string) (Ref[ExpandedlicensingCustomLicenseAddition], error) {
    return DecodeRef[ExpandedlicensingCustomLicenseAddition](data, path, context, expandedlicensingCustomLicenseAdditionType)
}

func (self ExpandedlicensingCustomLicenseAdditionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingCustomLicenseAddition)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingCustomLicenseAdditionObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingCustomLicenseAdditionObject(&ExpandedlicensingCustomLicenseAdditionObject{}, self)
}

func ConstructExpandedlicensingCustomLicenseAdditionObject(o *ExpandedlicensingCustomLicenseAdditionObject, typ SHACLType) *ExpandedlicensingCustomLicenseAdditionObject {
    ConstructExpandedlicensingLicenseAdditionObject(&o.ExpandedlicensingLicenseAdditionObject, typ)
    return o
}

type ExpandedlicensingCustomLicenseAddition interface {
    ExpandedlicensingLicenseAddition
}


func MakeExpandedlicensingCustomLicenseAddition() ExpandedlicensingCustomLicenseAddition {
    return ConstructExpandedlicensingCustomLicenseAdditionObject(&ExpandedlicensingCustomLicenseAdditionObject{}, expandedlicensingCustomLicenseAdditionType)
}

func MakeExpandedlicensingCustomLicenseAdditionRef() Ref[ExpandedlicensingCustomLicenseAddition] {
    o := MakeExpandedlicensingCustomLicenseAddition()
    return MakeObjectRef[ExpandedlicensingCustomLicenseAddition](o)
}

func (self *ExpandedlicensingCustomLicenseAdditionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingLicenseAdditionObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExpandedlicensingCustomLicenseAdditionObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingLicenseAdditionObject.Walk(path, visit)
}



func (self *ExpandedlicensingCustomLicenseAdditionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingLicenseAdditionObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Portion of an AnyLicenseInfo representing a set of licensing information where

// only one of the elements applies.
type ExpandedlicensingDisjunctiveLicenseSetObject struct {
    SimplelicensingAnyLicenseInfoObject

    // A license expression participating in a license set.
    expandedlicensingMember RefListProperty[SimplelicensingAnyLicenseInfo]
}


type ExpandedlicensingDisjunctiveLicenseSetObjectType struct {
    SHACLTypeBase
}
var expandedlicensingDisjunctiveLicenseSetType ExpandedlicensingDisjunctiveLicenseSetObjectType
var expandedlicensingDisjunctiveLicenseSetExpandedlicensingMemberContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense": "expandedlicensing_NoneLicense",
    "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense": "expandedlicensing_NoAssertionLicense",}

func DecodeExpandedlicensingDisjunctiveLicenseSet (data any, path Path, context map[string]string) (Ref[ExpandedlicensingDisjunctiveLicenseSet], error) {
    return DecodeRef[ExpandedlicensingDisjunctiveLicenseSet](data, path, context, expandedlicensingDisjunctiveLicenseSetType)
}

func (self ExpandedlicensingDisjunctiveLicenseSetObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingDisjunctiveLicenseSet)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/member", "expandedlicensing_member":
        val, err := DecodeList[Ref[SimplelicensingAnyLicenseInfo]](value, path, expandedlicensingDisjunctiveLicenseSetExpandedlicensingMemberContext, DecodeSimplelicensingAnyLicenseInfo)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingMember().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingDisjunctiveLicenseSetObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingDisjunctiveLicenseSetObject(&ExpandedlicensingDisjunctiveLicenseSetObject{}, self)
}

func ConstructExpandedlicensingDisjunctiveLicenseSetObject(o *ExpandedlicensingDisjunctiveLicenseSetObject, typ SHACLType) *ExpandedlicensingDisjunctiveLicenseSetObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    {
        validators := []Validator[Ref[SimplelicensingAnyLicenseInfo]]{}
        o.expandedlicensingMember = NewRefListProperty[SimplelicensingAnyLicenseInfo]("expandedlicensingMember", validators)
    }
    return o
}

type ExpandedlicensingDisjunctiveLicenseSet interface {
    SimplelicensingAnyLicenseInfo
    ExpandedlicensingMember() ListPropertyInterface[Ref[SimplelicensingAnyLicenseInfo]]
}


func MakeExpandedlicensingDisjunctiveLicenseSet() ExpandedlicensingDisjunctiveLicenseSet {
    return ConstructExpandedlicensingDisjunctiveLicenseSetObject(&ExpandedlicensingDisjunctiveLicenseSetObject{}, expandedlicensingDisjunctiveLicenseSetType)
}

func MakeExpandedlicensingDisjunctiveLicenseSetRef() Ref[ExpandedlicensingDisjunctiveLicenseSet] {
    o := MakeExpandedlicensingDisjunctiveLicenseSet()
    return MakeObjectRef[ExpandedlicensingDisjunctiveLicenseSet](o)
}

func (self *ExpandedlicensingDisjunctiveLicenseSetObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingMember")
        if ! self.expandedlicensingMember.Check(prop_path, handler) {
            valid = false
        }
        if len(self.expandedlicensingMember.Get()) < 2 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "expandedlicensingMember",
                    "Too few elements. Minimum of 2 required"},
                    prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingDisjunctiveLicenseSetObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
    self.expandedlicensingMember.Walk(path, visit)
}


func (self *ExpandedlicensingDisjunctiveLicenseSetObject) ExpandedlicensingMember() ListPropertyInterface[Ref[SimplelicensingAnyLicenseInfo]] { return &self.expandedlicensingMember }

func (self *ExpandedlicensingDisjunctiveLicenseSetObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingMember.IsSet() {
        val, err := EncodeList[Ref[SimplelicensingAnyLicenseInfo]](self.expandedlicensingMember.Get(), path.PushPath("expandedlicensingMember"), expandedlicensingDisjunctiveLicenseSetExpandedlicensingMemberContext, state, EncodeRef[SimplelicensingAnyLicenseInfo])
        if err != nil {
            return err
        }
        data["expandedlicensing_member"] = val
    }
    return nil
}

// Abstract class representing a License or an OrLaterOperator.
type ExpandedlicensingExtendableLicenseObject struct {
    SimplelicensingAnyLicenseInfoObject

}


type ExpandedlicensingExtendableLicenseObjectType struct {
    SHACLTypeBase
}
var expandedlicensingExtendableLicenseType ExpandedlicensingExtendableLicenseObjectType

func DecodeExpandedlicensingExtendableLicense (data any, path Path, context map[string]string) (Ref[ExpandedlicensingExtendableLicense], error) {
    return DecodeRef[ExpandedlicensingExtendableLicense](data, path, context, expandedlicensingExtendableLicenseType)
}

func (self ExpandedlicensingExtendableLicenseObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingExtendableLicense)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingExtendableLicenseObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingExtendableLicenseObject(&ExpandedlicensingExtendableLicenseObject{}, self)
}

func ConstructExpandedlicensingExtendableLicenseObject(o *ExpandedlicensingExtendableLicenseObject, typ SHACLType) *ExpandedlicensingExtendableLicenseObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    return o
}

type ExpandedlicensingExtendableLicense interface {
    SimplelicensingAnyLicenseInfo
}



func (self *ExpandedlicensingExtendableLicenseObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExpandedlicensingExtendableLicenseObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
}



func (self *ExpandedlicensingExtendableLicenseObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A concrete subclass of AnyLicenseInfo used by Individuals in the

// ExpandedLicensing profile.
type ExpandedlicensingIndividualLicensingInfoObject struct {
    SimplelicensingAnyLicenseInfoObject

}

// An Individual Value for License when no assertion can be made about its actual

// value.
const ExpandedlicensingIndividualLicensingInfoNoAssertionLicense = "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoAssertionLicense"
// An Individual Value for License where the SPDX data creator determines that no

// license is present.
const ExpandedlicensingIndividualLicensingInfoNoneLicense = "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/NoneLicense"

type ExpandedlicensingIndividualLicensingInfoObjectType struct {
    SHACLTypeBase
}
var expandedlicensingIndividualLicensingInfoType ExpandedlicensingIndividualLicensingInfoObjectType

func DecodeExpandedlicensingIndividualLicensingInfo (data any, path Path, context map[string]string) (Ref[ExpandedlicensingIndividualLicensingInfo], error) {
    return DecodeRef[ExpandedlicensingIndividualLicensingInfo](data, path, context, expandedlicensingIndividualLicensingInfoType)
}

func (self ExpandedlicensingIndividualLicensingInfoObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingIndividualLicensingInfo)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingIndividualLicensingInfoObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingIndividualLicensingInfoObject(&ExpandedlicensingIndividualLicensingInfoObject{}, self)
}

func ConstructExpandedlicensingIndividualLicensingInfoObject(o *ExpandedlicensingIndividualLicensingInfoObject, typ SHACLType) *ExpandedlicensingIndividualLicensingInfoObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    return o
}

type ExpandedlicensingIndividualLicensingInfo interface {
    SimplelicensingAnyLicenseInfo
}


func MakeExpandedlicensingIndividualLicensingInfo() ExpandedlicensingIndividualLicensingInfo {
    return ConstructExpandedlicensingIndividualLicensingInfoObject(&ExpandedlicensingIndividualLicensingInfoObject{}, expandedlicensingIndividualLicensingInfoType)
}

func MakeExpandedlicensingIndividualLicensingInfoRef() Ref[ExpandedlicensingIndividualLicensingInfo] {
    o := MakeExpandedlicensingIndividualLicensingInfo()
    return MakeObjectRef[ExpandedlicensingIndividualLicensingInfo](o)
}

func (self *ExpandedlicensingIndividualLicensingInfoObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExpandedlicensingIndividualLicensingInfoObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
}



func (self *ExpandedlicensingIndividualLicensingInfoObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Abstract class for the portion of an AnyLicenseInfo representing a license.
type ExpandedlicensingLicenseObject struct {
    ExpandedlicensingExtendableLicenseObject

    // Specifies whether a license or additional text identifier has been marked as
    // deprecated.
    expandedlicensingIsDeprecatedLicenseId Property[bool]
    // Specifies whether the License is listed as free by the
    // Free Software Foundation (FSF).
    expandedlicensingIsFsfLibre Property[bool]
    // Specifies whether the License is listed as approved by the
    // Open Source Initiative (OSI).
    expandedlicensingIsOsiApproved Property[bool]
    // Identifies all the text and metadata associated with a license in the license
    // XML format.
    expandedlicensingLicenseXml Property[string]
    // Specifies the licenseId that is preferred to be used in place of a deprecated
    // License or LicenseAddition.
    expandedlicensingObsoletedBy Property[string]
    // Contains a URL where the License or LicenseAddition can be found in use.
    expandedlicensingSeeAlso ListProperty[string]
    // Provides a License author's preferred text to indicate that a file is covered
    // by the License.
    expandedlicensingStandardLicenseHeader Property[string]
    // Identifies the full text of a License, in SPDX templating format.
    expandedlicensingStandardLicenseTemplate Property[string]
    // Identifies the full text of a License or Addition.
    simplelicensingLicenseText Property[string]
}


type ExpandedlicensingLicenseObjectType struct {
    SHACLTypeBase
}
var expandedlicensingLicenseType ExpandedlicensingLicenseObjectType
var expandedlicensingLicenseExpandedlicensingIsDeprecatedLicenseIdContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingIsFsfLibreContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingIsOsiApprovedContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingLicenseXmlContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingObsoletedByContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingSeeAlsoContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingStandardLicenseHeaderContext = map[string]string{}
var expandedlicensingLicenseExpandedlicensingStandardLicenseTemplateContext = map[string]string{}
var expandedlicensingLicenseSimplelicensingLicenseTextContext = map[string]string{}

func DecodeExpandedlicensingLicense (data any, path Path, context map[string]string) (Ref[ExpandedlicensingLicense], error) {
    return DecodeRef[ExpandedlicensingLicense](data, path, context, expandedlicensingLicenseType)
}

func (self ExpandedlicensingLicenseObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingLicense)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/isDeprecatedLicenseId", "expandedlicensing_isDeprecatedLicenseId":
        val, err := DecodeBoolean(value, path, expandedlicensingLicenseExpandedlicensingIsDeprecatedLicenseIdContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingIsDeprecatedLicenseId().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/isFsfLibre", "expandedlicensing_isFsfLibre":
        val, err := DecodeBoolean(value, path, expandedlicensingLicenseExpandedlicensingIsFsfLibreContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingIsFsfLibre().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/isOsiApproved", "expandedlicensing_isOsiApproved":
        val, err := DecodeBoolean(value, path, expandedlicensingLicenseExpandedlicensingIsOsiApprovedContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingIsOsiApproved().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/licenseXml", "expandedlicensing_licenseXml":
        val, err := DecodeString(value, path, expandedlicensingLicenseExpandedlicensingLicenseXmlContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingLicenseXml().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/obsoletedBy", "expandedlicensing_obsoletedBy":
        val, err := DecodeString(value, path, expandedlicensingLicenseExpandedlicensingObsoletedByContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingObsoletedBy().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/seeAlso", "expandedlicensing_seeAlso":
        val, err := DecodeList[string](value, path, expandedlicensingLicenseExpandedlicensingSeeAlsoContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingSeeAlso().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/standardLicenseHeader", "expandedlicensing_standardLicenseHeader":
        val, err := DecodeString(value, path, expandedlicensingLicenseExpandedlicensingStandardLicenseHeaderContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingStandardLicenseHeader().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/standardLicenseTemplate", "expandedlicensing_standardLicenseTemplate":
        val, err := DecodeString(value, path, expandedlicensingLicenseExpandedlicensingStandardLicenseTemplateContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingStandardLicenseTemplate().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/licenseText", "simplelicensing_licenseText":
        val, err := DecodeString(value, path, expandedlicensingLicenseSimplelicensingLicenseTextContext)
        if err != nil {
            return false, err
        }
        err = obj.SimplelicensingLicenseText().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingLicenseObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingLicenseObject(&ExpandedlicensingLicenseObject{}, self)
}

func ConstructExpandedlicensingLicenseObject(o *ExpandedlicensingLicenseObject, typ SHACLType) *ExpandedlicensingLicenseObject {
    ConstructExpandedlicensingExtendableLicenseObject(&o.ExpandedlicensingExtendableLicenseObject, typ)
    {
        validators := []Validator[bool]{}
        o.expandedlicensingIsDeprecatedLicenseId = NewProperty[bool]("expandedlicensingIsDeprecatedLicenseId", validators)
    }
    {
        validators := []Validator[bool]{}
        o.expandedlicensingIsFsfLibre = NewProperty[bool]("expandedlicensingIsFsfLibre", validators)
    }
    {
        validators := []Validator[bool]{}
        o.expandedlicensingIsOsiApproved = NewProperty[bool]("expandedlicensingIsOsiApproved", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingLicenseXml = NewProperty[string]("expandedlicensingLicenseXml", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingObsoletedBy = NewProperty[string]("expandedlicensingObsoletedBy", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingSeeAlso = NewListProperty[string]("expandedlicensingSeeAlso", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingStandardLicenseHeader = NewProperty[string]("expandedlicensingStandardLicenseHeader", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingStandardLicenseTemplate = NewProperty[string]("expandedlicensingStandardLicenseTemplate", validators)
    }
    {
        validators := []Validator[string]{}
        o.simplelicensingLicenseText = NewProperty[string]("simplelicensingLicenseText", validators)
    }
    return o
}

type ExpandedlicensingLicense interface {
    ExpandedlicensingExtendableLicense
    ExpandedlicensingIsDeprecatedLicenseId() PropertyInterface[bool]
    ExpandedlicensingIsFsfLibre() PropertyInterface[bool]
    ExpandedlicensingIsOsiApproved() PropertyInterface[bool]
    ExpandedlicensingLicenseXml() PropertyInterface[string]
    ExpandedlicensingObsoletedBy() PropertyInterface[string]
    ExpandedlicensingSeeAlso() ListPropertyInterface[string]
    ExpandedlicensingStandardLicenseHeader() PropertyInterface[string]
    ExpandedlicensingStandardLicenseTemplate() PropertyInterface[string]
    SimplelicensingLicenseText() PropertyInterface[string]
}



func (self *ExpandedlicensingLicenseObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingExtendableLicenseObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingIsDeprecatedLicenseId")
        if ! self.expandedlicensingIsDeprecatedLicenseId.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingIsFsfLibre")
        if ! self.expandedlicensingIsFsfLibre.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingIsOsiApproved")
        if ! self.expandedlicensingIsOsiApproved.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingLicenseXml")
        if ! self.expandedlicensingLicenseXml.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingObsoletedBy")
        if ! self.expandedlicensingObsoletedBy.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingSeeAlso")
        if ! self.expandedlicensingSeeAlso.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingStandardLicenseHeader")
        if ! self.expandedlicensingStandardLicenseHeader.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingStandardLicenseTemplate")
        if ! self.expandedlicensingStandardLicenseTemplate.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("simplelicensingLicenseText")
        if ! self.simplelicensingLicenseText.Check(prop_path, handler) {
            valid = false
        }
        if ! self.simplelicensingLicenseText.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"simplelicensingLicenseText", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingLicenseObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingExtendableLicenseObject.Walk(path, visit)
    self.expandedlicensingIsDeprecatedLicenseId.Walk(path, visit)
    self.expandedlicensingIsFsfLibre.Walk(path, visit)
    self.expandedlicensingIsOsiApproved.Walk(path, visit)
    self.expandedlicensingLicenseXml.Walk(path, visit)
    self.expandedlicensingObsoletedBy.Walk(path, visit)
    self.expandedlicensingSeeAlso.Walk(path, visit)
    self.expandedlicensingStandardLicenseHeader.Walk(path, visit)
    self.expandedlicensingStandardLicenseTemplate.Walk(path, visit)
    self.simplelicensingLicenseText.Walk(path, visit)
}


func (self *ExpandedlicensingLicenseObject) ExpandedlicensingIsDeprecatedLicenseId() PropertyInterface[bool] { return &self.expandedlicensingIsDeprecatedLicenseId }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingIsFsfLibre() PropertyInterface[bool] { return &self.expandedlicensingIsFsfLibre }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingIsOsiApproved() PropertyInterface[bool] { return &self.expandedlicensingIsOsiApproved }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingLicenseXml() PropertyInterface[string] { return &self.expandedlicensingLicenseXml }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingObsoletedBy() PropertyInterface[string] { return &self.expandedlicensingObsoletedBy }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingSeeAlso() ListPropertyInterface[string] { return &self.expandedlicensingSeeAlso }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingStandardLicenseHeader() PropertyInterface[string] { return &self.expandedlicensingStandardLicenseHeader }
func (self *ExpandedlicensingLicenseObject) ExpandedlicensingStandardLicenseTemplate() PropertyInterface[string] { return &self.expandedlicensingStandardLicenseTemplate }
func (self *ExpandedlicensingLicenseObject) SimplelicensingLicenseText() PropertyInterface[string] { return &self.simplelicensingLicenseText }

func (self *ExpandedlicensingLicenseObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingExtendableLicenseObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingIsDeprecatedLicenseId.IsSet() {
        val, err := EncodeBoolean(self.expandedlicensingIsDeprecatedLicenseId.Get(), path.PushPath("expandedlicensingIsDeprecatedLicenseId"), expandedlicensingLicenseExpandedlicensingIsDeprecatedLicenseIdContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_isDeprecatedLicenseId"] = val
    }
    if self.expandedlicensingIsFsfLibre.IsSet() {
        val, err := EncodeBoolean(self.expandedlicensingIsFsfLibre.Get(), path.PushPath("expandedlicensingIsFsfLibre"), expandedlicensingLicenseExpandedlicensingIsFsfLibreContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_isFsfLibre"] = val
    }
    if self.expandedlicensingIsOsiApproved.IsSet() {
        val, err := EncodeBoolean(self.expandedlicensingIsOsiApproved.Get(), path.PushPath("expandedlicensingIsOsiApproved"), expandedlicensingLicenseExpandedlicensingIsOsiApprovedContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_isOsiApproved"] = val
    }
    if self.expandedlicensingLicenseXml.IsSet() {
        val, err := EncodeString(self.expandedlicensingLicenseXml.Get(), path.PushPath("expandedlicensingLicenseXml"), expandedlicensingLicenseExpandedlicensingLicenseXmlContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_licenseXml"] = val
    }
    if self.expandedlicensingObsoletedBy.IsSet() {
        val, err := EncodeString(self.expandedlicensingObsoletedBy.Get(), path.PushPath("expandedlicensingObsoletedBy"), expandedlicensingLicenseExpandedlicensingObsoletedByContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_obsoletedBy"] = val
    }
    if self.expandedlicensingSeeAlso.IsSet() {
        val, err := EncodeList[string](self.expandedlicensingSeeAlso.Get(), path.PushPath("expandedlicensingSeeAlso"), expandedlicensingLicenseExpandedlicensingSeeAlsoContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["expandedlicensing_seeAlso"] = val
    }
    if self.expandedlicensingStandardLicenseHeader.IsSet() {
        val, err := EncodeString(self.expandedlicensingStandardLicenseHeader.Get(), path.PushPath("expandedlicensingStandardLicenseHeader"), expandedlicensingLicenseExpandedlicensingStandardLicenseHeaderContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_standardLicenseHeader"] = val
    }
    if self.expandedlicensingStandardLicenseTemplate.IsSet() {
        val, err := EncodeString(self.expandedlicensingStandardLicenseTemplate.Get(), path.PushPath("expandedlicensingStandardLicenseTemplate"), expandedlicensingLicenseExpandedlicensingStandardLicenseTemplateContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_standardLicenseTemplate"] = val
    }
    if self.simplelicensingLicenseText.IsSet() {
        val, err := EncodeString(self.simplelicensingLicenseText.Get(), path.PushPath("simplelicensingLicenseText"), expandedlicensingLicenseSimplelicensingLicenseTextContext, state)
        if err != nil {
            return err
        }
        data["simplelicensing_licenseText"] = val
    }
    return nil
}

// A license that is listed on the SPDX License List.
type ExpandedlicensingListedLicenseObject struct {
    ExpandedlicensingLicenseObject

    // Specifies the SPDX License List version in which this license or exception
    // identifier was deprecated.
    expandedlicensingDeprecatedVersion Property[string]
    // Specifies the SPDX License List version in which this ListedLicense or
    // ListedLicenseException identifier was first added.
    expandedlicensingListVersionAdded Property[string]
}


type ExpandedlicensingListedLicenseObjectType struct {
    SHACLTypeBase
}
var expandedlicensingListedLicenseType ExpandedlicensingListedLicenseObjectType
var expandedlicensingListedLicenseExpandedlicensingDeprecatedVersionContext = map[string]string{}
var expandedlicensingListedLicenseExpandedlicensingListVersionAddedContext = map[string]string{}

func DecodeExpandedlicensingListedLicense (data any, path Path, context map[string]string) (Ref[ExpandedlicensingListedLicense], error) {
    return DecodeRef[ExpandedlicensingListedLicense](data, path, context, expandedlicensingListedLicenseType)
}

func (self ExpandedlicensingListedLicenseObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingListedLicense)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/deprecatedVersion", "expandedlicensing_deprecatedVersion":
        val, err := DecodeString(value, path, expandedlicensingListedLicenseExpandedlicensingDeprecatedVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingDeprecatedVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/listVersionAdded", "expandedlicensing_listVersionAdded":
        val, err := DecodeString(value, path, expandedlicensingListedLicenseExpandedlicensingListVersionAddedContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingListVersionAdded().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingListedLicenseObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingListedLicenseObject(&ExpandedlicensingListedLicenseObject{}, self)
}

func ConstructExpandedlicensingListedLicenseObject(o *ExpandedlicensingListedLicenseObject, typ SHACLType) *ExpandedlicensingListedLicenseObject {
    ConstructExpandedlicensingLicenseObject(&o.ExpandedlicensingLicenseObject, typ)
    {
        validators := []Validator[string]{}
        o.expandedlicensingDeprecatedVersion = NewProperty[string]("expandedlicensingDeprecatedVersion", validators)
    }
    {
        validators := []Validator[string]{}
        o.expandedlicensingListVersionAdded = NewProperty[string]("expandedlicensingListVersionAdded", validators)
    }
    return o
}

type ExpandedlicensingListedLicense interface {
    ExpandedlicensingLicense
    ExpandedlicensingDeprecatedVersion() PropertyInterface[string]
    ExpandedlicensingListVersionAdded() PropertyInterface[string]
}


func MakeExpandedlicensingListedLicense() ExpandedlicensingListedLicense {
    return ConstructExpandedlicensingListedLicenseObject(&ExpandedlicensingListedLicenseObject{}, expandedlicensingListedLicenseType)
}

func MakeExpandedlicensingListedLicenseRef() Ref[ExpandedlicensingListedLicense] {
    o := MakeExpandedlicensingListedLicense()
    return MakeObjectRef[ExpandedlicensingListedLicense](o)
}

func (self *ExpandedlicensingListedLicenseObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingLicenseObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingDeprecatedVersion")
        if ! self.expandedlicensingDeprecatedVersion.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingListVersionAdded")
        if ! self.expandedlicensingListVersionAdded.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingListedLicenseObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingLicenseObject.Walk(path, visit)
    self.expandedlicensingDeprecatedVersion.Walk(path, visit)
    self.expandedlicensingListVersionAdded.Walk(path, visit)
}


func (self *ExpandedlicensingListedLicenseObject) ExpandedlicensingDeprecatedVersion() PropertyInterface[string] { return &self.expandedlicensingDeprecatedVersion }
func (self *ExpandedlicensingListedLicenseObject) ExpandedlicensingListVersionAdded() PropertyInterface[string] { return &self.expandedlicensingListVersionAdded }

func (self *ExpandedlicensingListedLicenseObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingLicenseObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingDeprecatedVersion.IsSet() {
        val, err := EncodeString(self.expandedlicensingDeprecatedVersion.Get(), path.PushPath("expandedlicensingDeprecatedVersion"), expandedlicensingListedLicenseExpandedlicensingDeprecatedVersionContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_deprecatedVersion"] = val
    }
    if self.expandedlicensingListVersionAdded.IsSet() {
        val, err := EncodeString(self.expandedlicensingListVersionAdded.Get(), path.PushPath("expandedlicensingListVersionAdded"), expandedlicensingListedLicenseExpandedlicensingListVersionAddedContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_listVersionAdded"] = val
    }
    return nil
}

// Portion of an AnyLicenseInfo representing this version, or any later version,

// of the indicated License.
type ExpandedlicensingOrLaterOperatorObject struct {
    ExpandedlicensingExtendableLicenseObject

    // A License participating in an 'or later' model.
    expandedlicensingSubjectLicense RefProperty[ExpandedlicensingLicense]
}


type ExpandedlicensingOrLaterOperatorObjectType struct {
    SHACLTypeBase
}
var expandedlicensingOrLaterOperatorType ExpandedlicensingOrLaterOperatorObjectType
var expandedlicensingOrLaterOperatorExpandedlicensingSubjectLicenseContext = map[string]string{}

func DecodeExpandedlicensingOrLaterOperator (data any, path Path, context map[string]string) (Ref[ExpandedlicensingOrLaterOperator], error) {
    return DecodeRef[ExpandedlicensingOrLaterOperator](data, path, context, expandedlicensingOrLaterOperatorType)
}

func (self ExpandedlicensingOrLaterOperatorObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingOrLaterOperator)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/subjectLicense", "expandedlicensing_subjectLicense":
        val, err := DecodeExpandedlicensingLicense(value, path, expandedlicensingOrLaterOperatorExpandedlicensingSubjectLicenseContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingSubjectLicense().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingOrLaterOperatorObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingOrLaterOperatorObject(&ExpandedlicensingOrLaterOperatorObject{}, self)
}

func ConstructExpandedlicensingOrLaterOperatorObject(o *ExpandedlicensingOrLaterOperatorObject, typ SHACLType) *ExpandedlicensingOrLaterOperatorObject {
    ConstructExpandedlicensingExtendableLicenseObject(&o.ExpandedlicensingExtendableLicenseObject, typ)
    {
        validators := []Validator[Ref[ExpandedlicensingLicense]]{}
        o.expandedlicensingSubjectLicense = NewRefProperty[ExpandedlicensingLicense]("expandedlicensingSubjectLicense", validators)
    }
    return o
}

type ExpandedlicensingOrLaterOperator interface {
    ExpandedlicensingExtendableLicense
    ExpandedlicensingSubjectLicense() RefPropertyInterface[ExpandedlicensingLicense]
}


func MakeExpandedlicensingOrLaterOperator() ExpandedlicensingOrLaterOperator {
    return ConstructExpandedlicensingOrLaterOperatorObject(&ExpandedlicensingOrLaterOperatorObject{}, expandedlicensingOrLaterOperatorType)
}

func MakeExpandedlicensingOrLaterOperatorRef() Ref[ExpandedlicensingOrLaterOperator] {
    o := MakeExpandedlicensingOrLaterOperator()
    return MakeObjectRef[ExpandedlicensingOrLaterOperator](o)
}

func (self *ExpandedlicensingOrLaterOperatorObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingExtendableLicenseObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingSubjectLicense")
        if ! self.expandedlicensingSubjectLicense.Check(prop_path, handler) {
            valid = false
        }
        if ! self.expandedlicensingSubjectLicense.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"expandedlicensingSubjectLicense", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingOrLaterOperatorObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingExtendableLicenseObject.Walk(path, visit)
    self.expandedlicensingSubjectLicense.Walk(path, visit)
}


func (self *ExpandedlicensingOrLaterOperatorObject) ExpandedlicensingSubjectLicense() RefPropertyInterface[ExpandedlicensingLicense] { return &self.expandedlicensingSubjectLicense }

func (self *ExpandedlicensingOrLaterOperatorObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingExtendableLicenseObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingSubjectLicense.IsSet() {
        val, err := EncodeRef[ExpandedlicensingLicense](self.expandedlicensingSubjectLicense.Get(), path.PushPath("expandedlicensingSubjectLicense"), expandedlicensingOrLaterOperatorExpandedlicensingSubjectLicenseContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_subjectLicense"] = val
    }
    return nil
}

// Portion of an AnyLicenseInfo representing a License which has additional

// text applied to it.
type ExpandedlicensingWithAdditionOperatorObject struct {
    SimplelicensingAnyLicenseInfoObject

    // A LicenseAddition participating in a 'with addition' model.
    expandedlicensingSubjectAddition RefProperty[ExpandedlicensingLicenseAddition]
    // A License participating in a 'with addition' model.
    expandedlicensingSubjectExtendableLicense RefProperty[ExpandedlicensingExtendableLicense]
}


type ExpandedlicensingWithAdditionOperatorObjectType struct {
    SHACLTypeBase
}
var expandedlicensingWithAdditionOperatorType ExpandedlicensingWithAdditionOperatorObjectType
var expandedlicensingWithAdditionOperatorExpandedlicensingSubjectAdditionContext = map[string]string{}
var expandedlicensingWithAdditionOperatorExpandedlicensingSubjectExtendableLicenseContext = map[string]string{}

func DecodeExpandedlicensingWithAdditionOperator (data any, path Path, context map[string]string) (Ref[ExpandedlicensingWithAdditionOperator], error) {
    return DecodeRef[ExpandedlicensingWithAdditionOperator](data, path, context, expandedlicensingWithAdditionOperatorType)
}

func (self ExpandedlicensingWithAdditionOperatorObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingWithAdditionOperator)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/subjectAddition", "expandedlicensing_subjectAddition":
        val, err := DecodeExpandedlicensingLicenseAddition(value, path, expandedlicensingWithAdditionOperatorExpandedlicensingSubjectAdditionContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingSubjectAddition().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/subjectExtendableLicense", "expandedlicensing_subjectExtendableLicense":
        val, err := DecodeExpandedlicensingExtendableLicense(value, path, expandedlicensingWithAdditionOperatorExpandedlicensingSubjectExtendableLicenseContext)
        if err != nil {
            return false, err
        }
        err = obj.ExpandedlicensingSubjectExtendableLicense().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingWithAdditionOperatorObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingWithAdditionOperatorObject(&ExpandedlicensingWithAdditionOperatorObject{}, self)
}

func ConstructExpandedlicensingWithAdditionOperatorObject(o *ExpandedlicensingWithAdditionOperatorObject, typ SHACLType) *ExpandedlicensingWithAdditionOperatorObject {
    ConstructSimplelicensingAnyLicenseInfoObject(&o.SimplelicensingAnyLicenseInfoObject, typ)
    {
        validators := []Validator[Ref[ExpandedlicensingLicenseAddition]]{}
        o.expandedlicensingSubjectAddition = NewRefProperty[ExpandedlicensingLicenseAddition]("expandedlicensingSubjectAddition", validators)
    }
    {
        validators := []Validator[Ref[ExpandedlicensingExtendableLicense]]{}
        o.expandedlicensingSubjectExtendableLicense = NewRefProperty[ExpandedlicensingExtendableLicense]("expandedlicensingSubjectExtendableLicense", validators)
    }
    return o
}

type ExpandedlicensingWithAdditionOperator interface {
    SimplelicensingAnyLicenseInfo
    ExpandedlicensingSubjectAddition() RefPropertyInterface[ExpandedlicensingLicenseAddition]
    ExpandedlicensingSubjectExtendableLicense() RefPropertyInterface[ExpandedlicensingExtendableLicense]
}


func MakeExpandedlicensingWithAdditionOperator() ExpandedlicensingWithAdditionOperator {
    return ConstructExpandedlicensingWithAdditionOperatorObject(&ExpandedlicensingWithAdditionOperatorObject{}, expandedlicensingWithAdditionOperatorType)
}

func MakeExpandedlicensingWithAdditionOperatorRef() Ref[ExpandedlicensingWithAdditionOperator] {
    o := MakeExpandedlicensingWithAdditionOperator()
    return MakeObjectRef[ExpandedlicensingWithAdditionOperator](o)
}

func (self *ExpandedlicensingWithAdditionOperatorObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SimplelicensingAnyLicenseInfoObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("expandedlicensingSubjectAddition")
        if ! self.expandedlicensingSubjectAddition.Check(prop_path, handler) {
            valid = false
        }
        if ! self.expandedlicensingSubjectAddition.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"expandedlicensingSubjectAddition", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("expandedlicensingSubjectExtendableLicense")
        if ! self.expandedlicensingSubjectExtendableLicense.Check(prop_path, handler) {
            valid = false
        }
        if ! self.expandedlicensingSubjectExtendableLicense.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"expandedlicensingSubjectExtendableLicense", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExpandedlicensingWithAdditionOperatorObject) Walk(path Path, visit Visit) {
    self.SimplelicensingAnyLicenseInfoObject.Walk(path, visit)
    self.expandedlicensingSubjectAddition.Walk(path, visit)
    self.expandedlicensingSubjectExtendableLicense.Walk(path, visit)
}


func (self *ExpandedlicensingWithAdditionOperatorObject) ExpandedlicensingSubjectAddition() RefPropertyInterface[ExpandedlicensingLicenseAddition] { return &self.expandedlicensingSubjectAddition }
func (self *ExpandedlicensingWithAdditionOperatorObject) ExpandedlicensingSubjectExtendableLicense() RefPropertyInterface[ExpandedlicensingExtendableLicense] { return &self.expandedlicensingSubjectExtendableLicense }

func (self *ExpandedlicensingWithAdditionOperatorObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SimplelicensingAnyLicenseInfoObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.expandedlicensingSubjectAddition.IsSet() {
        val, err := EncodeRef[ExpandedlicensingLicenseAddition](self.expandedlicensingSubjectAddition.Get(), path.PushPath("expandedlicensingSubjectAddition"), expandedlicensingWithAdditionOperatorExpandedlicensingSubjectAdditionContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_subjectAddition"] = val
    }
    if self.expandedlicensingSubjectExtendableLicense.IsSet() {
        val, err := EncodeRef[ExpandedlicensingExtendableLicense](self.expandedlicensingSubjectExtendableLicense.Get(), path.PushPath("expandedlicensingSubjectExtendableLicense"), expandedlicensingWithAdditionOperatorExpandedlicensingSubjectExtendableLicenseContext, state)
        if err != nil {
            return err
        }
        data["expandedlicensing_subjectExtendableLicense"] = val
    }
    return nil
}

// A type of extension consisting of a list of name value pairs.
type ExtensionCdxPropertiesExtensionObject struct {
    ExtensionExtensionObject

    // Provides a map of a property names to a values.
    extensionCdxProperty RefListProperty[ExtensionCdxPropertyEntry]
}


type ExtensionCdxPropertiesExtensionObjectType struct {
    SHACLTypeBase
}
var extensionCdxPropertiesExtensionType ExtensionCdxPropertiesExtensionObjectType
var extensionCdxPropertiesExtensionExtensionCdxPropertyContext = map[string]string{}

func DecodeExtensionCdxPropertiesExtension (data any, path Path, context map[string]string) (Ref[ExtensionCdxPropertiesExtension], error) {
    return DecodeRef[ExtensionCdxPropertiesExtension](data, path, context, extensionCdxPropertiesExtensionType)
}

func (self ExtensionCdxPropertiesExtensionObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExtensionCdxPropertiesExtension)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Extension/cdxProperty", "extension_cdxProperty":
        val, err := DecodeList[Ref[ExtensionCdxPropertyEntry]](value, path, extensionCdxPropertiesExtensionExtensionCdxPropertyContext, DecodeExtensionCdxPropertyEntry)
        if err != nil {
            return false, err
        }
        err = obj.ExtensionCdxProperty().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExtensionCdxPropertiesExtensionObjectType) Create() SHACLObject {
    return ConstructExtensionCdxPropertiesExtensionObject(&ExtensionCdxPropertiesExtensionObject{}, self)
}

func ConstructExtensionCdxPropertiesExtensionObject(o *ExtensionCdxPropertiesExtensionObject, typ SHACLType) *ExtensionCdxPropertiesExtensionObject {
    ConstructExtensionExtensionObject(&o.ExtensionExtensionObject, typ)
    {
        validators := []Validator[Ref[ExtensionCdxPropertyEntry]]{}
        o.extensionCdxProperty = NewRefListProperty[ExtensionCdxPropertyEntry]("extensionCdxProperty", validators)
    }
    return o
}

type ExtensionCdxPropertiesExtension interface {
    ExtensionExtension
    ExtensionCdxProperty() ListPropertyInterface[Ref[ExtensionCdxPropertyEntry]]
}


func MakeExtensionCdxPropertiesExtension() ExtensionCdxPropertiesExtension {
    return ConstructExtensionCdxPropertiesExtensionObject(&ExtensionCdxPropertiesExtensionObject{}, extensionCdxPropertiesExtensionType)
}

func MakeExtensionCdxPropertiesExtensionRef() Ref[ExtensionCdxPropertiesExtension] {
    o := MakeExtensionCdxPropertiesExtension()
    return MakeObjectRef[ExtensionCdxPropertiesExtension](o)
}

func (self *ExtensionCdxPropertiesExtensionObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExtensionExtensionObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("extensionCdxProperty")
        if ! self.extensionCdxProperty.Check(prop_path, handler) {
            valid = false
        }
        if len(self.extensionCdxProperty.Get()) < 1 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "extensionCdxProperty",
                    "Too few elements. Minimum of 1 required"},
                    prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *ExtensionCdxPropertiesExtensionObject) Walk(path Path, visit Visit) {
    self.ExtensionExtensionObject.Walk(path, visit)
    self.extensionCdxProperty.Walk(path, visit)
}


func (self *ExtensionCdxPropertiesExtensionObject) ExtensionCdxProperty() ListPropertyInterface[Ref[ExtensionCdxPropertyEntry]] { return &self.extensionCdxProperty }

func (self *ExtensionCdxPropertiesExtensionObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExtensionExtensionObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.extensionCdxProperty.IsSet() {
        val, err := EncodeList[Ref[ExtensionCdxPropertyEntry]](self.extensionCdxProperty.Get(), path.PushPath("extensionCdxProperty"), extensionCdxPropertiesExtensionExtensionCdxPropertyContext, state, EncodeRef[ExtensionCdxPropertyEntry])
        if err != nil {
            return err
        }
        data["extension_cdxProperty"] = val
    }
    return nil
}

// Provides a CVSS version 2.0 assessment for a vulnerability.
type SecurityCvssV2VulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Provides a numerical (0-10) representation of the severity of a vulnerability.
    securityScore Property[float64]
    // Specifies the CVSS vector string for a vulnerability.
    securityVectorString Property[string]
}


type SecurityCvssV2VulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityCvssV2VulnAssessmentRelationshipType SecurityCvssV2VulnAssessmentRelationshipObjectType
var securityCvssV2VulnAssessmentRelationshipSecurityScoreContext = map[string]string{}
var securityCvssV2VulnAssessmentRelationshipSecurityVectorStringContext = map[string]string{}

func DecodeSecurityCvssV2VulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityCvssV2VulnAssessmentRelationship], error) {
    return DecodeRef[SecurityCvssV2VulnAssessmentRelationship](data, path, context, securityCvssV2VulnAssessmentRelationshipType)
}

func (self SecurityCvssV2VulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityCvssV2VulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/score", "security_score":
        val, err := DecodeFloat(value, path, securityCvssV2VulnAssessmentRelationshipSecurityScoreContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityScore().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/vectorString", "security_vectorString":
        val, err := DecodeString(value, path, securityCvssV2VulnAssessmentRelationshipSecurityVectorStringContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityVectorString().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityCvssV2VulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityCvssV2VulnAssessmentRelationshipObject(&SecurityCvssV2VulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityCvssV2VulnAssessmentRelationshipObject(o *SecurityCvssV2VulnAssessmentRelationshipObject, typ SHACLType) *SecurityCvssV2VulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[float64]{}
        o.securityScore = NewProperty[float64]("securityScore", validators)
    }
    {
        validators := []Validator[string]{}
        o.securityVectorString = NewProperty[string]("securityVectorString", validators)
    }
    return o
}

type SecurityCvssV2VulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityScore() PropertyInterface[float64]
    SecurityVectorString() PropertyInterface[string]
}


func MakeSecurityCvssV2VulnAssessmentRelationship() SecurityCvssV2VulnAssessmentRelationship {
    return ConstructSecurityCvssV2VulnAssessmentRelationshipObject(&SecurityCvssV2VulnAssessmentRelationshipObject{}, securityCvssV2VulnAssessmentRelationshipType)
}

func MakeSecurityCvssV2VulnAssessmentRelationshipRef() Ref[SecurityCvssV2VulnAssessmentRelationship] {
    o := MakeSecurityCvssV2VulnAssessmentRelationship()
    return MakeObjectRef[SecurityCvssV2VulnAssessmentRelationship](o)
}

func (self *SecurityCvssV2VulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityScore")
        if ! self.securityScore.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityScore.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityScore", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityVectorString")
        if ! self.securityVectorString.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityVectorString.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityVectorString", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecurityCvssV2VulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityScore.Walk(path, visit)
    self.securityVectorString.Walk(path, visit)
}


func (self *SecurityCvssV2VulnAssessmentRelationshipObject) SecurityScore() PropertyInterface[float64] { return &self.securityScore }
func (self *SecurityCvssV2VulnAssessmentRelationshipObject) SecurityVectorString() PropertyInterface[string] { return &self.securityVectorString }

func (self *SecurityCvssV2VulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityScore.IsSet() {
        val, err := EncodeFloat(self.securityScore.Get(), path.PushPath("securityScore"), securityCvssV2VulnAssessmentRelationshipSecurityScoreContext, state)
        if err != nil {
            return err
        }
        data["security_score"] = val
    }
    if self.securityVectorString.IsSet() {
        val, err := EncodeString(self.securityVectorString.Get(), path.PushPath("securityVectorString"), securityCvssV2VulnAssessmentRelationshipSecurityVectorStringContext, state)
        if err != nil {
            return err
        }
        data["security_vectorString"] = val
    }
    return nil
}

// Provides a CVSS version 3 assessment for a vulnerability.
type SecurityCvssV3VulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Provides a numerical (0-10) representation of the severity of a vulnerability.
    securityScore Property[float64]
    // Specifies the CVSS qualitative severity rating of a vulnerability in relation to a piece of software.
    securitySeverity Property[string]
    // Specifies the CVSS vector string for a vulnerability.
    securityVectorString Property[string]
}


type SecurityCvssV3VulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityCvssV3VulnAssessmentRelationshipType SecurityCvssV3VulnAssessmentRelationshipObjectType
var securityCvssV3VulnAssessmentRelationshipSecurityScoreContext = map[string]string{}
var securityCvssV3VulnAssessmentRelationshipSecuritySeverityContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/critical": "critical",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/high": "high",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/low": "low",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/medium": "medium",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/none": "none",}
var securityCvssV3VulnAssessmentRelationshipSecurityVectorStringContext = map[string]string{}

func DecodeSecurityCvssV3VulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityCvssV3VulnAssessmentRelationship], error) {
    return DecodeRef[SecurityCvssV3VulnAssessmentRelationship](data, path, context, securityCvssV3VulnAssessmentRelationshipType)
}

func (self SecurityCvssV3VulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityCvssV3VulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/score", "security_score":
        val, err := DecodeFloat(value, path, securityCvssV3VulnAssessmentRelationshipSecurityScoreContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityScore().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/severity", "security_severity":
        val, err := DecodeIRI(value, path, securityCvssV3VulnAssessmentRelationshipSecuritySeverityContext)
        if err != nil {
            return false, err
        }
        err = obj.SecuritySeverity().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/vectorString", "security_vectorString":
        val, err := DecodeString(value, path, securityCvssV3VulnAssessmentRelationshipSecurityVectorStringContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityVectorString().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityCvssV3VulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityCvssV3VulnAssessmentRelationshipObject(&SecurityCvssV3VulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityCvssV3VulnAssessmentRelationshipObject(o *SecurityCvssV3VulnAssessmentRelationshipObject, typ SHACLType) *SecurityCvssV3VulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[float64]{}
        o.securityScore = NewProperty[float64]("securityScore", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/critical",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/high",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/low",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/medium",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/none",
        }})
        o.securitySeverity = NewProperty[string]("securitySeverity", validators)
    }
    {
        validators := []Validator[string]{}
        o.securityVectorString = NewProperty[string]("securityVectorString", validators)
    }
    return o
}

type SecurityCvssV3VulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityScore() PropertyInterface[float64]
    SecuritySeverity() PropertyInterface[string]
    SecurityVectorString() PropertyInterface[string]
}


func MakeSecurityCvssV3VulnAssessmentRelationship() SecurityCvssV3VulnAssessmentRelationship {
    return ConstructSecurityCvssV3VulnAssessmentRelationshipObject(&SecurityCvssV3VulnAssessmentRelationshipObject{}, securityCvssV3VulnAssessmentRelationshipType)
}

func MakeSecurityCvssV3VulnAssessmentRelationshipRef() Ref[SecurityCvssV3VulnAssessmentRelationship] {
    o := MakeSecurityCvssV3VulnAssessmentRelationship()
    return MakeObjectRef[SecurityCvssV3VulnAssessmentRelationship](o)
}

func (self *SecurityCvssV3VulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityScore")
        if ! self.securityScore.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityScore.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityScore", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securitySeverity")
        if ! self.securitySeverity.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securitySeverity.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securitySeverity", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityVectorString")
        if ! self.securityVectorString.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityVectorString.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityVectorString", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecurityCvssV3VulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityScore.Walk(path, visit)
    self.securitySeverity.Walk(path, visit)
    self.securityVectorString.Walk(path, visit)
}


func (self *SecurityCvssV3VulnAssessmentRelationshipObject) SecurityScore() PropertyInterface[float64] { return &self.securityScore }
func (self *SecurityCvssV3VulnAssessmentRelationshipObject) SecuritySeverity() PropertyInterface[string] { return &self.securitySeverity }
func (self *SecurityCvssV3VulnAssessmentRelationshipObject) SecurityVectorString() PropertyInterface[string] { return &self.securityVectorString }

func (self *SecurityCvssV3VulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityScore.IsSet() {
        val, err := EncodeFloat(self.securityScore.Get(), path.PushPath("securityScore"), securityCvssV3VulnAssessmentRelationshipSecurityScoreContext, state)
        if err != nil {
            return err
        }
        data["security_score"] = val
    }
    if self.securitySeverity.IsSet() {
        val, err := EncodeIRI(self.securitySeverity.Get(), path.PushPath("securitySeverity"), securityCvssV3VulnAssessmentRelationshipSecuritySeverityContext, state)
        if err != nil {
            return err
        }
        data["security_severity"] = val
    }
    if self.securityVectorString.IsSet() {
        val, err := EncodeString(self.securityVectorString.Get(), path.PushPath("securityVectorString"), securityCvssV3VulnAssessmentRelationshipSecurityVectorStringContext, state)
        if err != nil {
            return err
        }
        data["security_vectorString"] = val
    }
    return nil
}

// Provides a CVSS version 4 assessment for a vulnerability.
type SecurityCvssV4VulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Provides a numerical (0-10) representation of the severity of a vulnerability.
    securityScore Property[float64]
    // Specifies the CVSS qualitative severity rating of a vulnerability in relation to a piece of software.
    securitySeverity Property[string]
    // Specifies the CVSS vector string for a vulnerability.
    securityVectorString Property[string]
}


type SecurityCvssV4VulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityCvssV4VulnAssessmentRelationshipType SecurityCvssV4VulnAssessmentRelationshipObjectType
var securityCvssV4VulnAssessmentRelationshipSecurityScoreContext = map[string]string{}
var securityCvssV4VulnAssessmentRelationshipSecuritySeverityContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/critical": "critical",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/high": "high",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/low": "low",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/medium": "medium",
    "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/none": "none",}
var securityCvssV4VulnAssessmentRelationshipSecurityVectorStringContext = map[string]string{}

func DecodeSecurityCvssV4VulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityCvssV4VulnAssessmentRelationship], error) {
    return DecodeRef[SecurityCvssV4VulnAssessmentRelationship](data, path, context, securityCvssV4VulnAssessmentRelationshipType)
}

func (self SecurityCvssV4VulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityCvssV4VulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/score", "security_score":
        val, err := DecodeFloat(value, path, securityCvssV4VulnAssessmentRelationshipSecurityScoreContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityScore().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/severity", "security_severity":
        val, err := DecodeIRI(value, path, securityCvssV4VulnAssessmentRelationshipSecuritySeverityContext)
        if err != nil {
            return false, err
        }
        err = obj.SecuritySeverity().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/vectorString", "security_vectorString":
        val, err := DecodeString(value, path, securityCvssV4VulnAssessmentRelationshipSecurityVectorStringContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityVectorString().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityCvssV4VulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityCvssV4VulnAssessmentRelationshipObject(&SecurityCvssV4VulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityCvssV4VulnAssessmentRelationshipObject(o *SecurityCvssV4VulnAssessmentRelationshipObject, typ SHACLType) *SecurityCvssV4VulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[float64]{}
        o.securityScore = NewProperty[float64]("securityScore", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/critical",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/high",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/low",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/medium",
                "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType/none",
        }})
        o.securitySeverity = NewProperty[string]("securitySeverity", validators)
    }
    {
        validators := []Validator[string]{}
        o.securityVectorString = NewProperty[string]("securityVectorString", validators)
    }
    return o
}

type SecurityCvssV4VulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityScore() PropertyInterface[float64]
    SecuritySeverity() PropertyInterface[string]
    SecurityVectorString() PropertyInterface[string]
}


func MakeSecurityCvssV4VulnAssessmentRelationship() SecurityCvssV4VulnAssessmentRelationship {
    return ConstructSecurityCvssV4VulnAssessmentRelationshipObject(&SecurityCvssV4VulnAssessmentRelationshipObject{}, securityCvssV4VulnAssessmentRelationshipType)
}

func MakeSecurityCvssV4VulnAssessmentRelationshipRef() Ref[SecurityCvssV4VulnAssessmentRelationship] {
    o := MakeSecurityCvssV4VulnAssessmentRelationship()
    return MakeObjectRef[SecurityCvssV4VulnAssessmentRelationship](o)
}

func (self *SecurityCvssV4VulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityScore")
        if ! self.securityScore.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityScore.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityScore", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securitySeverity")
        if ! self.securitySeverity.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securitySeverity.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securitySeverity", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityVectorString")
        if ! self.securityVectorString.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityVectorString.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityVectorString", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecurityCvssV4VulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityScore.Walk(path, visit)
    self.securitySeverity.Walk(path, visit)
    self.securityVectorString.Walk(path, visit)
}


func (self *SecurityCvssV4VulnAssessmentRelationshipObject) SecurityScore() PropertyInterface[float64] { return &self.securityScore }
func (self *SecurityCvssV4VulnAssessmentRelationshipObject) SecuritySeverity() PropertyInterface[string] { return &self.securitySeverity }
func (self *SecurityCvssV4VulnAssessmentRelationshipObject) SecurityVectorString() PropertyInterface[string] { return &self.securityVectorString }

func (self *SecurityCvssV4VulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityScore.IsSet() {
        val, err := EncodeFloat(self.securityScore.Get(), path.PushPath("securityScore"), securityCvssV4VulnAssessmentRelationshipSecurityScoreContext, state)
        if err != nil {
            return err
        }
        data["security_score"] = val
    }
    if self.securitySeverity.IsSet() {
        val, err := EncodeIRI(self.securitySeverity.Get(), path.PushPath("securitySeverity"), securityCvssV4VulnAssessmentRelationshipSecuritySeverityContext, state)
        if err != nil {
            return err
        }
        data["security_severity"] = val
    }
    if self.securityVectorString.IsSet() {
        val, err := EncodeString(self.securityVectorString.Get(), path.PushPath("securityVectorString"), securityCvssV4VulnAssessmentRelationshipSecurityVectorStringContext, state)
        if err != nil {
            return err
        }
        data["security_vectorString"] = val
    }
    return nil
}

// Provides an EPSS assessment for a vulnerability.
type SecurityEpssVulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // The percentile of the current probability score.
    securityPercentile Property[float64]
    // A probability score between 0 and 1 of a vulnerability being exploited.
    securityProbability Property[float64]
}


type SecurityEpssVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityEpssVulnAssessmentRelationshipType SecurityEpssVulnAssessmentRelationshipObjectType
var securityEpssVulnAssessmentRelationshipSecurityPercentileContext = map[string]string{}
var securityEpssVulnAssessmentRelationshipSecurityProbabilityContext = map[string]string{}

func DecodeSecurityEpssVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityEpssVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityEpssVulnAssessmentRelationship](data, path, context, securityEpssVulnAssessmentRelationshipType)
}

func (self SecurityEpssVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityEpssVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/percentile", "security_percentile":
        val, err := DecodeFloat(value, path, securityEpssVulnAssessmentRelationshipSecurityPercentileContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityPercentile().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/probability", "security_probability":
        val, err := DecodeFloat(value, path, securityEpssVulnAssessmentRelationshipSecurityProbabilityContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityProbability().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityEpssVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityEpssVulnAssessmentRelationshipObject(&SecurityEpssVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityEpssVulnAssessmentRelationshipObject(o *SecurityEpssVulnAssessmentRelationshipObject, typ SHACLType) *SecurityEpssVulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[float64]{}
        o.securityPercentile = NewProperty[float64]("securityPercentile", validators)
    }
    {
        validators := []Validator[float64]{}
        o.securityProbability = NewProperty[float64]("securityProbability", validators)
    }
    return o
}

type SecurityEpssVulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityPercentile() PropertyInterface[float64]
    SecurityProbability() PropertyInterface[float64]
}


func MakeSecurityEpssVulnAssessmentRelationship() SecurityEpssVulnAssessmentRelationship {
    return ConstructSecurityEpssVulnAssessmentRelationshipObject(&SecurityEpssVulnAssessmentRelationshipObject{}, securityEpssVulnAssessmentRelationshipType)
}

func MakeSecurityEpssVulnAssessmentRelationshipRef() Ref[SecurityEpssVulnAssessmentRelationship] {
    o := MakeSecurityEpssVulnAssessmentRelationship()
    return MakeObjectRef[SecurityEpssVulnAssessmentRelationship](o)
}

func (self *SecurityEpssVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityPercentile")
        if ! self.securityPercentile.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityPercentile.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityPercentile", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityProbability")
        if ! self.securityProbability.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityProbability.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityProbability", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecurityEpssVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityPercentile.Walk(path, visit)
    self.securityProbability.Walk(path, visit)
}


func (self *SecurityEpssVulnAssessmentRelationshipObject) SecurityPercentile() PropertyInterface[float64] { return &self.securityPercentile }
func (self *SecurityEpssVulnAssessmentRelationshipObject) SecurityProbability() PropertyInterface[float64] { return &self.securityProbability }

func (self *SecurityEpssVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityPercentile.IsSet() {
        val, err := EncodeFloat(self.securityPercentile.Get(), path.PushPath("securityPercentile"), securityEpssVulnAssessmentRelationshipSecurityPercentileContext, state)
        if err != nil {
            return err
        }
        data["security_percentile"] = val
    }
    if self.securityProbability.IsSet() {
        val, err := EncodeFloat(self.securityProbability.Get(), path.PushPath("securityProbability"), securityEpssVulnAssessmentRelationshipSecurityProbabilityContext, state)
        if err != nil {
            return err
        }
        data["security_probability"] = val
    }
    return nil
}

// Provides an exploit assessment of a vulnerability.
type SecurityExploitCatalogVulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Specifies the exploit catalog type.
    securityCatalogType Property[string]
    // Describe that a CVE is known to have an exploit because it's been listed in an exploit catalog.
    securityExploited Property[bool]
    // Provides the location of an exploit catalog.
    securityLocator Property[string]
}


type SecurityExploitCatalogVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityExploitCatalogVulnAssessmentRelationshipType SecurityExploitCatalogVulnAssessmentRelationshipObjectType
var securityExploitCatalogVulnAssessmentRelationshipSecurityCatalogTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/kev": "kev",
    "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/other": "other",}
var securityExploitCatalogVulnAssessmentRelationshipSecurityExploitedContext = map[string]string{}
var securityExploitCatalogVulnAssessmentRelationshipSecurityLocatorContext = map[string]string{}

func DecodeSecurityExploitCatalogVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityExploitCatalogVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityExploitCatalogVulnAssessmentRelationship](data, path, context, securityExploitCatalogVulnAssessmentRelationshipType)
}

func (self SecurityExploitCatalogVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityExploitCatalogVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/catalogType", "security_catalogType":
        val, err := DecodeIRI(value, path, securityExploitCatalogVulnAssessmentRelationshipSecurityCatalogTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityCatalogType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/exploited", "security_exploited":
        val, err := DecodeBoolean(value, path, securityExploitCatalogVulnAssessmentRelationshipSecurityExploitedContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityExploited().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/locator", "security_locator":
        val, err := DecodeString(value, path, securityExploitCatalogVulnAssessmentRelationshipSecurityLocatorContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityLocator().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityExploitCatalogVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityExploitCatalogVulnAssessmentRelationshipObject(&SecurityExploitCatalogVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityExploitCatalogVulnAssessmentRelationshipObject(o *SecurityExploitCatalogVulnAssessmentRelationshipObject, typ SHACLType) *SecurityExploitCatalogVulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/kev",
                "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType/other",
        }})
        o.securityCatalogType = NewProperty[string]("securityCatalogType", validators)
    }
    {
        validators := []Validator[bool]{}
        o.securityExploited = NewProperty[bool]("securityExploited", validators)
    }
    {
        validators := []Validator[string]{}
        o.securityLocator = NewProperty[string]("securityLocator", validators)
    }
    return o
}

type SecurityExploitCatalogVulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityCatalogType() PropertyInterface[string]
    SecurityExploited() PropertyInterface[bool]
    SecurityLocator() PropertyInterface[string]
}


func MakeSecurityExploitCatalogVulnAssessmentRelationship() SecurityExploitCatalogVulnAssessmentRelationship {
    return ConstructSecurityExploitCatalogVulnAssessmentRelationshipObject(&SecurityExploitCatalogVulnAssessmentRelationshipObject{}, securityExploitCatalogVulnAssessmentRelationshipType)
}

func MakeSecurityExploitCatalogVulnAssessmentRelationshipRef() Ref[SecurityExploitCatalogVulnAssessmentRelationship] {
    o := MakeSecurityExploitCatalogVulnAssessmentRelationship()
    return MakeObjectRef[SecurityExploitCatalogVulnAssessmentRelationship](o)
}

func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityCatalogType")
        if ! self.securityCatalogType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityCatalogType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityCatalogType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityExploited")
        if ! self.securityExploited.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityExploited.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityExploited", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityLocator")
        if ! self.securityLocator.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityLocator.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityLocator", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityCatalogType.Walk(path, visit)
    self.securityExploited.Walk(path, visit)
    self.securityLocator.Walk(path, visit)
}


func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) SecurityCatalogType() PropertyInterface[string] { return &self.securityCatalogType }
func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) SecurityExploited() PropertyInterface[bool] { return &self.securityExploited }
func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) SecurityLocator() PropertyInterface[string] { return &self.securityLocator }

func (self *SecurityExploitCatalogVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityCatalogType.IsSet() {
        val, err := EncodeIRI(self.securityCatalogType.Get(), path.PushPath("securityCatalogType"), securityExploitCatalogVulnAssessmentRelationshipSecurityCatalogTypeContext, state)
        if err != nil {
            return err
        }
        data["security_catalogType"] = val
    }
    if self.securityExploited.IsSet() {
        val, err := EncodeBoolean(self.securityExploited.Get(), path.PushPath("securityExploited"), securityExploitCatalogVulnAssessmentRelationshipSecurityExploitedContext, state)
        if err != nil {
            return err
        }
        data["security_exploited"] = val
    }
    if self.securityLocator.IsSet() {
        val, err := EncodeString(self.securityLocator.Get(), path.PushPath("securityLocator"), securityExploitCatalogVulnAssessmentRelationshipSecurityLocatorContext, state)
        if err != nil {
            return err
        }
        data["security_locator"] = val
    }
    return nil
}

// Provides an SSVC assessment for a vulnerability.
type SecuritySsvcVulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Provide the enumeration of possible decisions in the
    // [Stakeholder-Specific Vulnerability Categorization (SSVC) decision tree](https://www.cisa.gov/stakeholder-specific-vulnerability-categorization-ssvc).
    securityDecisionType Property[string]
}


type SecuritySsvcVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securitySsvcVulnAssessmentRelationshipType SecuritySsvcVulnAssessmentRelationshipObjectType
var securitySsvcVulnAssessmentRelationshipSecurityDecisionTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/act": "act",
    "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/attend": "attend",
    "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/track": "track",
    "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/trackStar": "trackStar",}

func DecodeSecuritySsvcVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecuritySsvcVulnAssessmentRelationship], error) {
    return DecodeRef[SecuritySsvcVulnAssessmentRelationship](data, path, context, securitySsvcVulnAssessmentRelationshipType)
}

func (self SecuritySsvcVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecuritySsvcVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/decisionType", "security_decisionType":
        val, err := DecodeIRI(value, path, securitySsvcVulnAssessmentRelationshipSecurityDecisionTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityDecisionType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecuritySsvcVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecuritySsvcVulnAssessmentRelationshipObject(&SecuritySsvcVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecuritySsvcVulnAssessmentRelationshipObject(o *SecuritySsvcVulnAssessmentRelationshipObject, typ SHACLType) *SecuritySsvcVulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/act",
                "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/attend",
                "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/track",
                "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType/trackStar",
        }})
        o.securityDecisionType = NewProperty[string]("securityDecisionType", validators)
    }
    return o
}

type SecuritySsvcVulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityDecisionType() PropertyInterface[string]
}


func MakeSecuritySsvcVulnAssessmentRelationship() SecuritySsvcVulnAssessmentRelationship {
    return ConstructSecuritySsvcVulnAssessmentRelationshipObject(&SecuritySsvcVulnAssessmentRelationshipObject{}, securitySsvcVulnAssessmentRelationshipType)
}

func MakeSecuritySsvcVulnAssessmentRelationshipRef() Ref[SecuritySsvcVulnAssessmentRelationship] {
    o := MakeSecuritySsvcVulnAssessmentRelationship()
    return MakeObjectRef[SecuritySsvcVulnAssessmentRelationship](o)
}

func (self *SecuritySsvcVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityDecisionType")
        if ! self.securityDecisionType.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityDecisionType.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityDecisionType", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SecuritySsvcVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityDecisionType.Walk(path, visit)
}


func (self *SecuritySsvcVulnAssessmentRelationshipObject) SecurityDecisionType() PropertyInterface[string] { return &self.securityDecisionType }

func (self *SecuritySsvcVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityDecisionType.IsSet() {
        val, err := EncodeIRI(self.securityDecisionType.Get(), path.PushPath("securityDecisionType"), securitySsvcVulnAssessmentRelationshipSecurityDecisionTypeContext, state)
        if err != nil {
            return err
        }
        data["security_decisionType"] = val
    }
    return nil
}

// Abstract ancestor class for all VEX relationships
type SecurityVexVulnAssessmentRelationshipObject struct {
    SecurityVulnAssessmentRelationshipObject

    // Conveys information about how VEX status was determined.
    securityStatusNotes Property[string]
    // Specifies the version of a VEX statement.
    securityVexVersion Property[string]
}


type SecurityVexVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVexVulnAssessmentRelationshipType SecurityVexVulnAssessmentRelationshipObjectType
var securityVexVulnAssessmentRelationshipSecurityStatusNotesContext = map[string]string{}
var securityVexVulnAssessmentRelationshipSecurityVexVersionContext = map[string]string{}

func DecodeSecurityVexVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVexVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVexVulnAssessmentRelationship](data, path, context, securityVexVulnAssessmentRelationshipType)
}

func (self SecurityVexVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/statusNotes", "security_statusNotes":
        val, err := DecodeString(value, path, securityVexVulnAssessmentRelationshipSecurityStatusNotesContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityStatusNotes().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/vexVersion", "security_vexVersion":
        val, err := DecodeString(value, path, securityVexVulnAssessmentRelationshipSecurityVexVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityVexVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVexVulnAssessmentRelationshipObject(&SecurityVexVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVexVulnAssessmentRelationshipObject(o *SecurityVexVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVexVulnAssessmentRelationshipObject {
    ConstructSecurityVulnAssessmentRelationshipObject(&o.SecurityVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[string]{}
        o.securityStatusNotes = NewProperty[string]("securityStatusNotes", validators)
    }
    {
        validators := []Validator[string]{}
        o.securityVexVersion = NewProperty[string]("securityVexVersion", validators)
    }
    return o
}

type SecurityVexVulnAssessmentRelationship interface {
    SecurityVulnAssessmentRelationship
    SecurityStatusNotes() PropertyInterface[string]
    SecurityVexVersion() PropertyInterface[string]
}



func (self *SecurityVexVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityStatusNotes")
        if ! self.securityStatusNotes.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityVexVersion")
        if ! self.securityVexVersion.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SecurityVexVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityStatusNotes.Walk(path, visit)
    self.securityVexVersion.Walk(path, visit)
}


func (self *SecurityVexVulnAssessmentRelationshipObject) SecurityStatusNotes() PropertyInterface[string] { return &self.securityStatusNotes }
func (self *SecurityVexVulnAssessmentRelationshipObject) SecurityVexVersion() PropertyInterface[string] { return &self.securityVexVersion }

func (self *SecurityVexVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityStatusNotes.IsSet() {
        val, err := EncodeString(self.securityStatusNotes.Get(), path.PushPath("securityStatusNotes"), securityVexVulnAssessmentRelationshipSecurityStatusNotesContext, state)
        if err != nil {
            return err
        }
        data["security_statusNotes"] = val
    }
    if self.securityVexVersion.IsSet() {
        val, err := EncodeString(self.securityVexVersion.Get(), path.PushPath("securityVexVersion"), securityVexVulnAssessmentRelationshipSecurityVexVersionContext, state)
        if err != nil {
            return err
        }
        data["security_vexVersion"] = val
    }
    return nil
}

// Specifies a vulnerability and its associated information.
type SecurityVulnerabilityObject struct {
    ArtifactObject

    // Specifies a time when a vulnerability assessment was modified
    securityModifiedTime Property[time.Time]
    // Specifies the time when a vulnerability was published.
    securityPublishedTime Property[time.Time]
    // Specified the time and date when a vulnerability was withdrawn.
    securityWithdrawnTime Property[time.Time]
}


type SecurityVulnerabilityObjectType struct {
    SHACLTypeBase
}
var securityVulnerabilityType SecurityVulnerabilityObjectType
var securityVulnerabilitySecurityModifiedTimeContext = map[string]string{}
var securityVulnerabilitySecurityPublishedTimeContext = map[string]string{}
var securityVulnerabilitySecurityWithdrawnTimeContext = map[string]string{}

func DecodeSecurityVulnerability (data any, path Path, context map[string]string) (Ref[SecurityVulnerability], error) {
    return DecodeRef[SecurityVulnerability](data, path, context, securityVulnerabilityType)
}

func (self SecurityVulnerabilityObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVulnerability)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/modifiedTime", "security_modifiedTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnerabilitySecurityModifiedTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityModifiedTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/publishedTime", "security_publishedTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnerabilitySecurityPublishedTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityPublishedTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/withdrawnTime", "security_withdrawnTime":
        val, err := DecodeDateTimeStamp(value, path, securityVulnerabilitySecurityWithdrawnTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityWithdrawnTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVulnerabilityObjectType) Create() SHACLObject {
    return ConstructSecurityVulnerabilityObject(&SecurityVulnerabilityObject{}, self)
}

func ConstructSecurityVulnerabilityObject(o *SecurityVulnerabilityObject, typ SHACLType) *SecurityVulnerabilityObject {
    ConstructArtifactObject(&o.ArtifactObject, typ)
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityModifiedTime = NewProperty[time.Time]("securityModifiedTime", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityPublishedTime = NewProperty[time.Time]("securityPublishedTime", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityWithdrawnTime = NewProperty[time.Time]("securityWithdrawnTime", validators)
    }
    return o
}

type SecurityVulnerability interface {
    Artifact
    SecurityModifiedTime() PropertyInterface[time.Time]
    SecurityPublishedTime() PropertyInterface[time.Time]
    SecurityWithdrawnTime() PropertyInterface[time.Time]
}


func MakeSecurityVulnerability() SecurityVulnerability {
    return ConstructSecurityVulnerabilityObject(&SecurityVulnerabilityObject{}, securityVulnerabilityType)
}

func MakeSecurityVulnerabilityRef() Ref[SecurityVulnerability] {
    o := MakeSecurityVulnerability()
    return MakeObjectRef[SecurityVulnerability](o)
}

func (self *SecurityVulnerabilityObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ArtifactObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityModifiedTime")
        if ! self.securityModifiedTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityPublishedTime")
        if ! self.securityPublishedTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityWithdrawnTime")
        if ! self.securityWithdrawnTime.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SecurityVulnerabilityObject) Walk(path Path, visit Visit) {
    self.ArtifactObject.Walk(path, visit)
    self.securityModifiedTime.Walk(path, visit)
    self.securityPublishedTime.Walk(path, visit)
    self.securityWithdrawnTime.Walk(path, visit)
}


func (self *SecurityVulnerabilityObject) SecurityModifiedTime() PropertyInterface[time.Time] { return &self.securityModifiedTime }
func (self *SecurityVulnerabilityObject) SecurityPublishedTime() PropertyInterface[time.Time] { return &self.securityPublishedTime }
func (self *SecurityVulnerabilityObject) SecurityWithdrawnTime() PropertyInterface[time.Time] { return &self.securityWithdrawnTime }

func (self *SecurityVulnerabilityObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ArtifactObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityModifiedTime.IsSet() {
        val, err := EncodeDateTime(self.securityModifiedTime.Get(), path.PushPath("securityModifiedTime"), securityVulnerabilitySecurityModifiedTimeContext, state)
        if err != nil {
            return err
        }
        data["security_modifiedTime"] = val
    }
    if self.securityPublishedTime.IsSet() {
        val, err := EncodeDateTime(self.securityPublishedTime.Get(), path.PushPath("securityPublishedTime"), securityVulnerabilitySecurityPublishedTimeContext, state)
        if err != nil {
            return err
        }
        data["security_publishedTime"] = val
    }
    if self.securityWithdrawnTime.IsSet() {
        val, err := EncodeDateTime(self.securityWithdrawnTime.Get(), path.PushPath("securityWithdrawnTime"), securityVulnerabilitySecurityWithdrawnTimeContext, state)
        if err != nil {
            return err
        }
        data["security_withdrawnTime"] = val
    }
    return nil
}

// A distinct article or unit related to Software.
type SoftwareSoftwareArtifactObject struct {
    ArtifactObject

    // Provides additional purpose information of the software artifact.
    softwareAdditionalPurpose ListProperty[string]
    // Provides a place for the SPDX data creator to record acknowledgement text for
    // a software Package, File or Snippet.
    softwareAttributionText ListProperty[string]
    // A canonical, unique, immutable identifier of the artifact content, that may be
    // used for verifying its identity and/or integrity.
    softwareContentIdentifier RefListProperty[SoftwareContentIdentifier]
    // Identifies the text of one or more copyright notices for a software Package,
    // File or Snippet, if any.
    softwareCopyrightText Property[string]
    // Provides information about the primary purpose of the software artifact.
    softwarePrimaryPurpose Property[string]
}


type SoftwareSoftwareArtifactObjectType struct {
    SHACLTypeBase
}
var softwareSoftwareArtifactType SoftwareSoftwareArtifactObjectType
var softwareSoftwareArtifactSoftwareAdditionalPurposeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/application": "application",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/archive": "archive",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/bom": "bom",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/configuration": "configuration",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/container": "container",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/data": "data",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/device": "device",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/deviceDriver": "deviceDriver",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/diskImage": "diskImage",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/documentation": "documentation",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/evidence": "evidence",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/executable": "executable",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/file": "file",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/filesystemImage": "filesystemImage",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/firmware": "firmware",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/framework": "framework",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/install": "install",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/library": "library",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/manifest": "manifest",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/model": "model",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/module": "module",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/operatingSystem": "operatingSystem",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/patch": "patch",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/platform": "platform",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/requirement": "requirement",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/source": "source",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/specification": "specification",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/test": "test",}
var softwareSoftwareArtifactSoftwareAttributionTextContext = map[string]string{}
var softwareSoftwareArtifactSoftwareContentIdentifierContext = map[string]string{}
var softwareSoftwareArtifactSoftwareCopyrightTextContext = map[string]string{}
var softwareSoftwareArtifactSoftwarePrimaryPurposeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/application": "application",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/archive": "archive",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/bom": "bom",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/configuration": "configuration",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/container": "container",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/data": "data",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/device": "device",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/deviceDriver": "deviceDriver",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/diskImage": "diskImage",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/documentation": "documentation",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/evidence": "evidence",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/executable": "executable",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/file": "file",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/filesystemImage": "filesystemImage",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/firmware": "firmware",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/framework": "framework",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/install": "install",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/library": "library",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/manifest": "manifest",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/model": "model",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/module": "module",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/operatingSystem": "operatingSystem",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/patch": "patch",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/platform": "platform",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/requirement": "requirement",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/source": "source",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/specification": "specification",
    "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/test": "test",}

func DecodeSoftwareSoftwareArtifact (data any, path Path, context map[string]string) (Ref[SoftwareSoftwareArtifact], error) {
    return DecodeRef[SoftwareSoftwareArtifact](data, path, context, softwareSoftwareArtifactType)
}

func (self SoftwareSoftwareArtifactObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareSoftwareArtifact)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Software/additionalPurpose", "software_additionalPurpose":
        val, err := DecodeList[string](value, path, softwareSoftwareArtifactSoftwareAdditionalPurposeContext, DecodeIRI)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareAdditionalPurpose().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/attributionText", "software_attributionText":
        val, err := DecodeList[string](value, path, softwareSoftwareArtifactSoftwareAttributionTextContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareAttributionText().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/contentIdentifier", "software_contentIdentifier":
        val, err := DecodeList[Ref[SoftwareContentIdentifier]](value, path, softwareSoftwareArtifactSoftwareContentIdentifierContext, DecodeSoftwareContentIdentifier)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareContentIdentifier().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/copyrightText", "software_copyrightText":
        val, err := DecodeString(value, path, softwareSoftwareArtifactSoftwareCopyrightTextContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareCopyrightText().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/primaryPurpose", "software_primaryPurpose":
        val, err := DecodeIRI(value, path, softwareSoftwareArtifactSoftwarePrimaryPurposeContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwarePrimaryPurpose().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareSoftwareArtifactObjectType) Create() SHACLObject {
    return ConstructSoftwareSoftwareArtifactObject(&SoftwareSoftwareArtifactObject{}, self)
}

func ConstructSoftwareSoftwareArtifactObject(o *SoftwareSoftwareArtifactObject, typ SHACLType) *SoftwareSoftwareArtifactObject {
    ConstructArtifactObject(&o.ArtifactObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/application",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/archive",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/bom",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/configuration",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/container",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/data",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/device",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/deviceDriver",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/diskImage",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/documentation",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/evidence",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/executable",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/file",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/filesystemImage",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/firmware",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/framework",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/install",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/library",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/manifest",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/model",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/module",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/operatingSystem",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/other",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/patch",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/platform",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/requirement",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/source",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/specification",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/test",
        }})
        o.softwareAdditionalPurpose = NewListProperty[string]("softwareAdditionalPurpose", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwareAttributionText = NewListProperty[string]("softwareAttributionText", validators)
    }
    {
        validators := []Validator[Ref[SoftwareContentIdentifier]]{}
        o.softwareContentIdentifier = NewRefListProperty[SoftwareContentIdentifier]("softwareContentIdentifier", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwareCopyrightText = NewProperty[string]("softwareCopyrightText", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/application",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/archive",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/bom",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/configuration",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/container",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/data",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/device",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/deviceDriver",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/diskImage",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/documentation",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/evidence",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/executable",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/file",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/filesystemImage",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/firmware",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/framework",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/install",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/library",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/manifest",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/model",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/module",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/operatingSystem",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/other",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/patch",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/platform",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/requirement",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/source",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/specification",
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/test",
        }})
        o.softwarePrimaryPurpose = NewProperty[string]("softwarePrimaryPurpose", validators)
    }
    return o
}

type SoftwareSoftwareArtifact interface {
    Artifact
    SoftwareAdditionalPurpose() ListPropertyInterface[string]
    SoftwareAttributionText() ListPropertyInterface[string]
    SoftwareContentIdentifier() ListPropertyInterface[Ref[SoftwareContentIdentifier]]
    SoftwareCopyrightText() PropertyInterface[string]
    SoftwarePrimaryPurpose() PropertyInterface[string]
}



func (self *SoftwareSoftwareArtifactObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ArtifactObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("softwareAdditionalPurpose")
        if ! self.softwareAdditionalPurpose.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareAttributionText")
        if ! self.softwareAttributionText.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareContentIdentifier")
        if ! self.softwareContentIdentifier.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareCopyrightText")
        if ! self.softwareCopyrightText.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwarePrimaryPurpose")
        if ! self.softwarePrimaryPurpose.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SoftwareSoftwareArtifactObject) Walk(path Path, visit Visit) {
    self.ArtifactObject.Walk(path, visit)
    self.softwareAdditionalPurpose.Walk(path, visit)
    self.softwareAttributionText.Walk(path, visit)
    self.softwareContentIdentifier.Walk(path, visit)
    self.softwareCopyrightText.Walk(path, visit)
    self.softwarePrimaryPurpose.Walk(path, visit)
}


func (self *SoftwareSoftwareArtifactObject) SoftwareAdditionalPurpose() ListPropertyInterface[string] { return &self.softwareAdditionalPurpose }
func (self *SoftwareSoftwareArtifactObject) SoftwareAttributionText() ListPropertyInterface[string] { return &self.softwareAttributionText }
func (self *SoftwareSoftwareArtifactObject) SoftwareContentIdentifier() ListPropertyInterface[Ref[SoftwareContentIdentifier]] { return &self.softwareContentIdentifier }
func (self *SoftwareSoftwareArtifactObject) SoftwareCopyrightText() PropertyInterface[string] { return &self.softwareCopyrightText }
func (self *SoftwareSoftwareArtifactObject) SoftwarePrimaryPurpose() PropertyInterface[string] { return &self.softwarePrimaryPurpose }

func (self *SoftwareSoftwareArtifactObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ArtifactObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.softwareAdditionalPurpose.IsSet() {
        val, err := EncodeList[string](self.softwareAdditionalPurpose.Get(), path.PushPath("softwareAdditionalPurpose"), softwareSoftwareArtifactSoftwareAdditionalPurposeContext, state, EncodeIRI)
        if err != nil {
            return err
        }
        data["software_additionalPurpose"] = val
    }
    if self.softwareAttributionText.IsSet() {
        val, err := EncodeList[string](self.softwareAttributionText.Get(), path.PushPath("softwareAttributionText"), softwareSoftwareArtifactSoftwareAttributionTextContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["software_attributionText"] = val
    }
    if self.softwareContentIdentifier.IsSet() {
        val, err := EncodeList[Ref[SoftwareContentIdentifier]](self.softwareContentIdentifier.Get(), path.PushPath("softwareContentIdentifier"), softwareSoftwareArtifactSoftwareContentIdentifierContext, state, EncodeRef[SoftwareContentIdentifier])
        if err != nil {
            return err
        }
        data["software_contentIdentifier"] = val
    }
    if self.softwareCopyrightText.IsSet() {
        val, err := EncodeString(self.softwareCopyrightText.Get(), path.PushPath("softwareCopyrightText"), softwareSoftwareArtifactSoftwareCopyrightTextContext, state)
        if err != nil {
            return err
        }
        data["software_copyrightText"] = val
    }
    if self.softwarePrimaryPurpose.IsSet() {
        val, err := EncodeIRI(self.softwarePrimaryPurpose.Get(), path.PushPath("softwarePrimaryPurpose"), softwareSoftwareArtifactSoftwarePrimaryPurposeContext, state)
        if err != nil {
            return err
        }
        data["software_primaryPurpose"] = val
    }
    return nil
}

// A container for a grouping of SPDX-3.0 content characterizing details

// (provenence, composition, licensing, etc.) about a product.
type BomObject struct {
    BundleObject

}


type BomObjectType struct {
    SHACLTypeBase
}
var bomType BomObjectType

func DecodeBom (data any, path Path, context map[string]string) (Ref[Bom], error) {
    return DecodeRef[Bom](data, path, context, bomType)
}

func (self BomObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(Bom)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self BomObjectType) Create() SHACLObject {
    return ConstructBomObject(&BomObject{}, self)
}

func ConstructBomObject(o *BomObject, typ SHACLType) *BomObject {
    ConstructBundleObject(&o.BundleObject, typ)
    return o
}

type Bom interface {
    Bundle
}


func MakeBom() Bom {
    return ConstructBomObject(&BomObject{}, bomType)
}

func MakeBomRef() Ref[Bom] {
    o := MakeBom()
    return MakeObjectRef[Bom](o)
}

func (self *BomObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.BundleObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *BomObject) Walk(path Path, visit Visit) {
    self.BundleObject.Walk(path, visit)
}



func (self *BomObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.BundleObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// A license that is not listed on the SPDX License List.
type ExpandedlicensingCustomLicenseObject struct {
    ExpandedlicensingLicenseObject

}


type ExpandedlicensingCustomLicenseObjectType struct {
    SHACLTypeBase
}
var expandedlicensingCustomLicenseType ExpandedlicensingCustomLicenseObjectType

func DecodeExpandedlicensingCustomLicense (data any, path Path, context map[string]string) (Ref[ExpandedlicensingCustomLicense], error) {
    return DecodeRef[ExpandedlicensingCustomLicense](data, path, context, expandedlicensingCustomLicenseType)
}

func (self ExpandedlicensingCustomLicenseObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(ExpandedlicensingCustomLicense)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self ExpandedlicensingCustomLicenseObjectType) Create() SHACLObject {
    return ConstructExpandedlicensingCustomLicenseObject(&ExpandedlicensingCustomLicenseObject{}, self)
}

func ConstructExpandedlicensingCustomLicenseObject(o *ExpandedlicensingCustomLicenseObject, typ SHACLType) *ExpandedlicensingCustomLicenseObject {
    ConstructExpandedlicensingLicenseObject(&o.ExpandedlicensingLicenseObject, typ)
    return o
}

type ExpandedlicensingCustomLicense interface {
    ExpandedlicensingLicense
}


func MakeExpandedlicensingCustomLicense() ExpandedlicensingCustomLicense {
    return ConstructExpandedlicensingCustomLicenseObject(&ExpandedlicensingCustomLicenseObject{}, expandedlicensingCustomLicenseType)
}

func MakeExpandedlicensingCustomLicenseRef() Ref[ExpandedlicensingCustomLicense] {
    o := MakeExpandedlicensingCustomLicense()
    return MakeObjectRef[ExpandedlicensingCustomLicense](o)
}

func (self *ExpandedlicensingCustomLicenseObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.ExpandedlicensingLicenseObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *ExpandedlicensingCustomLicenseObject) Walk(path Path, visit Visit) {
    self.ExpandedlicensingLicenseObject.Walk(path, visit)
}



func (self *ExpandedlicensingCustomLicenseObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.ExpandedlicensingLicenseObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Connects a vulnerability and an element designating the element as a product

// affected by the vulnerability.
type SecurityVexAffectedVulnAssessmentRelationshipObject struct {
    SecurityVexVulnAssessmentRelationshipObject

    // Provides advise on how to mitigate or remediate a vulnerability when a VEX product
    // is affected by it.
    securityActionStatement Property[string]
    // Records the time when a recommended action was communicated in a VEX statement
    // to mitigate a vulnerability.
    securityActionStatementTime Property[time.Time]
}


type SecurityVexAffectedVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVexAffectedVulnAssessmentRelationshipType SecurityVexAffectedVulnAssessmentRelationshipObjectType
var securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementContext = map[string]string{}
var securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementTimeContext = map[string]string{}

func DecodeSecurityVexAffectedVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVexAffectedVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVexAffectedVulnAssessmentRelationship](data, path, context, securityVexAffectedVulnAssessmentRelationshipType)
}

func (self SecurityVexAffectedVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexAffectedVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/actionStatement", "security_actionStatement":
        val, err := DecodeString(value, path, securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityActionStatement().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/actionStatementTime", "security_actionStatementTime":
        val, err := DecodeDateTimeStamp(value, path, securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityActionStatementTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexAffectedVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVexAffectedVulnAssessmentRelationshipObject(&SecurityVexAffectedVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVexAffectedVulnAssessmentRelationshipObject(o *SecurityVexAffectedVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVexAffectedVulnAssessmentRelationshipObject {
    ConstructSecurityVexVulnAssessmentRelationshipObject(&o.SecurityVexVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[string]{}
        o.securityActionStatement = NewProperty[string]("securityActionStatement", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityActionStatementTime = NewProperty[time.Time]("securityActionStatementTime", validators)
    }
    return o
}

type SecurityVexAffectedVulnAssessmentRelationship interface {
    SecurityVexVulnAssessmentRelationship
    SecurityActionStatement() PropertyInterface[string]
    SecurityActionStatementTime() PropertyInterface[time.Time]
}


func MakeSecurityVexAffectedVulnAssessmentRelationship() SecurityVexAffectedVulnAssessmentRelationship {
    return ConstructSecurityVexAffectedVulnAssessmentRelationshipObject(&SecurityVexAffectedVulnAssessmentRelationshipObject{}, securityVexAffectedVulnAssessmentRelationshipType)
}

func MakeSecurityVexAffectedVulnAssessmentRelationshipRef() Ref[SecurityVexAffectedVulnAssessmentRelationship] {
    o := MakeSecurityVexAffectedVulnAssessmentRelationship()
    return MakeObjectRef[SecurityVexAffectedVulnAssessmentRelationship](o)
}

func (self *SecurityVexAffectedVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVexVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityActionStatement")
        if ! self.securityActionStatement.Check(prop_path, handler) {
            valid = false
        }
        if ! self.securityActionStatement.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"securityActionStatement", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityActionStatementTime")
        if ! self.securityActionStatementTime.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SecurityVexAffectedVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVexVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityActionStatement.Walk(path, visit)
    self.securityActionStatementTime.Walk(path, visit)
}


func (self *SecurityVexAffectedVulnAssessmentRelationshipObject) SecurityActionStatement() PropertyInterface[string] { return &self.securityActionStatement }
func (self *SecurityVexAffectedVulnAssessmentRelationshipObject) SecurityActionStatementTime() PropertyInterface[time.Time] { return &self.securityActionStatementTime }

func (self *SecurityVexAffectedVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVexVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityActionStatement.IsSet() {
        val, err := EncodeString(self.securityActionStatement.Get(), path.PushPath("securityActionStatement"), securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementContext, state)
        if err != nil {
            return err
        }
        data["security_actionStatement"] = val
    }
    if self.securityActionStatementTime.IsSet() {
        val, err := EncodeDateTime(self.securityActionStatementTime.Get(), path.PushPath("securityActionStatementTime"), securityVexAffectedVulnAssessmentRelationshipSecurityActionStatementTimeContext, state)
        if err != nil {
            return err
        }
        data["security_actionStatementTime"] = val
    }
    return nil
}

// Links a vulnerability and elements representing products (in the VEX sense) where

// a fix has been applied and are no longer affected.
type SecurityVexFixedVulnAssessmentRelationshipObject struct {
    SecurityVexVulnAssessmentRelationshipObject

}


type SecurityVexFixedVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVexFixedVulnAssessmentRelationshipType SecurityVexFixedVulnAssessmentRelationshipObjectType

func DecodeSecurityVexFixedVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVexFixedVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVexFixedVulnAssessmentRelationship](data, path, context, securityVexFixedVulnAssessmentRelationshipType)
}

func (self SecurityVexFixedVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexFixedVulnAssessmentRelationship)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexFixedVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVexFixedVulnAssessmentRelationshipObject(&SecurityVexFixedVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVexFixedVulnAssessmentRelationshipObject(o *SecurityVexFixedVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVexFixedVulnAssessmentRelationshipObject {
    ConstructSecurityVexVulnAssessmentRelationshipObject(&o.SecurityVexVulnAssessmentRelationshipObject, typ)
    return o
}

type SecurityVexFixedVulnAssessmentRelationship interface {
    SecurityVexVulnAssessmentRelationship
}


func MakeSecurityVexFixedVulnAssessmentRelationship() SecurityVexFixedVulnAssessmentRelationship {
    return ConstructSecurityVexFixedVulnAssessmentRelationshipObject(&SecurityVexFixedVulnAssessmentRelationshipObject{}, securityVexFixedVulnAssessmentRelationshipType)
}

func MakeSecurityVexFixedVulnAssessmentRelationshipRef() Ref[SecurityVexFixedVulnAssessmentRelationship] {
    o := MakeSecurityVexFixedVulnAssessmentRelationship()
    return MakeObjectRef[SecurityVexFixedVulnAssessmentRelationship](o)
}

func (self *SecurityVexFixedVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVexVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecurityVexFixedVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVexVulnAssessmentRelationshipObject.Walk(path, visit)
}



func (self *SecurityVexFixedVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVexVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Links a vulnerability and one or more elements designating the latter as products

// not affected by the vulnerability.
type SecurityVexNotAffectedVulnAssessmentRelationshipObject struct {
    SecurityVexVulnAssessmentRelationshipObject

    // Explains why a VEX product is not affected by a vulnerability. It is an
    // alternative in VexNotAffectedVulnAssessmentRelationship to the machine-readable
    // justification label.
    securityImpactStatement Property[string]
    // Timestamp of impact statement.
    securityImpactStatementTime Property[time.Time]
    // Impact justification label to be used when linking a vulnerability to an element
    // representing a VEX product with a VexNotAffectedVulnAssessmentRelationship
    // relationship.
    securityJustificationType Property[string]
}


type SecurityVexNotAffectedVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVexNotAffectedVulnAssessmentRelationshipType SecurityVexNotAffectedVulnAssessmentRelationshipObjectType
var securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementContext = map[string]string{}
var securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementTimeContext = map[string]string{}
var securityVexNotAffectedVulnAssessmentRelationshipSecurityJustificationTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/componentNotPresent": "componentNotPresent",
    "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/inlineMitigationsAlreadyExist": "inlineMitigationsAlreadyExist",
    "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeCannotBeControlledByAdversary": "vulnerableCodeCannotBeControlledByAdversary",
    "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotInExecutePath": "vulnerableCodeNotInExecutePath",
    "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotPresent": "vulnerableCodeNotPresent",}

func DecodeSecurityVexNotAffectedVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVexNotAffectedVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVexNotAffectedVulnAssessmentRelationship](data, path, context, securityVexNotAffectedVulnAssessmentRelationshipType)
}

func (self SecurityVexNotAffectedVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexNotAffectedVulnAssessmentRelationship)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Security/impactStatement", "security_impactStatement":
        val, err := DecodeString(value, path, securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityImpactStatement().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/impactStatementTime", "security_impactStatementTime":
        val, err := DecodeDateTimeStamp(value, path, securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementTimeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityImpactStatementTime().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Security/justificationType", "security_justificationType":
        val, err := DecodeIRI(value, path, securityVexNotAffectedVulnAssessmentRelationshipSecurityJustificationTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.SecurityJustificationType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexNotAffectedVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVexNotAffectedVulnAssessmentRelationshipObject(&SecurityVexNotAffectedVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVexNotAffectedVulnAssessmentRelationshipObject(o *SecurityVexNotAffectedVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVexNotAffectedVulnAssessmentRelationshipObject {
    ConstructSecurityVexVulnAssessmentRelationshipObject(&o.SecurityVexVulnAssessmentRelationshipObject, typ)
    {
        validators := []Validator[string]{}
        o.securityImpactStatement = NewProperty[string]("securityImpactStatement", validators)
    }
    {
        validators := []Validator[time.Time]{}
        validators = append(validators, RegexValidator[time.Time]{`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`})
        o.securityImpactStatementTime = NewProperty[time.Time]("securityImpactStatementTime", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/componentNotPresent",
                "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/inlineMitigationsAlreadyExist",
                "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeCannotBeControlledByAdversary",
                "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotInExecutePath",
                "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType/vulnerableCodeNotPresent",
        }})
        o.securityJustificationType = NewProperty[string]("securityJustificationType", validators)
    }
    return o
}

type SecurityVexNotAffectedVulnAssessmentRelationship interface {
    SecurityVexVulnAssessmentRelationship
    SecurityImpactStatement() PropertyInterface[string]
    SecurityImpactStatementTime() PropertyInterface[time.Time]
    SecurityJustificationType() PropertyInterface[string]
}


func MakeSecurityVexNotAffectedVulnAssessmentRelationship() SecurityVexNotAffectedVulnAssessmentRelationship {
    return ConstructSecurityVexNotAffectedVulnAssessmentRelationshipObject(&SecurityVexNotAffectedVulnAssessmentRelationshipObject{}, securityVexNotAffectedVulnAssessmentRelationshipType)
}

func MakeSecurityVexNotAffectedVulnAssessmentRelationshipRef() Ref[SecurityVexNotAffectedVulnAssessmentRelationship] {
    o := MakeSecurityVexNotAffectedVulnAssessmentRelationship()
    return MakeObjectRef[SecurityVexNotAffectedVulnAssessmentRelationship](o)
}

func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVexVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("securityImpactStatement")
        if ! self.securityImpactStatement.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityImpactStatementTime")
        if ! self.securityImpactStatementTime.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("securityJustificationType")
        if ! self.securityJustificationType.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVexVulnAssessmentRelationshipObject.Walk(path, visit)
    self.securityImpactStatement.Walk(path, visit)
    self.securityImpactStatementTime.Walk(path, visit)
    self.securityJustificationType.Walk(path, visit)
}


func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) SecurityImpactStatement() PropertyInterface[string] { return &self.securityImpactStatement }
func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) SecurityImpactStatementTime() PropertyInterface[time.Time] { return &self.securityImpactStatementTime }
func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) SecurityJustificationType() PropertyInterface[string] { return &self.securityJustificationType }

func (self *SecurityVexNotAffectedVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVexVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.securityImpactStatement.IsSet() {
        val, err := EncodeString(self.securityImpactStatement.Get(), path.PushPath("securityImpactStatement"), securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementContext, state)
        if err != nil {
            return err
        }
        data["security_impactStatement"] = val
    }
    if self.securityImpactStatementTime.IsSet() {
        val, err := EncodeDateTime(self.securityImpactStatementTime.Get(), path.PushPath("securityImpactStatementTime"), securityVexNotAffectedVulnAssessmentRelationshipSecurityImpactStatementTimeContext, state)
        if err != nil {
            return err
        }
        data["security_impactStatementTime"] = val
    }
    if self.securityJustificationType.IsSet() {
        val, err := EncodeIRI(self.securityJustificationType.Get(), path.PushPath("securityJustificationType"), securityVexNotAffectedVulnAssessmentRelationshipSecurityJustificationTypeContext, state)
        if err != nil {
            return err
        }
        data["security_justificationType"] = val
    }
    return nil
}

// Designates elements as products where the impact of a vulnerability is being

// investigated.
type SecurityVexUnderInvestigationVulnAssessmentRelationshipObject struct {
    SecurityVexVulnAssessmentRelationshipObject

}


type SecurityVexUnderInvestigationVulnAssessmentRelationshipObjectType struct {
    SHACLTypeBase
}
var securityVexUnderInvestigationVulnAssessmentRelationshipType SecurityVexUnderInvestigationVulnAssessmentRelationshipObjectType

func DecodeSecurityVexUnderInvestigationVulnAssessmentRelationship (data any, path Path, context map[string]string) (Ref[SecurityVexUnderInvestigationVulnAssessmentRelationship], error) {
    return DecodeRef[SecurityVexUnderInvestigationVulnAssessmentRelationship](data, path, context, securityVexUnderInvestigationVulnAssessmentRelationshipType)
}

func (self SecurityVexUnderInvestigationVulnAssessmentRelationshipObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SecurityVexUnderInvestigationVulnAssessmentRelationship)
    _ = obj
    switch name {
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SecurityVexUnderInvestigationVulnAssessmentRelationshipObjectType) Create() SHACLObject {
    return ConstructSecurityVexUnderInvestigationVulnAssessmentRelationshipObject(&SecurityVexUnderInvestigationVulnAssessmentRelationshipObject{}, self)
}

func ConstructSecurityVexUnderInvestigationVulnAssessmentRelationshipObject(o *SecurityVexUnderInvestigationVulnAssessmentRelationshipObject, typ SHACLType) *SecurityVexUnderInvestigationVulnAssessmentRelationshipObject {
    ConstructSecurityVexVulnAssessmentRelationshipObject(&o.SecurityVexVulnAssessmentRelationshipObject, typ)
    return o
}

type SecurityVexUnderInvestigationVulnAssessmentRelationship interface {
    SecurityVexVulnAssessmentRelationship
}


func MakeSecurityVexUnderInvestigationVulnAssessmentRelationship() SecurityVexUnderInvestigationVulnAssessmentRelationship {
    return ConstructSecurityVexUnderInvestigationVulnAssessmentRelationshipObject(&SecurityVexUnderInvestigationVulnAssessmentRelationshipObject{}, securityVexUnderInvestigationVulnAssessmentRelationshipType)
}

func MakeSecurityVexUnderInvestigationVulnAssessmentRelationshipRef() Ref[SecurityVexUnderInvestigationVulnAssessmentRelationship] {
    o := MakeSecurityVexUnderInvestigationVulnAssessmentRelationship()
    return MakeObjectRef[SecurityVexUnderInvestigationVulnAssessmentRelationship](o)
}

func (self *SecurityVexUnderInvestigationVulnAssessmentRelationshipObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SecurityVexVulnAssessmentRelationshipObject.Validate(path, handler) {
        valid = false
    }
    return valid
}

func (self *SecurityVexUnderInvestigationVulnAssessmentRelationshipObject) Walk(path Path, visit Visit) {
    self.SecurityVexVulnAssessmentRelationshipObject.Walk(path, visit)
}



func (self *SecurityVexUnderInvestigationVulnAssessmentRelationshipObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SecurityVexVulnAssessmentRelationshipObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    return nil
}

// Refers to any object that stores content on a computer.
type SoftwareFileObject struct {
    SoftwareSoftwareArtifactObject

    // Provides information about the content type of an Element or a Property.
    contentType Property[string]
    // Describes if a given file is a directory or non-directory kind of file.
    softwareFileKind Property[string]
}


type SoftwareFileObjectType struct {
    SHACLTypeBase
}
var softwareFileType SoftwareFileObjectType
var softwareFileContentTypeContext = map[string]string{}
var softwareFileSoftwareFileKindContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/directory": "directory",
    "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/file": "file",}

func DecodeSoftwareFile (data any, path Path, context map[string]string) (Ref[SoftwareFile], error) {
    return DecodeRef[SoftwareFile](data, path, context, softwareFileType)
}

func (self SoftwareFileObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareFile)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Core/contentType", "contentType":
        val, err := DecodeString(value, path, softwareFileContentTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.ContentType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/fileKind", "software_fileKind":
        val, err := DecodeIRI(value, path, softwareFileSoftwareFileKindContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareFileKind().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareFileObjectType) Create() SHACLObject {
    return ConstructSoftwareFileObject(&SoftwareFileObject{}, self)
}

func ConstructSoftwareFileObject(o *SoftwareFileObject, typ SHACLType) *SoftwareFileObject {
    ConstructSoftwareSoftwareArtifactObject(&o.SoftwareSoftwareArtifactObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators, RegexValidator[string]{`^[^\/]+\/[^\/]+$`})
        o.contentType = NewProperty[string]("contentType", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/directory",
                "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType/file",
        }})
        o.softwareFileKind = NewProperty[string]("softwareFileKind", validators)
    }
    return o
}

type SoftwareFile interface {
    SoftwareSoftwareArtifact
    ContentType() PropertyInterface[string]
    SoftwareFileKind() PropertyInterface[string]
}


func MakeSoftwareFile() SoftwareFile {
    return ConstructSoftwareFileObject(&SoftwareFileObject{}, softwareFileType)
}

func MakeSoftwareFileRef() Ref[SoftwareFile] {
    o := MakeSoftwareFile()
    return MakeObjectRef[SoftwareFile](o)
}

func (self *SoftwareFileObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SoftwareSoftwareArtifactObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("contentType")
        if ! self.contentType.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareFileKind")
        if ! self.softwareFileKind.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SoftwareFileObject) Walk(path Path, visit Visit) {
    self.SoftwareSoftwareArtifactObject.Walk(path, visit)
    self.contentType.Walk(path, visit)
    self.softwareFileKind.Walk(path, visit)
}


func (self *SoftwareFileObject) ContentType() PropertyInterface[string] { return &self.contentType }
func (self *SoftwareFileObject) SoftwareFileKind() PropertyInterface[string] { return &self.softwareFileKind }

func (self *SoftwareFileObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SoftwareSoftwareArtifactObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.contentType.IsSet() {
        val, err := EncodeString(self.contentType.Get(), path.PushPath("contentType"), softwareFileContentTypeContext, state)
        if err != nil {
            return err
        }
        data["contentType"] = val
    }
    if self.softwareFileKind.IsSet() {
        val, err := EncodeIRI(self.softwareFileKind.Get(), path.PushPath("softwareFileKind"), softwareFileSoftwareFileKindContext, state)
        if err != nil {
            return err
        }
        data["software_fileKind"] = val
    }
    return nil
}

// Refers to any unit of content that can be associated with a distribution of

// software.
type SoftwarePackageObject struct {
    SoftwareSoftwareArtifactObject

    // Identifies the download Uniform Resource Identifier for the package at the time
    // that the document was created.
    softwareDownloadLocation Property[string]
    // A place for the SPDX document creator to record a website that serves as the
    // package's home page.
    softwareHomePage Property[string]
    // Provides a place for the SPDX data creator to record the package URL string
    // (in accordance with the Package URL specification) for a software Package.
    softwarePackageUrl Property[string]
    // Identify the version of a package.
    softwarePackageVersion Property[string]
    // Records any relevant background information or additional comments
    // about the origin of the package.
    softwareSourceInfo Property[string]
}


type SoftwarePackageObjectType struct {
    SHACLTypeBase
}
var softwarePackageType SoftwarePackageObjectType
var softwarePackageSoftwareDownloadLocationContext = map[string]string{}
var softwarePackageSoftwareHomePageContext = map[string]string{}
var softwarePackageSoftwarePackageUrlContext = map[string]string{}
var softwarePackageSoftwarePackageVersionContext = map[string]string{}
var softwarePackageSoftwareSourceInfoContext = map[string]string{}

func DecodeSoftwarePackage (data any, path Path, context map[string]string) (Ref[SoftwarePackage], error) {
    return DecodeRef[SoftwarePackage](data, path, context, softwarePackageType)
}

func (self SoftwarePackageObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwarePackage)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Software/downloadLocation", "software_downloadLocation":
        val, err := DecodeString(value, path, softwarePackageSoftwareDownloadLocationContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareDownloadLocation().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/homePage", "software_homePage":
        val, err := DecodeString(value, path, softwarePackageSoftwareHomePageContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareHomePage().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/packageUrl", "software_packageUrl":
        val, err := DecodeString(value, path, softwarePackageSoftwarePackageUrlContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwarePackageUrl().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/packageVersion", "software_packageVersion":
        val, err := DecodeString(value, path, softwarePackageSoftwarePackageVersionContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwarePackageVersion().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/sourceInfo", "software_sourceInfo":
        val, err := DecodeString(value, path, softwarePackageSoftwareSourceInfoContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareSourceInfo().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwarePackageObjectType) Create() SHACLObject {
    return ConstructSoftwarePackageObject(&SoftwarePackageObject{}, self)
}

func ConstructSoftwarePackageObject(o *SoftwarePackageObject, typ SHACLType) *SoftwarePackageObject {
    ConstructSoftwareSoftwareArtifactObject(&o.SoftwareSoftwareArtifactObject, typ)
    {
        validators := []Validator[string]{}
        o.softwareDownloadLocation = NewProperty[string]("softwareDownloadLocation", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwareHomePage = NewProperty[string]("softwareHomePage", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwarePackageUrl = NewProperty[string]("softwarePackageUrl", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwarePackageVersion = NewProperty[string]("softwarePackageVersion", validators)
    }
    {
        validators := []Validator[string]{}
        o.softwareSourceInfo = NewProperty[string]("softwareSourceInfo", validators)
    }
    return o
}

type SoftwarePackage interface {
    SoftwareSoftwareArtifact
    SoftwareDownloadLocation() PropertyInterface[string]
    SoftwareHomePage() PropertyInterface[string]
    SoftwarePackageUrl() PropertyInterface[string]
    SoftwarePackageVersion() PropertyInterface[string]
    SoftwareSourceInfo() PropertyInterface[string]
}


func MakeSoftwarePackage() SoftwarePackage {
    return ConstructSoftwarePackageObject(&SoftwarePackageObject{}, softwarePackageType)
}

func MakeSoftwarePackageRef() Ref[SoftwarePackage] {
    o := MakeSoftwarePackage()
    return MakeObjectRef[SoftwarePackage](o)
}

func (self *SoftwarePackageObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SoftwareSoftwareArtifactObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("softwareDownloadLocation")
        if ! self.softwareDownloadLocation.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareHomePage")
        if ! self.softwareHomePage.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwarePackageUrl")
        if ! self.softwarePackageUrl.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwarePackageVersion")
        if ! self.softwarePackageVersion.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareSourceInfo")
        if ! self.softwareSourceInfo.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SoftwarePackageObject) Walk(path Path, visit Visit) {
    self.SoftwareSoftwareArtifactObject.Walk(path, visit)
    self.softwareDownloadLocation.Walk(path, visit)
    self.softwareHomePage.Walk(path, visit)
    self.softwarePackageUrl.Walk(path, visit)
    self.softwarePackageVersion.Walk(path, visit)
    self.softwareSourceInfo.Walk(path, visit)
}


func (self *SoftwarePackageObject) SoftwareDownloadLocation() PropertyInterface[string] { return &self.softwareDownloadLocation }
func (self *SoftwarePackageObject) SoftwareHomePage() PropertyInterface[string] { return &self.softwareHomePage }
func (self *SoftwarePackageObject) SoftwarePackageUrl() PropertyInterface[string] { return &self.softwarePackageUrl }
func (self *SoftwarePackageObject) SoftwarePackageVersion() PropertyInterface[string] { return &self.softwarePackageVersion }
func (self *SoftwarePackageObject) SoftwareSourceInfo() PropertyInterface[string] { return &self.softwareSourceInfo }

func (self *SoftwarePackageObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SoftwareSoftwareArtifactObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.softwareDownloadLocation.IsSet() {
        val, err := EncodeString(self.softwareDownloadLocation.Get(), path.PushPath("softwareDownloadLocation"), softwarePackageSoftwareDownloadLocationContext, state)
        if err != nil {
            return err
        }
        data["software_downloadLocation"] = val
    }
    if self.softwareHomePage.IsSet() {
        val, err := EncodeString(self.softwareHomePage.Get(), path.PushPath("softwareHomePage"), softwarePackageSoftwareHomePageContext, state)
        if err != nil {
            return err
        }
        data["software_homePage"] = val
    }
    if self.softwarePackageUrl.IsSet() {
        val, err := EncodeString(self.softwarePackageUrl.Get(), path.PushPath("softwarePackageUrl"), softwarePackageSoftwarePackageUrlContext, state)
        if err != nil {
            return err
        }
        data["software_packageUrl"] = val
    }
    if self.softwarePackageVersion.IsSet() {
        val, err := EncodeString(self.softwarePackageVersion.Get(), path.PushPath("softwarePackageVersion"), softwarePackageSoftwarePackageVersionContext, state)
        if err != nil {
            return err
        }
        data["software_packageVersion"] = val
    }
    if self.softwareSourceInfo.IsSet() {
        val, err := EncodeString(self.softwareSourceInfo.Get(), path.PushPath("softwareSourceInfo"), softwarePackageSoftwareSourceInfoContext, state)
        if err != nil {
            return err
        }
        data["software_sourceInfo"] = val
    }
    return nil
}

// A collection of SPDX Elements describing a single package.
type SoftwareSbomObject struct {
    BomObject

    // Provides information about the type of an SBOM.
    softwareSbomType ListProperty[string]
}


type SoftwareSbomObjectType struct {
    SHACLTypeBase
}
var softwareSbomType SoftwareSbomObjectType
var softwareSbomSoftwareSbomTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/analyzed": "analyzed",
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/build": "build",
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/deployed": "deployed",
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/design": "design",
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/runtime": "runtime",
    "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/source": "source",}

func DecodeSoftwareSbom (data any, path Path, context map[string]string) (Ref[SoftwareSbom], error) {
    return DecodeRef[SoftwareSbom](data, path, context, softwareSbomType)
}

func (self SoftwareSbomObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareSbom)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Software/sbomType", "software_sbomType":
        val, err := DecodeList[string](value, path, softwareSbomSoftwareSbomTypeContext, DecodeIRI)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareSbomType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareSbomObjectType) Create() SHACLObject {
    return ConstructSoftwareSbomObject(&SoftwareSbomObject{}, self)
}

func ConstructSoftwareSbomObject(o *SoftwareSbomObject, typ SHACLType) *SoftwareSbomObject {
    ConstructBomObject(&o.BomObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/analyzed",
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/build",
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/deployed",
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/design",
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/runtime",
                "https://spdx.org/rdf/3.0.1/terms/Software/SbomType/source",
        }})
        o.softwareSbomType = NewListProperty[string]("softwareSbomType", validators)
    }
    return o
}

type SoftwareSbom interface {
    Bom
    SoftwareSbomType() ListPropertyInterface[string]
}


func MakeSoftwareSbom() SoftwareSbom {
    return ConstructSoftwareSbomObject(&SoftwareSbomObject{}, softwareSbomType)
}

func MakeSoftwareSbomRef() Ref[SoftwareSbom] {
    o := MakeSoftwareSbom()
    return MakeObjectRef[SoftwareSbom](o)
}

func (self *SoftwareSbomObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.BomObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("softwareSbomType")
        if ! self.softwareSbomType.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *SoftwareSbomObject) Walk(path Path, visit Visit) {
    self.BomObject.Walk(path, visit)
    self.softwareSbomType.Walk(path, visit)
}


func (self *SoftwareSbomObject) SoftwareSbomType() ListPropertyInterface[string] { return &self.softwareSbomType }

func (self *SoftwareSbomObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.BomObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.softwareSbomType.IsSet() {
        val, err := EncodeList[string](self.softwareSbomType.Get(), path.PushPath("softwareSbomType"), softwareSbomSoftwareSbomTypeContext, state, EncodeIRI)
        if err != nil {
            return err
        }
        data["software_sbomType"] = val
    }
    return nil
}

// Describes a certain part of a file.
type SoftwareSnippetObject struct {
    SoftwareSoftwareArtifactObject

    // Defines the byte range in the original host file that the snippet information
    // applies to.
    softwareByteRange RefProperty[PositiveIntegerRange]
    // Defines the line range in the original host file that the snippet information
    // applies to.
    softwareLineRange RefProperty[PositiveIntegerRange]
    // Defines the original host file that the snippet information applies to.
    softwareSnippetFromFile RefProperty[SoftwareFile]
}


type SoftwareSnippetObjectType struct {
    SHACLTypeBase
}
var softwareSnippetType SoftwareSnippetObjectType
var softwareSnippetSoftwareByteRangeContext = map[string]string{}
var softwareSnippetSoftwareLineRangeContext = map[string]string{}
var softwareSnippetSoftwareSnippetFromFileContext = map[string]string{}

func DecodeSoftwareSnippet (data any, path Path, context map[string]string) (Ref[SoftwareSnippet], error) {
    return DecodeRef[SoftwareSnippet](data, path, context, softwareSnippetType)
}

func (self SoftwareSnippetObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(SoftwareSnippet)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Software/byteRange", "software_byteRange":
        val, err := DecodePositiveIntegerRange(value, path, softwareSnippetSoftwareByteRangeContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareByteRange().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/lineRange", "software_lineRange":
        val, err := DecodePositiveIntegerRange(value, path, softwareSnippetSoftwareLineRangeContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareLineRange().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Software/snippetFromFile", "software_snippetFromFile":
        val, err := DecodeSoftwareFile(value, path, softwareSnippetSoftwareSnippetFromFileContext)
        if err != nil {
            return false, err
        }
        err = obj.SoftwareSnippetFromFile().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self SoftwareSnippetObjectType) Create() SHACLObject {
    return ConstructSoftwareSnippetObject(&SoftwareSnippetObject{}, self)
}

func ConstructSoftwareSnippetObject(o *SoftwareSnippetObject, typ SHACLType) *SoftwareSnippetObject {
    ConstructSoftwareSoftwareArtifactObject(&o.SoftwareSoftwareArtifactObject, typ)
    {
        validators := []Validator[Ref[PositiveIntegerRange]]{}
        o.softwareByteRange = NewRefProperty[PositiveIntegerRange]("softwareByteRange", validators)
    }
    {
        validators := []Validator[Ref[PositiveIntegerRange]]{}
        o.softwareLineRange = NewRefProperty[PositiveIntegerRange]("softwareLineRange", validators)
    }
    {
        validators := []Validator[Ref[SoftwareFile]]{}
        o.softwareSnippetFromFile = NewRefProperty[SoftwareFile]("softwareSnippetFromFile", validators)
    }
    return o
}

type SoftwareSnippet interface {
    SoftwareSoftwareArtifact
    SoftwareByteRange() RefPropertyInterface[PositiveIntegerRange]
    SoftwareLineRange() RefPropertyInterface[PositiveIntegerRange]
    SoftwareSnippetFromFile() RefPropertyInterface[SoftwareFile]
}


func MakeSoftwareSnippet() SoftwareSnippet {
    return ConstructSoftwareSnippetObject(&SoftwareSnippetObject{}, softwareSnippetType)
}

func MakeSoftwareSnippetRef() Ref[SoftwareSnippet] {
    o := MakeSoftwareSnippet()
    return MakeObjectRef[SoftwareSnippet](o)
}

func (self *SoftwareSnippetObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SoftwareSoftwareArtifactObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("softwareByteRange")
        if ! self.softwareByteRange.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareLineRange")
        if ! self.softwareLineRange.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("softwareSnippetFromFile")
        if ! self.softwareSnippetFromFile.Check(prop_path, handler) {
            valid = false
        }
        if ! self.softwareSnippetFromFile.IsSet() {
            if handler != nil {
                handler.HandleError(&ValidationError{"softwareSnippetFromFile", "Value is required"}, prop_path)
            }
            valid = false
        }
    }
    return valid
}

func (self *SoftwareSnippetObject) Walk(path Path, visit Visit) {
    self.SoftwareSoftwareArtifactObject.Walk(path, visit)
    self.softwareByteRange.Walk(path, visit)
    self.softwareLineRange.Walk(path, visit)
    self.softwareSnippetFromFile.Walk(path, visit)
}


func (self *SoftwareSnippetObject) SoftwareByteRange() RefPropertyInterface[PositiveIntegerRange] { return &self.softwareByteRange }
func (self *SoftwareSnippetObject) SoftwareLineRange() RefPropertyInterface[PositiveIntegerRange] { return &self.softwareLineRange }
func (self *SoftwareSnippetObject) SoftwareSnippetFromFile() RefPropertyInterface[SoftwareFile] { return &self.softwareSnippetFromFile }

func (self *SoftwareSnippetObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SoftwareSoftwareArtifactObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.softwareByteRange.IsSet() {
        val, err := EncodeRef[PositiveIntegerRange](self.softwareByteRange.Get(), path.PushPath("softwareByteRange"), softwareSnippetSoftwareByteRangeContext, state)
        if err != nil {
            return err
        }
        data["software_byteRange"] = val
    }
    if self.softwareLineRange.IsSet() {
        val, err := EncodeRef[PositiveIntegerRange](self.softwareLineRange.Get(), path.PushPath("softwareLineRange"), softwareSnippetSoftwareLineRangeContext, state)
        if err != nil {
            return err
        }
        data["software_lineRange"] = val
    }
    if self.softwareSnippetFromFile.IsSet() {
        val, err := EncodeRef[SoftwareFile](self.softwareSnippetFromFile.Get(), path.PushPath("softwareSnippetFromFile"), softwareSnippetSoftwareSnippetFromFileContext, state)
        if err != nil {
            return err
        }
        data["software_snippetFromFile"] = val
    }
    return nil
}

// Specifies an AI package and its associated information.
type AiAIPackageObject struct {
    SoftwarePackageObject

    // Indicates whether the system can perform a decision or action without human
    // involvement or guidance.
    aiAutonomyType Property[string]
    // Captures the domain in which the AI package can be used.
    aiDomain ListProperty[string]
    // Indicates the amount of energy consumption incurred by an AI model.
    aiEnergyConsumption RefProperty[AiEnergyConsumption]
    // Records a hyperparameter used to build the AI model contained in the AI
    // package.
    aiHyperparameter RefListProperty[DictionaryEntry]
    // Provides relevant information about the AI software, not including the model
    // description.
    aiInformationAboutApplication Property[string]
    // Describes relevant information about different steps of the training process.
    aiInformationAboutTraining Property[string]
    // Captures a limitation of the AI software.
    aiLimitation Property[string]
    // Records the measurement of prediction quality of the AI model.
    aiMetric RefListProperty[DictionaryEntry]
    // Captures the threshold that was used for computation of a metric described in
    // the metric field.
    aiMetricDecisionThreshold RefListProperty[DictionaryEntry]
    // Describes all the preprocessing steps applied to the training data before the
    // model training.
    aiModelDataPreprocessing ListProperty[string]
    // Describes methods that can be used to explain the results from the AI model.
    aiModelExplainability ListProperty[string]
    // Records the results of general safety risk assessment of the AI system.
    aiSafetyRiskAssessment Property[string]
    // Captures a standard that is being complied with.
    aiStandardCompliance ListProperty[string]
    // Records the type of the model used in the AI software.
    aiTypeOfModel ListProperty[string]
    // Records if sensitive personal information is used during model training or
    // could be used during the inference.
    aiUseSensitivePersonalInformation Property[string]
}


type AiAIPackageObjectType struct {
    SHACLTypeBase
}
var aiAIPackageType AiAIPackageObjectType
var aiAIPackageAiAutonomyTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no": "no",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion": "noAssertion",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes": "yes",}
var aiAIPackageAiDomainContext = map[string]string{}
var aiAIPackageAiEnergyConsumptionContext = map[string]string{}
var aiAIPackageAiHyperparameterContext = map[string]string{}
var aiAIPackageAiInformationAboutApplicationContext = map[string]string{}
var aiAIPackageAiInformationAboutTrainingContext = map[string]string{}
var aiAIPackageAiLimitationContext = map[string]string{}
var aiAIPackageAiMetricContext = map[string]string{}
var aiAIPackageAiMetricDecisionThresholdContext = map[string]string{}
var aiAIPackageAiModelDataPreprocessingContext = map[string]string{}
var aiAIPackageAiModelExplainabilityContext = map[string]string{}
var aiAIPackageAiSafetyRiskAssessmentContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/high": "high",
    "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/low": "low",
    "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/medium": "medium",
    "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/serious": "serious",}
var aiAIPackageAiStandardComplianceContext = map[string]string{}
var aiAIPackageAiTypeOfModelContext = map[string]string{}
var aiAIPackageAiUseSensitivePersonalInformationContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no": "no",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion": "noAssertion",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes": "yes",}

func DecodeAiAIPackage (data any, path Path, context map[string]string) (Ref[AiAIPackage], error) {
    return DecodeRef[AiAIPackage](data, path, context, aiAIPackageType)
}

func (self AiAIPackageObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(AiAIPackage)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/AI/autonomyType", "ai_autonomyType":
        val, err := DecodeIRI(value, path, aiAIPackageAiAutonomyTypeContext)
        if err != nil {
            return false, err
        }
        err = obj.AiAutonomyType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/domain", "ai_domain":
        val, err := DecodeList[string](value, path, aiAIPackageAiDomainContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.AiDomain().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/energyConsumption", "ai_energyConsumption":
        val, err := DecodeAiEnergyConsumption(value, path, aiAIPackageAiEnergyConsumptionContext)
        if err != nil {
            return false, err
        }
        err = obj.AiEnergyConsumption().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/hyperparameter", "ai_hyperparameter":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, aiAIPackageAiHyperparameterContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.AiHyperparameter().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/informationAboutApplication", "ai_informationAboutApplication":
        val, err := DecodeString(value, path, aiAIPackageAiInformationAboutApplicationContext)
        if err != nil {
            return false, err
        }
        err = obj.AiInformationAboutApplication().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/informationAboutTraining", "ai_informationAboutTraining":
        val, err := DecodeString(value, path, aiAIPackageAiInformationAboutTrainingContext)
        if err != nil {
            return false, err
        }
        err = obj.AiInformationAboutTraining().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/limitation", "ai_limitation":
        val, err := DecodeString(value, path, aiAIPackageAiLimitationContext)
        if err != nil {
            return false, err
        }
        err = obj.AiLimitation().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/metric", "ai_metric":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, aiAIPackageAiMetricContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.AiMetric().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/metricDecisionThreshold", "ai_metricDecisionThreshold":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, aiAIPackageAiMetricDecisionThresholdContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.AiMetricDecisionThreshold().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/modelDataPreprocessing", "ai_modelDataPreprocessing":
        val, err := DecodeList[string](value, path, aiAIPackageAiModelDataPreprocessingContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.AiModelDataPreprocessing().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/modelExplainability", "ai_modelExplainability":
        val, err := DecodeList[string](value, path, aiAIPackageAiModelExplainabilityContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.AiModelExplainability().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/safetyRiskAssessment", "ai_safetyRiskAssessment":
        val, err := DecodeIRI(value, path, aiAIPackageAiSafetyRiskAssessmentContext)
        if err != nil {
            return false, err
        }
        err = obj.AiSafetyRiskAssessment().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/standardCompliance", "ai_standardCompliance":
        val, err := DecodeList[string](value, path, aiAIPackageAiStandardComplianceContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.AiStandardCompliance().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/typeOfModel", "ai_typeOfModel":
        val, err := DecodeList[string](value, path, aiAIPackageAiTypeOfModelContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.AiTypeOfModel().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/AI/useSensitivePersonalInformation", "ai_useSensitivePersonalInformation":
        val, err := DecodeIRI(value, path, aiAIPackageAiUseSensitivePersonalInformationContext)
        if err != nil {
            return false, err
        }
        err = obj.AiUseSensitivePersonalInformation().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self AiAIPackageObjectType) Create() SHACLObject {
    return ConstructAiAIPackageObject(&AiAIPackageObject{}, self)
}

func ConstructAiAIPackageObject(o *AiAIPackageObject, typ SHACLType) *AiAIPackageObject {
    ConstructSoftwarePackageObject(&o.SoftwarePackageObject, typ)
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes",
        }})
        o.aiAutonomyType = NewProperty[string]("aiAutonomyType", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiDomain = NewListProperty[string]("aiDomain", validators)
    }
    {
        validators := []Validator[Ref[AiEnergyConsumption]]{}
        o.aiEnergyConsumption = NewRefProperty[AiEnergyConsumption]("aiEnergyConsumption", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.aiHyperparameter = NewRefListProperty[DictionaryEntry]("aiHyperparameter", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiInformationAboutApplication = NewProperty[string]("aiInformationAboutApplication", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiInformationAboutTraining = NewProperty[string]("aiInformationAboutTraining", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiLimitation = NewProperty[string]("aiLimitation", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.aiMetric = NewRefListProperty[DictionaryEntry]("aiMetric", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.aiMetricDecisionThreshold = NewRefListProperty[DictionaryEntry]("aiMetricDecisionThreshold", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiModelDataPreprocessing = NewListProperty[string]("aiModelDataPreprocessing", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiModelExplainability = NewListProperty[string]("aiModelExplainability", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/high",
                "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/low",
                "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/medium",
                "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType/serious",
        }})
        o.aiSafetyRiskAssessment = NewProperty[string]("aiSafetyRiskAssessment", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiStandardCompliance = NewListProperty[string]("aiStandardCompliance", validators)
    }
    {
        validators := []Validator[string]{}
        o.aiTypeOfModel = NewListProperty[string]("aiTypeOfModel", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes",
        }})
        o.aiUseSensitivePersonalInformation = NewProperty[string]("aiUseSensitivePersonalInformation", validators)
    }
    return o
}

type AiAIPackage interface {
    SoftwarePackage
    AiAutonomyType() PropertyInterface[string]
    AiDomain() ListPropertyInterface[string]
    AiEnergyConsumption() RefPropertyInterface[AiEnergyConsumption]
    AiHyperparameter() ListPropertyInterface[Ref[DictionaryEntry]]
    AiInformationAboutApplication() PropertyInterface[string]
    AiInformationAboutTraining() PropertyInterface[string]
    AiLimitation() PropertyInterface[string]
    AiMetric() ListPropertyInterface[Ref[DictionaryEntry]]
    AiMetricDecisionThreshold() ListPropertyInterface[Ref[DictionaryEntry]]
    AiModelDataPreprocessing() ListPropertyInterface[string]
    AiModelExplainability() ListPropertyInterface[string]
    AiSafetyRiskAssessment() PropertyInterface[string]
    AiStandardCompliance() ListPropertyInterface[string]
    AiTypeOfModel() ListPropertyInterface[string]
    AiUseSensitivePersonalInformation() PropertyInterface[string]
}


func MakeAiAIPackage() AiAIPackage {
    return ConstructAiAIPackageObject(&AiAIPackageObject{}, aiAIPackageType)
}

func MakeAiAIPackageRef() Ref[AiAIPackage] {
    o := MakeAiAIPackage()
    return MakeObjectRef[AiAIPackage](o)
}

func (self *AiAIPackageObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SoftwarePackageObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("aiAutonomyType")
        if ! self.aiAutonomyType.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiDomain")
        if ! self.aiDomain.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiEnergyConsumption")
        if ! self.aiEnergyConsumption.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiHyperparameter")
        if ! self.aiHyperparameter.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiInformationAboutApplication")
        if ! self.aiInformationAboutApplication.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiInformationAboutTraining")
        if ! self.aiInformationAboutTraining.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiLimitation")
        if ! self.aiLimitation.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiMetric")
        if ! self.aiMetric.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiMetricDecisionThreshold")
        if ! self.aiMetricDecisionThreshold.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiModelDataPreprocessing")
        if ! self.aiModelDataPreprocessing.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiModelExplainability")
        if ! self.aiModelExplainability.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiSafetyRiskAssessment")
        if ! self.aiSafetyRiskAssessment.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiStandardCompliance")
        if ! self.aiStandardCompliance.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiTypeOfModel")
        if ! self.aiTypeOfModel.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("aiUseSensitivePersonalInformation")
        if ! self.aiUseSensitivePersonalInformation.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *AiAIPackageObject) Walk(path Path, visit Visit) {
    self.SoftwarePackageObject.Walk(path, visit)
    self.aiAutonomyType.Walk(path, visit)
    self.aiDomain.Walk(path, visit)
    self.aiEnergyConsumption.Walk(path, visit)
    self.aiHyperparameter.Walk(path, visit)
    self.aiInformationAboutApplication.Walk(path, visit)
    self.aiInformationAboutTraining.Walk(path, visit)
    self.aiLimitation.Walk(path, visit)
    self.aiMetric.Walk(path, visit)
    self.aiMetricDecisionThreshold.Walk(path, visit)
    self.aiModelDataPreprocessing.Walk(path, visit)
    self.aiModelExplainability.Walk(path, visit)
    self.aiSafetyRiskAssessment.Walk(path, visit)
    self.aiStandardCompliance.Walk(path, visit)
    self.aiTypeOfModel.Walk(path, visit)
    self.aiUseSensitivePersonalInformation.Walk(path, visit)
}


func (self *AiAIPackageObject) AiAutonomyType() PropertyInterface[string] { return &self.aiAutonomyType }
func (self *AiAIPackageObject) AiDomain() ListPropertyInterface[string] { return &self.aiDomain }
func (self *AiAIPackageObject) AiEnergyConsumption() RefPropertyInterface[AiEnergyConsumption] { return &self.aiEnergyConsumption }
func (self *AiAIPackageObject) AiHyperparameter() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.aiHyperparameter }
func (self *AiAIPackageObject) AiInformationAboutApplication() PropertyInterface[string] { return &self.aiInformationAboutApplication }
func (self *AiAIPackageObject) AiInformationAboutTraining() PropertyInterface[string] { return &self.aiInformationAboutTraining }
func (self *AiAIPackageObject) AiLimitation() PropertyInterface[string] { return &self.aiLimitation }
func (self *AiAIPackageObject) AiMetric() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.aiMetric }
func (self *AiAIPackageObject) AiMetricDecisionThreshold() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.aiMetricDecisionThreshold }
func (self *AiAIPackageObject) AiModelDataPreprocessing() ListPropertyInterface[string] { return &self.aiModelDataPreprocessing }
func (self *AiAIPackageObject) AiModelExplainability() ListPropertyInterface[string] { return &self.aiModelExplainability }
func (self *AiAIPackageObject) AiSafetyRiskAssessment() PropertyInterface[string] { return &self.aiSafetyRiskAssessment }
func (self *AiAIPackageObject) AiStandardCompliance() ListPropertyInterface[string] { return &self.aiStandardCompliance }
func (self *AiAIPackageObject) AiTypeOfModel() ListPropertyInterface[string] { return &self.aiTypeOfModel }
func (self *AiAIPackageObject) AiUseSensitivePersonalInformation() PropertyInterface[string] { return &self.aiUseSensitivePersonalInformation }

func (self *AiAIPackageObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SoftwarePackageObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.aiAutonomyType.IsSet() {
        val, err := EncodeIRI(self.aiAutonomyType.Get(), path.PushPath("aiAutonomyType"), aiAIPackageAiAutonomyTypeContext, state)
        if err != nil {
            return err
        }
        data["ai_autonomyType"] = val
    }
    if self.aiDomain.IsSet() {
        val, err := EncodeList[string](self.aiDomain.Get(), path.PushPath("aiDomain"), aiAIPackageAiDomainContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["ai_domain"] = val
    }
    if self.aiEnergyConsumption.IsSet() {
        val, err := EncodeRef[AiEnergyConsumption](self.aiEnergyConsumption.Get(), path.PushPath("aiEnergyConsumption"), aiAIPackageAiEnergyConsumptionContext, state)
        if err != nil {
            return err
        }
        data["ai_energyConsumption"] = val
    }
    if self.aiHyperparameter.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.aiHyperparameter.Get(), path.PushPath("aiHyperparameter"), aiAIPackageAiHyperparameterContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["ai_hyperparameter"] = val
    }
    if self.aiInformationAboutApplication.IsSet() {
        val, err := EncodeString(self.aiInformationAboutApplication.Get(), path.PushPath("aiInformationAboutApplication"), aiAIPackageAiInformationAboutApplicationContext, state)
        if err != nil {
            return err
        }
        data["ai_informationAboutApplication"] = val
    }
    if self.aiInformationAboutTraining.IsSet() {
        val, err := EncodeString(self.aiInformationAboutTraining.Get(), path.PushPath("aiInformationAboutTraining"), aiAIPackageAiInformationAboutTrainingContext, state)
        if err != nil {
            return err
        }
        data["ai_informationAboutTraining"] = val
    }
    if self.aiLimitation.IsSet() {
        val, err := EncodeString(self.aiLimitation.Get(), path.PushPath("aiLimitation"), aiAIPackageAiLimitationContext, state)
        if err != nil {
            return err
        }
        data["ai_limitation"] = val
    }
    if self.aiMetric.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.aiMetric.Get(), path.PushPath("aiMetric"), aiAIPackageAiMetricContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["ai_metric"] = val
    }
    if self.aiMetricDecisionThreshold.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.aiMetricDecisionThreshold.Get(), path.PushPath("aiMetricDecisionThreshold"), aiAIPackageAiMetricDecisionThresholdContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["ai_metricDecisionThreshold"] = val
    }
    if self.aiModelDataPreprocessing.IsSet() {
        val, err := EncodeList[string](self.aiModelDataPreprocessing.Get(), path.PushPath("aiModelDataPreprocessing"), aiAIPackageAiModelDataPreprocessingContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["ai_modelDataPreprocessing"] = val
    }
    if self.aiModelExplainability.IsSet() {
        val, err := EncodeList[string](self.aiModelExplainability.Get(), path.PushPath("aiModelExplainability"), aiAIPackageAiModelExplainabilityContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["ai_modelExplainability"] = val
    }
    if self.aiSafetyRiskAssessment.IsSet() {
        val, err := EncodeIRI(self.aiSafetyRiskAssessment.Get(), path.PushPath("aiSafetyRiskAssessment"), aiAIPackageAiSafetyRiskAssessmentContext, state)
        if err != nil {
            return err
        }
        data["ai_safetyRiskAssessment"] = val
    }
    if self.aiStandardCompliance.IsSet() {
        val, err := EncodeList[string](self.aiStandardCompliance.Get(), path.PushPath("aiStandardCompliance"), aiAIPackageAiStandardComplianceContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["ai_standardCompliance"] = val
    }
    if self.aiTypeOfModel.IsSet() {
        val, err := EncodeList[string](self.aiTypeOfModel.Get(), path.PushPath("aiTypeOfModel"), aiAIPackageAiTypeOfModelContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["ai_typeOfModel"] = val
    }
    if self.aiUseSensitivePersonalInformation.IsSet() {
        val, err := EncodeIRI(self.aiUseSensitivePersonalInformation.Get(), path.PushPath("aiUseSensitivePersonalInformation"), aiAIPackageAiUseSensitivePersonalInformationContext, state)
        if err != nil {
            return err
        }
        data["ai_useSensitivePersonalInformation"] = val
    }
    return nil
}

// Specifies a data package and its associated information.
type DatasetDatasetPackageObject struct {
    SoftwarePackageObject

    // Describes the anonymization methods used.
    datasetAnonymizationMethodUsed ListProperty[string]
    // Describes the confidentiality level of the data points contained in the dataset.
    datasetConfidentialityLevel Property[string]
    // Describes how the dataset was collected.
    datasetDataCollectionProcess Property[string]
    // Describes the preprocessing steps that were applied to the raw data to create the given dataset.
    datasetDataPreprocessing ListProperty[string]
    // The field describes the availability of a dataset.
    datasetDatasetAvailability Property[string]
    // Describes potentially noisy elements of the dataset.
    datasetDatasetNoise Property[string]
    // Captures the size of the dataset.
    datasetDatasetSize Property[int]
    // Describes the type of the given dataset.
    datasetDatasetType ListProperty[string]
    // Describes a mechanism to update the dataset.
    datasetDatasetUpdateMechanism Property[string]
    // Describes if any sensitive personal information is present in the dataset.
    datasetHasSensitivePersonalInformation Property[string]
    // Describes what the given dataset should be used for.
    datasetIntendedUse Property[string]
    // Records the biases that the dataset is known to encompass.
    datasetKnownBias ListProperty[string]
    // Describes a sensor used for collecting the data.
    datasetSensor RefListProperty[DictionaryEntry]
}


type DatasetDatasetPackageObjectType struct {
    SHACLTypeBase
}
var datasetDatasetPackageType DatasetDatasetPackageObjectType
var datasetDatasetPackageDatasetAnonymizationMethodUsedContext = map[string]string{}
var datasetDatasetPackageDatasetConfidentialityLevelContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/amber": "amber",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/clear": "clear",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/green": "green",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/red": "red",}
var datasetDatasetPackageDatasetDataCollectionProcessContext = map[string]string{}
var datasetDatasetPackageDatasetDataPreprocessingContext = map[string]string{}
var datasetDatasetPackageDatasetDatasetAvailabilityContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/clickthrough": "clickthrough",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/directDownload": "directDownload",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/query": "query",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/registration": "registration",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/scrapingScript": "scrapingScript",}
var datasetDatasetPackageDatasetDatasetNoiseContext = map[string]string{}
var datasetDatasetPackageDatasetDatasetSizeContext = map[string]string{}
var datasetDatasetPackageDatasetDatasetTypeContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/audio": "audio",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/categorical": "categorical",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/graph": "graph",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/image": "image",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/noAssertion": "noAssertion",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/numeric": "numeric",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/other": "other",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/sensor": "sensor",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/structured": "structured",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/syntactic": "syntactic",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/text": "text",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timeseries": "timeseries",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timestamp": "timestamp",
    "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/video": "video",}
var datasetDatasetPackageDatasetDatasetUpdateMechanismContext = map[string]string{}
var datasetDatasetPackageDatasetHasSensitivePersonalInformationContext = map[string]string{
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no": "no",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion": "noAssertion",
    "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes": "yes",}
var datasetDatasetPackageDatasetIntendedUseContext = map[string]string{}
var datasetDatasetPackageDatasetKnownBiasContext = map[string]string{}
var datasetDatasetPackageDatasetSensorContext = map[string]string{}

func DecodeDatasetDatasetPackage (data any, path Path, context map[string]string) (Ref[DatasetDatasetPackage], error) {
    return DecodeRef[DatasetDatasetPackage](data, path, context, datasetDatasetPackageType)
}

func (self DatasetDatasetPackageObjectType) DecodeProperty(o SHACLObject, name string, value interface{}, path Path) (bool, error) {
    obj := o.(DatasetDatasetPackage)
    _ = obj
    switch name {
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/anonymizationMethodUsed", "dataset_anonymizationMethodUsed":
        val, err := DecodeList[string](value, path, datasetDatasetPackageDatasetAnonymizationMethodUsedContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.DatasetAnonymizationMethodUsed().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/confidentialityLevel", "dataset_confidentialityLevel":
        val, err := DecodeIRI(value, path, datasetDatasetPackageDatasetConfidentialityLevelContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetConfidentialityLevel().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/dataCollectionProcess", "dataset_dataCollectionProcess":
        val, err := DecodeString(value, path, datasetDatasetPackageDatasetDataCollectionProcessContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDataCollectionProcess().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/dataPreprocessing", "dataset_dataPreprocessing":
        val, err := DecodeList[string](value, path, datasetDatasetPackageDatasetDataPreprocessingContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDataPreprocessing().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/datasetAvailability", "dataset_datasetAvailability":
        val, err := DecodeIRI(value, path, datasetDatasetPackageDatasetDatasetAvailabilityContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDatasetAvailability().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/datasetNoise", "dataset_datasetNoise":
        val, err := DecodeString(value, path, datasetDatasetPackageDatasetDatasetNoiseContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDatasetNoise().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/datasetSize", "dataset_datasetSize":
        val, err := DecodeInteger(value, path, datasetDatasetPackageDatasetDatasetSizeContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDatasetSize().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/datasetType", "dataset_datasetType":
        val, err := DecodeList[string](value, path, datasetDatasetPackageDatasetDatasetTypeContext, DecodeIRI)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDatasetType().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/datasetUpdateMechanism", "dataset_datasetUpdateMechanism":
        val, err := DecodeString(value, path, datasetDatasetPackageDatasetDatasetUpdateMechanismContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetDatasetUpdateMechanism().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/hasSensitivePersonalInformation", "dataset_hasSensitivePersonalInformation":
        val, err := DecodeIRI(value, path, datasetDatasetPackageDatasetHasSensitivePersonalInformationContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetHasSensitivePersonalInformation().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/intendedUse", "dataset_intendedUse":
        val, err := DecodeString(value, path, datasetDatasetPackageDatasetIntendedUseContext)
        if err != nil {
            return false, err
        }
        err = obj.DatasetIntendedUse().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/knownBias", "dataset_knownBias":
        val, err := DecodeList[string](value, path, datasetDatasetPackageDatasetKnownBiasContext, DecodeString)
        if err != nil {
            return false, err
        }
        err = obj.DatasetKnownBias().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    case "https://spdx.org/rdf/3.0.1/terms/Dataset/sensor", "dataset_sensor":
        val, err := DecodeList[Ref[DictionaryEntry]](value, path, datasetDatasetPackageDatasetSensorContext, DecodeDictionaryEntry)
        if err != nil {
            return false, err
        }
        err = obj.DatasetSensor().Set(val)
        if err != nil {
            return false, err
        }
        return true, nil
    default:
        found, err := self.SHACLTypeBase.DecodeProperty(o, name, value, path)
        if err != nil || found {
            return found, err
        }
    }

    return false, nil
}

func (self DatasetDatasetPackageObjectType) Create() SHACLObject {
    return ConstructDatasetDatasetPackageObject(&DatasetDatasetPackageObject{}, self)
}

func ConstructDatasetDatasetPackageObject(o *DatasetDatasetPackageObject, typ SHACLType) *DatasetDatasetPackageObject {
    ConstructSoftwarePackageObject(&o.SoftwarePackageObject, typ)
    {
        validators := []Validator[string]{}
        o.datasetAnonymizationMethodUsed = NewListProperty[string]("datasetAnonymizationMethodUsed", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/amber",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/clear",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/green",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType/red",
        }})
        o.datasetConfidentialityLevel = NewProperty[string]("datasetConfidentialityLevel", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetDataCollectionProcess = NewProperty[string]("datasetDataCollectionProcess", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetDataPreprocessing = NewListProperty[string]("datasetDataPreprocessing", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/clickthrough",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/directDownload",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/query",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/registration",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType/scrapingScript",
        }})
        o.datasetDatasetAvailability = NewProperty[string]("datasetDatasetAvailability", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetDatasetNoise = NewProperty[string]("datasetDatasetNoise", validators)
    }
    {
        validators := []Validator[int]{}
        validators = append(validators, IntegerMinValidator{0})
        o.datasetDatasetSize = NewProperty[int]("datasetDatasetSize", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/audio",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/categorical",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/graph",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/image",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/noAssertion",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/numeric",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/other",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/sensor",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/structured",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/syntactic",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/text",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timeseries",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/timestamp",
                "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType/video",
        }})
        o.datasetDatasetType = NewListProperty[string]("datasetDatasetType", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetDatasetUpdateMechanism = NewProperty[string]("datasetDatasetUpdateMechanism", validators)
    }
    {
        validators := []Validator[string]{}
        validators = append(validators,
            EnumValidator{[]string{
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/no",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/noAssertion",
                "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType/yes",
        }})
        o.datasetHasSensitivePersonalInformation = NewProperty[string]("datasetHasSensitivePersonalInformation", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetIntendedUse = NewProperty[string]("datasetIntendedUse", validators)
    }
    {
        validators := []Validator[string]{}
        o.datasetKnownBias = NewListProperty[string]("datasetKnownBias", validators)
    }
    {
        validators := []Validator[Ref[DictionaryEntry]]{}
        o.datasetSensor = NewRefListProperty[DictionaryEntry]("datasetSensor", validators)
    }
    return o
}

type DatasetDatasetPackage interface {
    SoftwarePackage
    DatasetAnonymizationMethodUsed() ListPropertyInterface[string]
    DatasetConfidentialityLevel() PropertyInterface[string]
    DatasetDataCollectionProcess() PropertyInterface[string]
    DatasetDataPreprocessing() ListPropertyInterface[string]
    DatasetDatasetAvailability() PropertyInterface[string]
    DatasetDatasetNoise() PropertyInterface[string]
    DatasetDatasetSize() PropertyInterface[int]
    DatasetDatasetType() ListPropertyInterface[string]
    DatasetDatasetUpdateMechanism() PropertyInterface[string]
    DatasetHasSensitivePersonalInformation() PropertyInterface[string]
    DatasetIntendedUse() PropertyInterface[string]
    DatasetKnownBias() ListPropertyInterface[string]
    DatasetSensor() ListPropertyInterface[Ref[DictionaryEntry]]
}


func MakeDatasetDatasetPackage() DatasetDatasetPackage {
    return ConstructDatasetDatasetPackageObject(&DatasetDatasetPackageObject{}, datasetDatasetPackageType)
}

func MakeDatasetDatasetPackageRef() Ref[DatasetDatasetPackage] {
    o := MakeDatasetDatasetPackage()
    return MakeObjectRef[DatasetDatasetPackage](o)
}

func (self *DatasetDatasetPackageObject) Validate(path Path, handler ErrorHandler) bool {
    var valid bool = true
    if ! self.SoftwarePackageObject.Validate(path, handler) {
        valid = false
    }
    {
        prop_path := path.PushPath("datasetAnonymizationMethodUsed")
        if ! self.datasetAnonymizationMethodUsed.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetConfidentialityLevel")
        if ! self.datasetConfidentialityLevel.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDataCollectionProcess")
        if ! self.datasetDataCollectionProcess.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDataPreprocessing")
        if ! self.datasetDataPreprocessing.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDatasetAvailability")
        if ! self.datasetDatasetAvailability.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDatasetNoise")
        if ! self.datasetDatasetNoise.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDatasetSize")
        if ! self.datasetDatasetSize.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDatasetType")
        if ! self.datasetDatasetType.Check(prop_path, handler) {
            valid = false
        }
        if len(self.datasetDatasetType.Get()) < 1 {
            if handler != nil {
                handler.HandleError(&ValidationError{
                    "datasetDatasetType",
                    "Too few elements. Minimum of 1 required"},
                    prop_path)
            }
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetDatasetUpdateMechanism")
        if ! self.datasetDatasetUpdateMechanism.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetHasSensitivePersonalInformation")
        if ! self.datasetHasSensitivePersonalInformation.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetIntendedUse")
        if ! self.datasetIntendedUse.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetKnownBias")
        if ! self.datasetKnownBias.Check(prop_path, handler) {
            valid = false
        }
    }
    {
        prop_path := path.PushPath("datasetSensor")
        if ! self.datasetSensor.Check(prop_path, handler) {
            valid = false
        }
    }
    return valid
}

func (self *DatasetDatasetPackageObject) Walk(path Path, visit Visit) {
    self.SoftwarePackageObject.Walk(path, visit)
    self.datasetAnonymizationMethodUsed.Walk(path, visit)
    self.datasetConfidentialityLevel.Walk(path, visit)
    self.datasetDataCollectionProcess.Walk(path, visit)
    self.datasetDataPreprocessing.Walk(path, visit)
    self.datasetDatasetAvailability.Walk(path, visit)
    self.datasetDatasetNoise.Walk(path, visit)
    self.datasetDatasetSize.Walk(path, visit)
    self.datasetDatasetType.Walk(path, visit)
    self.datasetDatasetUpdateMechanism.Walk(path, visit)
    self.datasetHasSensitivePersonalInformation.Walk(path, visit)
    self.datasetIntendedUse.Walk(path, visit)
    self.datasetKnownBias.Walk(path, visit)
    self.datasetSensor.Walk(path, visit)
}


func (self *DatasetDatasetPackageObject) DatasetAnonymizationMethodUsed() ListPropertyInterface[string] { return &self.datasetAnonymizationMethodUsed }
func (self *DatasetDatasetPackageObject) DatasetConfidentialityLevel() PropertyInterface[string] { return &self.datasetConfidentialityLevel }
func (self *DatasetDatasetPackageObject) DatasetDataCollectionProcess() PropertyInterface[string] { return &self.datasetDataCollectionProcess }
func (self *DatasetDatasetPackageObject) DatasetDataPreprocessing() ListPropertyInterface[string] { return &self.datasetDataPreprocessing }
func (self *DatasetDatasetPackageObject) DatasetDatasetAvailability() PropertyInterface[string] { return &self.datasetDatasetAvailability }
func (self *DatasetDatasetPackageObject) DatasetDatasetNoise() PropertyInterface[string] { return &self.datasetDatasetNoise }
func (self *DatasetDatasetPackageObject) DatasetDatasetSize() PropertyInterface[int] { return &self.datasetDatasetSize }
func (self *DatasetDatasetPackageObject) DatasetDatasetType() ListPropertyInterface[string] { return &self.datasetDatasetType }
func (self *DatasetDatasetPackageObject) DatasetDatasetUpdateMechanism() PropertyInterface[string] { return &self.datasetDatasetUpdateMechanism }
func (self *DatasetDatasetPackageObject) DatasetHasSensitivePersonalInformation() PropertyInterface[string] { return &self.datasetHasSensitivePersonalInformation }
func (self *DatasetDatasetPackageObject) DatasetIntendedUse() PropertyInterface[string] { return &self.datasetIntendedUse }
func (self *DatasetDatasetPackageObject) DatasetKnownBias() ListPropertyInterface[string] { return &self.datasetKnownBias }
func (self *DatasetDatasetPackageObject) DatasetSensor() ListPropertyInterface[Ref[DictionaryEntry]] { return &self.datasetSensor }

func (self *DatasetDatasetPackageObject) EncodeProperties(data map[string]interface{}, path Path, state *EncodeState) error {
    if err := self.SoftwarePackageObject.EncodeProperties(data, path, state); err != nil {
        return err
    }
    if self.datasetAnonymizationMethodUsed.IsSet() {
        val, err := EncodeList[string](self.datasetAnonymizationMethodUsed.Get(), path.PushPath("datasetAnonymizationMethodUsed"), datasetDatasetPackageDatasetAnonymizationMethodUsedContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["dataset_anonymizationMethodUsed"] = val
    }
    if self.datasetConfidentialityLevel.IsSet() {
        val, err := EncodeIRI(self.datasetConfidentialityLevel.Get(), path.PushPath("datasetConfidentialityLevel"), datasetDatasetPackageDatasetConfidentialityLevelContext, state)
        if err != nil {
            return err
        }
        data["dataset_confidentialityLevel"] = val
    }
    if self.datasetDataCollectionProcess.IsSet() {
        val, err := EncodeString(self.datasetDataCollectionProcess.Get(), path.PushPath("datasetDataCollectionProcess"), datasetDatasetPackageDatasetDataCollectionProcessContext, state)
        if err != nil {
            return err
        }
        data["dataset_dataCollectionProcess"] = val
    }
    if self.datasetDataPreprocessing.IsSet() {
        val, err := EncodeList[string](self.datasetDataPreprocessing.Get(), path.PushPath("datasetDataPreprocessing"), datasetDatasetPackageDatasetDataPreprocessingContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["dataset_dataPreprocessing"] = val
    }
    if self.datasetDatasetAvailability.IsSet() {
        val, err := EncodeIRI(self.datasetDatasetAvailability.Get(), path.PushPath("datasetDatasetAvailability"), datasetDatasetPackageDatasetDatasetAvailabilityContext, state)
        if err != nil {
            return err
        }
        data["dataset_datasetAvailability"] = val
    }
    if self.datasetDatasetNoise.IsSet() {
        val, err := EncodeString(self.datasetDatasetNoise.Get(), path.PushPath("datasetDatasetNoise"), datasetDatasetPackageDatasetDatasetNoiseContext, state)
        if err != nil {
            return err
        }
        data["dataset_datasetNoise"] = val
    }
    if self.datasetDatasetSize.IsSet() {
        val, err := EncodeInteger(self.datasetDatasetSize.Get(), path.PushPath("datasetDatasetSize"), datasetDatasetPackageDatasetDatasetSizeContext, state)
        if err != nil {
            return err
        }
        data["dataset_datasetSize"] = val
    }
    if self.datasetDatasetType.IsSet() {
        val, err := EncodeList[string](self.datasetDatasetType.Get(), path.PushPath("datasetDatasetType"), datasetDatasetPackageDatasetDatasetTypeContext, state, EncodeIRI)
        if err != nil {
            return err
        }
        data["dataset_datasetType"] = val
    }
    if self.datasetDatasetUpdateMechanism.IsSet() {
        val, err := EncodeString(self.datasetDatasetUpdateMechanism.Get(), path.PushPath("datasetDatasetUpdateMechanism"), datasetDatasetPackageDatasetDatasetUpdateMechanismContext, state)
        if err != nil {
            return err
        }
        data["dataset_datasetUpdateMechanism"] = val
    }
    if self.datasetHasSensitivePersonalInformation.IsSet() {
        val, err := EncodeIRI(self.datasetHasSensitivePersonalInformation.Get(), path.PushPath("datasetHasSensitivePersonalInformation"), datasetDatasetPackageDatasetHasSensitivePersonalInformationContext, state)
        if err != nil {
            return err
        }
        data["dataset_hasSensitivePersonalInformation"] = val
    }
    if self.datasetIntendedUse.IsSet() {
        val, err := EncodeString(self.datasetIntendedUse.Get(), path.PushPath("datasetIntendedUse"), datasetDatasetPackageDatasetIntendedUseContext, state)
        if err != nil {
            return err
        }
        data["dataset_intendedUse"] = val
    }
    if self.datasetKnownBias.IsSet() {
        val, err := EncodeList[string](self.datasetKnownBias.Get(), path.PushPath("datasetKnownBias"), datasetDatasetPackageDatasetKnownBiasContext, state, EncodeString)
        if err != nil {
            return err
        }
        data["dataset_knownBias"] = val
    }
    if self.datasetSensor.IsSet() {
        val, err := EncodeList[Ref[DictionaryEntry]](self.datasetSensor.Get(), path.PushPath("datasetSensor"), datasetDatasetPackageDatasetSensorContext, state, EncodeRef[DictionaryEntry])
        if err != nil {
            return err
        }
        data["dataset_sensor"] = val
    }
    return nil
}


func init() {
    objectTypes = make(map[string] SHACLType)
    aiEnergyConsumptionType = AiEnergyConsumptionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/AI/EnergyConsumption",
            compactTypeIRI: NewOptional[string]("ai_EnergyConsumption"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(aiEnergyConsumptionType)
    aiEnergyConsumptionDescriptionType = AiEnergyConsumptionDescriptionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/AI/EnergyConsumptionDescription",
            compactTypeIRI: NewOptional[string]("ai_EnergyConsumptionDescription"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(aiEnergyConsumptionDescriptionType)
    aiEnergyUnitTypeType = AiEnergyUnitTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/AI/EnergyUnitType",
            compactTypeIRI: NewOptional[string]("ai_EnergyUnitType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(aiEnergyUnitTypeType)
    aiSafetyRiskAssessmentTypeType = AiSafetyRiskAssessmentTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/AI/SafetyRiskAssessmentType",
            compactTypeIRI: NewOptional[string]("ai_SafetyRiskAssessmentType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(aiSafetyRiskAssessmentTypeType)
    annotationTypeType = AnnotationTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/AnnotationType",
            compactTypeIRI: NewOptional[string]("AnnotationType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(annotationTypeType)
    creationInfoType = CreationInfoObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/CreationInfo",
            compactTypeIRI: NewOptional[string]("CreationInfo"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(creationInfoType)
    dictionaryEntryType = DictionaryEntryObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/DictionaryEntry",
            compactTypeIRI: NewOptional[string]("DictionaryEntry"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(dictionaryEntryType)
    elementType = ElementObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            compactTypeIRI: NewOptional[string]("Element"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(elementType)
    elementCollectionType = ElementCollectionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ElementCollection",
            compactTypeIRI: NewOptional[string]("ElementCollection"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(elementCollectionType)
    externalIdentifierType = ExternalIdentifierObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifier",
            compactTypeIRI: NewOptional[string]("ExternalIdentifier"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(externalIdentifierType)
    externalIdentifierTypeType = ExternalIdentifierTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ExternalIdentifierType",
            compactTypeIRI: NewOptional[string]("ExternalIdentifierType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(externalIdentifierTypeType)
    externalMapType = ExternalMapObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ExternalMap",
            compactTypeIRI: NewOptional[string]("ExternalMap"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(externalMapType)
    externalRefType = ExternalRefObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRef",
            compactTypeIRI: NewOptional[string]("ExternalRef"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(externalRefType)
    externalRefTypeType = ExternalRefTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ExternalRefType",
            compactTypeIRI: NewOptional[string]("ExternalRefType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(externalRefTypeType)
    hashAlgorithmType = HashAlgorithmObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/HashAlgorithm",
            compactTypeIRI: NewOptional[string]("HashAlgorithm"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(hashAlgorithmType)
    individualElementType = IndividualElementObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/IndividualElement",
            compactTypeIRI: NewOptional[string]("IndividualElement"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(individualElementType)
    integrityMethodType = IntegrityMethodObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/IntegrityMethod",
            compactTypeIRI: NewOptional[string]("IntegrityMethod"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(integrityMethodType)
    lifecycleScopeTypeType = LifecycleScopeTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopeType",
            compactTypeIRI: NewOptional[string]("LifecycleScopeType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(lifecycleScopeTypeType)
    namespaceMapType = NamespaceMapObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/NamespaceMap",
            compactTypeIRI: NewOptional[string]("NamespaceMap"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(namespaceMapType)
    packageVerificationCodeType = PackageVerificationCodeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/PackageVerificationCode",
            compactTypeIRI: NewOptional[string]("PackageVerificationCode"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/IntegrityMethod",
            },
        },
    }
    RegisterType(packageVerificationCodeType)
    positiveIntegerRangeType = PositiveIntegerRangeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/PositiveIntegerRange",
            compactTypeIRI: NewOptional[string]("PositiveIntegerRange"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(positiveIntegerRangeType)
    presenceTypeType = PresenceTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/PresenceType",
            compactTypeIRI: NewOptional[string]("PresenceType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(presenceTypeType)
    profileIdentifierTypeType = ProfileIdentifierTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/ProfileIdentifierType",
            compactTypeIRI: NewOptional[string]("ProfileIdentifierType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(profileIdentifierTypeType)
    relationshipType = RelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Relationship",
            compactTypeIRI: NewOptional[string]("Relationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(relationshipType)
    relationshipCompletenessType = RelationshipCompletenessObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipCompleteness",
            compactTypeIRI: NewOptional[string]("RelationshipCompleteness"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(relationshipCompletenessType)
    relationshipTypeType = RelationshipTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType",
            compactTypeIRI: NewOptional[string]("RelationshipType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(relationshipTypeType)
    spdxDocumentType = SpdxDocumentObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/SpdxDocument",
            compactTypeIRI: NewOptional[string]("SpdxDocument"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/ElementCollection",
            },
        },
    }
    RegisterType(spdxDocumentType)
    supportTypeType = SupportTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/SupportType",
            compactTypeIRI: NewOptional[string]("SupportType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(supportTypeType)
    toolType = ToolObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Tool",
            compactTypeIRI: NewOptional[string]("Tool"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(toolType)
    datasetConfidentialityLevelTypeType = DatasetConfidentialityLevelTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Dataset/ConfidentialityLevelType",
            compactTypeIRI: NewOptional[string]("dataset_ConfidentialityLevelType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(datasetConfidentialityLevelTypeType)
    datasetDatasetAvailabilityTypeType = DatasetDatasetAvailabilityTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetAvailabilityType",
            compactTypeIRI: NewOptional[string]("dataset_DatasetAvailabilityType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(datasetDatasetAvailabilityTypeType)
    datasetDatasetTypeType = DatasetDatasetTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetType",
            compactTypeIRI: NewOptional[string]("dataset_DatasetType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(datasetDatasetTypeType)
    expandedlicensingLicenseAdditionType = ExpandedlicensingLicenseAdditionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/LicenseAddition",
            compactTypeIRI: NewOptional[string]("expandedlicensing_LicenseAddition"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(expandedlicensingLicenseAdditionType)
    expandedlicensingListedLicenseExceptionType = ExpandedlicensingListedLicenseExceptionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ListedLicenseException",
            compactTypeIRI: NewOptional[string]("expandedlicensing_ListedLicenseException"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/LicenseAddition",
            },
        },
    }
    RegisterType(expandedlicensingListedLicenseExceptionType)
    extensionCdxPropertyEntryType = ExtensionCdxPropertyEntryObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Extension/CdxPropertyEntry",
            compactTypeIRI: NewOptional[string]("extension_CdxPropertyEntry"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(extensionCdxPropertyEntryType)
    extensionExtensionType = ExtensionExtensionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Extension/Extension",
            compactTypeIRI: NewOptional[string]("extension_Extension"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            isExtensible: NewOptional[bool](true),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(extensionExtensionType)
    securityCvssSeverityTypeType = SecurityCvssSeverityTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/CvssSeverityType",
            compactTypeIRI: NewOptional[string]("security_CvssSeverityType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(securityCvssSeverityTypeType)
    securityExploitCatalogTypeType = SecurityExploitCatalogTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogType",
            compactTypeIRI: NewOptional[string]("security_ExploitCatalogType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(securityExploitCatalogTypeType)
    securitySsvcDecisionTypeType = SecuritySsvcDecisionTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/SsvcDecisionType",
            compactTypeIRI: NewOptional[string]("security_SsvcDecisionType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(securitySsvcDecisionTypeType)
    securityVexJustificationTypeType = SecurityVexJustificationTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexJustificationType",
            compactTypeIRI: NewOptional[string]("security_VexJustificationType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(securityVexJustificationTypeType)
    securityVulnAssessmentRelationshipType = SecurityVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VulnAssessmentRelationship"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Relationship",
            },
        },
    }
    RegisterType(securityVulnAssessmentRelationshipType)
    simplelicensingAnyLicenseInfoType = SimplelicensingAnyLicenseInfoObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            compactTypeIRI: NewOptional[string]("simplelicensing_AnyLicenseInfo"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(simplelicensingAnyLicenseInfoType)
    simplelicensingLicenseExpressionType = SimplelicensingLicenseExpressionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/LicenseExpression",
            compactTypeIRI: NewOptional[string]("simplelicensing_LicenseExpression"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(simplelicensingLicenseExpressionType)
    simplelicensingSimpleLicensingTextType = SimplelicensingSimpleLicensingTextObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/SimpleLicensingText",
            compactTypeIRI: NewOptional[string]("simplelicensing_SimpleLicensingText"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(simplelicensingSimpleLicensingTextType)
    softwareContentIdentifierType = SoftwareContentIdentifierObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifier",
            compactTypeIRI: NewOptional[string]("software_ContentIdentifier"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/IntegrityMethod",
            },
        },
    }
    RegisterType(softwareContentIdentifierType)
    softwareContentIdentifierTypeType = SoftwareContentIdentifierTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/ContentIdentifierType",
            compactTypeIRI: NewOptional[string]("software_ContentIdentifierType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(softwareContentIdentifierTypeType)
    softwareFileKindTypeType = SoftwareFileKindTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/FileKindType",
            compactTypeIRI: NewOptional[string]("software_FileKindType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(softwareFileKindTypeType)
    softwareSbomTypeType = SoftwareSbomTypeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/SbomType",
            compactTypeIRI: NewOptional[string]("software_SbomType"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(softwareSbomTypeType)
    softwareSoftwarePurposeType = SoftwareSoftwarePurposeObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose",
            compactTypeIRI: NewOptional[string]("software_SoftwarePurpose"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
            },
        },
    }
    RegisterType(softwareSoftwarePurposeType)
    buildBuildType = BuildBuildObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Build/Build",
            compactTypeIRI: NewOptional[string]("build_Build"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(buildBuildType)
    agentType = AgentObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Agent",
            compactTypeIRI: NewOptional[string]("Agent"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(agentType)
    annotationType = AnnotationObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Annotation",
            compactTypeIRI: NewOptional[string]("Annotation"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(annotationType)
    artifactType = ArtifactObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Artifact",
            compactTypeIRI: NewOptional[string]("Artifact"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Element",
            },
        },
    }
    RegisterType(artifactType)
    bundleType = BundleObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Bundle",
            compactTypeIRI: NewOptional[string]("Bundle"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/ElementCollection",
            },
        },
    }
    RegisterType(bundleType)
    hashType = HashObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Hash",
            compactTypeIRI: NewOptional[string]("Hash"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/IntegrityMethod",
            },
        },
    }
    RegisterType(hashType)
    lifecycleScopedRelationshipType = LifecycleScopedRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/LifecycleScopedRelationship",
            compactTypeIRI: NewOptional[string]("LifecycleScopedRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Relationship",
            },
        },
    }
    RegisterType(lifecycleScopedRelationshipType)
    organizationType = OrganizationObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Organization",
            compactTypeIRI: NewOptional[string]("Organization"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Agent",
            },
        },
    }
    RegisterType(organizationType)
    personType = PersonObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Person",
            compactTypeIRI: NewOptional[string]("Person"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Agent",
            },
        },
    }
    RegisterType(personType)
    softwareAgentType = SoftwareAgentObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/SoftwareAgent",
            compactTypeIRI: NewOptional[string]("SoftwareAgent"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Agent",
            },
        },
    }
    RegisterType(softwareAgentType)
    expandedlicensingConjunctiveLicenseSetType = ExpandedlicensingConjunctiveLicenseSetObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ConjunctiveLicenseSet",
            compactTypeIRI: NewOptional[string]("expandedlicensing_ConjunctiveLicenseSet"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(expandedlicensingConjunctiveLicenseSetType)
    expandedlicensingCustomLicenseAdditionType = ExpandedlicensingCustomLicenseAdditionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/CustomLicenseAddition",
            compactTypeIRI: NewOptional[string]("expandedlicensing_CustomLicenseAddition"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/LicenseAddition",
            },
        },
    }
    RegisterType(expandedlicensingCustomLicenseAdditionType)
    expandedlicensingDisjunctiveLicenseSetType = ExpandedlicensingDisjunctiveLicenseSetObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/DisjunctiveLicenseSet",
            compactTypeIRI: NewOptional[string]("expandedlicensing_DisjunctiveLicenseSet"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(expandedlicensingDisjunctiveLicenseSetType)
    expandedlicensingExtendableLicenseType = ExpandedlicensingExtendableLicenseObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ExtendableLicense",
            compactTypeIRI: NewOptional[string]("expandedlicensing_ExtendableLicense"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(expandedlicensingExtendableLicenseType)
    expandedlicensingIndividualLicensingInfoType = ExpandedlicensingIndividualLicensingInfoObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/IndividualLicensingInfo",
            compactTypeIRI: NewOptional[string]("expandedlicensing_IndividualLicensingInfo"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(expandedlicensingIndividualLicensingInfoType)
    expandedlicensingLicenseType = ExpandedlicensingLicenseObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/License",
            compactTypeIRI: NewOptional[string]("expandedlicensing_License"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ExtendableLicense",
            },
        },
    }
    RegisterType(expandedlicensingLicenseType)
    expandedlicensingListedLicenseType = ExpandedlicensingListedLicenseObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ListedLicense",
            compactTypeIRI: NewOptional[string]("expandedlicensing_ListedLicense"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/License",
            },
        },
    }
    RegisterType(expandedlicensingListedLicenseType)
    expandedlicensingOrLaterOperatorType = ExpandedlicensingOrLaterOperatorObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/OrLaterOperator",
            compactTypeIRI: NewOptional[string]("expandedlicensing_OrLaterOperator"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/ExtendableLicense",
            },
        },
    }
    RegisterType(expandedlicensingOrLaterOperatorType)
    expandedlicensingWithAdditionOperatorType = ExpandedlicensingWithAdditionOperatorObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/WithAdditionOperator",
            compactTypeIRI: NewOptional[string]("expandedlicensing_WithAdditionOperator"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/SimpleLicensing/AnyLicenseInfo",
            },
        },
    }
    RegisterType(expandedlicensingWithAdditionOperatorType)
    extensionCdxPropertiesExtensionType = ExtensionCdxPropertiesExtensionObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Extension/CdxPropertiesExtension",
            compactTypeIRI: NewOptional[string]("extension_CdxPropertiesExtension"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindBlankNodeOrIRI),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Extension/Extension",
            },
        },
    }
    RegisterType(extensionCdxPropertiesExtensionType)
    securityCvssV2VulnAssessmentRelationshipType = SecurityCvssV2VulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/CvssV2VulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_CvssV2VulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityCvssV2VulnAssessmentRelationshipType)
    securityCvssV3VulnAssessmentRelationshipType = SecurityCvssV3VulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/CvssV3VulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_CvssV3VulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityCvssV3VulnAssessmentRelationshipType)
    securityCvssV4VulnAssessmentRelationshipType = SecurityCvssV4VulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/CvssV4VulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_CvssV4VulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityCvssV4VulnAssessmentRelationshipType)
    securityEpssVulnAssessmentRelationshipType = SecurityEpssVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/EpssVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_EpssVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityEpssVulnAssessmentRelationshipType)
    securityExploitCatalogVulnAssessmentRelationshipType = SecurityExploitCatalogVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/ExploitCatalogVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_ExploitCatalogVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityExploitCatalogVulnAssessmentRelationshipType)
    securitySsvcVulnAssessmentRelationshipType = SecuritySsvcVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/SsvcVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_SsvcVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securitySsvcVulnAssessmentRelationshipType)
    securityVexVulnAssessmentRelationshipType = SecurityVexVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VexVulnAssessmentRelationship"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityVexVulnAssessmentRelationshipType)
    securityVulnerabilityType = SecurityVulnerabilityObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/Vulnerability",
            compactTypeIRI: NewOptional[string]("security_Vulnerability"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Artifact",
            },
        },
    }
    RegisterType(securityVulnerabilityType)
    softwareSoftwareArtifactType = SoftwareSoftwareArtifactObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/SoftwareArtifact",
            compactTypeIRI: NewOptional[string]("software_SoftwareArtifact"),
            isAbstract: true,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Artifact",
            },
        },
    }
    RegisterType(softwareSoftwareArtifactType)
    bomType = BomObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Core/Bom",
            compactTypeIRI: NewOptional[string]("Bom"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Bundle",
            },
        },
    }
    RegisterType(bomType)
    expandedlicensingCustomLicenseType = ExpandedlicensingCustomLicenseObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/CustomLicense",
            compactTypeIRI: NewOptional[string]("expandedlicensing_CustomLicense"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/ExpandedLicensing/License",
            },
        },
    }
    RegisterType(expandedlicensingCustomLicenseType)
    securityVexAffectedVulnAssessmentRelationshipType = SecurityVexAffectedVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexAffectedVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VexAffectedVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VexVulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityVexAffectedVulnAssessmentRelationshipType)
    securityVexFixedVulnAssessmentRelationshipType = SecurityVexFixedVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexFixedVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VexFixedVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VexVulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityVexFixedVulnAssessmentRelationshipType)
    securityVexNotAffectedVulnAssessmentRelationshipType = SecurityVexNotAffectedVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexNotAffectedVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VexNotAffectedVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VexVulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityVexNotAffectedVulnAssessmentRelationshipType)
    securityVexUnderInvestigationVulnAssessmentRelationshipType = SecurityVexUnderInvestigationVulnAssessmentRelationshipObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Security/VexUnderInvestigationVulnAssessmentRelationship",
            compactTypeIRI: NewOptional[string]("security_VexUnderInvestigationVulnAssessmentRelationship"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Security/VexVulnAssessmentRelationship",
            },
        },
    }
    RegisterType(securityVexUnderInvestigationVulnAssessmentRelationshipType)
    softwareFileType = SoftwareFileObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/File",
            compactTypeIRI: NewOptional[string]("software_File"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwareArtifact",
            },
        },
    }
    RegisterType(softwareFileType)
    softwarePackageType = SoftwarePackageObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/Package",
            compactTypeIRI: NewOptional[string]("software_Package"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwareArtifact",
            },
        },
    }
    RegisterType(softwarePackageType)
    softwareSbomType = SoftwareSbomObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/Sbom",
            compactTypeIRI: NewOptional[string]("software_Sbom"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Core/Bom",
            },
        },
    }
    RegisterType(softwareSbomType)
    softwareSnippetType = SoftwareSnippetObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Software/Snippet",
            compactTypeIRI: NewOptional[string]("software_Snippet"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Software/SoftwareArtifact",
            },
        },
    }
    RegisterType(softwareSnippetType)
    aiAIPackageType = AiAIPackageObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/AI/AIPackage",
            compactTypeIRI: NewOptional[string]("ai_AIPackage"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Software/Package",
            },
        },
    }
    RegisterType(aiAIPackageType)
    datasetDatasetPackageType = DatasetDatasetPackageObjectType{
        SHACLTypeBase: SHACLTypeBase{
            typeIRI: "https://spdx.org/rdf/3.0.1/terms/Dataset/DatasetPackage",
            compactTypeIRI: NewOptional[string]("dataset_DatasetPackage"),
            isAbstract: false,
            nodeKind: NewOptional[int](NodeKindIRI),
            idAlias: NewOptional[string]("spdxId"),
            parentIRIs: []string{
                "https://spdx.org/rdf/3.0.1/terms/Software/Package",
            },
        },
    }
    RegisterType(datasetDatasetPackageType)
}
