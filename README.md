# treeprint
a friendly tree print/dump interface for GO<br>
一个友好的树结构体打印实现，可用于调试、检查、学习树数据
### Usage
please refer to the example in treeprint_test.go<br>
使用方式请参考treeprint_test.go中的测试用例
```cmd
=== RUN   TestPrint
    TestPrint: treeprint_test.go:23: 
        23275(h:7)──────────────────────┐
        |                               |
        491(h:5)─┐                      67186(h:6)──────────────┐
                 |                      |                       |
                 11997(h:4)──┐          63557(h:4)───┐          99634(h:5)
                 |           |          |            |          |         
                 2729(h:3)─┐ 12340(h:1) 37570(h:3)─┐ 66002(h:1) 82050(h:4)─┐
                           |            |          |            |          |
                           3448(h:2)─┐  32756(h:1) 57629(h:2)   69080(h:1) 82509(h:3)─┐
                                     |             |                                  |
                                     8533(h:1)     42653(h:1)                         84362(h:2)─┐
                                                                                                 |
                                                                                                 85840(h:1)
                                                                                                           
    TestPrint: treeprint_test.go:40: 
        6592──┬─────┬─────┬─────┐
        |     |     |     |     |
        25046 16144 50057 32476 95663─┬─────────────────────────┬─────┬────┐
                                |     |                         |     |    |
                                42324 27110─┬─────┬─────┬─────┐ 78601 9704 2585
                                      |     |     |     |     |                
                                      12785 96844 96866 51530 77923
                                                                   
--- PASS: TestPrint (0.00s)
PASS

Process finished with exit code 0

```
