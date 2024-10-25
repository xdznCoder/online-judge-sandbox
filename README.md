# online-judge-sandbox 

## 在线判题系统沙盒

------

项目主体参考 [github.com/criyle/go-sandbox](https://github.com/criyle/go-sandbox) 与 [github.com/criyle/go-judge](https://github.com/criyle/go-judge) 的开源项目，为xdznOJ提供判题服务与代码运行沙箱。

### 需求分析：

判题系统的基本思路是为提交的代码提供运行环境，并记录其运行结果，与预期结果进行比较。除此之外，出于安全性考虑，需要限制提交代码的内容，防止会对系统造成危害的代码被运行。

### 为什么需要**沙箱 (sandbox)** ？：

1. 判题需求：用于限制提交代码的运行程序所占用cpu的时间与内存大小。
2. 安全性考虑：用于对提交代码的内容进行限制，防止其进行网络访问与未授权行为。
3. 分治：每次都为提交代码创建独立的进程，防止其对彼此产生干扰。

### 如何实现**沙箱 (sandbox)**？

本项目主要尝试使用`BPF seccomp` + `ptrace` + `setrlimit` 来实现 `sandbox`：

1. **seccomp**

   `seccomp`是`Linux`内核支持的一种安全机制，用于限制进程对系统调用的访问权限。通过`seccomp`，可以限制程序使用某些系统调用，从而减少系统的暴露面。

   `seccomp-BPF`使用`BPF`技术来实现系统调用过滤。`BPF`程序由一组`BPF`指令组成，可以在系统调用执行之前对其进行检查，以决定是否允许执行该系统调用。`seccomp-BPF`提供了两种模式：**白名单模式**和**黑名单模式**。**白名单模式**允许所有系统调用，除非明确指定不允许的系统调用；黑名单模式则禁止所有系统调用，除非明确指定允许的系统调用。

2. **ptrace**

   `ptrace`是一个系统调用，它允许父进程观察和控制其子进程的执行，以及读取和修改其子进程的核心映像和寄存器。由于`ptrace`可以使子进程处于受控状态，常被用于沙箱保护 。`ptrace`可以用于限制子进程可以使用的系统调用。当目标程序每次尝试进行系统调用时，`ptrace`会通知主程序。如果发现为危险的系统调用，主程序可以及时将程序杀死，从而防止恶意行为的发生。

3. **setrlimit**

   `setrlimit`是`Linux`系统中的一个函数，用于设置当前进程某一指定资源的限制。每个进程在操作系统内核中都是一个独立存在的运行实体，它们都有一组资源限制，这些限制决定了进程对于系统资源的申请量。通过`setrlimit`函数，系统管理员或进程本身可以对这些资源限制进行调整，以确保系统的稳定性和安全性。

   **函数原型：**

   ```c
   #include <sys/resource.h>
   int setrlimit(int resource, const struct rlimit *rlim);
   ```

   通过`setrlimit`函数，对进程所需的`cpu`运行时间与所需内存等资源进行限制，以达到判题需求并防止死循环和栈溢出等情况。

项目主要使用 [github.com/elastic/go-seccomp-bpf](github.com/elastic/go-seccomp-bpf) 来实现对`Linux`系统的`seccomp-BPF`调用。

------

~~(ps:先画饼，还在研究源码，具体实现还需等待........)~~
