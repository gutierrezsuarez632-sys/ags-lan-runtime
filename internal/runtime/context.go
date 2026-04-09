package runtime

type Project struct {
    Name     string
    Language string
    Contexts []*BoundedContext
}

type BoundedContext struct {
    Name         string
    Hexagonal    bool
    Subcontexts  []string
}

type Context struct {
    Project        *Project
    CurrentContext *BoundedContext
}