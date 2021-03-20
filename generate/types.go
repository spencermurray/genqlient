package generate

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type typeBuilder struct {
	typeName       string
	typeNamePrefix string
	strings.Builder
	*generator
}

func (g *generator) baseTypeForOperation(operation ast.Operation) *ast.Definition {
	switch operation {
	case ast.Query:
		return g.schema.Query
	case ast.Mutation:
		return g.schema.Mutation
	case ast.Subscription:
		return g.schema.Subscription
	default:
		panic(fmt.Sprintf("unexpected operation: %v", operation))
	}
}

func (g *generator) getTypeForOperation(operation *ast.OperationDefinition) (name string, err error) {
	// TODO: configure ResponseName format
	namePrefix := upperFirst(operation.Name)
	name = namePrefix + "Response"

	if def, ok := g.typeMap[name]; ok {
		// TODO: check for and handle conflicts a better way
		return "", fmt.Errorf("%s already defined:\n%s", name, def)
	}

	fields, err := selections(operation.SelectionSet)
	if err != nil {
		return "", err
	}

	return g.addTypeForDefinition(
		namePrefix, name, g.baseTypeForOperation(operation.Operation), fields)
}

var builtinTypes = map[string]string{
	"Int":     "int", // TODO: technically int32 is always enough, use that?
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"ID":      "string", // TODO: named type for IDs?
}

func (g *generator) addTypeForDefinition(namePrefix, nameOverride string, typ *ast.Definition, fields []field) (name string, err error) {
	// If this is a builtin type, just refer to it.
	goName, ok := builtinTypes[typ.Name]
	if ok {
		return goName, nil
	}

	if nameOverride != "" {
		// if we have an explicit name, the passed-in prefix is what we
		// propagate forward
		name = nameOverride
	} else {
		typeGoName := upperFirst(typ.Name)
		if strings.HasSuffix(namePrefix, typeGoName) {
			// If the field and type names are the same, we can avoid the
			// duplication.  (We include the field name in case there are
			// multiple fields with the same type, and the type name because
			// that's the actual name (the rest are really qualifiers); but if
			// they are the same then including it once suffices for both
			// purposes.)
			name = namePrefix
		} else {
			name = namePrefix + typeGoName
		}

		if typ.Kind != ast.Interface && typ.Kind != ast.Union {
			// for interface/union types, we do not add the type name to the
			// name prefix; we want to have QueryFieldType rather than
			// QueryFieldInterfaceType.  Otherwise, the name will also be the
			// prefix for the next type.
			namePrefix = name
		}

	}

	// Otherwise, build the type, put that in the type-map, and return its
	// name.
	builder := &typeBuilder{typeName: name, typeNamePrefix: namePrefix, generator: g}
	fmt.Fprintf(builder, "type %s ", name)
	err = builder.writeTypedef(typ, fields)
	if err != nil {
		return "", err
	}
	g.typeMap[name] = builder.String()
	return name, nil
}

func (g *generator) getTypeForInputType(typ *ast.Type) (string, error) {
	typeName := upperFirst(typ.Name())
	builder := &typeBuilder{typeName: typeName, typeNamePrefix: typeName, generator: g}
	err := builder.writeType("", typ, selectionsForType(g, typ))
	return builder.String(), err
}

type field interface {
	Alias() string
	Type() *ast.Type
	SubFields() ([]field, error)
}

type outputField struct{ field *ast.Field }

func (s outputField) Alias() string {
	if s.field.Alias != "" {
		return s.field.Alias
	}
	// TODO: is this case needed? tests don't seem to get here.
	return s.field.Name
}

func (s outputField) Type() *ast.Type {
	if s.field.Definition == nil {
		return nil
	}
	return s.field.Definition.Type
}

func (s outputField) SubFields() ([]field, error) {
	return selections(s.field.SelectionSet)
}

func selections(selectionSet ast.SelectionSet) ([]field, error) {
	retval := make([]field, len(selectionSet))
	for i, selection := range selectionSet {
		switch selection := selection.(type) {
		case *ast.Field:
			retval[i] = outputField{selection}
		case *ast.FragmentSpread, *ast.InlineFragment:
			return nil, fmt.Errorf("not implemented: %T", selection)
		default:
			return nil, fmt.Errorf("invalid selection type: %v", selection)
		}
	}
	return retval, nil
}

type inputField struct {
	*generator
	field *ast.FieldDefinition
}

func (s inputField) Alias() string   { return s.field.Name }
func (s inputField) Type() *ast.Type { return s.field.Type }

func (s inputField) SubFields() ([]field, error) {
	return selectionsForType(s.generator, s.field.Type), nil
}

