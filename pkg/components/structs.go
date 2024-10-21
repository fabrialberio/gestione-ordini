package components

type SidebarDest struct {
	DestURL     string
	FasIconName string
	Label       string
	Selected    bool
}

type Input struct {
	Label        string
	Name         string
	Type         string
	DefaultValue string
}

type Select struct {
	Label    string
	Name     string
	Selected int
	Options  []SelectOption
}

type SelectOption struct {
	Value int
	Text  string
}

type TableHead struct {
	URL       string
	OrderBy   int
	OrderDesc bool
	Headings  []TableHeading
}

type TableHeading struct {
	Index int
	Name  string
}
