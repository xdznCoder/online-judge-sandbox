# online-judge-sandbox 

## 在线判题系统沙盒

------

项目主体参考 [github.com/codeIIEST/test-runner](https://github.com/codeIIEST/test-runner)，以及 [github.com/criyle/go-sandbox](https://github.com/criyle/go-sandbox) 与 [github.com/criyle/go-judge](https://github.com/criyle/go-judge) 等 的开源项目，为xdznOJ提供判题服务与代码运行沙箱。

### 需求分析：

判题系统的基本思路是为提交的代码提供运行环境，并记录其运行结果，与预期结果进行比较。除此之外，出于安全性考虑，需要限制提交代码的内容，防止会对系统造成危害的代码被运行。

### 为什么需要**沙箱 (sandbox)** ？：

1. 判题需求：用于限制提交代码的运行程序所占用cpu的时间与内存大小。
2. 安全性考虑：用于对提交代码的内容进行限制，防止其进行网络访问与未授权行为。
3. 分治：每次都为提交代码创建独立的进程，防止其对彼此产生干扰。

------

~~(ps:先画饼，还在研究源码，具体实现还需等待........)~~

资源板：

1. 青岛大学 Online Judge 文档 https://opensource.qduoj.com/#/judger/how_it_works 
2. criyle 个人博客 https://criyle.github.io/
3. Online Judge 是如何解决判题端安全性问题的？ https://www.zhihu.com/question/23067497
