package sync

import "testing"

func benchWriteOnly(b *testing.B, c Counter) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func benchReadOnly(b *testing.B, c Counter) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = c.Get()
		}
	})
}

func benchMixed(b *testing.B, c Counter, readsPerWrite int) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if i%(readsPerWrite+1) == 0 {
				c.Inc()
			} else {
				_ = c.Get()
			}
		}
	})
}

/**
Benchmark<이름>-<NCPU>   				 <횟수>   		   <ns/op>  		   <B/op>   	   <allocs/op>
BenchmarkMCounter_WriteOnly-8         	 9213634	       118.2 ns/op	       0 B/op	       0 allocs/op


<ns/op>: 한 번의 루프 당 평균 실행 시간
<B/op>:  한 번의 루프당 평균 메모리 할당량
<allocs/op>: 한 번의 루프당 평균 메모리 할당 횟수
*/

func BenchmarkMCounter_WriteOnly(b *testing.B)  { benchWriteOnly(b, &MCounter{}) }
func BenchmarkRWCounter_WriteOnly(b *testing.B) { benchWriteOnly(b, &RWCounter{}) }

func BenchmarkMCounter_ReadOnly(b *testing.B)  { benchReadOnly(b, &MCounter{}) }
func BenchmarkRWCounter_ReadOnly(b *testing.B) { benchReadOnly(b, &RWCounter{}) }

// 읽기:쓰기 = 9:1
func BenchmarkMCounter_ReadWrite9to1(b *testing.B)  { benchMixed(b, &MCounter{}, 9) }
func BenchmarkRWCounter_ReadWrite9to1(b *testing.B) { benchMixed(b, &RWCounter{}, 9) }

// 읽기:쓰기 = 1:1
func BenchmarkMCounter_ReadWrite_1to1(b *testing.B)  { benchMixed(b, &MCounter{}, 1) }
func BenchmarkRWCounter_ReadWrite_1to1(b *testing.B) { benchMixed(b, &RWCounter{}, 1) }
