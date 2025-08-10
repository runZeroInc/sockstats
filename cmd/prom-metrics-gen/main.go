package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

const outputPath = "pkg/exporter/generated_exporter.go"

// Metric represents a single metric to be exported.
// It is used by the template to generate the exporter code.
// The template is in template.tmpl.
//
// The fields are:
// - Name: the name of the metric in Prometheus
// - FieldName: the name of the field in the TCPInfo struct
// - Help: the help text for the metric
// - Type: the Prometheus type of the metric (Gauge or Counter)
// - IsNullable: whether the field is a nullable type
// - IsBool: whether the field is a nullable boolean type
type Metric struct {
	Name       string
	FieldName  string
	Help       string
	Type       string
	IsNullable bool
	IsBool     bool
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "pkg/linux/tcpinfo.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var metrics []Metric
	ast.Inspect(node, func(n ast.Node) bool {
		s, ok := n.(*ast.StructType)
		if !ok {
			return true
		}

		for _, f := range s.Fields.List {
			if f.Tag == nil {
				continue
			}
			tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
			if tcpiTag, ok := tag.Lookup("tcpi"); ok {
				var metric Metric
				metric.FieldName = f.Names[0].Name
				tagString := tcpiTag
				for tagString != "" {
					i := strings.Index(tagString, "=")
					if i == -1 {
						log.Printf("malformed tag (missing =): %s [%s]", tagString, metric.FieldName)
						break // malformed tag
					}
					key := tagString[:i]
					tagString = tagString[i+1:]

					var value string
					if strings.HasPrefix(tagString, "'") {
						// value is quoted
						tagString = tagString[1:]
						j := strings.Index(tagString, "'")
						if j == -1 {
							log.Printf("malformed tag (missing '): %s [%s]", tagString, metric.FieldName)
							break // malformed tag
						}
						value = tagString[:j]
						tagString = tagString[j+1:]
						if strings.HasPrefix(tagString, ",") {
							tagString = tagString[1:]
						}
					} else {
						// value is not quoted
						j := strings.Index(tagString, ",")
						if j == -1 {
							value = tagString
							tagString = ""
						} else {
							value = tagString[:j]
							tagString = tagString[j+1:]
						}
					}

					switch key {
					case "name":
						metric.Name = value
					case "prom_type":
						switch value {
						case "gauge":
							metric.Type = "Gauge"
						case "counter":
							metric.Type = "Counter"
						}
					case "prom_help":
						metric.Help = value
					}
				}
				if ident, ok := f.Type.(*ast.Ident); ok {
					metric.IsNullable = strings.HasPrefix(ident.Name, "Nullable")
					metric.IsBool = ident.Name == "NullableBool"
				} else if selExpr, ok := f.Type.(*ast.SelectorExpr); ok {
					metric.IsNullable = strings.HasPrefix(selExpr.Sel.Name, "Nullable")
					metric.IsBool = selExpr.Sel.Name == "NullableBool"
				}
				metrics = append(metrics, metric)
			}
		}
		return false
	})

	t, err := template.ParseFiles("cmd/prom-metrics-gen/template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, struct{ Metrics []Metric }{Metrics: metrics}); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated %s\n", outputPath)
}
