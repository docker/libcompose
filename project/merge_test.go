package project

import "testing"

type NullLookup struct {
}

func (n *NullLookup) Lookup(file, relativeTo string) ([]byte, string, error) {
	return nil, "", nil
}

func TestExtendsInheritImage(t *testing.T) {
	p := NewProject(&Context{
		ConfigLookup: &NullLookup{},
	})

	config, err := mergeProject(p, []byte(`
parent:
  image: foo
child:
  extends:
    service: parent
`))

	if err != nil {
		t.Fatal(err)
	}

	parent := config["parent"]
	child := config["child"]

	if parent.Image != "foo" {
		t.Fatal("Invalid image", parent.Image)
	}

	if child.Build != "" {
		t.Fatal("Invalid build", child.Build)
	}

	if child.Image != "foo" {
		t.Fatal("Invalid image", child.Image)
	}
}

func TestExtendsInheritBuild(t *testing.T) {
	p := NewProject(&Context{
		ConfigLookup: &NullLookup{},
	})

	config, err := mergeProject(p, []byte(`
parent:
  build: .
child:
  extends:
    service: parent
`))

	if err != nil {
		t.Fatal(err)
	}

	parent := config["parent"]
	child := config["child"]

	if parent.Build != "." {
		t.Fatal("Invalid build", parent.Build)
	}

	if child.Build != "." {
		t.Fatal("Invalid build", child.Build)
	}

	if child.Image != "" {
		t.Fatal("Invalid image", child.Image)
	}
}

func TestExtendBuildOverImage(t *testing.T) {
	p := NewProject(&Context{
		ConfigLookup: &NullLookup{},
	})

	config, err := mergeProject(p, []byte(`
parent:
  image: foo
child:
  build: .
  extends:
    service: parent
`))

	if err != nil {
		t.Fatal(err)
	}

	parent := config["parent"]
	child := config["child"]

	if parent.Image != "foo" {
		t.Fatal("Invalid image", parent.Image)
	}

	if child.Build != "." {
		t.Fatal("Invalid build", child.Build)
	}

	if child.Image != "" {
		t.Fatal("Invalid image", child.Image)
	}
}

func TestExtendImageOverBuild(t *testing.T) {
	p := NewProject(&Context{
		ConfigLookup: &NullLookup{},
	})

	config, err := mergeProject(p, []byte(`
parent:
  build: .
child:
  image: foo
  extends:
    service: parent
`))

	if err != nil {
		t.Fatal(err)
	}

	parent := config["parent"]
	child := config["child"]

	if parent.Image != "" {
		t.Fatal("Invalid image", parent.Image)
	}

	if parent.Build != "." {
		t.Fatal("Invalid build", parent.Build)
	}

	if child.Build != "" {
		t.Fatal("Invalid build", child.Build)
	}

	if child.Image != "foo" {
		t.Fatal("Invalid image", child.Image)
	}
}
