package metric

/**
* created by mengqi on 2023/11/13
* 这里使用函数选项模式
 */

// A VectorOpts is a general configuration.
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}
