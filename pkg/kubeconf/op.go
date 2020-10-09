package kubeconf

import "io"

type Op interface {
	Run(stdout, stderr io.Writer) error
}
