package main


// A context contains anything interesting to the execution

// Right now, it's a global variable

// maybe we want to make it available via the ecosystem (thus during evaluation of forms)
// and passing it to primitives (so that they can use it to access, well, the context)

type Context struct {
	currentModule string
	nextCurrentModule string     // to switch modules, set nextCurrentModule != nil
	ecosystem *Ecosystem
}
	