func selectionsForType(g *generator, typ *ast.Type) []field {
	def := g.schema.Types[typ.Name()]
	fields := make([]field, len(def.Fields))
	for i, field := range def.Fields {
		fields[i] = inputField{g, field}
	}
	return fields
}

func (builder *typeBuilder) writeField(field field) error {
	jsonName := field.Alias()
	// We need an exportable name for JSON-marshaling.
	goName := upperFirst(jsonName)

	builder.WriteString(goName)
	builder.WriteRune(' ')

	typ := field.Type()
	if typ == nil {
		// Unclear why gqlparser hasn't already rejected this,
		// but empirically it might not.
		return fmt.Errorf("undefined field %v", field.Alias())
	}

	fields, err := field.SubFields()
	if err != nil {
		return err
	}

	err = builder.writeType(
		// Note we don't deduplicate here -- if our prefix is GetUser and the
		// field name is User, we do GetUserUser.  This is important because if
		// you have a field called user on a type called User we need
		// `query q { user { user { id } } }` to generate two types, QUser and
		// QUserUser.
		// Note also this is the alias, not the field-name, because if we have
		// `query q { a: f { b }, c: f { d } }` we need separate types for a
		// and c, even though they are the same type in GraphQL, because they
		// have different fields.
		builder.typeNamePrefix+upperFirst(field.Alias()), typ, fields)
	if err != nil {
		return err
	}

	if builder.schema.Types[typ.Name()].IsAbstractType() {
		// abstract types are handled in our UnmarshalJSON
		builder.WriteString(" `json:\"-\"`")
	} else if jsonName != goName {
		fmt.Fprintf(builder, " `json:\"%s\"`", jsonName)
	}
	builder.WriteRune('\n')
	return nil
}

func (builder *typeBuilder) writeType(namePrefix string, typ *ast.Type, fields []field) error {
	// gqlgen does slightly different things here, but its implementation may
	// be useful to crib from:
	// https://github.com/99designs/gqlgen/blob/master/plugin/modelgen/models.go#L113
	if typ.Elem != nil {
		// Type is a list.
		builder.WriteString("[]")
		typ = typ.Elem
	}
	if !typ.NonNull {
		builder.WriteString("*")
	}

	def := builder.schema.Types[typ.Name()]
	// Writes a typedef elsewhere (if not already defined)
	name, err := builder.addTypeForDefinition(namePrefix, "", def, fields)
	if err != nil {
		return err
	}

	builder.WriteString(name)
	return nil
}

func (builder *typeBuilder) writeTypedef(typedef *ast.Definition, fields []field) error {
	switch typedef.Kind {
	case ast.Object, ast.InputObject:
		builder.WriteString("struct {\n")
		for _, field := range fields {
			err := builder.writeField(field)
			if err != nil {
				return err
			}
		}
		builder.WriteString("}")

		// If any field is abstract, we need an UnmarshalJSON method to handle
		// it.
		return builder.maybeWriteUnmarshal(fields)

	case ast.Interface, ast.Union:
		// First, write the interface type.
		builder.WriteString("interface {\n")
		implementsMethodName := fmt.Sprintf("implementsGraphQLInterface%v", builder.typeName)
		// TODO: Also write GetX() accessor methods for fields of the interface
		builder.WriteString(implementsMethodName)
		builder.WriteString("()\n")
		builder.WriteString("}")

		// Then, write the implementations.
		// TODO(benkraft): Put a doc-comment somewhere with the list.
		for _, impldef := range builder.schema.GetPossibleTypes(typedef) {
			name, err := builder.addTypeForDefinition(builder.typeNamePrefix, "", impldef, fields)
			if err != nil {
				return err
			}

			// HACK HACK HACK
			builder.typeMap[name] += fmt.Sprintf(
				"\nfunc (v %v) %v() {}", name, implementsMethodName)
		}

		return nil

	case ast.Enum:
		// All GraphQL enums have underlying type string (in the Go sense).
		builder.WriteString("string\n")
		builder.WriteString("const (\n")
		for _, val := range typedef.EnumValues {
			// TODO: casing should be configurable
			fmt.Fprintf(builder, "%s %s = \"%s\"\n",
				builder.typeNamePrefix+goConstName(val.Name),
				builder.typeName, val.Name)
		}
		builder.WriteString(")\n")
		return nil
	case ast.Scalar:
		// TODO(benkraft): Handle custom scalars.
		return fmt.Errorf("not implemented: %v", typedef.Kind)
	default:
		return fmt.Errorf("unexpected kind: %v", typedef.Kind)
	}
}
