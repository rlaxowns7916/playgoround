package problem

import "reflect"

func Or(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		ch := make(chan interface{})
		close(ch)
		return ch
	case 1:
		return chs[0]
	}

	out := make(chan interface{})
	go func() {
		defer close(out)

		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		default:
			select {
			case <-chs[0]:
			case <-chs[1]:
			case <-Or(append(chs[2:], out)...):
			}

		}
	}()
	return out
}

func OrWithReflection(chs ...<-chan struct{}) <-chan struct{} {
	out := make(chan struct{})

	if len(chs) == 0 {
		close(out)
		return out
	}

	go func() {
		defer close(out)
		cases := make([]reflect.SelectCase, 0, len(chs))
		for _, ch := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}
		_, _, _ = reflect.Select(cases) // 하나라도 닫히면 리턴
	}()
	return out
}
